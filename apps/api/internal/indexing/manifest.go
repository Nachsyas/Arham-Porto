package indexing

import (
	"bytes"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

// MaxSourceFileBytes defines the upper bound on individual source files eligible for indexing (256 KiB).
const MaxSourceFileBytes = 256 * 1024

// ManifestEntry defines an approved source path within a repository.
type ManifestEntry struct {
	Path       string `json:"path"`
	Title      string `json:"title"`
	SourceType string `json:"source_type"` // doc, manifest, entrypoint, architecture
	Required   bool   `json:"required"`
}

// RepositorySourceManifest defines the deterministic source selection policy for a repository.
type RepositorySourceManifest struct {
	Repository string          `json:"repository"`
	Ref        string          `json:"ref"`
	Entries    []ManifestEntry `json:"entries"`
}

// GetRepositoryManifest returns the explicit source manifest policy for an approved repository.
func GetRepositoryManifest(repo string) (*RepositorySourceManifest, bool) {
	manifests := map[string]RepositorySourceManifest{
		"Nachsyas/EduTrace": {
			Repository: "Nachsyas/EduTrace",
			Ref:        "main",
			Entries: []ManifestEntry{
				{Path: "README.md", Title: "EduTrace Architecture Specification & README", SourceType: "doc", Required: true},
				{Path: "package.json", Title: "EduTrace Root Manifest", SourceType: "manifest", Required: false},
				{Path: "contracts/foundry.toml", Title: "Foundry Smart Contract Configuration", SourceType: "manifest", Required: false},
				{Path: "contracts/src/EduTraceSBT.sol", Title: "Soulbound Token Contract", SourceType: "architecture", Required: false},
			},
		},
		"Nachsyas/gdgoc-ecommerce": {
			Repository: "Nachsyas/gdgoc-ecommerce",
			Ref:        "main",
			Entries: []ManifestEntry{
				{Path: "README.md", Title: "GDGOC E-Commerce Clean Architecture Specification", SourceType: "doc", Required: true},
				{Path: "backend/go.mod", Title: "Go Backend Dependency Manifest", SourceType: "manifest", Required: false},
				{Path: "backend/cmd/api/main.go", Title: "Go REST API Server Entrypoint", SourceType: "entrypoint", Required: false},
				{Path: "backend/internal/domain/user.go", Title: "User Domain Entity Specification", SourceType: "architecture", Required: false},
			},
		},
		"Nachsyas/maritime-ai-dashboard": {
			Repository: "Nachsyas/maritime-ai-dashboard",
			Ref:        "main",
			Entries: []ManifestEntry{
				{Path: "README.md", Title: "Maritime AI Dashboard Project Documentation", SourceType: "doc", Required: false},
				{Path: "Backend/requirements.txt", Title: "FastAPI Backend Dependency Manifest", SourceType: "manifest", Required: false},
				{Path: "Frontend/package.json", Title: "Next.js Frontend Manifest", SourceType: "manifest", Required: false},
			},
		},
		"Nachsyas/smart-kitchen-backend": {
			Repository: "Nachsyas/smart-kitchen-backend",
			Ref:        "main",
			Entries: []ManifestEntry{
				{Path: "README.md", Title: "Smart Kitchen Backend Service Specification", SourceType: "doc", Required: true},
				{Path: "go.mod", Title: "Go Fiber Backend Manifest", SourceType: "manifest", Required: false},
				{Path: "main.go", Title: "Smart Kitchen Backend Entrypoint", SourceType: "entrypoint", Required: false},
				{Path: "database/connection.go", Title: "GORM PostgreSQL Database Connection", SourceType: "architecture", Required: false},
			},
		},
		"Nachsyas/smart-kitchen-frontend": {
			Repository: "Nachsyas/smart-kitchen-frontend",
			Ref:        "main",
			Entries: []ManifestEntry{
				{Path: "README.md", Title: "Smart Kitchen Frontend UI Specification", SourceType: "doc", Required: true},
				{Path: "package.json", Title: "Frontend Next.js Manifest", SourceType: "manifest", Required: false},
				{Path: "next.config.ts", Title: "Next.js Build Configuration", SourceType: "manifest", Required: false},
			},
		},
	}

	norm := strings.TrimSpace(repo)
	for k, m := range manifests {
		if strings.EqualFold(k, norm) {
			return &m, true
		}
	}
	return nil, false
}

// IsPathExcluded checks defense-in-depth security rules against sensitive, secret-bearing, or non-indexer paths.
func IsPathExcluded(path string) (bool, string) {
	lower := strings.ToLower(strings.TrimSpace(path))

	// Secrets, credentials, private keys
	secretPatterns := []string{
		".env",
		"id_rsa",
		"id_ecdsa",
		"id_ed25519",
	}
	for _, p := range secretPatterns {
		if lower == p || strings.HasPrefix(lower, p+".") || strings.HasSuffix(lower, "/"+p) || strings.Contains(lower, "/"+p+".") {
			return true, "secret-bearing pattern match"
		}
	}

	if strings.HasSuffix(lower, ".pem") || strings.HasSuffix(lower, ".key") || strings.HasSuffix(lower, ".crt") {
		return true, "certificate or private key file"
	}
	if strings.HasPrefix(filepath.Base(lower), "credentials.") || strings.HasPrefix(filepath.Base(lower), "service-account") {
		return true, "credential file pattern"
	}

	// Lockfiles & dependency caches
	lockfiles := []string{"package-lock.json", "pnpm-lock.yaml", "yarn.lock", "go.sum", "composer.lock", "cargo.lock"}
	for _, lf := range lockfiles {
		if filepath.Base(lower) == lf {
			return true, "dependency lockfile"
		}
	}

	// Directory exclusions
	excludedDirs := []string{"node_modules/", "vendor/", ".git/", ".github/", "dist/", "build/", ".next/", "out/"}
	for _, d := range excludedDirs {
		if strings.HasPrefix(lower, d) || strings.Contains(lower, "/"+d) {
			return true, "excluded build or dependency directory"
		}
	}

	// Binary assets
	binaryExts := []string{".png", ".jpg", ".jpeg", ".gif", ".ico", ".svg", ".pdf", ".zip", ".tar", ".gz", ".wasm", ".exe", ".bin"}
	for _, ext := range binaryExts {
		if strings.HasSuffix(lower, ext) {
			return true, "binary or media asset extension"
		}
	}

	return false, ""
}

// IsBinaryContent inspects byte slice for NUL bytes or invalid UTF-8 encoding.
func IsBinaryContent(data []byte) bool {
	if bytes.IndexByte(data, 0) != -1 {
		return true
	}
	return !utf8.Valid(data)
}
