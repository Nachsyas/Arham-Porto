# AI Specification: Safety & Prompt Injection Guardrails

- **Purpose**: Guard against hallucinations, prompt injection, and unauthorized data disclosure.
- **Current Phase**: Phase 0 (Bootstrap Foundation).

---

## 1. Prompt Injection Defense
- Strict context boundary separation:
  - System Instructions: Immutable rules and persona.
  - Approved Evidence: Grounded retrieved facts.
  - Untrusted Input: User questions and pasted Job Descriptions.
- User input is sanitized and wrapped in unambiguous XML/JSON boundary delimiters.

---

## 2. Hallucination Guardrail
- If evidence is absent, the model MUST explicitly respond:
  *"I couldn't find verified evidence of that in Nachsyas Arham Mumtaz Nashohi's approved portfolio sources."*
- Extrapolation or speculative suitability scoring is strictly prohibited.
