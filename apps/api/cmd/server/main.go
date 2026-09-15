package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/nachsyas/arham-porto/apps/api/internal/config"
	delivery "github.com/nachsyas/arham-porto/apps/api/internal/delivery/http"
	"github.com/nachsyas/arham-porto/apps/api/internal/embedding"
	"github.com/nachsyas/arham-porto/apps/api/internal/embedding/cloudflare"
	geminiEmb "github.com/nachsyas/arham-porto/apps/api/internal/embedding/gemini"
	"github.com/nachsyas/arham-porto/apps/api/internal/llm"
	geminiLLM "github.com/nachsyas/arham-porto/apps/api/internal/llm/gemini"
	"github.com/nachsyas/arham-porto/apps/api/internal/llm/groq"
	"github.com/nachsyas/arham-porto/apps/api/internal/repository/jsonfile"
	"github.com/nachsyas/arham-porto/apps/api/internal/repository/postgres"
	"github.com/nachsyas/arham-porto/apps/api/internal/retrieval"
	"github.com/nachsyas/arham-porto/apps/api/internal/usecase"
)

func main() {
	cfg := config.Load()
	log.Printf("[arham-porto-api] Initializing server (env: %s, port: %s, dbMode: %s)",
		cfg.AppEnv, cfg.Port, cfg.DatabaseMode)

	// 1. Resolve canonical data directory
	dataDir, err := cfg.ResolveDataDir()
	if err != nil {
		log.Fatalf("[arham-porto-api] Fatal: %v", err)
	}
	log.Printf("[arham-porto-api] Canonical data directory resolved: %s", dataDir)

	// 2. Load and index canonical portfolio content
	repo, err := jsonfile.LoadRepository(dataDir)
	if err != nil {
		log.Fatalf("[arham-porto-api] Fatal: failed to load canonical portfolio data: %v", err)
	}
	log.Println("[arham-porto-api] Canonical portfolio data loaded and indexed successfully")

	// 3. Initialize PostgreSQL client according to database mode
	pgClient := postgres.NewClient(cfg.DatabaseURL, cfg.DatabaseMode)
	if err := pgClient.Connect(context.Background()); err != nil {
		if cfg.DatabaseMode == config.DatabaseModeRequired {
			log.Fatalf("[arham-porto-api] Fatal: database connection required in mode %s: %v", cfg.DatabaseMode, err)
		}
		log.Printf("[arham-porto-api] Warning: database connection failed in mode %s: %v", cfg.DatabaseMode, err)
	}
	defer pgClient.Close()

	// 4. Initialize UseCases backed by the canonical repository
	profileUC := usecase.NewProfileUseCase(repo)
	projectUC := usecase.NewProjectUseCase(repo)
	skillUC := usecase.NewSkillUseCase(repo)
	evidenceUC := usecase.NewEvidenceUseCase(repo)
	journeyUC := usecase.NewJourneyUseCase(repo)

	// 5. Initialize bounded in-memory RateLimiter (120 req/min per client IP for general API)
	rateLimiter := delivery.NewRateLimiter(120, time.Minute)
	defer rateLimiter.Close()

	// 6. Initialize Ask Arham AI dependencies
	var askUC usecase.AskUseCase
	var aiLimiter *delivery.RateLimiter

	if cfg.AIMode == "remote" {
		if cfg.AIProvider == "fake" {
			log.Fatalf("[arham-porto-api] Fatal: AI_PROVIDER=fake is strictly prohibited on production server")
		}

		var llmProvider llm.LLMProvider
		switch cfg.AIProvider {
		case "groq":
			if cfg.GroqAPIKey == "" {
				log.Fatalf("[arham-porto-api] Fatal: GROQ_API_KEY is required when AI_MODE=remote and AI_PROVIDER=groq")
			}
			client, err := groq.NewClient(cfg.GroqAPIKey, cfg.AIModel, groq.WithReasoningEffort(cfg.AIThinkingLevel))
			if err != nil {
				log.Fatalf("[arham-porto-api] Fatal: failed to initialize Groq LLM provider: %v", err)
			}
			llmProvider = client
		case "gemini":
			if cfg.GeminiAPIKey == "" {
				log.Fatalf("[arham-porto-api] Fatal: GEMINI_API_KEY is required when AI_MODE=remote and AI_PROVIDER=gemini")
			}
			client, err := geminiLLM.NewClient(cfg.GeminiAPIKey, cfg.AIModel, geminiLLM.WithThinkingLevel(cfg.AIThinkingLevel))
			if err != nil {
				log.Fatalf("[arham-porto-api] Fatal: failed to initialize Gemini LLM provider: %v", err)
			}
			llmProvider = client
		default:
			log.Fatalf("[arham-porto-api] Fatal: unsupported AI_PROVIDER %q", cfg.AIProvider)
		}

		// Initialize embedding provider and retrieval service if postgres is available
		if pgClient.Pool() != nil && cfg.EmbeddingMode == "enabled" {
			var embClient embedding.Provider
			var embErr error

			switch cfg.EmbeddingProvider {
			case "cloudflare":
				if cfg.CloudflareAIToken == "" || cfg.CloudflareAccountID == "" {
					log.Printf("[arham-porto-api] Warning: CLOUDFLARE_AI_TOKEN and CLOUDFLARE_ACCOUNT_ID are required when EMBEDDING_PROVIDER=cloudflare")
				} else {
					embClient, embErr = cloudflare.NewProvider(cfg.CloudflareAccountID, cfg.CloudflareAIToken, cfg.EmbeddingModel, cfg.EmbeddingDimensions)
				}
			case "gemini":
				if cfg.GeminiAPIKey == "" {
					log.Printf("[arham-porto-api] Warning: GEMINI_API_KEY is required when EMBEDDING_PROVIDER=gemini")
				} else {
					embClient, embErr = geminiEmb.NewProvider(cfg.GeminiAPIKey, cfg.EmbeddingModel, cfg.EmbeddingDimensions)
				}
			case "fake":
				log.Fatalf("[arham-porto-api] Fatal: EMBEDDING_PROVIDER=fake is strictly prohibited on production server")
			default:
				log.Printf("[arham-porto-api] Warning: unsupported EMBEDDING_PROVIDER %q", cfg.EmbeddingProvider)
			}

			if embErr == nil && embClient != nil && llmProvider != nil {
				knowledgeRepo := postgres.NewKnowledgeRepository(pgClient)
				retrievalSvc := retrieval.NewService(knowledgeRepo, embClient)
				askUC = usecase.NewAskUseCase(retrievalSvc, llmProvider, cfg.AIMaxEvidenceChars)
				log.Printf("[arham-porto-api] Ask Arham AI initialized with llm_provider=%s model=%s emb_provider=%s emb_model=%s", cfg.AIProvider, cfg.AIModel, cfg.EmbeddingProvider, cfg.EmbeddingModel)
			} else if embErr != nil {
				log.Printf("[arham-porto-api] Warning: failed to initialize embedding client: %v", embErr)
			}
		} else {
			log.Printf("[arham-porto-api] Warning: database or embeddings not ready; Ask Arham AI will report unavailable")
		}
	} else {
		log.Println("[arham-porto-api] Ask Arham AI is disabled (AI_MODE=disabled)")
	}

	aiLimiter = delivery.NewRateLimiter(cfg.AIRateLimitPerMinute, time.Minute)
	defer aiLimiter.Close()

	// 7. Initialize delivery layer
	handler := delivery.NewHandler(profileUC, projectUC, skillUC, evidenceUC, journeyUC, pgClient)
	handler.EnableAI(askUC, cfg.AIMode, aiLimiter, cfg.AIMaxConcurrentRequests, cfg.AIRequestTimeoutSeconds)
	router := delivery.NewRouter(handler, cfg.AllowedOrigins, rateLimiter, cfg.TrustProxyMode)

	// 8. Configure hardened HTTP server
	writeTimeout := 35 * time.Second
	if cfg.AIRequestTimeoutSeconds > 30 {
		writeTimeout = time.Duration(cfg.AIRequestTimeoutSeconds+5) * time.Second
	}

	listenAddr := "0.0.0.0:" + cfg.Port
	if strings.Contains(cfg.Port, ":") {
		listenAddr = cfg.Port
	}

	srv := &http.Server{
		Addr:              listenAddr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MB
	}

	// Graceful shutdown channel
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("[arham-porto-api] Server listening on %s (trustProxyMode: %s)", listenAddr, cfg.TrustProxyMode)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[arham-porto-api] ListenAndServe error: %v", err)
		}
	}()

	<-stop
	log.Println("[arham-porto-api] Shutdown signal received. Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("[arham-porto-api] Server forced shutdown: %v", err)
	}

	log.Println("[arham-porto-api] Server stopped cleanly")
}
