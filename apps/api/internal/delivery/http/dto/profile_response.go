package dto

import "github.com/nachsyas/arham-porto/apps/api/internal/domain"

// ProfileResponse represents the public profile DTO.
// Internal metadata (TODOs, private fields) are strictly stripped.
type ProfileResponse struct {
	FullName     string  `json:"full_name"`
	Role         string  `json:"role"`
	ProjectName  string  `json:"project_name"`
	AIFeature    string  `json:"ai_feature"`
	Positioning  *string `json:"positioning,omitempty"`
	Bio          *string `json:"bio,omitempty"`
	GitHub       *string `json:"github,omitempty"`
	CurrentCity  *string `json:"current_city,omitempty"`
	Availability *string `json:"availability,omitempty"`
}

// FromDomainProfile maps a domain Profile to a safe public ProfileResponse.
func FromDomainProfile(p *domain.Profile) ProfileResponse {
	if p == nil {
		return ProfileResponse{}
	}
	return ProfileResponse{
		FullName:     p.FullName,
		Role:         p.Role,
		ProjectName:  p.ProjectName,
		AIFeature:    p.AIFeature,
		Positioning:  p.Positioning,
		Bio:          p.Bio,
		GitHub:       p.GitHub,
		CurrentCity:  p.CurrentCity,
		Availability: p.Availability,
	}
}
