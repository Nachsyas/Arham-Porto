package github

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/nachsyas/arham-porto/apps/api/internal/domain"
)

// DefaultSyncUsername is the portfolio owner whose repositories are synchronized.
const DefaultSyncUsername = "Nachsyas"

// SyncService orchestrates fetching repositories from GitHub and synchronizing them into the database.
type SyncService struct {
	client   *Client
	repo     domain.GithubProjectRepository
	username string
}

// NewSyncService constructs a new SyncService.
func NewSyncService(client *Client, repo domain.GithubProjectRepository, username ...string) *SyncService {
	user := DefaultSyncUsername
	if len(username) > 0 && strings.TrimSpace(username[0]) != "" {
		user = strings.TrimSpace(username[0])
	}
	return &SyncService{
		client:   client,
		repo:     repo,
		username: user,
	}
}

// SyncRepositories fetches all owner repositories for the configured user,
// normalizes the metadata, and upserts them into PostgreSQL.
// Historical records are strictly preserved (deleted repositories on GitHub are not removed).
func (s *SyncService) SyncRepositories(ctx context.Context) (*SyncResult, error) {
	if s.client == nil {
		return nil, ErrMissingToken
	}
	if s.repo == nil {
		return nil, errors.New("repository storage layer is not configured")
	}

	repos, err := s.client.FetchUserRepositories(ctx, s.username)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user repositories: %w", err)
	}

	result := &SyncResult{}

	for _, r := range repos {
		project := &domain.GithubProject{
			GithubID:    r.GithubID,
			Name:        strings.TrimSpace(r.Name),
			FullName:    strings.TrimSpace(r.FullName),
			Description: strings.TrimSpace(r.Description),
			HTMLURL:     strings.TrimSpace(r.HTMLURL),
			Homepage:    strings.TrimSpace(r.Homepage),
			Language:    strings.TrimSpace(r.Language),
			Stars:       r.Stars,
			Forks:       r.Forks,
			Topics:      r.Topics,
		}
		if project.Topics == nil {
			project.Topics = []string{}
		}

		isCreated, err := s.repo.Upsert(ctx, project)
		if err != nil {
			return nil, fmt.Errorf("failed to upsert repository %s (%d): %w", r.FullName, r.GithubID, err)
		}

		result.Synced++
		if isCreated {
			result.Created++
		} else {
			result.Updated++
		}
	}

	return result, nil
}
