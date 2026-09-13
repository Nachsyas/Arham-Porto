---
name: ai-eval
description: Golden dataset testing procedures for evaluating AI Reviewer groundedness, refusal accuracy, and citation validity.
---

# Skill: AI Golden Evaluation

## When to Use
Use during Phase 7 to run regression evaluations on Ask Arham AI responses.

## Required Reading
1. `docs/ai/evals.md`
2. `docs/ai/guardrails.md`
3. `docs/SOP/05-testing-strategy.md`

## Implementation Rules
- Run automated tests against the Golden Evaluation Dataset.
- Assert negative rejection on unevidenced technologies (e.g. Kubernetes, Terraform).
- Assert full citation presence for known skills (e.g. Go, PostgreSQL, Next.js).
- Assert absence of arbitrary fit scores (e.g. "95% match").

## Validation Checklist
- [ ] All golden questions pass assertion suite.
- [ ] No hallucinated claims detected.
