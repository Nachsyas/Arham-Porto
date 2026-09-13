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

	"github.com/nachsyas/arham-porto/apps/api/internal/embedding"
)

const (
	DefaultModel      = "gemini-embedding-2"
	DefaultDimensions = 768
	defaultTimeout    = 30 * time.Second
	defaultEndpoint   = "https://generativelanguage.googleapis.com/v1beta"
)

// Provider implements embedding.Provider for Google Gemini embedding models (Gate #26, #27).
type Provider struct {
	apiKey     string
	model      string
	dimensions int
	endpoint   string
	httpClient *http.Client
}

// NewProvider creates a production Gemini embedding provider.
func NewProvider(apiKey, model string, dimensions int) (*Provider, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, errors.New("gemini api key cannot be empty")
	}
	if strings.TrimSpace(model) == "" {
		model = DefaultModel
	}
	if dimensions <= 0 {
		dimensions = DefaultDimensions
	}

	return &Provider{
		apiKey:     strings.TrimSpace(apiKey),
		model:      model,
		dimensions: dimensions,
		endpoint:   defaultEndpoint,
		httpClient: &http.Client{Timeout: defaultTimeout},
	}, nil
}

// NewTestProvider creates a Gemini provider with a custom endpoint for mock testing.
func NewTestProvider(endpoint, apiKey, model string, dimensions int) *Provider {
	if strings.TrimSpace(model) == "" {
		model = DefaultModel
	}
	if dimensions <= 0 {
		dimensions = DefaultDimensions
	}
	return &Provider{
		apiKey:     apiKey,
		model:      model,
		dimensions: dimensions,
		endpoint:   strings.TrimRight(endpoint, "/"),
		httpClient: &http.Client{Timeout: defaultTimeout},
	}
}

func (p *Provider) Dimensions() int {
	return p.dimensions
}

func (p *Provider) Model() string {
	return p.model
}

func (p *Provider) ProviderName() string {
	return "gemini"
}

// Request and response payloads for batchEmbedContents
type batchEmbedRequest struct {
	Requests []embedRequestItem `json:"requests"`
}

type embedRequestItem struct {
	Model                string       `json:"model"`
	Content              embedContent `json:"content"`
	TaskType             string       `json:"taskType,omitempty"`
	OutputDimensionality int          `json:"outputDimensionality,omitempty"`
}

type embedContent struct {
	Parts []embedPart `json:"parts"`
}

type embedPart struct {
	Text string `json:"text"`
}

type batchEmbedResponse struct {
	Embeddings []struct {
		Values []float32 `json:"values"`
	} `json:"embeddings"`
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error,omitempty"`
}

// Embed generates vector embeddings for a list of texts using explicit Purpose (Gate #3, #28).
func (p *Provider) Embed(ctx context.Context, purpose embedding.Purpose, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, nil
	}

	// Translate Purpose to Gemini-native task types
	taskType := "RETRIEVAL_DOCUMENT"
	if purpose == embedding.PurposeQuery {
		taskType = "RETRIEVAL_QUERY"
	}

	modelRef := p.model
	if !strings.HasPrefix(modelRef, "models/") {
		modelRef = "models/" + modelRef
	}

	reqItems := make([]embedRequestItem, len(texts))
	for i, t := range texts {
		reqItems[i] = embedRequestItem{
			Model: modelRef,
			Content: embedContent{
				Parts: []embedPart{{Text: t}},
			},
			TaskType:             taskType,
			OutputDimensionality: p.dimensions,
		}
	}

	reqPayload := batchEmbedRequest{Requests: reqItems}
	reqBytes, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal gemini embedding request: %w", err)
	}

	url := fmt.Sprintf("%s/%s:batchEmbedContents", p.endpoint, modelRef)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini embedding request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	// Use header for API key to avoid query param leaking in URLs (Gate #27)
	req.Header.Set("x-goog-api-key", p.apiKey)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gemini embedding request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 10*1024*1024))
	if err != nil {
		return nil, fmt.Errorf("failed to read gemini embedding response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gemini api returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var batchResp batchEmbedResponse
	if err := json.Unmarshal(respBody, &batchResp); err != nil {
		return nil, fmt.Errorf("failed to decode gemini embedding response: %w", err)
	}

	if batchResp.Error != nil {
		return nil, fmt.Errorf("gemini embedding error (%s): %s", batchResp.Error.Status, batchResp.Error.Message)
	}

	// Gate #28: Ensure N input documents -> exactly N embedding vectors
	if len(batchResp.Embeddings) != len(texts) {
		return nil, fmt.Errorf("%w: expected %d vectors, got %d", embedding.ErrBatchSizeMismatch, len(texts), len(batchResp.Embeddings))
	}

	results := make([][]float32, len(batchResp.Embeddings))
	for i, emb := range batchResp.Embeddings {
		// Gate #4: Validate returned vector dimensions
		if err := embedding.ValidateVector(emb.Values, p.dimensions); err != nil {
			return nil, fmt.Errorf("vector at index %d validation failed: %w", i, err)
		}
		results[i] = emb.Values
	}

	return results, nil
}
