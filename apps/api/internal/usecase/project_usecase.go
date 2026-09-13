package usecase

import (
	"context"

	"github.com/nachsyas/arham-porto/apps/api/internal/domain"
)

// ProjectUseCase specifies operations for accessing project data.
type ProjectUseCase interface {
	ListProjects(ctx context.Context) ([]domain.Project, error)
	GetProjectBySlug(ctx context.Context, slug string) (*domain.Project, error)
}

type projectUseCase struct {
	repo domain.ProjectRepository
}

// NewProjectUseCase constructs a ProjectUseCase.
func NewProjectUseCase(repo domain.ProjectRepository) ProjectUseCase {
	return &projectUseCase{repo: repo}
}

func (u *projectUseCase) ListProjects(ctx context.Context) ([]domain.Project, error) {
	if u.repo == nil {
		return []domain.Project{}, nil
	}
	return u.repo.ListProjects(ctx)
}

func (u *projectUseCase) GetProjectBySlug(ctx context.Context, slug string) (*domain.Project, error) {
	if u.repo == nil {
		return nil, nil
	}
	return u.repo.GetProjectBySlug(ctx, slug)
}
