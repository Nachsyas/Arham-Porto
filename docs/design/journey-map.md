# Design Specification: Interactive Academic Journey Map

- **Project**: Arham Porto
- **Feature**: Interactive Geographic Journey Map (`features/journey/`)
- **Status**: Certified — Phase 3 Production Deliverable
- **Classification**: Stylized Geographic Visualization (Java Corridor, Indonesia)

---

## 1. Engine Selection: Precision / Stylized SVG Map

### 1.1 Rationale for SVG over MapLibre GL
During Phase 3 architectural evaluation, the mapping engine was selected as a **Precision / Stylized SVG Geographic Map** rather than MapLibre GL JS:

1. **Focused Geographic Scope**: The narrative focuses strictly on 3 unique locations across Central and East Java. Free pan/zoom across the entire globe is unnecessary and adds cognitive friction for recruiters.
2. **Zero External Infrastructure**: Zero tile server requests, zero API keys, and zero third-party telemetry dependencies. The asset is 100% bundled locally at runtime.
3. **No WebGL / GPU Fragility**: Operates deterministically across all devices, mobile browsers, battery-saver modes, and headless test runners without WebGL context loss.
4. **Lightweight Bundle Footprint**: Avoided hundreds of kilobytes of mapping engine JavaScript. The entire production bundle increase for Phase 3 was just **7 kB** First Load JS (from 163 kB baseline to 170 kB).
5. **Deterministic Testing & Accessibility**: SVG markers map directly to real HTML `<button>` elements with keyboard focus, visible rings, and screen-reader labels.

---

## 2. Geographic Asset, Projection & License

- **Source**: Natural Earth Vector (`ne_50m_admin_0_countries` for Java, Madura, Bali; `ne_110m_admin_0_countries` for ambient Indonesia archipelago).
- **License**: Public Domain / CC0 (unrestricted commercial and non-commercial use).
- **Cartographic Disclaimer**: Explicitly labeled in the UI as a **Stylized Geographic Visualization** rather than survey-level cadastral cartography.
- **Coordinate Projection**: Coordinates from WGS84 are pre-projected into an `800x360` SVG viewBox coordinate space:
  - **Karanganyar** (Origin): `(478.3, 203.2)`
  - **Salatiga** (Secondary Education): `(445.3, 178.2)`
  - **Malang** (University & Current Base): `(603.4, 239.2)`

---

## 3. Core Architectural Model: 3 Unique Locations != 4 Milestones

A fundamental requirement of the Journey architecture is the distinction between geographic points and narrative milestones:

```
Canonical Journey Stops:
1. Origin (Karanganyar)
2. Madrasah Aliyah Tahfizhul Qur'an As-Surkati (Salatiga)
3. Universitas Islam Negeri Maulana Malik Ibrahim Malang (Malang)
4. Current Engineering Base (Malang)

Geographic Locations:
Karanganyar (1) ───► Salatiga (2) ───► Malang (3)
```

- **Exactly 3 Geographic Pins**: Karanganyar, Salatiga, Malang.
- **Shared Malang Location**: Both Higher Education (Computer Science) and Current Base reference the same Malang geographic pin.
- **No Malang → Malang Route Segment**: The SVG corridor connects Karanganyar $\rightarrow$ Salatiga $\rightarrow$ Malang as a smooth quadratic curve (`M 478.3 203.2 Q 460 185 445.3 178.2 Q 520 180 603.4 239.2`). There is zero redundant Malang $\rightarrow$ Malang segment.
- **Shared Location State Transition**: Transitioning between University and Current Base preserves geographic focus on Malang, subtly updating the marker badge (`[University]` $\leftrightarrow$ `[Current Base]`) while updating the storytelling card and timeline without replaying route animations.

---

## 4. Privacy & Truthfulness Protocol

1. **Public vs. Non-Public Filtering**:
   - `data/journey/journey.json` maintains the full canonical schema.
   - Early education records (`journey-tk`, `journey-sd`, `journey-smp`) have `public: false` and are strictly filtered out before rendering (`filterPublicMilestones`). They are never injected into the DOM.
2. **Birth Year & Personal Data Guard**:
   - No birth year (e.g. 2004) or birth date is rendered.
   - No residential street address, accommodation, or telephone number is present.
   - Strictly verified by automated unit tests asserting absence of `2004`, `TODO_USER`, etc.
3. **No Raw Coordinates Telemetry**:
   - Recruiters see clear administrative titles: `Karanganyar, Jawa Tengah`, `Salatiga, Jawa Tengah`, `Malang, Jawa Timur`.
   - Decorative UI badge uses neutral geographic text: `JAVA // INDONESIA`.

---

## 5. Interaction Model & Native Scrolling

- **Default Active Milestone**: Origin — Karanganyar (`origin-karanganyar`). The journey unfolds chronologically.
- **Single Source of Active State**: `activeMilestoneId` is canonical across map pins, storytelling card, and timeline.
- **Multiple Navigation Channels**:
  1. Map marker clicks (Karanganyar, Salatiga, Malang toggle).
  2. Timeline step button clicks (01 Origin, 02 MA, 03 University, 04 Current).
  3. Previous / Next buttons with boundary disabling.
  4. Keyboard Arrow keys (Left/Right) for sequential navigation.
  5. Native keyboard Enter/Space on any marker or timeline button.
- **Zero Scroll Hijacking**: Document scrolling is 100% native. Viewport is never locked; mouse wheel is never intercepted.

---

## 6. Accessibility & Motion Preferences

- **Semantic HTML Markers**: Real `<button type="button">` elements positioned over the SVG canvas with `aria-pressed`, `aria-label`, visible cyan focus rings, and touch-friendly tap targets.
- **Reduced Motion (`prefers-reduced-motion: reduce`)**:
  - Route line renders immediately (`strokeDasharray: none`, `strokeDashoffset: 0`).
  - Marker radar ping animation (`animate-ping`) is hidden.
  - Scale transitions are disabled.
  - Full interaction remains completely operational.

---

## 7. Responsive Layout Specifications

- **Desktop (1440px)**:
  - 12-column grid: 7 columns for SVG Map, 5 columns for Milestone Card.
  - Full-width interactive horizontal timeline below with step numbers, titles, and active glow fill.
- **Mobile (390px)**:
  - Vertical stack: Header $\rightarrow$ SVG Map $\rightarrow$ Milestone Card $\rightarrow$ Timeline.
  - Salatiga pin label positioned above the pin to prevent any text collision with Karanganyar on small viewports.
  - Timeline adapts with concise labels (`Origin`, `MA`, `University`, `Current`).
  - Zero horizontal overflow or page clipping.
