# Arham Porto

> Production-quality, AI-powered personal engineering portfolio platform for **Nachsyas Arham Mumtaz Nashohi** (Software Engineer).

---

## 1. Overview

**Arham Porto** is an evidence-first technical platform designed for engineering recruiters, hiring managers, and technical collaborators. Rather than an ordinary static CV or animation showcase, it functions as an interactive engineering case-study platform, skill explorer, and AI-assisted portfolio reviewer.

- **Portfolio Owner**: Nachsyas Arham Mumtaz Nashohi
- **Primary Role**: Software Engineer
- **AI Feature**: Ask Arham AI (Internal: AI Reviewer Copilot)
- **Single Source of Truth**: Master Prompt (v1.0.1)

---

## 2. Shell Navigation

Because the root folder name contains a space, always quote the directory path in shell commands:

```bash
cd "Arham Porto"
```

---

## 3. Technology Stack

- **Frontend**: Next.js (App Router), TypeScript (Strict Mode), Tailwind CSS
- **3D & Visual**: React Three Fiber, Three.js (Seated wireframe hologram)
- **Scroll Orchestration**: GSAP ScrollTrigger
- **Interactive Map**: MapLibre GL JS (Indonesia journey exploration)
- **Backend**: Go (1.22+ modern standard library `net/http`, Clean Architecture)
- **Database**: PostgreSQL 16 + pgvector (for knowledge retrieval in later phases)
- **AI Engine**: Evidence Retrieval + Provider-Agnostic LLM / Embeddings + SSE Streaming
- **Testing**: Vitest, React Testing Library, Go `testing` & `httptest`

---

## 4. Repository Structure

```
Arham Porto/
├── apps/
│   ├── web/                     # Next.js App Router frontend application
│   └── api/                     # Go clean-architecture REST API
├── packages/
│   ├── content-schema/          # TypeScript types & Zod runtime schemas
│   ├── portfolio-data/          # Data access loader with Zod validation
│   └── design-tokens/           # Canonical visual design tokens
├── data/                        # Canonical structured portfolio content (JSON)
├── docs/                        # Architecture, SOPs, specs, and master todo
├── .agents/skills/              # 12 Core Agent Skills
├── .github/workflows/           # CI/CD pipelines
├── AGENTS.md                    # Operational AI rules & boundaries
├── docker-compose.yml           # PostgreSQL + pgvector development service
├── .env.example                 # Environment configuration template
└── README.md                    # Root project guide
```

---

## 5. Development Roadmap

- **Phase 0 — Bootstrap Foundation**: Monorepo workspace, content schemas, tokens, clean architecture skeletons, documentation, skills *(Current)*
- **Phase 1 — Reviewer-First Static Portfolio**: Hero, Quick Review, Projects, Skills, Experience, Education, Contact
- **Phase 2 — 3D Wireframe Hologram**: Seated human wireframe with subtle scroll rotation
- **Phase 3 — Journey Map**: Interactive Indonesia map and storytelling cards
- **Phase 4 — Go Backend**: REST API, Clean Architecture, health checks
- **Phase 5 — Knowledge Indexer**: GitHub repository indexing pipeline, embeddings, pgvector
- **Phase 6 — Ask Arham AI**: Evidence retrieval, RAG, citations, SSE streaming
- **Phase 7 — Reviewer Intelligence**: 60-Second Brief, Role Explorer, Job Description Explorer
- **Phase 8 — Production Hardening**: Audits, accessibility, performance, security, SEO

---

## 6. Local Quickstart (Phase 0)

### Prerequisites
- Node.js >= 20.0.0 (pinned in `.nvmrc`)
- npm >= 10.0.0
- Go >= 1.22
- Docker & Docker Compose (optional in Phase 0)

### Commands
```bash
# Install monorepo dependencies
npm install

# Validate all canonical portfolio data against Zod schemas
npm run validate:data

# Test and build frontend smoke page
npm run test:web
npm run build:web

# Test and verify Go API skeleton
cd apps/api && go test ./... && go vet ./...
```
