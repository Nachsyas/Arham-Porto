package http_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	delivery "github.com/nachsyas/arham-porto/apps/api/internal/delivery/http"
	"github.com/nachsyas/arham-porto/apps/api/internal/delivery/http/dto"
	"github.com/nachsyas/arham-porto/apps/api/internal/domain"
	"github.com/nachsyas/arham-porto/apps/api/internal/repository/postgres"
)

type mockGithubProjectsReaderRepo struct {
	projects []*domain.GithubProject
	err      error
	lastPage int
	lastLimit int
}

func (m *mockGithubProjectsReaderRepo) Create(ctx context.Context, p *domain.GithubProject) error {
	return nil
}
func (m *mockGithubProjectsReaderRepo) Update(ctx context.Context, p *domain.GithubProject) error {
	return nil
}
func (m *mockGithubProjectsReaderRepo) Upsert(ctx context.Context, p *domain.GithubProject) (bool, error) {
	return false, nil
}
func (m *mockGithubProjectsReaderRepo) List(ctx context.Context) ([]*domain.GithubProject, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.projects, nil
}
func (m *mockGithubProjectsReaderRepo) ListGithubProjects(ctx context.Context, page, limit int) ([]*domain.GithubProject, error) {
	m.lastPage = page
	m.lastLimit = limit
	if m.err != nil {
		return nil, m.err
	}
	return m.projects, nil
}
func (m *mockGithubProjectsReaderRepo) GetByGithubID(ctx context.Context, githubID int64) (*domain.GithubProject, error) {
	return nil, nil
}

func setupProjectsReaderTestRouter(repo domain.GithubProjectRepository) http.Handler {
	handler := delivery.NewHandler(nil, nil, nil, nil, nil, nil)
	if repo != nil {
		handler.SetGithubProjectRepository(repo)
	}
	generalLimiter := delivery.NewRateLimiter(120, time.Minute)
	return delivery.NewRouter(handler, []string{"*"}, generalLimiter)
}

func TestListGithubProjects_Success(t *testing.T) {
	syncedTime := time.Date(2026, 9, 14, 0, 0, 0, 0, time.UTC)
	mockProjects := []*domain.GithubProject{
		{
			ID:          "uuid-101",
			GithubID:    101,
			Name:        "EduTrace",
			FullName:    "Nachsyas/EduTrace",
			Description: "Educational platform",
			HTMLURL:     "https://github.com/Nachsyas/EduTrace",
			Homepage:    "https://edutrace.example.com",
			Language:    "TypeScript",
			Stars:       42,
			Forks:       8,
			Topics:      []string{"nextjs", "ai"},
			SyncedAt:    syncedTime,
		},
		{
			ID:          "uuid-102",
			GithubID:    102,
			Name:        "maritime-ai-dashboard",
			FullName:    "Nachsyas/maritime-ai-dashboard",
			Description: "Maritime monitoring system",
			HTMLURL:     "https://github.com/Nachsyas/maritime-ai-dashboard",
			Homepage:    "",
			Language:    "Python",
			Stars:       12,
			Forks:       3,
			Topics:      []string{"python", "fastapi"},
			SyncedAt:    syncedTime,
		},
	}

	mockRepo := &mockGithubProjectsReaderRepo{
		projects: mockProjects,
	}
	router := setupProjectsReaderTestRouter(mockRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/github/projects?page=1&limit=10", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Body: %s", rec.Code, rec.Body.String())
	}

	contentType := rec.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("expected application/json, got %s", contentType)
	}

	var res dto.GithubProjectsListResponse
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if res.Status != "success" {
		t.Errorf("expected status 'success', got %s", res.Status)
	}
	if len(res.Data) != 2 {
		t.Fatalf("expected 2 items, got %d", len(res.Data))
	}

	item := res.Data[0]
	if item.ID != "uuid-101" || item.Name != "EduTrace" || item.FullName != "Nachsyas/EduTrace" {
		t.Errorf("item 1 identity mismatch: %+v", item)
	}
	if item.Description != "Educational platform" || item.HTMLURL != "https://github.com/Nachsyas/EduTrace" {
		t.Errorf("item 1 urls mismatch: %+v", item)
	}
	if item.Language != "TypeScript" || item.Stars != 42 || item.Forks != 8 {
		t.Errorf("item 1 metrics mismatch: %+v", item)
	}
	if len(item.Topics) != 2 || item.Topics[0] != "nextjs" {
		t.Errorf("item 1 topics mismatch: %+v", item.Topics)
	}
	if item.SyncedAt != "2026-09-14T00:00:00Z" {
		t.Errorf("expected synced_at '2026-09-14T00:00:00Z', got %s", item.SyncedAt)
	}

	// Verify pagination parameters passed to repo
	if mockRepo.lastPage != 1 || mockRepo.lastLimit != 10 {
		t.Errorf("expected page=1 limit=10 passed to repo, got page=%d limit=%d", mockRepo.lastPage, mockRepo.lastLimit)
	}
}

func TestListGithubProjects_DefaultPagination(t *testing.T) {
	mockRepo := &mockGithubProjectsReaderRepo{
		projects: []*domain.GithubProject{},
	}
	router := setupProjectsReaderTestRouter(mockRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/github/projects", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if mockRepo.lastPage != 1 || mockRepo.lastLimit != 20 {
		t.Errorf("expected default page=1 limit=20, got page=%d limit=%d", mockRepo.lastPage, mockRepo.lastLimit)
	}
}

func TestListGithubProjects_EmptyDatabase(t *testing.T) {
	mockRepo := &mockGithubProjectsReaderRepo{
		projects: []*domain.GithubProject{},
	}
	router := setupProjectsReaderTestRouter(mockRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/github/projects", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var res dto.GithubProjectsListResponse
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if res.Status != "success" {
		t.Errorf("expected status 'success', got %s", res.Status)
	}
	if res.Data == nil || len(res.Data) != 0 {
		t.Errorf("expected empty data array [], got %+v", res.Data)
	}
}

func TestListGithubProjects_InvalidPagination(t *testing.T) {
	mockRepo := &mockGithubProjectsReaderRepo{}
	router := setupProjectsReaderTestRouter(mockRepo)

	testCases := []struct {
		name        string
		queryString string
		expectedMsg string
	}{
		{"page is zero", "?page=0", "invalid 'page' query parameter"},
		{"page is negative", "?page=-5", "invalid 'page' query parameter"},
		{"page is non-numeric", "?page=abc", "invalid 'page' query parameter"},
		{"repeated page param", "?page=1&page=2", "ambiguous repeated 'page' query parameter"},
		{"limit is zero", "?limit=0", "invalid 'limit' query parameter"},
		{"limit is negative", "?limit=-1", "invalid 'limit' query parameter"},
		{"limit exceeds max 100", "?limit=101", "invalid 'limit' query parameter"},
		{"limit is non-numeric", "?limit=xyz", "invalid 'limit' query parameter"},
		{"repeated limit param", "?limit=10&limit=20", "ambiguous repeated 'limit' query parameter"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/github/projects"+tc.queryString, nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("expected status 400, got %d for %s", rec.Code, tc.queryString)
			}

			var errResp dto.ErrorEnvelope
			if err := json.NewDecoder(rec.Body).Decode(&errResp); err != nil {
				t.Fatalf("failed to decode error: %v", err)
			}
			if errResp.Error.Code != "bad_request" {
				t.Errorf("expected code 'bad_request', got %s", errResp.Error.Code)
			}
			if !strings.Contains(errResp.Error.Message, tc.expectedMsg) {
				t.Errorf("expected message to contain %q, got %q", tc.expectedMsg, errResp.Error.Message)
			}
		})
	}
}

func TestListGithubProjects_DatabaseError(t *testing.T) {
	t.Run("internal database query error returns 500", func(t *testing.T) {
		mockRepo := &mockGithubProjectsReaderRepo{
			err: errors.New("sql: syntax or connection failure"),
		}
		router := setupProjectsReaderTestRouter(mockRepo)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/github/projects", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("expected status 500, got %d", rec.Code)
		}

		var errResp dto.ErrorEnvelope
		if err := json.NewDecoder(rec.Body).Decode(&errResp); err != nil {
			t.Fatalf("failed to decode error: %v", err)
		}
		if errResp.Error.Code != "internal_error" {
			t.Errorf("expected code 'internal_error', got %s", errResp.Error.Code)
		}
	})

	t.Run("database disabled returns 503", func(t *testing.T) {
		mockRepo := &mockGithubProjectsReaderRepo{
			err: postgres.ErrDatabaseDisabled,
		}
		router := setupProjectsReaderTestRouter(mockRepo)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/github/projects", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusServiceUnavailable {
			t.Fatalf("expected status 503, got %d", rec.Code)
		}
	})

	t.Run("unconfigured repo layer returns 503", func(t *testing.T) {
		router := setupProjectsReaderTestRouter(nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/github/projects", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusServiceUnavailable {
			t.Fatalf("expected status 503, got %d", rec.Code)
		}
	})
}

func TestListGithubProjects_MethodNotAllowed(t *testing.T) {
	router := setupProjectsReaderTestRouter(&mockGithubProjectsReaderRepo{})

	disallowedMethods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}
	for _, method := range disallowedMethods {
		req := httptest.NewRequest(method, "/api/v1/github/projects", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected 405 for %s, got %d", method, rec.Code)
		}
		if !strings.Contains(rec.Header().Get("Content-Type"), "application/json") {
			t.Errorf("expected json error response for %s, got %s", method, rec.Header().Get("Content-Type"))
		}
	}
}
