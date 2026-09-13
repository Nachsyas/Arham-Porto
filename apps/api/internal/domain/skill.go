package domain

import "context"

// Skill represents a verified technical area linked to evidence.
// Arbitrary percentage scores are strictly prohibited.
type Skill struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Category    string   `json:"category"`
	Claim       *string  `json:"claim,omitempty"`
	EvidenceIDs []string `json:"evidenceIds"`
	TODO        []string `json:"todo,omitempty"`
}

// SkillRepository defines storage access for skills.
type SkillRepository interface {
	ListSkills(ctx context.Context, category *string) ([]Skill, error)
}

