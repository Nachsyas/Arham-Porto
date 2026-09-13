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

// ResolveDataDir resolves and validates the path to canonical data/ directory.
func (c *Config) ResolveDataDir() (string, error) {
	candidates := []string{}
	if c.PortfolioDataDir != "" {
		candidates = append(candidates, c.PortfolioDataDir)
	}
	candidates = append(candidates,
		"data",
		"../data",
		"../../data",
		"../../../data",
		"../../../../data",
	)

	// Also check current working directory and walk up
	if cwd, err := os.Getwd(); err == nil {
		curr := cwd
		for i := 0; i < 5; i++ {
			candidates = append(candidates, filepath.Join(curr, "data"))
			parent := filepath.Dir(curr)
			if parent == curr {
				break
			}
			curr = parent
		}
	}

	requiredFiles := []string{
		"profile/profile.json",
		"projects/projects.json",
		"skills/skills.json",
		"evidence/evidence.json",
		"journey/journey.json",
	}

	for _, dir := range candidates {
		absDir, err := filepath.Abs(dir)
		if err != nil {
			continue
		}

		allFound := true
		for _, reqFile := range requiredFiles {
			fullPath := filepath.Join(absDir, reqFile)
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				allFound = false
				break
			}
		}

		if allFound {
			return absDir, nil
		}
	}

	return "", fmt.Errorf("canonical portfolio data directory not found in candidates: %v", candidates)
}


