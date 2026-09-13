package domain

import "context"

// Evidence represents a verifiable technical artifact (repository, commit, architecture doc).
type Evidence struct {
	ID         string   `json:"id"`
	Type       string   `json:"type"` // project, github, cv, experience, education
	Title      string   `json:"title"`
	SkillIDs   []string `json:"skill_ids"`
	SourceURL  *string  `json:"source_url,omitempty"`
	SourcePath *string  `json:"source_path,omitempty"`
	Summary    string   `json:"summary"`
	Verified   bool     `json:"verified"`
	TODO       []string `json:"todo,omitempty"`
}

// EvidenceRepository defines storage access for evidence.
type EvidenceRepository interface {
	ListEvidence(ctx context.Context) ([]Evidence, error)
	GetEvidenceByID(ctx context.Context, id string) (*Evidence, error)
}
