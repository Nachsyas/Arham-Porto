---
name: build-journey-map
description: MapLibre GL JS integration and camera choreography guidelines for the interactive Indonesia Journey Map.
---

# Skill: Build Journey Map

## When to Use
Use during Phase 3 when implementing or updating the interactive Journey section.

## Required Reading
1. `docs/design/journey-map.md`
2. `docs/content/privacy-classification.md`
3. `docs/SOP/06-accessibility.md`

## Implementation Rules
- Start with Indonesia overview; smoothly pan and zoom to active stops.
- Never hijack or lock browser mouse wheel or global scrolling.
- Privacy compliance: never use exact home coordinates. Use approved city-center approximate coordinates (`TODO_USER_JOURNEY_COORDINATES`).
- Provide fully accessible text/timeline fallback for screen readers and failed WebGL/map tile loads.

## Validation Checklist
- [ ] Global page scroll operates naturally without trapping.
- [ ] Timeline fallback renders when map fails to initialize.
- [ ] No private coordinates exposed in client payloads.
