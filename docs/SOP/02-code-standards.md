# SOP 02: Code Standards & Clean Architecture

- **Purpose**: Enforce strict TypeScript typing and Go clean architecture principles.
- **Scope**: Frontend (`apps/web`), Backend (`apps/api`), and Shared Packages (`packages/*`).
- **Current Phase**: Phase 0 (Bootstrap Foundation).

---

## Rules
1. **Strict TypeScript**: `noImplicitAny: true`, strict null checks. `any` is strictly prohibited.
2. **Runtime Schema Validation**: Use Zod schemas in `packages/content-schema` for all JSON and API payloads.
3. **Pure Go Domain**: `apps/api/internal/domain/` must contain only standard library types. Zero database or delivery imports.
4. **Dependency Direction**: Delivery -> UseCase -> Domain. Interfaces defined in Domain or UseCase, implemented in Infrastructure.
5. **Safe Self-Healing**: When checks fail, fix the root cause with the smallest justified patch. Never weaken validation or remove tests.

---

## Future Implementation Boundary
- Concrete pgx queries and HTTP handlers arrive in Phase 4.
- AI provider SDK adapters arrive in Phase 6.

---

## Validation Checklist
- [ ] TypeScript typecheck passes with `--noEmit`.
- [ ] Go code passes `go vet ./...`.
- [ ] No `any` types present in TypeScript code.
- [ ] Domain packages contain zero external dependencies.
