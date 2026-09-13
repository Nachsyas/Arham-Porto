package dto

import "github.com/nachsyas/arham-porto/apps/api/internal/domain"

// ProjectResponse represents the public project DTO.
// Internal metadata (TODOs) are strictly excluded.
type ProjectResponse struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Slug          string   `json:"slug"`
	Summary       *string  `json:"summary,omitempty"`
	Problem       *string  `json:"problem,omitempty"`
	Solution      *string  `json:"solution,omitempty"`
	Role          []string `json:"role"`
	Contributions []string `json:"contributions"`
	Technologies  []string `json:"technologies"`
	GitHubURL     *string  `json:"github_url,omitempty"`
	DemoURL       *string  `json:"demo_url,omitempty"`
	Image         *string  `json:"image,omitempty"`
	Featured      bool     `json:"featured"`
	Category      *string  `json:"category,omitempty"`
	EvidenceIDs   []string `json:"evidence_ids"`
}

// FromDomainProject maps a domain Project entity to a safe public ProjectResponse DTO.
func FromDomainProject(p domain.Project) ProjectResponse {
	role := p.Role
	if role == nil {
		role = []string{}
	}
	contributions := p.Contributions
	if contributions == nil {
		contributions = []string{}
	}
	technologies := p.Technologies
	if technologies == nil {
		technologies = []string{}
	}
	evidenceIDs := p.EvidenceIDs
	if evidenceIDs == nil {
		evidenceIDs = []string{}
	}

	return ProjectResponse{
		ID:            p.ID,
		Title:         p.Title,
		Slug:          p.Slug,
		Summary:       SanitizeStringPtr(p.Summary),
		Problem:       SanitizeStringPtr(p.Problem),
		Solution:      SanitizeStringPtr(p.Solution),
		Role:          role,
		Contributions: contributions,
		Technologies:  technologies,
		GitHubURL:     SanitizeExternalURL(p.GitHubURL),
		DemoURL:       SanitizeExternalURL(p.DemoURL),
		Image:         SanitizeImageURL(p.Image),
		Featured:      p.Featured,
		Category:      SanitizeStringPtr(p.Category),
		EvidenceIDs:   evidenceIDs,
	}
}

// FromDomainProjects maps a slice of domain Projects to public ProjectResponse DTOs.
func FromDomainProjects(projects []domain.Project) []ProjectResponse {
	if projects == nil {
		return []ProjectResponse{}
	}
	result := make([]ProjectResponse, len(projects))
	for i, p := range projects {
		result[i] = FromDomainProject(p)
	}
	return result
}
