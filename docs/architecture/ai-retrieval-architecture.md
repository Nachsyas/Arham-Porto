# Architecture: AI & Retrieval Pipeline

- **Purpose**: Define the evidence retrieval, indexing, and LLM orchestration flow for Ask Arham AI.
- **Current Phase**: Phase 0 (Bootstrap Foundation).

---

## 1. Grounded RAG Pipeline (Phases 5–6)

```
[ Approved GitHub Repositories ] ──> [ Go Indexer ] ──> [ Chunker & Embedding ] ──> [ pgvector ]
                                                                                           │
[ Reviewer Question / JD ] ──> [ Go API ] ──> [ Semantic Vector Search ] <─────────────────┘
                                     │
                                     ▼
                            [ LLM Orchestrator ]
                                     │
                                     ▼
                       [ Evidence-Backed Response + Citations (SSE) ]
```

---

## 2. Guardrails
- If no evidence matches with high confidence, the system explicitly reports that no verified evidence exists.
- System instructions strictly isolate untrusted queries and external documentation.
