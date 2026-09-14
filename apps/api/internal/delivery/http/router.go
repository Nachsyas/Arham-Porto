package http

import (
	"net/http"
)

// NewRouter registers HTTP routes using Go standard library net/http.ServeMux and wires the middleware stack.
func NewRouter(h *Handler, allowedOrigins []string, rl *RateLimiter, trustProxyMode ...string) http.Handler {
	mode := "direct"
	if len(trustProxyMode) > 0 && trustProxyMode[0] != "" {
		mode = trustProxyMode[0]
	}
	h.SetTrustProxyMode(mode)

	mux := http.NewServeMux()

	// Probes
	mux.HandleFunc("GET /health", h.Healthz)
	mux.HandleFunc("GET /healthz", h.Healthz)
	mux.HandleFunc("GET /readyz", h.Readyz)

	// API v1 Endpoints
	mux.HandleFunc("GET /api/v1/profile", h.GetProfile)
	mux.HandleFunc("GET /api/v1/projects", h.ListProjects)
	mux.HandleFunc("GET /api/v1/projects/{slug}", h.GetProjectBySlug)
	mux.HandleFunc("GET /api/v1/skills", h.ListSkills)
	mux.HandleFunc("GET /api/v1/evidence", h.ListEvidence)
	mux.HandleFunc("GET /api/v1/evidence/{id}", h.GetEvidenceByID)
	mux.HandleFunc("GET /api/v1/journey", h.ListJourney)
	mux.HandleFunc("POST /api/v1/ai/ask", h.AskHandler)

	// Middleware composition (outer -> inner)
	// 1. Recovery
	// 2. Logger
	// 3. SecurityHeaders
	// 4. CORS
	// 5. RateLimiter
	// 6. MethodNotAllowed
	var handler http.Handler = mux

	handler = MethodNotAllowedMiddleware(handler)

	if rl != nil {
		handler = RateLimiterMiddleware(rl, mode)(handler)
	}

	handler = CORSMiddleware(allowedOrigins)(handler)
	handler = SecurityHeadersMiddleware(handler)
	handler = LoggerMiddleware(handler)
	handler = RecoveryMiddleware(handler)

	return handler
}
