---
name: performance-audit
description: Procedures for auditing Core Web Vitals, dynamic imports, bundle size, and rendering efficiency.
---

# Skill: Performance Audit

## When to Use
Use to verify Core Web Vitals and optimize asset loading in `apps/web`.

## Required Reading
1. `docs/SOP/07-performance.md`
2. `docs/design/motion-system.md`

## Implementation Rules
- Verify that 3D canvas and MapLibre are dynamically imported with SSR disabled.
- Check bundle composition: AI drawer and heavy libraries excluded from hero bundle.
- Ensure LCP <= 2.5s, CLS <= 0.1, INP <= 200ms.
- Enforce adaptive DPR on high-density displays.

## Validation Checklist
- [ ] Initial bundle size within budget.
- [ ] No unoptimized full-resolution textures or images.
- [ ] DPR clamped at 2.0.
