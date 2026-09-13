package usecase

import (
	"context"

	"github.com/nachsyas/arham-porto/apps/api/internal/domain"
)

// ProfileUseCase specifies operations for accessing profile data.
type ProfileUseCase interface {
	GetProfile(ctx context.Context) (*domain.Profile, error)
}

type profileUseCase struct {
	repo domain.ProfileRepository
}

// NewProfileUseCase constructs a ProfileUseCase.
func NewProfileUseCase(repo domain.ProfileRepository) ProfileUseCase {
	return &profileUseCase{repo: repo}
}

func (u *profileUseCase) GetProfile(ctx context.Context) (*domain.Profile, error) {
	if u.repo == nil {
		return &domain.Profile{
			FullName:    "Nachsyas Arham Mumtaz Nashohi",
			Role:        "Software Engineer",
			ProjectName: "Arham Porto",
			AIFeature:   "Ask Arham AI",
		}, nil
	}
	return u.repo.GetProfile(ctx)
}

