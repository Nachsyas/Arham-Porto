package indexing

import (
	"context"
)

// Pipeline defines the interface for indexing approved repositories into pgvector.
// Implementation scheduled for Phase 5.
type Pipeline interface {
	IndexRepository(ctx context.Context, repo string) error
	ReindexAll(ctx context.Context) error
}
