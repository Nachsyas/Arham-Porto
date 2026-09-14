package indexing

import (
	"strings"
	"testing"

	"github.com/nachsyas/arham-porto/apps/api/internal/domain"
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

func TestBuildCanonicalKnowledgeSources_ProfileApprovedFieldsOnly(t *testing.T) {
	dataDir := "../../../../data"

	sources, err := BuildCanonicalKnowledgeSources(dataDir)
	if err != nil {
		t.Fatalf("BuildCanonicalKnowledgeSources failed: %v", err)
	}

	var profileSource *domain.SourceWithChunks
	for i, s := range sources {
		if s.Source.SourceType == "canonical_profile" {
			profileSource = &sources[i]
			break
		}
	}

	if profileSource == nil {
		t.Fatal("expected canonical_profile source to be found")
	}

	for _, c := range profileSource.Chunks {
		// Regression: Must NOT place unapproved TODO-backed positioning sentence into chunk content (Corrections 14, 15, 16)
		unapprovedPositioning := "Building intelligent, scalable, and human-centered digital systems"
		if strings.Contains(c.Content, unapprovedPositioning) {
			t.Fatalf("profile chunk leaked unapproved positioning sentence: %s", c.Content)
		}
		if strings.Contains(c.Content, "Positioning:") {
			t.Fatalf("profile chunk leaked unapproved Positioning field: %s", c.Content)
		}

		// Must contain approved public baseline fields
		if !strings.Contains(c.Content, "Nachsyas Arham Mumtaz Nashohi") {
			t.Errorf("expected chunk to contain FullName, got: %s", c.Content)
		}
		if !strings.Contains(c.Content, "Software Engineer") {
			t.Errorf("expected chunk to contain Role, got: %s", c.Content)
		}
		if !strings.Contains(c.Content, "Arham Porto") {
			t.Errorf("expected chunk to contain ProjectName, got: %s", c.Content)
		}
		if !strings.Contains(c.Content, "Ask Arham AI") {
			t.Errorf("expected chunk to contain AIFeature, got: %s", c.Content)
		}
		if !strings.Contains(c.Content, "https://github.com/Nachsyas") {
			t.Errorf("expected chunk to contain GitHub URL, got: %s", c.Content)
		}
	}
}

