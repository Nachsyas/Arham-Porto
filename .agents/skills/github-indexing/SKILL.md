---
name: github-indexing
description: Guidelines for safely ingesting approved candidate repositories, chunking documents, and generating embeddings.
---

# Skill: GitHub Knowledge Indexing

## When to Use
Use during Phase 5 when building the ingestion pipeline in `apps/api`.

## Required Reading
1. `docs/ai/indexing.md`
2. `docs/ai/knowledge-sources.md`
3. `docs/SOP/08-security-privacy.md`

## Implementation Rules
- Ingest only repositories explicitly present and enabled in `data/ai/allowlist.json`.
- Ignore files containing sensitive patterns (`.env`, `credentials`, `secrets`, `.pem`).
- Store metadata alongside chunks (repository, commit SHA, branch, file path, timestamp).
- Provider-agnostic embedding interface (`EmbeddingProvider`).

## Validation Checklist
- [ ] No unapproved repositories ingested.
- [ ] Source traceability recorded for every stored chunk.
- [ ] Chunking preserves markdown code blocks and headers cleanly.
