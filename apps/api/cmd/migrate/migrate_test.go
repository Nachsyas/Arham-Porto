package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveMigrationsDir(t *testing.T) {
	// 1. Candidate paths from apps/api
	dir, err := ResolveMigrationsDir()
	if err != nil {
		t.Fatalf("expected to resolve migrations dir, got error: %v", err)
	}

	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		t.Fatalf("expected resolved path to be a directory, got %s", dir)
	}

	// 2. Explicit MIGRATIONS_DIR env override
	tempDir := t.TempDir()
	t.Setenv("MIGRATIONS_DIR", tempDir)

	resolved, err := ResolveMigrationsDir()
	if err != nil {
		t.Fatalf("expected to resolve explicit MIGRATIONS_DIR, got error: %v", err)
	}
	expected, _ := filepath.Abs(tempDir)
	if resolved != expected {
		t.Errorf("expected %s, got %s", expected, resolved)
	}

	// 3. Non-existent MIGRATIONS_DIR fails
	t.Setenv("MIGRATIONS_DIR", filepath.Join(tempDir, "non_existent"))
	_, err = ResolveMigrationsDir()
	if err == nil {
		t.Errorf("expected error for non-existent MIGRATIONS_DIR")
	}
}
