# AI Specification: Knowledge Indexing Pipeline

- **Purpose**: Outline the Go backend knowledge indexing pipeline (Phase 5).
- **Current Phase**: Phase 0 (Bootstrap Foundation).

---

## 1. Indexing Flow (Phase 5 Target)
1. Read `data/ai/allowlist.json`.
2. Ingest approved repository metadata and README files via GitHub API.
3. Normalize and chunk markdown documents (chunk size ~512 tokens with overlap).
4. Generate vector embeddings using provider-agnostic embedding interface.
5. Persist chunks and vectors to PostgreSQL `knowledge_chunks` table with pgvector index.
6. Record full source traceability (repository, branch, commit SHA, file path, indexed timestamp).
