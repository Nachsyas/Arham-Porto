package http

import (
	"encoding/json"
	"net/http"

	"github.com/nachsyas/arham-porto/apps/api/internal/usecase"
)

// Handler bundles HTTP handler dependencies.
type Handler struct {
	profileUC  usecase.ProfileUseCase
	projectUC  usecase.ProjectUseCase
	evidenceUC usecase.EvidenceUseCase
}

// NewHandler constructs a Handler.
func NewHandler(profileUC usecase.ProfileUseCase, projectUC usecase.ProjectUseCase, evidenceUC usecase.EvidenceUseCase) *Handler {
	return &Handler{
		profileUC:  profileUC,
		projectUC:  projectUC,
		evidenceUC: evidenceUC,
	}
}

// Healthz checks process liveness.
func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
		"scope":  "process_alive",
	})
}

// Readyz checks application readiness to serve traffic.
// In Phase 0, validates server initialization. In Phase 4, validates DB connectivity.
func (h *Handler) Readyz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ready",
		"scope":  "server_initialized",
	})
}

// GetProfile returns the portfolio owner's profile.
func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	profile, err := h.profileUC.GetProfile(ctx)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, profile)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
