package usecase

import (
	"context"

	"github.com/nachsyas/arham-porto/apps/api/internal/domain"
)

// SkillUseCase specifies operations for accessing skill data.
type SkillUseCase interface {
	ListSkills(ctx context.Context, category *string) ([]domain.Skill, error)
}

type skillUseCase struct {
	repo domain.SkillRepository
}

// NewSkillUseCase constructs a SkillUseCase.
func NewSkillUseCase(repo domain.SkillRepository) SkillUseCase {
	return &skillUseCase{repo: repo}
}

func (u *skillUseCase) ListSkills(ctx context.Context, category *string) ([]domain.Skill, error) {
	if u.repo == nil {
		return []domain.Skill{}, nil
	}
	return u.repo.ListSkills(ctx, category)
}
