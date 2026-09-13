/**
 * Canonical Spacing Scale for Arham Porto (Master Prompt Section 50)
 */
export const spacing = {
  1: "0.25rem", // 4px
  2: "0.5rem",  // 8px
  3: "0.75rem", // 12px
  4: "1rem",    // 16px
  6: "1.5rem",  // 24px
  8: "2rem",    // 32px
  12: "3rem",   // 48px
  16: "4rem",   // 64px
  24: "6rem",   // 96px
  28: "7rem",   // 112px
} as const;

export type SpacingTokens = typeof spacing;
