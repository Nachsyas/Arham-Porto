/**
 * Canonical Typography Tokens for Arham Porto (Master Prompt Section 49)
 */
export const typography = {
  fonts: {
    display: "var(--font-space-grotesk), 'Space Grotesk', system-ui, sans-serif",
    body: "var(--font-geist-sans), 'Geist Sans', 'Inter', system-ui, sans-serif",
    code: "var(--font-jetbrains-mono), 'JetBrains Mono', monospace",
  },
  fontSize: {
    xs: ["0.75rem", { lineHeight: "1rem" }],
    sm: ["0.875rem", { lineHeight: "1.25rem" }],
    base: ["1rem", { lineHeight: "1.5rem" }],
    lg: ["1.125rem", { lineHeight: "1.75rem" }],
    xl: ["1.25rem", { lineHeight: "1.75rem" }],
    "2xl": ["1.5rem", { lineHeight: "2rem" }],
    "3xl": ["1.875rem", { lineHeight: "2.25rem" }],
    "4xl": ["2.25rem", { lineHeight: "2.5rem" }],
    "5xl": ["3rem", { lineHeight: "1.16" }],
    "6xl": ["3.75rem", { lineHeight: "1.1" }],
  },
} as const;

export type TypographyTokens = typeof typography;
