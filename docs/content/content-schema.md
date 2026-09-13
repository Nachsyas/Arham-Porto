# Content Specification: Schema Model

- **Purpose**: Document the data models and schemas used for portfolio content.
- **Current Phase**: Phase 0 (Bootstrap Foundation).

---

## 1. Schemas in `packages/content-schema`
- `Profile`: Identity, contact links, positioning.
- `Project`: Title, slug, summary, problem, solution, contributions, tech stack, evidence IDs.
- `Skill`: Name, category, claim, evidence IDs (no arbitrary percentages).
- `Evidence`: Type, source URL/path, summary, verification status.
- `ExperienceItem`: Role, company, period, description, skills.
- `EducationItem`: Institution, degree, field, period, description.
- `JourneyStop`: Category, title, institution, city, region, country, coordinates, period.
- `GitHubAllowlist`: List of permitted repositories for AI indexing.
