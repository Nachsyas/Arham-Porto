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
	"github.com/nachsyas/arham-porto/apps/api/internal/github"
	"github.com/nachsyas/arham-porto/apps/api/internal/repository/postgres"
)

type mockSyncService struct {
	result *github.SyncResult
	err    error
}

func (m *mockSyncService) SyncRepositories(ctx context.Context) (*github.SyncResult, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.result, nil
}

func setupSyncTestRouter(svc delivery.GitHubSyncService) http.Handler {
	handler := delivery.NewHandler(nil, nil, nil, nil, nil, nil)
	if svc != nil {
		handler.SetGitHubSyncService(svc)
	}
	generalLimiter := delivery.NewRateLimiter(120, time.Minute)
	return delivery.NewRouter(handler, []string{"*"}, generalLimiter)
}

func TestSyncGitHubRepositories_Success(t *testing.T) {
	mockSvc := &mockSyncService{
		result: &github.SyncResult{
			Synced:  25,
			Created: 5,
			Updated: 20,
		},
	}
	router := setupSyncTestRouter(mockSvc)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/github/sync", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Body: %s", rec.Code, rec.Body.String())
	}

	contentType := rec.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("expected application/json content-type, got %s", contentType)
	}

	var res map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}

	if res["status"] != "success" {
		t.Errorf("expected status 'success', got %v", res["status"])
	}
	if synced, ok := res["synced"].(float64); !ok || int(synced) != 25 {
		t.Errorf("expected synced=25, got %v", res["synced"])
	}
	if created, ok := res["created"].(float64); !ok || int(created) != 5 {
		t.Errorf("expected created=5, got %v", res["created"])
	}
	if updated, ok := res["updated"].(float64); !ok || int(updated) != 20 {
		t.Errorf("expected updated=20, got %v", res["updated"])
	}
}

func TestSyncGitHubRepositories_Errors(t *testing.T) {
	t.Run("missing GITHUB_TOKEN returns 401 unauthorized", func(t *testing.T) {
		mockSvc := &mockSyncService{
			err: github.ErrMissingToken,
		}
		router := setupSyncTestRouter(mockSvc)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/github/sync", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected status 401, got %d", rec.Code)
		}

		var errEnvelope dto.ErrorEnvelope
		if err := json.NewDecoder(rec.Body).Decode(&errEnvelope); err != nil {
			t.Fatalf("failed to decode error envelope: %v", err)
		}

		if errEnvelope.Error.Code != "unauthorized" {
			t.Errorf("expected error code 'unauthorized', got %q", errEnvelope.Error.Code)
		}
	})

	t.Run("rate limit exceeded returns 429 rate_limited", func(t *testing.T) {
		mockSvc := &mockSyncService{
			err: github.ErrRateLimited,
		}
		router := setupSyncTestRouter(mockSvc)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/github/sync", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusTooManyRequests {
			t.Fatalf("expected status 429, got %d", rec.Code)
		}

		var errEnvelope dto.ErrorEnvelope
		if err := json.NewDecoder(rec.Body).Decode(&errEnvelope); err != nil {
			t.Fatalf("failed to decode error envelope: %v", err)
		}

		if errEnvelope.Error.Code != "rate_limited" {
			t.Errorf("expected error code 'rate_limited', got %q", errEnvelope.Error.Code)
		}
	})

	t.Run("github API failure returns 502 bad_gateway", func(t *testing.T) {
		mockSvc := &mockSyncService{
			err: errors.New("github fetch repos failed: upstream connection reset"),
		}
		router := setupSyncTestRouter(mockSvc)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/github/sync", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadGateway {
			t.Fatalf("expected status 502, got %d", rec.Code)
		}

		var errEnvelope dto.ErrorEnvelope
		if err := json.NewDecoder(rec.Body).Decode(&errEnvelope); err != nil {
			t.Fatalf("failed to decode error envelope: %v", err)
		}

		if errEnvelope.Error.Code != "bad_gateway" {
			t.Errorf("expected error code 'bad_gateway', got %q", errEnvelope.Error.Code)
		}
	})

	t.Run("database failure returns 503 service_unavailable", func(t *testing.T) {
		mockSvc := &mockSyncService{
			err: postgres.ErrDatabaseDisabled,
		}
		router := setupSyncTestRouter(mockSvc)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/github/sync", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusServiceUnavailable {
			t.Fatalf("expected status 503, got %d", rec.Code)
		}

		var errEnvelope dto.ErrorEnvelope
		if err := json.NewDecoder(rec.Body).Decode(&errEnvelope); err != nil {
			t.Fatalf("failed to decode error envelope: %v", err)
		}

		if errEnvelope.Error.Code != "service_unavailable" {
			t.Errorf("expected error code 'service_unavailable', got %q", errEnvelope.Error.Code)
		}
	})

	t.Run("unconfigured sync service returns 503 service_unavailable", func(t *testing.T) {
		router := setupSyncTestRouter(nil)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/github/sync", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusServiceUnavailable {
			t.Fatalf("expected status 503, got %d", rec.Code)
		}
	})
}

func TestSyncGitHubRepositories_MethodNotAllowed(t *testing.T) {
	router := setupSyncTestRouter(&mockSyncService{})

	disallowedMethods := []string{http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodPatch}
	for _, method := range disallowedMethods {
		req := httptest.NewRequest(method, "/api/v1/github/sync", nil)
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
