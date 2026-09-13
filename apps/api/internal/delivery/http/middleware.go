package http

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"github.com/nachsyas/arham-porto/apps/api/internal/delivery/http/dto"
)

// responseRecorder captures the status code of the HTTP response.
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

// RecoveryMiddleware catches panics, logs diagnostics on the server, and returns generic 500 JSON.
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("[PANIC RECOVERED] path=%s err=%v\nstack:\n%s", r.URL.Path, rec, debug.Stack())
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Cache-Control", "no-store")
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(dto.NewErrorEnvelope("internal_error", "an internal server error occurred"))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// LoggerMiddleware logs basic privacy-safe HTTP request metrics without logging credentials or bodies.
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rec, r)

		remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			remoteIP = r.RemoteAddr
		}

		log.Printf("[HTTP] %s %s %d %s (client: %s)",
			r.Method,
			r.URL.Path,
			rec.statusCode,
			time.Since(start),
			remoteIP,
		)
	})
}

// SecurityHeadersMiddleware attaches modern, non-deprecated security headers to every response.
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
		w.Header().Set("X-Frame-Options", "DENY")
		next.ServeHTTP(w, r)
	})
}

// CORSMiddleware handles cross-origin requests using an allowlist with Vary: Origin support.
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
			w.Header().Add("Vary", "Origin")

			origin := r.Header.Get("Origin")
			if origin != "" {
				if allowAll {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				} else if originSet[origin] {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				}
				w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// clientBucket tracks request counts for a single IP.
type clientBucket struct {
	count     int
	resetTime time.Time
}

// RateLimiter manages in-memory rate limiting with bounded memory and periodic cleanup.
type RateLimiter struct {
	mu          sync.Mutex
	limit       int
	window      time.Duration
	clients     map[string]*clientBucket
	stopCleanup chan struct{}
}

// NewRateLimiter creates a RateLimiter instance.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		limit:       limit,
		window:      window,
		clients:     make(map[string]*clientBucket),
		stopCleanup: make(chan struct{}),
	}
	go rl.cleanupRoutine()
	return rl
}

// Allow checks if a request from the given IP is allowed.
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	bucket, exists := rl.clients[ip]
	if !exists || now.After(bucket.resetTime) {
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

			remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				remoteIP = r.RemoteAddr
			}

			if !rl.Allow(remoteIP) {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Cache-Control", "no-store")
				w.WriteHeader(http.StatusTooManyRequests)
				_ = json.NewEncoder(w).Encode(dto.NewErrorEnvelope("rate_limited", "too many requests, please slow down"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
