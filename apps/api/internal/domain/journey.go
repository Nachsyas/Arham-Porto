package domain

import "context"

// JourneyStop represents an educational or geographical checkpoint.
type JourneyStop struct {
	ID          string    `json:"id"`
	Category    string    `json:"category"`
	Title       *string   `json:"title,omitempty"`
	Institution *string   `json:"institution,omitempty"`
	City        *string   `json:"city,omitempty"`
	Region      *string   `json:"region,omitempty"`
	Country     string    `json:"country"`
	Coordinates []float64 `json:"coordinates,omitempty"`
	Period      *string   `json:"period,omitempty"`
	Description *string   `json:"description,omitempty"`
	Image       *string   `json:"image,omitempty"`
	Public      bool      `json:"public"`
	TODO        []string  `json:"todo,omitempty"`
}

// JourneyRepository defines storage access for journey stops.
type JourneyRepository interface {
	ListStops(ctx context.Context) ([]JourneyStop, error)
}
