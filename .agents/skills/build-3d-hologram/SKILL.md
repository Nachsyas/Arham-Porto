---
name: build-3d-hologram
description: Implementation guidelines for the Three.js and React Three Fiber seated wireframe hologram scene.
---

# Skill: Build 3D Hologram

## When to Use
Use during Phase 2 when developing or modifying the 3D hologram component.

## Required Reading
1. `docs/design/hologram-3d.md`
2. `docs/design/motion-system.md`
3. `docs/SOP/07-performance.md`

## Implementation Rules
- Asset path: `apps/web/public/models/arham-wireframe.glb`. Tagged with `TODO_USER_3D_MODEL` if using temporary placeholder.
- Seated wireframe human mesh style with glowing cyan vertices (`#61D8FF`).
- Rotation range strictly constrained to Y: `[-10°, +10°]`, X: `[-3°, +3°]`. Continuous 360° spinning is prohibited.
- Disable rotation and complex shaders when `prefers-reduced-motion` is active. Provide static poster fallback for WebGL failure.

## Validation Checklist
- [ ] No WebGL context crashes on low-end devices.
- [ ] Reduced motion check tested and verified.
- [ ] Scroll rotation smoothly interpolated and bounded.
