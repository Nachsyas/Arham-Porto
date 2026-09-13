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

	bodyStr := rec.Body.String()
	t.Logf("Profile serialized body: %s", bodyStr)

	var resp dto.DataEnvelope[dto.ProfileResponse]
	if err := json.NewDecoder(strings.NewReader(bodyStr)).Decode(&resp); err != nil {
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

	bodyStr := rec.Body.String()
	t.Logf("Projects serialized body: %s", bodyStr)

	var resp dto.ListEnvelope[dto.ProjectResponse]
	if err := json.NewDecoder(strings.NewReader(bodyStr)).Decode(&resp); err != nil {
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

	// Verify absence of journey-tk and verify public field is not exposed in public DTO
	if strings.Contains(rec.Body.String(), `"public":`) {
		t.Error("redundant internal publication-control flag 'public' must not appear in public journey transport DTO")
	}
	for _, stop := range resp.Data {
		if stop.ID == "journey-tk" {
			t.Error("private stop journey-tk must not be present in public journey endpoint")
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

func TestEndpoints_RateLimiter_PortNormalization(t *testing.T) {
	// Limiter allows 1 request per minute
	rl := delivery.NewRateLimiter(1, time.Minute)
	defer rl.Close()

	router := setupFullTestRouter(t, []string{"*"}, rl)

	// First request from 192.0.2.10:51432 (allowed)
	req1 := httptest.NewRequest(http.MethodGet, "/api/v1/projects", nil)
	req1.RemoteAddr = "192.0.2.10:51432"
	rec1 := httptest.NewRecorder()
	router.ServeHTTP(rec1, req1)
	if rec1.Code != http.StatusOK {
		t.Fatalf("first request expected status 200, got %d", rec1.Code)
	}

	// Second request from SAME IP but DIFFERENT source port 192.0.2.10:51433
	// Must share the same rate-limit bucket and be rejected with 429
	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/projects", nil)
	req2.RemoteAddr = "192.0.2.10:51433"
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusTooManyRequests {
		t.Fatalf("second request from same IP with different port must be throttled with 429, got %d", rec2.Code)
	}

	// Different IP 192.0.2.20:51432 must have its own bucket (allowed)
	req3 := httptest.NewRequest(http.MethodGet, "/api/v1/projects", nil)
	req3.RemoteAddr = "192.0.2.20:51432"
	rec3 := httptest.NewRecorder()
	router.ServeHTTP(rec3, req3)
	if rec3.Code != http.StatusOK {
		t.Fatalf("request from different IP expected status 200, got %d", rec3.Code)
	}
}

func TestEndpoints_MethodNotAllowed(t *testing.T) {
	router := setupFullTestRouter(t, []string{"*"}, nil)

	// POST on GET-only endpoint
	req := httptest.NewRequest(http.MethodPost, "/api/v1/projects", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405 Method Not Allowed for POST /api/v1/projects, got %d", rec.Code)
	}

	// PUT on GET-only endpoint
	req = httptest.NewRequest(http.MethodPut, "/api/v1/profile", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405 Method Not Allowed for PUT /api/v1/profile, got %d", rec.Code)
	}
}

func TestEndpoints_QueryParameterStrictness(t *testing.T) {
	router := setupFullTestRouter(t, []string{"*"}, nil)

	// Valid boolean
	req := httptest.NewRequest(http.MethodGet, "/api/v1/projects?featured=true", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 for featured=true, got %d", rec.Code)
	}

	// Invalid boolean value
	req = httptest.NewRequest(http.MethodGet, "/api/v1/projects?featured=banana", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for featured=banana, got %d", rec.Code)
	}

	// Repeated ambiguous parameter: ?featured=true&featured=false
	req = httptest.NewRequest(http.MethodGet, "/api/v1/projects?featured=true&featured=false", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for repeated featured param, got %d", rec.Code)
	}

	// Repeated ambiguous category: ?category=AI&category=Backend
	req = httptest.NewRequest(http.MethodGet, "/api/v1/projects?category=AI&category=Backend", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for repeated category param on projects, got %d", rec.Code)
	}

	// Repeated ambiguous category on skills: ?category=Backend&category=Frontend
	req = httptest.NewRequest(http.MethodGet, "/api/v1/skills?category=Backend&category=Frontend", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for repeated category param on skills, got %d", rec.Code)
	}
}

func TestEndpoints_SecurityHeaders_OnErrors(t *testing.T) {
	router := setupFullTestRouter(t, []string{"*"}, nil)

	assertSecurityHeaders := func(t *testing.T, rec *httptest.ResponseRecorder, context string) {
		t.Helper()
		if nosniff := rec.Header().Get("X-Content-Type-Options"); nosniff != "nosniff" {
			t.Errorf("[%s] expected X-Content-Type-Options 'nosniff', got '%s'", context, nosniff)
		}
		if ref := rec.Header().Get("Referrer-Policy"); ref != "no-referrer" {
			t.Errorf("[%s] expected Referrer-Policy 'no-referrer', got '%s'", context, ref)
		}
		if csp := rec.Header().Get("Content-Security-Policy"); !strings.Contains(csp, "default-src 'none'") {
			t.Errorf("[%s] expected CSP default-src none, got '%s'", context, csp)
		}
		if xfo := rec.Header().Get("X-Frame-Options"); xfo != "DENY" {
			t.Errorf("[%s] expected X-Frame-Options 'DENY', got '%s'", context, xfo)
		}
	}

	// 400 Bad Request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/projects?featured=invalid", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
	assertSecurityHeaders(t, rec, "400 Bad Request")
	if cc := rec.Header().Get("Cache-Control"); cc != "no-store" {
		t.Errorf("expected Cache-Control 'no-store' on 400, got '%s'", cc)
	}

	// 404 Not Found (valid route, missing resource)
	req = httptest.NewRequest(http.MethodGet, "/api/v1/projects/non-existent-slug-xyz", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
	assertSecurityHeaders(t, rec, "404 Not Found")
	if cc := rec.Header().Get("Cache-Control"); cc != "no-store" {
		t.Errorf("expected Cache-Control 'no-store' on 404, got '%s'", cc)
	}

	// 500 Panic Recovery
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("simulated disaster")
	})
	recovered := delivery.RecoveryMiddleware(panicHandler)
	req = httptest.NewRequest(http.MethodGet, "/panic", nil)
	rec = httptest.NewRecorder()
	recovered.ServeHTTP(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
	assertSecurityHeaders(t, rec, "500 Internal Error")
	if cc := rec.Header().Get("Cache-Control"); cc != "no-store" {
		t.Errorf("expected Cache-Control 'no-store' on 500, got '%s'", cc)
	}
}

func TestEndpoints_ContentTypeCharset(t *testing.T) {
	router := setupFullTestRouter(t, []string{"*"}, nil)

	routes := []string{
		"/api/v1/profile",
		"/api/v1/projects",
		"/api/v1/skills",
		"/api/v1/evidence",
		"/api/v1/journey",
	}

	for _, route := range routes {
		req := httptest.NewRequest(http.MethodGet, route, nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 for %s, got %d", route, rec.Code)
		}

		ct := rec.Header().Get("Content-Type")
		if !strings.Contains(ct, "application/json") || !strings.Contains(ct, "charset=utf-8") {
			t.Errorf("expected 'application/json; charset=utf-8' on %s, got '%s'", route, ct)
		}
	}
}

func TestDTO_Sanitizer_Security(t *testing.T) {
	// 1. Source Path Security
	safeRel := "apps/api/cmd/server/main.go"
	if res := dto.SanitizeSourcePath(&safeRel); res == nil || *res != safeRel {
		t.Errorf("expected safe relative path to be preserved, got %v", res)
	}

	unsafePaths := []string{
		"/Users/user/Documents/secret.txt",
		"/home/runner/work/arham-porto/secret.key",
		"C:\\Windows\\system32\\cmd.exe",
		"../../etc/passwd",
		"../secret.json",
		"/etc/shadow",
		"users/private/key.pem",
	}
	for _, p := range unsafePaths {
		pCopy := p
		if res := dto.SanitizeSourcePath(&pCopy); res != nil {
			t.Errorf("SECURITY FAILURE: Unsafe source path %q was not rejected, got %s", p, *res)
		}
	}

	// 2. URL Output Security
	safeURL := "https://github.com/Nachsyas/arham-porto"
	if res := dto.SanitizeExternalURL(&safeURL); res == nil || *res != safeURL {
		t.Errorf("expected safe HTTPS URL to be preserved, got %v", res)
	}

	unsafeURLs := []string{
		"javascript:alert(document.cookie)",
		"file:///etc/passwd",
		"data:text/html,<script>alert(1)</script>",
		"vbscript:msgbox(1)",
		"blob:https://example.com/xyz",
		"http://insecure-site.com",
		"TODO_USER_URL",
	}
	for _, u := range unsafeURLs {
		uCopy := u
		if res := dto.SanitizeExternalURL(&uCopy); res != nil {
			t.Errorf("SECURITY FAILURE: Unsafe external URL %q was not rejected, got %s", u, *res)
		}
	}

	// 3. Image URL Security
	safeImageRel := "/images/hero.webp"
	if res := dto.SanitizeImageURL(&safeImageRel); res == nil || *res != safeImageRel {
		t.Errorf("expected safe relative image URL to be preserved, got %v", res)
	}

	unsafeImages := []string{
		"javascript:alert(1)",
		"file:///root/image.png",
		"data:image/svg+xml;base64,...",
	}
	for _, img := range unsafeImages {
		imgCopy := img
		if res := dto.SanitizeImageURL(&imgCopy); res != nil {
			t.Errorf("SECURITY FAILURE: Unsafe image URL %q was not rejected, got %s", img, *res)
		}
	}
}

func TestEndpoints_ProfileApprovedFieldsOnly(t *testing.T) {
	router := setupFullTestRouter(t, []string{"*"}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var raw map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&raw); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	data, ok := raw["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected data object in response, got %v", raw)
	}

	// Approved fields that MUST exist
	approvedKeys := []string{"full_name", "role", "project_name", "ai_feature", "github"}
	for _, k := range approvedKeys {
		val, exists := data[k]
		if !exists || val == nil || val == "" {
			t.Errorf("expected approved field '%s' to be populated, got %v", k, val)
		}
	}

	// Unapproved, null, or TODO-backed fields that MUST NOT exist in public payload
	forbiddenKeys := []string{"positioning", "bio", "current_city", "availability", "email", "linkedin", "cvUrl", "todo"}
	for _, k := range forbiddenKeys {
		if val, exists := data[k]; exists && val != nil {
			t.Errorf("SECURITY/PRIVACY VIOLATION: Unapproved field '%s' is present in public profile response with value: %v", k, val)
		}
	}
}

func TestEndpoints_DumpPublicSnapshots(t *testing.T) {
	router := setupFullTestRouter(t, []string{"*"}, nil)

	endpoints := []string{
		"/api/v1/profile",
		"/api/v1/projects",
		"/api/v1/skills",
		"/api/v1/evidence",
		"/api/v1/journey",
	}

	for _, ep := range endpoints {
		req := httptest.NewRequest(http.MethodGet, ep, nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("endpoint %s returned status %d", ep, rec.Code)
		}

		t.Logf("=== SNAPSHOT [%s] ===\n%s\n", ep, rec.Body.String())
	}
}
