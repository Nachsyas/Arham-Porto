/**
 * Canonical constants and interfaces for the 3D Hologram Feature.
 * Derived from Arham Porto design tokens.
 */

export const HOLOGRAM_CONSTANTS = {
  // Brand & Hologram Colors
  COLORS: {
    PRIMARY: "#26B8FF",
    HOLOGRAM: "#61D8FF",
    HOLOGRAM_SOFT: "#2E8FB5",
    HOLOGRAM_FAINT: "#14384A",
    CANVAS: "#02060B",
    SURFACE: "#07111C",
  },

  // Rotation Constraints (Strictly Bounded per Design Specification)
  ROTATION: {
    MAX_YAW_RAD: 0.14,    // ~8 degrees, strictly bounded under 10 deg
    MAX_PITCH_RAD: 0.035, // ~2 degrees, strictly bounded under 3 deg
    DAMPING_FACTOR: 4.5,  // Smooth interpolation rate
  },

  // Camera & Viewport Defaults
  CAMERA: {
    FOV: 42,
    POSITION: [0, 0.35, 3.4] as [number, number, number],
  },

  // Performance Configuration
  PERFORMANCE: {
    DPR_DESKTOP: [1, 1.5] as [number, number],
    DPR_MOBILE: [1, 1.2] as [number, number],
    PARTICLES_DESKTOP: 45,
    PARTICLES_MOBILE: 20,
  },

  // Asset Loading Paths
  ASSETS: {
    PRODUCTION_MODEL: "/models/arham-wireframe.glb",
    STATUS_FLAG: "TODO_USER_3D_MODEL",
  },
} as const;

export type HologramMode = "idle" | "active";

export interface HologramStateContract {
  /** Future AI State: "idle" during passive browsing, "active" when interacting with AI reviewer */
  mode?: HologramMode;
  /** Future Journey Interaction: orientation offset in radians */
  orientationBias?: number;
  /** Future Intensity Multiplier: for audio or scan highlights */
  intensity?: number;
}
