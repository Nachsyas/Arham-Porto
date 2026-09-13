---
name: accessibility-review
description: Protocol for auditing WCAG 2.2 AA compliance, semantic hierarchy, keyboard navigation, and screen reader announcements.
---

# Skill: Accessibility Review

## When to Use
Use before completing any UI feature milestone or during Phase 8 hardening.

## Required Reading
1. `docs/SOP/06-accessibility.md`
2. `docs/design/design-system.md`

## Implementation Rules
- Verify color contrast against tokens (`#F3F8FC` on `#02060B` > 14:1).
- Test keyboard tab navigation: no focus traps, visible focus rings.
- Ensure all interactive buttons have accessible names (`aria-label` or visible text).
- Check `prefers-reduced-motion` handling in animations and 3D scenes.

## Validation Checklist
- [ ] No contrast violations detected.
- [ ] Tab order verified across desktop and mobile sheets.
- [ ] Skip-to-content functional.
