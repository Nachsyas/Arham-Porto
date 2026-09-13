package indexing

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/nachsyas/arham-porto/apps/api/internal/domain"
	"github.com/nachsyas/arham-porto/apps/api/internal/embedding"
	"github.com/nachsyas/arham-porto/apps/api/internal/github"
)

var (
	ErrNotAllowlisted = errors.New("repository is not in canonical allowlist or is disabled")
	ErrRepoNotFound   = errors.New("repository manifest not found")
)

// RepositoryPlanReport represents offline policy inspection (--plan-only).
type RepositoryPlanReport struct {
	Repository string          `json:"repository"`
	Ref        string          `json:"ref"`
	Entries    []ManifestEntry `json:"entries"`
	Note       string          `json:"note"`
}

// FileProcessingDetail records the outcome of inspecting or fetching a single manifest entry.
type FileProcessingDetail struct {
	Path       string `json:"path"`
	Title      string `json:"title"`
	SourceType string `json:"source_type"`
	Status     string `json:"status"` // selected, skipped, missing_required
	Reason     string `json:"reason,omitempty"`
	Bytes      int    `json:"bytes,omitempty"`
	Chunks     int    `json:"chunks,omitempty"`
	Checksum   string `json:"checksum,omitempty"`
}

// RepositoryDryRunReport records read-only network acquisition, chunking, and validation (--dry-run).
type RepositoryDryRunReport struct {
	Repository    string                 `json:"repository"`
	CommitSHA     string                 `json:"commit_sha"`
	SelectedFiles int                    `json:"selected_files"`
	SkippedFiles  int                    `json:"skipped_files"`
	TotalBytes    int                    `json:"total_bytes"`
	TotalChunks   int                    `json:"total_chunks"`
	Files         []FileProcessingDetail `json:"files"`
}

// RepositoryIndexReport records the result of full indexing and atomic snapshot replacement.
type RepositoryIndexReport struct {
	Repository        string `json:"repository"`
	CommitSHA         string `json:"commit_sha"`
	SelectedFiles     int    `json:"selected_files"`
	SkippedFiles      int    `json:"skipped_files"`
	TotalBytes        int    `json:"total_bytes"`
	TotalChunks       int    `json:"total_chunks"`
	EmbeddingStatus   string `json:"embedding_status"`   // embedded, disabled
	PersistenceStatus string `json:"persistence_status"` // committed, dry_run, failed
}

// Service orchestrates canonical knowledge and approved repository ingestion.
type Service struct {
	dataDir           string
	allowlist         []AllowlistEntry
	ghClient          *github.Client
	embeddingProvider embedding.Provider
	knowledgeRepo     domain.KnowledgeRepository
}

// NewService creates a new indexing service instance.
func NewService(
	dataDir string,
	allowlist []AllowlistEntry,
	ghClient *github.Client,
	embeddingProvider embedding.Provider,
	knowledgeRepo domain.KnowledgeRepository,
) *Service {
	return &Service{
		dataDir:           dataDir,
		allowlist:         allowlist,
		ghClient:          ghClient,
		embeddingProvider: embeddingProvider,
		knowledgeRepo:     knowledgeRepo,
	}
}

// PlanOnly returns offline repository manifests and policy paths without network or database access (Gate #8).
func (s *Service) PlanOnly(targetRepo string) ([]RepositoryPlanReport, error) {
	var targetRepos []string
	if targetRepo != "" {
		if !IsRepositoryApproved(s.allowlist, targetRepo) {
			return nil, fmt.Errorf("%w: %s", ErrNotAllowlisted, targetRepo)
		}
		targetRepos = []string{targetRepo}
	} else {
		targetRepos = GetApprovedRepositories(s.allowlist)
	}

	var reports []RepositoryPlanReport
	for _, repo := range targetRepos {
		manifest, ok := GetRepositoryManifest(repo)
		if !ok {
			reports = append(reports, RepositoryPlanReport{
				Repository: repo,
				Note:       "No static manifest defined for this repository.",
			})
			continue
		}

		reports = append(reports, RepositoryPlanReport{
			Repository: manifest.Repository,
			Ref:        manifest.Ref,
			Entries:    manifest.Entries,
			Note:       "Offline plan-only mode: commit SHA, file availability, chunk counts, and checksums are not resolved.",
		})
	}

	return reports, nil
}

// DryRun performs read-only network acquisition, normalization, chunking, and validation (Gate #8, #9).
// Does NOT require PostgreSQL and makes ZERO database mutations and ZERO embedding API calls.
func (s *Service) DryRun(ctx context.Context, targetRepo string) ([]RepositoryDryRunReport, error) {
	if s.ghClient == nil {
		return nil, errors.New("github client is required for dry-run network acquisition")
	}

	var targetRepos []string
	if targetRepo != "" {
		if !IsRepositoryApproved(s.allowlist, targetRepo) {
			return nil, fmt.Errorf("%w: %s", ErrNotAllowlisted, targetRepo)
		}
		targetRepos = []string{targetRepo}
	} else {
		targetRepos = GetApprovedRepositories(s.allowlist)
	}

	var reports []RepositoryDryRunReport

	for _, repoName := range targetRepos {
		manifest, ok := GetRepositoryManifest(repoName)
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrRepoNotFound, repoName)
		}

		parts := strings.Split(repoName, "/")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid repository format: %s", repoName)
		}
		owner, repo := parts[0], parts[1]

		// Resolve immutable commit SHA (Gate #39)
		commitSHA, err := s.ghClient.ResolveCommitSHA(ctx, owner, repo, manifest.Ref)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve commit SHA for %s: %w", repoName, err)
		}

		rep := RepositoryDryRunReport{
			Repository: repoName,
			CommitSHA:  commitSHA,
		}

		for _, entry := range manifest.Entries {
			// Check defense-in-depth exclusions
			if excluded, reason := IsPathExcluded(entry.Path); excluded {
				rep.SkippedFiles++
				rep.Files = append(rep.Files, FileProcessingDetail{
					Path:       entry.Path,
					Title:      entry.Title,
					SourceType: entry.SourceType,
					Status:     "skipped",
					Reason:     reason,
				})
				continue
			}

			// Fetch file bytes at resolved commit SHA (Gate #14)
			contentBytes, err := s.ghClient.FetchRawFile(ctx, owner, repo, commitSHA, entry.Path)
			if err != nil {
				if errors.Is(err, github.ErrFileNotFound) {
					if entry.Required {
						rep.SkippedFiles++
						rep.Files = append(rep.Files, FileProcessingDetail{
							Path:       entry.Path,
							Title:      entry.Title,
							SourceType: entry.SourceType,
							Status:     "missing_required",
							Reason:     "required anchor file not found on github at commit",
						})
					} else {
						rep.SkippedFiles++
						rep.Files = append(rep.Files, FileProcessingDetail{
							Path:       entry.Path,
							Title:      entry.Title,
							SourceType: entry.SourceType,
							Status:     "skipped",
							Reason:     "optional path not present in repository",
						})
					}
					continue
				}

				rep.SkippedFiles++
				rep.Files = append(rep.Files, FileProcessingDetail{
					Path:       entry.Path,
					Title:      entry.Title,
					SourceType: entry.SourceType,
					Status:     "skipped",
					Reason:     err.Error(),
				})
				continue
			}

			// Normalize and chunk
			normText := NormalizeText(string(contentBytes))
			sourceID := GenerateSourceID(repoName, entry.Path, commitSHA)
			source := domain.KnowledgeSource{
				ID:         sourceID,
				Repository: repoName,
				Path:       entry.Path,
				Ref:        manifest.Ref,
				CommitSHA:  commitSHA,
				SourceType: entry.SourceType,
			}

			chunks := ChunkDocument(source, normText, nil, nil, nil)
			checksum := CalculateStringSHA256(normText)

			rep.SelectedFiles++
			rep.TotalBytes += len(contentBytes)
			rep.TotalChunks += len(chunks)
			rep.Files = append(rep.Files, FileProcessingDetail{
				Path:       entry.Path,
				Title:      entry.Title,
				SourceType: entry.SourceType,
				Status:     "selected",
				Bytes:      len(contentBytes),
				Chunks:     len(chunks),
				Checksum:   checksum,
			})
		}

		reports = append(reports, rep)
	}

	return reports, nil
}

// PrepareRepositorySnapshot fetches, normalizes, chunks, and optionally embeds a repository's manifest sources.
func (s *Service) PrepareRepositorySnapshot(ctx context.Context, repoName string) (*domain.RepositoryIndexStatus, []domain.SourceWithChunks, error) {
	if !IsRepositoryApproved(s.allowlist, repoName) {
		return nil, nil, fmt.Errorf("%w: %s", ErrNotAllowlisted, repoName)
	}

	manifest, ok := GetRepositoryManifest(repoName)
	if !ok {
		return nil, nil, fmt.Errorf("%w: %s", ErrRepoNotFound, repoName)
	}

	parts := strings.Split(repoName, "/")
	if len(parts) != 2 {
		return nil, nil, fmt.Errorf("invalid repository format: %s", repoName)
	}
	owner, repo := parts[0], parts[1]

	commitSHA, err := s.ghClient.ResolveCommitSHA(ctx, owner, repo, manifest.Ref)
	if err != nil {
		return nil, nil, fmt.Errorf("failed resolving commit SHA for %s: %w", repoName, err)
	}

	now := time.Now().UTC()
	var preparedItems []domain.SourceWithChunks
	var allChunksForEmbedding []*domain.KnowledgeChunk

	for _, entry := range manifest.Entries {
		if excluded, _ := IsPathExcluded(entry.Path); excluded {
			continue
		}

		contentBytes, err := s.ghClient.FetchRawFile(ctx, owner, repo, commitSHA, entry.Path)
		if err != nil {
			if errors.Is(err, github.ErrFileNotFound) && !entry.Required {
				continue
			}
			if entry.Required {
				return nil, nil, fmt.Errorf("missing critical anchor file %s for repo %s: %w", entry.Path, repoName, err)
			}
			continue
		}

		normText := NormalizeText(string(contentBytes))
		sourceURL, _ := github.GenerateCitationURL(owner, repo, commitSHA, entry.Path)
		sourceID := GenerateSourceID(repoName, entry.Path, commitSHA)
		source := domain.KnowledgeSource{
			ID:             sourceID,
			Repository:     repoName,
			RepositoryURL:  sourceURL,
			Ref:            manifest.Ref,
			CommitSHA:      commitSHA,
			Path:           entry.Path,
			Title:          entry.Title,
			SourceType:     entry.SourceType,
			Checksum:       CalculateStringSHA256(normText),
			ApprovalStatus: "approved",
			IndexedAt:      now,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		chunks := ChunkDocument(source, normText, nil, nil, nil)
		preparedItems = append(preparedItems, domain.SourceWithChunks{
			Source: source,
			Chunks: chunks,
		})

		for i := range preparedItems[len(preparedItems)-1].Chunks {
			allChunksForEmbedding = append(allChunksForEmbedding, &preparedItems[len(preparedItems)-1].Chunks[i])
		}
	}

	// Generate embeddings if provider is enabled
	statusStr := "metadata_only"
	provName := "disabled"
	modelName := "none"
	if s.embeddingProvider != nil && s.embeddingProvider.Model() != "disabled" {
		provName = s.embeddingProvider.ProviderName()
		modelName = s.embeddingProvider.Model()
		dims := s.embeddingProvider.Dimensions()

		texts := make([]string, len(allChunksForEmbedding))
		for i, c := range allChunksForEmbedding {
			texts[i] = c.Content
		}

		// Embed using explicit PurposeDocument (Gate #3, #29)
		vectors, err := s.embeddingProvider.Embed(ctx, embedding.PurposeDocument, texts)
		if err != nil {
			return nil, nil, fmt.Errorf("embedding generation failed for repo %s: %w", repoName, err)
		}

		for i, vec := range vectors {
			allChunksForEmbedding[i].Embedding = vec
			allChunksForEmbedding[i].EmbeddingProvider = provName
			allChunksForEmbedding[i].EmbeddingModel = modelName
			allChunksForEmbedding[i].EmbeddingDimensions = dims
		}
		statusStr = "vector_ready"
	}

	status := &domain.RepositoryIndexStatus{
		Repository:        repoName,
		CommitSHA:         commitSHA,
		SourceCount:       len(preparedItems),
		ChunkCount:        len(allChunksForEmbedding),
		EmbeddingProvider: provName,
		EmbeddingModel:    modelName,
		Status:            statusStr,
		LastIndexedAt:     &now,
	}

	return status, preparedItems, nil
}

// Index executes complete acquisition, optional embedding, and atomic repository snapshot persistence (Gate #19, #20).
func (s *Service) Index(ctx context.Context, targetRepo string) ([]RepositoryIndexReport, error) {
	if s.knowledgeRepo == nil {
		return nil, errors.New("knowledge repository is required for indexing persistence")
	}

	var targetRepos []string
	if targetRepo != "" {
		if !IsRepositoryApproved(s.allowlist, targetRepo) {
			return nil, fmt.Errorf("%w: %s", ErrNotAllowlisted, targetRepo)
		}
		targetRepos = []string{targetRepo}
	} else {
		targetRepos = GetApprovedRepositories(s.allowlist)
	}

	var reports []RepositoryIndexReport

	// 1. Index approved GitHub repositories
	for _, repoName := range targetRepos {
		status, items, err := s.PrepareRepositorySnapshot(ctx, repoName)
		if err != nil {
			return nil, fmt.Errorf("failed preparing snapshot for %s: %w", repoName, err)
		}

		// Atomic snapshot replacement (Gate #19)
		err = s.knowledgeRepo.ReplaceRepositorySnapshot(ctx, repoName, status.CommitSHA, items)
		if err != nil {
			return nil, fmt.Errorf("failed persisting snapshot for %s: %w", repoName, err)
		}

		embStatus := "disabled"
		if status.Status == "vector_ready" {
			embStatus = fmt.Sprintf("%s (%s, dim=%d)", status.EmbeddingModel, status.EmbeddingProvider, s.embeddingProvider.Dimensions())
		}

		reports = append(reports, RepositoryIndexReport{
			Repository:        repoName,
			CommitSHA:         status.CommitSHA,
			SelectedFiles:     status.SourceCount,
			TotalChunks:       status.ChunkCount,
			EmbeddingStatus:   embStatus,
			PersistenceStatus: "committed",
		})
	}

	// 2. Index Canonical Portfolio Knowledge Sources (Gate #16, #17, #18)
	canonicalItems, err := BuildCanonicalKnowledgeSources(s.dataDir)
	if err == nil && len(canonicalItems) > 0 {
		var canonicalChunks []*domain.KnowledgeChunk
		for i := range canonicalItems {
			for j := range canonicalItems[i].Chunks {
				canonicalChunks = append(canonicalChunks, &canonicalItems[i].Chunks[j])
			}
		}

		if s.embeddingProvider != nil && s.embeddingProvider.Model() != "disabled" {
			texts := make([]string, len(canonicalChunks))
			for i, c := range canonicalChunks {
				texts[i] = c.Content
			}
			vectors, err := s.embeddingProvider.Embed(ctx, embedding.PurposeDocument, texts)
			if err == nil {
				provName := s.embeddingProvider.ProviderName()
				modelName := s.embeddingProvider.Model()
				dims := s.embeddingProvider.Dimensions()
				for i, vec := range vectors {
					canonicalChunks[i].Embedding = vec
					canonicalChunks[i].EmbeddingProvider = provName
					canonicalChunks[i].EmbeddingModel = modelName
					canonicalChunks[i].EmbeddingDimensions = dims
				}
			}
		}

		_ = s.knowledgeRepo.ReplaceRepositorySnapshot(ctx, "canonical", "canonical", canonicalItems)
		reports = append(reports, RepositoryIndexReport{
			Repository:        "canonical",
			CommitSHA:         "canonical",
			SelectedFiles:     len(canonicalItems),
			TotalChunks:       len(canonicalChunks),
			EmbeddingStatus:   s.embeddingProvider.Model(),
			PersistenceStatus: "committed",
		})
	}

	return reports, nil
}

// DeleteRepository removes all stored index state for an exact repository (Gate #22).
// Does NOT require the repository to remain in the canonical allowlist.
func (s *Service) DeleteRepository(ctx context.Context, repository string) error {
	if s.knowledgeRepo == nil {
		return errors.New("knowledge repository is required for repository deletion")
	}
	return s.knowledgeRepo.DeleteRepositorySources(ctx, repository)
}

// GetStatus returns the current index readiness status across all indexed sources (Gate #51).
func (s *Service) GetStatus(ctx context.Context) ([]domain.RepositoryIndexStatus, error) {
	if s.knowledgeRepo == nil {
		return nil, errors.New("knowledge repository is required for status inspection")
	}
	return s.knowledgeRepo.GetIndexStatus(ctx)
}
