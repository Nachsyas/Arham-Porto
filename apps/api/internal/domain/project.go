package domain

import "context"

// Project represents an engineering case study or showcased system.
type Project struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Slug          string   `json:"slug"`
	Summary       *string  `json:"summary,omitempty"`
	Problem       *string  `json:"problem,omitempty"`
	Solution      *string  `json:"solution,omitempty"`
	Role          []string `json:"role"`
	Contributions []string `json:"contributions"`
	Technologies  []string `json:"technologies"`
	GitHubURL     *string  `json:"githubUrl,omitempty"`
	DemoURL       *string  `json:"demoUrl,omitempty"`
	Image         *string  `json:"image,omitempty"`
	Featured      bool     `json:"featured"`
	Category      *string  `json:"category,omitempty"`
	EvidenceIDs   []string `json:"evidenceIds"`
	TODO          []string `json:"todo,omitempty"`
}

// ProjectRepository defines storage access for projects.
type ProjectRepository interface {
	ListProjects(ctx context.Context, category *string, featured *bool) ([]Project, error)
	GetProjectBySlug(ctx context.Context, slug string) (*Project, error)
}

