package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	delivery "github.com/nachsyas/arham-porto/apps/api/internal/delivery/http"
	"github.com/nachsyas/arham-porto/apps/api/internal/domain"
	"github.com/nachsyas/arham-porto/apps/api/internal/llm"
	"github.com/nachsyas/arham-porto/apps/api/internal/retrieval"
	"github.com/nachsyas/arham-porto/apps/api/internal/usecase"
)

// mockRetriever for HTTP handler tests
type mockHTTPRetriever struct {
	results []domain.RetrievalResult
	err     error
	delay   time.Duration
}

func (m *mockHTTPRetriever) Retrieve(ctx context.Context, query string, opts retrieval.SearchOptions) ([]domain.RetrievalResult, error) {
	if m.delay > 0 {
		select {
		case <-time.After(m.delay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	if m.err != nil {
		return nil, m.err
	}
	return m.results, nil
}

// setupTestServer wires a test server with configurable AI settings
func setupTestServer(
	aiMode string,
	retriever usecase.Retriever,
	fakeLLM llm.LLMProvider,
	maxConcurrent int,
	rateLimit int,
) http.Handler {
	handler := delivery.NewHandler(nil, nil, nil, nil, nil, nil)
	var askUC usecase.AskUseCase
	if retriever != nil && fakeLLM != nil {
		askUC = usecase.NewAskUseCase(retriever, fakeLLM, 24000)
	}

	aiLimiter := delivery.NewRateLimiter(rateLimit, time.Minute)
	handler.EnableAI(askUC, aiMode, aiLimiter, maxConcurrent, 30)

	generalLimiter := delivery.NewRateLimiter(120, time.Minute)
	return delivery.NewRouter(handler, []string{"*"}, generalLimiter)
}

func TestAskHandler_MethodNotAllowed(t *testing.T) {
	router := setupTestServer("disabled", nil, nil, 4, 5)

	methods := []string{http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodPatch}
	for _, m := range methods {
		req := httptest.NewRequest(m, "/api/v1/ai/ask", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected 405 for %s, got %d", m, rec.Code)
		}
		if !strings.Contains(rec.Header().Get("Content-Type"), "application/json") {
			t.Errorf("expected json error envelope for 405, got %s", rec.Header().Get("Content-Type"))
		}
	}
}

func TestAskHandler_BadRequest(t *testing.T) {
	fakeLLM := llm.NewDeterministicFakeProvider()
	retriever := &mockHTTPRetriever{}
	router := setupTestServer("remote", retriever, fakeLLM, 4, 10)

	cases := []struct {
		name string
		body string
	}{
		{"invalid json", `{"question": `},
		{"empty question", `{"question": ""}`},
		{"whitespace question", `{"question": "    "}`},
		{"single char", `{"question": "a"}`},
		{"oversized question", `{"question": "` + strings.Repeat("x", 1001) + `"}`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/ai/ask", bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("expected 400 bad request, got %d", rec.Code)
			}
			if rec.Header().Get("Content-Type") != "application/json; charset=utf-8" {
				t.Fatalf("expected json content type, got %s", rec.Header().Get("Content-Type"))
			}
		})
	}
}

func TestAskHandler_AIDisabled(t *testing.T) {
	router := setupTestServer("disabled", nil, nil, 4, 5)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/ai/ask", bytes.NewBufferString(`{"question": "What is EduTrace?"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 service unavailable, got %d", rec.Code)
	}

	var errBody struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &errBody); err != nil {
		t.Fatalf("failed to decode error body: %v", err)
	}

	if errBody.Error.Code != "service_unavailable" {
		t.Errorf("expected code service_unavailable, got %s", errBody.Error.Code)
	}
	if errBody.Error.Message != "Ask Arham AI is currently unavailable." {
		t.Errorf("expected specific disabled message, got %s", errBody.Error.Message)
	}
}

func TestAskHandler_RateLimiter(t *testing.T) {
	fakeLLM := llm.NewDeterministicFakeProvider()
	retriever := &mockHTTPRetriever{
		results: []domain.RetrievalResult{
			{ChunkID: "chk_1", SourceTitle: "Title", Content: "Content"},
		},
	}
	// Strict limit: 2 requests per minute
	router := setupTestServer("remote", retriever, fakeLLM, 4, 2)

	// Call 1: OK
	req1 := httptest.NewRequest(http.MethodPost, "/api/v1/ai/ask", bytes.NewBufferString(`{"question": "Question 1"}`))
	req1.RemoteAddr = "192.0.2.1:12345"
	rec1 := httptest.NewRecorder()
	router.ServeHTTP(rec1, req1)
	if rec1.Code != http.StatusOK {
		t.Fatalf("call 1 failed with status %d", rec1.Code)
	}

	// Call 2: OK
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/ai/ask", bytes.NewBufferString(`{"question": "Question 2"}`))
	req2.RemoteAddr = "192.0.2.1:12346"
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("call 2 failed with status %d", rec2.Code)
	}

	// Call 3: Exceeded -> 429
	req3 := httptest.NewRequest(http.MethodPost, "/api/v1/ai/ask", bytes.NewBufferString(`{"question": "Question 3"}`))
	req3.RemoteAddr = "192.0.2.1:12347"
	rec3 := httptest.NewRecorder()
	router.ServeHTTP(rec3, req3)
	if rec3.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 too many requests on call 3, got %d", rec3.Code)
	}
}

func TestAskHandler_ConcurrencySemaphore(t *testing.T) {
	fakeLLM := llm.NewDeterministicFakeProvider()
	retriever := &mockHTTPRetriever{
		delay: 100 * time.Millisecond,
		results: []domain.RetrievalResult{
			{ChunkID: "chk_1", SourceTitle: "Title", Content: "Content"},
		},
	}
	// Max concurrent: 1
	router := setupTestServer("remote", retriever, fakeLLM, 1, 100)

	var wg sync.WaitGroup
	var status1, status2 int

	wg.Add(1)
	go func() {
		defer wg.Done()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/ai/ask", bytes.NewBufferString(`{"question": "Concurrent 1"}`))
		req.RemoteAddr = "10.0.0.1:12345"
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		status1 = rec.Code
	}()

	time.Sleep(10 * time.Millisecond) // Ensure request 1 has acquired the permit

	wg.Add(1)
	go func() {
		defer wg.Done()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/ai/ask", bytes.NewBufferString(`{"question": "Concurrent 2"}`))
		req.RemoteAddr = "10.0.0.2:12345"
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		status2 = rec.Code
	}()

	wg.Wait()

	if status1 != http.StatusOK {
		t.Errorf("expected request 1 status 200, got %d", status1)
	}
	if status2 != http.StatusServiceUnavailable {
		t.Errorf("expected request 2 status 503 when concurrency saturated, got %d", status2)
	}
}

func TestAskHandler_SuccessfulSSEStream(t *testing.T) {
	fakeLLM := llm.NewDeterministicFakeProvider()
	fakeLLM.SetAnswer(llm.GeneratedAnswer{
		Status: "supported",
		Segments: []llm.GroundedSegment{
			{Text: "EduTrace uses Soulbound Tokens.", EvidenceIDs: []string{"E1"}},
		},
		SuggestedActionIDs: []string{"view-project-edutrace"},
	})

	projID := "edutrace"
	retriever := &mockHTTPRetriever{
		results: []domain.RetrievalResult{
			{
				ChunkID:     "chk_1",
				SourceTitle: "EduTrace SBT",
				SourceType:  "github",
				Repository:  "Nachsyas/EduTrace",
				Path:        "contracts/EduTraceSBT.sol",
				CommitSHA:   "1234567890123456789012345678901234567890",
				SourceURL:   "https://github.com/Nachsyas/EduTrace/blob/1234567890123456789012345678901234567890/contracts/EduTraceSBT.sol",
				Content:     "contract EduTraceSBT is ERC5192 {}",
				ProjectID:   &projID,
			},
		},
	}

	router := setupTestServer("remote", retriever, fakeLLM, 4, 10)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/ai/ask", bytes.NewBufferString(`{"question": "How does EduTrace work?"}`))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}
	if !strings.Contains(rec.Header().Get("Content-Type"), "text/event-stream") {
		t.Fatalf("expected text/event-stream, got %s", rec.Header().Get("Content-Type"))
	}
	if rec.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("expected no-store, got %s", rec.Header().Get("Cache-Control"))
	}
	if rec.Header().Get("X-Request-ID") == "" {
		t.Fatal("expected X-Request-ID header to be present")
	}

	bodyStr := rec.Body.String()
	// Verify SSE event sequence (Correction 6)
	expectedEvents := []string{
		"event: status\ndata: {\"status\":\"retrieving\"}",
		"event: evidence\n",
		"event: status\ndata: {\"status\":\"generating\"}",
		"event: result\n",
		"event: done\n",
	}

	for _, ev := range expectedEvents {
		if !strings.Contains(bodyStr, ev) {
			t.Errorf("missing expected SSE event fragment: %s\nGot body:\n%s", ev, bodyStr)
		}
	}
}

func TestAskHandler_ErrorAfterSSEStart(t *testing.T) {
	fakeLLM := llm.NewDeterministicFakeProvider()
	retriever := &mockHTTPRetriever{
		err: errors.New("database connection terminated"),
	}

	router := setupTestServer("remote", retriever, fakeLLM, 4, 10)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/ai/ask", bytes.NewBufferString(`{"question": "How does EduTrace work?"}`))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Since SSE started, HTTP status is 200 but body contains event: error then event: done (Correction 10)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 for SSE stream start, got %d", rec.Code)
	}

	bodyStr := rec.Body.String()
	if !strings.Contains(bodyStr, "event: error\n") {
		t.Fatalf("expected event: error, got:\n%s", bodyStr)
	}
	if !strings.Contains(bodyStr, "event: done\n") {
		t.Fatalf("expected event: done, got:\n%s", bodyStr)
	}
	if strings.Contains(bodyStr, "event: result\n") {
		t.Fatalf("result should NOT be emitted after error")
	}
}
