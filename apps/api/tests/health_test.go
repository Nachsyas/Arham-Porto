package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	delivery "github.com/nachsyas/arham-porto/apps/api/internal/delivery/http"
	"github.com/nachsyas/arham-porto/apps/api/internal/usecase"
)

func TestHealthzEndpoint(t *testing.T) {
	profileUC := usecase.NewProfileUseCase(nil)
	projectUC := usecase.NewProjectUseCase(nil)
	evidenceUC := usecase.NewEvidenceUseCase(nil)
	h := delivery.NewHandler(profileUC, projectUC, evidenceUC)
	router := delivery.NewRouter(h, "*")

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["status"] != "ok" {
		t.Errorf("expected status 'ok', got '%s'", resp["status"])
	}
}

func TestReadyzEndpoint(t *testing.T) {
	profileUC := usecase.NewProfileUseCase(nil)
	projectUC := usecase.NewProjectUseCase(nil)
	evidenceUC := usecase.NewEvidenceUseCase(nil)
	h := delivery.NewHandler(profileUC, projectUC, evidenceUC)
	router := delivery.NewRouter(h, "*")

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

	if resp["status"] != "ready" {
		t.Errorf("expected status 'ready', got '%s'", resp["status"])
	}
}
