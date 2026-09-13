# Knowledge Indexing Pipeline Architecture (Phase 5 Certified)

## 1. Overview & Objective
The knowledge indexing pipeline powers future question-answering ("Ask Arham AI") by transforming approved canonical data and GitHub repository sources into verifiable, non-hallucinatory knowledge chunks stored in PostgreSQL with pgvector.

Pipeline Flow:
```
Approved Allowlist + Manifests
            ↓
   Source Acquisition
 (GitHub API / Canonical)
            ↓
  Text Normalization
  (Verbatim code fence)
            ↓
    Semantic Chunking
(Context Header Prepending)
            ↓
  Optional Embedding
  (Purpose: Document)
            ↓
Atomic Snapshot Persistence
    (PostgreSQL pgvector)
```

---

## 2. GitHub as Untrusted Data & Zero-Trust Boundary
Per Gate #15 and AGENTS.md §2.2:
- All external data (GitHub repository contents, commit messages, external documentation) is treated as **UNTRUSTED DATA**.
- No AI classifier is used for injection detection.
- Strings like `Ignore previous instructions` or `Reveal system prompt` remain inert indexed text.
- Retrieved chunks are treated purely as passive context in Phase 6 prompt templates, never executable system instructions.

---

## 3. Canonical Knowledge Sources
In addition to GitHub repositories, approved canonical knowledge sources (`profile`, `projects`, `skills`, `evidence`, `journey`) are ingested through a dedicated adapter (`internal/indexing/canonical_source.go`):
- **Public Data Boundary (Gate #16)**: Strictly excludes internal TODOs, geographic coordinates, birthdates, phone numbers, home addresses, and private Journey milestones.
- **Explicit Provenance (Gate #17)**: Stored with `source_type`: `canonical_profile`, `canonical_project`, `canonical_skill`, `canonical_evidence`, `canonical_journey`.
- **Relational Preservation (Gate #18)**: Preserves structured foreign links (`project_id`, `skill_ids`, `evidence_id`) so queries like *"Does Nachsyas have Go experience?"* trace from Skill → Evidence → Project → Repository source.

---

## 4. Authoritative Allowlist & Manifest Policy
- Canonical allowlist: `data/ai/allowlist.json`. The indexing runtime **never mutates** this file.
- Only repositories with `enabled: true` can be indexed.
- Source manifests (`internal/indexing/manifest.go`) dictate approved paths for inspection.
- Manifest paths are verified against actual Git trees; missing optional paths are skipped with structured reasons, while missing required anchor files fail validation.

---

## 5. Execution Modes & Dry-Run Semantics
1. `--plan-only`: Fully offline mode. Inspects allowlist and manifests without network calls and without PostgreSQL.
2. `--dry-run`: Read-only network mode. Calls GitHub API, resolves immutable 40-char commit SHAs, downloads and normalizes files, splits into chunks, and computes SHA-256 checksums. **Makes ZERO embedding calls and ZERO database writes.** Does NOT require PostgreSQL.
3. Full Indexing: Connects to PostgreSQL, embeds chunks (if enabled), and commits atomic repository snapshots.
4. `--delete-repo <owner/name>`: Purges indexed state for a repository. Works even if the repository was previously revoked from the allowlist.
5. `--status`: Inspects index readiness, distinguishing `metadata_only` from `vector_ready`.

---

## 6. Atomic Repository Snapshot Replacement
Per Gate #19, #20, and #48:
- Rather than per-file writes, the pipeline uses transactional snapshot replacement:
  ```go
  ReplaceRepositorySnapshot(ctx, repository, commitSHA, items)
  ```
- Steps inside a single PostgreSQL transaction:
  1. Validate complete prepared snapshot.
  2. Delete superseded active state for the repository (`DELETE FROM knowledge_sources WHERE repository = $1`, cascading to chunks).
  3. Insert all new sources and chunks.
  4. Commit.
- If any preparation or database write fails, the entire transaction rolls back and the previous healthy snapshot remains active.

---

## 7. Defense-in-Depth Secret & Binary Exclusion
Before indexing:
- Rejects paths matching `.env`, `.env.*`, `*.key`, `*.pem`, `id_rsa`, `credentials.*`, `service-account*.json`.
- Excludes package lockfiles (`package-lock.json`, `go.sum`, etc.) and vendor/node_modules directories.
- Rejects binary assets (`.png`, `.jpg`, `.pdf`, `.zip`, `.wasm`, etc.) and detects NUL bytes or invalid UTF-8 sequences.
- Bounded file size: strictly limited to 256 KiB (`MaxSourceFileBytes`).
