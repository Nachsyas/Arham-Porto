/**
 * Canonical Color Tokens for Arham Porto (Master Prompt Section 47)
 * Single Source of Truth for all color constants.
 */
export const colors = {
  canvas: {
    DEFAULT: "#F6FAFD",
    soft: "#EEF6FB",
  },
  surface: {
    DEFAULT: "#FFFFFF",
    elevated: "#F9FCFE",
    strong: "#EAF4F9",
  },
  primary: {
    DEFAULT: "#0799E6",
    active: "#007CC3",
    muted: "#E1F4FD",
  },
  hologram: {
    DEFAULT: "#0799E6",
    soft: "#38BDF8",
    faint: "#E0F2FE",
  },
  text: {
    primary: "#0B1F2A",
    body: "#40515C",
    muted: "#71838E",
    mutedSoft: "#96A5AE",
  },
  border: {
    DEFAULT: "#D8E6EE",
    strong: "#BED3DF",
  },
  status: {
    success: "#2E9D70",
    warning: "#B87918",
    error: "#D94C57",
  },
} as const;

export type ColorTokens = typeof colors;

