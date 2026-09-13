import type { Config } from "tailwindcss";
import { colors, radius, spacing } from "arham-porto-tokens";

const config: Config = {
  content: [
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
    "./features/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        canvas: {
          DEFAULT: colors.canvas.DEFAULT,
          soft: colors.canvas.soft,
        },
        surface: {
          DEFAULT: colors.surface.DEFAULT,
          elevated: colors.surface.elevated,
          strong: colors.surface.strong,
        },
        primary: {
          DEFAULT: colors.primary.DEFAULT,
          active: colors.primary.active,
          muted: colors.primary.muted,
        },
        hologram: {
          DEFAULT: colors.hologram.DEFAULT,
          soft: colors.hologram.soft,
          faint: colors.hologram.faint,
        },
        themeText: {
          primary: colors.text.primary,
          body: colors.text.body,
          muted: colors.text.muted,
          mutedSoft: colors.text.mutedSoft,
        },
        border: {
          DEFAULT: colors.border.DEFAULT,
          strong: colors.border.strong,
        },
        status: {
          success: colors.status.success,
          warning: colors.status.warning,
          error: colors.status.error,
        },
      },
      borderRadius: {
        sm: radius.sm,
        md: radius.md,
        card: radius.card,
        lg: radius.lg,
        xl: radius.xl,
      },
      spacing: {
        "1": spacing[1],
        "2": spacing[2],
        "3": spacing[3],
        "4": spacing[4],
        "6": spacing[6],
        "8": spacing[8],
        "12": spacing[12],
        "16": spacing[16],
        "24": spacing[24],
        "28": spacing[28],
      },
      fontFamily: {
        display: ["var(--font-space-grotesk)", "Space Grotesk", "sans-serif"],
        body: ["var(--font-geist-sans)", "Geist Sans", "Inter", "sans-serif"],
        code: ["var(--font-jetbrains-mono)", "JetBrains Mono", "monospace"],
      },
    },
  },
  plugins: [],
};

export default config;
