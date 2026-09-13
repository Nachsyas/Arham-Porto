package llm

import (
	"context"
	"sync"
)

// DeterministicFakeProvider provides predictable responses for unit tests and local harness.
// It NEVER makes real network calls or invokes external AI APIs.
type DeterministicFakeProvider struct {
	mu           sync.Mutex
	customAnswer *GeneratedAnswer
	customErr    error
	calls        []GenerateRequest
}

// NewDeterministicFakeProvider creates a new fake provider.
func NewDeterministicFakeProvider() *DeterministicFakeProvider {
	return &DeterministicFakeProvider{
		calls: make([]GenerateRequest, 0),
	}
}

// SetAnswer overrides the next returned answer.
func (p *DeterministicFakeProvider) SetAnswer(ans GeneratedAnswer) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.customAnswer = &ans
}

// SetError overrides the next returned error.
func (p *DeterministicFakeProvider) SetError(err error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.customErr = err
}

// Calls returns a copy of all requests received.
func (p *DeterministicFakeProvider) Calls() []GenerateRequest {
	p.mu.Lock()
	defer p.mu.Unlock()
	copied := make([]GenerateRequest, len(p.calls))
	copy(copied, p.calls)
	return copied
}

// Generate returns the configured deterministic answer or a default grounded answer.
func (p *DeterministicFakeProvider) Generate(ctx context.Context, req GenerateRequest) (GeneratedAnswer, error) {
	select {
	case <-ctx.Done():
		return GeneratedAnswer{}, ctx.Err()
	default:
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	p.calls = append(p.calls, req)

	if p.customErr != nil {
		return GeneratedAnswer{}, p.customErr
	}

	if p.customAnswer != nil {
		return *p.customAnswer, nil
	}

	// Default behavior: if evidence provided, return supported segments
	if len(req.Evidence) > 0 {
		return GeneratedAnswer{
			Status: "supported",
			Segments: []GroundedSegment{
				{
					Text:        "Nachsyas Arham Mumtaz Nashohi demonstrated engineering expertise in this area.",
					EvidenceIDs: []string{req.Evidence[0].ID},
				},
			},
			SuggestedActionIDs: []string{"go-to-projects"},
		}, nil
	}

	// No evidence -> insufficient evidence refusal
	return GeneratedAnswer{
		Status: "insufficient_evidence",
		Segments: []GroundedSegment{
			{
				Text:        "I couldn't find verified evidence of that in Nachsyas Arham Mumtaz Nashohi's approved portfolio sources.",
				EvidenceIDs: []string{},
			},
		},
		SuggestedActionIDs: []string{"go-to-projects"},
	}, nil
}

// ProviderName returns the provider identifier.
func (p *DeterministicFakeProvider) ProviderName() string {
	return "fake"
}

// Model returns the model identifier.
func (p *DeterministicFakeProvider) Model() string {
	return "deterministic-test-fake"
}
