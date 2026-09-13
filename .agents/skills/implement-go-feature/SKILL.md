---
name: implement-go-feature
description: Clean architecture patterns and standards for implementing Go backend endpoints, usecases, and domain models.
---

# Skill: Implement Go Feature

## When to Use
Use when implementing REST endpoints, business interactors, or data persistence in `apps/api`.

## Required Reading
1. `docs/architecture/backend-clean-architecture.md`
2. `docs/SOP/02-code-standards.md`
3. `docs/SOP/08-security-privacy.md`

## Implementation Rules
- Domain entities in `internal/domain/` must use Go standard library only.
- Routing uses modern Go 1.22+ `net/http` standard library. Frameworks like Gin/Echo/Fiber are prohibited.
- Implement strict input validation and sanitization before passing parameters to use cases.
- Write unit tests using standard `testing` and `httptest`.

## Validation Checklist
- [ ] `go test ./...` passes in `apps/api`.
- [ ] `go vet ./...` passes with zero issues.
- [ ] Domain layer has no database or transport imports.
