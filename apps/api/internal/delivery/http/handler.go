package http

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/nachsyas/arham-porto/apps/api/internal/delivery/http/dto"
	"github.com/nachsyas/arham-porto/apps/api/internal/repository/jsonfile"
	"github.com/nachsyas/arham-porto/apps/api/internal/usecase"
)

var validSlugRegex = regexp.MustCompile(`^[a-z0-9-]+$`)

// DBStatusChecker defines the interface for obtaining current database health status.
type DBStatusChecker interface {
	Status(ctx context.Context) string
}

// Handler bundles all HTTP usecases and dependencies.
type Handler struct {
	profileUC  usecase.ProfileUseCase
	projectUC  usecase.ProjectUseCase
	skillUC    usecase.SkillUseCase
	evidenceUC usecase.EvidenceUseCase
	journeyUC  usecase.JourneyUseCase
	dbChecker  DBStatusChecker

	// Phase 6 Ask Arham AI dependencies
	askUC            usecase.AskUseCase
	aiMode           string
	aiRateLimiter    *RateLimiter
	aiSemaphore      chan struct{}
	aiTimeoutSeconds int
}

// NewHandler constructs a Handler with all usecases.
func NewHandler(
	profileUC usecase.ProfileUseCase,
	projectUC usecase.ProjectUseCase,
	skillUC usecase.SkillUseCase,
	evidenceUC usecase.EvidenceUseCase,
	journeyUC usecase.JourneyUseCase,
	dbChecker DBStatusChecker,
) *Handler {
	return &Handler{
		profileUC:        profileUC,
		projectUC:        projectUC,
		skillUC:          skillUC,
		evidenceUC:       evidenceUC,
		journeyUC:        journeyUC,
		dbChecker:        dbChecker,
		aiMode:           "disabled",
		aiSemaphore:      make(chan struct{}, 4),
		aiTimeoutSeconds: 30,
	}
}

// EnableAI configures the Ask Arham AI copilot dependencies.
func (h *Handler) EnableAI(askUC usecase.AskUseCase, aiMode string, aiLimiter *RateLimiter, maxConcurrent, timeoutSecs int) {
	h.askUC = askUC
	h.aiMode = aiMode
	h.aiRateLimiter = aiLimiter
	if maxConcurrent <= 0 {
		maxConcurrent = 4
	}
	h.aiSemaphore = make(chan struct{}, maxConcurrent)
	h.aiTimeoutSeconds = timeoutSecs
}

// Healthz handles liveness probes.
func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
		"scope":  "process_alive",
	}, false)
}

// Readyz handles readiness probes according to database mode.
func (h *Handler) Readyz(w http.ResponseWriter, r *http.Request) {
	status := "disabled"
	if h.dbChecker != nil {
		status = h.dbChecker.Status(r.Context())
	}

	switch status {
	case "connected":
		writeJSON(w, http.StatusOK, map[string]string{
			"status":   "ready",
			"database": "connected",
		}, false)
	case "degraded":
		writeJSON(w, http.StatusOK, map[string]string{
			"status":   "ready",
			"database": "degraded",
		}, false)
	case "unavailable":
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{
			"status":   "not_ready",
			"database": "unavailable",
		}, false)
	default:
		writeJSON(w, http.StatusOK, map[string]string{
			"status":   "ready",
			"database": "disabled",
		}, false)
	}
}

// GetProfile returns the safe public profile.
func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	profile, err := h.profileUC.GetProfile(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "failed to retrieve profile")
		return
	}

	res := dto.FromDomainProfile(profile)
	writeJSON(w, http.StatusOK, dto.NewDataEnvelope(res), true)
}

// ListProjects handles GET /api/v1/projects with optional query filters.
func (h *Handler) ListProjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query()

	if len(q["category"]) > 1 {
		writeError(w, http.StatusBadRequest, "bad_request", "ambiguous repeated 'category' query parameter")
		return
	}
	if len(q["featured"]) > 1 {
		writeError(w, http.StatusBadRequest, "bad_request", "ambiguous repeated 'featured' query parameter")
		return
	}

	var category *string
	if cat := strings.TrimSpace(q.Get("category")); cat != "" {
		if !jsonfile.ValidProjectCategory(cat) {
			writeError(w, http.StatusBadRequest, "bad_request", "unsupported project category filter")
			return
		}
		category = &cat
	}

	var featured *bool
	if featStr := strings.TrimSpace(q.Get("featured")); featStr != "" {
		val, err := strconv.ParseBool(featStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "bad_request", "invalid 'featured' query parameter, must be 'true' or 'false'")
			return
		}
		featured = &val
	}

	projects, err := h.projectUC.ListProjects(ctx, category, featured)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "failed to list projects")
		return
	}

	res := dto.FromDomainProjects(projects)
	writeJSON(w, http.StatusOK, dto.NewListEnvelope(res), true)
}

// GetProjectBySlug handles GET /api/v1/projects/{slug}.
func (h *Handler) GetProjectBySlug(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slug := strings.TrimSpace(r.PathValue("slug"))

	if slug == "" || len(slug) > 64 || !validSlugRegex.MatchString(slug) {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid project slug")
		return
	}

	project, err := h.projectUC.GetProjectBySlug(ctx, slug)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "failed to retrieve project")
		return
	}
	if project == nil {
		writeError(w, http.StatusNotFound, "not_found", "project not found")
		return
	}

	res := dto.FromDomainProject(*project)
	writeJSON(w, http.StatusOK, dto.NewDataEnvelope(res), true)
}

// ListSkills handles GET /api/v1/skills with optional category filter.
func (h *Handler) ListSkills(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query()

	if len(q["category"]) > 1 {
		writeError(w, http.StatusBadRequest, "bad_request", "ambiguous repeated 'category' query parameter")
		return
	}

	var category *string
	if cat := strings.TrimSpace(q.Get("category")); cat != "" {
		if !jsonfile.ValidSkillCategory(cat) {
			writeError(w, http.StatusBadRequest, "bad_request", "unsupported skill category filter")
			return
		}
		category = &cat
	}

	skills, err := h.skillUC.ListSkills(ctx, category)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "failed to list skills")
		return
	}

	res := dto.FromDomainSkills(skills)
	writeJSON(w, http.StatusOK, dto.NewListEnvelope(res), true)
}

// ListEvidence handles GET /api/v1/evidence with optional skill_id and project_id filters.
func (h *Handler) ListEvidence(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query()

	var skillID *string
	if sid := strings.TrimSpace(q.Get("skill_id")); sid != "" {
		if len(sid) > 64 {
			writeError(w, http.StatusBadRequest, "bad_request", "invalid skill_id filter")
			return
		}
		skillID = &sid
	}

	var projectID *string
	if pid := strings.TrimSpace(q.Get("project_id")); pid != "" {
		if len(pid) > 64 {
			writeError(w, http.StatusBadRequest, "bad_request", "invalid project_id filter")
			return
		}
		projectID = &pid
	}

	evidence, err := h.evidenceUC.ListEvidence(ctx, skillID, projectID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "failed to list evidence")
		return
	}

	res := dto.FromDomainEvidenceList(evidence)
	writeJSON(w, http.StatusOK, dto.NewListEnvelope(res), true)
}

// GetEvidenceByID handles GET /api/v1/evidence/{id}.
func (h *Handler) GetEvidenceByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := strings.TrimSpace(r.PathValue("id"))

	if id == "" || len(id) > 64 {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid evidence id")
		return
	}

	evidence, err := h.evidenceUC.GetEvidenceByID(ctx, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "failed to retrieve evidence")
		return
	}
	if evidence == nil {
		writeError(w, http.StatusNotFound, "not_found", "evidence not found")
		return
	}

	res := dto.FromDomainEvidence(*evidence)
	writeJSON(w, http.StatusOK, dto.NewDataEnvelope(res), true)
}

// ListJourney handles GET /api/v1/journey returning public milestones only.
func (h *Handler) ListJourney(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	stops, err := h.journeyUC.ListPublicStops(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "failed to list journey stops")
		return
	}

	res := dto.FromDomainJourneyStops(stops)
	writeJSON(w, http.StatusOK, dto.NewListEnvelope(res), true)
}

// NotFound handles requests to unregistered routes with a safe JSON error envelope.
func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotFound, "not_found", "the requested resource was not found")
}

func writeJSON(w http.ResponseWriter, status int, data any, cacheable bool) {
	setSecurityHeaders(w)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if cacheable {
		w.Header().Set("Cache-Control", "public, max-age=60, stale-while-revalidate=300")
	} else {
		w.Header().Set("Cache-Control", "no-store")
	}
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	setSecurityHeaders(w)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(dto.NewErrorEnvelope(code, message))
}
