# Arham Porto — GitHub Repository Sync Service (Sprint 1)

## 1. Architecture Overview

The GitHub Repository Sync Service establishes the foundational ingestion pipeline for Arham Porto, transforming the portfolio from a static presentation into an automated living engineering portfolio grounded in verifiable GitHub activity.

### Architectural Hierarchy

```
GitHub REST API (users/Nachsyas/repos)
               │
               ▼ [HTTPS / Bearer Auth / Pagination]
┌──────────────────────────────────────────────┐
│        apps/api/internal/github/             │
│  - client.go   : Hardened GitHub REST client │
│  - models.go   : DTOs & Domain normalization │
│  - service.go  : Synchronization workflow    │
└──────────────────────┬───────────────────────┘
                       │
                       ▼ [Pure Go Domain Interface]
┌──────────────────────────────────────────────┐
│  domain.GithubProjectRepository (Upsert)     │
└──────────────────────┬───────────────────────┘
                       │
                       ▼ [pgxpool / SQL]
┌──────────────────────────────────────────────┐
│        apps/api/internal/repository/postgres │
│  - Table: github_projects (PostgreSQL)       │
└──────────────────────────────────────────────┘
                       ▲
                       │ [HTTP POST /api/v1/github/sync]
┌──────────────────────────────────────────────┐
│        apps/api/internal/delivery/http       │
│  - router.go & github_sync_handler.go        │
└──────────────────────────────────────────────┘
```

### Clean Architecture Boundaries
In strict alignment with `AGENTS.md` and repository standards:
1. **Domain Layer Isolation (`internal/domain/github_project.go`)**:
   - Strictly uses Go standard library (`context`, `time`). Zero third-party dependencies.
   - Defines the `GithubProject` entity and `GithubProjectRepository` interface.
2. **Integration Layer (`internal/github/`)**:
   - `client.go`: Hardened standard library HTTP client with 15s timeout, redirect protection, rate limit detection, and pagination support.
   - `models.go`: External GitHub API DTOs and clean mapping methods to standard structures.
   - `service.go`: Orchestrates sync, normalization, and metrics calculation (`Synced`, `Created`, `Updated`).
3. **Database Repository Layer (`internal/repository/postgres/github_project_repo.go`)**:
   - Implements `domain.GithubProjectRepository` using `jackc/pgx/v5/pgxpool`.
   - Executes atomic SQL upserts with `(xmax = 0)` conflict detection to accurately distinguish `INSERT` vs `UPDATE`.
   - Never exposes raw database connections to HTTP handlers.
4. **Delivery Layer (`internal/delivery/http/github_sync_handler.go`)**:
   - Maps HTTP `POST /api/v1/github/sync` through the security, logging, CORS, and rate limiting middleware stack.
   - Converts domain errors into standard JSON error envelopes.

---

## 2. Synchronization Flow

```
[Trigger] POST /api/v1/github/sync
    │
    ├── 1. Validate GITHUB_TOKEN configuration (fail with 401 if missing)
    │
    ├── 2. Fetch public owner repositories from GitHub API:
    │      GET https://api.github.com/users/Nachsyas/repos?per_page=100&page={n}&type=owner
    │      - Respects pagination (Link header rel="next")
    │      - Handles GitHub API rate limits (HTTP 429)
    │      - Context cancellation & 15s timeout
    │
    ├── 3. Normalize external repository schema:
    │      - Nil description/homepage/language converted to empty strings
    │      - Topics normalized to JSONB array
    │
    ├── 4. Upsert into PostgreSQL (github_projects table):
    │      - Match on github_id UNIQUE
    │      - New repository: INSERT (created count +1)
    │      - Existing repository: UPDATE fields & updated_at (updated count +1)
    │      - Preserves existing readme & preview_image columns
    │
    └── 5. Historical Record Retention:
           - Repositories deleted on GitHub are NEVER deleted from the database
           - Preserves full audit trail and historical evidence
```

---

## 3. Environment Variables & Secret Configuration

| Variable | Required | Description | Example / Secret Name |
| :--- | :--- | :--- | :--- |
| `GITHUB_TOKEN` | Yes | Personal access token with `public_repo` or `read:user` scope | `ghp_xxxxxxxxxxxxxxxxxxxx` |
| `DATABASE_URL` | Yes (in required mode) | PostgreSQL connection string | `postgres://user:pass@host:5432/arham_porto` |
| `DATABASE_MODE` | No | `disabled`, `optional`, or `required` (default: `optional`) | `optional` |

> [!IMPORTANT]
> **Production Secret Management**:
> In Google Cloud Run deployment, the token is provisioned via Secret Manager using secret name:
> `arham-porto-github-token` mapped to environment variable `GITHUB_TOKEN`.
> The token is never logged, printed, or exposed in HTTP responses or version control.

---

## 4. Database Design

### Migration Files
- `apps/api/migrations/000003_github_projects.up.sql`
- `apps/api/migrations/000003_github_projects.down.sql`
- `apps/api/migrations/000003_github_projects.sql`

### Schema Definition
```sql
CREATE TABLE IF NOT EXISTS github_projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    github_id BIGINT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    full_name TEXT NOT NULL,
    description TEXT,
    html_url TEXT,
    homepage TEXT,
    language TEXT,
    stars INTEGER DEFAULT 0,
    forks INTEGER DEFAULT 0,
    topics JSONB DEFAULT '[]',
    readme TEXT,
    preview_image TEXT,
    synced_at TIMESTAMP DEFAULT now(),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

-- Fast lookup indexes
CREATE INDEX IF NOT EXISTS idx_github_projects_github_id ON github_projects(github_id);
CREATE INDEX IF NOT EXISTS idx_github_projects_updated_at ON github_projects(updated_at);
```

### Conflict Resolution Strategy
```sql
INSERT INTO github_projects (
    github_id, name, full_name, description, html_url, homepage, language,
    stars, forks, topics, readme, preview_image, synced_at, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10::jsonb, $11, $12, NOW(), NOW(), NOW()
)
ON CONFLICT (github_id) DO UPDATE SET
    name = EXCLUDED.name,
    full_name = EXCLUDED.full_name,
    description = EXCLUDED.description,
    html_url = EXCLUDED.html_url,
    homepage = EXCLUDED.homepage,
    language = EXCLUDED.language,
    stars = EXCLUDED.stars,
    forks = EXCLUDED.forks,
    topics = EXCLUDED.topics,
    synced_at = NOW(),
    updated_at = NOW()
RETURNING id, (xmax = 0) AS is_created, created_at, updated_at, synced_at;
```

---

## 5. HTTP API Endpoint

### `POST /api/v1/github/sync`

Trigger repository ingestion and synchronization from GitHub to PostgreSQL.

#### Request
- **Method**: `POST`
- **Path**: `/api/v1/github/sync`
- **Headers**:
  - `Accept: application/json`

#### Success Response (`200 OK`)
```json
{
  "status": "success",
  "synced": 25,
  "created": 5,
  "updated": 20
}
```

#### Error Responses
- **`401 Unauthorized`** (missing or invalid `GITHUB_TOKEN`):
  ```json
  {
    "error": {
      "code": "unauthorized",
      "message": "GITHUB_TOKEN is missing or not configured"
    }
  }
  ```
- **`429 Too Many Requests`** (GitHub rate limit exceeded):
  ```json
  {
    "error": {
      "code": "rate_limited",
      "message": "GitHub API rate limit exceeded"
    }
  }
  ```
- **`502 Bad Gateway`** (GitHub upstream communication failure):
  ```json
  {
    "error": {
      "code": "bad_gateway",
      "message": "failed to fetch repositories from GitHub API"
    }
  }
  ```
- **`503 Service Unavailable`** (Database disconnected or sync service unconfigured):
  ```json
  {
    "error": {
      "code": "service_unavailable",
      "message": "database is unavailable"
    }
  }
  ```
- **`405 Method Not Allowed`** (Disallowed HTTP methods, e.g. GET/PUT/DELETE):
  ```json
  {
    "error": {
      "code": "method_not_allowed",
      "message": "method not allowed"
    }
  }
  ```

---

## 6. Future Expansion Roadmap

The GitHub Repository Sync Service serves as the trusted evidence layer for upcoming automated capabilities:

```
┌────────────────────────────────────────────────────────┐
│           Sprint 1: GitHub Ingestion Service           │  <-- COMPLETED (Current)
│   (Polls GitHub REST API, upserts github_projects)     │
└───────────────────────────┬────────────────────────────┘
                            │
                            ▼
┌────────────────────────────────────────────────────────┐
│        Phase 2: GitHub AI Project Analyzer             │
│   (Analyzes code structure, commits, and technologies) │
└───────────────────────────┬────────────────────────────┘
                            │
                            ▼
┌────────────────────────────────────────────────────────┐
│  Phase 3: Automatic Screenshot & Preview Generator    │
│   (Populates preview_image using headless browser)     │
└───────────────────────────┬────────────────────────────┘
                            │
                            ▼
┌────────────────────────────────────────────────────────┐
│        Phase 4: Automatic Portfolio Publishing         │
│   (Generates verified project cards & showcase pages)  │
└───────────────────────────┬────────────────────────────┘
                            │
                            ▼
┌────────────────────────────────────────────────────────┐
│        Phase 5: RAG Knowledge Integration              │
│   (Embeds synced repositories into pgvector knowledge) │
└────────────────────────────────────────────────────────┘
```
