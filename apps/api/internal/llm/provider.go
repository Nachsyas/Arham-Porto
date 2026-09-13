package llm

import (
	"context"
	"errors"
)

var (
	// ErrAIFeatureDisabled is returned when generation is attempted while AI_MODE=disabled.
	ErrAIFeatureDisabled = errors.New("ai feature is disabled")
	// ErrProviderUnavailable is returned when the remote LLM provider fails or is unreachable.
	ErrProviderUnavailable = errors.New("ai provider is temporarily unavailable")
	// ErrContextLengthExceeded is returned when the request/evidence context exceeds bounds.
	ErrContextLengthExceeded = errors.New("context length exceeded")
)

// EvidenceContext holds the server-constructed evidence unit passed into generation.
type EvidenceContext struct {
	ID          string   `json:"id"` // Server-issued ID: E1, E2, etc.
	Kind        string   `json:"kind"` // "github" or "portfolio"
	Title       string   `json:"title"`
	Repository  *string  `json:"repository,omitempty"`
	Path        *string  `json:"path,omitempty"`
	CommitSHA   *string  `json:"commit_sha,omitempty"`
	CitationURL *string  `json:"citation_url,omitempty"`
	Content     string   `json:"content"`
	ProjectID   *string  `json:"project_id,omitempty"`
	SkillIDs    []string `json:"skill_ids,omitempty"`
	EvidenceID  *string  `json:"evidence_id,omitempty"`
	Similarity  float32  `json:"-"` // Internal only — NEVER serialized to public clients
}

// GroundedSegment pairs a specific claim sentence/clause with its supporting evidence IDs.
type GroundedSegment struct {
	Text        string   `json:"text"`
	EvidenceIDs []string `json:"evidence_ids"`
}

// GeneratedAnswer represents the structured generation output from the LLM provider.
type GeneratedAnswer struct {
	Status             string            `json:"status"` // supported, insufficient_evidence, privacy_refusal, scope_refusal
	Segments           []GroundedSegment `json:"segments"`
	SuggestedActionIDs []string          `json:"suggested_action_ids"`
}

// GenerateRequest defines the input payload for the LLM provider.
type GenerateRequest struct {
	SystemInstruction string            `json:"system_instruction"`
	UserQuestion      string            `json:"user_question"`
	Evidence          []EvidenceContext `json:"evidence"`
}

// SourceCitation represents a server-certified, immutable source of truth.
type SourceCitation struct {
	ID         string  `json:"id"`                   // S1, S2, etc.
	Kind       string  `json:"kind"`                 // "github" or "portfolio"
	Label      string  `json:"label"`                // Public friendly label (e.g. "EduTrace / README.md" or "Project: EduTrace")
	URL        *string `json:"url,omitempty"`        // Verified URL (nil if canonical portfolio stop without web link)
	Repository *string `json:"repository,omitempty"` // For GitHub citations
	Path       *string `json:"path,omitempty"`       // For GitHub citations
	CommitSHA  *string `json:"commit_sha,omitempty"` // For GitHub citations
}

// PublicEvidenceItem represents an evidence chunk presented to the reviewer.
// Similarity scores and raw internal paths are strictly omitted.
type PublicEvidenceItem struct {
	ID         string  `json:"id"`                   // E1, E2, etc.
	Kind       string  `json:"kind"`                 // "github" or "portfolio"
	Title      string  `json:"title"`                // Document / section title
	Repository *string `json:"repository,omitempty"` // If applicable
	Path       *string `json:"path,omitempty"`       // If applicable
	Excerpt    string  `json:"excerpt"`              // Bounded plain text snippet (max 800 chars)
	CitationID string  `json:"citation_id"`          // Links to SourceCitation.ID
}

// GroundedResponse represents the final validated, safe public payload returned to the client.
type GroundedResponse struct {
	Status   string               `json:"status"` // supported, insufficient_evidence, privacy_refusal, scope_refusal
	Answer   string               `json:"answer"`
	Segments []GroundedSegment    `json:"segments"`
	Evidence []PublicEvidenceItem `json:"evidence"`
	Sources  []SourceCitation     `json:"sources"`
	Actions  []PublicAction       `json:"actions"`
}

// PublicAction is a safe action sent in the public response.
type PublicAction struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

// LLMProvider defines the provider-neutral interface for structured evidence generation.
type LLMProvider interface {
	Generate(ctx context.Context, req GenerateRequest) (GeneratedAnswer, error)
	ProviderName() string
	Model() string
}
