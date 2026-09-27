package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/nachsyas/arham-porto/apps/api/internal/config"
	"github.com/nachsyas/arham-porto/apps/api/internal/domain"
)

func TestGithubProjectRepository_DatabaseDisabled(t *testing.T) {
	client := NewClient("", config.DatabaseModeDisabled)
	repo := NewGithubProjectRepository(client)
	ctx := context.Background()

	project := &domain.GithubProject{
		GithubID: 100,
		Name:     "TestRepo",
		FullName: "Nachsyas/TestRepo",
	}

	t.Run("Create fails with ErrDatabaseDisabled", func(t *testing.T) {
		err := repo.Create(ctx, project)
		if !errors.Is(err, ErrDatabaseDisabled) {
			t.Errorf("expected ErrDatabaseDisabled, got %v", err)
		}
	})

	t.Run("Update fails with ErrDatabaseDisabled", func(t *testing.T) {
		err := repo.Update(ctx, project)
		if !errors.Is(err, ErrDatabaseDisabled) {
			t.Errorf("expected ErrDatabaseDisabled, got %v", err)
		}
	})

	t.Run("Upsert fails with ErrDatabaseDisabled", func(t *testing.T) {
		_, err := repo.Upsert(ctx, project)
		if !errors.Is(err, ErrDatabaseDisabled) {
			t.Errorf("expected ErrDatabaseDisabled, got %v", err)
		}
	})

	t.Run("List fails with ErrDatabaseDisabled", func(t *testing.T) {
		_, err := repo.List(ctx)
		if !errors.Is(err, ErrDatabaseDisabled) {
			t.Errorf("expected ErrDatabaseDisabled, got %v", err)
		}
	})

	t.Run("ListGithubProjects fails with ErrDatabaseDisabled", func(t *testing.T) {
		_, err := repo.ListGithubProjects(ctx, 1, 20)
		if !errors.Is(err, ErrDatabaseDisabled) {
			t.Errorf("expected ErrDatabaseDisabled, got %v", err)
		}
	})

	t.Run("GetByGithubID fails with ErrDatabaseDisabled", func(t *testing.T) {
		_, err := repo.GetByGithubID(ctx, 100)
		if !errors.Is(err, ErrDatabaseDisabled) {
			t.Errorf("expected ErrDatabaseDisabled, got %v", err)
		}
	})
}
