package dto

import "github.com/nachsyas/arham-porto/apps/api/internal/domain"

// SkillResponse represents the public skill DTO.
// Arbitrary percentage scores are strictly prohibited. Internal TODOs are excluded.
type SkillResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Category    string   `json:"category"`
	Claim       *string  `json:"claim,omitempty"`
	EvidenceIDs []string `json:"evidence_ids"`
}

// FromDomainSkill maps a domain Skill entity to a safe public SkillResponse DTO.
func FromDomainSkill(s domain.Skill) SkillResponse {
	evidenceIDs := s.EvidenceIDs
	if evidenceIDs == nil {
		evidenceIDs = []string{}
	}
	return SkillResponse{
		ID:          s.ID,
		Name:        s.Name,
		Category:    s.Category,
		Claim:       s.Claim,
		EvidenceIDs: evidenceIDs,
	}
}

// FromDomainSkills maps a slice of domain Skills to public SkillResponse DTOs.
func FromDomainSkills(skills []domain.Skill) []SkillResponse {
	if skills == nil {
		return []SkillResponse{}
	}
	result := make([]SkillResponse, len(skills))
	for i, s := range skills {
		result[i] = FromDomainSkill(s)
	}
	return result
}
