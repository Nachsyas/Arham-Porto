package http

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/nachsyas/arham-porto/apps/api/internal/delivery/http/dto"
)

// ExtractClientIP extracts the normalized IP or host from RemoteAddr, discarding ephemeral TCP source ports.
// It safely handles malformed addresses, missing ports, and IPv6 brackets.
// It does NOT trust raw spoofable headers (X-Forwarded-For) without a trusted reverse proxy configuration.
func ExtractClientIP(remoteAddr string) string {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err == nil {
		return strings.Trim(host, "[]")
	}
	trimmed := strings.TrimSpace(remoteAddr)
	return strings.Trim(trimmed, "[]")
}

// addVaryHeader adds a Vary header value without duplicating existing entries or overwriting.
func addVaryHeader(w http.ResponseWriter, value string) {
	existing := w.Header().Get("Vary")
	if existing == "" {
		w.Header().Set("Vary", value)
		return
	}
	for _, p := range strings.Split(existing, ",") {
		if strings.EqualFold(strings.TrimSpace(p), value) {
			return
		}
	}
	w.Header().Set("Vary", existing+", "+value)
}

// setSecurityHeaders writes standard security headers to the ResponseWriter.
func setSecurityHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
	w.Header().Set("X-Frame-Options", "DENY")
}

// responseRecorder captures the status code of the HTTP response.
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *responseRecorder) Flush() {
	if flusher, ok := r.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// RecoveryMiddleware catches panics, logs server-side diagnostics, and returns generic 500 JSON with security headers.
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("[PANIC RECOVERED] path=%s err=%v\nstack:\n%s", r.URL.Path, rec, debug.Stack())
				setSecurityHeaders(w)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.Header().Set("Cache-Control", "no-store")
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(dto.NewErrorEnvelope("internal_error", "an internal server error occurred"))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// LoggerMiddleware logs privacy-safe HTTP request metrics without logging credentials or request bodies.
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rec, r)

		clientIP := ExtractClientIP(r.RemoteAddr)
		log.Printf("[HTTP] %s %s %d %s (client: %s)",
			r.Method,
			r.URL.Path,
			rec.statusCode,
			time.Since(start),
			clientIP,
		)
	})
}

// SecurityHeadersMiddleware attaches modern, non-deprecated security headers to every response.
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setSecurityHeaders(w)
		next.ServeHTTP(w, r)
	})
}

// CORSMiddleware handles cross-origin requests using an explicit allowlist and safe Vary: Origin handling.
func CORSMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	originSet := make(map[string]bool)
	allowAll := false
	for _, o := range allowedOrigins {
		if o == "*" {
			allowAll = true
		} else {
			originSet[o] = true
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			addVaryHeader(w, "Origin")

			origin := r.Header.Get("Origin")
			if origin != "" {
				if allowAll {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				} else if originSet[origin] {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				}
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, X-Request-ID")
			}

			if r.Method == http.MethodOptions {
				setSecurityHeaders(w)
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// clientBucket tracks request counts for a single normalized client IP.
type clientBucket struct {
	count     int
	resetTime time.Time
}

const defaultMaxRateLimiterEntries = 10000

// RateLimiter manages in-memory rate limiting with bounded memory and periodic cleanup.
type RateLimiter struct {
	mu          sync.Mutex
	limit       int
	window      time.Duration
	maxEntries  int
	clients     map[string]*clientBucket
	stopCleanup chan struct{}
}

// NewRateLimiter creates a RateLimiter instance with bounded memory.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		limit:       limit,
		window:      window,
		maxEntries:  defaultMaxRateLimiterEntries,
		clients:     make(map[string]*clientBucket),
		stopCleanup: make(chan struct{}),
	}
	go rl.cleanupRoutine()
	return rl
}

// Allow checks if a request from the given IP is allowed.
// If the memory boundary is reached under high client cardinality, it prunes expired entries immediately.
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	bucket, exists := rl.clients[ip]
	if !exists || now.After(bucket.resetTime) {
		// Enforce bounded memory size
		if len(rl.clients) >= rl.maxEntries {
			for k, b := range rl.clients {
				if now.After(b.resetTime) {
					delete(rl.clients, k)
				}
			}
			// If still full after pruning, reject new allocations to prevent unbounded growth
			if len(rl.clients) >= rl.maxEntries {
				return false
			}
		}

		rl.clients[ip] = &clientBucket{
			count:     1,
			resetTime: now.Add(rl.window),
		}
		return true
	}

	if bucket.count < rl.limit {
		bucket.count++
		return true
	}

	return false
}

// cleanupRoutine periodically removes expired IP buckets to prevent memory unbounded growth.
func (rl *RateLimiter) cleanupRoutine() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.mu.Lock()
			now := time.Now()
			for ip, bucket := range rl.clients {
				if now.After(bucket.resetTime) {
					delete(rl.clients, ip)
				}
			}
			rl.mu.Unlock()
		case <-rl.stopCleanup:
			return
		}
	}
}

// Close terminates the rate limiter's background cleanup goroutine.
func (rl *RateLimiter) Close() {
	select {
	case <-rl.stopCleanup:
	default:
		close(rl.stopCleanup)
	}
}

// RateLimiterMiddleware throttles incoming requests while excluding health and readiness probes.
func RateLimiterMiddleware(rl *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Probes bypass rate limiting
			if r.URL.Path == "/healthz" || r.URL.Path == "/readyz" {
				next.ServeHTTP(w, r)
				return
			}

			clientIP := ExtractClientIP(r.RemoteAddr)

			if !rl.Allow(clientIP) {
				setSecurityHeaders(w)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.Header().Set("Cache-Control", "no-store")
				w.WriteHeader(http.StatusTooManyRequests)
				_ = json.NewEncoder(w).Encode(dto.NewErrorEnvelope("rate_limited", "too many requests, please slow down"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// methodNotAllowedInterceptor intercepts 405 Method Not Allowed responses generated by http.ServeMux
// and suppresses default text/plain output so that standard JSON error envelopes can be rendered.
type methodNotAllowedInterceptor struct {
	http.ResponseWriter
	statusCode int
	is405      bool
}

func (w *methodNotAllowedInterceptor) WriteHeader(code int) {
	w.statusCode = code
	if code == http.StatusMethodNotAllowed {
		w.is405 = true
		return
	}
	w.ResponseWriter.WriteHeader(code)
}

func (w *methodNotAllowedInterceptor) Write(b []byte) (int, error) {
	if w.is405 {
		// Suppress the default text/plain body from http.Error
		return len(b), nil
	}
	return w.ResponseWriter.Write(b)
}

func (w *methodNotAllowedInterceptor) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// MethodNotAllowedMiddleware intercepts 405 Method Not Allowed responses and transforms them
// into the certified public JSON error envelope while preserving Allow, Vary, and security headers.
func MethodNotAllowedMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		interceptor := &methodNotAllowedInterceptor{ResponseWriter: w}
		next.ServeHTTP(interceptor, r)

		if interceptor.is405 {
			setSecurityHeaders(w)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Header().Set("Cache-Control", "no-store")
			interceptor.ResponseWriter.WriteHeader(http.StatusMethodNotAllowed)
			_ = json.NewEncoder(interceptor.ResponseWriter).Encode(dto.NewErrorEnvelope(dto.ErrCodeMethodNotAllowed, "method not allowed"))
		}
	})
}

