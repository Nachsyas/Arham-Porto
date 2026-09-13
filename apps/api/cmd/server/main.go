package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nachsyas/arham-porto/apps/api/internal/config"
	delivery "github.com/nachsyas/arham-porto/apps/api/internal/delivery/http"
	"github.com/nachsyas/arham-porto/apps/api/internal/repository/jsonfile"
	"github.com/nachsyas/arham-porto/apps/api/internal/repository/postgres"
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
			log.Fatalf("[arham-porto-api] Fatal: PostgreSQL is required but failed to connect: %v", err)
		}
	}
	defer pgClient.Close()

	// 4. Initialize UseCases backed by the canonical repository
	profileUC := usecase.NewProfileUseCase(repo)
	projectUC := usecase.NewProjectUseCase(repo)
	skillUC := usecase.NewSkillUseCase(repo)
	evidenceUC := usecase.NewEvidenceUseCase(repo)
	journeyUC := usecase.NewJourneyUseCase(repo)

	// 5. Initialize bounded in-memory RateLimiter (120 req/min per client IP)
	rateLimiter := delivery.NewRateLimiter(120, time.Minute)
	defer rateLimiter.Close()

	// 6. Initialize delivery layer
	handler := delivery.NewHandler(profileUC, projectUC, skillUC, evidenceUC, journeyUC, pgClient)
	router := delivery.NewRouter(handler, cfg.AllowedOrigins, rateLimiter)

	// 7. Configure hardened HTTP server
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MB
	}

	// Graceful shutdown channel
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("[arham-porto-api] Server listening on :%s", cfg.Port)
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
