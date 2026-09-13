# HTTP API Documentation — Arham Porto

Backend REST API documentation for `apps/api`.

---

## 1. Overview & Architecture Standards

- **Server Runtime**: Go 1.22+ Standard Library `net/http.ServeMux` (Zero third-party routing frameworks).
- **Clean Architecture**: Delivery $\rightarrow$ UseCase $\rightarrow$ Domain $\leftarrow$ Infrastructure/Repository.
- **Purity Rule**: `internal/domain` contains zero external dependencies.
- **DTO Isolation**: Domain entities are never serialized directly to HTTP. Responses pass through public DTO mappers in `internal/delivery/http/dto/`.
- **Operating Modes**: Explicit `DATABASE_MODE=disabled|optional|required`.
- **Content-Type**: `application/json; charset=utf-8` on all JSON responses and errors.
- **Server Timeouts**: ReadHeaderTimeout 5s, ReadTimeout 10s, WriteTimeout 10s, IdleTimeout 60s, MaxHeaderBytes 1MB.

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
    "code": "bad_request | not_found | method_not_allowed | rate_limited | internal_error | service_unavailable",
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
  - When `DATABASE_MODE=optional` & DB offline: `200 OK` `{"status":"ready","database":"degraded"}` (remains degraded until restart)
  - When `DATABASE_MODE=required` & DB offline: Server aborts startup cleanly. In test/mock harness: `503 Service Unavailable` `{"status":"not_ready","database":"unavailable"}`

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
      "github": "https://github.com/Nachsyas"
    }
  }
  ```
- **Privacy Boundary**: All unapproved or TODO-backed fields (`positioning`, `bio`, `current_city`, `availability`, `email`, `linkedin`) are safely stripped/omitted. Excludes phone numbers, birth dates, home addresses, and private coordinates.

#### `GET /api/v1/projects`
- **Description**: Lists engineering case studies and systems.
- **Cache-Control**: `public, max-age=60, stale-while-revalidate=300`
- **Query Parameters**:
  - `category` (optional): Filter by category (`AI`, `Full-Stack`, `Backend`, `Systems`). Repeated values return `400 Bad Request`.
  - `featured` (optional): Filter by boolean (`true` or `false`). Invalid values or repeated values return `400 Bad Request`.
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
      "count": 4
    }
  }
  ```

#### `GET /api/v1/projects/{slug}`
- **Description**: Retrieves a single project by slug.
- **Slug Validation**: Alphanumeric and hyphens only (`^[a-z0-9-]+$`), maximum 64 characters.
- **Response**: `200 OK` or `404 Not Found`.

#### `GET /api/v1/skills`
- **Description**: Lists verified engineering skills and claims.
- **Cache-Control**: `public, max-age=60, stale-while-revalidate=300`
- **Query Parameters**:
  - `category` (optional): Filter by category (`Backend`, `Frontend`, `AI / ML`, `Systems`, `DevOps`, `Database`). Repeated values return `400 Bad Request`.
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
      "count": 7
    }
  }
  ```

#### `GET /api/v1/evidence`
- **Description**: Lists verifiable artifacts and citation anchors.
- **Cache-Control**: `public, max-age=60, stale-while-revalidate=300`
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
      "count": 4
    }
  }
  ```
- **Path Safety**: `source_path` contains strictly repository-relative paths. Absolute local filesystem paths (`/Users/...`, `C:\...`) are never serialized.

#### `GET /api/v1/evidence/{id}`
- **Description**: Retrieves evidence artifact by ID.
- **Response**: `200 OK` or `404 Not Found`.

#### `GET /api/v1/journey`
- **Description**: Lists public educational and geographical milestones (strictly `public == true`).
- **Cache-Control**: `public, max-age=60, stale-while-revalidate=300`
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
        "country": "Indonesia"
      }
    ],
    "meta": {
      "count": 7
    }
  }
  ```
- **Privacy Guard**: Raw coordinates, internal TODOs, private records (e.g. `journey-tk`), redundant `public` flag, and unapproved narrative descriptions are omitted.

---

## 4. HTTP Method Contract & Status Codes

| Method | Registered Path | Status Code | Notes |
| :--- | :--- | :--- | :--- |
| `GET` | All endpoints above | `200 OK` | Public cache header where applicable |
| `HEAD` | All endpoints above | `200 OK` | Body discarded by standard library |
| `OPTIONS` | Any registered path | `204 No Content` | CORS preflight with ACAO and Vary headers |
| `POST` | Any GET-only endpoint | `405 Method Not Allowed` | Header `Allow: GET, HEAD`, JSON error envelope |
| `PUT` | Any GET-only endpoint | `405 Method Not Allowed` | Header `Allow: GET, HEAD`, JSON error envelope |
| `DELETE` | Any GET-only endpoint | `405 Method Not Allowed` | Header `Allow: GET, HEAD`, JSON error envelope |
| `GET` | Unregistered path | `404 Not Found` | JSON ErrorEnvelope |
| `GET` | Ambiguous repeated param | `400 Bad Request` | JSON ErrorEnvelope |

### 4.1 Certified 405 Method Not Allowed Response
```http
HTTP/1.1 405 Method Not Allowed
Allow: GET, HEAD
Content-Type: application/json; charset=utf-8
Cache-Control: no-store
X-Content-Type-Options: nosniff
Referrer-Policy: no-referrer
Content-Security-Policy: default-src 'none'; frame-ancestors 'none'
X-Frame-Options: DENY

{
  "error": {
    "code": "method_not_allowed",
    "message": "method not allowed"
  }
}
```

---

## 5. Security & Rate Limiting

- **Security Headers**: Injected on all responses (200, 204, 400, 404, 405, 429, 500):
  - `X-Content-Type-Options: nosniff`
  - `Referrer-Policy: no-referrer`
  - `Content-Security-Policy: default-src 'none'; frame-ancestors 'none'`
  - `X-Frame-Options: DENY`
- **CORS**:
  - Validates `Origin` header against `ALLOWED_ORIGINS`.
  - Appends `Vary: Origin` cleanly without header clobbering.
- **Rate Limiter**:
  - Default: 120 req/min per normalized client IP (ephemeral TCP ports stripped via `net.SplitHostPort`).
  - Max entries bounded at 10,000 to prevent memory exhaustion under spoofed client cardinality.
  - `/healthz` and `/readyz` bypass rate limiting unconditionally.
