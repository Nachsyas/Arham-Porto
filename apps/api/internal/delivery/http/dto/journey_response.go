package dto

import "github.com/nachsyas/arham-porto/apps/api/internal/domain"

// JourneyStopResponse represents the safe public journey stop DTO.
// Raw coordinates, internal TODOs, private metadata, and publication-control flags are strictly excluded.
type JourneyStopResponse struct {
	ID          string  `json:"id"`
	Category    string  `json:"category"`
	Title       *string `json:"title,omitempty"`
	Institution *string `json:"institution,omitempty"`
	City        *string `json:"city,omitempty"`
	Region      *string `json:"region,omitempty"`
	Country     string  `json:"country"`
	Period      *string `json:"period,omitempty"`
	Description *string `json:"description,omitempty"`
	Image       *string `json:"image,omitempty"`
}

// FromDomainJourneyStop maps a domain JourneyStop entity to a safe public JourneyStopResponse DTO.
func FromDomainJourneyStop(s domain.JourneyStop) JourneyStopResponse {
	return JourneyStopResponse{
		ID:          s.ID,
		Category:    s.Category,
		Title:       SanitizeStringPtr(s.Title),
		Institution: SanitizeStringPtr(s.Institution),
		City:        SanitizeStringPtr(s.City),
		Region:      SanitizeStringPtr(s.Region),
		Country:     s.Country,
		Period:      SanitizeStringPtr(s.Period),
		Description: SanitizeStringPtr(s.Description),
		Image:       SanitizeImageURL(s.Image),
	}
}

// FromDomainJourneyStops maps a slice of domain JourneyStops to public JourneyStopResponse DTOs.
func FromDomainJourneyStops(stops []domain.JourneyStop) []JourneyStopResponse {
	if stops == nil {
		return []JourneyStopResponse{}
	}
	result := make([]JourneyStopResponse, len(stops))
	for i, s := range stops {
		result[i] = FromDomainJourneyStop(s)
	}
	return result
}
