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
	defaultBaseURL     = "https://generativelanguage.googleapis.com/v1beta"
	interactionsPath   = "/interactions"
	generateContentFmt = "/models/%s:generateContent"
)

// Client implements llm.LLMProvider using Google Gemini.
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
		c.thinkingLevel = level
	}
}

// NewClient constructs a new Gemini LLMProvider.
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

	return c, nil
}

func (c *Client) ProviderName() string {
	return "gemini"
}

func (c *Client) Model() string {
	return c.model
}

// Generate executes structured evidence generation.
func (c *Client) Generate(ctx context.Context, req llm.GenerateRequest) (llm.GeneratedAnswer, error) {
	// Attempt 1: Gemini Interactions API
	ans, err := c.callInteractionsAPI(ctx, req)
	if err == nil {
		return ans, nil
	}

	// If interactions endpoint is not found or unsupported, fallback to generateContent
	var httpErr *httpStatusError
	if errors.As(err, &httpErr) && (httpErr.StatusCode == http.StatusNotFound || httpErr.StatusCode == http.StatusMethodNotAllowed) {
		return c.callGenerateContentAPI(ctx, req)
	}

	return ans, err
}

type httpStatusError struct {
	StatusCode int
	Body       string
}

func (e *httpStatusError) Error() string {
	return fmt.Sprintf("gemini api returned status %d", e.StatusCode)
}

// callInteractionsAPI calls POST /v1beta/interactions
func (c *Client) callInteractionsAPI(ctx context.Context, req llm.GenerateRequest) (llm.GeneratedAnswer, error) {
	endpoint := c.baseURL + interactionsPath

	fullPrompt := req.UserQuestion
	payload := map[string]any{
		"model":              c.model,
		"system_instruction": req.SystemInstruction,
		"input":              fullPrompt,
		"generation_config": map[string]any{
			"response_mime_type": "application/json",
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

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 512*1024))
	if err != nil {
		return llm.GeneratedAnswer{}, fmt.Errorf("failed to read interactions response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return llm.GeneratedAnswer{}, &httpStatusError{
			StatusCode: resp.StatusCode,
			Body:       string(respBody),
		}
	}

	return parseStructuredAnswer(respBody)
}

// callGenerateContentAPI calls POST /v1beta/models/{model}:generateContent as compatibility fallback
func (c *Client) callGenerateContentAPI(ctx context.Context, req llm.GenerateRequest) (llm.GeneratedAnswer, error) {
	modelPath := c.model
	if !strings.HasPrefix(modelPath, "models/") {
		modelPath = fmt.Sprintf(generateContentFmt, c.model)
	}
	endpoint := c.baseURL + modelPath

	payload := map[string]any{
		"system_instruction": map[string]any{
			"parts": []map[string]any{
				{"text": req.SystemInstruction},
			},
		},
		"contents": []map[string]any{
			{
				"role": "user",
				"parts": []map[string]any{
					{"text": req.UserQuestion},
				},
			},
		},
		"generationConfig": map[string]any{
			"responseMimeType": "application/json",
		},
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return llm.GeneratedAnswer{}, fmt.Errorf("failed to marshal generateContent request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return llm.GeneratedAnswer{}, fmt.Errorf("failed to create generateContent request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-goog-api-key", c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return llm.GeneratedAnswer{}, fmt.Errorf("generateContent network request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 512*1024))
	if err != nil {
		return llm.GeneratedAnswer{}, fmt.Errorf("failed to read generateContent response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return llm.GeneratedAnswer{}, &httpStatusError{
			StatusCode: resp.StatusCode,
			Body:       string(respBody),
		}
	}

	return parseStructuredAnswer(respBody)
}

// parseStructuredAnswer extracts the JSON output from either Gemini response structure
func parseStructuredAnswer(body []byte) (llm.GeneratedAnswer, error) {
	// 1. Try parsing directly if the response is directly the structured JSON or has output field
	var directAnswer llm.GeneratedAnswer
	if err := json.Unmarshal(body, &directAnswer); err == nil && directAnswer.Status != "" {
		return directAnswer, nil
	}

	// 2. Check for Interactions API response envelope {"output": "..."} or {"result": "..."}
	var envelope struct {
		Output string `json:"output"`
		Result string `json:"result"`
		Text   string `json:"text"`
	}
	if err := json.Unmarshal(body, &envelope); err == nil {
		raw := envelope.Output
		if raw == "" {
			raw = envelope.Result
		}
		if raw == "" {
			raw = envelope.Text
		}
		if raw != "" {
			if parsed, err := parseJSONText(raw); err == nil {
				return parsed, nil
			}
		}
	}

	// 3. Check for generateContent envelope {"candidates": [{"content": {"parts": [{"text": "..."}]}}]}
	var gcEnvelope struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.Unmarshal(body, &gcEnvelope); err == nil && len(gcEnvelope.Candidates) > 0 {
		for _, part := range gcEnvelope.Candidates[0].Content.Parts {
			if part.Text != "" {
				if parsed, err := parseJSONText(part.Text); err == nil {
					return parsed, nil
				}
			}
		}
	}

	return llm.GeneratedAnswer{}, fmt.Errorf("failed to parse structured model response from Gemini payload")
}

func parseJSONText(raw string) (llm.GeneratedAnswer, error) {
	trimmed := strings.TrimSpace(raw)
	// Strip markdown fences if present
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
	return ans, nil
}
