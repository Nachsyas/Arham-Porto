package gemini_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nachsyas/arham-porto/apps/api/internal/llm"
	"github.com/nachsyas/arham-porto/apps/api/internal/llm/gemini"
)

func TestNewClientValidation(t *testing.T) {
	_, err := gemini.NewClient("", "gemini-3.8-flash")
	if err == nil {
		t.Fatal("expected error when api key is empty")
	}

	client, err := gemini.NewClient("test-key", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.Model() != "gemini-3.8-flash" {
		t.Fatalf("expected default model gemini-3.8-flash, got %s", client.Model())
	}
	if client.ProviderName() != "gemini" {
		t.Fatalf("expected provider name gemini, got %s", client.ProviderName())
	}
}

func TestInteractionsAPISuccess(t *testing.T) {
	expectedAnswer := llm.GeneratedAnswer{
		Status: "supported",
		Segments: []llm.GroundedSegment{
			{Text: "EduTrace uses Soulbound Tokens.", EvidenceIDs: []string{"E1"}},
		},
		SuggestedActionIDs: []string{"view-project-edutrace"},
	}

	answerJSON, _ := json.Marshal(expectedAnswer)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/interactions" {
			t.Errorf("expected /interactions endpoint, got %s", r.URL.Path)
		}
		if r.Header.Get("x-goog-api-key") != "test-api-key" {
			t.Errorf("expected x-goog-api-key header, got %s", r.Header.Get("x-goog-api-key"))
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"output": string(answerJSON),
		})
	}))
	defer server.Close()

	client, err := gemini.NewClient(
		"test-api-key",
		"gemini-3.8-flash",
		gemini.WithBaseURL(server.URL),
		gemini.WithHTTPClient(server.Client()),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	res, err := client.Generate(context.Background(), llm.GenerateRequest{
		SystemInstruction: "You are Ask Arham AI.",
		UserQuestion:      "How does EduTrace work?",
	})
	if err != nil {
		t.Fatalf("unexpected generate error: %v", err)
	}

	if res.Status != "supported" {
		t.Fatalf("expected supported status, got %s", res.Status)
	}
	if len(res.Segments) != 1 || res.Segments[0].EvidenceIDs[0] != "E1" {
		t.Fatalf("unexpected segments: %+v", res.Segments)
	}
}

func TestGenerateContentFallback(t *testing.T) {
	expectedAnswer := llm.GeneratedAnswer{
		Status: "supported",
		Segments: []llm.GroundedSegment{
			{Text: "Fallback answer works.", EvidenceIDs: []string{"E1"}},
		},
		SuggestedActionIDs: []string{"go-to-projects"},
	}
	answerJSON, _ := json.Marshal(expectedAnswer)
	fencedJSON := "```json\n" + string(answerJSON) + "\n```"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/interactions" {
			// Simulate interactions API not implemented / 404
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error": "not found"}`))
			return
		}

		// Fallback should hit /models/gemini-3.8-flash:generateContent
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"candidates": []map[string]any{
				{
					"content": map[string]any{
						"parts": []map[string]any{
							{"text": fencedJSON},
						},
					},
				},
			},
		})
	}))
	defer server.Close()

	client, err := gemini.NewClient(
		"test-api-key",
		"gemini-3.8-flash",
		gemini.WithBaseURL(server.URL),
		gemini.WithHTTPClient(server.Client()),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	res, err := client.Generate(context.Background(), llm.GenerateRequest{
		SystemInstruction: "You are Ask Arham AI.",
		UserQuestion:      "How does fallback work?",
	})
	if err != nil {
		t.Fatalf("unexpected generate error: %v", err)
	}

	if res.Status != "supported" {
		t.Fatalf("expected supported status, got %s", res.Status)
	}
	if res.Segments[0].Text != "Fallback answer works." {
		t.Fatalf("expected text from fenced markdown, got %s", res.Segments[0].Text)
	}
}

func TestGeminiAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error": "internal provider failure"}`))
	}))
	defer server.Close()

	client, err := gemini.NewClient(
		"test-api-key",
		"gemini-3.8-flash",
		gemini.WithBaseURL(server.URL),
		gemini.WithHTTPClient(server.Client()),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Generate(context.Background(), llm.GenerateRequest{
		UserQuestion: "Any question?",
	})
	if err == nil {
		t.Fatal("expected error on 500 response, got nil")
	}
}

func TestGeminiContextTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := gemini.NewClient(
		"test-api-key",
		"gemini-3.8-flash",
		gemini.WithBaseURL(server.URL),
		gemini.WithHTTPClient(server.Client()),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	_, err = client.Generate(ctx, llm.GenerateRequest{UserQuestion: "Timeout test"})
	if err == nil {
		t.Fatal("expected error on timed out context, got nil")
	}
}
