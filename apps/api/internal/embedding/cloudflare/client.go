package cloudflare

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

	"github.com/nachsyas/arham-porto/apps/api/internal/embedding"
)

const (
	DefaultModel      = "@cf/baai/bge-base-en-v1.5"
	DefaultDimensions = 768
	defaultTimeout    = 30 * time.Second
	defaultBaseURL   = "https://api.cloudflare.com/client/v4"
	maxResponseBytes  = 5 * 1024 * 1024 // 5 MiB read budget
)

// ErrOversizedResponse is returned when the Cloudflare response exceeds the read budget.
var ErrOversizedResponse = errors.New("cloudflare workers ai response exceeded maximum allowed size of 5 MiB")

// Provider implements embedding.Provider using Cloudflare Workers AI.
type Provider struct {
	accountID  string
	apiToken   string
	model      string
	dimensions int
	baseURL    string
	httpClient *http.Client
}

// Option configures the Cloudflare provider.
type Option func(*Provider)

// WithBaseURL overrides the base URL (useful for httptest.Server).
func WithBaseURL(url string) Option {
	return func(p *Provider) {
		p.baseURL = strings.TrimRight(url, "/")
	}
}

// WithHTTPClient overrides the HTTP client.
func WithHTTPClient(client *http.Client) Option {
	return func(p *Provider) {
		p.httpClient = client
	}
}

// NewProvider creates a new Cloudflare Workers AI embedding provider.
func NewProvider(accountID, apiToken, model string, dimensions int, opts ...Option) (*Provider, error) {
	cleanAccountID := strings.TrimSpace(accountID)
	if cleanAccountID == "" {
		return nil, errors.New("cloudflare account id cannot be empty")
	}

	cleanToken := strings.TrimSpace(apiToken)
	if cleanToken == "" {
		return nil, errors.New("cloudflare api token cannot be empty")
	}

	cleanModel := strings.TrimSpace(model)
	if cleanModel == "" {
		cleanModel = DefaultModel
	}

	if dimensions <= 0 {
		dimensions = DefaultDimensions
	}

	p := &Provider{
		accountID:  cleanAccountID,
		apiToken:   cleanToken,
		model:      cleanModel,
		dimensions: dimensions,
		baseURL:    defaultBaseURL,
		httpClient: &http.Client{Timeout: defaultTimeout},
	}

	for _, opt := range opts {
		opt(p)
	}

	return p, nil
}

// ProviderName returns the identifier for this provider.
func (p *Provider) ProviderName() string {
	return "cloudflare"
}

// Model returns the configured model identifier.
func (p *Provider) Model() string {
	return p.model
}

// Dimensions returns the configured embedding vector dimension.
func (p *Provider) Dimensions() int {
	return p.dimensions
}

// Embed generates embedding vectors for a batch of input texts using Cloudflare Workers AI.
func (p *Provider) Embed(ctx context.Context, _ embedding.Purpose, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, nil
	}

	endpoint := fmt.Sprintf("%s/accounts/%s/ai/run/%s", p.baseURL, p.accountID, p.model)

	reqPayload := map[string]any{
		"text": texts,
	}

	jsonBytes, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, fmt.Errorf("failed marshaling embedding request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("failed creating embedding request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+p.apiToken)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("cloudflare embedding request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyReader := io.LimitReader(resp.Body, maxResponseBytes+1)
	body, err := io.ReadAll(bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed reading cloudflare response body: %w", err)
	}

	if len(body) > maxResponseBytes {
		return nil, ErrOversizedResponse
	}

	type cfErrorItem struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	type cfResponseEnvelope struct {
		Success  bool            `json:"success"`
		Errors   []cfErrorItem   `json:"errors"`
		Messages []string        `json:"messages"`
		Result   json.RawMessage `json:"result"`
	}

	var envelope cfResponseEnvelope
	if unmarshalErr := json.Unmarshal(body, &envelope); unmarshalErr != nil {
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("cloudflare workers ai returned status %d (malformed response body)", resp.StatusCode)
		}
		return nil, fmt.Errorf("failed unmarshaling cloudflare response envelope: %w", unmarshalErr)
	}

	if resp.StatusCode != http.StatusOK || !envelope.Success {
		var errMsgs []string
		for _, e := range envelope.Errors {
			if e.Message != "" {
				errMsgs = append(errMsgs, fmt.Sprintf("[%d] %s", e.Code, e.Message))
			}
		}
		details := strings.Join(errMsgs, "; ")
		if details == "" {
			details = "unknown provider error"
		}
		return nil, fmt.Errorf("cloudflare workers ai returned status %d: %s", resp.StatusCode, details)
	}

	vectors, err := parseEmbeddingResult(envelope.Result, len(texts), p.dimensions)
	if err != nil {
		return nil, fmt.Errorf("failed parsing embedding result: %w", err)
	}

	if len(vectors) != len(texts) {
		return nil, fmt.Errorf("%w: expected %d, got %d", embedding.ErrBatchSizeMismatch, len(texts), len(vectors))
	}

	for i, vec := range vectors {
		if err := embedding.ValidateVector(vec, p.dimensions); err != nil {
			return nil, fmt.Errorf("vector validation failed for item %d: %w", i, err)
		}
	}

	return vectors, nil
}

// parseEmbeddingResult decodes Cloudflare Workers AI result payloads defensively.
func parseEmbeddingResult(raw json.RawMessage, expectedCount, expectedDims int) ([][]float32, error) {
	if len(raw) == 0 {
		return nil, errors.New("empty result in response envelope")
	}

	// Format 1: Standard shape/data 2D envelope: {"shape": [N, 768], "data": [[...], [...]]}
	var res2D struct {
		Shape []int       `json:"shape"`
		Data  [][]float32 `json:"data"`
	}
	if err := json.Unmarshal(raw, &res2D); err == nil && len(res2D.Data) > 0 {
		return res2D.Data, nil
	}

	// Format 2: Shape/data 1D envelope when single text: {"shape": [768], "data": [...]}
	if expectedCount == 1 {
		var res1D struct {
			Shape []int     `json:"shape"`
			Data  []float32 `json:"data"`
		}
		if err := json.Unmarshal(raw, &res1D); err == nil && len(res1D.Data) == expectedDims {
			return [][]float32{res1D.Data}, nil
		}
	}

	// Format 3: Direct 2D array: [[...], [...]]
	var direct2D [][]float32
	if err := json.Unmarshal(raw, &direct2D); err == nil && len(direct2D) > 0 {
		return direct2D, nil
	}

	// Format 4: Direct 1D array when single text: [...]
	if expectedCount == 1 {
		var direct1D []float32
		if err := json.Unmarshal(raw, &direct1D); err == nil && len(direct1D) == expectedDims {
			return [][]float32{direct1D}, nil
		}
	}

	return nil, errors.New("unrecognized embedding result structure from cloudflare workers ai")
}

var _ embedding.Provider = (*Provider)(nil)
