package embedding

import (
	"context"
	"errors"
	"fmt"
	"math"
)

// Purpose defines the semantic role of an embedding vector per Gate #3.
type Purpose string

const (
	PurposeDocument Purpose = "document"
	PurposeQuery    Purpose = "query"
)

var (
	ErrInvalidDimensions = errors.New("vector dimension mismatch")
	ErrInvalidFloatValue = errors.New("vector contains NaN or Infinity")
	ErrBatchSizeMismatch = errors.New("embedding provider returned unexpected vector count")
	ErrEmbeddingDisabled = errors.New("embedding provider is disabled")
)

// Provider defines the provider-agnostic vector generation contract per Gate #3.
type Provider interface {
	Embed(ctx context.Context, purpose Purpose, texts []string) ([][]float32, error)
	Dimensions() int
	Model() string
	ProviderName() string
}

// ValidateVector strictly validates vector dimensionality and float validity per Gate #4.
// Rejects shorter, longer, NaN, or Inf vectors without truncation, padding, or silent coercion.
func ValidateVector(v []float32, expectedDimensions int) error {
	if len(v) != expectedDimensions {
		return fmt.Errorf("%w: expected %d, got %d", ErrInvalidDimensions, expectedDimensions, len(v))
	}

	for i, val := range v {
		f := float64(val)
		if math.IsNaN(f) || math.IsInf(f, 0) {
			return fmt.Errorf("%w at index %d: %v", ErrInvalidFloatValue, i, val)
		}
	}

	return nil
}
