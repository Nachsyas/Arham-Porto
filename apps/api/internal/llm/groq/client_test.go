package groq

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/nachsyas/arham-porto/apps/api/internal/llm"
)

func TestNewClient(t *testing.T) {
	t.Run("missing api key", func(t *testing.T) {
		_, err := NewClient("", "openai/gpt-oss-20b")
		if err == nil {
			t.Fatal("expected error for missing api key")
		}
	})

	t.Run("default model and provider name", func(t *testing.T) {
		client, err := NewClient("gsk_test123", "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if client.ProviderName() != "groq" {
			t.Errorf("expected provider 'groq', got %q", client.ProviderName())
		}
		if client.Model() != "openai/gpt-oss-20b" {
			t.Errorf("expected default model 'openai/gpt-oss-20b', got %q", client.Model())
		}
		if client.ReasoningEffort() != "low" {
			t.Errorf("expected default reasoning effort 'low', got %q", client.ReasoningEffort())
		}
	})

	t.Run("custom options", func(t *testing.T) {
		client, err := NewClient("gsk_test123", "custom-model",
			WithReasoningEffort("medium"),
			WithBaseURL("https://example.com/v1/"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if client.Model() != "custom-model" {
			t.Errorf("expected model 'custom-model', got %q", client.Model())
		}
		if client.ReasoningEffort() != "medium" {
			t.Errorf("expected reasoning effort 'medium', got %q", client.ReasoningEffort())
		}
	})
}

func TestGenerate_Success(t *testing.T) {
	expectedAnswer := llm.GeneratedAnswer{
		Status: "supported",
		Segments: []llm.GroundedSegment{
			{
				Text:        "Nachsyas Arham Mumtaz Nashohi specializes in backend engineering with Go.",
				EvidenceIDs: []string{"E1", "E2"},
			},
		},
		SuggestedActionIDs: []string{"go-to-projects"},
	}

	contentBytes, err := json.Marshal(expectedAnswer)
	if err != nil {
		t.Fatalf("failed to marshal expected answer: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/chat/completions" {
			t.Errorf("expected /chat/completions, got %s", r.URL.Path)
		}
		if auth := r.Header.Get("Authorization"); auth != "Bearer gsk_valid_key" {
			t.Errorf("expected Bearer gsk_valid_key, got %q", auth)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json, got %q", ct)
		}

		var reqBody chatCompletionRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		if reqBody.Model != "openai/gpt-oss-20b" {
			t.Errorf("expected model openai/gpt-oss-20b, got %s", reqBody.Model)
		}
		if reqBody.ReasoningEffort != "low" {
			t.Errorf("expected reasoning_effort low, got %s", reqBody.ReasoningEffort)
		}
		if reqBody.Stream {
			t.Errorf("expected stream false")
		}
		if len(reqBody.Messages) != 2 {
			t.Fatalf("expected 2 messages (system, user), got %d", len(reqBody.Messages))
		}
		if reqBody.Messages[0].Role != "system" || reqBody.Messages[0].Content != "System Prompt" {
			t.Errorf("unexpected system message: %+v", reqBody.Messages[0])
		}
		if reqBody.Messages[1].Role != "user" || reqBody.Messages[1].Content != "User Question" {
			t.Errorf("unexpected user message: %+v", reqBody.Messages[1])
		}
		if reqBody.ResponseFormat.Type != "json_schema" {
			t.Errorf("expected response_format type json_schema, got %s", reqBody.ResponseFormat.Type)
		}
		if !reqBody.ResponseFormat.JSONSchema.Strict {
			t.Errorf("expected strict json_schema")
		}

		resp := chatCompletionResponse{
			ID: "chatcmpl-123",
			Choices: []struct {
				Index   int `json:"index"`
				Message struct {
					Role    string `json:"role"`
					Content string `json:"content"`
				} `json:"message"`
				FinishReason string `json:"finish_reason"`
			}{
				{
					Index: 0,
					Message: struct {
						Role    string `json:"role"`
						Content string `json:"content"`
					}{
						Role:    "assistant",
						Content: string(contentBytes),
					},
					FinishReason: "stop",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, err := NewClient("gsk_valid_key", "openai/gpt-oss-20b", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ans, err := client.Generate(context.Background(), llm.GenerateRequest{
		SystemInstruction: "System Prompt",
		UserQuestion:      "User Question",
	})
	if err != nil {
		t.Fatalf("unexpected generation error: %v", err)
	}

	if ans.Status != "supported" {
		t.Errorf("expected status 'supported', got %q", ans.Status)
	}
	if len(ans.Segments) != 1 {
		t.Fatalf("expected 1 segment, got %d", len(ans.Segments))
	}
	if ans.Segments[0].Text != expectedAnswer.Segments[0].Text {
		t.Errorf("segment text mismatch: got %q", ans.Segments[0].Text)
	}
	if len(ans.Segments[0].EvidenceIDs) != 2 {
		t.Errorf("evidence IDs mismatch: got %+v", ans.Segments[0].EvidenceIDs)
	}
}

func TestGenerate_Errors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		respBody   string
		expectErr  string
	}{
		{
			name:       "zero choices",
			statusCode: http.StatusOK,
			respBody:   `{"id":"chat-1","choices":[]}`,
			expectErr:  "zero choices",
		},
		{
			name:       "empty content",
			statusCode: http.StatusOK,
			respBody:   `{"id":"chat-1","choices":[{"message":{"role":"assistant","content":""}}]}`,
			expectErr:  "content is empty",
		},
		{
			name:       "malformed json content",
			statusCode: http.StatusOK,
			respBody:   `{"id":"chat-1","choices":[{"message":{"role":"assistant","content":"not-json"}}]}`,
			expectErr:  "invalid json in groq message content",
		},
		{
			name:       "missing status field",
			statusCode: http.StatusOK,
			respBody:   `{"id":"chat-1","choices":[{"message":{"role":"assistant","content":"{\"segments\":[]}"}}]}`,
			expectErr:  "missing status field",
		},
		{
			name:       "unknown status",
			statusCode: http.StatusOK,
			respBody:   `{"id":"chat-1","choices":[{"message":{"role":"assistant","content":"{\"status\":\"invalid_status\",\"segments\":[]}"}}]}`,
			expectErr:  "unknown answer status",
		},
		{
			name:       "provider 400 error",
			statusCode: http.StatusBadRequest,
			respBody:   `{"error":{"message":"bad request"}}`,
			expectErr:  "returned status 400",
		},
		{
			name:       "provider 401 unauthorized",
			statusCode: http.StatusUnauthorized,
			respBody:   `{"error":{"message":"invalid api key"}}`,
			expectErr:  "returned status 401",
		},
		{
			name:       "provider 429 rate limited",
			statusCode: http.StatusTooManyRequests,
			respBody:   `{"error":{"message":"rate limit reached"}}`,
			expectErr:  "returned status 429",
		},
		{
			name:       "provider 500 internal server error",
			statusCode: http.StatusInternalServerError,
			respBody:   `{"error":{"message":"internal error"}}`,
			expectErr:  "returned status 500",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
				io.WriteString(w, tc.respBody)
			}))
			defer server.Close()

			client, err := NewClient("gsk_key", "openai/gpt-oss-20b", WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			_, err = client.Generate(context.Background(), llm.GenerateRequest{
				UserQuestion: "test",
			})
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tc.expectErr)
			}
			if !strings.Contains(err.Error(), tc.expectErr) {
				t.Errorf("expected error to contain %q, got %q", tc.expectErr, err.Error())
			}
		})
	}
}

func TestGenerate_OversizedResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// Write 1.5MB of data
		padding := strings.Repeat("x", 1024)
		for i := 0; i < 1500; i++ {
			io.WriteString(w, padding)
		}
	}))
	defer server.Close()

	client, err := NewClient("gsk_key", "openai/gpt-oss-20b", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Generate(context.Background(), llm.GenerateRequest{UserQuestion: "test"})
	if err == nil {
		t.Fatal("expected error for oversized response")
	}
	if !strings.Contains(err.Error(), "exceeded maximum allowed size") {
		t.Errorf("expected oversized response error, got %v", err)
	}
}

func TestGenerate_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClient("gsk_key", "openai/gpt-oss-20b", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	_, err = client.Generate(ctx, llm.GenerateRequest{UserQuestion: "test"})
	if err == nil {
		t.Fatal("expected error on cancelled context")
	}
}

func TestGenerate_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClient("gsk_key", "openai/gpt-oss-20b",
		WithBaseURL(server.URL),
		WithHTTPClient(&http.Client{Timeout: 50 * time.Millisecond}),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Generate(context.Background(), llm.GenerateRequest{UserQuestion: "test"})
	if err == nil {
		t.Fatal("expected timeout error")
	}
}
