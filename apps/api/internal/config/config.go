package config

import (
	"os"
)

// Config holds the application configuration.
type Config struct {
	AppEnv      string
	Port        string
	WebOrigin   string
	DatabaseURL string
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

	webOrigin := os.Getenv("WEB_ORIGIN")
	if webOrigin == "" {
		webOrigin = "http://localhost:3000"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://arham_user:arham_password@localhost:5432/arham_porto?sslmode=disable"
	}

	return &Config{
		AppEnv:      appEnv,
		Port:        port,
		WebOrigin:   webOrigin,
		DatabaseURL: dbURL,
	}
}
