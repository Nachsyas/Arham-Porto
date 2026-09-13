# HTTP API Documentation — Arham Porto

Backend REST API documentation for `apps/api`.

---

## 1. Overview & Architecture Standards

- **Server Runtime**: Go 1.22+ Standard Library `net/http.ServeMux` (Zero third-party routing frameworks).
- **Clean Architecture**: Delivery $\rightarrow$ UseCase $\rightarrow$ Domain $\leftarrow$ Infrastructure/Repository.
- **Purity Rule**: `internal/domain` contains zero external dependencies.
- **DTO Isolation**: Domain entities are never serialized directly to HTTP. Responses pass through public DTO mappers in `internal/delivery/http/dto/`.
- **Operating Modes**: Explicit `DATABASE_MODE=disabled|optional|required`.

---

## 2. Global Envelopes

### Success Envelope (Single Resource)
```json
{
  "data": { ... }
}
```

### Success Envelope (Collection)
```json
{
  "data": [ ... ],
  "meta": {
    "count": 4
  }
}
```

### Error Envelope
```json
{
  "error": {
    "code": "bad_request | not_found | rate_limited | internal_error | service_unavailable",
    "message": "Human-readable description of error"
  }
}
```

---

## 3. Endpoints

### 3.1 Probes

#### `GET /healthz`
- **Purpose**: Process liveness probe.
- **Rate Limit**: Bypasses rate limiting.
- **Cache-Control**: `no-store`
- **Response** (`200 OK`):
  ```json
  {
    "status": "ok",
    "scope": "process_alive"
  }
  ```

#### `GET /readyz`
- **Purpose**: Readiness probe indicating dependency readiness.
- **Rate Limit**: Bypasses rate limiting.
- **Cache-Control**: `no-store`
- **Response** (`200 OK` or `503 Service Unavailable`):
  - When `DATABASE_MODE=disabled`: `200 OK` `{"status":"ready","database":"disabled"}`
  - When `DATABASE_MODE=optional` & DB connected: `200 OK` `{"status":"ready","database":"connected"}`
  - When `DATABASE_MODE=optional` & DB offline: `200 OK` `{"status":"ready","database":"degraded"}`
  - When `DATABASE_MODE=required` & DB offline: `503 Service Unavailable` `{"status":"not_ready","database":"unavailable"}`

---

### 3.2 Portfolio Resources

#### `GET /api/v1/profile`
- **Description**: Returns verified professional profile details.
- **Cache-Control**: `public, max-age=60, stale-while-revalidate=300`
- **Response** (`200 OK`):
  ```json
  {
    "data": {
      "full_name": "Nachsyas Arham Mumtaz Nashohi",
      "role": "Software Engineer",
      "project_name": "Arham Porto",
      "ai_feature": "Ask Arham AI",
      "positioning": "Building intelligent, scalable, and human-centered digital systems.",
      "github": "https://github.com/Nachsyas"
    }
  }
  ```
- **Privacy Guard**: Excludes phone numbers, birth dates, home addresses, coordinates, and internal TODO metadata.

#### `GET /api/v1/projects`
- **Description**: Lists engineering case studies and systems.
- **Query Parameters**:
  - `category` (optional): Filter by category (`AI`, `Full-Stack`, `Backend`, `Systems`). Returns `400 Bad Request` if unsupported.
  - `featured` (optional): Filter by boolean (`true` or `false`). Returns `400 Bad Request` if not a valid boolean.
- **Cache-Control**: `public, max-age=60, stale-while-revalidate=300`
- **Response** (`200 OK`):
  ```json
  {
    "data": [
      {
        "id": "edutrace",
        "title": "EduTrace",
        "slug": "edutrace",
        "summary": "Decentralized academic record ledger using Soulbound Tokens (ERC-5192) and predictive time-series performance analysis.",
        "problem": "Academic credential verification challenges and lack of early warning indicators for student retention.",
        "solution": "Immutable credential minting on an EVM ledger paired with an analytics dashboard and Go event indexing pipeline.",
        "role": [],
        "contributions": [],
        "technologies": ["Solidity", "Foundry", "Next.js", "TypeScript", "Go", "Python", "Tailwind CSS"],
        "github_url": "https://github.com/Nachsyas/EduTrace",
        "demo_url": "https://edu-trace-nine.vercel.app",
        "featured": true,
        "category": "Full-Stack",
        "evidence_ids": ["evidence-edutrace-repo"]
      }
    ],
    "meta": {
      "count": 1
    }
  }
  ```

#### `GET /api/v1/projects/{slug}`
- **Description**: Retrieves single project by slug.
- **Path Parameter**: `slug` (alphanumeric with hyphens, $\le$ 64 chars). Returns `400 Bad Request` if malformed.
- **Response**: `200 OK` or `404 Not Found`.

#### `GET /api/v1/skills`
- **Description**: Lists verified skills.
- **Query Parameters**:
  - `category` (optional): Filter by category (`Backend`, `Frontend`, `AI / ML`, `Systems`, `DevOps`, `Database`).
- **Response** (`200 OK`):
  ```json
  {
    "data": [
      {
        "id": "backend-go",
        "name": "Backend Engineering with Go",
        "category": "Backend",
        "claim": "Building modular backend services and REST APIs with modern Go and Clean Architecture.",
        "evidence_ids": ["evidence-gdgoc-ecommerce-repo", "evidence-edutrace-repo", "evidence-smart-kitchen-backend-repo"]
      }
    ],
    "meta": {
      "count": 1
    }
  }
  ```

#### `GET /api/v1/evidence`
- **Description**: Lists verifiable artifacts.
- **Query Parameters**:
  - `skill_id` (optional): Filter by associated skill ID.
  - `project_id` (optional): Filter by associated project ID.
- **Response** (`200 OK`):
  ```json
  {
    "data": [
      {
        "id": "evidence-edutrace-repo",
        "type": "github",
        "title": "EduTrace Repository & Architecture Specification",
        "skill_ids": ["backend-go", "frontend-nextjs", "database-postgres"],
        "source_url": "https://github.com/Nachsyas/EduTrace",
        "source_path": "README.md",
        "summary": "Decentralized academic achievement system combining Soulbound Tokens (ERC-5192) via Solidity, Next.js 15 frontend, and analytical backend services with Go and Python.",
        "verified": true
      }
    ],
    "meta": {
      "count": 1
    }
  }
  ```

#### `GET /api/v1/evidence/{id}`
- **Description**: Retrieves evidence artifact by ID.
- **Response**: `200 OK` or `404 Not Found`.

#### `GET /api/v1/journey`
- **Description**: Lists public educational and geographical milestones (strictly `public == true`).
- **Response** (`200 OK`):
  ```json
  {
    "data": [
      {
        "id": "origin-karanganyar",
        "category": "birthplace",
        "title": "Origin",
        "city": "Karanganyar",
        "region": "Jawa Tengah",
        "country": "Indonesia",
        "description": "Early formative origin in Karanganyar, Central Java.",
        "public": true
      }
    ],
    "meta": {
      "count": 7
    }
  }
  ```
- **Privacy Guard**: Raw coordinates, internal TODOs, and private records (e.g. `journey-tk`) are strictly omitted.

---

## 4. Middleware Pipeline Order

Outer $\rightarrow$ Inner:
1. `Recovery`: Catches panics, logs diagnostics on the server, returns generic 500 JSON without leaking stack traces.
2. `Logger`: Privacy-safe logging (timestamp, method, path, status, latency, IP). No credentials/headers/bodies logged.
3. `SecurityHeaders`:
   - `X-Content-Type-Options: nosniff`
   - `Referrer-Policy: no-referrer`
   - `Content-Security-Policy: default-src 'none'; frame-ancestors 'none'`
   - `X-Frame-Options: DENY`
4. `CORS`: Allowed origins parsed from `ALLOWED_ORIGINS` (default `http://localhost:3000`). Adds `Vary: Origin`. Preflight `OPTIONS` returns `204 No Content`.
5. `RateLimiter`: Bounded in-memory sliding window (120 req/min per IP) with TTL cleanup. `/healthz` and `/readyz` strictly bypass rate limiting.
6. `ServeMux`: Standard library route multiplexer.
