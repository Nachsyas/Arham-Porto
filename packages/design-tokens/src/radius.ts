/**
 * Canonical Border Radius Tokens for Arham Porto (Master Prompt Section 51)
 */
export const radius = {
  none: "0px",
  sm: "8px",
  md: "12px",
  card: "16px",
  lg: "20px",
  xl: "24px",
  full: "9999px",
} as const;

export type RadiusTokens = typeof radius;
