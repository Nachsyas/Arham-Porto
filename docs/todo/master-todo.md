# Master Development Roadmap & Progress Tracker

> Single source of truth for implementation progress of **Arham Porto**.

---

## Phases Overview

| Phase | Description | Status |
| :--- | :--- | :--- |
| **Phase 0** | **Bootstrap Foundation** | **CERTIFIED — Foundation** |
| **Phase 1** | **Reviewer-First Static Portfolio** | **CERTIFIED — Reviewer-First Portfolio** |
| **Phase 2** | **Hero Visual: Authentic Portrait** | **CERTIFIED — Authentic Portrait Hero** *(3D Hologram: CANCELLED / SUPERSEDED)* |
| **Phase 3** | **Journey Map** | **AWAITING USER-APPROVED JOURNEY DATA** |
| **Phase 4** | **Go Backend** | NOT STARTED |
| **Phase 5** | **AI Indexing & Retrieval** | NOT STARTED |
| **Phase 6** | **AI Reviewer Copilot** | NOT STARTED |
| **Phase 7** | **Integration & Polish** | NOT STARTED |
| **Phase 8** | **Production & Launch** | NOT STARTED |

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

## Phase 2 Checklist — Hero Visual: Authentic Portrait (Completed)
> **Product Decision Note**: The 3D Wireframe Hologram concept has been cancelled and superseded by an authentic recruiter-first portrait photo of Nachsyas Arham Mumtaz Nashohi. All WebGL/Three.js dependencies and features have been removed from production.

- [x] **Hologram Cancellation & Cleanup**:
  - Cancelled 3D wireframe hologram and procedural mesh (`features/hologram/` deleted from production).
  - Uninstalled Three.js / R3F production packages (`three`, `@types/three`, `@react-three/fiber`, `@react-three/drei`).
  - Archived Phase 2 hologram specification to `docs/archive/hologram-3d-cancelled.md` (and updated `docs/design/hologram-3d.md` with cancellation notice).
  - Purged all public holographic terminology (`NODE:AP-01`, `SYS:MATRIX`, `MODE:2D_VECTOR`, "seated development mesh").
- [x] **Portrait Asset Extraction & Privacy Guard**:
  - Extracted authentic portrait of Nachsyas Arham Mumtaz Nashohi from user CV into `apps/web/public/images/profile/nachsyas-arham.jpg`.
  - Zero-Trust Privacy: Absolutely no personal details (full address, phone number, date of birth, or signature) extracted or published.
  - Tagged temporary extraction with `TODO_USER_HIGH_RES_PROFILE_PHOTO` for future high-resolution replacement.
- [x] **Hero Portrait Component (`apps/web/features/hero/HeroPortrait.tsx`)**:
  - Built with Next.js `<Image />` (`priority`, `fill`, `sizes`, `object-fit: cover`).
  - Semantic, meaningful alt text (`"Portrait of Nachsyas Arham Mumtaz Nashohi"`).
  - Sleek dark technical frame with `#02060B` canvas, `#07111C` surface, `#173247` border, and subtle `#26B8FF` cyan accents.
  - Technical corner markers (`┌ PROFILE`, `┐`, `└`, `┘`) and clean identity strip (`NACHSYAS ARHAM` / `SOFTWARE ENGINEER`).
  - Graceful image fallback state with user icon and technical prompt if image fails to load.
  - Restrained entrance motion respecting `(prefers-reduced-motion: reduce)`.
- [x] **Hero Section Integration (`apps/web/features/hero/HeroSection.tsx`)**:
  - Desktop 2-column layout (Left: Identity, positioning, CTAs, links; Right: Portrait frame).
  - Mobile responsive single-column layout (Name -> Role -> Positioning -> CTAs -> Portrait -> Quick Review drawer).
  - Preserved canonical design tokens and existing Quick Review drawer functionality.
- [x] **Testing & Verification**:
  - Replaced hologram tests with comprehensive `apps/web/tests/hero-portrait.test.tsx` (5 tests covering image render, alt text, role pill, corner accents, and fallback behavior).
  - 100% PASS across full unit test suite (17/17 tests passing in `apps/web`).
  - Verified zero Three.js/R3F imports in production bundle.
  - Performance: Lightweight First Load JS without WebGL overhead.
- [x] **Documentation**: Created `docs/design/hero-portrait.md` and updated `docs/todo/master-todo.md`.
- [x] **Stop Rule Enforcement**: Halt execution upon completion; await explicit user approval before Phase 3 (Journey Map).

