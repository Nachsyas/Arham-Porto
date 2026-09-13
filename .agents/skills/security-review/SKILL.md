---
name: security-review
description: Security verification for prompt injection defense, secret leakage prevention, CORS, and rate limiting.
---

# Skill: Security Review

## When to Use
Use before releasing code modifications to Go API endpoints or AI integrations.

## Required Reading
1. `docs/SOP/08-security-privacy.md`
2. `docs/ai/guardrails.md`

## Implementation Rules
- Scan for hardcoded credentials, API keys, or private tokens.
- Verify that untrusted inputs (questions, JDs) cannot alter system instructions.
- Ensure CORS middleware enforces allowed origins only.
- Validate request payload size limits and rate limiting on public endpoints.

## Validation Checklist
- [ ] No secrets committed.
- [ ] Prompt injection boundaries verified with test inputs.
- [ ] CORS allowlist strictly configured.
