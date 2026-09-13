# SOP 04: Code Review Protocol

- **Purpose**: Provide structured criteria for reviewing pull requests and phase milestones.
- **Scope**: All code modifications across the portfolio platform.
- **Current Phase**: Phase 0 (Bootstrap Foundation).

---

## Review Gates
1. **Architectural Alignment**: Does the code adhere to the active phase without premature feature bloat?
2. **Design System Adherence**: Are colors, typography, and spacing consumed from `packages/design-tokens`?
3. **Evidence Integrity**: Are skills backed by verified sources? Are arbitrary percentages absent?
4. **Security & Privacy**: Are private coordinates or home addresses protected? Is input sanitized?
5. **Test Coverage**: Are tests added for new units and handlers?

---

## Validation Checklist
- [ ] Code passes automated CI (`npm ci`, lint, typecheck, test, `go test`, `go vet`).
- [ ] No regression introduced to existing components.
