package groq

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
	defaultBaseURL     = "https://api.groq.com/openai/v1"
	chatCompletionsURI = "/chat/completions"
	defaultModel       = "openai/gpt-oss-20b"
	maxResponseBytes   = 1024 * 1024 // 1 MiB
)

var (
	// ErrOversizedResponse is returned when Groq response exceeds maximum allowed size.
	ErrOversizedResponse = errors.New("groq response exceeded maximum allowed size of 1 MiB")
	// ErrMissingAPIKey is returned when the API key is empty.
	ErrMissingAPIKey = errors.New("groq api key is required")
)

// Client implements llm.LLMProvider using Groq's OpenAI-compatible chat completions API.
type Client struct {
	apiKey          string
	model           string
	reasoningEffort string
	baseURL         string
	httpClient      *http.Client
}

// Option configures the Groq client.
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

// WithReasoningEffort sets the reasoning_effort parameter ("low", "medium", "high").
func WithReasoningEffort(effort string) Option {
	return func(c *Client) {
		c.reasoningEffort = normalizeReasoningEffort(effort)
	}
}

func normalizeReasoningEffort(effort string) string {
	switch strings.ToLower(strings.TrimSpace(effort)) {
	case "low", "medium", "high":
		return strings.ToLower(strings.TrimSpace(effort))
	default:
		return "low"
	}
}

// NewClient constructs a new Groq LLMProvider.
func NewClient(apiKey, model string, opts ...Option) (*Client, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, ErrMissingAPIKey
	}
	if strings.TrimSpace(model) == "" {
		model = defaultModel
	}

	c := &Client{
		apiKey:          strings.TrimSpace(apiKey),
		model:           strings.TrimSpace(model),
		reasoningEffort: "low",
		baseURL:         defaultBaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	c.reasoningEffort = normalizeReasoningEffort(c.reasoningEffort)
	return c, nil
}

func (c *Client) ProviderName() string {
	return "groq"
}

func (c *Client) Model() string {
	return c.model
}

func (c *Client) ReasoningEffort() string {
	return c.reasoningEffort
}

type httpStatusError struct {
	StatusCode int
	Body       string
}

func (e *httpStatusError) Error() string {
	return fmt.Sprintf("groq api returned status %d: %s", e.StatusCode, e.Body)
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type jsonSchemaField struct {
	Name   string         `json:"name"`
	Strict bool           `json:"strict"`
	Schema map[string]any `json:"schema"`
}

type responseFormatField struct {
	Type       string          `json:"type"`
	JSONSchema jsonSchemaField `json:"json_schema"`
}

type chatCompletionRequest struct {
	Model           string              `json:"model"`
	Messages        []chatMessage       `json:"messages"`
	Stream          bool                `json:"stream"`
	ReasoningEffort string              `json:"reasoning_effort,omitempty"`
	ResponseFormat  responseFormatField `json:"response_format"`
}

type chatCompletionResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

// Generate executes structured evidence generation using Groq.
func (c *Client) Generate(ctx context.Context, req llm.GenerateRequest) (llm.GeneratedAnswer, error) {
	endpoint := c.baseURL + chatCompletionsURI

	// Strict JSON schema representing llm.GeneratedAnswer
	schema := map[string]any{
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
					"required":             []string{"text", "evidence_ids"},
					"additionalProperties": false,
				},
			},
			"suggested_action_ids": map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "string",
				},
			},
		},
		"required":             []string{"status", "segments", "suggested_action_ids"},
		"additionalProperties": false,
	}

	messages := []chatMessage{
		{
			Role:    "system",
			Content: req.SystemInstruction,
		},
		{
			Role:    "user",
			Content: req.UserQuestion,
		},
	}

	payload := chatCompletionRequest{
		Model:           c.model,
		Messages:        messages,
		Stream:          false,
		ReasoningEffort: c.reasoningEffort,
		ResponseFormat: responseFormatField{
			Type: "json_schema",
			JSONSchema: jsonSchemaField{
				Name:   "ask_arham_grounded_answer",
				Strict: true,
				Schema: schema,
			},
		},
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return llm.GeneratedAnswer{}, fmt.Errorf("failed to marshal groq chat completion request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return llm.GeneratedAnswer{}, fmt.Errorf("failed to create groq http request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return llm.GeneratedAnswer{}, fmt.Errorf("groq network request failed: %w", err)
	}
	defer resp.Body.Close()

	// Bound response size with overflow detection
	lr := io.LimitReader(resp.Body, int64(maxResponseBytes+1))
	respBody, err := io.ReadAll(lr)
	if err != nil {
		return llm.GeneratedAnswer{}, fmt.Errorf("failed to read groq response body: %w", err)
	}
	if len(respBody) > maxResponseBytes {
		return llm.GeneratedAnswer{}, ErrOversizedResponse
	}

	if resp.StatusCode != http.StatusOK {
		// Bound error snippet to avoid unbounded or sensitive logging
		errSnippet := string(respBody)
		if len(errSnippet) > 500 {
			errSnippet = errSnippet[:500] + "... (truncated)"
		}
		return llm.GeneratedAnswer{}, &httpStatusError{
			StatusCode: resp.StatusCode,
			Body:       errSnippet,
		}
	}

	return parseChatCompletionResponse(respBody)
}

func parseChatCompletionResponse(body []byte) (llm.GeneratedAnswer, error) {
	var respObj chatCompletionResponse
	if err := json.Unmarshal(body, &respObj); err != nil {
		return llm.GeneratedAnswer{}, fmt.Errorf("failed to decode groq chat completion response: %w", err)
	}

	if len(respObj.Choices) == 0 {
		return llm.GeneratedAnswer{}, errors.New("groq response contained zero choices")
	}

	content := strings.TrimSpace(respObj.Choices[0].Message.Content)
	if content == "" {
		return llm.GeneratedAnswer{}, errors.New("groq response choice content is empty")
	}

	return parseJSONContent(content)
}

func parseJSONContent(raw string) (llm.GeneratedAnswer, error) {
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
		return llm.GeneratedAnswer{}, fmt.Errorf("invalid json in groq message content: %w", err)
	}

	if ans.Status == "" {
		return llm.GeneratedAnswer{}, errors.New("model output missing status field")
	}

	// Validate status against allowed values
	switch ans.Status {
	case "supported", "insufficient_evidence", "privacy_refusal", "scope_refusal":
		// valid
	default:
		return llm.GeneratedAnswer{}, fmt.Errorf("unknown answer status: %s", ans.Status)
	}

	// Ensure segments is non-nil slice
	if ans.Segments == nil {
		ans.Segments = []llm.GroundedSegment{}
	}
	if ans.SuggestedActionIDs == nil {
		ans.SuggestedActionIDs = []string{}
	}

	return ans, nil
}
