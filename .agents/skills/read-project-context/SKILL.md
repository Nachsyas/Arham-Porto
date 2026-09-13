---
name: read-project-context
description: Mandatory pre-execution skill for initializing context, verifying the active phase, and reviewing relevant SOPs before executing tasks.
---

# Skill: Read Project Context & SOPs

## When to Use
Execute or read this skill **MANDATORILY** prior to starting any development or refactoring task in the Arham Porto repository.

## Required Reading
1. `AGENTS.md`: Review deployment constraints, strict type-safety, and honest AI principles.
2. `docs/todo/master-todo.md`: Verify the currently approved active phase.
3. `docs/SOP/`: Read the specific SOP relevant to your task (e.g. `01-development-workflow.md`, `02-code-standards.md`, etc.).

## Implementation Rules
- Never proceed if the task belongs to a future unapproved phase.
- Do not fabricate missing user data; use `TODO_USER_*` markers.
- Confirm that boundaries between packages, client/server, and architecture layers are respected.

## Validation Checklist
- [ ] Active phase confirmed.
- [ ] Relevant SOP read and understood.
- [ ] No out-of-scope files touched.
