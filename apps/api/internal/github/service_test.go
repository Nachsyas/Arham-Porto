package github

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/nachsyas/arham-porto/apps/api/internal/domain"
)

type mockGithubProjectRepo struct {
	mu           sync.Mutex
	projects     map[int64]*domain.GithubProject
	upsertFunc   func(ctx context.Context, p *domain.GithubProject) (bool, error)
	createCalls  int
	updateCalls  int
	upsertCalls  int
	listCalls    int
}

func newMockGithubProjectRepo() *mockGithubProjectRepo {
	return &mockGithubProjectRepo{
		projects: make(map[int64]*domain.GithubProject),
	}
}

func (m *mockGithubProjectRepo) Create(ctx context.Context, p *domain.GithubProject) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.createCalls++
	m.projects[p.GithubID] = p
	return nil
}

func (m *mockGithubProjectRepo) Update(ctx context.Context, p *domain.GithubProject) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.updateCalls++
	m.projects[p.GithubID] = p
	return nil
}

func (m *mockGithubProjectRepo) Upsert(ctx context.Context, p *domain.GithubProject) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.upsertCalls++
	if m.upsertFunc != nil {
		return m.upsertFunc(ctx, p)
	}

	_, exists := m.projects[p.GithubID]
	m.projects[p.GithubID] = p
	return !exists, nil
}

func (m *mockGithubProjectRepo) List(ctx context.Context) ([]*domain.GithubProject, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.listCalls++
	var res []*domain.GithubProject
	for _, p := range m.projects {
		res = append(res, p)
	}
	return res, nil
}

func (m *mockGithubProjectRepo) GetByGithubID(ctx context.Context, githubID int64) (*domain.GithubProject, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	p, ok := m.projects[githubID]
	if !ok {
		return nil, nil
	}
	return p, nil
}

func TestSyncService(t *testing.T) {
	t.Run("insert new repository", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{
					"id": 101,
					"name": "EduTrace",
					"full_name": "Nachsyas/EduTrace",
					"description": "Evidence-grounded portfolio sync",
					"html_url": "https://github.com/Nachsyas/EduTrace",
					"stargazers_count": 15,
					"forks_count": 3,
					"topics": ["go", "security"]
				}
			]`))
		}))
		defer ts.Close()

		client := NewTestClient(ts.URL, "valid-token")
		repo := newMockGithubProjectRepo()
		svc := NewSyncService(client, repo, "Nachsyas")

		result, err := svc.SyncRepositories(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.Synced != 1 || result.Created != 1 || result.Updated != 0 {
			t.Errorf("unexpected sync result: %+v", result)
		}

		if repo.upsertCalls != 1 {
			t.Errorf("expected 1 upsert call, got %d", repo.upsertCalls)
		}

		stored, _ := repo.GetByGithubID(context.Background(), 101)
		if stored == nil || stored.Name != "EduTrace" || stored.Stars != 15 {
			t.Errorf("stored project mismatch: %+v", stored)
		}
	})

	t.Run("update existing repository", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{
					"id": 101,
					"name": "EduTrace",
					"full_name": "Nachsyas/EduTrace",
					"description": "Updated description",
					"html_url": "https://github.com/Nachsyas/EduTrace",
					"stargazers_count": 25,
					"forks_count": 5,
					"topics": ["go", "security", "v2"]
				}
			]`))
		}))
		defer ts.Close()

		client := NewTestClient(ts.URL, "valid-token")
		repo := newMockGithubProjectRepo()

		// Pre-populate existing repository
		repo.projects[101] = &domain.GithubProject{
			GithubID: 101,
			Name:     "EduTrace",
			Stars:    15,
		}

		svc := NewSyncService(client, repo, "Nachsyas")
		result, err := svc.SyncRepositories(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.Synced != 1 || result.Created != 0 || result.Updated != 1 {
			t.Errorf("unexpected sync result: %+v", result)
		}

		stored, _ := repo.GetByGithubID(context.Background(), 101)
		if stored.Stars != 25 || stored.Description != "Updated description" {
			t.Errorf("updated project fields mismatch: %+v", stored)
		}
	})

	t.Run("mixed insert and update", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{
					"id": 201,
					"name": "Project-A",
					"full_name": "Nachsyas/Project-A",
					"html_url": "https://github.com/Nachsyas/Project-A"
				},
				{
					"id": 202,
					"name": "Project-B",
					"full_name": "Nachsyas/Project-B",
					"html_url": "https://github.com/Nachsyas/Project-B"
				}
			]`))
		}))
		defer ts.Close()

		client := NewTestClient(ts.URL, "valid-token")
		repo := newMockGithubProjectRepo()

		// 201 already exists, 202 is new
		repo.projects[201] = &domain.GithubProject{
			GithubID: 201,
			Name:     "Project-A",
		}

		svc := NewSyncService(client, repo, "Nachsyas")
		result, err := svc.SyncRepositories(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.Synced != 2 || result.Created != 1 || result.Updated != 1 {
			t.Errorf("expected 2 synced, 1 created, 1 updated; got %+v", result)
		}
	})

	t.Run("historical records retention - deleted repo in github is not deleted", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// GitHub returns only project 301
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{"id": 301, "name": "active-repo", "full_name": "Nachsyas/active-repo"}
			]`))
		}))
		defer ts.Close()

		client := NewTestClient(ts.URL, "valid-token")
		repo := newMockGithubProjectRepo()

		// Pre-populate active-repo (301) and historical deleted-repo (999)
		repo.projects[301] = &domain.GithubProject{GithubID: 301, Name: "active-repo"}
		repo.projects[999] = &domain.GithubProject{GithubID: 999, Name: "historical-deleted-repo"}

		svc := NewSyncService(client, repo, "Nachsyas")
		result, err := svc.SyncRepositories(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.Synced != 1 {
			t.Errorf("expected 1 synced, got %d", result.Synced)
		}

		// Historical repo 999 must still exist!
		historical, _ := repo.GetByGithubID(context.Background(), 999)
		if historical == nil || historical.Name != "historical-deleted-repo" {
			t.Fatalf("historical record was unexpectedly deleted!")
		}
	})

	t.Run("synchronization failure - GitHub API error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()

		client := NewTestClient(ts.URL, "valid-token")
		repo := newMockGithubProjectRepo()
		svc := NewSyncService(client, repo, "Nachsyas")

		_, err := svc.SyncRepositories(context.Background())
		if err == nil {
			t.Fatalf("expected error from failed github API, got nil")
		}

		if repo.upsertCalls != 0 {
			t.Errorf("expected 0 upsert calls on fetch failure, got %d", repo.upsertCalls)
		}
	})

	t.Run("synchronization failure - database upsert error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{"id": 401, "name": "db-fail-repo", "full_name": "Nachsyas/db-fail-repo"}
			]`))
		}))
		defer ts.Close()

		client := NewTestClient(ts.URL, "valid-token")
		repo := newMockGithubProjectRepo()
		repo.upsertFunc = func(ctx context.Context, p *domain.GithubProject) (bool, error) {
			return false, errors.New("postgres: connection pool timeout")
		}

		svc := NewSyncService(client, repo, "Nachsyas")
		_, err := svc.SyncRepositories(context.Background())
		if err == nil {
			t.Fatalf("expected error from failed DB upsert, got nil")
		}
	})

	t.Run("synchronization failure - missing token", func(t *testing.T) {
		client := NewTestClient("https://api.github.com", "")
		repo := newMockGithubProjectRepo()
		svc := NewSyncService(client, repo, "Nachsyas")

		_, err := svc.SyncRepositories(context.Background())
		if err == nil || !errors.Is(err, ErrMissingToken) {
			t.Errorf("expected ErrMissingToken, got %v", err)
		}
	})
}
