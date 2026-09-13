package indexing

import (
	"strings"
	"testing"
)

func TestBuildCanonicalKnowledgeSources(t *testing.T) {
	// The repo root is 4 directories up from apps/api/internal/indexing
	dataDir := "../../../../data"

	sources, err := BuildCanonicalKnowledgeSources(dataDir)
	if err != nil {
		t.Fatalf("BuildCanonicalKnowledgeSources failed: %v", err)
	}

	if len(sources) == 0 {
		t.Fatalf("expected canonical sources, got 0")
	}

	sourceTypes := make(map[string]int)
	for _, s := range sources {
		sourceTypes[s.Source.SourceType]++

		// Check for forbidden fields
		for _, c := range s.Chunks {
			if strings.Contains(c.Content, "TODO") {
				t.Errorf("chunk content leaked TODO field: %s", c.Content)
			}
			if strings.Contains(c.Content, "coordinates") || strings.Contains(c.Content, "-7.59") {
				t.Errorf("chunk content leaked geographic coordinates: %s", c.Content)
			}
		}
	}

	expectedTypes := []string{
		"canonical_profile",
		"canonical_project",
		"canonical_skill",
		"canonical_evidence",
		"canonical_journey",
	}

	for _, et := range expectedTypes {
		if sourceTypes[et] == 0 {
			t.Errorf("expected at least 1 source of type %s, found %d", et, sourceTypes[et])
		}
	}
}
