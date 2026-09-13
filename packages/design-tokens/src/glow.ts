/**
 * Canonical Glow Tokens for Arham Porto (Master Prompt Section 52)
 * Restricted to interactive focus, hologram visual, active markers, and AI states.
 * Ubiquitous card or text glows are strictly prohibited.
 */
export const glow = {
  hologram: "0 0 25px rgba(97, 216, 255, 0.4)",
  markerActive: "0 0 15px rgba(38, 184, 255, 0.6)",
  focusButton: "0 0 12px rgba(85, 200, 255, 0.5)",
  aiActive: "0 0 20px rgba(97, 216, 255, 0.35)",
  evidenceSelected: "0 0 10px rgba(38, 184, 255, 0.4)",
} as const;

export type GlowTokens = typeof glow;
