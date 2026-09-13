# SOP 05: Testing Strategy

- **Purpose**: Establish automated testing requirements across unit, integration, end-to-end, and AI evaluation layers.
- **Scope**: Frontend, Backend, Data Schema, and AI.
- **Current Phase**: Phase 0 (Bootstrap Foundation).

---

## Testing Matrix
1. **Content Schema & Data**:
   - Runtime Zod validation testing via `npm run validate:data`.
2. **Frontend**:
   - Unit/Component: Vitest + React Testing Library (`npm run test --workspace=apps/web`).
   - E2E (Phase 8): Playwright for critical flows (homepage, quick review, journey, drawer).
3. **Backend**:
   - Unit/Integration: Go `testing` + `httptest` (`go test ./...` in `apps/api`).
4. **AI (Phase 7)**:
   - Golden Dataset evaluation tests for hallucination detection and citation grounding.

---

## Validation Checklist
- [ ] `npm run validate:data` passes with 0 errors.
- [ ] `npm run test --workspace=apps/web` passes.
- [ ] `go test ./...` passes in `apps/api`.
