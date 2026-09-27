package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/nachsyas/arham-porto/apps/api/internal/github"
	"github.com/nachsyas/arham-porto/apps/api/internal/repository/postgres"
)

// SyncGitHubRepositories handles POST /api/v1/github/sync.
// Ingests public repositories from GitHub for Nachsyas and updates the database.
func (h *Handler) SyncGitHubRepositories(w http.ResponseWriter, r *http.Request) {
	if h.githubSyncService == nil {
		writeError(w, http.StatusServiceUnavailable, "service_unavailable", "github synchronization service is not configured")
		return
	}

	ctx := r.Context()
	result, err := h.githubSyncService.SyncRepositories(ctx)
	if err != nil {
		if errors.Is(err, github.ErrMissingToken) {
			writeError(w, http.StatusUnauthorized, "unauthorized", "GITHUB_TOKEN is missing or not configured")
			return
		}
		if errors.Is(err, github.ErrRateLimited) {
			writeError(w, http.StatusTooManyRequests, "rate_limited", "GitHub API rate limit exceeded")
			return
		}
		if errors.Is(err, postgres.ErrDatabaseDisabled) || strings.Contains(err.Error(), "postgres:") || strings.Contains(err.Error(), "database") {
			writeError(w, http.StatusServiceUnavailable, "service_unavailable", "database is unavailable")
			return
		}
		if strings.Contains(err.Error(), "github fetch repos failed") || strings.Contains(err.Error(), "unexpected status") || errors.Is(err, github.ErrUserNotFound) {
			writeError(w, http.StatusBadGateway, "bad_gateway", "failed to fetch repositories from GitHub API")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal_error", "failed to synchronize github repositories")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"status":  "success",
		"synced":  result.Synced,
		"created": result.Created,
		"updated": result.Updated,
	}, false)
}
