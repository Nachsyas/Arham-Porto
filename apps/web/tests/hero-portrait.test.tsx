import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import HeroSection from "../features/hero/HeroSection";
import HeroPortrait from "../features/hero/HeroPortrait";
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
  todo: ["TODO_USER_HIGH_RES_PROFILE_PHOTO"],
};

describe("Hero Portrait Feature", () => {
  const originalMatchMedia = window.matchMedia;

  beforeEach(() => {
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
    vi.restoreAllMocks();
  });

  it("renders the real portrait image with accurate alt text", () => {
    render(<HeroPortrait fullName={mockProfile.fullName} role={mockProfile.role} />);

    const img = screen.getByRole("img", {
      name: `Portrait of ${mockProfile.fullName}`,
    });
    expect(img).toBeDefined();
    expect(img.getAttribute("src")).toContain("nachsyas-arham.jpg");
  });

  it("preserves all critical semantic Hero identity, positioning, and actions", () => {
    const onOpenQuickReview = vi.fn();
    render(<HeroSection profile={mockProfile} onOpenQuickReview={onOpenQuickReview} />);

    // Name typography
    expect(screen.getByText("NACHSYAS ARHAM")).toBeDefined();
    expect(screen.getByText("MUMTAZ NASHOHI")).toBeDefined();

    // Role and positioning
    expect(screen.getAllByText("Software Engineer").length).toBeGreaterThan(0);
    expect(
      screen.getByText("Building intelligent, scalable, and human-centered digital systems.")
    ).toBeDefined();

    // CTAs
    expect(screen.getByText("Explore My Work")).toBeDefined();
    const quickReviewBtn = screen.getByText("60-Second Quick Review");
    expect(quickReviewBtn).toBeDefined();

    // Quick review trigger
    fireEvent.click(quickReviewBtn);
    expect(onOpenQuickReview).toHaveBeenCalledTimes(1);
  });

  it("renders the graceful fallback card when image fails to load", () => {
    render(<HeroPortrait fullName={mockProfile.fullName} role={mockProfile.role} />);

    const img = screen.getByRole("img", {
      name: `Portrait of ${mockProfile.fullName}`,
    });

    // Simulate image error event
    fireEvent.error(img);

    // Fallback UI should render
    const fallback = screen.getByRole("img", {
      name: `Visual profile placeholder for ${mockProfile.fullName}`,
    });
    expect(fallback).toBeDefined();
    expect(screen.getByText(mockProfile.fullName)).toBeDefined();
  });

  it("honors prefers-reduced-motion without breaking portrait layout", () => {
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

    render(<HeroPortrait fullName={mockProfile.fullName} role={mockProfile.role} />);
    const img = screen.getByRole("img", {
      name: `Portrait of ${mockProfile.fullName}`,
    });
    expect(img).toBeDefined();
  });

  it("guarantees no raw TODO markers leak in rendered output", () => {
    const { container } = render(
      <HeroSection profile={mockProfile} onOpenQuickReview={vi.fn()} />
    );
    expect(container.innerHTML).not.toContain("TODO_USER");
  });
});
