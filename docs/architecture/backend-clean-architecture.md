# Architecture: Backend Clean Architecture

- **Purpose**: Specify the clean architecture layers, repository strategies, and delivery contracts for the Go backend service (`apps/api`).
- **Current Status**: Phase 4 Certified (Backend Contract & Security Hardened).
- **Minimum Go Runtime**: Go 1.22+ (`go 1.22` in `go.mod`).

---

## 1. Clean Architecture Layers & Data Flow

```
HTTP Request
    │
    ▼
[ Delivery Layer ] (apps/api/internal/delivery/http)
  ├── Hardened Server Configuration
  │     ├── ReadHeaderTimeout: 5s
  │     ├── ReadTimeout: 10s
  │     ├── WriteTimeout: 10s
  │     ├── IdleTimeout: 60s
  │     └── MaxHeaderBytes: 1MB (1 << 20)
  ├── Middleware Pipeline (Recovery -> Logger -> SecurityHeaders -> CORS -> RateLimiter)
  ├── net/http.ServeMux Router (Standard Library Go 1.22+ Path & Method Routing, Zero Frameworks)
  └── HTTP Handlers (Single-value query parameter validation, Path value regex, DTO mapping)
    │
    ▼ (invokes)
[ Public DTO Layer ] (apps/api/internal/delivery/http/dto)
  ├── Strict Isolation: Domain Entity != Public API Response
  ├── ProfileResponse: Excludes unapproved/TODO-backed fields (positioning, bio, city, availability)
  ├── ProjectResponse: Sanitizes URLs (https:// required), strips internal TODOs
  ├── SkillResponse: Pure verified claims and associated evidence IDs
  ├── EvidenceResponse: Sanitizes SourceURL (https://) and SourcePath (repository-relative only)
  ├── JourneyStopResponse: Excludes redundant "public" boolean, coordinates, and private records
  ├── Uniform Envelopes: DataEnvelope[T] & ListEnvelope[T] with Meta
  └── Standard ErrorEnvelope: code (bad_request, not_found, rate_limited, internal_error, service_unavailable)
    ▲
    │ (maps from)
[ UseCase Layer ] (apps/api/internal/usecase)
  ├── ProfileUseCase: GetProfile
  ├── ProjectUseCase: ListProjects (category & featured filters), GetProjectBySlug
  ├── SkillUseCase: ListSkills (category filter)
  ├── EvidenceUseCase: ListEvidence (skill_id & project_id filters), GetEvidenceByID
  └── JourneyUseCase: ListPublicStops (strictly enforces public == true, zero private leakage)
    │
    ▼ (interacts through domain interfaces)
[ Domain Layer ] (apps/api/internal/domain)
  ├── Pure Go standard library ONLY (context, errors, time)
  ├── Zero third-party packages (no pgx, no net/http, no AI SDKs, no frameworks)
  ├── Entities: Profile, Project, Skill, Evidence, JourneyStop, KnowledgeChunk
  └── Repository Interfaces: ProfileRepository, ProjectRepository, SkillRepository, EvidenceRepository, JourneyRepository
    ▲
    │ (implemented by)
[ Infrastructure / Repository Layer ] (apps/api/internal/repository)
  ├── jsonfile.Repository (Canonical in-memory store loaded once at startup with fail-fast invariant checks)
  └── postgres.Client (Connection pool wrapper via jackc/pgx/v5/pgxpool v5.6.0 with explicit lifecycle & ping)
```

---

## 2. Purity of the Domain Layer

- `apps/api/internal/domain/` contains **ZERO** external dependencies.
- No `pgx`, no `net/http`, no frameworks, no database driver packages.
- Repository contracts accept standard `context.Context` and return pure Go domain entities or standard library `error`.

---

## 3. Public Transport DTO Boundary & Sanitizers

Domain entities represent internal canonical structures and may contain internal metadata (such as `TODO` arrays, private geographic coordinates, internal verification metadata).
- **Rule**: Domain entities must **NEVER** be serialized directly to HTTP responses.
- All HTTP responses pass through public DTO mappers in `internal/delivery/http/dto/`:
  - `ProfileResponse`: Strips all internal TODOs, unapproved positioning, unpublished bio fields, and private contact notes. Only verified approved identity fields populate (`full_name`, `role`, `project_name`, `ai_feature`, `github`).
  - `JourneyStopResponse`: Strips raw coordinates, internal TODOs, and the redundant publication-control flag (`public`). Returns public resource attributes only.
  - `EvidenceResponse`: Verifies that `source_path` can contain **ONLY** repository-relative paths (e.g. `README.md`, `backend/cmd/api/main.go`). Any absolute paths (`/Users/...`, `C:\...`) or traversal attempts are strictly rejected and stripped to `nil`.
  - `URL Sanitization`: All external URLs (`github_url`, `demo_url`, `source_url`, `image`) must strictly use `https://` (or `/` for relative images). Schemes like `file://`, `javascript:`, `data:` are rejected and return `nil`.

---

## 4. Repository Architecture & Deterministic Data Resolution

Canonical portfolio data lives in `data/` and is validated via TypeScript Zod (`packages/content-schema`) at authoring time.
- **Deterministic Resolution**:
  - If `PORTFOLIO_DATA_DIR` is set: strictly uses that directory.
  - If unset: uses documented deterministic relative defaults (`../../data` when running from `apps/api`, `../../../data` when running from `apps/api/tests`, or `data` when running from repository root).
  - **No Directory Walking**: Does NOT recursively walk parent directories.
  - **Startup Fail-Fast**: Validates all canonical files (`profile/profile.json`, `projects/projects.json`, `skills/skills.json`, `evidence/evidence.json`, `journey/journey.json`). If missing or invalid, startup aborts immediately.
- **Zero Disk I/O During Serving**: Content is indexed in memory at startup. Reads are concurrent and thread-safe via `sync.RWMutex`.

---

## 5. Database Modes & Readiness Semantics

PostgreSQL connectivity is managed explicitly via `DATABASE_MODE`:
- `disabled`: PostgreSQL connection is skipped. `/readyz` returns HTTP 200 `{"status":"ready","database":"disabled"}`.
- `optional`: Attempts connection with a 3-second timeout.
  - If online: `/readyz` returns HTTP 200 `{"status":"ready","database":"connected"}`.
  - If offline: logs a warning, continues serving from canonical in-memory repository, and `/readyz` returns HTTP 200 `{"status":"ready","database":"degraded"}`.
  - **Recovery Semantics**: The service remains degraded until restart (no background reconnection polling in Phase 4).
- `required`: Mandatory PostgreSQL. If offline, the process fails startup cleanly and terminates immediately. In mock/test diagnostic environments, `/readyz` returns HTTP 503 `{"status":"not_ready","database":"unavailable"}`.
- **Liveness Probe**:
  - `/healthz`: Always returns HTTP 200 `{"status":"ok","scope":"process_alive"}` with `Cache-Control: no-store`.

---

## 6. Rate Limiter Security & Trust Model

- **Client Key Normalization**: Client identity is normalized using `net.SplitHostPort(r.RemoteAddr)` to strip ephemeral TCP source ports. Requests from `192.0.2.10:51432` and `192.0.2.10:51433` share the exact same rate-limit bucket.
- **Zero-Trust Header Model**: Does NOT blindly trust `X-Forwarded-For` or `X-Real-IP` until a verified trusted reverse proxy configuration is provisioned.
- **Memory Boundedness**: Maximum client entries are capped (`defaultMaxRateLimiterEntries = 10,000`). If client cardinality reaches capacity, expired buckets are pruned immediately. If capacity remains saturated, new allocations are rejected to guarantee finite memory usage.
- **Lifecycle Cleanliness**: Background cleanup goroutine runs on a ticker and stops cleanly via `Close()` during graceful server shutdown. Tests terminate the limiter cleanly without goroutine leaks.
- **Probe Bypass**: `/healthz` and `/readyz` strictly bypass rate limiting.

---

## 7. HTTP Contracts & Query Strictness

- **Content-Type**: Every JSON response (including error envelopes) emits `Content-Type: application/json; charset=utf-8`.
- **Method Enforcement**: Unsupported HTTP methods on registered routes (e.g. `POST /api/v1/projects`) return `405 Method Not Allowed` with `Allow: GET, HEAD`.
- **HEAD Support**: Standard Go `http.ServeMux` matches HEAD requests against GET endpoints.
- **Query Parameter Strictness**:
  - Recognized parameters with ambiguous duplicate keys (e.g. `?featured=true&featured=false` or `?category=AI&category=Backend`) are rejected with `400 Bad Request`.
  - Invalid filter values (e.g. `?featured=banana`) are rejected with `400 Bad Request`.
  - Unknown query parameters are safely ignored.
- **Caching**:
  - Public portfolio GET endpoints emit: `Cache-Control: public, max-age=60, stale-while-revalidate=300`.
  - Errors (400, 404, 405, 429, 500) and probes (`/healthz`, `/readyz`) emit: `Cache-Control: no-store`.
- **Security Headers on All Responses**: All responses (200, 204, 400, 404, 405, 429, 500) carry:
  - `X-Content-Type-Options: nosniff`
  - `Referrer-Policy: no-referrer`
  - `Content-Security-Policy: default-src 'none'; frame-ancestors 'none'`
  - `X-Frame-Options: DENY`
- **CORS**:
  - Evaluates request origin against `ALLOWED_ORIGINS`.
  - Emits `Vary: Origin` cleanly without overwriting existing `Vary` headers.
  - No wildcard `*` in production mode unless explicitly configured.
