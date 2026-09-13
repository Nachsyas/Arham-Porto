# Architecture: Frontend Architecture

- **Purpose**: Outline Next.js App Router structure, component segregation, and state management.
- **Current Phase**: Phase 0 (Bootstrap Foundation).

---

## 1. Core Principles
- **Reviewer-First**: Immediate clarity within 5–30 seconds. No critical information hidden behind complex interactions.
- **Component Segregation**:
  - `app/`: Next.js routing and server page wrappers.
  - `features/`: Domain-specific UI features (`hero`, `quick-review`, `projects`, `journey`, `ai`).
  - `components/`: Atomic, reusable design components (buttons, cards, badges, modals).
  - `hooks/`: Reusable hooks for media queries, scroll progress, reduced motion.
  - `lib/`: Utility functions and formatting.

---

## 2. Animation & 3D Boundaries
- **GSAP ScrollTrigger**: Manages scroll orchestration, sticky hologram rotation, and section pinning.
- **Motion**: Manages local component transitions (drawers, modals, tab switching).
- **React Three Fiber (R3F)**: Renders the seated human wireframe hologram.

---

## 3. Server vs Client Boundaries
- `packages/portfolio-data` is consumed on the server (Server Components).
- Interactive 3D and MapLibre components are client components marked with `'use client'` and dynamically imported.
