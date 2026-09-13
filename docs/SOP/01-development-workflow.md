# SOP 01: Development Workflow

- **Purpose**: Define the phased development methodology and boundary rules for Arham Porto.
- **Scope**: All AI agents and human contributors working across the monorepo.
- **Current Phase**: Phase 0 (Bootstrap Foundation).

---

## Rules
1. **Single Source of Truth**: Follow Master Prompt v1.0.1. Decisions not covered must be flagged as `TODO_USER` or `DECISION_REQUIRED`.
2. **Phase Gating**: Never begin work on a future phase without explicit user approval.
3. **Context First**: Read `.agents/skills/read-project-context/SKILL.md` before initiating any task.
4. **No Fake Completeness**: Never fabricate personal bio, unapproved project contributions, or metrics. Use typed `null` with explicit `todo` tags.

---

## Future Implementation Boundary
- Phase 1: Reviewer-First Static Portfolio.
- Phase 2: 3D Wireframe Hologram.
- Phase 3: Journey Map.
- Phase 4–8: Go backend, Indexing, AI Copilot, Evals, Hardening.

---

## Validation Checklist
- [ ] Active phase confirmed in `docs/todo/master-todo.md`.
- [ ] No unapproved cross-phase code written.
- [ ] All unknown data explicitly marked with `TODO_USER_*`.
