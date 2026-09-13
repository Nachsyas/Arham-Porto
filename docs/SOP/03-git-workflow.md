# SOP 03: Git Workflow & Version Control

- **Purpose**: Maintain git hygiene and clear commit boundaries.
- **Scope**: Repository-wide version control.
- **Current Phase**: Phase 0 (Bootstrap Foundation).

---

## Rules
1. **Never Commit Secrets**: Ensure `.env`, credentials, tokens, and private documents are strictly ignored in `.gitignore`.
2. **Deterministic Locking**: Commit `package-lock.json`. CI must always execute `npm ci`.
3. **Clean Commit Scope**: Commits must reflect discrete logical steps (e.g. `feat(phase0): scaffold monorepo and content schemas`).
4. **Untracked Check**: Verify `git status` contains only intentional, tracked repository files.

---

## Validation Checklist
- [ ] No `.env` or credentials tracked.
- [ ] No `node_modules` or `.next` tracked.
- [ ] Lockfile updated and synchronized.
