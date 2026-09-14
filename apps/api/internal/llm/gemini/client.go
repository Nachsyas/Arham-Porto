package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/nachsyas/arham-porto/apps/api/internal/llm"
)

const (
	defaultBaseURL   = "https://generativelanguage.googleapis.com/v1beta"
	interactionsPath = "/interactions"
	maxResponseBytes = 512 * 1024 // 512 KiB
)

// ErrOversizedResponse is returned when the provider response exceeds the read budget.
var ErrOversizedResponse = errors.New("gemini interactions response exceeded maximum allowed size of 512 KiB")

// Client implements llm.LLMProvider using the Google Gemini Interactions API.
type Client struct {
	apiKey        string
	model         string
	thinkingLevel string
	baseURL       string
	httpClient    *http.Client
}

// Option configures the Gemini client.
type Option func(*Client)

// WithBaseURL overrides the base URL (useful for httptest.Server).
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = strings.TrimRight(url, "/")
	}
}

// WithHTTPClient overrides the HTTP client.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.httpClient = client
	}
}

// WithThinkingLevel sets the thinking level (low, medium, high).
func WithThinkingLevel(level string) Option {
	return func(c *Client) {
		c.thinkingLevel = normalizeThinkingLevel(level)
	}
}

func normalizeThinkingLevel(level string) string {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "low", "medium", "high":
		return strings.ToLower(strings.TrimSpace(level))
	default:
		return "low"
	}
}

// NewClient constructs a new Gemini LLMProvider using the Interactions API.
func NewClient(apiKey, model string, opts ...Option) (*Client, error) {
	if apiKey == "" {
		return nil, errors.New("gemini api key is required")
	}
	if model == "" {
		model = "gemini-3.8-flash"
	}

	c := &Client{
		apiKey:        apiKey,
		model:         model,
		thinkingLevel: "low",
		baseURL:       defaultBaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	c.thinkingLevel = normalizeThinkingLevel(c.thinkingLevel)
	return c, nil
}

func (c *Client) ProviderName() string {
	return "gemini"
}

func (c *Client) Model() string {
	return c.model
}

func (c *Client) ThinkingLevel() string {
	return c.thinkingLevel
}

// Generate executes structured evidence generation using the Gemini Interactions API.
// Automatic legacy fallback to generateContent has been removed per Phase 6 Gate Rule #5.
func (c *Client) Generate(ctx context.Context, req llm.GenerateRequest) (llm.GeneratedAnswer, error) {
	return c.callInteractionsAPI(ctx, req)
}

type httpStatusError struct {
	StatusCode int
	Body       string
}

func (e *httpStatusError) Error() string {
	return fmt.Sprintf("gemini api returned status %d: %s", e.StatusCode, e.Body)
}

// callInteractionsAPI calls POST /v1beta/interactions with the current 2026 schema.
func (c *Client) callInteractionsAPI(ctx context.Context, req llm.GenerateRequest) (llm.GeneratedAnswer, error) {
	endpoint := c.baseURL + interactionsPath

	// Top-level response_format with structured JSON schema per Correction 1
	payload := map[string]any{
		"model":              c.model,
		"system_instruction": req.SystemInstruction,
		"input":              req.UserQuestion,
		"generation_config": map[string]any{
			"thinking_level":     c.thinkingLevel,
			"thinking_summaries": "none",
		},
		"response_format": map[string]any{
			"type":      "text",
			"mime_type": "application/json",
			"schema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"status": map[string]any{
						"type": "string",
						"enum": []string{
							"supported",
							"insufficient_evidence",
							"privacy_refusal",
							"scope_refusal",
						},
					},
					"segments": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"text": map[string]any{
									"type": "string",
								},
								"evidence_ids": map[string]any{
									"type": "array",
									"items": map[string]any{
										"type": "string",
									},
								},
							},
							"required": []string{"text", "evidence_ids"},
						},
					},
					"suggested_action_ids": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "string",
						},
					},
				},
				"required": []string{"status", "segments", "suggested_action_ids"},
			},
		},
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return llm.GeneratedAnswer{}, fmt.Errorf("failed to marshal interactions request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return llm.GeneratedAnswer{}, fmt.Errorf("failed to create interactions request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-goog-api-key", c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return llm.GeneratedAnswer{}, fmt.Errorf("interactions network request failed: %w", err)
	}
	defer resp.Body.Close()

	// Bound response size with overflow detection per Correction 19
	lr := io.LimitReader(resp.Body, int64(maxResponseBytes+1))
	respBody, err := io.ReadAll(lr)
	if err != nil {
		return llm.GeneratedAnswer{}, fmt.Errorf("failed to read interactions response: %w", err)
	}
	if len(respBody) > maxResponseBytes {
		return llm.GeneratedAnswer{}, ErrOversizedResponse
	}

	if resp.StatusCode != http.StatusOK {
		return llm.GeneratedAnswer{}, &httpStatusError{
			StatusCode: resp.StatusCode,
			Body:       string(respBody),
		}
	}

	return parseInteractionsResponse(respBody)
}

// interactionsResponse reflects the current 2026 Gemini Interactions API response structure.
type interactionsResponse struct {
	Status string `json:"status"` // "completed", "failed", "cancelled", "incomplete"
	Steps  []struct {
		Type    string `json:"type"` // "model_output"
		Content []struct {
			Type string `json:"type"` // "text"
			Text string `json:"text"`
		} `json:"content"`
	} `json:"steps"`
}

// parseInteractionsResponse extracts the model output from the steps array.
func parseInteractionsResponse(body []byte) (llm.GeneratedAnswer, error) {
	var respObj interactionsResponse
	if err := json.Unmarshal(body, &respObj); err != nil {
		return llm.GeneratedAnswer{}, fmt.Errorf("failed to decode interactions response envelope: %w", err)
	}

	// Verify interaction status per Correction 18
	switch respObj.Status {
	case "completed":
		// Normal successful completion
	case "failed", "cancelled", "incomplete":
		return llm.GeneratedAnswer{}, fmt.Errorf("gemini interaction not completed: status=%s", respObj.Status)
	default:
		if respObj.Status != "" {
			return llm.GeneratedAnswer{}, fmt.Errorf("gemini interaction returned unexpected status: %s", respObj.Status)
		}
	}

	// Search steps for model_output -> text per Correction 3
	for _, step := range respObj.Steps {
		if step.Type == "model_output" {
			for _, content := range step.Content {
				if content.Type == "text" && strings.TrimSpace(content.Text) != "" {
					return parseJSONText(content.Text)
				}
			}
		}
	}

	return llm.GeneratedAnswer{}, fmt.Errorf("gemini interaction response contained no model_output text step")
}

func parseJSONText(raw string) (llm.GeneratedAnswer, error) {
	trimmed := strings.TrimSpace(raw)
	if strings.HasPrefix(trimmed, "```json") {
		trimmed = strings.TrimPrefix(trimmed, "```json")
		trimmed = strings.TrimSuffix(trimmed, "```")
		trimmed = strings.TrimSpace(trimmed)
	} else if strings.HasPrefix(trimmed, "```") {
		trimmed = strings.TrimPrefix(trimmed, "```")
		trimmed = strings.TrimSuffix(trimmed, "```")
		trimmed = strings.TrimSpace(trimmed)
	}

	var ans llm.GeneratedAnswer
	if err := json.Unmarshal([]byte(trimmed), &ans); err != nil {
		return llm.GeneratedAnswer{}, fmt.Errorf("invalid json in model text: %w", err)
	}

	if ans.Status == "" {
		return llm.GeneratedAnswer{}, errors.New("model output missing status field")
	}

	return ans, nil
}
