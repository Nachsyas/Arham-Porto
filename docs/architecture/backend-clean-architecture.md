# Architecture: Backend Clean Architecture

- **Purpose**: Specify the clean architecture layers, repository strategies, and delivery contracts for the Go backend service (`apps/api`).
- **Current Phase**: Phase 4 (Go Backend — Certified).

---

## 1. Clean Architecture Layers & Data Flow

```
HTTP Request
    │
    ▼
[ Delivery Layer ] (apps/api/internal/delivery/http)
  ├── Hardened Server (ReadHeaderTimeout, WriteTimeout, MaxHeaderBytes)
  ├── Middleware Pipeline (Recovery -> Logger -> SecurityHeaders -> CORS -> RateLimiter)
  ├── net/http.ServeMux Router (Standard Library Go 1.22+ Path & Method Routing)
  └── HTTP Handlers (Query & Path Sanitization, Error Mapping)
    │
    ▼ (invokes)
[ Public DTO Layer ] (apps/api/internal/delivery/http/dto)
  ├── Strict Isolation: Domain Entity != Public API Response
  ├── ProfileResponse, ProjectResponse, SkillResponse, EvidenceResponse, JourneyResponse
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
  ├── Entities: Profile, Project, Skill, Evidence, JourneyStop, KnowledgeChunk
  └── Repository Interfaces: ProfileRepository, ProjectRepository, SkillRepository, EvidenceRepository, JourneyRepository
    ▲
    │ (implemented by)
[ Infrastructure / Repository Layer ] (apps/api/internal/repository)
  ├── jsonfile.Repository (Canonical in-memory store loaded once at startup with fail-fast invariant checks)
  └── postgres.Client (Connection pool wrapper via jackc/pgx/v5/pgxpool with explicit lifecycle & ping)
```

---

## 2. Purity of the Domain Layer

- `apps/api/internal/domain/` contains **ZERO** external dependencies.
- No `pgx`, no `net/http`, no frameworks, no database driver packages.
- Repository contracts accept standard `context.Context` and return pure Go domain entities or standard library `error`.

---

## 3. Public Transport DTO Boundary

Domain entities represent internal canonical structures and may contain internal metadata (such as `TODO` arrays, private geographic coordinates, internal verification metadata).
- **Rule**: Domain entities must **NEVER** be serialized directly to HTTP responses.
- All HTTP responses are passed through public DTO mappers in `internal/delivery/http/dto/`.
- Privacy Guard:
  - `ProfileResponse`: Strips all internal TODOs, unpublished bio fields, and private contact notes.
  - `JourneyStopResponse`: Strips raw coordinates, internal TODOs, and private metadata.
  - `ProjectResponse`, `SkillResponse`, `EvidenceResponse`: Strips internal TODOs.

---

## 4. Repository Architecture & Startup Loading

Canonical portfolio data lives in `data/` and is validated via Zod (`packages/content-schema`).
- **Load Once**: `jsonfile.LoadRepository` reads, unmarshals, and builds immutable in-memory lookup indexes at startup.
- **Fail-Fast Invariants**: Duplicate slugs, duplicate IDs, or invalid categories cause startup to abort immediately.
- Zero disk I/O on active HTTP requests.
- Concurrency-safe: reads are protected by `sync.RWMutex`.

---

## 5. Database Modes & Readiness Contract

The API explicitly manages PostgreSQL connectivity via `DATABASE_MODE`:
- `disabled`: PostgreSQL connection is skipped entirely. `/readyz` returns HTTP 200 `{"status":"ready","database":"disabled"}`.
- `optional`: Attempts connection with short timeout. If connected: HTTP 200 `{"status":"ready","database":"connected"}`. If offline: logs warning, continues serving canonical JSON, and `/readyz` returns HTTP 200 `{"status":"ready","database":"degraded"}`.
- `required`: Mandatory PostgreSQL. If unavailable, startup or `/readyz` fails with HTTP 503 `{"status":"not_ready","database":"unavailable"}`.

Health Probe:
- `/healthz`: Always returns HTTP 200 `{"status":"ok","scope":"process_alive"}` with `Cache-Control: no-store`.

---

## 6. Middleware Composition Order

Requests flow outer $\rightarrow$ inner:
1. **Recovery**: Catches panics, logs diagnostics server-side, returns generic HTTP 500 JSON without leaking stack traces.
2. **Logger**: Privacy-safe structured logging (timestamp, method, path, status, duration, client IP). No headers or payloads logged.
3. **SecurityHeaders**: `X-Content-Type-Options: nosniff`, `Referrer-Policy: no-referrer`, `Content-Security-Policy: default-src 'none'; frame-ancestors 'none'`, `X-Frame-Options: DENY`.
4. **CORS**: Origin validation against allowlist (`ALLOWED_ORIGINS`), `Vary: Origin`, OPTIONS preflight 204.
5. **RateLimiter**: In-memory token bucket / sliding window with TTL cleanup. Throttles at configured threshold with HTTP 429 (`rate_limited`). `/healthz` and `/readyz` strictly bypass throttling.

