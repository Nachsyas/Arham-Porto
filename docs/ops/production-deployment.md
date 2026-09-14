# Production Deployment Guide — Arham Porto

Project: **Arham Porto**  
Owner: **Nachsyas Arham Mumtaz Nashohi**  
Certified Baseline: **Phase 8 Production Architecture**

---

## 1. Production Architecture Overview

```
[ Visitor / Reviewer Browser ]
              │
              │ HTTPS
              ▼
   [ Vercel Edge Network ]
   Next.js 15 App Router Frontend (apps/web)
   Static Assets, Case Studies, Ask Arham UI
              │
              │ HTTPS + CORS (Restricted Origin)
              ▼
    [ Google Cloud Run ]
    Go API Service (apps/api)
    Multi-stage Container (/app/server)
    Trusted Proxy (TRUST_PROXY_MODE=cloudrun)
              │
      ┌───────┴────────────────────────┐
      ▼                                ▼
[ Supabase PostgreSQL 16 ]     [ Google Gemini API ]
  - pgvector extension          - gemini-3.8-flash (Reasoning)
  - schema_migrations           - gemini-embedding-2 (Retrieval)
  - knowledge_sources/chunks    - Zero-Trust Grounding Guards
```

---

## 2. Infrastructure Services Separation

The production deployment maintains strict isolation between presentation, application logic, persistence, and external AI providers:

| Service | Hosting Provider | Deployment Unit | Responsibility |
| :--- | :--- | :--- | :--- |
| **Frontend** | Vercel | Next.js 15 Monorepo Workspace (`apps/web`) | SSR/SSG Portfolio, Case Studies, SSE Chat Client |
| **API Backend** | Google Cloud Run | Docker Image (`apps/api/Dockerfile`) | Clean Architecture REST API, Proxy Rate Limiter, SSE Streaming |
| **Database** | Supabase | PostgreSQL 16 with `pgvector` extension | Relational storage & exact cosine vector similarity search |
| **AI Generation** | Google AI Studio | Gemini Interactions API (`v1beta/interactions`) | Structured, evidence-grounded responses |

---

## 3. Environment Variable Specifications (Names Only)

> [!IMPORTANT]
> Never store plaintext credentials or secrets in git or documentation. Only secret variable names are listed below.

### 3.1 Vercel Frontend (`apps/web`)

| Variable Name | Required | Description | Example / Allowed Values |
| :--- | :--- | :--- | :--- |
| `NEXT_PUBLIC_API_BASE_URL` | **Yes** | Public HTTPS URL of the Cloud Run Go API | `https://<service-hash>-<region>.a.run.app` |
| `NEXT_PUBLIC_SITE_URL` | **Yes** | Public canonical origin of the frontend | `https://arham-porto.vercel.app` |

### 3.2 Google Cloud Run Backend Go API (`apps/api`)

| Variable Name | Required | Description | Production Value |
| :--- | :--- | :--- | :--- |
| `PORT` | **Yes** | Port injected by Cloud Run runtime | Default `8080` (dynamically set by Cloud Run) |
| `APP_ENV` | **Yes** | Application environment | `production` |
| `ALLOWED_ORIGINS` | **Yes** | Exact allowed frontend origins for CORS | Exact Vercel domain (e.g. `https://arham-porto.vercel.app`) |
| `DATABASE_URL` | **Yes** | PostgreSQL connection string | Supabase connection secret |
| `DATABASE_MODE` | **Yes** | Database readiness requirement | `required` |
| `TRUST_PROXY_MODE` | **Yes** | Edge proxy client IP resolution | `cloudrun` |
| `PORTFOLIO_DATA_DIR` | **Yes** | Path to canonical portfolio content | `/app/data` |
| `MIGRATIONS_DIR` | **Yes** | Path to database schema migrations | `/app/migrations` |
| `AI_MODE` | **Yes** | Ask Arham AI operational mode | `remote` (or `disabled` if key pending) |
| `AI_PROVIDER` | **Yes** | Generative AI provider | `gemini` |
| `AI_MODEL` | **Yes** | Gemini generation model | `gemini-3.8-flash` |
| `AI_THINKING_LEVEL` | **Yes** | Gemini thinking level | `low` |
| `AI_REQUEST_TIMEOUT_SECONDS`| No | Timeout for Ask Arham SSE responses | `30` |
| `AI_MAX_CONCURRENT_REQUESTS`| No | Server concurrency semaphore | `4` |
| `AI_RATE_LIMIT_PER_MINUTE` | No | Per-client IP rate limit | `5` |
| `AI_MAX_EVIDENCE_CHARS` | No | Retrieval context character cap | `24000` |
| `EMBEDDING_MODE` | **Yes** | Vector embedding subsystem | `enabled` |
| `EMBEDDING_PROVIDER` | **Yes** | Embedding provider | `gemini` |
| `EMBEDDING_MODEL` | **Yes** | Embedding model name | `gemini-embedding-2` |
| `EMBEDDING_DIMENSIONS` | **Yes** | Vector dimensions | `768` |
| `GEMINI_API_KEY` | **Yes** | Google Gemini Auth API Key | Injected via Secret Manager / Cloud Run secret |
| `GITHUB_TOKEN` | Optional | GitHub API read-only token | Injected via Secret Manager / Cloud Run secret |

---

## 4. Supabase Database & pgvector Provisioning

1. Use Supabase PostgreSQL 16.
2. Enable the `vector` extension:
   ```sql
   CREATE EXTENSION IF NOT EXISTS vector;
   ```
3. Verify that `vector` is active:
   ```sql
   SELECT extname, extversion FROM pg_extension WHERE extname = 'vector';
   ```
4. Do not continue production indexing until `pgvector` extension existence is confirmed.

---

## 5. Pre-Deploy Migrations

The Go API container includes a dedicated, zero-dependency migration binary:
- **Binary Path**: `/app/migrate`
- **Workflow**:
  1. Build container image from `apps/api/Dockerfile`.
  2. Execute `/app/migrate` against production Supabase before routing public traffic.
  3. Execute migration command twice; the second execution must be idempotent and report that schema is up-to-date.
  4. Records successful migrations in the `schema_migrations` table.
  5. If any migration fails, deployment halts immediately.

---

## 6. Administrative One-Time Production Indexing

> [!WARNING]
> Indexing is strictly an administrative command and NEVER executes on application startup.

Once the Supabase database is healthy, `pgvector` confirmed, migrations completed, and `GEMINI_API_KEY` configured:
1. Execute a controlled one-time administrative indexing job:
   ```bash
   /app/indexer --all-approved
   ```
2. The indexer will:
   - Load canonical portfolio entities from `/app/data/`.
   - Ingest approved allowlisted repositories from GitHub.
   - Embed chunks using `gemini-embedding-2` (768 dimensions).
   - Upsert records into `knowledge_sources` and `knowledge_chunks`.
   - Confirm status as `vector_ready`.

---

## 7. Trusted Proxy & Rate Limiting

- **Google Cloud Run Forwarding Behavior**:
  - Google Cloud Run terminates TLS at Google Front End (GFE) and appends the connecting client IP to the `X-Forwarded-For` HTTP header chain.
  - Setting `TRUST_PROXY_MODE=cloudrun` instructs the API to parse the rightmost valid, non-internal client IP from `X-Forwarded-For`, preventing header spoofing from malicious initial values.
  - In `direct` mode (local testing), proxy headers are strictly ignored and `RemoteAddr` is used.
- **Rate Limiting**:
  - General API endpoints: 120 requests/minute per client IP.
  - Ask Arham AI copilot: 5 requests/minute per client IP.

---

## 8. Health & Readiness Verification

- **`/health` (External Liveness)**:
  - Verifies HTTP process responsiveness and server liveness.
  - Returns `200 OK` (`{"status":"ok","scope":"process_alive"}`).
  - Recommended endpoint for external uptime monitors and synthetic checks.
- **`/readyz` (Readiness)**:
  - Verifies PostgreSQL database connectivity.
  - Returns `200 OK` when `DATABASE_MODE=required` and database is connected (`{"status":"ready","database":"connected"}`).
  - Used for Cloud Run startup and readiness probes.
- **`/healthz` (Local Compatibility Probe)**:
  - Retained for local compatibility and standard container liveness probes.
  - Returns `200 OK` when accessed directly on the container.
  - Note: In Cloud Run production, an edge anomaly has been observed where requests directly to `/healthz` return a Google edge 404 before reaching the container. While not an officially documented Google reservation rule, `/health` serves as the public liveness alias.

---

## 9. AI Degraded State & Failure Recovery

- **Ask Arham AI Status**: `PENDING PRODUCTION ACTIVATION`.
  - Core portfolio deployment is active and fully functional on Vercel and Google Cloud Run.
  - The API serves all canonical portfolio data (`/api/v1/profile`, `/api/v1/projects`, `/api/v1/skills`, `/api/v1/evidence`, `/api/v1/journey`).
  - `/api/v1/ai/ask` responds with HTTP 503 Service Unavailable (`{"error":{"code":"service_unavailable","message":"Ask Arham AI is currently unavailable."}}`).
  - Real embedding indexing via Cloud Run Job `arham-porto-indexer` requires an active Google AI Studio API key with available prepayment credits / quota.
- If PostgreSQL becomes unreachable:
  - `/readyz` fails (HTTP 503), preventing unready instances from receiving traffic.
