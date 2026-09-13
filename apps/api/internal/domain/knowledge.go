package domain

import (
	"context"
	"time"
)

// KnowledgeSource represents an ingested canonical or repository document.
type KnowledgeSource struct {
	ID             string    `json:"id"`
	Repository     string    `json:"repository"`
	RepositoryURL  string    `json:"repository_url"`
	Ref            string    `json:"ref"`
	CommitSHA      string    `json:"commit_sha"`
	Path           string    `json:"path"`
	Title          string    `json:"title"`
	SourceType     string    `json:"source_type"`
	Checksum       string    `json:"checksum"`
	ApprovalStatus string    `json:"approval_status"`
	IndexedAt      time.Time `json:"indexed_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// KnowledgeChunk represents an atomic text segment for embedding and retrieval.
type KnowledgeChunk struct {
	ID                  string    `json:"id"`
	SourceID            string    `json:"source_id"`
	ChunkIndex          int       `json:"chunk_index"`
	Content             string    `json:"content"`
	Checksum            string    `json:"checksum"`
	Embedding           []float32 `json:"embedding,omitempty"`
	EmbeddingProvider   string    `json:"embedding_provider,omitempty"`
	EmbeddingModel      string    `json:"embedding_model,omitempty"`
	EmbeddingDimensions int       `json:"embedding_dimensions,omitempty"`
	ProjectID           *string   `json:"project_id,omitempty"`
	SkillIDs            []string  `json:"skill_ids,omitempty"`
	EvidenceID          *string   `json:"evidence_id,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
}

// SourceWithChunks binds a source document to its constituent chunks for atomic repository snapshots.
type SourceWithChunks struct {
	Source KnowledgeSource `json:"source"`
	Chunks []KnowledgeChunk `json:"chunks"`
}

// RetrievalResult encapsulates a similarity match with complete provenance metadata.
type RetrievalResult struct {
	ChunkID     string   `json:"chunk_id"`
	SourceID    string   `json:"source_id"`
	Repository  string   `json:"repository"`
	Path        string   `json:"path"`
	Ref         string   `json:"ref"`
	CommitSHA   string   `json:"commit_sha"`
	SourceURL   string   `json:"source_url"`
	SourceTitle string   `json:"source_title"`
	SourceType  string   `json:"source_type"`
	Content     string   `json:"content"`
	Distance    float64  `json:"distance"`
	Similarity  float64  `json:"similarity"`
	ProjectID   *string  `json:"project_id,omitempty"`
	SkillIDs    []string `json:"skill_ids,omitempty"`
	EvidenceID  *string  `json:"evidence_id,omitempty"`
}

// AICitation represents a structured reference returned by AI services.
type AICitation struct {
	SourceType string  `json:"source_type"`
	Title      string  `json:"title"`
	URL        *string `json:"url,omitempty"`
	Confidence string  `json:"confidence"` // verified, supporting
}

// RepositoryIndexStatus records the operational health and completeness of a repository's index.
type RepositoryIndexStatus struct {
	Repository        string     `json:"repository"`
	CommitSHA         string     `json:"commit_sha"`
	SourceCount       int        `json:"source_count"`
	ChunkCount        int        `json:"chunk_count"`
	EmbeddingProvider string     `json:"embedding_provider"`
	EmbeddingModel    string     `json:"embedding_model"`
	Status            string     `json:"status"` // metadata_only, vector_ready, revoked, failed
	LastIndexedAt     *time.Time `json:"last_indexed_at,omitempty"`
}

// KnowledgeRepository defines vector and metadata storage access for knowledge indexing and retrieval.
// Pure Go domain interface without third-party dependencies.
type KnowledgeRepository interface {
	ReplaceRepositorySnapshot(ctx context.Context, repository string, commitSHA string, items []SourceWithChunks) error
	SearchSimilar(ctx context.Context, queryVector []float32, provider, model string, dimensions int, repoFilter, sourceTypeFilter string, limit int) ([]RetrievalResult, error)
	DeleteRepositorySources(ctx context.Context, repository string) error
	GetIndexStatus(ctx context.Context) ([]RepositoryIndexStatus, error)
}
