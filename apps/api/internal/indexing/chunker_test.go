package indexing

import (
	"strings"
	"testing"

	"github.com/nachsyas/arham-porto/apps/api/internal/domain"
)

func TestChunkDocumentMarkdown(t *testing.T) {
	source := domain.KnowledgeSource{
		ID:         "src-123",
		Repository: "Nachsyas/EduTrace",
		Path:       "README.md",
		SourceType: "doc",
	}

	markdown := `# EduTrace

EduTrace is a verifiable credential system on Ethereum.

## Architecture

It uses Soulbound Tokens (ERC-5173) for non-transferable academic badges.

## Smart Contracts

The contract is implemented in Solidity using Foundry.`

	projectID := "edutrace"
	skillIDs := []string{"solidity", "ethereum"}
	evidenceID := "evidence-edutrace"

	chunks := ChunkDocument(source, markdown, &projectID, skillIDs, &evidenceID)
	if len(chunks) == 0 {
		t.Fatalf("expected chunks, got 0")
	}

	for _, c := range chunks {
		if c.SourceID != "src-123" {
			t.Errorf("expected source ID src-123, got %s", c.SourceID)
		}
		if *c.ProjectID != "edutrace" {
			t.Errorf("expected project ID edutrace, got %v", c.ProjectID)
		}
		if len(c.SkillIDs) != 2 {
			t.Errorf("expected 2 skill IDs, got %d", len(c.SkillIDs))
		}
		if *c.EvidenceID != "evidence-edutrace" {
			t.Errorf("expected evidence ID evidence-edutrace, got %v", c.EvidenceID)
		}
		if !strings.Contains(c.Content, "Repository: Nachsyas/EduTrace") {
			t.Errorf("expected content to contain repository header, got:\n%s", c.Content)
		}
		if !strings.Contains(c.Content, "Path: README.md") {
			t.Errorf("expected content to contain path header, got:\n%s", c.Content)
		}
		if c.Checksum == "" {
			t.Errorf("expected non-empty checksum")
		}
		if c.ID == "" {
			t.Errorf("expected non-empty chunk ID")
		}
	}
}

func TestChunkDocumentCode(t *testing.T) {
	source := domain.KnowledgeSource{
		ID:         "src-456",
		Repository: "Nachsyas/smart-kitchen-backend",
		Path:       "main.go",
		SourceType: "entrypoint",
	}

	codeLines := make([]string, 60)
	for i := range codeLines {
		codeLines[i] = "fmt.Println(\"code line\")"
	}
	code := strings.Join(codeLines, "\n")

	chunks := ChunkDocument(source, code, nil, nil, nil)
	if len(chunks) < 2 {
		t.Fatalf("expected at least 2 chunks for 60 lines, got %d", len(chunks))
	}

	if !strings.Contains(chunks[0].Content, "Section: Lines 1-40") {
		t.Errorf("expected section Lines 1-40 in first chunk, got:\n%s", chunks[0].Content)
	}
}
