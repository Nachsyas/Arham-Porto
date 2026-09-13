package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/nachsyas/arham-porto/apps/api/internal/config"
	"github.com/nachsyas/arham-porto/apps/api/internal/domain"
	"github.com/nachsyas/arham-porto/apps/api/internal/embedding"
	"github.com/nachsyas/arham-porto/apps/api/internal/embedding/gemini"
	"github.com/nachsyas/arham-porto/apps/api/internal/github"
	"github.com/nachsyas/arham-porto/apps/api/internal/indexing"
	"github.com/nachsyas/arham-porto/apps/api/internal/repository/postgres"
)

func main() {
	var (
		repoFlag        = flag.String("repo", "", "Target approved repository to index (e.g. Nachsyas/EduTrace)")
		allApprovedFlag = flag.Bool("all-approved", false, "Index all approved repositories from canonical allowlist")
		dryRunFlag      = flag.Bool("dry-run", false, "Read-only network mode: resolves commits, fetches, chunks, validates (NO DB, NO EMBEDDINGS)")
		planOnlyFlag    = flag.Bool("plan-only", false, "Offline mode: prints source manifest policies with ZERO network and ZERO DB")
		statusFlag      = flag.Bool("status", false, "Display index readiness status across all indexed sources")
		deleteRepoFlag  = flag.String("delete-repo", "", "Delete all index state for a repository (works even if revoked from allowlist)")
	)
	flag.Parse()

	// 1. Validate mutually incompatible flags (Gate #49)
	if err := validateFlags(*repoFlag, *allApprovedFlag, *dryRunFlag, *planOnlyFlag, *statusFlag, *deleteRepoFlag); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
		flag.Usage()
		os.Exit(1)
	}

	cfg := config.Load()
	dataDir, err := cfg.ResolveDataDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: failed to resolve data dir: %v\n", err)
		os.Exit(1)
	}

	// 2. Load canonical allowlist (authoritative, read-only per Gate #52, #53)
	allowlist, err := indexing.LoadAllowlist(dataDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading canonical allowlist: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()

	// Mode A: Offline Plan-Only Mode (--plan-only)
	if *planOnlyFlag {
		target := strings.TrimSpace(*repoFlag)
		svc := indexing.NewService(dataDir, allowlist, nil, nil, nil)
		reports, err := svc.PlanOnly(target)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Plan-only inspection failed: %v\n", err)
			os.Exit(1)
		}

		printPlanReport(reports)
		os.Exit(0)
	}

	// Mode B: Read-Only Network Dry-Run (--dry-run) (Gate #8, #9)
	if *dryRunFlag {
		ghClient := github.NewClient(cfg.GitHubToken)
		svc := indexing.NewService(dataDir, allowlist, ghClient, nil, nil)
		target := strings.TrimSpace(*repoFlag)

		reports, err := svc.DryRun(ctx, target)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Dry-run execution failed: %v\n", err)
			os.Exit(1)
		}

		printDryRunReport(reports)
		os.Exit(0)
	}

	// For database-backed modes (--status, --delete-repo, full indexing)
	pgClient := postgres.NewClient(cfg.DatabaseURL, cfg.DatabaseMode)
	if err := pgClient.Connect(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Database connection failed: %v\n", err)
		os.Exit(1)
	}
	defer pgClient.Close()

	knowledgeRepo := postgres.NewKnowledgeRepository(pgClient)

	// Mode C: Status Inspection (--status) (Gate #51)
	if *statusFlag {
		statuses, err := knowledgeRepo.GetIndexStatus(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed querying index status: %v\n", err)
			os.Exit(1)
		}

		printStatusReport(statuses)
		os.Exit(0)
	}

	// Mode D: Repository Deletion (--delete-repo) (Gate #22)
	if *deleteRepoFlag != "" {
		target := strings.TrimSpace(*deleteRepoFlag)
		if err := knowledgeRepo.DeleteRepositorySources(ctx, target); err != nil {
			fmt.Fprintf(os.Stderr, "Failed deleting repository %s: %v\n", target, err)
			os.Exit(1)
		}
		fmt.Printf("Successfully purged index state for repository: %s\n", target)
		os.Exit(0)
	}

	// Mode E: Full Ingestion & Indexing
	var embProvider embedding.Provider
	if cfg.EmbeddingMode == "enabled" {
		if cfg.EmbeddingProvider == "gemini" {
			p, err := gemini.NewProvider(cfg.GeminiAPIKey, cfg.EmbeddingModel, cfg.EmbeddingDimensions)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed initializing Gemini provider: %v\n", err)
				os.Exit(1)
			}
			embProvider = p
		} else if cfg.EmbeddingProvider == "fake" {
			embProvider = embedding.NewDeterministicFakeProvider(cfg.EmbeddingDimensions)
		} else {
			fmt.Fprintf(os.Stderr, "Unknown EMBEDDING_PROVIDER: %q (expected 'gemini', 'fake', or 'disabled')\n", cfg.EmbeddingProvider)
			os.Exit(1)
		}
	} else {
		embProvider = embedding.NewDisabledProvider(cfg.EmbeddingDimensions)
	}

	ghClient := github.NewClient(cfg.GitHubToken)
	svc := indexing.NewService(dataDir, allowlist, ghClient, embProvider, knowledgeRepo)
	target := strings.TrimSpace(*repoFlag)

	fmt.Printf("Starting indexing pipeline (target: %s, embedding: %s)...\n", targetLabel(target), embProvider.Model())
	reports, err := svc.Index(ctx, target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Indexing pipeline failed: %v\n", err)
		os.Exit(1)
	}

	printIndexReport(reports)
}

func validateFlags(repo string, allApproved, dryRun, planOnly, status bool, deleteRepo string) error {
	actionCount := 0
	if dryRun {
		actionCount++
	}
	if planOnly {
		actionCount++
	}
	if status {
		actionCount++
	}
	if deleteRepo != "" {
		actionCount++
	}

	if repo != "" && allApproved {
		return errors.New("cannot specify both --repo and --all-approved")
	}

	if status && (repo != "" || allApproved || dryRun || planOnly || deleteRepo != "") {
		return errors.New("--status cannot be combined with other target or execution flags")
	}

	if deleteRepo != "" && (repo != "" || allApproved || dryRun || planOnly || status) {
		return errors.New("--delete-repo cannot be combined with other indexing flags")
	}

	if dryRun && planOnly {
		return errors.New("cannot specify both --dry-run and --plan-only")
	}

	if !status && deleteRepo == "" && repo == "" && !allApproved {
		return errors.New("must specify either --repo <name>, --all-approved, --status, or --delete-repo <name>")
	}

	return nil
}

func targetLabel(target string) string {
	if target == "" {
		return "all approved repositories"
	}
	return target
}

func printPlanReport(reports []indexing.RepositoryPlanReport) {
	fmt.Println("================================================================================")
	fmt.Println("PHASE 5 OFFLINE SOURCE MANIFEST POLICY (--plan-only)")
	fmt.Println("Zero Network Calls | Zero Database Writes")
	fmt.Println("================================================================================")
	for _, r := range reports {
		fmt.Printf("\nRepository: %s (Ref: %s)\n", r.Repository, r.Ref)
		if r.Note != "" {
			fmt.Printf("Note: %s\n", r.Note)
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "  PATH\tTYPE\tREQUIRED\tTITLE")
		for _, e := range r.Entries {
			reqStr := "optional"
			if e.Required {
				reqStr = "REQUIRED"
			}
			fmt.Fprintf(w, "  %s\t%s\t%s\t%s\n", e.Path, e.SourceType, reqStr, e.Title)
		}
		w.Flush()
	}
	fmt.Println()
}

func printDryRunReport(reports []indexing.RepositoryDryRunReport) {
	fmt.Println("================================================================================")
	fmt.Println("PHASE 5 DRY-RUN ACQUISITION & VALIDATION REPORT")
	fmt.Println("Read-Only Network Mode | ZERO Embedding Calls | ZERO Database Writes")
	fmt.Println("================================================================================")

	totalRepos := len(reports)
	totalFiles := 0
	totalChunks := 0
	totalBytes := 0

	for _, r := range reports {
		totalFiles += r.SelectedFiles
		totalChunks += r.TotalChunks
		totalBytes += r.TotalBytes

		fmt.Printf("\nRepository: %s\n", r.Repository)
		fmt.Printf("Resolved Commit SHA: %s\n", r.CommitSHA)
		fmt.Printf("Selected Files: %d | Skipped Files: %d | Chunks: %d | Bytes: %d\n",
			r.SelectedFiles, r.SkippedFiles, r.TotalChunks, r.TotalBytes)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "  STATUS\tPATH\tBYTES\tCHUNKS\tCHECKSUM\tREASON")
		for _, f := range r.Files {
			cs := f.Checksum
			if len(cs) > 12 {
				cs = cs[:12] + "..."
			}
			reason := f.Reason
			if reason == "" {
				reason = "-"
			}
			fmt.Fprintf(w, "  %s\t%s\t%d\t%d\t%s\t%s\n", f.Status, f.Path, f.Bytes, f.Chunks, cs, reason)
		}
		w.Flush()
	}

	fmt.Println("\n--------------------------------------------------------------------------------")
	fmt.Printf("SUMMARY: %d Repositories | %d Files Selected | %d Chunks Formed | %d Total Bytes\n",
		totalRepos, totalFiles, totalChunks, totalBytes)
	fmt.Println("STATE: ZERO database mutations | ZERO embedding provider calls")
	fmt.Println("================================================================================")
}

func printStatusReport(statuses []domain.RepositoryIndexStatus) {
	fmt.Println("================================================================================")
	fmt.Println("PHASE 5 KNOWLEDGE INDEX OPERATIONAL STATUS")
	fmt.Println("================================================================================")
	if len(statuses) == 0 {
		fmt.Println("No indexed repositories found in database.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "REPOSITORY\tCOMMIT SHA\tSOURCES\tCHUNKS\tEMBEDDING MODEL\tSTATUS")
	for _, s := range statuses {
		sha := s.CommitSHA
		if len(sha) > 8 {
			sha = sha[:8]
		}
		fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%s\t%s\n",
			s.Repository, sha, s.SourceCount, s.ChunkCount, s.EmbeddingModel, s.Status)
	}
	w.Flush()
	fmt.Println("================================================================================")
}

func printIndexReport(reports []indexing.RepositoryIndexReport) {
	fmt.Println("================================================================================")
	fmt.Println("PHASE 5 INDEXING PIPELINE EXECUTION REPORT")
	fmt.Println("================================================================================")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "REPOSITORY\tCOMMIT SHA\tSOURCES\tCHUNKS\tEMBEDDING\tPERSISTENCE")
	for _, r := range reports {
		sha := r.CommitSHA
		if len(sha) > 8 {
			sha = sha[:8]
		}
		fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%s\t%s\n",
			r.Repository, sha, r.SelectedFiles, r.TotalChunks, r.EmbeddingStatus, r.PersistenceStatus)
	}
	w.Flush()
	fmt.Println("================================================================================")
}
