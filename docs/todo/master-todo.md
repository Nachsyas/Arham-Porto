# Master Development Roadmap & Progress Tracker

> Single source of truth for implementation progress of **Arham Porto**.

---

## Phases Overview

| Phase | Description | Status |
| :--- | :--- | :--- |
| **Phase 0** | **Bootstrap Foundation** | **CERTIFIED — Foundation** |
| **Phase 1** | **Reviewer-First Static Portfolio** | **CERTIFIED — Reviewer-First Portfolio** |
| **Phase 2** | **Hero Visual: Authentic Portrait** | **CERTIFIED — Authentic Portrait Hero** *(3D Hologram: CANCELLED / SUPERSEDED)* |
| **Phase 3** | **Journey Map** | **RECERTIFIED — Interactive Journey Map (Expanded 7 Milestones)** |
| **Phase 4** | **Go Backend** | **CERTIFIED — Clean Architecture REST API** |
| **Phase 5** | **AI Indexing & Retrieval** | **CERTIFIED — Evidence Indexing & pgvector Retrieval Foundation** |
| **Phase 6** | **AI Reviewer Copilot** | **CERTIFIED — Grounded Ask Arham AI Reviewer Copilot** |
| **Phase 7** | **Integration & Polish** | **DEFERRED — Post-Launch Enhancement** |
| **Phase 8** | **Production & Launch** | **CORE PORTFOLIO LIVE — ASK ARHAM AI PENDING** |

---

## Phase 0 Checklist (Certified Baseline)
- [x] Root project folder setup (`Arham Porto`)
- [x] Monorepo architecture with npm workspaces (`apps/*`, `packages/*`)
- [x] Pinned toolchain versions (`.nvmrc`, `package.json` engines)
- [x] Root `.gitignore` and `.env.example`
- [x] Docker Compose with PostgreSQL 16 + pgvector (`pgvector/pgvector:pg16`) and DRY healthcheck
- [x] AI Agent Deployment Rules (`AGENTS.md`)
- [x] Canonical design tokens package (`packages/design-tokens`)
- [x] Content schema package with Zod runtime validation (`packages/content-schema`)
- [x] Canonical JSON content files (`data/*`) with zero-trust validation script
- [x] Portfolio data access layer (`packages/portfolio-data`)
- [x] Next.js + TypeScript strict frontend workspace (`apps/web`) with smoke test page
- [x] Go API module skeleton (`apps/api`) with clean architecture and health endpoints
- [x] GitHub Actions CI workflow skeleton (`.github/workflows/ci.yml`)
- [x] Baseline git commit: `chore: bootstrap Arham Porto foundation` (`c31b393`)
- [x] Validation suite executed and certified:
  - `npm run validate:data`: 100% PASS (8/8 schemas)
  - `npm run lint --workspace=apps/web`: 100% PASS (0 warnings/errors)
  - `npm run typecheck --workspace=apps/web`: 100% PASS (0 errors)
  - `npm run test --workspace=apps/web`: 100% PASS
  - `npm run build --workspace=apps/web`: 100% PASS (Next.js production build static generation)
  - `go test -v ./...` and `go vet ./...`: 100% PASS
  - `docker compose config`: 100% PASS
  - `git status --short`: Clean working tree

---

## Phase 1 Checklist (Completed)
- [x] Global Navigation (Navbar, identity monogram, nav links, Quick Review trigger, mobile responsive drawer)
- [x] Hero Identity (Name `NACHSYAS ARHAM MUMTAZ NASHOHI`, Software Engineer role, neutral approved positioning in data layer, tags, dual CTAs, lightweight 3D wireframe stage placeholder)
- [x] Quick Review (Drawer with verified candidate profile, verified core stack, project links, verified contact; graceful omission of unverified fields; zero `TODO_USER` leaks)
- [x] Selected Work (Data-driven cards for EduTrace, GDGOC E-Commerce, Maritime AI Dashboard, Smart Kitchen)
- [x] Project Filtering (All, AI, Full-Stack, Backend, Systems with client-side reactive filtering)
- [x] Skill -> Evidence Explorer (`SKILL -> CLAIM -> EVIDENCE` architecture without arbitrary percentage bars)
- [x] Experience (Section architecture with dignified empty state when unverified)
- [x] Education (Section architecture with dignified empty state when unverified)
- [x] Contact (Verified channels only: GitHub `https://github.com/Nachsyas`, privacy note, zero fake email)
- [x] Project Case Study Dynamic Route (`/projects/[slug]`) with SSG `generateStaticParams`
- [x] Accessibility & Responsive verification (WCAG 2.2 AA, semantic HTML, skip-to-content, keyboard navigation, mobile touch targets >= 44px)
- [x] Unit & Component test suite for all Phase 1 components (10 tests passing across 3 test suites)
- [x] Final Phase 1 Walkthrough & Stop Rule enforcement (halt before Phase 2)

---

## Phase 2 Checklist — Hero Visual: Authentic Portrait (Completed)
> **Product Decision Note**: The 3D Wireframe Hologram concept has been cancelled and superseded by an authentic recruiter-first portrait photo of Nachsyas Arham Mumtaz Nashohi. All WebGL/Three.js dependencies and features have been removed from production.

- [x] **Hologram Cancellation & Cleanup**:
  - Cancelled 3D wireframe hologram and procedural mesh (`features/hologram/` deleted from production).
  - Uninstalled Three.js / R3F production packages (`three`, `@types/three`, `@react-three/fiber`, `@react-three/drei`).
  - Archived Phase 2 hologram specification to `docs/archive/hologram-3d-cancelled.md` (and updated `docs/design/hologram-3d.md` with cancellation notice).
  - Purged all public holographic terminology (`NODE:AP-01`, `SYS:MATRIX`, `MODE:2D_VECTOR`, "seated development mesh").
- [x] **Portrait Asset Extraction & Privacy Guard**:
  - Extracted authentic portrait of Nachsyas Arham Mumtaz Nashohi from user CV into `apps/web/public/images/profile/nachsyas-arham.jpg`.
  - Zero-Trust Privacy: Absolutely no personal details (full address, phone number, date of birth, or signature) extracted or published.
  - Tagged temporary extraction with `TODO_USER_HIGH_RES_PROFILE_PHOTO` for future high-resolution replacement.
- [x] **Hero Portrait Component (`apps/web/features/hero/HeroPortrait.tsx`)**:
  - Built with Next.js `<Image />` (`priority`, `fill`, `sizes`, `object-fit: cover`).
  - Semantic, meaningful alt text (`"Portrait of Nachsyas Arham Mumtaz Nashohi"`).
  - Sleek dark technical frame with `#02060B` canvas, `#07111C` surface, `#173247` border, and subtle `#26B8FF` cyan accents.
  - Technical corner markers (`┌ PROFILE`, `┐`, `└`, `┘`) and clean identity strip (`NACHSYAS ARHAM` / `SOFTWARE ENGINEER`).
  - Graceful image fallback state with user icon and technical prompt if image fails to load.
  - Restrained entrance motion respecting `(prefers-reduced-motion: reduce)`.
- [x] **Hero Section Integration (`apps/web/features/hero/HeroSection.tsx`)**:
  - Desktop 2-column layout (Left: Identity, positioning, CTAs, links; Right: Portrait frame).
  - Mobile responsive single-column layout (Name -> Role -> Positioning -> CTAs -> Portrait -> Quick Review drawer).
  - Preserved canonical design tokens and existing Quick Review drawer functionality.
- [x] **Testing & Verification**:
  - Replaced hologram tests with comprehensive `apps/web/tests/hero-portrait.test.tsx` (5 tests covering image render, alt text, role pill, corner accents, and fallback behavior).
  - 100% PASS across full unit test suite (17/17 tests passing in `apps/web`).
  - Verified zero Three.js/R3F imports in production bundle.
  - Performance: Lightweight First Load JS without WebGL overhead.
- [x] **Documentation**: Created `docs/design/hero-portrait.md` and updated `docs/todo/master-todo.md`.
- [x] **Stop Rule Enforcement**: Halt execution upon completion; await explicit user approval before Phase 3 (Journey Map).

---

## Phase 3 Checklist — Interactive Indonesia Journey Map (RECERTIFIED)
> **Recertified Scope (Phase 3 Reopen)**: Expanded to 7 public milestones across 4 unique geographic locations (Karanganyar -> Jakarta -> Salatiga -> Malang). Includes user-approved Jakarta Residence, MI Terpadu Al-Hamid, and MTsN 30 Jakarta Timur. Preserves city-level privacy, null period omission, and zero intra-city route hops.

- [x] **Approved Canonical Data Synchronization**:
  - Updated `packages/content-schema/src/journey.schema.ts` to support `"residence"` category.
  - Synced approved 7 milestones in `data/journey/journey.json`:
    1. `origin-karanganyar`: Karanganyar, Jawa Tengah (Birth year hidden)
    2. `residence-jakarta`: Jakarta, DKI Jakarta (City-level only, no home address)
    3. `mi-al-hamid-jakarta`: MI Terpadu Al-Hamid, Jakarta Timur (`period: null`)
    4. `mtsn30-jakarta`: MTsN 30 Jakarta Timur, Jakarta Timur (`period: null`)
    5. `ma-assurkati-salatiga`: MA Tahfizhul Qur'an As-Surkati, Salatiga (2019–2023)
    6. `university-uin-malang`: UIN Maulana Malik Ibrahim Malang, Computer Science (2023–Present)
    7. `current-base-malang`: Malang, East Java (Present)
  - Retained `journey-tk` as private (`public: false`).
  - Zero-trust validation passed: `npm run validate:data` (100% PASS).
- [x] **Journey Map Architecture & Implementation**:
  - Designed interactive Stylized SVG Geographic Journey Map section (`apps/web/features/journey/`).
  - Implemented 4 unique geographic pins: Karanganyar `(478.3, 203.2)`, Jakarta `(172.0, 72.0)`, Salatiga `(445.3, 178.2)`, Malang `(603.4, 239.2)`.
  - Shared Jakarta marker (Residence, MI, MTs) with dynamic badge updates: `[Residence]`, `[MI Al-Hamid]`, `[MTsN 30]`.
  - Shared Malang marker (University, Current Base) with dynamic badge updates: `[University]`, `[Current Base]`.
  - Exactly 3 geographic route segments: Karanganyar $\rightarrow$ Jakarta $\rightarrow$ Salatiga $\rightarrow$ Malang; zero `Jakarta -> Jakarta` and zero `Malang -> Malang` route hops.
  - Initial active state defaults strictly to Origin — Karanganyar (`origin-karanganyar`).
  - Storytelling milestone card gracefully omits the period row for MI and MTs without placeholder strings ("Unknown", "TBD", "TODO").
  - Responsive 7-step timeline (Desktop: horizontal progress bar; Mobile: vertical stacked compact timeline, zero horizontal document overflow).
  - Recruiter-facing copy updated: *"From Karanganyar to Jakarta, Salatiga, and Malang — a journey through formative education and computer science."* (0 instances of "founder").
  - Full keyboard accessibility (Arrow Left/Right, Tab, Enter) and screen reader labels.
  - Reduced-motion support (`prefers-reduced-motion: reduce`).
- [x] **Testing & Quality Gates**:
  - Updated `apps/web/tests/journey-map.test.tsx` covering all 22 mandatory criteria.
  - All 39 tests pass across test suite (`npm run test --workspace=apps/web`).
  - Lint (`npm run lint`), typecheck (`npm run typecheck`), data validation (`npm run validate:data`), and production build (`npm run build`) all pass cleanly.
  - First Load JS impact remains lean at **171 kB** (Route `/` size: 64 kB).
- [x] **Phase 3 Verification & Stop Rule**:
  - Captured 7 desktop (1440px) and 3 mobile (390px) screenshots to `docs/screenshots/phase3/`.
  - Updated design documentation: `docs/design/journey-map.md`.
  - Phase 3 RECERTIFIED. Awaiting explicit user approval before starting Phase 4 (Go Backend).

---

## Phase 4 Checklist — Go Backend (CERTIFIED)
> **Certified Scope**: Production-ready Go Clean Architecture REST API (`apps/api`), Go 1.22+ runtime compatibility (`go 1.22` in `go.mod`), standard library `net/http.ServeMux` routing, pure Go standard library domain entities, public transport DTO isolation boundary, in-memory indexed canonical JSON repository with deterministic path resolution and zero directory-walking, `jackc/pgx/v5` v5.6.0 connection pool client, explicit `DATABASE_MODE=disabled|optional|required`, readiness semantics, hardened middleware pipeline (Recovery, Logger, SecurityHeaders, CORS with safe Vary header, Bounded RateLimiter with IP port stripping), query parameter strictness, safe source paths, safe URL schemes, security headers on all error responses, and 100% passing concurrency-safe race test suite.

- [x] **Pure Go Domain Layer (`apps/api/internal/domain/`)**:
  - Zero external dependencies: `context`, `errors`, `time` only.
  - Canonical entities: `Profile`, `Project`, `Skill`, `Evidence`, `JourneyStop`.
  - Defined storage contracts: `ProfileRepository`, `ProjectRepository`, `SkillRepository`, `EvidenceRepository`, `JourneyRepository`.
  - Maintained canonical Journey categories (`birthplace`, `residence`, `tk`, `sd`, `smp`, `sma`, `university`, `current`).
- [x] **Public Transport DTO Boundary (`apps/api/internal/delivery/http/dto/`)**:
  - Strict isolation: domain entities are never serialized directly to HTTP.
  - Safe DTOs: `ProfileResponse`, `ProjectResponse`, `SkillResponse`, `EvidenceResponse`, `JourneyStopResponse`.
  - Standard response envelopes: `DataEnvelope[T]` and `ListEnvelope[T]` with `Meta.Count`.
  - Standard error envelope: `ErrorEnvelope` with certified `code` set (`bad_request`, `not_found`, `method_not_allowed`, `rate_limited`, `internal_error`, `service_unavailable`) and `message`.
  - **Privacy Guard & Minimization**:
    - `ProfileResponse`: Strips all internal `TODO` arrays, private contact fields, unapproved positioning, unpublished bio notes. Serializes strictly approved fields: `full_name`, `role`, `project_name`, `ai_feature`, `github`.
    - `JourneyStopResponse`: Strictly strips raw `coordinates`, `todo`, private records (`journey-tk`), redundant `public` boolean flag, and unapproved narrative descriptions. Retains purely factual structured fields (`id`, `category`, `title`, `institution`, `city`, `region`, `country`, `period`).
    - `EvidenceResponse`: Enforces repository-relative paths for `source_path` (e.g. `README.md`, `backend/cmd/api/main.go`). Rejects absolute paths (`/Users/...`, `C:\...`) and path traversal (`..`).
    - `URL Output Safety`: Enforces `https://` on all external URLs. Reject `javascript:`, `file:`, `data:`.
- [x] **Canonical JSON Repository & Deterministic Path (`apps/api/internal/repository/jsonfile/`)**:
  - Deterministic data path resolution (`PORTFOLIO_DATA_DIR` or fixed documented defaults: `../../data`, `../../../data`, `data`).
  - Completely removed recursive directory-walking data discovery.
  - Load once at startup, build immutable thread-safe in-memory indexes (slug, ID, category, public).
  - Fail-fast validation of invariants (duplicate project slug/ID, duplicate skill ID, duplicate evidence ID, invalid category, missing required files).
  - Zero disk I/O on active HTTP requests.
- [x] **PostgreSQL Client & Database Modes (`apps/api/internal/repository/postgres/`)**:
  - Connection pooling using `jackc/pgx/v5/pgxpool` v5.6.0 (Go 1.22 compatible).
  - Explicit operating mode: `DATABASE_MODE=disabled|optional|required`.
  - Readiness semantics:
    - `disabled`: HTTP 200 `{"status":"ready","database":"disabled"}`
    - `optional + connected`: HTTP 200 `{"status":"ready","database":"connected"}`
    - `optional + degraded`: HTTP 200 `{"status":"ready","database":"degraded"}` (remains degraded until service restart)
    - `required + unavailable`: Fails startup cleanly. In test/mock harness: HTTP 503 `{"status":"not_ready","database":"unavailable"}`
- [x] **Hardened HTTP Delivery & Middleware (`apps/api/internal/delivery/http/`)**:
  - Standard Go 1.22+ `net/http.ServeMux` with pattern matching (`GET /api/v1/projects/{slug}`).
  - Unsupported methods return `405 Method Not Allowed` with `Allow: GET, HEAD`, `Cache-Control: no-store`, full security headers, and the standard JSON error envelope (`{"error":{"code":"method_not_allowed","message":"method not allowed"}}`).
  - Content-Type: `application/json; charset=utf-8` on all JSON responses and errors.
  - Server timeouts configured: ReadHeaderTimeout 5s, ReadTimeout 10s, WriteTimeout 10s, IdleTimeout 60s, MaxHeaderBytes 1MB.
  - Query parameter strictness: ambiguous repeated keys (`?featured=true&featured=false`, `?category=AI&category=Backend`) return HTTP 400 `bad_request`. Invalid boolean values return HTTP 400.
  - Cache policy: `Cache-Control: public, max-age=60, stale-while-revalidate=300` on public portfolio GETs; `no-store` on probes and errors (400, 404, 405, 429, 500).
  - Middleware order: `Recovery` $\rightarrow$ `Logger` $\rightarrow$ `SecurityHeaders` $\rightarrow$ `CORS` $\rightarrow$ `RateLimiter` $\rightarrow$ `MethodNotAllowed` $\rightarrow$ `mux`.
  - Security headers present on ALL responses (200, 204, 400, 404, 405, 429, 500).
  - `RateLimiter`: Normalized client key via `net.SplitHostPort` (stripping ephemeral TCP source ports), bounded memory (`maxEntries = 10,000`), cleanup lifecycle `Close()`, bypass on `/healthz` and `/readyz`.
- [x] **Testing & Quality Gates**:
  - `apps/api/tests/health_test.go`: Probes across all database modes.
  - `apps/api/tests/repository_test.go`: Canonical JSON loading, invariant enforcement, duplicate slug rejection, filtering.
  - `apps/api/tests/endpoints_test.go`: All endpoints, filtering, 400/404/405 handling, CORS, security headers on errors, rate limiting port normalization, panic recovery, privacy leakage audit (0 `TODO_` tokens), exact public profile snapshot verification.
  - `go test -v -race ./...`: 100% PASS (0 data races).
  - `go vet ./...`: 100% PASS (0 warnings).
  - `npm run validate:data`: 100% PASS (8/8 schemas).
  - `npm run test --workspace=apps/web`: 100% PASS (39/39 tests).
  - `docker compose config`: 100% PASS.
  - Live server probe audit verifying all endpoints, headers, and status codes (200, 204, 400, 404, 405, 429, 500).
- [x] **Documentation**:
  - Updated `docs/architecture/backend-clean-architecture.md`.
  - Updated `docs/api/http-api.md`.
- [x] **Stop Rule Enforcement**: Halt execution upon completion; await explicit user approval before Phase 5 (AI Indexing & Retrieval).

---

## Phase 5 Checklist — AI Indexing & Retrieval Foundation (CERTIFIED)
> **Certified Scope**: Zero-hallucination evidence acquisition, normalization, chunking, and internal retrieval system. Completely removed deprecated `text-embedding-004`; provider-agnostic embedding abstraction with explicit `Purpose` (`document` vs `query`), verified dimension contract (768 dimensions), and Gemini default model `gemini-embedding-2`. Atomic repository snapshot replacement in PostgreSQL transactions, exact nearest-neighbor search (`<=>` cosine distance) without premature HNSW/ANN, strict model identity filtering, defense-in-depth secret and binary exclusion, read-only network `--dry-run` and offline `--plan-only` modes operating without PostgreSQL, canonical portfolio knowledge adapter with explicit provenance and relational link preservation (`project_id`, `skill_ids`, `evidence_id`), immutable citation generation, and zero public retrieval routes or generative chat features.

- [x] **Database Migrations (`apps/api/migrations/`)**:
  - `000002_knowledge_index.up.sql`: Created `knowledge_sources` and `knowledge_chunks` with `embedding vector NULL`, `embedding_provider`, `embedding_model`, `embedding_dimensions`, `project_id`, `skill_ids`, `evidence_id`, and exact model filter indexes. Zero HNSW/IVFFlat indexes created (exact cosine search first).
  - `000002_knowledge_index.down.sql`: Cascading clean rollback.
- [x] **Pure Go Domain Layer (`apps/api/internal/domain/knowledge.go`)**:
  - Standard library only (`context`, `time`).
  - Entities: `KnowledgeSource`, `KnowledgeChunk`, `SourceWithChunks`, `RetrievalResult`, `AICitation`, `RepositoryIndexStatus`.
  - Interface: `KnowledgeRepository` with `ReplaceRepositorySnapshot`, `SearchSimilar`, `DeleteRepositorySources`, `GetIndexStatus`.
- [x] **Canonical Allowlist & Manifests (`apps/api/internal/indexing/`)**:
  - `allowlist.go`: Loads canonical `data/ai/allowlist.json`. Authoritative and strictly read-only (zero runtime mutation).
  - `manifest.go`: Defined verified source manifests for all 5 allowlisted repositories (`Nachsyas/EduTrace`, `gdgoc-ecommerce`, `maritime-ai-dashboard`, `smart-kitchen-backend`, `smart-kitchen-frontend`).
  - Defense-in-depth secret exclusion (`.env`, `*.key`, `*.pem`, lockfiles, node_modules, etc.) and binary rejection (`MaxSourceFileBytes = 256 KiB`, NUL bytes, UTF-8 validity).
- [x] **Text Normalization & Chunking (`apps/api/internal/indexing/`)**:
  - `normalizer.go`: Global line ending normalization; trims trailing whitespace and collapses excessive blank lines outside code blocks; strictly preserves fenced code blocks verbatim.
  - `chunker.go`: Heading-aware Markdown chunking and bounded window code chunking; prepends deterministic context header (`Repository`, `Path`, `Section`); generates stable SHA-256 chunk IDs.
  - `checksum.go`: SHA-256 hex utilities and deterministic ID generators.
- [x] **Canonical Knowledge Source Adapter (`apps/api/internal/indexing/canonical_source.go`)**:
  - Adapts `profile`, `projects`, `skills`, `evidence`, and public `journey` stops.
  - Strict privacy minimization: zero coordinates, zero TODOs, zero private contact details.
  - Explicit provenance (`canonical_profile`, `canonical_project`, `canonical_skill`, `canonical_evidence`, `canonical_journey`).
  - Preserves relational identifiers (`project_id`, `skill_ids`, `evidence_id`).
- [x] **Hardened Standard Library GitHub Client (`apps/api/internal/github/`)**:
  - Standard library `net/http` client with 15s timeout.
  - Base URL pinned to `https://api.github.com` in production; test constructor accepts `httptest.Server`.
  - Redirect safety: allowlist check against GitHub hosts (`api.github.com`, `raw.githubusercontent.com`, `github.com`).
  - Commit SHA validation: strictly enforces 40 hex characters.
  - Bounded responses (256 KiB), binary detection, secret exclusion, 403/429 rate limit detection.
  - Immutable citation URL generator: `https://github.com/<owner>/<repo>/blob/<sha>/<path>`.
- [x] **Embedding Abstraction & Gemini Provider (`apps/api/internal/embedding/`)**:
  - `provider.go`: Explicit `Purpose` contract (`PurposeDocument`, `PurposeQuery`) and `ValidateVector` (768 dimensions, rejects NaN/Inf).
  - `disabled.go`: `DisabledProvider` for `EMBEDDING_MODE=disabled`.
  - `fake.go`: `DeterministicFakeProvider` generating unit-normalized vectors for CI structural testing.
  - `gemini/client.go`: Standard library HTTP client for `gemini-embedding-2` with explicit 768 dimensions, Purpose translation (`RETRIEVAL_DOCUMENT` vs `RETRIEVAL_QUERY`), $N \rightarrow N$ batch validation, and API key header authentication (never logged).
- [x] **PostgreSQL pgvector Repository (`apps/api/internal/repository/postgres/knowledge_repo.go`)**:
  - `ReplaceRepositorySnapshot`: Atomic single-transaction snapshot replacement (delete superseded state + insert new sources & chunks). Rollback on any failure preserves previous healthy state.
  - `SearchSimilar`: Exact nearest-neighbor cosine distance (`<=>`) with strict model identity filtering (`embedding_provider = $2 AND embedding_model = $3 AND embedding_dimensions = $4`).
  - Score semantics: `Distance` (`<=>`) and `Similarity` (`1.0 - Distance`). Bounded limit (default 5, max 20).
  - `DeleteRepositorySources`: Purges stored index state for a repository even if revoked from the allowlist.
  - `GetIndexStatus`: Distinguishes `vector_ready` vs `metadata_only`.
- [x] **Indexer CLI & Services (`apps/api/cmd/indexer/`, `internal/indexing/`, `internal/retrieval/`)**:
  - Indexer CLI supporting `--all-approved`, `--repo`, `--dry-run`, `--plan-only`, `--status`, `--delete-repo`.
  - Mutually incompatible flag validations and non-zero exit codes on failure.
  - `--plan-only`: Fully offline mode (0 network, 0 DB).
  - `--dry-run`: Read-only network mode (0 embedding calls, 0 DB writes, does not require PostgreSQL).
  - Internal retrieval service (`internal/retrieval/service.go`) with query embedding and model identity filtering.
  - Strictly zero public retrieval routes (`/api/v1/search`, `/api/v1/ai`) and zero generative LLM operations.
- [x] **Testing & Validation Gates**:
  - `apps/api/internal/indexing/normalizer_test.go`: Verbatim fenced code preservation, blank line collapsing.
  - `apps/api/internal/indexing/chunker_test.go`: Markdown heading segmentation, code windowing, context prepending.
  - `apps/api/internal/indexing/canonical_source_test.go`: Public canonical extraction, zero coordinates, zero TODOs.
  - `apps/api/internal/indexing/service_test.go`: Repository snapshot atomicity, idempotency, transaction failure rollback, atomic commit replacement, deletion.
  - `apps/api/internal/github/client_test.go`: 11 test cases covering SHA resolution, 404, rate limit, timeout, oversized, malformed SHA, path traversal, secret exclusion, binary rejection, unauthorized redirect.
  - `apps/api/internal/embedding/provider_test.go`: Vector dimension validation (768, NaN, Inf, short, long) and fake provider determinism.
  - `apps/api/internal/embedding/gemini/client_test.go`: Batch $N \rightarrow N$ semantics, purpose translation, dimensions.
  - `apps/api/internal/retrieval/service_test.go`: Model identity filtering, purpose query invocation, empty query rejection.
  - `go test -count=1 -race ./...`: 100% PASS (0 data races).
  - `go vet ./...`: 100% PASS (0 warnings).
  - `npm run validate:data`: 100% PASS (8/8 schemas).
  - `npm run test --workspace=apps/web`: 100% PASS (39/39 tests).
  - `docker compose config`: 100% PASS.
  - Live dry-run verified: 5 repositories, 17 files selected, 85 chunks formed, 40,694 bytes, 0 DB mutations, 0 embedding calls.
- [x] **Documentation**:
  - `docs/ai/indexing-pipeline.md`: Architecture, untrusted external boundary, canonical adapter, dry-run, atomic snapshots.
  - `docs/ai/retrieval-architecture.md`: Model identity vector space segregation, explicit purpose contract, 768 dimensions, exact search first, citations, Phase 6 boundary.
  - `docs/ai/retrieval-evals.md`: Groundedness benchmark, positive/negative queries, honest reporting separation (structural vs real).
- [x] **Stop Rule Enforcement**: Halt execution upon Phase 5 certification; await explicit user approval before Phase 6 (Ask Arham AI).

---

## Phase 6 Checklist — Ask Arham AI Grounded Reviewer Copilot (CERTIFIED)
> **Certified Scope**: Reviewer-oriented copilot grounded exclusively in Phase 5 approved evidence. Gemini Interactions API (`POST /v1beta/interactions`) with fallback to `generateContent` and configurable model `gemini-3.8-flash`. Structured generation schema (`status`, `segments` with `evidence_ids`, `suggested_action_ids`). Claim-level evidence validation (dropping unverified/unknown IDs). No raw model token streaming before server validation. Preflight failure checks returning normal HTTP error envelopes (400, 405, 429, 503) before SSE headers are committed. SSE lifecycle (`status: retrieving` -> `evidence` -> `status: generating` -> `result` -> `done`) with structured in-stream error handling. Server-owned immutable GitHub HTTPS and canonical portfolio citations (zero local file paths). Safe action registry matching verified routes. Strict removal of similarity scores from recruiter-facing UI. Mobile-first accessible panel (`100dvh`, `safe-area-inset-bottom`, polite aria-live announcer, Alt+A shortcut). 100% test passes across Go and TypeScript with zero data races.

- [x] **Server Configuration & Safe Action Registry (`apps/api/internal/config/`, `internal/ai/`)**:
  - `AIMode`, `AIProvider`, `AIModel`, `AIThinkingLevel`, `AIRequestTimeoutSeconds` (30), `AIMaxConcurrentRequests` (4), `AIRateLimitPerMinute` (5), `AIMaxEvidenceChars` (24000).
  - Safe Action Registry with 8 verified actions mapping to real internal targets (`/projects/edutrace`, `/projects/gdgoc-ecommerce`, `/projects/maritime-ai-dashboard`, `/projects/smart-kitchen`, `/#projects`, `/#skills`, `/#journey`, `/#contact`).
- [x] **LLM Provider Abstraction & Gemini Interactions Client (`apps/api/internal/llm/`)**:
  - `provider.go`: Domain contracts (`EvidenceContext`, `GroundedSegment`, `GeneratedAnswer`, `GenerateRequest`, `SourceCitation`, `PublicEvidenceItem`, `GroundedResponse`, `SafeAction`, `LLMProvider`).
  - `disabled.go`: `DisabledProvider` returning `ErrAIFeatureDisabled`.
  - `fake.go`: `DeterministicFakeProvider` for automated CI and local harnesses.
  - `gemini/client.go`: Standard library HTTP client targeting Google Interactions API with fallback to `generateContent`, JSON schema parsing, header authentication (`x-goog-api-key`). Zero Gemini tools enabled.
- [x] **Grounded Reviewer UseCase (`apps/api/internal/usecase/ask_usecase.go`)**:
  - Privacy Pre-Guard: Deterministic refusal for requests targeting private personal information (phone number, home address, private coordinates).
  - Prompt Injection Defense: System prompt boundary ensuring evidence blocks are inert reference data. Rejects prompt leak attempts.
  - Zero-Evidence Short-Circuit & Top-K Context Budgeting (`AI_MAX_EVIDENCE_CHARS = 24000`).
  - Server-side Claim Validation: Every segment's `evidence_ids` verified against retrieved evidence; unknown IDs dropped.
  - Citation Mapping: Builds immutable GitHub HTTPS links or safe canonical labels without local file paths.
  - Safe Action Resolution: Validates against certified action registry.
- [x] **Hardened SSE Delivery & HTTP Handler (`apps/api/internal/delivery/http/ai_handler.go`)**:
  - Strict preflight validation before committing SSE headers (400 bad JSON / length 2..1000, 405 non-POST, 503 AI disabled or concurrency permit exhausted, 429 rate limit).
  - Concurrency Semaphore Permit: 4 permits, acquired before SSE headers, released on success/failure/timeout/client abort.
  - Rate Limiter: 5 req/min per normalized client identity.
  - Flusher implementation on ResponseWriter wrappers in middleware.
  - SSE headers committed only on success (`text/event-stream`, `Cache-Control: no-store`, `X-Content-Type-Options: nosniff`, `X-Request-ID`).
  - SSE Event Lifecycle: `status: retrieving` -> `evidence` -> `status: generating` -> `result` -> `done`.
  - In-stream error handling: Emits safe JSON `event: error` -> `event: done`.
- [x] **Frontend Reviewer Copilot UI (`apps/web/features/ask-arham/`)**:
  - `AskArhamLauncher`: Floating button with ping animation, accessible label, and Alt+A shortcut.
  - `AskArhamPanel`: Responsive drawer / mobile bottom-sheet (`100dvh`, `safe-area-inset-bottom`), Escape key handler, focus management, polite `aria-live` announcer.
  - `SuggestedQuestions`: 7 recruiter-oriented questions.
  - `EvidenceList`: Collapsible accordion displaying safe excerpts; strictly ZERO similarity scores.
  - `SourceList`: Immutable GitHub links and certified canonical portfolio labels without local paths.
  - `SafeActionButtons`: Verified action buttons with client-side safe routing.
  - `AskArhamAnswer`: Renders status badges, claim segments with `[E1]` badges, sources, and actions.
  - `AskArhamInput`: Character counter (2..1000), submit button, stream abort button.
  - `ask-arham.api.ts`: Fetch ReadableStream SSE consumer handling fragmented chunks, CRLF/LF line endings, and AbortController.
  - Dynamic import in `PortfolioApp.tsx` (`ssr: false`) keeping initial load overhead minimal.
- [x] **Testing & Validation Gates**:
  - `apps/api/internal/llm/provider_test.go`: PASS
  - `apps/api/internal/llm/gemini/client_test.go`: PASS
  - `apps/api/internal/usecase/ask_usecase_test.go`: PASS
  - `apps/api/internal/delivery/http/ai_handler_test.go`: PASS
  - `apps/web/tests/ask-arham-sse.test.ts`: PASS (11 tests)
  - `apps/web/tests/ask-arham-ui.test.tsx`: PASS (10 tests)
  - `go test -count=1 -race ./...`: 100% PASS (0 data races)
  - `go vet ./...`: 100% PASS (0 warnings)
  - `npm run validate:data`: 100% PASS (8/8 schemas)
  - `npm run lint --workspace=apps/web`: 100% PASS
  - `npm run typecheck --workspace=apps/web`: 100% PASS
  - `npm run test --workspace=apps/web`: 100% PASS (7 test files, 60 tests)
  - `npm run build --workspace=apps/web`: 100% PASS
  - `docker compose config`: 100% PASS
- [x] **Screenshots & Visual Evidence (`docs/screenshots/phase6/`)**:
  - `01_desktop_launcher.png`: Floating launcher button on desktop.
  - `02_desktop_empty_panel.png`: Slide-over panel with suggested recruiter questions.
  - `03_desktop_grounded_response.png`: Grounded answer with verified evidence badge and claim citations.
  - `04_desktop_evidence_sources.png`: Expanded evidence excerpts (zero similarity scores) and primary sources.
  - `05_desktop_refusal.png`: Refusal state for ungrounded question (insufficient evidence guard).
  - `06_desktop_ai_unavailable.png`: Graceful banner when AI is disabled (HTTP 503).
  - `07_mobile_panel.png`: Mobile panel sheet at 390px viewport.
  - `08_mobile_grounded_response.png`: Mobile grounded answer and claim citations.
  - `09_mobile_sources.png`: Mobile certified sources and contextual actions.
  - `10_mobile_input.png`: Mobile input field and submit button.
- [x] **Documentation**:
  - `docs/ai/ask-arham-architecture.md`: Interactions API, structured output, SSE lifecycle, security boundaries.
  - `docs/ai/answer-evals.md`: Evaluation methodology, 8 canonical questions, honest reporting separation.
  - `docs/api/http-api.md`: `POST /api/v1/ai/ask` specification, preflight error contracts, in-stream recovery.
  - `.env.example`: Non-secret placeholders for all Phase 6 environment variables.
- [x] **Stop Rule Enforcement**: Halt execution upon Phase 6 certification and push; strictly do NOT proceed to Phase 7 without explicit user approval.

---

## Phase 8 Checklist — Production & Launch (CORE PORTFOLIO LIVE — ASK ARHAM AI PENDING)
> **Production Architecture**:
> - **Frontend**: Vercel (`https://arham-porto.vercel.app`) running Next.js 15 App Router (`apps/web`).
> - **Backend API**: Google Cloud Run (`https://arham-porto-api-6ivuekmjva-et.a.run.app`, region `asia-southeast2`).
> - **Database**: Supabase PostgreSQL 16 with `pgvector` extension.
> - **Liveness Hardening**: `GET /health` returns 200 OK (`{"status":"ok","scope":"process_alive"}`). Cloud Run edge anomaly affecting `/healthz` documented without claiming Google officially reserves it.
> - **Readiness**: `GET /readyz` returns 200 OK (`{"status":"ready","database":"connected"}`).
> - **CORS**: Finalized to exact Vercel production origin (`https://arham-porto.vercel.app`), unauthorized origins denied.
> - **Indexer Job**: Cloud Run Job `arham-porto-indexer` created targeting `gemini-embedding-2` (768 dimensions).
> - **AI Copilot Status**: `PENDING PRODUCTION ACTIVATION` — blocked on Google AI Studio prepayment credits / quota (HTTP 429 RESOURCE_EXHAUSTED). Core portfolio is 100% operational; `/api/v1/ai/ask` fails gracefully with HTTP 503 Service Unavailable.
> - **Phase 7**: DEFERRED.

