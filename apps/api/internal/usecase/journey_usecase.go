package usecase

import (
	"context"

	"github.com/nachsyas/arham-porto/apps/api/internal/domain"
)

// JourneyUseCase specifies operations for accessing journey data.
type JourneyUseCase interface {
	ListPublicStops(ctx context.Context) ([]domain.JourneyStop, error)
}

type journeyUseCase struct {
	repo domain.JourneyRepository
}

// NewJourneyUseCase constructs a JourneyUseCase.
func NewJourneyUseCase(repo domain.JourneyRepository) JourneyUseCase {
	return &journeyUseCase{repo: repo}
}

func (u *journeyUseCase) ListPublicStops(ctx context.Context) ([]domain.JourneyStop, error) {
	if u.repo == nil {
		return []domain.JourneyStop{}, nil
	}
	// Strictly request public-only stops
	return u.repo.ListStops(ctx, true)
}
