package tests

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/nachsyas/arham-porto/apps/api/internal/config"
	"github.com/nachsyas/arham-porto/apps/api/internal/repository/jsonfile"
)

func getCanonicalDataDir(t *testing.T) string {
	t.Helper()
	cfg := config.Load()
	dataDir, err := cfg.ResolveDataDir()
	if err != nil {
		t.Fatalf("failed to resolve canonical data dir: %v", err)
	}
	return dataDir
}

func TestRepository_LoadCanonicalData(t *testing.T) {
	dataDir := getCanonicalDataDir(t)
	repo, err := jsonfile.LoadRepository(dataDir)
	if err != nil {
		t.Fatalf("expected canonical data to load cleanly, got error: %v", err)
	}

	ctx := context.Background()

	// 1. Verify Profile
	profile, err := repo.GetProfile(ctx)
	if err != nil {
		t.Fatalf("GetProfile failed: %v", err)
	}
	if profile.FullName != "Nachsyas Arham Mumtaz Nashohi" {
		t.Errorf("unexpected full name: %s", profile.FullName)
	}
	if profile.Role != "Software Engineer" {
		t.Errorf("unexpected role: %s", profile.Role)
	}

	// 2. Verify Projects
	projects, err := repo.ListProjects(ctx, nil, nil)
	if err != nil {
		t.Fatalf("ListProjects failed: %v", err)
	}
	if len(projects) != 4 {
		t.Errorf("expected 4 canonical projects, got %d", len(projects))
	}

	// Test GetProjectBySlug
	edutrace, err := repo.GetProjectBySlug(ctx, "edutrace")
	if err != nil || edutrace == nil {
		t.Fatalf("expected to find edutrace project, got: %v", err)
	}
	if edutrace.Title != "EduTrace" {
		t.Errorf("expected title 'EduTrace', got '%s'", edutrace.Title)
	}

	// 3. Verify Skills
	skills, err := repo.ListSkills(ctx, nil)
	if err != nil {
		t.Fatalf("ListSkills failed: %v", err)
	}
	if len(skills) == 0 {
		t.Error("expected non-empty skills list")
	}

	// 4. Verify Evidence
	evidence, err := repo.ListEvidence(ctx, nil, nil)
	if err != nil {
		t.Fatalf("ListEvidence failed: %v", err)
	}
	if len(evidence) != 4 {
		t.Errorf("expected 4 canonical evidence items, got %d", len(evidence))
	}

	// Test GetEvidenceByID
	ev, err := repo.GetEvidenceByID(ctx, "evidence-edutrace-repo")
	if err != nil || ev == nil {
		t.Fatalf("expected to find evidence-edutrace-repo, got: %v", err)
	}
	if !ev.Verified {
		t.Error("expected evidence to be verified")
	}

	// 5. Verify Journey
	allStops, err := repo.ListStops(ctx, false)
	if err != nil {
		t.Fatalf("ListStops failed: %v", err)
	}
	if len(allStops) != 8 {
		t.Errorf("expected 8 total stops in canonical data (including private TK), got %d", len(allStops))
	}

	publicStops, err := repo.ListStops(ctx, true)
	if err != nil {
		t.Fatalf("ListStops(publicOnly=true) failed: %v", err)
	}
	if len(publicStops) != 7 {
		t.Errorf("expected exactly 7 public stops, got %d", len(publicStops))
	}

	for _, s := range publicStops {
		if !s.Public {
			t.Errorf("stop %s is not marked public", s.ID)
		}
		if s.ID == "journey-tk" {
			t.Errorf("private stop journey-tk must not appear in public stops")
		}
	}
}

func TestRepository_Filtering(t *testing.T) {
	dataDir := getCanonicalDataDir(t)
	repo, err := jsonfile.LoadRepository(dataDir)
	if err != nil {
		t.Fatalf("failed to load repo: %v", err)
	}

	ctx := context.Background()

	// Project filter by category
	catBackend := "Backend"
	backendProjects, err := repo.ListProjects(ctx, &catBackend, nil)
	if err != nil {
		t.Fatalf("ListProjects error: %v", err)
	}
	if len(backendProjects) != 1 || backendProjects[0].Slug != "gdgoc-ecommerce" {
		t.Errorf("expected 1 backend project (gdgoc-ecommerce), got %d", len(backendProjects))
	}

	// Project filter by featured
	featTrue := true
	featuredProjects, err := repo.ListProjects(ctx, nil, &featTrue)
	if err != nil {
		t.Fatalf("ListProjects error: %v", err)
	}
	if len(featuredProjects) != 4 {
		t.Errorf("expected 4 featured projects, got %d", len(featuredProjects))
	}

	// Evidence filter by skillID
	sid := "backend-go"
	goEvidence, err := repo.ListEvidence(ctx, &sid, nil)
	if err != nil {
		t.Fatalf("ListEvidence error: %v", err)
	}
	if len(goEvidence) < 2 {
		t.Errorf("expected at least 2 Go evidence items, got %d", len(goEvidence))
	}

	// Evidence filter by projectID
	pid := "edutrace"
	edutraceEvidence, err := repo.ListEvidence(ctx, nil, &pid)
	if err != nil {
		t.Fatalf("ListEvidence error: %v", err)
	}
	if len(edutraceEvidence) != 1 || edutraceEvidence[0].ID != "evidence-edutrace-repo" {
		t.Errorf("expected 1 evidence item for edutrace, got %d", len(edutraceEvidence))
	}
}

func TestRepository_DuplicateProjectSlugFails(t *testing.T) {
	tempDir := t.TempDir()

	// Setup valid profile, skills, evidence, journey
	copyDir(t, getCanonicalDataDir(t), tempDir)

	// Overwrite projects.json with duplicate slugs
	badProjects := `{
		"projects": [
			{"id": "p1", "title": "Project 1", "slug": "dup-slug", "category": "AI"},
			{"id": "p2", "title": "Project 2", "slug": "dup-slug", "category": "AI"}
		]
	}`
	err := os.WriteFile(filepath.Join(tempDir, "projects", "projects.json"), []byte(badProjects), 0644)
	if err != nil {
		t.Fatalf("failed to write bad projects.json: %v", err)
	}

	_, err = jsonfile.LoadRepository(tempDir)
	if err == nil {
		t.Fatal("expected error on duplicate project slug, but got nil")
	}
}

func copyDir(t *testing.T, src, dst string) {
	t.Helper()
	dirs := []string{"profile", "projects", "skills", "evidence", "journey"}
	files := map[string]string{
		"profile":  "profile.json",
		"projects": "projects.json",
		"skills":   "skills.json",
		"evidence": "evidence.json",
		"journey":  "journey.json",
	}

	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(dst, d), 0755); err != nil {
			t.Fatalf("failed to create dir: %v", err)
		}
		f := files[d]
		srcFile := filepath.Join(src, d, f)
		dstFile := filepath.Join(dst, d, f)
		data, err := os.ReadFile(srcFile)
		if err != nil {
			t.Fatalf("failed to read %s: %v", srcFile, err)
		}
		if err := os.WriteFile(dstFile, data, 0644); err != nil {
			t.Fatalf("failed to write %s: %v", dstFile, err)
		}
	}
}
