package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

// ResolveMigrationsDir determines the migrations directory path deterministically.
func ResolveMigrationsDir() (string, error) {
	if envDir := os.Getenv("MIGRATIONS_DIR"); envDir != "" {
		if info, err := os.Stat(envDir); err == nil && info.IsDir() {
			return filepath.Abs(envDir)
		}
		return "", fmt.Errorf("configured MIGRATIONS_DIR %q not found", envDir)
	}

	// Documented candidate paths:
	candidates := []string{
		"/app/migrations",
		"apps/api/migrations",
		"migrations",
		"../../migrations",
	}

	for _, cand := range candidates {
		if info, err := os.Stat(cand); err == nil && info.IsDir() {
			return filepath.Abs(cand)
		}
	}

	return "", fmt.Errorf("migrations directory not found in candidate paths: %v", candidates)
}

// RunMigrations connects to the database, initializes schema_migrations if needed,
// and applies pending up-migrations in deterministic version order inside transactions.
func RunMigrations(ctx context.Context, dbURL, migrationsDir string) error {
	log.Printf("[migrate] Starting migration runner (dir: %s)", migrationsDir)

	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer conn.Close(ctx)

	// Ensure schema_migrations table exists
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`
	if _, err := conn.Exec(ctx, createTableQuery); err != nil {
		return fmt.Errorf("failed to ensure schema_migrations table: %w", err)
	}

	// Query already applied versions
	rows, err := conn.Query(ctx, "SELECT version FROM schema_migrations")
	if err != nil {
		return fmt.Errorf("failed to query schema_migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return fmt.Errorf("failed to scan migration version: %w", err)
		}
		applied[v] = true
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error reading applied migrations: %w", err)
	}

	// Find all *.up.sql files in migrationsDir
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var upFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".up.sql") {
			upFiles = append(upFiles, entry.Name())
		}
	}

	sort.Strings(upFiles) // Deterministic order: 000001_..., 000002_...

	if len(upFiles) == 0 {
		log.Println("[migrate] No .up.sql migration files found")
		return nil
	}

	var newlyApplied int
	for _, file := range upFiles {
		version := strings.TrimSuffix(file, ".up.sql")
		if applied[version] {
			log.Printf("[migrate] Version %s already applied (skipping)", version)
			continue
		}

		filePath := filepath.Join(migrationsDir, file)
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		log.Printf("[migrate] Applying migration %s...", version)

		tx, err := conn.Begin(ctx)
		if err != nil {
			return fmt.Errorf("failed to begin transaction for %s: %w", version, err)
		}

		if _, err := tx.Exec(ctx, string(content)); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("failed to execute migration %s: %w", version, err)
		}

		if _, err := tx.Exec(ctx, "INSERT INTO schema_migrations (version, applied_at) VALUES ($1, NOW())", version); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("failed to record applied migration %s: %w", version, err)
		}

		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", version, err)
		}

		log.Printf("[migrate] Successfully applied migration %s", version)
		newlyApplied++
	}

	log.Printf("[migrate] Complete. Newly applied: %d, Total tracked: %d", newlyApplied, len(applied)+newlyApplied)
	return nil
}

func main() {
	dbURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if dbURL == "" {
		log.Fatalf("[migrate] Fatal: DATABASE_URL environment variable is required")
	}

	migrationsDir, err := ResolveMigrationsDir()
	if err != nil {
		log.Fatalf("[migrate] Fatal: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	if err := RunMigrations(ctx, dbURL, migrationsDir); err != nil {
		log.Fatalf("[migrate] Fatal: migration failed: %v", err)
	}
}
