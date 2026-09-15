# Production Runbook — Arham Porto

Project: **Arham Porto**  
Owner: **Nachsyas Arham Mumtaz Nashohi**  
Certified Baseline: **Phase 8 Operations**

---

## 1. Routine Deployment Workflow

### 1.1 Triggering a Deployment
1. Ensure all local tests pass:
   ```bash
   npm run validate:data
   npm run test --workspace=apps/web
   cd apps/api && go test -count=1 -race ./...
   ```
2. Commit and push changes to `main`:
   ```bash
   git push origin main
   ```
3. GitHub Actions CI automatically triggers and runs validation jobs:
   - **Frontend & Data Schema**
   - **Backend Go API**
4. Once CI passes:
   - Vercel automatically deploys the frontend from `apps/web`.
   - Google Cloud Run receives the updated container image, verifies `/readyz` healthcheck, and shifts traffic.

---

## 2. Rollback Procedures

### 2.1 Frontend Rollback (Vercel)
If a frontend defect is detected in production:
1. Open the **Vercel Dashboard** $\rightarrow$ **Arham Porto** $\rightarrow$ **Deployments**.
2. Locate the previous successful production deployment.
3. Click **Instant Rollback** (or click the three dots $\rightarrow$ **Promote to Production**).
4. Traffic is immediately redirected to the previous healthy build.

### 2.2 Backend API Rollback (Google Cloud Run)
If a backend regression or container failure occurs:
1. Open Google Cloud Console $\rightarrow$ **Cloud Run** $\rightarrow$ `arham-porto-api` $\rightarrow$ **Revisions**.
2. Identify the last certified healthy revision.
3. Select **Manage Traffic** $\rightarrow$ route 100% of traffic to that revision.
4. Traffic shifts instantly with zero downtime.

### 2.3 Database Migration Safeguards
- Schema migrations must **NEVER** be rolled back destructively in production.
- If a migration fails during the pre-deploy or release step:
  1. The migration command `/app/migrate` exits with non-zero status.
  2. The Cloud Run revision is not activated.
  3. Inspect logs via Cloud Logging:
     ```bash
     gcloud logging read 'resource.type="cloud_run_revision" AND severity>=ERROR' --limit 20
     ```
  4. Fix the migration script forward in a new commit, test locally against a clone, and push to `main`.

---

## 3. Incident Response & Troubleshooting

### 3.1 Database Unavailable (`/readyz` returning 503)
**Symptom**: Cloud Run instances fail healthchecks; `/readyz` returns `{"status":"not_ready","database":"unavailable"}`.
**Actions**:
1. Check Supabase project status:
   - Ensure the database is active and not paused.
   - Verify connection pooler settings (Transaction pooler port 6543 or Session pooler port 5432).
2. Inspect connection pool saturation:
   - Check Supabase metrics for active connections.
   - Ensure Cloud Run container max concurrency is aligned with database pool limits.
3. If PostgreSQL is recovering:
   - The Go API automatically retries connections with exponential backoff.
   - Once the database accepts connections, `/readyz` returns HTTP 200 without requiring container restarts.

### 3.2 Groq / Cloudflare Outage or Quota Exhaustion
**Symptom**: Ask Arham returns HTTP 503 with code `service_unavailable`; core portfolio pages remain operational.
**Actions**:
1. Check Groq operational status (status.groq.com) and Cloudflare Workers AI operational status.
2. Check Cloud Run logs for provider rate limit or quota responses.
3. If necessary, temporarily disable the AI copilot without redeploying code:
   - Update Cloud Run environment variable `AI_MODE=disabled`.
   - The API will gracefully inform visitors that Ask Arham is resting while preserving 100% of portfolio navigation.

### 3.3 Provider Credential Rotation
To rotate `GROQ_API_KEY` or `CLOUDFLARE_AI_TOKEN` without service downtime:
1. Generate a new API key/token in Groq Console or Cloudflare Dashboard.
2. In Google Cloud Secret Manager:
   - Add a new secret version to `arham-porto-groq-api-key` or `arham-porto-cloudflare-ai-token`.
3. Cloud Run automatically picks up the `:latest` version on new revision deployment or restart.
4. Verify functionality:
   ```bash
   curl -X POST https://arham-porto-api-6ivuekmjva-et.a.run.app/api/v1/ai/ask \
     -H "Content-Type: application/json" \
     -H "Accept: text/event-stream" \
     -d '{"question":"What does Arham specialize in?"}'
   ```
5. Revoke the retired key in the provider console only after the new deployment is verified.

---

## 4. Administrative Indexing Procedures

### 4.1 Running Production Indexing via Cloud Run Job
When new projects or skills are approved in canonical data, execute the dedicated administrative Cloud Run Job:
```bash
gcloud run jobs execute arham-porto-indexer \
  --project=project-c3c283e3-36a1-43ec-a7f \
  --region=asia-southeast2 \
  --wait
```
The job executes `/app/indexer --all-approved` using:
- `EMBEDDING_PROVIDER=cloudflare`
- `EMBEDDING_MODEL=@cf/baai/bge-base-en-v1.5` (768 dimensions)
- Secret bindings for `DATABASE_URL` and `CLOUDFLARE_AI_TOKEN`
- Persists embeddings into Supabase PostgreSQL + `pgvector`

### 4.2 Inspecting Vector Index Status
To query the operational index status without running ingestion:
```bash
gcloud run jobs execute arham-porto-indexer \
  --project=project-c3c283e3-36a1-43ec-a7f \
  --region=asia-southeast2 \
  --args=--status \
  --wait
```
Or directly in Supabase SQL Editor:
```sql
-- Check total embedded chunks
SELECT embedding_model, COUNT(*) 
FROM knowledge_chunks 
GROUP BY embedding_model;

-- Check latest indexed sources
SELECT repository, source_type, title, indexed_at 
FROM knowledge_sources 
ORDER BY indexed_at DESC 
LIMIT 10;
```
