# Design Specification: 3D Wireframe Hologram

- **Purpose**: Technical and aesthetic specification for the signature seated hologram.
- **Current Phase**: Phase 0 (Bootstrap Foundation).

---

## 1. Aesthetic Definition
- Real 3D geometric wireframe mesh (not 2D image cutout, sprite, or CSS fake).
- Seated pose with glowing vertices and semi-transparent geometric lines (`#61D8FF`).
- Projection platform and subtle scan plane.
- Rotation range: Y-axis `[-10°, +10°]`, X-axis `[-3°, +3°]`. Continuous 360° spinning is strictly forbidden.

---

## 2. Asset Pipeline
- Target asset: `apps/web/public/models/arham-wireframe.glb`.
- Placeholder tag: `TODO_USER_3D_MODEL`.
- Fallback: Static wireframe poster on devices without WebGL support or with low GPU performance.
