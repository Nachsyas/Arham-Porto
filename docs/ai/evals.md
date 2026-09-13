# AI Specification: Golden Evaluation Dataset

- **Purpose**: Define standard evaluation test cases for AI Reviewer Copilot.
- **Current Phase**: Phase 0 (Bootstrap Foundation).

---

## 1. Golden Evaluation Questions (Phase 7 Target)
1. *"What does Arham specialize in?"* -> Must cite confirmed Software Engineer role and core focus areas.
2. *"Show evidence of Go usage."* -> Must cite backend projects and repository evidence.
3. *"Which projects involve AI?"* -> Must reference Maritime AI Dashboard / EduTrace.
4. *"Does Arham use Kubernetes?"* -> Must state no verified evidence was found.
5. *"Summarize Nachsyas Arham Mumtaz Nashohi in 60 seconds."* -> Must generate concise, grounded brief.
6. *"Explain EduTrace architecture."* -> Grounded in approved architectural docs.

---

## 2. Evaluation Metrics
- Groundedness (Citations must match retrieved content).
- Faithfulness (Zero invented claims).
- Negative Rejection (Proper refusal when evidence is missing).
