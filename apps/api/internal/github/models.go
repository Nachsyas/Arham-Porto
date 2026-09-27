package github

// Repository represents normalized repository metadata fetched from GitHub API.
type Repository struct {
	GithubID    int64    `json:"github_id"`
	Name        string   `json:"name"`
	FullName    string   `json:"full_name"`
	Description string   `json:"description"`
	HTMLURL     string   `json:"html_url"`
	Homepage    string   `json:"homepage"`
	Language    string   `json:"language"`
	Stars       int      `json:"stars"`
	Forks       int      `json:"forks"`
	Topics      []string `json:"topics"`
}

// githubRepoResponse maps the raw GitHub REST API v3 JSON repository schema.
type githubRepoResponse struct {
	ID              int64    `json:"id"`
	Name            string   `json:"name"`
	FullName        string   `json:"full_name"`
	Description     *string  `json:"description"`
	HTMLURL         string   `json:"html_url"`
	Homepage        *string  `json:"homepage"`
	Language        *string  `json:"language"`
	StargazersCount int      `json:"stargazers_count"`
	ForksCount      int      `json:"forks_count"`
	Topics          []string `json:"topics"`
}

// ToRepository normalizes nullable and optional fields into a clean Repository model.
func (r githubRepoResponse) ToRepository() Repository {
	desc := ""
	if r.Description != nil {
		desc = *r.Description
	}

	hp := ""
	if r.Homepage != nil {
		hp = *r.Homepage
	}

	lang := ""
	if r.Language != nil {
		lang = *r.Language
	}

	topics := r.Topics
	if topics == nil {
		topics = []string{}
	}

	return Repository{
		GithubID:    r.ID,
		Name:        r.Name,
		FullName:    r.FullName,
		Description: desc,
		HTMLURL:     r.HTMLURL,
		Homepage:    hp,
		Language:    lang,
		Stars:       r.StargazersCount,
		Forks:       r.ForksCount,
		Topics:      topics,
	}
}

// SyncResult summarizes the outcome of repository synchronization.
type SyncResult struct {
	Synced  int `json:"synced"`
	Created int `json:"created"`
	Updated int `json:"updated"`
}
