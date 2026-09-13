package embedding

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"math"
	"sync"
)

// DeterministicFakeProvider produces reproducible unit-normalized vectors for CI and structural tests (Gate #34).
// It cannot and must not be used to claim production semantic retrieval quality.
type DeterministicFakeProvider struct {
	mu           sync.Mutex
	dimensions   int
	model        string
	providerName string
	lastPurpose  Purpose
	callCount    int
}

func NewDeterministicFakeProvider(dimensions int) *DeterministicFakeProvider {
	if dimensions <= 0 {
		dimensions = 768
	}
	return &DeterministicFakeProvider{
		dimensions:   dimensions,
		model:        "deterministic-fake-v1",
		providerName: "deterministic_fake",
	}
}

func (p *DeterministicFakeProvider) Embed(ctx context.Context, purpose Purpose, texts []string) ([][]float32, error) {
	p.mu.Lock()
	p.lastPurpose = purpose
	p.callCount++
	p.mu.Unlock()

	results := make([][]float32, len(texts))
	for i, text := range texts {
		vec := make([]float32, p.dimensions)

		// Seed generator with SHA256 of text + purpose
		h := sha256.New()
		h.Write([]byte(string(purpose)))
		h.Write([]byte(":"))
		h.Write([]byte(text))
		seed := h.Sum(nil)

		var sumSq float64
		for d := 0; d < p.dimensions; d++ {
			byteIdx := (d * 4) % len(seed)
			var seedVal uint32
			if byteIdx+4 <= len(seed) {
				seedVal = binary.BigEndian.Uint32(seed[byteIdx : byteIdx+4])
			} else {
				seedVal = uint32(seed[byteIdx]) + uint32(d)
			}

			// Map to [-1.0, 1.0]
			val := float32(float64(seedVal%10000)/5000.0 - 1.0)
			vec[d] = val
			sumSq += float64(val * val)
		}

		// L2 unit normalization for exact cosine metric
		norm := math.Sqrt(sumSq)
		if norm > 0 {
			for d := 0; d < p.dimensions; d++ {
				vec[d] = float32(float64(vec[d]) / norm)
			}
		}

		if err := ValidateVector(vec, p.dimensions); err != nil {
			return nil, err
		}

		results[i] = vec
	}

	return results, nil
}

func (p *DeterministicFakeProvider) Dimensions() int {
	return p.dimensions
}

func (p *DeterministicFakeProvider) Model() string {
	return p.model
}

func (p *DeterministicFakeProvider) ProviderName() string {
	return p.providerName
}

func (p *DeterministicFakeProvider) GetLastPurpose() Purpose {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.lastPurpose
}

func (p *DeterministicFakeProvider) GetCallCount() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.callCount
}
