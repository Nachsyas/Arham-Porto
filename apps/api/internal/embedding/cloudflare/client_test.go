package cloudflare

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/nachsyas/arham-porto/apps/api/internal/embedding"
)

func makeDummyVector(dim int, base float32) []float32 {
	v := make([]float32, dim)
	for i := range v {
		v[i] = base + float32(i)*0.001
	}
	return v
}

func TestCloudflareProvider_ConstructorValidation(t *testing.T) {
	_, err := NewProvider("", "valid-token", DefaultModel, 768)
	if err == nil || !strings.Contains(err.Error(), "account id") {
		t.Errorf("expected error for empty account id, got %v", err)
	}

	_, err = NewProvider("valid-account", "", DefaultModel, 768)
	if err == nil || !strings.Contains(err.Error(), "api token") {
		t.Errorf("expected error for empty api token, got %v", err)
	}

	p, err := NewProvider("account-123", "token-abc", "", 0)
	if err != nil {
		t.Fatalf("unexpected constructor error: %v", err)
	}

	if p.ProviderName() != "cloudflare" {
		t.Errorf("expected provider name 'cloudflare', got %q", p.ProviderName())
	}
	if p.Model() != DefaultModel {
		t.Errorf("expected default model %q, got %q", DefaultModel, p.Model())
	}
	if p.Dimensions() != DefaultDimensions {
		t.Errorf("expected dimensions %d, got %d", DefaultDimensions, p.Dimensions())
	}
}

func TestCloudflareProvider_Success_Batch(t *testing.T) {
	v1 := makeDummyVector(768, 0.1)
	v2 := makeDummyVector(768, 0.2)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if auth := r.Header.Get("Authorization"); auth != "Bearer secret-token" {
			t.Errorf("expected Bearer secret-token, got %s", auth)
		}
		if !strings.Contains(r.URL.Path, "/accounts/acc-123/ai/run/@cf/baai/bge-base-en-v1.5") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		resp := map[string]any{
			"success":  true,
			"errors":   []any{},
			"messages": []any{},
			"result": map[string]any{
				"shape": []int{2, 768},
				"data":  [][]float32{v1, v2},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	p, err := NewProvider("acc-123", "secret-token", DefaultModel, 768, WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewProvider failed: %v", err)
	}

	vectors, err := p.Embed(context.Background(), embedding.PurposeDocument, []string{"text1", "text2"})
	if err != nil {
		t.Fatalf("Embed failed: %v", err)
	}

	if len(vectors) != 2 {
		t.Fatalf("expected 2 vectors, got %d", len(vectors))
	}
	if len(vectors[0]) != 768 || len(vectors[1]) != 768 {
		t.Fatalf("expected 768 dimensions per vector")
	}
}

func TestCloudflareProvider_Success_SingleText1D(t *testing.T) {
	v1 := makeDummyVector(768, 0.5)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"success": true,
			"result": map[string]any{
				"shape": []int{768},
				"data":  v1,
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	p, err := NewProvider("acc-123", "token", DefaultModel, 768, WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewProvider failed: %v", err)
	}

	vectors, err := p.Embed(context.Background(), embedding.PurposeQuery, []string{"single text"})
	if err != nil {
		t.Fatalf("Embed failed: %v", err)
	}
	if len(vectors) != 1 || len(vectors[0]) != 768 {
		t.Fatalf("unexpected vector shape: %v", vectors)
	}
}

func TestCloudflareProvider_Errors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantErrStr string
	}{
		{
			name:       "400 bad request",
			statusCode: http.StatusBadRequest,
			body:       `{"success":false,"errors":[{"code":1000,"message":"Invalid model input"}]}`,
			wantErrStr: "status 400: [1000] Invalid model input",
		},
		{
			name:       "401 unauthorized",
			statusCode: http.StatusUnauthorized,
			body:       `{"success":false,"errors":[{"code":10000,"message":"Authentication error"}]}`,
			wantErrStr: "status 401: [10000] Authentication error",
		},
		{
			name:       "403 forbidden",
			statusCode: http.StatusForbidden,
			body:       `{"success":false,"errors":[{"code":10001,"message":"Unauthorized"}]}`,
			wantErrStr: "status 403: [10001] Unauthorized",
		},
		{
			name:       "429 rate limited",
			statusCode: http.StatusTooManyRequests,
			body:       `{"success":false,"errors":[{"code":10014,"message":"Rate limit exceeded"}]}`,
			wantErrStr: "status 429: [10014] Rate limit exceeded",
		},
		{
			name:       "500 internal error",
			statusCode: http.StatusInternalServerError,
			body:       `{"success":false,"errors":[{"code":9999,"message":"Worker crashed"}]}`,
			wantErrStr: "status 500: [9999] Worker crashed",
		},
		{
			name:       "malformed json non-200",
			statusCode: http.StatusBadGateway,
			body:       `<html>502 Bad Gateway</html>`,
			wantErrStr: "status 502 (malformed response body)",
		},
		{
			name:       "success true but malformed result",
			statusCode: http.StatusOK,
			body:       `{"success":true,"result":"invalid"}`,
			wantErrStr: "failed parsing embedding result",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
				_, _ = w.Write([]byte(tc.body))
			}))
			defer server.Close()

			p, err := NewProvider("acc-123", "token", DefaultModel, 768, WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("NewProvider failed: %v", err)
			}

			_, err = p.Embed(context.Background(), embedding.PurposeDocument, []string{"test"})
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.wantErrStr) {
				t.Errorf("expected error containing %q, got %q", tc.wantErrStr, err.Error())
			}
		})
	}
}

func TestCloudflareProvider_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{"success": true})
	}))
	defer server.Close()

	p, err := NewProvider("acc-123", "token", DefaultModel, 768, WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewProvider failed: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = p.Embed(ctx, embedding.PurposeDocument, []string{"test"})
	if err == nil {
		t.Fatalf("expected error on cancelled context, got nil")
	}
}

func TestCloudflareProvider_OversizedResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// Write 6 MiB of whitespace
		chunk := make([]byte, 1024*1024)
		for i := range chunk {
			chunk[i] = ' '
		}
		for i := 0; i < 6; i++ {
			_, _ = w.Write(chunk)
		}
	}))
	defer server.Close()

	p, err := NewProvider("acc-123", "token", DefaultModel, 768, WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewProvider failed: %v", err)
	}

	_, err = p.Embed(context.Background(), embedding.PurposeDocument, []string{"test"})
	if err == nil || !strings.Contains(err.Error(), "exceeded maximum allowed size") {
		t.Errorf("expected oversized response error, got %v", err)
	}
}

func TestCloudflareProvider_EmptyTexts(t *testing.T) {
	p, _ := NewProvider("acc-123", "token", DefaultModel, 768)
	vectors, err := p.Embed(context.Background(), embedding.PurposeDocument, nil)
	if err != nil {
		t.Errorf("expected nil error on empty texts, got %v", err)
	}
	if vectors != nil {
		t.Errorf("expected nil vectors, got %v", vectors)
	}
}
