package ai

import (
	"context"
	"io"
)

// LLMProvider defines the provider-agnostic interface for AI completions.
// Concrete implementations (Google Gemini, etc.) will be added in Phase 6.
type LLMProvider interface {
	GenerateResponse(ctx context.Context, prompt string) (string, error)
	StreamResponse(ctx context.Context, prompt string) (io.ReadCloser, error)
}

// EmbeddingProvider defines the provider-agnostic interface for generating vector embeddings.
type EmbeddingProvider interface {
	GenerateEmbedding(ctx context.Context, text string) ([]float32, error)
}
