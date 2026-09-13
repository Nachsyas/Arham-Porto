package llm

import (
	"context"
)

// DisabledProvider implements LLMProvider when AI_MODE=disabled.
type DisabledProvider struct{}

// NewDisabledProvider creates an instance of DisabledProvider.
func NewDisabledProvider() *DisabledProvider {
	return &DisabledProvider{}
}

// Generate immediately returns ErrAIFeatureDisabled without making network calls.
func (p *DisabledProvider) Generate(ctx context.Context, req GenerateRequest) (GeneratedAnswer, error) {
	return GeneratedAnswer{}, ErrAIFeatureDisabled
}

// ProviderName returns the provider identifier.
func (p *DisabledProvider) ProviderName() string {
	return "disabled"
}

// Model returns the model identifier.
func (p *DisabledProvider) Model() string {
	return "none"
}
