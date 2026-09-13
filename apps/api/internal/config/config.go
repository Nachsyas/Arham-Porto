package config

import (
	"fmt"
	"os"
	"path/filepath"
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
	DatabaseURL      string
	DatabaseMode     string // disabled, optional, required
	PortfolioDataDir string
}

// Load loads environment variables with safe development fallbacks.
func Load() *Config {
	port := os.Getenv("API_PORT")
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

	return &Config{
		AppEnv:           appEnv,
		Port:             port,
		AllowedOrigins:   allowedOrigins,
		DatabaseURL:      dbURL,
		DatabaseMode:     dbMode,
		PortfolioDataDir: dataDir,
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


