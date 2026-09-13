package gemini

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nachsyas/arham-porto/apps/api/internal/embedding"
)

func TestGeminiProviderBatchAndPurpose(t *testing.T) {
	var capturedTaskType string
	var capturedBatchSize int
	var capturedOutputDim int

	// Mock server mimicking Gemini batchEmbedContents API
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, ":batchEmbedContents") {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var req batchEmbedRequest
		if err := json.Unmarshal(body, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		capturedBatchSize = len(req.Requests)
		if capturedBatchSize > 0 {
			capturedTaskType = req.Requests[0].TaskType
			capturedOutputDim = req.Requests[0].OutputDimensionality
		}

		// Return N embeddings with dimension 768
		embeddings := make([]struct {
			Values []float32 `json:"values"`
		}, len(req.Requests))

		for i := range embeddings {
			vals := make([]float32, 768)
			for d := range vals {
				vals[d] = float32(d) / 768.0
			}
			embeddings[i].Values = vals
		}

		resp := batchEmbedResponse{Embeddings: embeddings}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	provider := NewTestProvider(ts.URL, "fake-api-key", "gemini-embedding-2", 768)
	ctx := context.Background()

	t.Run("3 texts -> exactly 3 vectors with 768 dimensions (Gate #28)", func(t *testing.T) {
		texts := []string{
			"Go backend microservice architecture",
			"EduTrace Soulbound Token implementation",
			"Next.js App Router interactive UI components",
		}

		vectors, err := provider.Embed(ctx, embedding.PurposeDocument, texts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(vectors) != 3 {
			t.Fatalf("expected exactly 3 vectors, got %d", len(vectors))
		}

		for i, vec := range vectors {
			if len(vec) != 768 {
				t.Errorf("vector %d: expected dimension 768, got %d", i, len(vec))
			}
		}

		if capturedBatchSize != 3 {
			t.Errorf("expected server to receive 3 requests, got %d", capturedBatchSize)
		}
		if capturedOutputDim != 768 {
			t.Errorf("expected output dimensionality 768, got %d", capturedOutputDim)
		}
	})

	t.Run("Document purpose translates to RETRIEVAL_DOCUMENT (Gate #29)", func(t *testing.T) {
		_, err := provider.Embed(ctx, embedding.PurposeDocument, []string{"document text"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if capturedTaskType != "RETRIEVAL_DOCUMENT" {
			t.Errorf("expected RETRIEVAL_DOCUMENT, got %s", capturedTaskType)
		}
	})

	t.Run("Query purpose translates to RETRIEVAL_QUERY (Gate #29)", func(t *testing.T) {
		_, err := provider.Embed(ctx, embedding.PurposeQuery, []string{"query text"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if capturedTaskType != "RETRIEVAL_QUERY" {
			t.Errorf("expected RETRIEVAL_QUERY, got %s", capturedTaskType)
		}
	})
}
