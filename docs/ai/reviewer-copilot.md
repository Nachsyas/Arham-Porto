# AI Specification: Reviewer Copilot (Ask Arham AI)

- **Purpose**: Architecture and interaction design for the evidence-backed reviewer assistant.
- **Current Phase**: Phase 0 (Bootstrap Foundation).

---

## 1. Feature Identity
- **Public Name**: Ask Arham AI
- **Internal System Name**: AI Reviewer Copilot
- **Goal**: Assist recruiters, technical leads, and engineering managers in exploring verified skills, architectural decisions, and repository evidence.

---

## 2. Response Contract
Every answer provides structured citations:
- `answer`: Natural language summary grounded in evidence.
- `citations`: Array of objects (`source_type`, `title`, `url`, `confidence`).
- `actions`: Deep-link navigation within the portfolio (e.g. view project, inspect case study).
