package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// DatabaseMode defines how the API interacts with PostgreSQL.
const (
	DatabaseModeDisabled = "disabled"
	DatabaseModeOptional = "optional"
	DatabaseModeRequired = "required"
)

// Config holds the application configuration.
type Config struct {
	AppEnv           string
	Port             string
	AllowedOrigins   []string
	DatabaseURL         string
	DatabaseMode        string // disabled, optional, required
	PortfolioDataDir    string
	GitHubToken         string
	EmbeddingMode       string // disabled, enabled
	EmbeddingProvider   string // gemini, fake, disabled
	EmbeddingModel      string
	EmbeddingDimensions int
	GeminiAPIKey        string

	// Production Edge Proxy Mode
	TrustProxyMode string // direct, cloudrun, railway (default: direct)

	// Phase 6 Ask Arham AI Configuration
	AIMode                  string // disabled, remote
	AIProvider              string // gemini (rejects fake in production)
	AIModel                 string // default: gemini-3.8-flash
	AIThinkingLevel         string // low, medium, high (default: low)
	AIRequestTimeoutSeconds int    // default: 30
	AIMaxConcurrentRequests int    // default: 4
	AIRateLimitPerMinute    int    // default: 5
	AIMaxEvidenceChars      int    // default: 24000
}

// Load loads environment variables with safe development fallbacks.
func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("API_PORT")
	}
	if port == "" {
		port = "8080"
	}

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}

	allowedOriginsStr := os.Getenv("ALLOWED_ORIGINS")
	if allowedOriginsStr == "" {
		allowedOriginsStr = os.Getenv("WEB_ORIGIN")
	}
	if allowedOriginsStr == "" {
		allowedOriginsStr = "http://localhost:3000"
	}

	var allowedOrigins []string
	for _, origin := range strings.Split(allowedOriginsStr, ",") {
		trimmed := strings.TrimSpace(origin)
		if trimmed != "" {
			allowedOrigins = append(allowedOrigins, trimmed)
		}
	}
	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{"http://localhost:3000"}
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://arham_user:arham_password@localhost:5432/arham_porto?sslmode=disable"
	}

	dbMode := strings.ToLower(os.Getenv("DATABASE_MODE"))
	if dbMode == "" {
		dbMode = DatabaseModeOptional
	}
	if dbMode != DatabaseModeDisabled && dbMode != DatabaseModeOptional && dbMode != DatabaseModeRequired {
		dbMode = DatabaseModeOptional
	}

	dataDir := os.Getenv("PORTFOLIO_DATA_DIR")
	githubToken := os.Getenv("GITHUB_TOKEN")

	embeddingMode := strings.ToLower(strings.TrimSpace(os.Getenv("EMBEDDING_MODE")))
	if embeddingMode == "" {
		embeddingMode = "disabled"
	}

	embeddingProvider := strings.ToLower(strings.TrimSpace(os.Getenv("EMBEDDING_PROVIDER")))
	embeddingModel := strings.TrimSpace(os.Getenv("EMBEDDING_MODEL"))
	if embeddingModel == "" {
		embeddingModel = "gemini-embedding-2"
	}

	embeddingDimensions := 768
	if dimsStr := os.Getenv("EMBEDDING_DIMENSIONS"); dimsStr != "" {
		if d, err := strconv.Atoi(dimsStr); err == nil && d > 0 {
			embeddingDimensions = d
		}
	}

	geminiAPIKey := os.Getenv("GEMINI_API_KEY")

	// Phase 6 Ask Arham AI Configuration
	aiMode := strings.ToLower(strings.TrimSpace(os.Getenv("AI_MODE")))
	if aiMode == "" {
		aiMode = "disabled"
	}

	aiProvider := strings.ToLower(strings.TrimSpace(os.Getenv("AI_PROVIDER")))
	aiModel := strings.TrimSpace(os.Getenv("AI_MODEL"))
	if aiModel == "" {
		aiModel = "gemini-3.8-flash"
	}

	aiThinkingLevel := strings.ToLower(strings.TrimSpace(os.Getenv("AI_THINKING_LEVEL")))
	if aiThinkingLevel == "" {
		aiThinkingLevel = "low"
	}

	aiTimeoutSeconds := 30
	if timeoutStr := os.Getenv("AI_REQUEST_TIMEOUT_SECONDS"); timeoutStr != "" {
		if t, err := strconv.Atoi(timeoutStr); err == nil && t > 0 {
			aiTimeoutSeconds = t
		}
	}

	aiMaxConcurrent := 4
	if concurrentStr := os.Getenv("AI_MAX_CONCURRENT_REQUESTS"); concurrentStr != "" {
		if c, err := strconv.Atoi(concurrentStr); err == nil && c > 0 {
			aiMaxConcurrent = c
		}
	}

	aiRateLimit := 5
	if rateLimitStr := os.Getenv("AI_RATE_LIMIT_PER_MINUTE"); rateLimitStr != "" {
		if r, err := strconv.Atoi(rateLimitStr); err == nil && r > 0 {
			aiRateLimit = r
		}
	}

	aiMaxEvidenceChars := 24000
	if maxCharsStr := os.Getenv("AI_MAX_EVIDENCE_CHARS"); maxCharsStr != "" {
		if m, err := strconv.Atoi(maxCharsStr); err == nil && m > 0 {
			aiMaxEvidenceChars = m
		}
	}

	trustProxyMode := strings.ToLower(strings.TrimSpace(os.Getenv("TRUST_PROXY_MODE")))
	if trustProxyMode != "cloudrun" && trustProxyMode != "railway" {
		trustProxyMode = "direct"
	}

	return &Config{
		AppEnv:                  appEnv,
		Port:                    port,
		AllowedOrigins:          allowedOrigins,
		DatabaseURL:             dbURL,
		DatabaseMode:            dbMode,
		PortfolioDataDir:        dataDir,
		GitHubToken:             githubToken,
		EmbeddingMode:           embeddingMode,
		EmbeddingProvider:       embeddingProvider,
		EmbeddingModel:          embeddingModel,
		EmbeddingDimensions:     embeddingDimensions,
		GeminiAPIKey:            geminiAPIKey,
		TrustProxyMode:          trustProxyMode,
		AIMode:                  aiMode,
		AIProvider:              aiProvider,
		AIModel:                 aiModel,
		AIThinkingLevel:         aiThinkingLevel,
		AIRequestTimeoutSeconds: aiTimeoutSeconds,
		AIMaxConcurrentRequests: aiMaxConcurrent,
		AIRateLimitPerMinute:    aiRateLimit,
		AIMaxEvidenceChars:      aiMaxEvidenceChars,
	}
}

// ResolveDataDir resolves and validates the path to canonical data/ directory deterministically.
// If PORTFOLIO_DATA_DIR is set, it uses exactly that directory.
// If unset, it uses documented deterministic defaults based on server working directory (../../data or data).
// It does NOT walk upwards searching for data directories.
func (c *Config) ResolveDataDir() (string, error) {
	var targetDir string

	if c.PortfolioDataDir != "" {
		targetDir = c.PortfolioDataDir
	} else {
		// Documented development fallbacks:
		// 1. "../../data" (when working directory is apps/api)
		// 2. "../../../data" (when working directory is apps/api/tests)
		// 3. "data" (when working directory is repository root)
		if _, err := os.Stat("../../data/profile/profile.json"); err == nil {
			targetDir = "../../data"
		} else if _, err := os.Stat("../../../data/profile/profile.json"); err == nil {
			targetDir = "../../../data"
		} else {
			targetDir = "data"
		}
	}

	absDir, err := filepath.Abs(targetDir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path for portfolio data directory %q: %w", targetDir, err)
	}

	requiredFiles := []string{
		"profile/profile.json",
		"projects/projects.json",
		"skills/skills.json",
		"evidence/evidence.json",
		"journey/journey.json",
	}

	for _, reqFile := range requiredFiles {
		fullPath := filepath.Join(absDir, reqFile)
		info, err := os.Stat(fullPath)
		if err != nil || info.IsDir() {
			return "", fmt.Errorf("canonical portfolio data directory at %s is invalid: missing required file %s", absDir, reqFile)
		}
	}

	return absDir, nil
}


