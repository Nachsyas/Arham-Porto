# Testing Specification: Quality Plan

- **Purpose**: Detail automated testing suites across monorepo layers.
- **Current Phase**: Phase 0 (Bootstrap Foundation).

---

## 1. Test Layers
- **Schema & Data Tests**: Validated at runtime using Zod via `npm run validate:data`.
- **Frontend Unit Tests**: Component behavior, rendering, accessibility, and smoke verification via Vitest and React Testing Library.
- **Backend Tests**: Go unit tests for domain logic and HTTP handlers via `httptest` (`go test -v ./...`).
- **End-to-End Tests (Phase 8)**: Playwright testing for critical user flows.
- **AI Evaluation (Phase 7)**: Golden test evaluation for citation fidelity and groundedness.
