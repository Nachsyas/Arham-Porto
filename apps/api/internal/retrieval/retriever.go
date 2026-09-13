package retrieval

import (
	"context"

	"github.com/nachsyas/arham-porto/apps/api/internal/domain"
)

// EvidenceRetriever defines the boundary for similarity search and evidence assembly.
// Full implementation using pgvector cosine distance scheduled for Phase 6.
type EvidenceRetriever interface {
	RetrieveEvidence(ctx context.Context, query string, limit int) ([]domain.Evidence, []domain.AICitation, error)
}
