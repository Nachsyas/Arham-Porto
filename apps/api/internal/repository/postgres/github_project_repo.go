package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/nachsyas/arham-porto/apps/api/internal/domain"
)

// GithubProjectRepository implements domain.GithubProjectRepository using PostgreSQL.
type GithubProjectRepository struct {
	client *Client
}

// NewGithubProjectRepository constructs a PostgreSQL-backed GithubProjectRepository.
func NewGithubProjectRepository(client *Client) *GithubProjectRepository {
	return &GithubProjectRepository{client: client}
}

// Create inserts a new GitHub repository record into github_projects.
func (r *GithubProjectRepository) Create(ctx context.Context, project *domain.GithubProject) error {
	pool := r.client.Pool()
	if pool == nil {
		return ErrDatabaseDisabled
	}

	topicsJSON, err := json.Marshal(project.Topics)
	if err != nil {
		return fmt.Errorf("failed to marshal topics: %w", err)
	}

	query := `
		INSERT INTO github_projects (
			github_id, name, full_name, description, html_url, homepage, language,
			stars, forks, topics, readme, preview_image, synced_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10::jsonb, $11, $12, NOW(), NOW(), NOW()
		)
		RETURNING id, created_at, updated_at, synced_at;
	`

	var (
		id        string
		createdAt time.Time
		updatedAt time.Time
		syncedAt  time.Time
	)

	err = pool.QueryRow(ctx, query,
		project.GithubID,
		project.Name,
		project.FullName,
		project.Description,
		project.HTMLURL,
		project.Homepage,
		project.Language,
		project.Stars,
		project.Forks,
		topicsJSON,
		project.Readme,
		project.PreviewImage,
	).Scan(&id, &createdAt, &updatedAt, &syncedAt)
	if err != nil {
		return fmt.Errorf("failed to insert github project %s: %w", project.FullName, err)
	}

	project.ID = id
	project.CreatedAt = createdAt
	project.UpdatedAt = updatedAt
	project.SyncedAt = syncedAt

	return nil
}

// Update updates an existing GitHub repository record by github_id.
func (r *GithubProjectRepository) Update(ctx context.Context, project *domain.GithubProject) error {
	pool := r.client.Pool()
	if pool == nil {
		return ErrDatabaseDisabled
	}

	topicsJSON, err := json.Marshal(project.Topics)
	if err != nil {
		return fmt.Errorf("failed to marshal topics: %w", err)
	}

	query := `
		UPDATE github_projects SET
			name = $2,
			full_name = $3,
			description = $4,
			html_url = $5,
			homepage = $6,
			language = $7,
			stars = $8,
			forks = $9,
			topics = $10::jsonb,
			synced_at = NOW(),
			updated_at = NOW()
		WHERE github_id = $1
		RETURNING id, updated_at, synced_at;
	`

	var (
		id        string
		updatedAt time.Time
		syncedAt  time.Time
	)

	err = pool.QueryRow(ctx, query,
		project.GithubID,
		project.Name,
		project.FullName,
		project.Description,
		project.HTMLURL,
		project.Homepage,
		project.Language,
		project.Stars,
		project.Forks,
		topicsJSON,
	).Scan(&id, &updatedAt, &syncedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("github project with id %d not found for update", project.GithubID)
		}
		return fmt.Errorf("failed to update github project %s: %w", project.FullName, err)
	}

	project.ID = id
	project.UpdatedAt = updatedAt
	project.SyncedAt = syncedAt

	return nil
}

// Upsert inserts a new GitHub repository or updates an existing repository matching github_id.
// Returns isCreated=true when a new row is inserted, or isCreated=false when updated.
func (r *GithubProjectRepository) Upsert(ctx context.Context, project *domain.GithubProject) (bool, error) {
	pool := r.client.Pool()
	if pool == nil {
		return false, ErrDatabaseDisabled
	}

	topicsJSON, err := json.Marshal(project.Topics)
	if err != nil {
		return false, fmt.Errorf("failed to marshal topics: %w", err)
	}

	query := `
		INSERT INTO github_projects (
			github_id, name, full_name, description, html_url, homepage, language,
			stars, forks, topics, readme, preview_image, synced_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10::jsonb, $11, $12, NOW(), NOW(), NOW()
		)
		ON CONFLICT (github_id) DO UPDATE SET
			name = EXCLUDED.name,
			full_name = EXCLUDED.full_name,
			description = EXCLUDED.description,
			html_url = EXCLUDED.html_url,
			homepage = EXCLUDED.homepage,
			language = EXCLUDED.language,
			stars = EXCLUDED.stars,
			forks = EXCLUDED.forks,
			topics = EXCLUDED.topics,
			synced_at = NOW(),
			updated_at = NOW()
		RETURNING id, (xmax = 0) AS is_created, created_at, updated_at, synced_at;
	`

	var (
		id        string
		isCreated bool
		createdAt time.Time
		updatedAt time.Time
		syncedAt  time.Time
	)

	err = pool.QueryRow(ctx, query,
		project.GithubID,
		project.Name,
		project.FullName,
		project.Description,
		project.HTMLURL,
		project.Homepage,
		project.Language,
		project.Stars,
		project.Forks,
		topicsJSON,
		project.Readme,
		project.PreviewImage,
	).Scan(&id, &isCreated, &createdAt, &updatedAt, &syncedAt)
	if err != nil {
		return false, fmt.Errorf("failed to upsert github project %s: %w", project.FullName, err)
	}

	project.ID = id
	project.CreatedAt = createdAt
	project.UpdatedAt = updatedAt
	project.SyncedAt = syncedAt

	return isCreated, nil
}

// List returns all synchronized GitHub repositories ordered by stars descending.
func (r *GithubProjectRepository) List(ctx context.Context) ([]*domain.GithubProject, error) {
	pool := r.client.Pool()
	if pool == nil {
		return nil, ErrDatabaseDisabled
	}

	query := `
		SELECT
			id, github_id, name, full_name, COALESCE(description, ''), COALESCE(html_url, ''),
			COALESCE(homepage, ''), COALESCE(language, ''), stars, forks,
			COALESCE(topics, '[]'::jsonb), readme, preview_image, synced_at, created_at, updated_at
		FROM github_projects
		ORDER BY stars DESC, updated_at DESC;
	`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query github projects: %w", err)
	}
	defer rows.Close()

	var projects []*domain.GithubProject

	for rows.Next() {
		var (
			p           domain.GithubProject
			topicsBytes []byte
		)

		err := rows.Scan(
			&p.ID,
			&p.GithubID,
			&p.Name,
			&p.FullName,
			&p.Description,
			&p.HTMLURL,
			&p.Homepage,
			&p.Language,
			&p.Stars,
			&p.Forks,
			&topicsBytes,
			&p.Readme,
			&p.PreviewImage,
			&p.SyncedAt,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan github project row: %w", err)
		}

		if len(topicsBytes) > 0 {
			if err := json.Unmarshal(topicsBytes, &p.Topics); err != nil {
				p.Topics = []string{}
			}
		} else {
			p.Topics = []string{}
		}

		projects = append(projects, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during github projects iteration: %w", err)
	}

	return projects, nil
}

// ListGithubProjects returns synchronized GitHub repositories ordered by newest updated first, with pagination support.
func (r *GithubProjectRepository) ListGithubProjects(ctx context.Context, page, limit int) ([]*domain.GithubProject, error) {
	pool := r.client.Pool()
	if pool == nil {
		return nil, ErrDatabaseDisabled
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	} else if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	query := `
		SELECT
			id, github_id, name, full_name, COALESCE(description, ''), COALESCE(html_url, ''),
			COALESCE(homepage, ''), COALESCE(language, ''), stars, forks,
			COALESCE(topics, '[]'::jsonb), readme, preview_image, synced_at, created_at, updated_at
		FROM github_projects
		ORDER BY updated_at DESC
		LIMIT $1 OFFSET $2;
	`

	rows, err := pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query paginated github projects: %w", err)
	}
	defer rows.Close()

	var projects []*domain.GithubProject

	for rows.Next() {
		var (
			p           domain.GithubProject
			topicsBytes []byte
		)

		err := rows.Scan(
			&p.ID,
			&p.GithubID,
			&p.Name,
			&p.FullName,
			&p.Description,
			&p.HTMLURL,
			&p.Homepage,
			&p.Language,
			&p.Stars,
			&p.Forks,
			&topicsBytes,
			&p.Readme,
			&p.PreviewImage,
			&p.SyncedAt,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan github project row: %w", err)
		}

		if len(topicsBytes) > 0 {
			if err := json.Unmarshal(topicsBytes, &p.Topics); err != nil {
				p.Topics = []string{}
			}
		} else {
			p.Topics = []string{}
		}

		projects = append(projects, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during paginated github projects iteration: %w", err)
	}

	return projects, nil
}

// GetByGithubID retrieves a single GitHub repository by its unique GitHub ID.
func (r *GithubProjectRepository) GetByGithubID(ctx context.Context, githubID int64) (*domain.GithubProject, error) {
	pool := r.client.Pool()
	if pool == nil {
		return nil, ErrDatabaseDisabled
	}

	query := `
		SELECT
			id, github_id, name, full_name, COALESCE(description, ''), COALESCE(html_url, ''),
			COALESCE(homepage, ''), COALESCE(language, ''), stars, forks,
			COALESCE(topics, '[]'::jsonb), readme, preview_image, synced_at, created_at, updated_at
		FROM github_projects
		WHERE github_id = $1;
	`

	var (
		p           domain.GithubProject
		topicsBytes []byte
	)

	err := pool.QueryRow(ctx, query, githubID).Scan(
		&p.ID,
		&p.GithubID,
		&p.Name,
		&p.FullName,
		&p.Description,
		&p.HTMLURL,
		&p.Homepage,
		&p.Language,
		&p.Stars,
		&p.Forks,
		&topicsBytes,
		&p.Readme,
		&p.PreviewImage,
		&p.SyncedAt,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get github project %d: %w", githubID, err)
	}

	if len(topicsBytes) > 0 {
		if err := json.Unmarshal(topicsBytes, &p.Topics); err != nil {
			p.Topics = []string{}
		}
	} else {
		p.Topics = []string{}
	}

	return &p, nil
}
