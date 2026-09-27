package dto

import (
	"time"

	"github.com/nachsyas/arham-porto/apps/api/internal/domain"
)

// GithubProjectResponse represents the serialized public GitHub repository DTO.
type GithubProjectResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	FullName    string   `json:"full_name"`
	Description string   `json:"description"`
	HTMLURL     string   `json:"html_url"`
	Homepage    string   `json:"homepage"`
	Language    string   `json:"language"`
	Stars       int      `json:"stars"`
	Forks       int      `json:"forks"`
	Topics      []string `json:"topics"`
	SyncedAt    string   `json:"synced_at"`
}

// GithubProjectsListResponse envelopes a list of GitHub project DTOs.
type GithubProjectsListResponse struct {
	Status string                  `json:"status"`
	Data   []GithubProjectResponse `json:"data"`
}

// FromDomainGithubProject converts an internal domain GithubProject into a sanitized public DTO.
func FromDomainGithubProject(p domain.GithubProject) GithubProjectResponse {
	topics := p.Topics
	if topics == nil {
		topics = []string{}
	}

	return GithubProjectResponse{
		ID:          p.ID,
		Name:        p.Name,
		FullName:    p.FullName,
		Description: p.Description,
		HTMLURL:     p.HTMLURL,
		Homepage:    p.Homepage,
		Language:    p.Language,
		Stars:       p.Stars,
		Forks:       p.Forks,
		Topics:      topics,
		SyncedAt:    p.SyncedAt.UTC().Format(time.RFC3339),
	}
}

// FromDomainGithubProjects converts a slice of domain projects into DTO responses.
func FromDomainGithubProjects(projects []*domain.GithubProject) []GithubProjectResponse {
	if len(projects) == 0 {
		return []GithubProjectResponse{}
	}

	res := make([]GithubProjectResponse, len(projects))
	for i, p := range projects {
		res[i] = FromDomainGithubProject(*p)
	}
	return res
}
