---
name: implement-frontend-section
description: Guidelines for implementing accessible, high-performance UI sections using Next.js App Router, React, and shared design tokens.
---

# Skill: Implement Frontend Section

## When to Use
Use when developing or refactoring UI components or feature sections in `apps/web`.

## Required Reading
1. `docs/design/design-system.md`
2. `docs/SOP/06-accessibility.md`
3. `docs/SOP/07-performance.md`

## Implementation Rules
- Consume all visual styling from `packages/design-tokens` (Tailwind classes mapped to canonical tokens).
- Maintain strict TypeScript (`strict: true`, no `any`).
- Ensure WCAG 2.2 AA accessibility (visible focus rings, contrast, keyboard navigability, semantic tags).
- Separate UI components, feature logic, and server data fetching cleanly.

## Validation Checklist
- [ ] TypeScript typecheck passes (`npm run typecheck --workspace=apps/web`).
- [ ] Visual styling matches canonical tokens.
- [ ] Component is fully accessible via keyboard.
