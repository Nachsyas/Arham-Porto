# Design Specification: Journey Map

- **Purpose**: Choreography and privacy boundaries for the interactive Indonesia Journey Map.
- **Current Phase**: Phase 0 (Bootstrap Foundation).

---

## 1. Interaction Flow
- Map Engine: MapLibre GL JS.
- Initial view: Panoramic Indonesia overview.
- Scroll-driven camera movement: Pan, zoom, and fly smoothly to active journey stops (`birthplace`, `tk`, `sd`, `smp`, `sma`, `university`, `current`).
- Global page scroll remains native; mouse wheel is never hijacked.

---

## 2. Privacy Constraints
- Never display exact residential coordinates or personal street addresses.
- Current location is restricted to city-level approximate center (`TODO_USER_CURRENT_CITY`, `TODO_USER_JOURNEY_COORDINATES`).
