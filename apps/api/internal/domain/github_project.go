package domain

import (
	"context"
	"time"
)

// GithubProject represents a synchronized GitHub repository record in the database.
type GithubProject struct {
	ID           string    `json:"id"`
	GithubID     int64     `json:"github_id"`
	Name         string    `json:"name"`
	FullName     string    `json:"full_name"`
	Description  string    `json:"description"`
	HTMLURL      string    `json:"html_url"`
	Homepage     string    `json:"homepage"`
	Language     string    `json:"language"`
	Stars        int       `json:"stars"`
	Forks        int       `json:"forks"`
	Topics       []string  `json:"topics"`
	Readme       *string   `json:"readme,omitempty"`
	PreviewImage *string   `json:"preview_image,omitempty"`
	SyncedAt     time.Time `json:"synced_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// GithubProjectRepository defines persistence operations for GitHub projects.
type GithubProjectRepository interface {
	Create(ctx context.Context, project *GithubProject) error
	Update(ctx context.Context, project *GithubProject) error
	Upsert(ctx context.Context, project *GithubProject) (isCreated bool, err error)
	List(ctx context.Context) ([]*GithubProject, error)
	ListGithubProjects(ctx context.Context, page, limit int) ([]*GithubProject, error)
	GetByGithubID(ctx context.Context, githubID int64) (*GithubProject, error)
}

