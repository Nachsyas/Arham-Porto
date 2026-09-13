/**
 * Canonical Color Tokens for Arham Porto (Master Prompt Section 47)
 * Single Source of Truth for all color constants.
 */
export const colors = {
  canvas: {
    DEFAULT: "#02060B",
    soft: "#040A12",
  },
  surface: {
    DEFAULT: "#07111C",
    elevated: "#0A1825",
    strong: "#0D2030",
  },
  primary: {
    DEFAULT: "#26B8FF",
    active: "#55C8FF",
    muted: "#123D55",
  },
  hologram: {
    DEFAULT: "#61D8FF",
    soft: "#2E8FB5",
    faint: "#14384A",
  },
  text: {
    primary: "#F3F8FC",
    body: "#C7D4DD",
    muted: "#8395A3",
    mutedSoft: "#5F7180",
  },
  border: {
    DEFAULT: "#173247",
    strong: "#25516E",
  },
  status: {
    success: "#5AD6A0",
    warning: "#F0B65A",
    error: "#EF6A73",
  },
} as const;

export type ColorTokens = typeof colors;
