# Design Specification: Motion System

- **Purpose**: Define animation boundaries and performance guidelines.
- **Current Phase**: Phase 0 (Bootstrap Foundation).

---

## 1. Responsibilities
- **GSAP + ScrollTrigger**: Scroll-driven section transitions, sticky 3D hologram subtle rotation, and Journey map scroll orchestration.
- **Motion**: Component-level transitions (dialogs, drawers, tabs, hover states).
- **Rule**: Never duplicate animation logic between libraries.

---

## 2. Scroll Integrity & Accessibility
- Native browser scroll must never be locked, hijacked, or accelerated.
- If `prefers-reduced-motion: reduce` is active, all dynamic scroll rotations and complex transitions must gracefully degrade to static states.
