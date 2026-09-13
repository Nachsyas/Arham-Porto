# SOP 07: Performance Engineering

- **Purpose**: Ensure sub-second interactivity and optimal Core Web Vitals.
- **Scope**: Frontend (`apps/web`) and Backend (`apps/api`).
- **Current Phase**: Phase 0 (Bootstrap Foundation).

---

## Core Targets
1. **Core Web Vitals**:
   - Largest Contentful Paint (LCP): <= 2.5s.
   - Cumulative Layout Shift (CLS): <= 0.1.
   - Interaction to Next Paint (INP): <= 200ms.
2. **Asset Loading Strategy**:
   - 3D Hologram (`apps/web/public/models/`): Lazy-load via dynamic import (`next/dynamic` with `ssr: false`).
   - Map Engine (MapLibre): Lazy-load when viewport approaches Journey section.
   - AI Assistant: Excluded from initial hero bundle; loaded on user drawer trigger.
3. **Adaptive DPR**: Cap canvas Device Pixel Ratio at `Math.min(window.devicePixelRatio, 2)` to protect GPU performance on mobile devices.

---

## Validation Checklist
- [ ] No heavy 3D or map packages bundled into initial page load.
- [ ] Images optimized with `next/image` when implemented in Phase 1+.
