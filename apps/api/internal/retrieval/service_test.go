package retrieval

import (
	"context"
	"errors"
	"testing"

	"github.com/nachsyas/arham-porto/apps/api/internal/domain"
	"github.com/nachsyas/arham-porto/apps/api/internal/embedding"
)

// mockKnowledgeRepo implements domain.KnowledgeRepository for unit testing.
type mockKnowledgeRepo struct {
	searchCalled      bool
	capturedProvider  string
	capturedModel     string
	capturedDims      int
	capturedRepo      string
	capturedType      string
	capturedLimit     int
	resultsToReturn   []domain.RetrievalResult
	searchErrorToReturn error
}

func (m *mockKnowledgeRepo) ReplaceRepositorySnapshot(ctx context.Context, repository string, commitSHA string, items []domain.SourceWithChunks) error {
	return nil
}

func (m *mockKnowledgeRepo) SearchSimilar(
	ctx context.Context,
	queryVector []float32,
	provider, model string,
	dimensions int,
	repoFilter, sourceTypeFilter string,
	limit int,
) ([]domain.RetrievalResult, error) {
	m.searchCalled = true
	m.capturedProvider = provider
	m.capturedModel = model
	m.capturedDims = dimensions
	m.capturedRepo = repoFilter
	m.capturedType = sourceTypeFilter
	m.capturedLimit = limit
	if m.searchErrorToReturn != nil {
		return nil, m.searchErrorToReturn
	}
	return m.resultsToReturn, nil
}

func (m *mockKnowledgeRepo) DeleteRepositorySources(ctx context.Context, repository string) error {
	return nil
}

func (m *mockKnowledgeRepo) GetIndexStatus(ctx context.Context) ([]domain.RepositoryIndexStatus, error) {
	return nil, nil
}

func TestRetrievalService(t *testing.T) {
	fakeProvider := embedding.NewDeterministicFakeProvider(768)
	ctx := context.Background()

	t.Run("empty query returns error", func(t *testing.T) {
		repo := &mockKnowledgeRepo{}
		svc := NewService(repo, fakeProvider)
		_, err := svc.Retrieve(ctx, "   ", SearchOptions{})
		if err == nil || !errors.Is(err, ErrQueryEmpty) {
			t.Fatalf("expected ErrQueryEmpty, got %v", err)
		}
	})

	t.Run("disabled provider returns error", func(t *testing.T) {
		repo := &mockKnowledgeRepo{}
		svc := NewService(repo, embedding.NewDisabledProvider(768))
		_, err := svc.Retrieve(ctx, "valid query", SearchOptions{})
		if err == nil || !errors.Is(err, ErrEmbeddingMissing) {
			t.Fatalf("expected ErrEmbeddingMissing, got %v", err)
		}
	})

	t.Run("successful retrieval invokes query purpose and passes model identity", func(t *testing.T) {
		projID := "edutrace"
		mockResults := []domain.RetrievalResult{
			{
				ChunkID:     "chunk-1",
				SourceID:    "src-1",
				Repository:  "Nachsyas/EduTrace",
				Path:        "contracts/src/EduTraceSBT.sol",
				SourceTitle: "Soulbound Token Contract",
				Content:     "Soulbound Token implementation",
				Distance:    0.15,
				Similarity:  0.85,
				ProjectID:   &projID,
			},
		}

		repo := &mockKnowledgeRepo{resultsToReturn: mockResults}
		svc := NewService(repo, fakeProvider)

		opts := SearchOptions{
			RepositoryFilter: "Nachsyas/EduTrace",
			SourceTypeFilter: "architecture",
			Limit:            5,
		}

		results, err := svc.Retrieve(ctx, "Soulbound Token experience", opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		if results[0].ChunkID != "chunk-1" {
			t.Errorf("expected chunk-1, got %s", results[0].ChunkID)
		}
		if results[0].Similarity != 0.85 {
			t.Errorf("expected similarity 0.85, got %f", results[0].Similarity)
		}

		// Verify model identity filtering passed to repository (Gate #2, #7)
		if repo.capturedProvider != fakeProvider.ProviderName() {
			t.Errorf("expected provider %s, got %s", fakeProvider.ProviderName(), repo.capturedProvider)
		}
		if repo.capturedModel != fakeProvider.Model() {
			t.Errorf("expected model %s, got %s", fakeProvider.Model(), repo.capturedModel)
		}
		if repo.capturedDims != 768 {
			t.Errorf("expected dimensions 768, got %d", repo.capturedDims)
		}
		if repo.capturedRepo != "Nachsyas/EduTrace" {
			t.Errorf("expected repo filter Nachsyas/EduTrace, got %s", repo.capturedRepo)
		}
		if repo.capturedType != "architecture" {
			t.Errorf("expected source type filter architecture, got %s", repo.capturedType)
		}

		// Verify embedding purpose was PurposeQuery (Gate #3, #29)
		if fakeProvider.GetLastPurpose() != embedding.PurposeQuery {
			t.Errorf("expected purpose %s, got %s", embedding.PurposeQuery, fakeProvider.GetLastPurpose())
		}
	})
}
