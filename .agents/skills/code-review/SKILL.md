---
name: code-review
description: Comprehensive review protocol for architecture compliance, regression prevention, and strict type safety.
---

# Skill: Code Review

## When to Use
Use when completing feature tasks or preparing a phase completion report.

## Required Reading
1. `docs/SOP/04-code-review.md`
2. `docs/SOP/02-code-standards.md`
3. `AGENTS.md`

## Implementation Rules
- Verify end-to-end type safety: no `any` in TypeScript; pure Go standard library in domain layer.
- Confirm zero regression on existing features.
- Ensure safe self-healing: never weaken tests or broaden scope to pass checks.
- Verify that Stop Rule is respected upon phase completion.

## Validation Checklist
- [ ] All linters and typechecks pass with zero warnings/errors.
- [ ] All unit and integration tests pass.
- [ ] Phase boundaries strictly preserved.
