# Golden Retrieval Evaluation & Groundedness Benchmark (Phase 5 Certified)

## 1. Methodology & Honest Reporting Protocol
Per Gate #34 and AGENTS.md §2.4:
Evaluations are strictly separated into two distinct categories:
1. **Structural Retrieval Tests (Automated & CI-Safe)**:
   - Uses `DeterministicFakeProvider` with unit-normalized vectors.
   - Verifies vector dimensionality (768), strict model identity filtering, relational metadata propagation (`project_id`, `skill_ids`, `evidence_id`), SQL query generation, and exact cosine distance ranking.
   - **Explicit Boundary**: Does NOT evaluate semantic language understanding and does NOT claim semantic Recall@K.
2. **Real Retrieval Quality Evaluation**:
   - Requires live credentials (`GEMINI_API_KEY`) and active PostgreSQL + pgvector.
   - Evaluates natural-language queries against indexed corpus using `gemini-embedding-2` (768 dimensions).
   - Reports empirical Recall@K and query latency.
   - **Certification Status**: If `GEMINI_API_KEY` is omitted in development/CI, Phase 5 honestly reports:
     ```
     REAL PROVIDER EVALUATION: NOT EXECUTED — CREDENTIAL NOT PROVIDED
     ```
     Fabricating semantic evaluation scores from fake vectors is strictly forbidden.

---

## 2. Golden Question Dataset (Gate #35)

### A. Positive Benchmark Queries
| Query | Expected Authoritative Source | Key Evidence Identifiers | Target Technologies |
| :--- | :--- | :--- | :--- |
| *"What project demonstrates Go backend experience?"* | `gdgoc-ecommerce`, `smart-kitchen-backend`, `EduTrace` | `backend/cmd/api/main.go`, `main.go`, `evidence-gdgoc-ecommerce-repo` | Go, Clean Architecture, Fiber |
| *"Which project uses Soulbound Tokens?"* | `EduTrace` | `contracts/src/EduTraceSBT.sol`, `README.md` | Solidity, ERC-5192 / ERC-5173, Foundry |
| *"What database is used by GDGOC E-Commerce?"* | `gdgoc-ecommerce` | `backend/go.mod`, `evidence-gdgoc-ecommerce-repo` | MongoDB, Go Driver |
| *"What technologies power Smart Kitchen backend?"* | `smart-kitchen-backend` | `go.mod`, `database/connection.go` | Go, Fiber, GORM, PostgreSQL |
| *"Where did Arham study Computer Science?"* | `canonical_journey` | `journey/university-uin-malang.json` | UIN Maulana Malik Ibrahim Malang |

### B. Negative & Guardrail Benchmark Queries
| Query | Expected Authoritative Evidence | Evaluation Rule |
| :--- | :--- | :--- |
| *"Does Maritime AI prove RAG experience?"* | **No authoritative evidence** | Chunks may mention Python/FastAPI/Next.js, but must NOT confirm RAG claims without explicit proof. |
| *"Does Nachsyas have production Kubernetes evidence?"* | **No authoritative evidence** | Docker containerization exists, but production Kubernetes is unverified. Must not extrapolate certainty. |
| *"Did Nachsyas personally implement every module in GDGOC E-Commerce?"* | **No authoritative evidence** | Contributions specify backend architecture and REST API; does not claim entire monolithic/multi-tenant suite alone. |
| *"What is Nachsyas's home address and phone number?"* | **Redacted / Omitted** | Private personal data is completely excluded from knowledge ingestion (Gate #16). |

---

## 3. Threshold Calibration Strategy
Per Gate #36:
- Arbitrary similarity thresholds (e.g. `0.80` or `0.85`) are **not hardcoded** prior to measuring real embedding models on the complete corpus.
- Top-K retrieval returns the ranked candidates with exact cosine distances.
- Threshold determination will be calibrated during Phase 6 groundedness tuning based on positive vs negative query distance separation.
