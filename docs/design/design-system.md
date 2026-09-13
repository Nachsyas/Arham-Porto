# Design Specification: Design System

- **Purpose**: Document color palettes, typography, spacing, radii, and glow rules.
- **Current Phase**: Phase 0 (Bootstrap Foundation).

---

## 1. Palette & Theme Tokens
All visual tokens derive directly from `packages/design-tokens`:
- **Canvas**: `#02060B` (Canvas Soft: `#040A12`)
- **Surface**: `#07111C` (Elevated: `#0A1825`, Strong: `#0D2030`)
- **Primary Cyan**: `#26B8FF` (Active: `#55C8FF`, Muted: `#123D55`)
- **Hologram**: `#61D8FF` (Soft: `#2E8FB5`, Faint: `#14384A`)
- **Text**: Primary `#F3F8FC`, Body `#C7D4DD`, Muted `#8395A3`
- **Border**: `#173247` (Strong: `#25516E`)
- **Status**: Success `#5AD6A0`, Warning `#F0B65A`, Error `#EF6A73`

---

## 2. Restrictive Glow Policy
- Glow is strictly an active accent indicator (hologram mesh, active map marker, focused action, AI active panel, selected evidence).
- Ubiquitous card or text glows are prohibited to avoid gaming/sci-fi visual degradation.
