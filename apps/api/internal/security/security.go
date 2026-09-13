package security

import (
	"context"
	"strings"
)

// InputSanitizer defines the boundary for cleaning untrusted user input.
type InputSanitizer interface {
	Sanitize(input string) string
}

// SimpleSanitizer provides a basic string trimmer/sanitizer stub for Phase 0.
type SimpleSanitizer struct{}

func (s *SimpleSanitizer) Sanitize(input string) string {
	return strings.TrimSpace(input)
}

// RateLimiter defines the interface for throttling API consumers.
type RateLimiter interface {
	Allow(ctx context.Context, key string) bool
}
