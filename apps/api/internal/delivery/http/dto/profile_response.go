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
// Any field backed by unapproved TODO markers or containing invalid schemes is stripped to nil.
func FromDomainProfile(p *domain.Profile) ProfileResponse {
	if p == nil {
		return ProfileResponse{}
	}

	var positioning *string
	if !IsTodoBacked(p.TODO, "POSITIONING") {
		positioning = SanitizeStringPtr(p.Positioning)
	}

	var bio *string
	if !IsTodoBacked(p.TODO, "BIO") {
		bio = SanitizeStringPtr(p.Bio)
	}

	var currentCity *string
	if !IsTodoBacked(p.TODO, "CURRENT_CITY") {
		currentCity = SanitizeStringPtr(p.CurrentCity)
	}

	var availability *string
	if !IsTodoBacked(p.TODO, "AVAILABILITY") {
		availability = SanitizeStringPtr(p.Availability)
	}

	return ProfileResponse{
		FullName:     p.FullName,
		Role:         p.Role,
		ProjectName:  p.ProjectName,
		AIFeature:    p.AIFeature,
		Positioning:  positioning,
		Bio:          bio,
		GitHub:       SanitizeExternalURL(p.GitHub),
		CurrentCity:  currentCity,
		Availability: availability,
	}
}
