package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	delivery "github.com/nachsyas/arham-porto/apps/api/internal/delivery/http"
)

func TestResolveClientIP(t *testing.T) {
	testCases := []struct {
		name           string
		trustProxyMode string
		remoteAddr     string
		headers        map[string]string
		expectedIP     string
	}{
		{
			name:           "direct mode ignores spoofed X-Real-IP",
			trustProxyMode: "direct",
			remoteAddr:     "192.168.1.50:54321",
			headers: map[string]string{
				"X-Real-IP":       "203.0.113.195",
				"X-Forwarded-For": "198.51.100.1",
			},
			expectedIP: "192.168.1.50",
		},
		{
			name:           "direct mode ignores X-Forwarded-For",
			trustProxyMode: "direct",
			remoteAddr:     "10.0.0.1:8080",
			headers: map[string]string{
				"X-Forwarded-For": "203.0.113.5",
			},
			expectedIP: "10.0.0.1",
		},
		{
			name:           "railway mode uses valid X-Real-IP",
			trustProxyMode: "railway",
			remoteAddr:     "10.0.0.1:8080",
			headers: map[string]string{
				"X-Real-IP": "203.0.113.195",
			},
			expectedIP: "203.0.113.195",
		},
		{
			name:           "railway mode uses valid IPv6 X-Real-IP",
			trustProxyMode: "railway",
			remoteAddr:     "10.0.0.1:8080",
			headers: map[string]string{
				"X-Real-IP": "2001:db8::1",
			},
			expectedIP: "2001:db8::1",
		},
		{
			name:           "railway mode ignores spoofed X-Forwarded-For when X-Real-IP is absent",
			trustProxyMode: "railway",
			remoteAddr:     "10.0.0.2:8080",
			headers: map[string]string{
				"X-Forwarded-For": "203.0.113.99",
			},
			expectedIP: "10.0.0.2",
		},
		{
			name:           "railway mode safely falls back on malformed X-Real-IP",
			trustProxyMode: "railway",
			remoteAddr:     "10.0.0.3:8080",
			headers: map[string]string{
				"X-Real-IP": "not-an-ip-address",
			},
			expectedIP: "10.0.0.3",
		},
		{
			name:           "railway mode safely falls back on empty X-Real-IP",
			trustProxyMode: "railway",
			remoteAddr:     "10.0.0.4:8080",
			headers: map[string]string{
				"X-Real-IP": "   ",
			},
			expectedIP: "10.0.0.4",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/profile", nil)
			req.RemoteAddr = tc.remoteAddr
			for k, v := range tc.headers {
				req.Header.Set(k, v)
			}

			actualIP := delivery.ResolveClientIP(req, tc.trustProxyMode)
			if actualIP != tc.expectedIP {
				t.Errorf("expected IP %q, got %q", tc.expectedIP, actualIP)
			}
		})
	}
}

func TestRateLimiterMiddleware_TrustedProxyMode(t *testing.T) {
	// Create rate limiter allowing 1 request per minute
	rl := delivery.NewRateLimiter(1, time.Minute)
	defer rl.Close()

	// Railway proxy mode middleware
	mw := delivery.RateLimiterMiddleware(rl, "railway")

	var hitCount int
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hitCount++
		w.WriteHeader(http.StatusOK)
	}))

	// Request 1 from Client A (X-Real-IP: 203.0.113.1) through Railway edge proxy (RemoteAddr: 10.0.0.1:1234)
	reqA1 := httptest.NewRequest("GET", "/api/v1/projects", nil)
	reqA1.RemoteAddr = "10.0.0.1:1234"
	reqA1.Header.Set("X-Real-IP", "203.0.113.1")
	recA1 := httptest.NewRecorder()
	handler.ServeHTTP(recA1, reqA1)
	if recA1.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for first request from Client A, got %d", recA1.Code)
	}

	// Request 2 from Client A should be rate limited (429)
	reqA2 := httptest.NewRequest("GET", "/api/v1/projects", nil)
	reqA2.RemoteAddr = "10.0.0.1:1234" // Same edge proxy IP
	reqA2.Header.Set("X-Real-IP", "203.0.113.1")
	recA2 := httptest.NewRecorder()
	handler.ServeHTTP(recA2, reqA2)
	if recA2.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 Too Many Requests for second request from Client A, got %d", recA2.Code)
	}

	// Request 1 from Client B (X-Real-IP: 203.0.113.2) through SAME Railway edge proxy (RemoteAddr: 10.0.0.1:1234)
	// Must SUCCEED because Client B is a distinct client IP and not collapsed into the proxy IP!
	reqB1 := httptest.NewRequest("GET", "/api/v1/projects", nil)
	reqB1.RemoteAddr = "10.0.0.1:1234" // Same edge proxy IP
	reqB1.Header.Set("X-Real-IP", "203.0.113.2")
	recB1 := httptest.NewRecorder()
	handler.ServeHTTP(recB1, reqB1)
	if recB1.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for Client B under Railway mode, got %d (clients collapsed into proxy IP)", recB1.Code)
	}
}
