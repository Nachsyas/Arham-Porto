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
    [ Railway Platform ]
    Go API Service (apps/api)
    Multi-stage Container (/app/server)
    Trusted Proxy (TRUST_PROXY_MODE=railway)
              │
      ┌───────┴────────────────────────┐
      ▼                                ▼
[ Railway PostgreSQL 16 ]     [ Google Gemini API ]
  - pgvector extension         - gemini-3.8-flash (Reasoning)
  - schema_migrations          - gemini-embedding-2 (Retrieval)
  - knowledge_sources/chunks   - Zero-Trust Grounding Guards
```

---

## 2. Infrastructure Services Separation

The production deployment maintains strict isolation between presentation, application logic, persistence, and external AI providers:

| Service | Hosting Provider | Deployment Unit | Responsibility |
| :--- | :--- | :--- | :--- |
| **Frontend** | Vercel | Next.js 15 Monorepo Workspace (`apps/web`) | SSR/SSG Portfolio, Case Studies, SSE Chat Client |
| **API Backend** | Railway | Docker Image (`apps/api/Dockerfile`) | Clean Architecture REST API, Proxy Rate Limiter, SSE Streaming |
| **Database** | Railway | PostgreSQL 16 with `pgvector` template | Relational storage & exact cosine vector similarity search |
| **AI Generation** | Google AI Studio | Gemini Interactions API (`v1beta/interactions`) | Structured, evidence-grounded responses |

---

## 3. Environment Variable Specifications (Names Only)

> [!IMPORTANT]
> Never store plaintext credentials or secrets in git or documentation. Only secret variable names are listed below.

### 3.1 Vercel Frontend (`apps/web`)

| Variable Name | Required | Description | Example / Allowed Values |
| :--- | :--- | :--- | :--- |
| `NEXT_PUBLIC_API_BASE_URL` | **Yes** | Public HTTPS URL of the Railway Go API | `https://api.arhamporto.com` or Railway generated domain |
| `NEXT_PUBLIC_SITE_URL` | **Yes** | Public canonical origin of the frontend | `https://arhamporto.com` or `https://arham-porto.vercel.app` |

### 3.2 Railway Backend Go API (`apps/api`)

| Variable Name | Required | Description | Production Value |
| :--- | :--- | :--- | :--- |
| `PORT` | **Yes** | Port injected by Railway runtime | Injected dynamically by Railway |
| `APP_ENV` | **Yes** | Application environment | `production` |
| `ALLOWED_ORIGINS` | **Yes** | Exact allowed frontend origins for CORS | Exact Vercel domain (e.g. `https://arham-porto.vercel.app`) |
| `DATABASE_URL` | **Yes** | PostgreSQL connection string | Railway private reference `${{Postgres.DATABASE_URL}}` |
| `DATABASE_MODE` | **Yes** | Database readiness requirement | `required` |
| `TRUST_PROXY_MODE` | **Yes** | Edge proxy client IP resolution | `railway` |
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
| `GEMINI_API_KEY` | **Yes** | Google Gemini Auth API Key | Injected as a Railway secret |
| `GITHUB_TOKEN` | Optional | GitHub API read-only token | Injected as a Railway secret (for indexing) |

---

## 4. Railway PostgreSQL & pgvector Provisioning

1. Provision PostgreSQL using the **pgvector** template on Railway.
2. Ensure the database resides in the same Railway project and region as the API service.
3. Validate that `pgvector` is installed by running:
   ```sql
   SELECT extname FROM pg_extension WHERE extname = 'vector';
   ```
   If not yet installed, the initial migration automatically executes:
   ```sql
   CREATE EXTENSION IF NOT EXISTS vector;
   ```

---

## 5. Automated Pre-Deploy Migrations

The Go API container includes a dedicated, dependency-free migration binary:
- **Binary Path**: `/app/migrate`
- **Execution**: Configured as Railway's `preDeployCommand` in `railway.json`.
- **Workflow**:
  1. Railway builds the production image from `apps/api/Dockerfile`.
  2. Prior to routing traffic to the new deployment, `/app/migrate` runs.
  3. Applies migrations from `/app/migrations` in deterministic version order (`000001_...`, `000002_...`).
  4. Records successful migrations in the `schema_migrations` table.
  5. If any migration fails, deployment halts non-zero and the previous healthy deployment remains active.

---

## 6. Administrative One-Time Production Indexing

> [!WARNING]
> Indexing is strictly an administrative command and NEVER executes on application startup.

Once the database is migrated and `GEMINI_API_KEY` is configured:
1. Execute a one-off administrative task in Railway:
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

- **Railway Proxy Configuration**:
  - Railway edge proxies terminate SSL and forward the real client IP in the `X-Real-IP` HTTP header.
  - Setting `TRUST_PROXY_MODE=railway` instructs the API to validate and extract the client identity from `X-Real-IP`.
  - In `direct` mode (local testing), `X-Real-IP` is strictly ignored to prevent header spoofing.
- **Rate Limiting**:
  - General API endpoints: 120 requests/minute per client IP.
  - Ask Arham AI copilot: 5 requests/minute per client IP.

---

## 8. Health & Readiness Verification

- **`/healthz` (Liveness)**:
  - Verifies HTTP process responsiveness.
  - Returns `200 OK` (`{"status":"ok","scope":"process_alive"}`).
- **`/readyz` (Readiness)**:
  - Verifies PostgreSQL connectivity.
  - Returns `200 OK` when `DATABASE_MODE=required` and database is connected.
  - Used as Railway's deployment healthcheck (`healthcheckPath: "/readyz"`).

---

## 9. AI Degraded State & Failure Recovery

- If the Gemini API key is missing or invalid:
  - The API continues serving core portfolio data (`/api/v1/profile`, `/api/v1/projects`, `/api/v1/skills`, etc.).
  - `/api/v1/ai/ask` responds with HTTP 503 Service Unavailable and a clear JSON error envelope.
  - The frontend UI displays an informative "Ask Arham is temporarily unavailable" notice without breaking page navigation.
- If PostgreSQL becomes unreachable:
  - `/readyz` fails, preventing unready deployments from receiving traffic.
