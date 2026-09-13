package http

import (
	"net/http"
)

// NewRouter registers HTTP routes using Go 1.22+ standard library net/http.
func NewRouter(h *Handler, allowedOrigin string) http.Handler {
	mux := http.NewServeMux()

	// Health and readiness probes
	mux.HandleFunc("GET /healthz", h.Healthz)
	mux.HandleFunc("GET /readyz", h.Readyz)

	// API v1 endpoints (Phase 0 skeleton)
	mux.HandleFunc("GET /api/v1/profile", h.GetProfile)

	// Apply middleware stack
	handler := RecoveryMiddleware(mux)
	handler = CORSMiddleware(allowedOrigin)(handler)
	handler = LoggerMiddleware(handler)

	return handler
}
