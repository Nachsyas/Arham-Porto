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
| **Phase 5** | **AI Indexing & Retrieval** | NOT STARTED |
| **Phase 6** | **AI Reviewer Copilot** | NOT STARTED |
| **Phase 7** | **Integration & Polish** | NOT STARTED |
| **Phase 8** | **Production & Launch** | NOT STARTED |

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
> **Certified Scope**: Production-ready Go Clean Architecture REST API (`apps/api`), Go 1.22+ standard library `net/http.ServeMux` routing, pure Go standard library domain entities, public transport DTO isolation boundary, in-memory indexed canonical JSON repository with fail-fast startup invariant checks, `jackc/pgx/v5` connection pool client, explicit `DATABASE_MODE=disabled|optional|required`, readiness semantics, hardened middleware pipeline (Recovery, Logger, SecurityHeaders, CORS, RateLimiter), and 100% passing concurrency-safe race test suite.

- [x] **Pure Go Domain Layer (`apps/api/internal/domain/`)**:
  - Zero external dependencies: `context`, `errors`, `time` only.
  - Canonical entities: `Profile`, `Project`, `Skill`, `Evidence`, `JourneyStop`.
  - Defined storage contracts: `ProfileRepository`, `ProjectRepository`, `SkillRepository`, `EvidenceRepository`, `JourneyRepository`.
  - Maintained canonical Journey categories (`birthplace`, `residence`, `tk`, `sd`, `smp`, `sma`, `university`, `current`).
- [x] **Public Transport DTO Boundary (`apps/api/internal/delivery/http/dto/`)**:
  - Strict isolation: domain entities are never serialized directly to HTTP.
  - Safe DTOs: `ProfileResponse`, `ProjectResponse`, `SkillResponse`, `EvidenceResponse`, `JourneyStopResponse`.
  - Standard response envelopes: `DataEnvelope[T]` and `ListEnvelope[T]` with `Meta.Count`.
  - Standard error envelope: `ErrorEnvelope` with `code` (`bad_request`, `not_found`, `rate_limited`, `internal_error`, `service_unavailable`) and `message`.
  - **Privacy Guard**:
    - `ProfileResponse`: Strips all internal `TODO` arrays, private contact fields, unpublished bio notes.
    - `JourneyStopResponse`: Strictly strips raw `coordinates`, `todo`, and private metadata; `journey-tk` excluded.
- [x] **Canonical JSON Repository (`apps/api/internal/repository/jsonfile/`)**:
  - Configurable deterministic data path (`PORTFOLIO_DATA_DIR` with auto-resolution).
  - Load once at startup, build immutable thread-safe in-memory indexes (slug, ID, category, public).
  - Fail-fast validation of invariants (duplicate project slug/ID, duplicate skill ID, duplicate evidence ID, invalid category).
  - Zero disk I/O on active HTTP requests.
- [x] **PostgreSQL Client & Database Modes (`apps/api/internal/repository/postgres/`)**:
  - Connection pooling using `jackc/pgx/v5/pgxpool`.
  - Explicit operating mode: `DATABASE_MODE=disabled|optional|required`.
  - Readiness semantics:
    - `disabled`: HTTP 200 `{"status":"ready","database":"disabled"}`
    - `optional + connected`: HTTP 200 `{"status":"ready","database":"connected"}`
    - `optional + degraded`: HTTP 200 `{"status":"ready","database":"degraded"}`
    - `required + unavailable`: HTTP 503 `{"status":"not_ready","database":"unavailable"}`
- [x] **Hardened HTTP Delivery & Middleware (`apps/api/internal/delivery/http/`)**:
  - Standard Go 1.22+ `net/http.ServeMux` with pattern matching (`GET /api/v1/projects/{slug}`).
  - Strict query & path validation: invalid category or boolean parameter returns HTTP 400 `bad_request`.
  - Cache policy: `Cache-Control: public, max-age=60, stale-while-revalidate=300` on public portfolio GETs; `no-store` on probes and errors.
  - Middleware order: `Recovery` $\rightarrow$ `Logger` $\rightarrow$ `SecurityHeaders` $\rightarrow$ `CORS` $\rightarrow$ `RateLimiter` $\rightarrow$ `Router`.
  - `Recovery`: Catches panics, logs server-side diagnostics, returns generic 500 JSON without stack trace leaks.
  - `SecurityHeaders`: `X-Content-Type-Options: nosniff`, `Referrer-Policy: no-referrer`, `Content-Security-Policy: default-src 'none'; frame-ancestors 'none'`, `X-Frame-Options: DENY` (no deprecated `X-XSS-Protection`).
  - `CORS`: Origin checking against `ALLOWED_ORIGINS`, `Vary: Origin`, OPTIONS 204.
  - `RateLimiter`: In-memory token bucket / sliding window with TTL cleanup; `/healthz` and `/readyz` strictly bypass throttling.
- [x] **Testing & Quality Gates**:
  - `apps/api/tests/health_test.go`: Probes across all database modes.
  - `apps/api/tests/repository_test.go`: Canonical JSON loading, invariant enforcement, duplicate slug rejection, filtering.
  - `apps/api/tests/endpoints_test.go`: All 9 endpoints, filtering, 400/404 handling, CORS, security headers, rate limiting, panic recovery, privacy leakage audit (0 `TODO_` tokens).
  - `go test -v -race ./...`: 100% PASS (0 data races).
  - `go vet ./...`: 100% PASS (0 warnings).
  - `npm run validate:data`: 100% PASS (8/8 schemas).
  - `npm run test --workspace=apps/web`: 100% PASS (39/39 tests).
  - `docker compose config`: 100% PASS.
  - Manual `curl` audit verifying all endpoints and status codes (200, 204, 400, 404, 429, 500).
- [x] **Documentation**:
  - Updated `docs/architecture/backend-clean-architecture.md`.
  - Created `docs/api/http-api.md`.
- [x] **Stop Rule Enforcement**: Halt execution upon completion; await explicit user approval before Phase 5 (AI Indexing & Retrieval).



