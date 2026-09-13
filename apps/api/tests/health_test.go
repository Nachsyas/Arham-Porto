package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	delivery "github.com/nachsyas/arham-porto/apps/api/internal/delivery/http"
	"github.com/nachsyas/arham-porto/apps/api/internal/usecase"
)

type mockDBChecker struct {
	status string
}

func (m *mockDBChecker) Status(ctx context.Context) string {
	return m.status
}

func setupTestRouter(dbStatus string) http.Handler {
	profileUC := usecase.NewProfileUseCase(nil)
	projectUC := usecase.NewProjectUseCase(nil)
	skillUC := usecase.NewSkillUseCase(nil)
	evidenceUC := usecase.NewEvidenceUseCase(nil)
	journeyUC := usecase.NewJourneyUseCase(nil)
	checker := &mockDBChecker{status: dbStatus}

	h := delivery.NewHandler(profileUC, projectUC, skillUC, evidenceUC, journeyUC, checker)
	return delivery.NewRouter(h, []string{"*"}, nil)
}

func TestHealthzEndpoint(t *testing.T) {
	router := setupTestRouter("disabled")

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if cc := rec.Header().Get("Cache-Control"); cc != "no-store" {
		t.Errorf("expected Cache-Control 'no-store', got '%s'", cc)
	}

	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["status"] != "ok" || resp["scope"] != "process_alive" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestReadyzEndpoint_Disabled(t *testing.T) {
	router := setupTestRouter("disabled")

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["status"] != "ready" || resp["database"] != "disabled" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestReadyzEndpoint_Connected(t *testing.T) {
	router := setupTestRouter("connected")

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["status"] != "ready" || resp["database"] != "connected" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestReadyzEndpoint_Degraded(t *testing.T) {
	router := setupTestRouter("degraded")

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["status"] != "ready" || resp["database"] != "degraded" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestReadyzEndpoint_Unavailable(t *testing.T) {
	router := setupTestRouter("unavailable")

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status 503, got %d", rec.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["status"] != "not_ready" || resp["database"] != "unavailable" {
		t.Errorf("unexpected response: %+v", resp)
	}
}
