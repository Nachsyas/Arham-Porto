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
   - Railway builds the container, executes `/app/migrate`, and activates `/app/server` once `/readyz` responds with HTTP 200.

---

## 2. Rollback Procedures

### 2.1 Frontend Rollback (Vercel)
If a frontend defect is detected in production:
1. Open the **Vercel Dashboard** $\rightarrow$ **Arham Porto** $\rightarrow$ **Deployments**.
2. Locate the previous successful production deployment.
3. Click **Instant Rollback** (or click the three dots $\rightarrow$ **Promote to Production**).
4. Traffic is immediately redirected to the previous healthy build within seconds.

### 2.2 Backend API Rollback (Railway)
If a backend regression or container failure occurs:
1. Open the **Railway Dashboard** $\rightarrow$ **Arham Porto** $\rightarrow$ **API Service** $\rightarrow$ **Deployments**.
2. Find the last healthy deployment.
3. Click the deployment $\rightarrow$ **Rollback to this deployment**.
4. Railway will immediately route incoming HTTPS traffic back to the chosen deployment replica.

### 2.3 Database Migration Safeguards
- Schema migrations must **NEVER** be rolled back destructively in production.
- If a migration fails during the Railway pre-deploy step:
  1. The pre-deploy command `/app/migrate` exits with non-zero status.
  2. Railway cancels the new deployment and keeps the previous container online.
  3. Inspect logs via Railway or CLI:
     ```bash
     railway logs --service api
     ```
  4. Fix the migration script forward in a new commit, test locally against a clone, and push to `main`.

---

## 3. Incident Response & Troubleshooting

### 3.1 Database Unavailable (`/readyz` returning 503)
**Symptom**: Railway container restarts or healthcheck fails; `/readyz` returns `{"status":"not_ready","database":"unavailable"}`.
**Actions**:
1. Check Railway PostgreSQL service status:
   - Ensure the database container is in the `Running` state.
   - Verify connection credentials in `DATABASE_URL`.
2. Inspect connection pool saturation:
   - Check Railway CPU/Memory metrics for the database.
   - Verify that client connections are within limits (default pool size: 10 connections).
3. If PostgreSQL is recovering:
   - The Go API automatically retries connections using exponential backoff.
   - Once the database accepts connections, `/readyz` returns HTTP 200 without requiring an API restart.

### 3.2 Gemini API Outage or Quota Exhaustion
**Symptom**: Ask Arham returns HTTP 503 with code `service_unavailable`; core portfolio pages remain operational.
**Actions**:
1. Verify Google AI Studio operational status and billing quota.
2. Check API container logs for rate limit or quota responses:
   ```bash
   railway logs --service api | grep "Ask Arham"
   ```
3. If necessary, temporarily disable the AI copilot without redeploying code:
   - In Railway API service settings, set `AI_MODE=disabled`.
   - The API will gracefully inform visitors that Ask Arham is resting while preserving 100% of portfolio navigation.

### 3.3 Gemini API Key Rotation
To rotate `GEMINI_API_KEY` without service downtime:
1. Generate a new API key in Google AI Studio.
2. In the Railway Dashboard:
   - Navigate to **Variables** for the API service.
   - Update `GEMINI_API_KEY` with the new value.
3. Railway performs a rolling restart of the container.
4. Verify functionality:
   ```bash
   curl -X POST https://<railway-domain>/api/v1/ai/ask \
     -H "Content-Type: application/json" \
     -d '{"question":"What does Arham specialize in?"}'
   ```
5. Revoke the retired key in Google AI Studio only after the new deployment is verified.

---

## 4. Administrative Indexing Procedures

### 4.1 Running a Full Re-Index
When new projects or skills are approved in canonical data:
1. Run a one-off task in Railway:
   ```bash
   railway run /app/indexer --all-approved
   ```
2. Monitor output logs:
   - Confirms canonical source extraction.
   - Generates embeddings with `gemini-embedding-2`.
   - Validates that `knowledge_sources` and `knowledge_chunks` are in `vector_ready` state.

### 4.2 Checking Vector Index Health
To inspect the production vector index directly via SQL:
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
