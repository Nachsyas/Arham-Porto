package llm_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nachsyas/arham-porto/apps/api/internal/llm"
)

func TestDisabledProvider(t *testing.T) {
	provider := llm.NewDisabledProvider()
	if provider.ProviderName() != "disabled" {
		t.Fatalf("expected provider name 'disabled', got %s", provider.ProviderName())
	}
	if provider.Model() != "none" {
		t.Fatalf("expected model 'none', got %s", provider.Model())
	}

	_, err := provider.Generate(context.Background(), llm.GenerateRequest{})
	if !errors.Is(err, llm.ErrAIFeatureDisabled) {
		t.Fatalf("expected ErrAIFeatureDisabled, got %v", err)
	}
}

func TestDeterministicFakeProvider(t *testing.T) {
	provider := llm.NewDeterministicFakeProvider()
	if provider.ProviderName() != "fake" {
		t.Fatalf("expected provider name 'fake', got %s", provider.ProviderName())
	}

	// 1. Default with evidence -> supported
	reqWithEvidence := llm.GenerateRequest{
		UserQuestion: "What demonstrates Go experience?",
		Evidence: []llm.EvidenceContext{
			{ID: "E1", Title: "GDGOC E-Commerce", Content: "Go backend with Clean Architecture"},
		},
	}
	ans, err := provider.Generate(context.Background(), reqWithEvidence)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ans.Status != "supported" {
		t.Fatalf("expected status supported, got %s", ans.Status)
	}
	if len(ans.Segments) == 0 || len(ans.Segments[0].EvidenceIDs) == 0 || ans.Segments[0].EvidenceIDs[0] != "E1" {
		t.Fatalf("expected segment with evidence E1, got %+v", ans.Segments)
	}

	// 2. Default with no evidence -> insufficient_evidence
	reqNoEvidence := llm.GenerateRequest{
		UserQuestion: "Does Arham know Kubernetes?",
		Evidence:     []llm.EvidenceContext{},
	}
	ansNoEv, err := provider.Generate(context.Background(), reqNoEvidence)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ansNoEv.Status != "insufficient_evidence" {
		t.Fatalf("expected status insufficient_evidence, got %s", ansNoEv.Status)
	}

	// 3. Custom answer override
	customAns := llm.GeneratedAnswer{
		Status: "privacy_refusal",
		Segments: []llm.GroundedSegment{
			{Text: "Private info protected.", EvidenceIDs: []string{}},
		},
		SuggestedActionIDs: []string{"go-to-contact"},
	}
	provider.SetAnswer(customAns)
	gotCustom, err := provider.Generate(context.Background(), reqWithEvidence)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotCustom.Status != "privacy_refusal" {
		t.Fatalf("expected custom status, got %s", gotCustom.Status)
	}

	// 4. Custom error override
	expectedErr := errors.New("simulated provider error")
	provider.SetError(expectedErr)
	_, err = provider.Generate(context.Background(), reqWithEvidence)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected simulated error, got %v", err)
	}

	// 5. Context cancellation
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	provider.SetError(nil)
	_, err = provider.Generate(ctx, reqWithEvidence)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}

	// 6. Request call recording
	calls := provider.Calls()
	if len(calls) < 4 {
		t.Fatalf("expected at least 4 recorded calls, got %d", len(calls))
	}
}

func TestContextTimeoutPropagation(t *testing.T) {
	provider := llm.NewDeterministicFakeProvider()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	time.Sleep(2 * time.Millisecond)

	_, err := provider.Generate(ctx, llm.GenerateRequest{UserQuestion: "test"})
	if !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled) {
		t.Fatalf("expected deadline exceeded, got %v", err)
	}
}
