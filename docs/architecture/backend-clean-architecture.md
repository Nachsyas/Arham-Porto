# Architecture: Backend Clean Architecture

- **Purpose**: Specify the clean architecture layers for the Go backend service (`apps/api`).
- **Current Phase**: Phase 0 (Bootstrap Foundation).

---

## 1. Clean Architecture Layers

```
Delivery (net/http handlers, routing, middleware)
   │
   ▼
UseCase (Business interactors & orchestration)
   │
   ▼
Domain (Pure business entities & repository interfaces)
   ▲
   │
Infrastructure / Repository (pgx, pgvector, GitHub API, LLM Providers)
```

---

## 2. Purity of the Domain Layer
- `apps/api/internal/domain/` has zero third-party dependencies.
- It defines entities (`Profile`, `Project`, `Skill`, `Evidence`, `KnowledgeChunk`) and contract interfaces (`ProfileRepository`, `KnowledgeRepository`).
- Database adapters implement these interfaces in `internal/repository/postgres/`.

---

## 3. HTTP Delivery
- Standard Go 1.22+ `net/http` router with method matching (`GET /healthz`, `GET /readyz`).
- Heavy frameworks (Gin, Echo, Fiber) are prohibited.
