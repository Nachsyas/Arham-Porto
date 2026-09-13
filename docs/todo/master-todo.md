# Master Development Roadmap & Progress Tracker

> Single source of truth for implementation progress of **Arham Porto**.

---

## Phases Overview

| Phase | Description | Status |
| :--- | :--- | :--- |
| **Phase 0** | **Bootstrap Foundation** | **Completed / Certified** |
| **Phase 1** | **Reviewer-First Static Portfolio** | **Completed / Certified** |
| **Phase 2** | **3D Wireframe Hologram** | **Completed** |
| **Phase 3** | **Journey Map** | **Awaiting User Approval** |
| **Phase 4** | **Go Backend** | Not Started |
| **Phase 5** | **AI Indexing & Retrieval** | Not Started |
| **Phase 6** | **AI Reviewer Copilot** | Not Started |
| **Phase 7** | **Integration & Polish** | Not Started |
| **Phase 8** | **Production & Launch** | Not Started |

---

## Phase 0 Checklist (Certified Baseline)
- [x] Root project folder setup (`Arham Porto`)
- [x] Monorepo architecture with npm workspaces (`apps/*`, `packages/*`)
- [x] Pinned toolchain versions (`.nvmrc`, `package.json` engines)
- [x] Root `.gitignore` and `.env.example`
- [x] Docker Compose with PostgreSQL 16 + pgvector (`pgvector/pgvector:pg16`) and DRY healthcheck
- [x] AI Agent Deployment Rules (`AGENTS.md`)
- [x] Canonical design tokens package (`packages/design-tokens`)
- [x] Content schema package with Zod runtime validation (`packages/content-schema`)
- [x] Canonical JSON content files (`data/*`) with zero-trust validation script
- [x] Portfolio data access layer (`packages/portfolio-data`)
- [x] Next.js + TypeScript strict frontend workspace (`apps/web`) with smoke test page
- [x] Go API module skeleton (`apps/api`) with clean architecture and health endpoints
- [x] GitHub Actions CI workflow skeleton (`.github/workflows/ci.yml`)
- [x] Baseline git commit: `chore: bootstrap Arham Porto foundation` (`c31b393`)
- [x] Validation suite executed and certified:
  - `npm run validate:data`: 100% PASS (8/8 schemas)
  - `npm run lint --workspace=apps/web`: 100% PASS (0 warnings/errors)
  - `npm run typecheck --workspace=apps/web`: 100% PASS (0 errors)
  - `npm run test --workspace=apps/web`: 100% PASS
  - `npm run build --workspace=apps/web`: 100% PASS (Next.js production build static generation)
  - `go test -v ./...` and `go vet ./...`: 100% PASS
  - `docker compose config`: 100% PASS
  - `git status --short`: Clean working tree

---

## Phase 1 Checklist (Completed)
- [x] Global Navigation (Navbar, identity monogram, nav links, Quick Review trigger, mobile responsive drawer)
- [x] Hero Identity (Name `NACHSYAS ARHAM MUMTAZ NASHOHI`, Software Engineer role, neutral approved positioning in data layer, tags, dual CTAs, lightweight 3D wireframe stage placeholder)
- [x] Quick Review (Drawer with verified candidate profile, verified core stack, project links, verified contact; graceful omission of unverified fields; zero `TODO_USER` leaks)
- [x] Selected Work (Data-driven cards for EduTrace, GDGOC E-Commerce, Maritime AI Dashboard, Smart Kitchen)
- [x] Project Filtering (All, AI, Full-Stack, Backend, Systems with client-side reactive filtering)
- [x] Skill -> Evidence Explorer (`SKILL -> CLAIM -> EVIDENCE` architecture without arbitrary percentage bars)
- [x] Experience (Section architecture with dignified empty state when unverified)
- [x] Education (Section architecture with dignified empty state when unverified)
- [x] Contact (Verified channels only: GitHub `https://github.com/Nachsyas`, privacy note, zero fake email)
- [x] Project Case Study Dynamic Route (`/projects/[slug]`) with SSG `generateStaticParams`
- [x] Accessibility & Responsive verification (WCAG 2.2 AA, semantic HTML, skip-to-content, keyboard navigation, mobile touch targets >= 44px)
- [x] Unit & Component test suite for all Phase 1 components (10 tests passing across 3 test suites)
- [x] Final Phase 1 Walkthrough & Stop Rule enforcement (halt before Phase 2)

---

## Phase 2 Checklist (Completed)
- [x] Evidence mapping correction: remapped `evidence-maritime-dashboard-repo` to `ai-ml`, reserved `ai-rag` with empty evidence array, validated with `npm run validate:data`
- [x] R3F & Three.js dependency integration (`three`, `@types/three`, `@react-three/fiber`, `@react-three/drei`)
- [x] Modular Hologram feature boundary in `apps/web/features/hologram/`
- [x] Canonical design tokens & color constants (`#61D8FF`, `#2E8FB5`, `#14384A`, `#26B8FF`)
- [x] Procedural Seated Human Wireframe Mesh (`SeatedDevelopmentMesh`) with zero-trust posture truthfulness (no fake biometric/facial reconstruction)
- [x] Future GLB model loading boundary (`/models/arham-wireframe.glb`) isolating `TODO_USER_3D_MODEL`
- [x] Technical ground projection platform with concentric rings (`RingGeometry`) and cardinal crosshairs
- [x] Sparse floating depth particles (45 desktop / 20 mobile) with slow vertical drift
- [x] R3F `useFrame()` render loop with smooth `THREE.MathUtils.damp()` rotation
- [x] Scroll-driven rotation constrained strictly to [-8°, +8°] Yaw and [-2°, +2°] Pitch
- [x] Subtle pointer parallax on desktop (disabled on mobile)
- [x] Lazy loading via Next.js dynamic import (`ssr: false`) with instant HTML Hero text rendering
- [x] WebGL error boundary and static SVG/2D vector fallback (`HologramFallback`)
- [x] Accessibility support with ARIA label and reduced motion handling (`prefers-reduced-motion: reduce`)
- [x] Extensible future AI/Journey state contracts (`mode`, `isAIActive`, `orientationBias`)
- [x] Comprehensive automated test suite (`apps/web/tests/phase2-hologram.test.tsx`, 6/6 tests passing, total 18/18 tests passing in `apps/web`)
- [x] Visual validation across Desktop (1440px) and Mobile (390px) viewports with zero console errors
- [x] Production build verification and bundle analysis
- [x] Stop Rule enforcement before Phase 3 (Journey Map)

