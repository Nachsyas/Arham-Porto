package http

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/nachsyas/arham-porto/apps/api/internal/llm"
)

// generateRequestID produces a cryptographically random, privacy-safe request ID.
func generateRequestID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("req_%d", time.Now().UnixNano())
	}
	return "req_" + hex.EncodeToString(b)
}

// writeSSE writes a securely JSON-serialized SSE event and flushes immediately.
func writeSSE(w io.Writer, flusher http.Flusher, event string, data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal sse data: %w", err)
	}
	if _, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, string(payload)); err != nil {
		return err
	}
	flusher.Flush()
	return nil
}

// AskHandler processes reviewer questions with preflight validation and SSE streaming.
func (h *Handler) AskHandler(w http.ResponseWriter, r *http.Request) {
	reqID := generateRequestID()
	w.Header().Set("X-Request-ID", reqID)

	// ==========================================
	// 1. Preflight Validation (BEFORE SSE headers)
	// ==========================================

	// Body decode with bounded read limit (64 KB)
	var body struct {
		Question string `json:"question"`
	}
	if err := json.NewDecoder(io.LimitReader(r.Body, 64*1024)).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid request body")
		return
	}

	trimmedQ := strings.TrimSpace(body.Question)
	if len(trimmedQ) < 2 || len(trimmedQ) > 1000 {
		writeError(w, http.StatusBadRequest, "bad_request", "question must be between 2 and 1000 characters")
		return
	}

	// AI Mode check (Correction 52)
	if h.aiMode == "disabled" {
		writeError(w, http.StatusServiceUnavailable, "service_unavailable", "Ask Arham AI is currently unavailable.")
		return
	}

	// Readiness check (Correction 53)
	if h.askUC == nil {
		writeError(w, http.StatusServiceUnavailable, "service_unavailable", "Ask Arham AI is currently unavailable.")
		return
	}

	// AI Rate Limiter check (Correction 15, 37, Gate #15)
	clientIP := ResolveClientIP(r, h.trustProxyMode)
	if h.aiRateLimiter != nil && !h.aiRateLimiter.Allow(clientIP) {
		writeError(w, http.StatusTooManyRequests, "rate_limited", "AI rate limit exceeded. Please try again in a minute.")
		return
	}

	// Concurrency Semaphore acquisition (Correction 16, 36)
	if h.aiSemaphore != nil {
		select {
		case h.aiSemaphore <- struct{}{}:
			defer func() { <-h.aiSemaphore }()
		default:
			writeError(w, http.StatusServiceUnavailable, "service_unavailable", "Ask Arham AI is currently busy. Please try again in a moment.")
			return
		}
	}

	// Flusher verification
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, http.StatusInternalServerError, "internal_error", "streaming not supported")
		return
	}

	// ==========================================
	// 2. Commit SSE Headers (Correction 40)
	// ==========================================
	setSecurityHeaders(w)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	// ==========================================
	// 3. SSE Stream Execution
	// ==========================================
	timeoutDuration := 30 * time.Second
	if h.aiTimeoutSeconds > 0 {
		timeoutDuration = time.Duration(h.aiTimeoutSeconds) * time.Second
	}
	ctx, cancel := context.WithTimeout(r.Context(), timeoutDuration)
	defer cancel()

	// Status: retrieving
	_ = writeSSE(w, flusher, "status", map[string]string{"status": "retrieving"})

	// Fetch evidence to expose verified metadata early to reviewer
	evidenceList, err := h.askUC.GetRetrievedEvidence(ctx, trimmedQ)
	if err != nil {
		_ = writeSSE(w, flusher, "error", map[string]string{
			"code":       "service_unavailable",
			"message":    "Ask Arham AI is temporarily unavailable.",
			"request_id": reqID,
		})
		_ = writeSSE(w, flusher, "done", map[string]any{})
		return
	}

	// Format public evidence metadata (safe excerpts, zero similarity scores, zero raw paths)
	publicEvList := make([]llm.PublicEvidenceItem, 0, len(evidenceList))
	for _, e := range evidenceList {
		excerpt := e.Content
		if len(excerpt) > 600 {
			excerpt = excerpt[:600] + "..."
		}
		publicEvList = append(publicEvList, llm.PublicEvidenceItem{
			ID:         e.ID,
			Kind:       e.Kind,
			Title:      e.Title,
			Repository: e.Repository,
			Path:       e.Path,
			Excerpt:    excerpt,
		})
	}
	_ = writeSSE(w, flusher, "evidence", map[string]any{"evidence": publicEvList})

	// Status: generating
	_ = writeSSE(w, flusher, "status", map[string]string{"status": "generating"})

	// Execute full grounded Q&A with claim-level validation
	groundedResp, err := h.askUC.Ask(ctx, trimmedQ)
	if err != nil {
		_ = writeSSE(w, flusher, "error", map[string]string{
			"code":       "service_unavailable",
			"message":    "Ask Arham AI is temporarily unavailable.",
			"request_id": reqID,
		})
		_ = writeSSE(w, flusher, "done", map[string]any{})
		return
	}

	// Result: validated GroundedResponse (Correction 6)
	_ = writeSSE(w, flusher, "result", groundedResp)

	// Done event
	_ = writeSSE(w, flusher, "done", map[string]any{})
}
