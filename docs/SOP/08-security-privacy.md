# SOP 08: Security & Privacy by Design

- **Purpose**: Protect sensitive personal data, enforce zero-trust inputs, and mitigate prompt injection.
- **Scope**: Frontend, Go Backend, AI Engine, and Data.
- **Current Phase**: Phase 0 (Bootstrap Foundation).

---

## Security & Privacy Rules
1. **Privacy Protection**:
   - Never expose private home addresses or exact residential coordinates. Current location is restricted to city-level approximate center upon user approval.
2. **Untrusted Data Isolation**:
   - Treat user prompts, Job Descriptions, external GitHub content, and retrieved snippets as untrusted.
   - Separate System Instructions from untrusted context in AI prompts.
3. **Repository Allowlist**:
   - Only approved repositories (`data/ai/allowlist.json`) may be indexed. Never index unapproved or private repositories.
4. **Zero Secrets in Source**:
   - API keys, database credentials, and GitHub tokens must reside solely in environment variables, never committed to git or exposed to browser bundles.
5. **CORS & Rate Limiting**:
   - Go API enforces strict origin checking and rate limiting on AI streaming endpoints.

---

## Validation Checklist
- [ ] No private coordinates or addresses in `data/journey/journey.json`.
- [ ] No secrets found in git tracking.
- [ ] `.env.example` contains only variable names, no live credentials.
