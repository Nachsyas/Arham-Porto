package gemini_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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

	client, err := gemini.NewClient("test-key", "", gemini.WithThinkingLevel("medium"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.Model() != "gemini-3.8-flash" {
		t.Fatalf("expected default model gemini-3.8-flash, got %s", client.Model())
	}
	if client.ProviderName() != "gemini" {
		t.Fatalf("expected provider name gemini, got %s", client.ProviderName())
	}
	if client.ThinkingLevel() != "medium" {
		t.Fatalf("expected thinking level medium, got %s", client.ThinkingLevel())
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
		// Endpoint check
		if r.URL.Path != "/interactions" {
			t.Errorf("expected /interactions endpoint, got %s", r.URL.Path)
		}
		// Auth header check
		if r.Header.Get("x-goog-api-key") != "test-api-key" {
			t.Errorf("expected x-goog-api-key header, got %s", r.Header.Get("x-goog-api-key"))
		}

		// Inspect incoming request payload per Correction 4
		var reqBody map[string]any
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		// Assert response_format schema
		respFormat, ok := reqBody["response_format"].(map[string]any)
		if !ok {
			t.Fatal("expected response_format object in request")
		}
		if respFormat["type"] != "text" {
			t.Errorf("expected response_format.type == 'text', got %v", respFormat["type"])
		}
		if respFormat["mime_type"] != "application/json" {
			t.Errorf("expected response_format.mime_type == 'application/json', got %v", respFormat["mime_type"])
		}
		schema, ok := respFormat["schema"].(map[string]any)
		if !ok || schema["type"] != "object" {
			t.Fatalf("expected response_format.schema object, got %v", respFormat["schema"])
		}
		props, ok := schema["properties"].(map[string]any)
		if !ok || props["status"] == nil || props["segments"] == nil || props["suggested_action_ids"] == nil {
			t.Fatalf("expected required properties in schema, got %v", props)
		}

		// Assert generation_config thinking level and thinking_summaries
		genConfig, ok := reqBody["generation_config"].(map[string]any)
		if !ok {
			t.Fatal("expected generation_config in request")
		}
		if genConfig["thinking_level"] != "low" {
			t.Errorf("expected thinking_level == 'low', got %v", genConfig["thinking_level"])
		}
		if genConfig["thinking_summaries"] != "none" {
			t.Errorf("expected thinking_summaries == 'none', got %v", genConfig["thinking_summaries"])
		}

		// Assert tools are absent / not passed
		if reqBody["tools"] != nil {
			t.Errorf("expected tools to be absent, got %v", reqBody["tools"])
		}

		// Return current 2026 response structure
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status": "completed",
			"steps": []map[string]any{
				{
					"type": "model_output",
					"content": []map[string]any{
						{
							"type": "text",
							"text": string(answerJSON),
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
		gemini.WithThinkingLevel("low"),
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
	if len(res.SuggestedActionIDs) != 1 || res.SuggestedActionIDs[0] != "view-project-edutrace" {
		t.Fatalf("unexpected suggested actions: %+v", res.SuggestedActionIDs)
	}
}

func TestInteractionsStatusFailed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status": "failed",
			"steps":  []map[string]any{},
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

	_, err = client.Generate(context.Background(), llm.GenerateRequest{
		UserQuestion: "Will this fail?",
	})
	if err == nil {
		t.Fatal("expected error on interaction status=failed, got nil")
	}
	if !strings.Contains(err.Error(), "status=failed") {
		t.Fatalf("expected error mentioning status=failed, got: %v", err)
	}
}

func TestInteractionsOversizedResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Write 600 KiB response (exceeding 512 KiB limit)
		hugePayload := strings.Repeat("A", 600*1024)
		_, _ = io.WriteString(w, hugePayload)
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
		UserQuestion: "Oversized test",
	})
	if err == nil {
		t.Fatal("expected error on oversized response, got nil")
	}
	if !errors.Is(err, gemini.ErrOversizedResponse) && !strings.Contains(err.Error(), "exceeded maximum allowed size") {
		t.Fatalf("expected ErrOversizedResponse, got: %v", err)
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
