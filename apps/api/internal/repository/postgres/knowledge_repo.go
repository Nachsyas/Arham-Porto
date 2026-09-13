package postgres

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/nachsyas/arham-porto/apps/api/internal/domain"
)

// KnowledgeRepository implements domain.KnowledgeRepository using PostgreSQL + pgvector.
type KnowledgeRepository struct {
	client *Client
}

// NewKnowledgeRepository creates a new PostgreSQL knowledge repository instance.
func NewKnowledgeRepository(client *Client) *KnowledgeRepository {
	return &KnowledgeRepository{client: client}
}

// FormatVector converts a float32 slice into pgvector string format "[v1,v2,...]".
func FormatVector(vec []float32) string {
	if len(vec) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteByte('[')
	for i, v := range vec {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(strconv.FormatFloat(float64(v), 'f', -1, 32))
	}
	b.WriteByte(']')
	return b.String()
}

// ReplaceRepositorySnapshot atomically replaces all sources and chunks for a repository in one transaction (Gate #19, #20, #48).
// If any validation or insert fails, rollback occurs and the previous healthy snapshot remains intact.
func (r *KnowledgeRepository) ReplaceRepositorySnapshot(ctx context.Context, repository string, commitSHA string, items []domain.SourceWithChunks) error {
	pool := r.client.Pool()
	if pool == nil {
		return ErrDatabaseDisabled
	}

	if strings.TrimSpace(repository) == "" {
		return errors.New("repository name cannot be empty")
	}

	// 1. Pre-validation of complete prepared snapshot
	for _, item := range items {
		if item.Source.ID == "" || item.Source.Repository == "" || item.Source.Path == "" {
			return fmt.Errorf("invalid source in snapshot: missing required fields for source %s", item.Source.Path)
		}
		for _, chunk := range item.Chunks {
			if chunk.ID == "" || chunk.SourceID == "" || chunk.Content == "" {
				return fmt.Errorf("invalid chunk in snapshot: missing required fields for chunk %s", chunk.ID)
			}
		}
	}

	// 2. Begin atomic transaction
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin snapshot transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// 3. Delete superseded active state for this repository (cascades to knowledge_chunks)
	deleteQuery := `DELETE FROM knowledge_sources WHERE repository = $1`
	if _, err := tx.Exec(ctx, deleteQuery, repository); err != nil {
		return fmt.Errorf("failed to remove existing snapshot for repo %s: %w", repository, err)
	}

	// 4. Insert all new sources and chunks
	insertSourceQuery := `
		INSERT INTO knowledge_sources (
			id, repository, repository_url, ref, commit_sha, path, title,
			source_type, checksum, approval_status, indexed_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)
	`

	insertChunkQuery := `
		INSERT INTO knowledge_chunks (
			id, source_id, chunk_index, content, checksum,
			embedding, embedding_provider, embedding_model, embedding_dimensions,
			project_id, skill_ids, evidence_id, created_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6::vector, $7, $8, $9,
			$10, $11, $12, $13
		)
	`

	for _, item := range items {
		s := item.Source
		_, err := tx.Exec(ctx, insertSourceQuery,
			s.ID, s.Repository, s.RepositoryURL, s.Ref, s.CommitSHA, s.Path, s.Title,
			s.SourceType, s.Checksum, s.ApprovalStatus, s.IndexedAt, s.CreatedAt, s.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert knowledge source %s: %w", s.Path, err)
		}

		for _, c := range item.Chunks {
			var vecParam *string
			var provParam, modelParam *string
			var dimParam *int

			if len(c.Embedding) > 0 {
				formatted := FormatVector(c.Embedding)
				vecParam = &formatted
				provParam = &c.EmbeddingProvider
				modelParam = &c.EmbeddingModel
				dimParam = &c.EmbeddingDimensions
			}

			_, err := tx.Exec(ctx, insertChunkQuery,
				c.ID, c.SourceID, c.ChunkIndex, c.Content, c.Checksum,
				vecParam, provParam, modelParam, dimParam,
				c.ProjectID, c.SkillIDs, c.EvidenceID, c.CreatedAt,
			)
			if err != nil {
				return fmt.Errorf("failed to insert chunk %s for source %s: %w", c.ID, s.Path, err)
			}
		}
	}

	// 5. Commit atomic transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit snapshot transaction: %w", err)
	}

	return nil
}

// SearchSimilar performs exact pgvector cosine distance nearest-neighbor retrieval (Gate #5, #7, #40, #41).
// Strictly filters by embedding identity (provider, model, dimensions) before ranking to ensure compatible vector space.
func (r *KnowledgeRepository) SearchSimilar(
	ctx context.Context,
	queryVector []float32,
	provider, model string,
	dimensions int,
	repoFilter, sourceTypeFilter string,
	limit int,
) ([]domain.RetrievalResult, error) {
	pool := r.client.Pool()
	if pool == nil {
		return nil, ErrDatabaseDisabled
	}

	// Bounded limits (Gate #40): default 5, max 20
	if limit <= 0 {
		limit = 5
	}
	if limit > 20 {
		limit = 20
	}

	formattedVector := FormatVector(queryVector)
	if formattedVector == "" {
		return nil, errors.New("empty query vector")
	}

	query := `
		SELECT
			c.id, c.source_id, s.repository, s.path, s.ref, s.commit_sha,
			s.repository_url, s.title, s.source_type, c.content,
			c.project_id, c.skill_ids, c.evidence_id,
			(c.embedding <=> $1::vector) AS distance,
			(1.0 - (c.embedding <=> $1::vector)) AS similarity
		FROM knowledge_chunks c
		JOIN knowledge_sources s ON c.source_id = s.id
		WHERE c.embedding IS NOT NULL
		  AND c.embedding_provider = $2
		  AND c.embedding_model = $3
		  AND c.embedding_dimensions = $4
		  AND ($5 = '' OR s.repository = $5)
		  AND ($6 = '' OR s.source_type = $6)
		ORDER BY c.embedding <=> $1::vector ASC
		LIMIT $7
	`

	rows, err := pool.Query(ctx, query,
		formattedVector,
		provider,
		model,
		dimensions,
		strings.TrimSpace(repoFilter),
		strings.TrimSpace(sourceTypeFilter),
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("exact vector search query failed: %w", err)
	}
	defer rows.Close()

	var results []domain.RetrievalResult
	for rows.Next() {
		var res domain.RetrievalResult
		var skillIDs []string
		err := rows.Scan(
			&res.ChunkID, &res.SourceID, &res.Repository, &res.Path, &res.Ref, &res.CommitSHA,
			&res.SourceURL, &res.SourceTitle, &res.SourceType, &res.Content,
			&res.ProjectID, &skillIDs, &res.EvidenceID,
			&res.Distance, &res.Similarity,
		)
		if err != nil {
			return nil, fmt.Errorf("failed scanning retrieval result row: %w", err)
		}
		res.SkillIDs = skillIDs
		results = append(results, res)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error in vector search: %w", err)
	}

	return results, nil
}

// DeleteRepositorySources removes all sources and chunks for a given repository (Gate #22).
// Must succeed even if the repository is already revoked from the canonical allowlist.
func (r *KnowledgeRepository) DeleteRepositorySources(ctx context.Context, repository string) error {
	pool := r.client.Pool()
	if pool == nil {
		return ErrDatabaseDisabled
	}

	cleanRepo := strings.TrimSpace(repository)
	if cleanRepo == "" {
		return errors.New("repository name cannot be empty")
	}

	query := `DELETE FROM knowledge_sources WHERE repository = $1`
	_, err := pool.Exec(ctx, query, cleanRepo)
	if err != nil {
		return fmt.Errorf("failed deleting repository %s sources: %w", cleanRepo, err)
	}

	return nil
}

// GetIndexStatus reports the operational index state per repository, distinguishing metadata_only vs vector_ready (Gate #51).
func (r *KnowledgeRepository) GetIndexStatus(ctx context.Context) ([]domain.RepositoryIndexStatus, error) {
	pool := r.client.Pool()
	if pool == nil {
		return nil, ErrDatabaseDisabled
	}

	query := `
		SELECT
			s.repository,
			COALESCE(MAX(s.commit_sha), '') AS commit_sha,
			COUNT(DISTINCT s.id) AS source_count,
			COUNT(c.id) AS chunk_count,
			COALESCE(MAX(c.embedding_provider), 'none') AS embedding_provider,
			COALESCE(MAX(c.embedding_model), 'none') AS embedding_model,
			COUNT(c.embedding) AS embedded_chunk_count,
			MAX(s.indexed_at) AS last_indexed_at
		FROM knowledge_sources s
		LEFT JOIN knowledge_chunks c ON s.id = c.source_id
		GROUP BY s.repository
		ORDER BY s.repository ASC
	`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query index status: %w", err)
	}
	defer rows.Close()

	var statuses []domain.RepositoryIndexStatus
	for rows.Next() {
		var st domain.RepositoryIndexStatus
		var embeddedCount int
		err := rows.Scan(
			&st.Repository,
			&st.CommitSHA,
			&st.SourceCount,
			&st.ChunkCount,
			&st.EmbeddingProvider,
			&st.EmbeddingModel,
			&embeddedCount,
			&st.LastIndexedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed scanning index status: %w", err)
		}

		if st.ChunkCount > 0 && embeddedCount == st.ChunkCount {
			st.Status = "vector_ready"
		} else if st.ChunkCount > 0 && embeddedCount == 0 {
			st.Status = "metadata_only"
		} else if st.ChunkCount > 0 {
			st.Status = "partially_embedded"
		} else {
			st.Status = "empty"
		}

		statuses = append(statuses, st)
	}

	return statuses, nil
}

// Compile-time check that KnowledgeRepository implements domain.KnowledgeRepository.
var _ domain.KnowledgeRepository = (*KnowledgeRepository)(nil)
