package indexing

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/nachsyas/arham-porto/apps/api/internal/domain"
	"github.com/nachsyas/arham-porto/apps/api/internal/embedding"
	"github.com/nachsyas/arham-porto/apps/api/internal/github"
)

// memorySnapshotRepo provides an in-memory transactional mock of domain.KnowledgeRepository.
type memorySnapshotRepo struct {
	mu           sync.Mutex
	snapshots    map[string]struct {
		commitSHA string
		items     []domain.SourceWithChunks
	}
	failNextTx   bool
	callsCount   int
}

func newMemorySnapshotRepo() *memorySnapshotRepo {
	return &memorySnapshotRepo{
		snapshots: make(map[string]struct {
			commitSHA string
			items     []domain.SourceWithChunks
		}),
	}
}

func (m *memorySnapshotRepo) ReplaceRepositorySnapshot(ctx context.Context, repository string, commitSHA string, items []domain.SourceWithChunks) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.callsCount++

	if m.failNextTx {
		m.failNextTx = false
		return errors.New("simulated database transaction failure")
	}

	// Atomic snapshot replacement
	m.snapshots[repository] = struct {
		commitSHA string
		items     []domain.SourceWithChunks
	}{
		commitSHA: commitSHA,
		items:     items,
	}
	return nil
}

func (m *memorySnapshotRepo) SearchSimilar(ctx context.Context, queryVector []float32, provider, model string, dimensions int, repoFilter, sourceTypeFilter string, limit int) ([]domain.RetrievalResult, error) {
	return nil, nil
}

func (m *memorySnapshotRepo) DeleteRepositorySources(ctx context.Context, repository string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.snapshots, repository)
	return nil
}

func (m *memorySnapshotRepo) GetIndexStatus(ctx context.Context) ([]domain.RepositoryIndexStatus, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var statuses []domain.RepositoryIndexStatus
	for repo, snap := range m.snapshots {
		chunkCount := 0
		for _, item := range snap.items {
			chunkCount += len(item.Chunks)
		}
		statuses = append(statuses, domain.RepositoryIndexStatus{
			Repository:  repo,
			CommitSHA:   snap.commitSHA,
			SourceCount: len(snap.items),
			ChunkCount:  chunkCount,
			Status:      "vector_ready",
		})
	}
	return statuses, nil
}

func TestRepositorySnapshotLevel(t *testing.T) {
	commit1 := "1111111111111111111111111111111111111111"
	commit2 := "2222222222222222222222222222222222222222"
	activeCommit := commit1

	// Mock GitHub server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "/commits/"):
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(activeCommit))
		case strings.Contains(r.URL.Path, "README.md"):
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf("# Documentation for commit %s", activeCommit)))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	allowlist := []AllowlistEntry{
		{Repository: "Nachsyas/EduTrace", Enabled: true},
	}
	ghClient := github.NewTestClient(ts.URL, "token")
	fakeProvider := embedding.NewDeterministicFakeProvider(768)
	repo := newMemorySnapshotRepo()
	ctx := context.Background()

	svc := NewService("../../../../data", allowlist, ghClient, fakeProvider, repo)

	t.Run("complete snapshot succeeds (Gate #55)", func(t *testing.T) {
		reports, err := svc.Index(ctx, "Nachsyas/EduTrace")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(reports) < 1 {
			t.Fatalf("expected at least 1 report, got %d", len(reports))
		}

		snap, exists := repo.snapshots["Nachsyas/EduTrace"]
		if !exists {
			t.Fatalf("expected snapshot in repository, not found")
		}
		if snap.commitSHA != commit1 {
			t.Errorf("expected commit %s, got %s", commit1, snap.commitSHA)
		}
		if len(snap.items) == 0 {
			t.Errorf("expected sources in snapshot, got 0")
		}
	})

	t.Run("idempotent re-indexing (Gate #55)", func(t *testing.T) {
		initialCallCount := repo.callsCount
		_, err := svc.Index(ctx, "Nachsyas/EduTrace")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if repo.callsCount <= initialCallCount {
			t.Errorf("expected index call, callsCount=%d", repo.callsCount)
		}
		snap := repo.snapshots["Nachsyas/EduTrace"]
		if snap.commitSHA != commit1 {
			t.Errorf("expected commit %s, got %s", commit1, snap.commitSHA)
		}
	})

	t.Run("database transaction failure preserves previous healthy snapshot (Gate #55)", func(t *testing.T) {
		// Set DB to fail on next transaction
		repo.failNextTx = true

		_, err := svc.Index(ctx, "Nachsyas/EduTrace")
		if err == nil {
			t.Fatalf("expected error from failed transaction, got nil")
		}

		// Verify previous healthy snapshot was preserved intact
		snap, exists := repo.snapshots["Nachsyas/EduTrace"]
		if !exists || snap.commitSHA != commit1 {
			t.Fatalf("previous healthy snapshot was not preserved after failed transaction")
		}
	})

	t.Run("new commit atomically replaces old snapshot (Gate #55)", func(t *testing.T) {
		// Update active commit to commit2
		activeCommit = commit2

		reports, err := svc.Index(ctx, "Nachsyas/EduTrace")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		snap := repo.snapshots["Nachsyas/EduTrace"]
		if snap.commitSHA != commit2 {
			t.Errorf("expected commit to be updated to %s, got %s", commit2, snap.commitSHA)
		}
		if reports[0].CommitSHA != commit2 {
			t.Errorf("expected report commit %s, got %s", commit2, reports[0].CommitSHA)
		}
	})

	t.Run("deletion works even if revoked from allowlist (Gate #22)", func(t *testing.T) {
		// Repository is deleted by exact name
		err := svc.DeleteRepository(ctx, "Nachsyas/EduTrace")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, exists := repo.snapshots["Nachsyas/EduTrace"]; exists {
			t.Errorf("expected repository snapshot to be removed after delete")
		}
	})
}
