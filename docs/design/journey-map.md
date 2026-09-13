# Design Specification: Interactive Academic Journey Map

- **Project**: Arham Porto
- **Portfolio Owner**: Nachsyas Arham Mumtaz Nashohi
- **Feature**: Interactive Geographic Journey Map (`features/journey/`)
- **Status**: Recertified — Phase 3 Journey Expansion Deliverable
- **Classification**: Stylized Geographic Visualization (Java Corridor, Indonesia)

---

## 1. Engine Selection: Precision / Stylized SVG Map

### 1.1 Rationale for SVG over MapLibre GL
During Phase 3 architectural evaluation and subsequent data expansion, the mapping engine remains a **Precision / Stylized SVG Geographic Map**:

1. **Focused Geographic Scope**: The narrative focuses strictly on 4 unique locations across DKI Jakarta, Central Java, and East Java. Free pan/zoom across the entire globe is unnecessary and adds cognitive friction for recruiters.
2. **Zero External Infrastructure**: Zero tile server requests, zero API keys, and zero third-party telemetry dependencies. The asset is 100% bundled locally at runtime.
3. **No WebGL / GPU Fragility**: Operates deterministically across all devices, mobile browsers, battery-saver modes, and headless test runners without WebGL context loss.
4. **Lightweight Bundle Footprint**: Avoids heavy mapping engine dependencies. The entire First Load JS remains lean at **171 kB** (only +1 kB over the initial Phase 3 baseline).
5. **Deterministic Testing & Accessibility**: SVG markers map directly to real semantic HTML `<button>` elements with full keyboard focus, visible rings, and accessible screen-reader labels.

---

## 2. Geographic Asset, Projection & License

- **Source**: Natural Earth Vector (`ne_50m_admin_0_countries` for Java, Madura, Bali; `ne_110m_admin_0_countries` for ambient Indonesia archipelago).
- **License**: Public Domain / CC0 (unrestricted commercial and non-commercial use).
- **Cartographic Disclaimer**: Explicitly labeled in the UI as a **Stylized Geographic Visualization** rather than survey-level cadastral cartography.
- **Coordinate Projection**: Coordinates from WGS84 are pre-projected into an `800x360` SVG viewBox coordinate space:
  - **Jakarta** (Residence, MI, MTs): `(172.0, 72.0)`
  - **Karanganyar** (Origin): `(478.3, 203.2)`
  - **Salatiga** (Secondary Education): `(445.3, 178.2)`
  - **Malang** (University & Current Base): `(603.4, 239.2)`

---

## 3. Core Architectural Model: 7 Story Milestones vs. 4 Unique Locations

A fundamental architectural tenet of the Journey system is the explicit separation between **Geographic Locations** and **Narrative Story Milestones**:

```
Chronological Story Milestones:
01. Origin (Karanganyar)
02. Residence / Formative Base (Jakarta)
03. Primary Education — MI Terpadu Al-Hamid (Jakarta Timur)
04. Lower Secondary Education — MTsN 30 (Jakarta Timur)
05. Secondary Education — MA Tahfizhul Qur'an As-Surkati (Salatiga)
06. University — UIN Maulana Malik Ibrahim (Malang)
07. Current Base (Malang)

Geographic Route Corridor:
Karanganyar (1) ───► Jakarta (2) ───► Salatiga (3) ───► Malang (4)
```

### 3.1 Exactly 4 Unique Geographic Pins
The interactive map renders exactly 4 geographic pins on the Java corridor:
1. `karanganyar`: Origin
2. `jakarta`: Shared marker for Residence, MI Al-Hamid, and MTsN 30
3. `salatiga`: Secondary Education (MA As-Surkati)
4. `malang`: Shared marker for University (CS) and Current Base

### 3.2 Shared Marker Behavior
- **Jakarta Shared Marker**:
  - Clicking the Jakarta marker when coming from another location selects the first Jakarta milestone (`residence-jakarta`).
  - If a Jakarta milestone is already active, clicking the marker preserves the current Jakarta milestone.
  - The marker badge dynamically reflects the active milestone: `[Residence]`, `[MI Al-Hamid]`, or `[MTsN 30]`.
  - Switching between Jakarta milestones (Residence $\leftrightarrow$ MI $\leftrightarrow$ MTs) updates the card and timeline state without replaying the full geographic route animation.
- **Malang Shared Marker**:
  - Shared between University (`university-uin-malang`) and Current Base (`current-base-malang`).
  - Dynamically updates badge: `[University]` $\leftrightarrow$ `[Current Base]`.

### 3.3 Route Rule (Zero Intra-City Hops)
- The geographic route comprises exactly **3 curved route segments**:
  $$\text{Karanganyar} \longrightarrow \text{Jakarta} \longrightarrow \text{Salatiga} \longrightarrow \text{Malang}$$
- SVG Bezier corridor:
  `M 478.3 203.2 Q 320 100 172 72 Q 310 160 445.3 178.2 Q 520 185 603.4 239.2`
- **Strictly Forbidden**:
  - Zero `Jakarta → Jakarta` hops.
  - Zero `Malang → Malang` hops.

---

## 4. Privacy & Truthfulness Protocol

### 4.1 Safe Jakarta City-Level Location
- The Residence milestone uses city-level Jakarta coordinates `[106.8272, -6.1754]` (`DKI Jakarta, Indonesia`).
- **Zero Exposure**: No home address, housing complex (e.g. Bambu Kuning), district/neighborhood, or private GPS coordinates are ever exposed publicly or committed.
- Public display is strictly confined to administrative region: `Jakarta, DKI Jakarta`.

### 4.2 Period Omission for Unconfirmed Education (MI & MTs)
- Schooling years for MI Terpadu Al-Hamid and MTsN 30 Jakarta Timur are not yet user-confirmed.
- Both records store `period: null` in `data/journey/journey.json`.
- **Honest AI & Rendering Guard**:
  - The UI gracefully omits the period metadata row completely for these stops.
  - Never renders placeholders such as `"Unknown"`, `"TBD"`, `"TODO"`, or inferred dates.

### 4.3 Birth Year & Identity Protection
- Birth year (2004) remains strictly hidden.
- Early childhood education (`journey-tk`) remains `public: false` and is excluded from the DOM.
- No personal contact numbers, residential details, or unauthorized claims are rendered.

### 4.4 Tone & Recruiter Copy Guard
- Neutral, evidence-backed copy:
  *"From Karanganyar to Jakarta, Salatiga, and Malang — a journey through formative education and computer science."*
- Zero instances of unsupported titles (e.g. "founder").

---

## 5. Interaction Model & Navigation Channels

- **Default State**: Chronological start at Origin (`origin-karanganyar`).
- **Canonical State Management**: `activeMilestoneId` is the single source of truth across SVG map pins, milestone story card, and timeline.
- **Multiple Navigation Modalities**:
  1. Map marker clicks (Karanganyar, Jakarta, Salatiga, Malang).
  2. Timeline step button clicks (01 Origin through 07 Current Base).
  3. Previous / Next pagination buttons with boundary disabling.
  4. Keyboard Left / Right arrow navigation.
  5. Native keyboard Tab + Enter/Space on any interactive control.
- **Zero Scroll Hijacking**: Document scrolling is completely native with zero wheel interception.

---

## 6. Accessibility & Motion Preferences

- **Full Screen-Reader Names**:
  - Timeline buttons include complete accessible labels via `aria-label`:
    - *"Origin in Karanganyar"*
    - *"Residence in Jakarta"*
    - *"Primary education at Madrasah Ibtidaiyah Terpadu Al-Hamid"*
    - *"Lower secondary education at Madrasah Tsanawiyah Negeri 30 Jakarta Timur"*
    - *"Secondary education at Madrasah Aliyah Tahfizhul Qur'an As-Surkati"*
    - *"Undergraduate education at Universitas Islam Negeri Maulana Malik Ibrahim Malang"*
    - *"Current base in Malang"*
- **Reduced Motion Support (`prefers-reduced-motion: reduce`)**:
  - Connecting corridor stroke renders statically without dash animation.
  - Radar ping animations (`animate-ping`) are suppressed.
  - Scale transforms are disabled.

---

## 7. Responsive Layout Specifications

- **Desktop (1440px)**:
  - 12-column grid: 7 columns for SVG Map with glowing corridor, 5 columns for Milestone Card.
  - Horizontal progress bar timeline with 7 steps and glowing cyan progression line.
- **Mobile (390px)**:
  - Vertical stack layout avoiding horizontal document overflow.
  - Vertical compact timeline with connected vertical line and step badges.
  - Salatiga marker label positioned above the pin to eliminate visual overlap.
