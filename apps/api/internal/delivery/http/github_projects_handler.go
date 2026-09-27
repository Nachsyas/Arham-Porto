package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/nachsyas/arham-porto/apps/api/internal/delivery/http/dto"
	"github.com/nachsyas/arham-porto/apps/api/internal/repository/postgres"
)

// ListGithubProjects handles GET /api/v1/github/projects.
// Returns paginated synchronized GitHub projects ordered by newest updated first.
func (h *Handler) ListGithubProjects(w http.ResponseWriter, r *http.Request) {
	if h.githubProjectRepo == nil {
		writeError(w, http.StatusServiceUnavailable, "service_unavailable", "github project repository is not configured")
		return
	}

	query := r.URL.Query()

	// Strict query parameter validation: reject duplicate keys
	if len(query["page"]) > 1 {
		writeError(w, http.StatusBadRequest, "bad_request", "ambiguous repeated 'page' query parameter")
		return
	}
	if len(query["limit"]) > 1 {
		writeError(w, http.StatusBadRequest, "bad_request", "ambiguous repeated 'limit' query parameter")
		return
	}

	page := 1
	if pageStr := query.Get("page"); pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil || p < 1 {
			writeError(w, http.StatusBadRequest, "bad_request", "invalid 'page' query parameter, must be a positive integer")
			return
		}
		page = p
	}

	limit := 20
	if limitStr := query.Get("limit"); limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil || l < 1 || l > 100 {
			writeError(w, http.StatusBadRequest, "bad_request", "invalid 'limit' query parameter, must be between 1 and 100")
			return
		}
		limit = l
	}

	ctx := r.Context()
	projects, err := h.githubProjectRepo.ListGithubProjects(ctx, page, limit)
	if err != nil {
		if errors.Is(err, postgres.ErrDatabaseDisabled) {
			writeError(w, http.StatusServiceUnavailable, "service_unavailable", "database is unavailable")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal_error", "failed to retrieve github projects")
		return
	}

	response := dto.GithubProjectsListResponse{
		Status: "success",
		Data:   dto.FromDomainGithubProjects(projects),
	}

	writeJSON(w, http.StatusOK, response, true)
}
