# Architecture: System Overview

- **Purpose**: High-level structural overview of the Arham Porto platform.
- **Current Phase**: Phase 0 (Bootstrap Foundation).

---

## 1. System Architecture Diagram

```
[ Reviewer / User ]
       │
       ▼
[ Next.js Web App (apps/web) ]
  ├── UI Sections (Hero, Quick Review, Projects, Skills, Experience, Education)
  ├── 3D Seated Hologram (Three.js / R3F)
  ├── Journey Map (MapLibre GL JS)
  └── AI Reviewer Panel (SSE Client)
       │
       ▼ (REST / SSE)
[ Go Backend (apps/api) ]
  ├── Delivery (HTTP router net/http, CORS, Rate Limiter)
  ├── UseCases (Profile, Projects, Evidence, AI Orchestrator)
  ├── Domain (Pure Go entities)
  └── Infrastructure Adapters
       ├── PostgreSQL 16 + pgvector (Embeddings & Knowledge Chunks)
       ├── GitHub Indexer Adapter (Approved Repositories)
       └── Provider-Agnostic LLM / Embedding Client
```

---

## 2. Canonical Source of Truth (Phase 0–3)
- Canonical professional content resides in `data/` and is validated via `packages/portfolio-data`.
- PostgreSQL with `pgvector` serves as the runtime storage engine for AI knowledge indexing and retrieval (Phases 4–7).
