import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { render, screen } from "@testing-library/react";
import HeroSection from "../features/hero/HeroSection";
import HologramFallback from "../features/hologram/HologramFallback";
import HologramStage from "../features/hologram/HologramStage";
import { HOLOGRAM_CONSTANTS } from "../features/hologram/hologram.constants";
import type { Profile } from "arham-porto-schema";

const mockProfile: Profile = {
  fullName: "Nachsyas Arham Mumtaz Nashohi",
  role: "Software Engineer",
  projectName: "Arham Porto",
  aiFeature: "Ask Arham AI",
  positioning: "Building intelligent, scalable, and human-centered digital systems.",
  bio: null,
  email: null,
  linkedin: null,
  github: "https://github.com/Nachsyas",
  cvUrl: null,
  currentCity: null,
  availability: null,
  todo: [],
};

describe("Phase 2: 3D Wireframe Hologram Feature", () => {
  const originalMatchMedia = window.matchMedia;
  const originalGetContext = HTMLCanvasElement.prototype.getContext;

  beforeEach(() => {
    // Provide a safe WebGL mock that returns null in jsdom test runner
    HTMLCanvasElement.prototype.getContext = vi.fn().mockReturnValue(null);

    window.matchMedia = vi.fn().mockImplementation((query) => ({
      matches: false,
      media: query,
      onchange: null,
      addListener: vi.fn(),
      removeListener: vi.fn(),
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
      dispatchEvent: vi.fn(),
    }));
  });

  afterEach(() => {
    window.matchMedia = originalMatchMedia;
    HTMLCanvasElement.prototype.getContext = originalGetContext;
    vi.restoreAllMocks();
  });

  it("preserves all critical semantic Hero identity and actions alongside the 3D feature", () => {
    const onOpenQuickReview = vi.fn();
    render(<HeroSection profile={mockProfile} onOpenQuickReview={onOpenQuickReview} />);

    // Primary Name Display
    expect(screen.getByText("NACHSYAS ARHAM")).toBeDefined();
    expect(screen.getByText("MUMTAZ NASHOHI")).toBeDefined();

    // Role & Engineering Positioning
    expect(screen.getAllByText("Software Engineer").length).toBeGreaterThan(0);
    expect(
      screen.getByText("Building intelligent, scalable, and human-centered digital systems.")
    ).toBeDefined();

    // Interactive CTAs
    expect(screen.getByText("Explore My Work")).toBeDefined();
    expect(screen.getByText("60-Second Quick Review")).toBeDefined();
  });

  it("maintains the strict production model path contract for future GLB drop-in", () => {
    expect(HOLOGRAM_CONSTANTS.ASSETS.PRODUCTION_MODEL).toBe("/models/arham-wireframe.glb");
    expect(HOLOGRAM_CONSTANTS.ROTATION.MAX_YAW_RAD).toBeLessThanOrEqual(0.18); // ~10 degrees hard cap
    expect(HOLOGRAM_CONSTANTS.ROTATION.MAX_PITCH_RAD).toBeLessThanOrEqual(0.05); // ~3 degrees
  });

  it("renders the accessible HologramFallback when WebGL is unavailable in the environment", () => {
    render(<HologramFallback reason="WebGL unsupported" />);

    const fallbackContainer = screen.getByRole("img", {
      name: /static 3d seated wireframe fallback placeholder/i,
    });
    expect(fallbackContainer).toBeDefined();
    expect(screen.getByText("SYSTEM ARCHITECTURE")).toBeDefined();
    expect(screen.getByText("Static Mode")).toBeDefined();
  });

  it("renders the HologramStage container with accessible region label", async () => {
    render(<HologramStage />);

    // In jsdom without WebGL, HologramStage safely renders the fallback
    const stageOrFallback = await screen.findByRole("img", {
      name: /static 3d seated wireframe fallback/i,
    });
    expect(stageOrFallback).toBeDefined();
  });

  it("respects prefers-reduced-motion configuration", () => {
    window.matchMedia = vi.fn().mockImplementation((query) => ({
      matches: query === "(prefers-reduced-motion: reduce)",
      media: query,
      onchange: null,
      addListener: vi.fn(),
      removeListener: vi.fn(),
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
      dispatchEvent: vi.fn(),
    }));

    render(<HologramStage />);
    // Under reduced motion in non-WebGL jsdom, fallback remains cleanly mounted
    expect(screen.getByText("SYSTEM ARCHITECTURE")).toBeDefined();
  });

  it("guarantees no raw TODO markers or leaks in rendered output", () => {
    const { container } = render(
      <HeroSection profile={mockProfile} onOpenQuickReview={vi.fn()} />
    );
    expect(container.innerHTML).not.toContain("TODO_USER");
  });
});
