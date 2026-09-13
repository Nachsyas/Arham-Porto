package domain

import "context"

// Profile represents the professional identity of the portfolio owner.
// Standard library types only — zero external dependencies.
type Profile struct {
	FullName     string   `json:"fullName"`
	Role         string   `json:"role"`
	ProjectName  string   `json:"projectName"`
	AIFeature    string   `json:"aiFeature"`
	Positioning  *string  `json:"positioning,omitempty"`
	Bio          *string  `json:"bio,omitempty"`
	Email        *string  `json:"email,omitempty"`
	LinkedIn     *string  `json:"linkedin,omitempty"`
	GitHub       *string  `json:"github,omitempty"`
	CVURL        *string  `json:"cvUrl,omitempty"`
	CurrentCity  *string  `json:"currentCity,omitempty"`
	Availability *string  `json:"availability,omitempty"`
	TODO         []string `json:"todo,omitempty"`
}

// ProfileRepository defines the storage interface for profile data.
type ProfileRepository interface {
	GetProfile(ctx context.Context) (*Profile, error)
}

