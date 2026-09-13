package jsonfile

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/nachsyas/arham-porto/apps/api/internal/domain"
)

// Supported canonical categories and types
var (
	validProjectCategories = map[string]bool{
		"AI":         true,
		"Full-Stack": true,
		"Backend":    true,
		"Systems":    true,
	}

	validSkillCategories = map[string]bool{
		"Backend":  true,
		"Frontend": true,
		"AI / ML":  true,
		"Systems":  true,
		"DevOps":   true,
		"Database": true,
	}

	validEvidenceTypes = map[string]bool{
		"project":    true,
		"github":     true,
		"cv":         true,
		"experience": true,
		"education":  true,
	}

	validJourneyCategories = map[string]bool{
		"birthplace": true,
		"residence":  true,
		"tk":         true,
		"sd":         true,
		"smp":        true,
		"sma":        true,
		"university": true,
		"current":    true,
	}
)

// Repository manages in-memory canonical data loaded once at startup.
type Repository struct {
	mu            sync.RWMutex
	profile       *domain.Profile
	projects      []domain.Project
	projectBySlug map[string]domain.Project
	projectByID   map[string]domain.Project
	skills        []domain.Skill
	skillByID     map[string]domain.Skill
	evidence      []domain.Evidence
	evidenceByID  map[string]domain.Evidence
	journey       []domain.JourneyStop
	publicJourney []domain.JourneyStop
}

// LoadRepository loads and strictly validates canonical JSON files from dataDir.
func LoadRepository(dataDir string) (*Repository, error) {
	repo := &Repository{
		projectBySlug: make(map[string]domain.Project),
		projectByID:   make(map[string]domain.Project),
		skillByID:     make(map[string]domain.Skill),
		evidenceByID:  make(map[string]domain.Evidence),
	}

	// 1. Load Profile
	profilePath := filepath.Join(dataDir, "profile", "profile.json")
	profileData, err := os.ReadFile(profilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read profile.json at %s: %w", profilePath, err)
	}
	var profile domain.Profile
	if err := json.Unmarshal(profileData, &profile); err != nil {
		return nil, fmt.Errorf("failed to unmarshal profile.json: %w", err)
	}
	if strings.TrimSpace(profile.FullName) == "" {
		return nil, errors.New("invalid profile: fullName cannot be empty")
	}
	if strings.TrimSpace(profile.Role) == "" {
		return nil, errors.New("invalid profile: role cannot be empty")
	}
	repo.profile = &profile

	// 2. Load Projects
	projectsPath := filepath.Join(dataDir, "projects", "projects.json")
	projectsData, err := os.ReadFile(projectsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read projects.json at %s: %w", projectsPath, err)
	}
	var projectsContainer struct {
		Projects []domain.Project `json:"projects"`
	}
	if err := json.Unmarshal(projectsData, &projectsContainer); err != nil {
		return nil, fmt.Errorf("failed to unmarshal projects.json: %w", err)
	}
	for _, p := range projectsContainer.Projects {
		if strings.TrimSpace(p.ID) == "" {
			return nil, errors.New("invalid project: id cannot be empty")
		}
		if strings.TrimSpace(p.Slug) == "" {
			return nil, fmt.Errorf("invalid project %s: slug cannot be empty", p.ID)
		}
		if _, exists := repo.projectByID[p.ID]; exists {
			return nil, fmt.Errorf("duplicate project id detected: %s", p.ID)
		}
		if _, exists := repo.projectBySlug[p.Slug]; exists {
			return nil, fmt.Errorf("duplicate project slug detected: %s", p.Slug)
		}
		if p.Category != nil && !validProjectCategories[*p.Category] {
			return nil, fmt.Errorf("invalid category '%s' on project %s", *p.Category, p.ID)
		}

		repo.projects = append(repo.projects, p)
		repo.projectByID[p.ID] = p
		repo.projectBySlug[p.Slug] = p
	}

	// 3. Load Skills
	skillsPath := filepath.Join(dataDir, "skills", "skills.json")
	skillsData, err := os.ReadFile(skillsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read skills.json at %s: %w", skillsPath, err)
	}
	var skillsContainer struct {
		Skills []domain.Skill `json:"skills"`
	}
	if err := json.Unmarshal(skillsData, &skillsContainer); err != nil {
		return nil, fmt.Errorf("failed to unmarshal skills.json: %w", err)
	}
	for _, s := range skillsContainer.Skills {
		if strings.TrimSpace(s.ID) == "" {
			return nil, errors.New("invalid skill: id cannot be empty")
		}
		if strings.TrimSpace(s.Name) == "" {
			return nil, fmt.Errorf("invalid skill %s: name cannot be empty", s.ID)
		}
		if _, exists := repo.skillByID[s.ID]; exists {
			return nil, fmt.Errorf("duplicate skill id detected: %s", s.ID)
		}
		if !validSkillCategories[s.Category] {
			return nil, fmt.Errorf("invalid category '%s' on skill %s", s.Category, s.ID)
		}

		repo.skills = append(repo.skills, s)
		repo.skillByID[s.ID] = s
	}

	// 4. Load Evidence
	evidencePath := filepath.Join(dataDir, "evidence", "evidence.json")
	evidenceData, err := os.ReadFile(evidencePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read evidence.json at %s: %w", evidencePath, err)
	}
	var evidenceContainer struct {
		Items []domain.Evidence `json:"items"`
	}
	if err := json.Unmarshal(evidenceData, &evidenceContainer); err != nil {
		return nil, fmt.Errorf("failed to unmarshal evidence.json: %w", err)
	}
	for _, e := range evidenceContainer.Items {
		if strings.TrimSpace(e.ID) == "" {
			return nil, errors.New("invalid evidence: id cannot be empty")
		}
		if strings.TrimSpace(e.Title) == "" {
			return nil, fmt.Errorf("invalid evidence %s: title cannot be empty", e.ID)
		}
		if _, exists := repo.evidenceByID[e.ID]; exists {
			return nil, fmt.Errorf("duplicate evidence id detected: %s", e.ID)
		}
		if !validEvidenceTypes[e.Type] {
			return nil, fmt.Errorf("invalid type '%s' on evidence %s", e.Type, e.ID)
		}

		repo.evidence = append(repo.evidence, e)
		repo.evidenceByID[e.ID] = e
	}

	// 5. Load Journey
	journeyPath := filepath.Join(dataDir, "journey", "journey.json")
	journeyData, err := os.ReadFile(journeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read journey.json at %s: %w", journeyPath, err)
	}
	var journeyContainer struct {
		Stops []domain.JourneyStop `json:"stops"`
	}
	if err := json.Unmarshal(journeyData, &journeyContainer); err != nil {
		return nil, fmt.Errorf("failed to unmarshal journey.json: %w", err)
	}
	journeyIDs := make(map[string]bool)
	for _, j := range journeyContainer.Stops {
		if strings.TrimSpace(j.ID) == "" {
			return nil, errors.New("invalid journey stop: id cannot be empty")
		}
		if journeyIDs[j.ID] {
			return nil, fmt.Errorf("duplicate journey stop id detected: %s", j.ID)
		}
		if !validJourneyCategories[j.Category] {
			return nil, fmt.Errorf("invalid category '%s' on journey stop %s", j.Category, j.ID)
		}
		journeyIDs[j.ID] = true

		repo.journey = append(repo.journey, j)
		if j.Public {
			repo.publicJourney = append(repo.publicJourney, j)
		}
	}

	return repo, nil
}

// GetProfile returns the domain profile.
func (r *Repository) GetProfile(ctx context.Context) (*domain.Profile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.profile == nil {
		return nil, errors.New("profile not initialized")
	}
	return r.profile, nil
}

// ListProjects lists projects filtered optionally by category and featured status.
func (r *Repository) ListProjects(ctx context.Context, category *string, featured *bool) ([]domain.Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []domain.Project
	for _, p := range r.projects {
		if category != nil && (p.Category == nil || *p.Category != *category) {
			continue
		}
		if featured != nil && p.Featured != *featured {
			continue
		}
		result = append(result, p)
	}
	if result == nil {
		result = []domain.Project{}
	}
	return result, nil
}

// GetProjectBySlug retrieves a project by slug. Returns nil, nil if not found.
func (r *Repository) GetProjectBySlug(ctx context.Context, slug string) (*domain.Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if p, ok := r.projectBySlug[slug]; ok {
		return &p, nil
	}
	return nil, nil
}

// ListSkills lists verified skills filtered optionally by category.
func (r *Repository) ListSkills(ctx context.Context, category *string) ([]domain.Skill, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []domain.Skill
	for _, s := range r.skills {
		if category != nil && s.Category != *category {
			continue
		}
		result = append(result, s)
	}
	if result == nil {
		result = []domain.Skill{}
	}
	return result, nil
}

// ListEvidence lists evidence filtered optionally by skillID or projectID.
func (r *Repository) ListEvidence(ctx context.Context, skillID *string, projectID *string) ([]domain.Evidence, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// If projectID is specified, find the project's evidence IDs
	var allowedIDs map[string]bool
	if projectID != nil {
		if p, ok := r.projectByID[*projectID]; ok {
			allowedIDs = make(map[string]bool)
			for _, id := range p.EvidenceIDs {
				allowedIDs[id] = true
			}
		} else {
			return []domain.Evidence{}, nil
		}
	}

	var result []domain.Evidence
	for _, e := range r.evidence {
		if allowedIDs != nil && !allowedIDs[e.ID] {
			continue
		}
		if skillID != nil {
			hasSkill := false
			for _, sid := range e.SkillIDs {
				if sid == *skillID {
					hasSkill = true
					break
				}
			}
			if !hasSkill {
				continue
			}
		}
		result = append(result, e)
	}
	if result == nil {
		result = []domain.Evidence{}
	}
	return result, nil
}

// GetEvidenceByID retrieves evidence by ID. Returns nil, nil if not found.
func (r *Repository) GetEvidenceByID(ctx context.Context, id string) (*domain.Evidence, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if e, ok := r.evidenceByID[id]; ok {
		return &e, nil
	}
	return nil, nil
}

// ListStops returns journey stops. If publicOnly is true, returns only public stops.
func (r *Repository) ListStops(ctx context.Context, publicOnly bool) ([]domain.JourneyStop, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if publicOnly {
		result := make([]domain.JourneyStop, len(r.publicJourney))
		copy(result, r.publicJourney)
		return result, nil
	}

	result := make([]domain.JourneyStop, len(r.journey))
	copy(result, r.journey)
	return result, nil
}

// ValidProjectCategory checks if a category string is valid.
func ValidProjectCategory(cat string) bool {
	return validProjectCategories[cat]
}

// ValidSkillCategory checks if a skill category string is valid.
func ValidSkillCategory(cat string) bool {
	return validSkillCategories[cat]
}
