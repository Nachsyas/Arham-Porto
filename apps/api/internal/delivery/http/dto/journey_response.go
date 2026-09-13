package dto

import "github.com/nachsyas/arham-porto/apps/api/internal/domain"

// JourneyStopResponse represents the safe public journey stop DTO.
// Raw coordinates, internal TODOs, and private metadata are strictly excluded.
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
	Public      bool    `json:"public"`
}

// FromDomainJourneyStop maps a domain JourneyStop entity to a safe public JourneyStopResponse DTO.
func FromDomainJourneyStop(s domain.JourneyStop) JourneyStopResponse {
	return JourneyStopResponse{
		ID:          s.ID,
		Category:    s.Category,
		Title:       s.Title,
		Institution: s.Institution,
		City:        s.City,
		Region:      s.Region,
		Country:     s.Country,
		Period:      s.Period,
		Description: s.Description,
		Image:       s.Image,
		Public:      s.Public,
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
