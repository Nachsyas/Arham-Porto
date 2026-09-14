package indexing

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/nachsyas/arham-porto/apps/api/internal/delivery/http/dto"
	"github.com/nachsyas/arham-porto/apps/api/internal/domain"
	"github.com/nachsyas/arham-porto/apps/api/internal/repository/jsonfile"
)

// BuildCanonicalKnowledgeSources extracts sanitized, public-only canonical entities
// (profile, projects, skills, evidence, journey) and transforms them into SourceWithChunks
// preserving explicit provenance and structured relationships per Gates #16, #17, #18.
func BuildCanonicalKnowledgeSources(dataDir string) ([]domain.SourceWithChunks, error) {
	repo, err := jsonfile.LoadRepository(dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load canonical repository: %w", err)
	}

	ctx := context.Background()
	var results []domain.SourceWithChunks
	now := time.Now().UTC()

	// 1. Profile (Sanitized Approved Public Fields Only - Centralized via DTO)
	profile, err := repo.GetProfile(ctx)
	if err == nil && profile != nil {
		pub := dto.FromDomainProfile(profile)
		var b strings.Builder
		b.WriteString(fmt.Sprintf("Name: %s\n", pub.FullName))
		b.WriteString(fmt.Sprintf("Role: %s\n", pub.Role))
		b.WriteString(fmt.Sprintf("Project: %s\n", pub.ProjectName))
		b.WriteString(fmt.Sprintf("AI Feature: %s\n", pub.AIFeature))
		if pub.Positioning != nil {
			b.WriteString(fmt.Sprintf("Positioning: %s\n", *pub.Positioning))
		}
		if pub.Bio != nil {
			b.WriteString(fmt.Sprintf("Bio: %s\n", *pub.Bio))
		}
		if pub.GitHub != nil {
			b.WriteString(fmt.Sprintf("GitHub: %s\n", *pub.GitHub))
		}
		if pub.CurrentCity != nil {
			b.WriteString(fmt.Sprintf("Current City: %s\n", *pub.CurrentCity))
		}
		if pub.Availability != nil {
			b.WriteString(fmt.Sprintf("Availability: %s\n", *pub.Availability))
		}

		profileText := NormalizeText(b.String())
		sourceID := GenerateSourceID("canonical", "profile/profile.json", "canonical")
		source := domain.KnowledgeSource{
			ID:             sourceID,
			Repository:     "canonical",
			RepositoryURL:  "file://data/profile/profile.json",
			Ref:            "canonical",
			CommitSHA:      "canonical",
			Path:           "profile/profile.json",
			Title:          fmt.Sprintf("Canonical Profile: %s", profile.FullName),
			SourceType:     "canonical_profile",
			Checksum:       CalculateStringSHA256(profileText),
			ApprovalStatus: "approved",
			IndexedAt:      now,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		chunks := ChunkDocument(source, profileText, nil, nil, nil)
		results = append(results, domain.SourceWithChunks{
			Source: source,
			Chunks: chunks,
		})
	}

	// 2. Projects (Sanitized Public Fields + Relational Metadata)
	projects, err := repo.ListProjects(ctx, nil, nil)
	if err == nil {
		for _, p := range projects {
			var b strings.Builder
			b.WriteString(fmt.Sprintf("Title: %s\n", p.Title))
			b.WriteString(fmt.Sprintf("Slug: %s\n", p.Slug))
			if p.Category != nil {
				b.WriteString(fmt.Sprintf("Category: %s\n", *p.Category))
			}
			if p.Summary != nil {
				b.WriteString(fmt.Sprintf("Summary: %s\n", *p.Summary))
			}
			if p.Problem != nil {
				b.WriteString(fmt.Sprintf("Problem: %s\n", *p.Problem))
			}
			if p.Solution != nil {
				b.WriteString(fmt.Sprintf("Solution: %s\n", *p.Solution))
			}
			if len(p.Role) > 0 {
				b.WriteString(fmt.Sprintf("Roles: %s\n", strings.Join(p.Role, ", ")))
			}
			if len(p.Contributions) > 0 {
				b.WriteString("Contributions:\n")
				for _, c := range p.Contributions {
					b.WriteString(fmt.Sprintf("- %s\n", c))
				}
			}
			if len(p.Technologies) > 0 {
				b.WriteString(fmt.Sprintf("Technologies: %s\n", strings.Join(p.Technologies, ", ")))
			}
			if p.GitHubURL != nil {
				b.WriteString(fmt.Sprintf("GitHub Repository: %s\n", *p.GitHubURL))
			}
			if p.DemoURL != nil {
				b.WriteString(fmt.Sprintf("Live Demo: %s\n", *p.DemoURL))
			}

			normText := NormalizeText(b.String())
			path := fmt.Sprintf("projects/%s.json", p.Slug)
			sourceID := GenerateSourceID("canonical", path, "canonical")
			source := domain.KnowledgeSource{
				ID:             sourceID,
				Repository:     "canonical",
				RepositoryURL:  fmt.Sprintf("file://data/projects/projects.json#%s", p.ID),
				Ref:            "canonical",
				CommitSHA:      "canonical",
				Path:           path,
				Title:          fmt.Sprintf("Canonical Project: %s", p.Title),
				SourceType:     "canonical_project",
				Checksum:       CalculateStringSHA256(normText),
				ApprovalStatus: "approved",
				IndexedAt:      now,
				CreatedAt:      now,
				UpdatedAt:      now,
			}

			projectID := p.ID
			var evidenceID *string
			if len(p.EvidenceIDs) > 0 {
				evidenceID = &p.EvidenceIDs[0]
			}

			chunks := ChunkDocument(source, normText, &projectID, nil, evidenceID)
			results = append(results, domain.SourceWithChunks{
				Source: source,
				Chunks: chunks,
			})
		}
	}

	// 3. Skills (Sanitized Public Fields + Relational Metadata)
	skills, err := repo.ListSkills(ctx, nil)
	if err == nil {
		for _, s := range skills {
			var b strings.Builder
			b.WriteString(fmt.Sprintf("Skill: %s\n", s.Name))
			b.WriteString(fmt.Sprintf("Category: %s\n", s.Category))
			if s.Claim != nil {
				b.WriteString(fmt.Sprintf("Claim: %s\n", *s.Claim))
			}
			if len(s.EvidenceIDs) > 0 {
				b.WriteString(fmt.Sprintf("Evidence IDs: %s\n", strings.Join(s.EvidenceIDs, ", ")))
			}

			normText := NormalizeText(b.String())
			path := fmt.Sprintf("skills/%s.json", s.ID)
			sourceID := GenerateSourceID("canonical", path, "canonical")
			source := domain.KnowledgeSource{
				ID:             sourceID,
				Repository:     "canonical",
				RepositoryURL:  fmt.Sprintf("file://data/skills/skills.json#%s", s.ID),
				Ref:            "canonical",
				CommitSHA:      "canonical",
				Path:           path,
				Title:          fmt.Sprintf("Canonical Skill: %s (%s)", s.Name, s.Category),
				SourceType:     "canonical_skill",
				Checksum:       CalculateStringSHA256(normText),
				ApprovalStatus: "approved",
				IndexedAt:      now,
				CreatedAt:      now,
				UpdatedAt:      now,
			}

			skillIDs := []string{s.ID}
			var evidenceID *string
			if len(s.EvidenceIDs) > 0 {
				evidenceID = &s.EvidenceIDs[0]
			}

			chunks := ChunkDocument(source, normText, nil, skillIDs, evidenceID)
			results = append(results, domain.SourceWithChunks{
				Source: source,
				Chunks: chunks,
			})
		}
	}

	// 4. Evidence (Sanitized Public Fields + Relational Metadata)
	evidenceList, err := repo.ListEvidence(ctx, nil, nil)
	if err == nil {
		for _, e := range evidenceList {
			var b strings.Builder
			b.WriteString(fmt.Sprintf("Title: %s\n", e.Title))
			b.WriteString(fmt.Sprintf("Type: %s\n", e.Type))
			b.WriteString(fmt.Sprintf("Summary: %s\n", e.Summary))
			b.WriteString(fmt.Sprintf("Verified: %t\n", e.Verified))
			if e.SourceURL != nil {
				b.WriteString(fmt.Sprintf("Source URL: %s\n", *e.SourceURL))
			}
			if e.SourcePath != nil {
				b.WriteString(fmt.Sprintf("Source Path: %s\n", *e.SourcePath))
			}
			if len(e.SkillIDs) > 0 {
				b.WriteString(fmt.Sprintf("Related Skills: %s\n", strings.Join(e.SkillIDs, ", ")))
			}

			normText := NormalizeText(b.String())
			path := fmt.Sprintf("evidence/%s.json", e.ID)
			sourceID := GenerateSourceID("canonical", path, "canonical")
			source := domain.KnowledgeSource{
				ID:             sourceID,
				Repository:     "canonical",
				RepositoryURL:  fmt.Sprintf("file://data/evidence/evidence.json#%s", e.ID),
				Ref:            "canonical",
				CommitSHA:      "canonical",
				Path:           path,
				Title:          fmt.Sprintf("Canonical Evidence: %s", e.Title),
				SourceType:     "canonical_evidence",
				Checksum:       CalculateStringSHA256(normText),
				ApprovalStatus: "approved",
				IndexedAt:      now,
				CreatedAt:      now,
				UpdatedAt:      now,
			}

			evidenceID := e.ID
			chunks := ChunkDocument(source, normText, nil, e.SkillIDs, &evidenceID)
			results = append(results, domain.SourceWithChunks{
				Source: source,
				Chunks: chunks,
			})
		}
	}

	// 5. Journey (PUBLIC STOPS ONLY - Strictly No Coordinates, No Private Data)
	journeyStops, err := repo.ListStops(ctx, true) // true = publicOnly
	if err == nil {
		for _, j := range journeyStops {
			if !j.Public {
				continue
			}

			var b strings.Builder
			b.WriteString(fmt.Sprintf("Category: %s\n", j.Category))
			if j.Title != nil {
				b.WriteString(fmt.Sprintf("Milestone: %s\n", *j.Title))
			}
			if j.Institution != nil {
				b.WriteString(fmt.Sprintf("Institution: %s\n", *j.Institution))
			}
			if j.City != nil {
				b.WriteString(fmt.Sprintf("City: %s\n", *j.City))
			}
			if j.Region != nil {
				b.WriteString(fmt.Sprintf("Region: %s\n", *j.Region))
			}
			b.WriteString(fmt.Sprintf("Country: %s\n", j.Country))
			if j.Period != nil {
				b.WriteString(fmt.Sprintf("Period: %s\n", *j.Period))
			}
			if j.Description != nil {
				b.WriteString(fmt.Sprintf("Description: %s\n", *j.Description))
			}

			normText := NormalizeText(b.String())
			path := fmt.Sprintf("journey/%s.json", j.ID)
			sourceID := GenerateSourceID("canonical", path, "canonical")
			title := "Canonical Journey Stop"
			if j.Title != nil {
				title = fmt.Sprintf("Canonical Journey: %s", *j.Title)
			}
			source := domain.KnowledgeSource{
				ID:             sourceID,
				Repository:     "canonical",
				RepositoryURL:  fmt.Sprintf("file://data/journey/journey.json#%s", j.ID),
				Ref:            "canonical",
				CommitSHA:      "canonical",
				Path:           path,
				Title:          title,
				SourceType:     "canonical_journey",
				Checksum:       CalculateStringSHA256(normText),
				ApprovalStatus: "approved",
				IndexedAt:      now,
				CreatedAt:      now,
				UpdatedAt:      now,
			}

			chunks := ChunkDocument(source, normText, nil, nil, nil)
			results = append(results, domain.SourceWithChunks{
				Source: source,
				Chunks: chunks,
			})
		}
	}

	return results, nil
}
