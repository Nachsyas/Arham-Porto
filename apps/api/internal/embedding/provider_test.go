package embedding

import (
	"context"
	"math"
	"testing"
)

func TestValidateVector(t *testing.T) {
	t.Run("valid vector 768 dimensions", func(t *testing.T) {
		vec := make([]float32, 768)
		for i := range vec {
			vec[i] = 0.5
		}
		if err := ValidateVector(vec, 768); err != nil {
			t.Fatalf("expected valid vector, got: %v", err)
		}
	})

	t.Run("short vector rejected", func(t *testing.T) {
		vec := make([]float32, 767)
		if err := ValidateVector(vec, 768); err == nil {
			t.Fatalf("expected error for short vector, got nil")
		}
	})

	t.Run("long vector rejected", func(t *testing.T) {
		vec := make([]float32, 769)
		if err := ValidateVector(vec, 768); err == nil {
			t.Fatalf("expected error for long vector, got nil")
		}
	})

	t.Run("NaN vector rejected", func(t *testing.T) {
		vec := make([]float32, 768)
		vec[42] = float32(math.NaN())
		if err := ValidateVector(vec, 768); err == nil {
			t.Fatalf("expected error for NaN vector, got nil")
		}
	})

	t.Run("Inf vector rejected", func(t *testing.T) {
		vec := make([]float32, 768)
		vec[100] = float32(math.Inf(1))
		if err := ValidateVector(vec, 768); err == nil {
			t.Fatalf("expected error for Inf vector, got nil")
		}
	})
}

func TestDeterministicFakeProvider(t *testing.T) {
	provider := NewDeterministicFakeProvider(768)
	ctx := context.Background()

	vecs, err := provider.Embed(ctx, PurposeDocument, []string{"hello", "world"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(vecs) != 2 {
		t.Fatalf("expected 2 vectors, got %d", len(vecs))
	}
	if len(vecs[0]) != 768 || len(vecs[1]) != 768 {
		t.Fatalf("expected 768 dimensions")
	}
	if provider.GetLastPurpose() != PurposeDocument {
		t.Fatalf("expected PurposeDocument")
	}

	// Verify determinism
	vecs2, err := provider.Embed(ctx, PurposeDocument, []string{"hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vecs[0][0] != vecs2[0][0] {
		t.Fatalf("expected deterministic output for identical input")
	}
}
