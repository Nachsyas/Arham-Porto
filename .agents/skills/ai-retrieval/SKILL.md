---
name: ai-retrieval
description: Patterns for grounded vector search, similarity matching, and structured citation generation.
---

# Skill: AI Evidence Retrieval

## When to Use
Use during Phase 6 when implementing RAG search logic and citation assembly in Go.

## Required Reading
1. `docs/ai/reviewer-copilot.md`
2. `docs/ai/guardrails.md`
3. `docs/SOP/08-security-privacy.md`

## Implementation Rules
- Query pgvector index using cosine distance against indexed knowledge chunks.
- Validate similarity threshold; if confidence is low or evidence missing, explicitly trigger the missing evidence response.
- Attach structured citations (`source_type`, `title`, `url`, `confidence: "verified" | "supporting"`).
- Stream responses using Server-Sent Events (SSE).

## Validation Checklist
- [ ] Hallucination prevention active on ungrounded queries.
- [ ] Structured citations include source URLs and verified status.
- [ ] SSE streams flush tokens without intentional synthetic delays.
