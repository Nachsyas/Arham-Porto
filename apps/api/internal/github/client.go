package github

import (
	"context"

	"github.com/nachsyas/arham-porto/apps/api/internal/domain"
)

// Client defines the interface for interacting with the GitHub API.
// Implementation scheduled for Phase 5.
type Client interface {
	FetchRepositoryReadme(ctx context.Context, repo string) (string, error)
	FetchRepositoryMetadata(ctx context.Context, repo string) (*domain.Project, error)
}
