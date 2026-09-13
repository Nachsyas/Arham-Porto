-- Phase 5 Migration: Knowledge Indexing & Retrieval Tables
-- Zero-Hallucination Evidence Storage with pgvector Exact Retrieval

CREATE TABLE IF NOT EXISTS knowledge_sources (
    id VARCHAR(64) PRIMARY KEY,
    repository VARCHAR(255) NOT NULL,
    repository_url TEXT NOT NULL,
    ref VARCHAR(100) NOT NULL DEFAULT '',
    commit_sha VARCHAR(40) NOT NULL DEFAULT '',
    path TEXT NOT NULL,
    title TEXT NOT NULL,
    source_type VARCHAR(50) NOT NULL,
    checksum VARCHAR(64) NOT NULL,
    approval_status VARCHAR(50) NOT NULL DEFAULT 'approved',
    indexed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_knowledge_sources_repo_path_commit UNIQUE (repository, path, commit_sha)
);

CREATE TABLE IF NOT EXISTS knowledge_chunks (
    id VARCHAR(64) PRIMARY KEY,
    source_id VARCHAR(64) NOT NULL REFERENCES knowledge_sources(id) ON DELETE CASCADE,
    chunk_index INT NOT NULL,
    content TEXT NOT NULL,
    checksum VARCHAR(64) NOT NULL,
    embedding vector NULL,
    embedding_provider VARCHAR(50) NULL,
    embedding_model VARCHAR(100) NULL,
    embedding_dimensions INT NULL,
    project_id VARCHAR(100) NULL,
    skill_ids TEXT[] NULL,
    evidence_id VARCHAR(100) NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_knowledge_chunks_source_index UNIQUE (source_id, chunk_index)
);

-- Fast lookup indexes
CREATE INDEX IF NOT EXISTS idx_knowledge_sources_repo ON knowledge_sources(repository);
CREATE INDEX IF NOT EXISTS idx_knowledge_sources_type ON knowledge_sources(source_type);
CREATE INDEX IF NOT EXISTS idx_knowledge_chunks_source_id ON knowledge_chunks(source_id);
CREATE INDEX IF NOT EXISTS idx_knowledge_chunks_checksum ON knowledge_chunks(checksum);
CREATE INDEX IF NOT EXISTS idx_knowledge_chunks_model_filter ON knowledge_chunks(embedding_provider, embedding_model, embedding_dimensions);
