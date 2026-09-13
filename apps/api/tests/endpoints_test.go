package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	delivery "github.com/nachsyas/arham-porto/apps/api/internal/delivery/http"
	"github.com/nachsyas/arham-porto/apps/api/internal/delivery/http/dto"
	"github.com/nachsyas/arham-porto/apps/api/internal/repository/jsonfile"
	"github.com/nachsyas/arham-porto/apps/api/internal/usecase"
)

func setupFullTestRouter(t *testing.T, allowedOrigins []string, rl *delivery.RateLimiter) http.Handler {
	t.Helper()
	dataDir := getCanonicalDataDir(t)
	repo, err := jsonfile.LoadRepository(dataDir)
	if err != nil {
		t.Fatalf("failed to load canonical repository: %v", err)
	}

	profileUC := usecase.NewProfileUseCase(repo)
	projectUC := usecase.NewProjectUseCase(repo)
	skillUC := usecase.NewSkillUseCase(repo)
	evidenceUC := usecase.NewEvidenceUseCase(repo)
	journeyUC := usecase.NewJourneyUseCase(repo)
	checker := &mockDBChecker{status: "connected"}

	h := delivery.NewHandler(profileUC, projectUC, skillUC, evidenceUC, journeyUC, checker)
	return delivery.NewRouter(h, allowedOrigins, rl)
}

func TestEndpoints_GetProfile(t *testing.T) {
	router := setupFullTestRouter(t, []string{"http://localhost:3000"}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if cc := rec.Header().Get("Cache-Control"); !strings.Contains(cc, "public") {
		t.Errorf("expected public Cache-Control header, got '%s'", cc)
	}

	var resp dto.DataEnvelope[dto.ProfileResponse]
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.FullName != "Nachsyas Arham Mumtaz Nashohi" {
		t.Errorf("expected FullName 'Nachsyas Arham Mumtaz Nashohi', got '%s'", resp.Data.FullName)
	}
	if resp.Data.Role != "Software Engineer" {
		t.Errorf("expected Role 'Software Engineer', got '%s'", resp.Data.Role)
	}
}

func TestEndpoints_PrivacyLeakageAssertion(t *testing.T) {
	router := setupFullTestRouter(t, []string{"http://localhost:3000"}, nil)

	endpoints := []string{
		"/api/v1/profile",
		"/api/v1/projects",
		"/api/v1/skills",
		"/api/v1/evidence",
		"/api/v1/journey",
	}

	forbiddenSubstrings := []string{
		"TODO_USER",
		"TODO_",
		"coordinates",
		"phone",
		"birth_date",
		"birthDate",
		"home_address",
		"homeAddress",
		"residential",
	}

	for _, ep := range endpoints {
		req := httptest.NewRequest(http.MethodGet, ep, nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("endpoint %s returned status %d", ep, rec.Code)
		}

		bodyStr := rec.Body.String()
		for _, forbidden := range forbiddenSubstrings {
			if strings.Contains(bodyStr, forbidden) {
				t.Errorf("SECURITY/PRIVACY VIOLATION: Endpoint %s contains forbidden token '%s'. Response body: %s",
					ep, forbidden, bodyStr)
			}
		}
	}
}

func TestEndpoints_ListProjects(t *testing.T) {
	router := setupFullTestRouter(t, []string{"http://localhost:3000"}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/projects", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp dto.ListEnvelope[dto.ProjectResponse]
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Meta.Count != 4 || len(resp.Data) != 4 {
		t.Errorf("expected count 4, got %d", resp.Meta.Count)
	}
}

func TestEndpoints_ListProjects_Filtering(t *testing.T) {
	router := setupFullTestRouter(t, []string{"http://localhost:3000"}, nil)

	// Valid filter: category=Backend
	req := httptest.NewRequest(http.MethodGet, "/api/v1/projects?category=Backend", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	var resp dto.ListEnvelope[dto.ProjectResponse]
	_ = json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Meta.Count != 1 || resp.Data[0].Slug != "gdgoc-ecommerce" {
		t.Errorf("expected 1 project (gdgoc-ecommerce), got %d", resp.Meta.Count)
	}

	// Valid filter: featured=true
	req = httptest.NewRequest(http.MethodGet, "/api/v1/projects?featured=true", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	// Invalid filter: category=Invalid
	req = httptest.NewRequest(http.MethodGet, "/api/v1/projects?category=Invalid", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
	var errResp dto.ErrorEnvelope
	_ = json.NewDecoder(rec.Body).Decode(&errResp)
	if errResp.Error.Code != "bad_request" {
		t.Errorf("expected error code 'bad_request', got '%s'", errResp.Error.Code)
	}

	// Invalid filter: featured=banana
	req = httptest.NewRequest(http.MethodGet, "/api/v1/projects?featured=banana", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for featured=banana, got %d", rec.Code)
	}
}

func TestEndpoints_GetProjectBySlug(t *testing.T) {
	router := setupFullTestRouter(t, []string{"http://localhost:3000"}, nil)

	// Found: edutrace
	req := httptest.NewRequest(http.MethodGet, "/api/v1/projects/edutrace", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	var resp dto.DataEnvelope[dto.ProjectResponse]
	_ = json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Data.Slug != "edutrace" {
		t.Errorf("expected slug 'edutrace', got '%s'", resp.Data.Slug)
	}

	// Not Found
	req = httptest.NewRequest(http.MethodGet, "/api/v1/projects/nonexistent-slug", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
	var errResp dto.ErrorEnvelope
	_ = json.NewDecoder(rec.Body).Decode(&errResp)
	if errResp.Error.Code != "not_found" {
		t.Errorf("expected code 'not_found', got '%s'", errResp.Error.Code)
	}

	// Bad Request (invalid slug characters)
	req = httptest.NewRequest(http.MethodGet, "/api/v1/projects/INVALID_SLUG!!", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for invalid slug format, got %d", rec.Code)
	}
}

func TestEndpoints_ListSkills(t *testing.T) {
	router := setupFullTestRouter(t, []string{"http://localhost:3000"}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/skills", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	var resp dto.ListEnvelope[dto.SkillResponse]
	_ = json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Meta.Count == 0 {
		t.Error("expected non-empty skills list")
	}

	// Valid category filter
	req = httptest.NewRequest(http.MethodGet, "/api/v1/skills?category=Backend", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	// Invalid category filter
	req = httptest.NewRequest(http.MethodGet, "/api/v1/skills?category=InvalidCat", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestEndpoints_ListEvidence(t *testing.T) {
	router := setupFullTestRouter(t, []string{"http://localhost:3000"}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/evidence", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	var resp dto.ListEnvelope[dto.EvidenceResponse]
	_ = json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Meta.Count != 4 {
		t.Errorf("expected 4 evidence items, got %d", resp.Meta.Count)
	}

	// Filter by skill_id
	req = httptest.NewRequest(http.MethodGet, "/api/v1/evidence?skill_id=backend-go", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	// Get by ID: found
	req = httptest.NewRequest(http.MethodGet, "/api/v1/evidence/evidence-edutrace-repo", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	// Get by ID: not found
	req = httptest.NewRequest(http.MethodGet, "/api/v1/evidence/unknown-id", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
}

func TestEndpoints_ListJourney(t *testing.T) {
	router := setupFullTestRouter(t, []string{"http://localhost:3000"}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/journey", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp dto.ListEnvelope[dto.JourneyStopResponse]
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Strictly 7 public milestones
	if resp.Meta.Count != 7 || len(resp.Data) != 7 {
		t.Fatalf("expected exactly 7 public milestones, got %d", resp.Meta.Count)
	}

	// Verify absence of journey-tk
	for _, stop := range resp.Data {
		if stop.ID == "journey-tk" {
			t.Error("private stop journey-tk must not be present in public journey endpoint")
		}
		if !stop.Public {
			t.Errorf("stop %s is not marked public", stop.ID)
		}
	}
}

func TestEndpoints_CORS(t *testing.T) {
	allowed := []string{"http://localhost:3000", "https://portfolio.example.com"}
	router := setupFullTestRouter(t, allowed, nil)

	// Approved origin
	req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if origin := rec.Header().Get("Access-Control-Allow-Origin"); origin != "http://localhost:3000" {
		t.Errorf("expected Access-Control-Allow-Origin 'http://localhost:3000', got '%s'", origin)
	}
	if vary := rec.Header().Get("Vary"); !strings.Contains(vary, "Origin") {
		t.Errorf("expected Vary header to contain 'Origin', got '%s'", vary)
	}

	// Unapproved origin
	req = httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
	req.Header.Set("Origin", "http://unapproved.com")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if origin := rec.Header().Get("Access-Control-Allow-Origin"); origin != "" {
		t.Errorf("expected empty Access-Control-Allow-Origin for unapproved origin, got '%s'", origin)
	}

	// OPTIONS preflight
	req = httptest.NewRequest(http.MethodOptions, "/api/v1/projects", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204 No Content for OPTIONS, got %d", rec.Code)
	}
}

func TestEndpoints_SecurityHeaders(t *testing.T) {
	router := setupFullTestRouter(t, []string{"*"}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if nosniff := rec.Header().Get("X-Content-Type-Options"); nosniff != "nosniff" {
		t.Errorf("expected X-Content-Type-Options 'nosniff', got '%s'", nosniff)
	}
	if ref := rec.Header().Get("Referrer-Policy"); ref != "no-referrer" {
		t.Errorf("expected Referrer-Policy 'no-referrer', got '%s'", ref)
	}
	if csp := rec.Header().Get("Content-Security-Policy"); !strings.Contains(csp, "default-src 'none'") {
		t.Errorf("expected CSP 'default-src none', got '%s'", csp)
	}
	if xfo := rec.Header().Get("X-Frame-Options"); xfo != "DENY" {
		t.Errorf("expected X-Frame-Options 'DENY', got '%s'", xfo)
	}
	// Obsolete X-XSS-Protection must NOT be set
	if xxss := rec.Header().Get("X-XSS-Protection"); xxss != "" {
		t.Errorf("obsolete X-XSS-Protection header should not be set, got '%s'", xxss)
	}
}

func TestEndpoints_RateLimiting(t *testing.T) {
	// Limiter allowing 3 requests per minute
	rl := delivery.NewRateLimiter(3, time.Minute)
	defer rl.Close()

	router := setupFullTestRouter(t, []string{"*"}, rl)

	// Fire 3 allowed requests
	for i := 1; i <= 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/projects", nil)
		req.RemoteAddr = "192.168.1.100:12345"
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("request %d expected status 200, got %d", i, rec.Code)
		}
	}

	// 4th request must be throttled with 429
	req := httptest.NewRequest(http.MethodGet, "/api/v1/projects", nil)
	req.RemoteAddr = "192.168.1.100:12345"
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("expected status 429 for throttled request, got %d", rec.Code)
	}

	var errResp dto.ErrorEnvelope
	_ = json.NewDecoder(rec.Body).Decode(&errResp)
	if errResp.Error.Code != "rate_limited" {
		t.Errorf("expected code 'rate_limited', got '%s'", errResp.Error.Code)
	}

	// Critical check: Probes must bypass rate limiting even when client is throttled!
	healthReq := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	healthReq.RemoteAddr = "192.168.1.100:12345"
	healthRec := httptest.NewRecorder()
	router.ServeHTTP(healthRec, healthReq)
	if healthRec.Code != http.StatusOK {
		t.Fatalf("expected /healthz to bypass rate limiting with 200, got %d", healthRec.Code)
	}
}

func TestEndpoints_PanicRecovery(t *testing.T) {
	// Create an inner handler that panics
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("unexpected critical failure simulation")
	})

	recovered := delivery.RecoveryMiddleware(panicHandler)

	req := httptest.NewRequest(http.MethodGet, "/simulate-panic", nil)
	rec := httptest.NewRecorder()

	recovered.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}

	var errResp dto.ErrorEnvelope
	if err := json.NewDecoder(rec.Body).Decode(&errResp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if errResp.Error.Code != "internal_error" {
		t.Errorf("expected error code 'internal_error', got '%s'", errResp.Error.Code)
	}

	// Ensure stack trace is NOT exposed to client
	if strings.Contains(rec.Body.String(), "panic") || strings.Contains(rec.Body.String(), "goroutine") {
		t.Error("error response must never expose internal stack trace or panic details")
	}
}
