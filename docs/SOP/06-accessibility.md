# SOP 06: Accessibility (a11y) Standards

- **Purpose**: Guarantee WCAG 2.2 AA compliance across all user interfaces.
- **Scope**: Frontend (`apps/web`).
- **Current Phase**: Phase 0 (Bootstrap Foundation).

---

## Core Guidelines
1. **Semantic HTML**: Proper heading hierarchy (`<h1>` through `<h3>`), `<main>`, `<nav>`, `<header>`, `<footer>`.
2. **Keyboard Navigation**: All interactive elements must be focusable with visible focus rings (`focus-visible:ring-2`).
3. **Contrast Ratio**: Minimum 4.5:1 for normal text, 3:1 for large text against dark surfaces (`#02060B` / `#07111C`).
4. **Reduced Motion**: Respect `prefers-reduced-motion: reduce` by disabling non-essential animations, continuous 3D rotation, and map transitions.
5. **Touch Targets**: Minimum 44x44px clickable area on mobile viewports.
6. **Screen Reader Alternatives**: Textual descriptions for 3D canvas and interactive map components.

---

## Validation Checklist
- [ ] Skip-to-content link present in root layout.
- [ ] Color contrast verified against token definitions.
- [ ] Tab order logical and non-trapping.
