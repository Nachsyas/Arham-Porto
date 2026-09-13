package usecase

import (
	"context"

	"github.com/nachsyas/arham-porto/apps/api/internal/domain"
)

// EvidenceUseCase specifies operations for accessing evidence data.
type EvidenceUseCase interface {
	ListEvidence(ctx context.Context) ([]domain.Evidence, error)
	GetEvidenceByID(ctx context.Context, id string) (*domain.Evidence, error)
}

type evidenceUseCase struct {
	repo domain.EvidenceRepository
}

// NewEvidenceUseCase constructs an EvidenceUseCase.
func NewEvidenceUseCase(repo domain.EvidenceRepository) EvidenceUseCase {
	return &evidenceUseCase{repo: repo}
}

func (u *evidenceUseCase) ListEvidence(ctx context.Context) ([]domain.Evidence, error) {
	if u.repo == nil {
		return []domain.Evidence{}, nil
	}
	return u.repo.ListEvidence(ctx)
}

func (u *evidenceUseCase) GetEvidenceByID(ctx context.Context, id string) (*domain.Evidence, error) {
	if u.repo == nil {
		return nil, nil
	}
	return u.repo.GetEvidenceByID(ctx, id)
}
