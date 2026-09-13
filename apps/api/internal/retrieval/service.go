package retrieval

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/nachsyas/arham-porto/apps/api/internal/domain"
	"github.com/nachsyas/arham-porto/apps/api/internal/embedding"
)

var (
	ErrQueryEmpty       = errors.New("search query cannot be empty")
	ErrEmbeddingMissing = errors.New("embedding provider not configured for retrieval")
)

// SearchOptions parametrizes internal vector retrieval queries.
type SearchOptions struct {
	RepositoryFilter string `json:"repository_filter,omitempty"`
	SourceTypeFilter string `json:"source_type_filter,omitempty"`
	Limit            int    `json:"limit,omitempty"`
}

// Service provides internal semantic evidence retrieval with exact pgvector cosine distance.
// Phase 5: Internal service, tests, and evals only — strictly NO public HTTP routes (Gate #42)
// and NO generative chat/answer synthesis (Gate #43).
type Service struct {
	repo     domain.KnowledgeRepository
	provider embedding.Provider
}

// NewService creates a new internal retrieval service.
func NewService(repo domain.KnowledgeRepository, provider embedding.Provider) *Service {
	return &Service{
		repo:     repo,
		provider: provider,
	}
}

// Retrieve searches for the most relevant evidence chunks matching a query.
// Embeds the query with PurposeQuery (Gate #3) and filters strictly by active model identity (Gate #2, #7).
func (s *Service) Retrieve(ctx context.Context, query string, opts SearchOptions) ([]domain.RetrievalResult, error) {
	trimmedQuery := strings.TrimSpace(query)
	if trimmedQuery == "" {
		return nil, ErrQueryEmpty
	}
	if s.provider == nil || s.provider.Model() == "disabled" {
		return nil, ErrEmbeddingMissing
	}

	// 1. Embed query with PurposeQuery (Gate #3, #29)
	vectors, err := s.provider.Embed(ctx, embedding.PurposeQuery, []string{trimmedQuery})
	if err != nil {
		return nil, fmt.Errorf("failed generating query embedding: %w", err)
	}
	if len(vectors) != 1 {
		return nil, errors.New("expected 1 query vector from embedding provider")
	}
	queryVec := vectors[0]

	// 2. Validate vector dimensions (Gate #4)
	if err := embedding.ValidateVector(queryVec, s.provider.Dimensions()); err != nil {
		return nil, fmt.Errorf("invalid query vector: %w", err)
	}

	// 3. Exact search with strict model filtering (Gate #2, #5, #7, #40)
	results, err := s.repo.SearchSimilar(
		ctx,
		queryVec,
		s.provider.ProviderName(),
		s.provider.Model(),
		s.provider.Dimensions(),
		opts.RepositoryFilter,
		opts.SourceTypeFilter,
		opts.Limit,
	)
	if err != nil {
		return nil, fmt.Errorf("retrieval search failed: %w", err)
	}

	return results, nil
}
