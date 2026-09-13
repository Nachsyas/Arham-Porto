# Design Specification: Hero Portrait

**Status**: ACTIVE PRODUCTION VISUAL  
**Owner**: Nachsyas Arham Mumtaz Nashohi  
**Component**: [`apps/web/features/hero/HeroPortrait.tsx`](file:///Users/user/Documents/Arham%20Porto/apps/web/features/hero/HeroPortrait.tsx)  

---

## 1. Overview & Purpose
The Hero visual features an authentic portrait photo of **Nachsyas Arham Mumtaz Nashohi** to establish immediate credibility, engineering professionalism, and a recruiter-first personal connection. This replaces the previous 3D wireframe hologram concept.

---

## 2. Photo Source & Integrity
- **Current Canonical Asset**: [`apps/web/public/images/profile/nachsyas-arham.jpg`](file:///Users/user/Documents/Arham%20Porto/apps/web/public/images/profile/nachsyas-arham.jpg)
- **Source**: Extracted cleanly from the owner's approved curriculum vitae document.
- **Replacement Tracker**: `TODO_USER_HIGH_RES_PROFILE_PHOTO` is maintained in the profile metadata for when a studio-grade high-resolution portrait is provided.
- **Natural Color Grading**: The photo retains its natural skin tones and lighting. Cyan tinting over the face is strictly forbidden.

---

## 3. Privacy & CV Data Safeguards
- **Zero Document Exposure**: Only the cropped portrait itself is published.
- **Strict Privacy Isolation**: Private personal details present in CV documents (home address, phone number, date of birth, identity numbers) must **NEVER** be published or exposed in the client bundle.

---

## 4. Visual Styling & Frame Architecture
- **Frame Aesthetic**: Dark navy / near-black container (`bg-gradient-to-b from-surfaceElevated/70 to-surface/90`, `border-border`).
- **Restrained Accents**:
  - Subtle corner coordinate brackets (`┌ PROFILE`, `VERIFIED ┐`, `└ ID:NASHOHI`, `LIVE ┘`).
  - Minimal telemetry strip (`Software Engineer | Reviewer-Ready`).
  - Subtle grounding gradient (`bg-gradient-to-t from-canvas/80 via-transparent to-transparent`) for seamless theme integration.
- **Forbidden Treatments**: No fake face detection boxes, no biometric labels, no fake radar sweeps, and no heavy neon glows.

---

## 5. Responsive Cropping & Layout
- **Desktop (1440px)**:
  - Left column: Executive typography (`NACHSYAS ARHAM MUMTAZ NASHOHI`), role badge, positioning statement, primary CTAs (`Explore My Work`, `60-Second Quick Review`), and GitHub link.
  - Right column: Portrait card (`max-w-[420px]`, `aspect-[3/4]`). Head to upper torso framing.
- **Mobile (390px)**:
  - Single column sequence: Name -> Role -> Positioning -> CTAs -> Portrait (`max-w-[340px]`).
  - Mobile crop focuses closely on face and shoulders with zero horizontal overflow.

---

## 6. Accessibility & Motion
- **Accessible Identification**: Next.js `<Image />` component with explicit, semantic alt text:  
  `alt="Portrait of Nachsyas Arham Mumtaz Nashohi"`.
  `aria-hidden` is strictly forbidden on the primary portrait.
- **Motion**:
  - Subtle initial fade and translation: `translateY: 12px -> 0px` on page entrance.
  - Gentle scale on hover: `scale: 1.00 -> 1.02` with 700ms ease.
  - Automatically respects `(prefers-reduced-motion: reduce)`.

---

## 7. Performance & LCP Optimization
- **Next.js `<Image />`**: Configured with `priority` attribute to ensure optimal Largest Contentful Paint (LCP).
- **Responsive Delivery**: Explicit `sizes="(max-width: 640px) 90vw, (max-width: 1024px) 45vw, 420px"`.
- **Zero WebGL Overhead**: Removing Three.js returns the homepage First Load JS to baseline static weight (~159 kB).

---

## 8. High-Resolution Replacement Process
When a high-resolution portrait photograph is provided by the portfolio owner:
1. Ensure the asset is cropped to a 3:4 or 4:5 portrait orientation (recommended resolution: at least 1200x1600 px).
2. Save the optimized file as `apps/web/public/images/profile/nachsyas-arham.webp` (or `.jpg`).
3. If using `.webp`, update the `src` attribute in [`apps/web/features/hero/HeroPortrait.tsx`](file:///Users/user/Documents/Arham%20Porto/apps/web/features/hero/HeroPortrait.tsx).
4. Remove the `TODO_USER_HIGH_RES_PROFILE_PHOTO` marker from the profile documentation once verified.
5. Re-run `npm run build --workspace=apps/web` and verify that the layout and LCP metrics remain optimal.
