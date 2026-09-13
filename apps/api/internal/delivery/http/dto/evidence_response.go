package dto

import "github.com/nachsyas/arham-porto/apps/api/internal/domain"

// EvidenceResponse represents the public evidence DTO.
// Internal metadata (TODOs) are excluded.
type EvidenceResponse struct {
	ID         string   `json:"id"`
	Type       string   `json:"type"`
	Title      string   `json:"title"`
	SkillIDs   []string `json:"skill_ids"`
	SourceURL  *string  `json:"source_url,omitempty"`
	SourcePath *string  `json:"source_path,omitempty"`
	Summary    string   `json:"summary"`
	Verified   bool     `json:"verified"`
}

// FromDomainEvidence maps a domain Evidence entity to a safe public EvidenceResponse DTO.
func FromDomainEvidence(e domain.Evidence) EvidenceResponse {
	skillIDs := e.SkillIDs
	if skillIDs == nil {
		skillIDs = []string{}
	}
	return EvidenceResponse{
		ID:         e.ID,
		Type:       e.Type,
		Title:      e.Title,
		SkillIDs:   skillIDs,
		SourceURL:  e.SourceURL,
		SourcePath: e.SourcePath,
		Summary:    e.Summary,
		Verified:   e.Verified,
	}
}

// FromDomainEvidenceList maps a slice of domain Evidence entities to public EvidenceResponse DTOs.
func FromDomainEvidenceList(items []domain.Evidence) []EvidenceResponse {
	if items == nil {
		return []EvidenceResponse{}
	}
	result := make([]EvidenceResponse, len(items))
	for i, item := range items {
		result[i] = FromDomainEvidence(item)
	}
	return result
}
