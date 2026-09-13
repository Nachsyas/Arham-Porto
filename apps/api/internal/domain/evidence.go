package domain

import "context"

// Evidence represents a verifiable technical artifact (repository, commit, architecture doc).
type Evidence struct {
	ID         string   `json:"id"`
	Type       string   `json:"type"` // project, github, cv, experience, education
	Title      string   `json:"title"`
	SkillIDs   []string `json:"skillIds"`
	SourceURL  *string  `json:"sourceUrl,omitempty"`
	SourcePath *string  `json:"sourcePath,omitempty"`
	Summary    string   `json:"summary"`
	Verified   bool     `json:"verified"`
	TODO       []string `json:"todo,omitempty"`
}

// EvidenceRepository defines storage access for evidence.
type EvidenceRepository interface {
	ListEvidence(ctx context.Context, skillID *string, projectID *string) ([]Evidence, error)
	GetEvidenceByID(ctx context.Context, id string) (*Evidence, error)
}

