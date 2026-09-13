package dto

import (
	"net/url"
	"path/filepath"
	"strings"
)

// SanitizeStringPtr safely handles nullable strings.
// Returns nil if s is nil, empty, or contains TODO markers.
func SanitizeStringPtr(s *string) *string {
	if s == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*s)
	if trimmed == "" || strings.Contains(trimmed, "TODO_") {
		return nil
	}
	return &trimmed
}

// SanitizeExternalURL validates that an external link uses a safe scheme (strictly https://, or http://localhost for local dev).
// Reject file://, javascript:, data:, and malformed URLs.
func SanitizeExternalURL(u *string) *string {
	if u == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*u)
	if trimmed == "" || strings.Contains(trimmed, "TODO_") {
		return nil
	}

	lower := strings.ToLower(trimmed)
	if strings.HasPrefix(lower, "javascript:") ||
		strings.HasPrefix(lower, "file:") ||
		strings.HasPrefix(lower, "data:") ||
		strings.HasPrefix(lower, "vbscript:") ||
		strings.HasPrefix(lower, "blob:") {
		return nil
	}

	parsed, err := url.ParseRequestURI(trimmed)
	if err != nil {
		return nil
	}

	// Must be https:// (or http://localhost:PORT for local development)
	if parsed.Scheme == "https" {
		return &trimmed
	}
	if parsed.Scheme == "http" && (parsed.Hostname() == "localhost" || parsed.Hostname() == "127.0.0.1") {
		return &trimmed
	}

	return nil
}

// SanitizeImageURL validates that an image URL uses safe schemes (https:// or repository-relative /images/...).
// Rejects javascript:, file:, data:.
func SanitizeImageURL(u *string) *string {
	if u == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*u)
	if trimmed == "" || strings.Contains(trimmed, "TODO_") {
		return nil
	}

	lower := strings.ToLower(trimmed)
	if strings.HasPrefix(lower, "javascript:") ||
		strings.HasPrefix(lower, "file:") ||
		strings.HasPrefix(lower, "data:") ||
		strings.HasPrefix(lower, "vbscript:") ||
		strings.HasPrefix(lower, "blob:") {
		return nil
	}

	// Permitted: relative asset paths starting with "/" (e.g. /images/hero.webp)
	if strings.HasPrefix(trimmed, "/") && !strings.HasPrefix(trimmed, "//") {
		clean := filepath.Clean(trimmed)
		return &clean
	}

	// External image: must be https://
	parsed, err := url.ParseRequestURI(trimmed)
	if err != nil {
		return nil
	}
	if parsed.Scheme == "https" {
		return &trimmed
	}

	return nil
}

// SanitizeSourcePath validates that source_path is strictly repository-relative.
// Absolute filesystem paths (Unix /Users/... or Windows C:\...), path traversal (..),
// and user home directories are strictly rejected and return nil.
func SanitizeSourcePath(p *string) *string {
	if p == nil {
		return nil
	}
	clean := filepath.ToSlash(strings.TrimSpace(*p))
	if clean == "" || strings.Contains(clean, "TODO_") {
		return nil
	}

	// Reject Windows drive letter (e.g. C:/, D:/)
	if len(clean) >= 2 && clean[1] == ':' {
		return nil
	}

	// Reject absolute Unix paths
	if strings.HasPrefix(clean, "/") || strings.HasPrefix(clean, "\\") {
		return nil
	}

	// Reject path traversal components
	parts := strings.Split(clean, "/")
	for _, part := range parts {
		if part == ".." {
			return nil
		}
	}

	// Reject user home directory references
	lower := strings.ToLower(clean)
	if strings.HasPrefix(lower, "users/") || strings.HasPrefix(lower, "home/") || strings.HasPrefix(lower, "root/") {
		return nil
	}

	cleaned := filepath.Clean(clean)
	cleaned = filepath.ToSlash(cleaned)
	if strings.HasPrefix(cleaned, "/") || strings.HasPrefix(cleaned, "../") {
		return nil
	}

	return &cleaned
}

// IsTodoBacked checks whether a field is marked with a TODO indicator in the domain entity metadata.
func IsTodoBacked(todos []string, field string) bool {
	upperField := strings.ToUpper(field)
	for _, t := range todos {
		if strings.Contains(strings.ToUpper(t), upperField) {
			return true
		}
	}
	return false
}
