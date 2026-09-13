package domain

import (
	"context"
	"time"
)

// KnowledgeChunk represents an indexed documentation or code snippet.
type KnowledgeChunk struct {
	ID         string    `json:"id"`
	Repository string    `json:"repository"`
	FilePath   string    `json:"file_path"`
	CommitSHA  *string   `json:"commit_sha,omitempty"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	SourceURL  string    `json:"source_url"`
	IndexedAt  time.Time `json:"indexed_at"`
	Approved   bool      `json:"approved"`
}

// AICitation represents a structured reference returned by Ask Arham AI.
type AICitation struct {
	SourceType string  `json:"source_type"`
	Title      string  `json:"title"`
	URL        *string `json:"url,omitempty"`
	Confidence string  `json:"confidence"` // verified, supporting
}

// KnowledgeRepository defines vector storage access for knowledge chunks.
type KnowledgeRepository interface {
	SearchSimilar(ctx context.Context, queryEmbedding []float32, limit int) ([]KnowledgeChunk, error)
	StoreChunk(ctx context.Context, chunk *KnowledgeChunk, embedding []float32) error
}
