package embedding

import (
	"context"
)

// DisabledProvider represents a no-op embedding provider for EMBEDDING_MODE=disabled.
type DisabledProvider struct {
	dimensions int
}

func NewDisabledProvider(dimensions int) *DisabledProvider {
	if dimensions <= 0 {
		dimensions = 768
	}
	return &DisabledProvider{dimensions: dimensions}
}

func (p *DisabledProvider) Embed(ctx context.Context, purpose Purpose, texts []string) ([][]float32, error) {
	return nil, ErrEmbeddingDisabled
}

func (p *DisabledProvider) Dimensions() int {
	return p.dimensions
}

func (p *DisabledProvider) Model() string {
	return "disabled"
}

func (p *DisabledProvider) ProviderName() string {
	return "disabled"
}
