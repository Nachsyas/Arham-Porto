# Architecture: Deployment & Containerization

- **Purpose**: Outline local and production container deployment strategies.
- **Current Phase**: Phase 0 (Bootstrap Foundation).

---

## 1. Local Infrastructure
- PostgreSQL 16 + pgvector container (`pgvector/pgvector:pg16`) via `docker-compose.yml`.
- Next.js development server running on `http://localhost:3000`.
- Go backend server running on `http://localhost:8080`.

---

## 2. Production Topology
- Web frontend deployable to edge/containerized Node.js environments.
- Backend deployable as a lightweight, single-binary container.
- Managed PostgreSQL with pgvector extension enabled.
