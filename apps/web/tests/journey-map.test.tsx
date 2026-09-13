import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import JourneySection from "../features/journey/JourneySection";
import { getCorridorRoutePath, filterPublicMilestones } from "../features/journey/journey.utils";
import type { JourneyStop } from "arham-porto-schema";

// Canonical Scenario A mock data matching data/journey/journey.json
const mockStops: JourneyStop[] = [
  {
    id: "origin-karanganyar",
    category: "birthplace",
    title: "Origin",
    institution: null,
    city: "Karanganyar",
    region: "Jawa Tengah",
    country: "Indonesia",
    coordinates: [-7.5968, 110.9515],
    period: null, // Birth year strictly omitted for privacy
    description: "Formative upbringing and roots in Karanganyar, Jawa Tengah.",
    image: null,
    public: true,
    todo: [],
  },
  {
    id: "journey-tk",
    category: "tk",
    title: "Pendidikan Kanak-Kanak",
    institution: null,
    city: "Karanganyar",
    region: "Jawa Tengah",
    country: "Indonesia",
    coordinates: null,
    period: null,
    description: null,
    image: null,
    public: false, // Early education skipped from public Journey
    todo: ["TODO_USER_TK"],
  },
  {
    id: "journey-sd",
    category: "sd",
    title: "Pendidikan Dasar",
    institution: null,
    city: "Karanganyar",
    region: "Jawa Tengah",
    country: "Indonesia",
    coordinates: null,
    period: null,
    description: null,
    image: null,
    public: false,
    todo: ["TODO_USER_SD"],
  },
  {
    id: "journey-smp",
    category: "smp",
    title: "Pendidikan Menengah Pertama",
    institution: null,
    city: "Karanganyar",
    region: "Jawa Tengah",
    country: "Indonesia",
    coordinates: null,
    period: null,
    description: null,
    image: null,
    public: false,
    todo: ["TODO_USER_SMP"],
  },
  {
    id: "ma-assurkati-salatiga",
    category: "sma",
    title: "Tahfizh & Academic Foundation",
    institution: "Madrasah Aliyah Tahfizhul Qur'an As-Surkati",
    city: "Salatiga",
    region: "Jawa Tengah",
    country: "Indonesia",
    coordinates: [-7.3305, 110.5084],
    period: "2019–2023",
    description: "Rigorous tahfizh curriculum and foundational secondary education in Salatiga, Jawa Tengah.",
    image: null,
    public: true,
    todo: [],
  },
  {
    id: "university-uin-malang",
    category: "university",
    title: "Computer Science Undergraduate",
    institution: "Universitas Islam Negeri Maulana Malik Ibrahim Malang",
    city: "Malang",
    region: "Jawa Timur",
    country: "Indonesia",
    coordinates: [-7.9525, 112.6079],
    period: "2023–Present",
    description: "Undergraduate studies in Computer Science at UIN Maulana Malik Ibrahim Malang, focusing on software engineering, backend systems, and algorithm design.",
    image: null,
    public: true,
    todo: [],
  },
  {
    id: "current-base-malang",
    category: "current",
    title: "Current Engineering Base",
    institution: null,
    city: "Malang",
    region: "Jawa Timur",
    country: "Indonesia",
    coordinates: [-7.9839, 112.6214],
    period: "Present",
    description: "Active software engineering base in Malang, Jawa Timur.",
    image: null,
    public: true,
    todo: [],
  },
];

describe("Phase 3 Interactive Journey Map", () => {
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

  // 1. Public filtering
  it("renders only public milestones and completely excludes private stops (TK, SD, SMP)", () => {
    const publicStops = filterPublicMilestones(mockStops);
    expect(publicStops).toHaveLength(4);

    const { container } = render(<JourneySection stops={mockStops} />);
    expect(container.innerHTML).not.toContain("Pendidikan Kanak-Kanak");
    expect(container.innerHTML).not.toContain("Pendidikan Dasar");
    expect(container.innerHTML).not.toContain("Pendidikan Menengah Pertama");
  });

  // 2. Exactly three unique geographic markers
  it("renders exactly three unique geographic markers on the map", () => {
    render(<JourneySection stops={mockStops} />);

    const karanganyarMarker = screen.getByTestId("journey-marker-karanganyar");
    const salatigaMarker = screen.getByTestId("journey-marker-salatiga");
    const malangMarker = screen.getByTestId("journey-marker-malang");

    expect(karanganyarMarker).toBeDefined();
    expect(salatigaMarker).toBeDefined();
    expect(malangMarker).toBeDefined();

    // Exactly 3 marker buttons exist
    const markers = screen.getAllByTestId(/^journey-marker-/);
    expect(markers).toHaveLength(3);
  });

  // 3. Four timeline milestones exist
  it("renders exactly four timeline milestones representing the chronological narrative", () => {
    render(<JourneySection stops={mockStops} />);

    const timelineSteps = screen.getAllByTestId(/^timeline-step-/);
    expect(timelineSteps).toHaveLength(4);

    expect(screen.getByTestId("timeline-step-origin-karanganyar")).toBeDefined();
    expect(screen.getByTestId("timeline-step-ma-assurkati-salatiga")).toBeDefined();
    expect(screen.getByTestId("timeline-step-university-uin-malang")).toBeDefined();
    expect(screen.getByTestId("timeline-step-current-base-malang")).toBeDefined();
  });

  // 4. University and Current Base share Malang location
  it("ensures University and Current Base share the same geographic Malang pin", () => {
    render(<JourneySection stops={mockStops} />);

    const malangMarker = screen.getByTestId("journey-marker-malang");

    // Click timeline step for University
    const uinTimelineStep = screen.getByTestId("timeline-step-university-uin-malang");
    fireEvent.click(uinTimelineStep);
    expect(malangMarker.getAttribute("aria-pressed")).toBe("true");

    // Click timeline step for Current Base
    const currentTimelineStep = screen.getByTestId("timeline-step-current-base-malang");
    fireEvent.click(currentTimelineStep);
    expect(malangMarker.getAttribute("aria-pressed")).toBe("true");
  });

  // 5. No Malang -> Malang route segment
  it("does not construct a Malang -> Malang route segment in the corridor SVG curve", () => {
    const routePath = getCorridorRoutePath();
    // Route must connect Karanganyar -> Salatiga -> Malang
    expect(routePath).toContain("M 478.3 203.2");
    expect(routePath).toContain("445.3 178.2");
    expect(routePath).toContain("603.4 239.2");

    // Must not have a redundant second Malang coordinate pair in the corridor curve
    const malangCoordMatches = (routePath.match(/603\.4 239\.2/g) || []).length;
    expect(malangCoordMatches).toBe(1);
  });

  // 6. Clicking marker updates milestone card
  it("updates active milestone card when clicking map markers", () => {
    render(<JourneySection stops={mockStops} />);

    // Initially defaults to Origin
    expect(screen.getByText("ORIGIN")).toBeDefined();
    expect(screen.getByText("Karanganyar")).toBeDefined();

    // Click Salatiga marker
    const salatigaMarker = screen.getByTestId("journey-marker-salatiga");
    fireEvent.click(salatigaMarker);

    expect(screen.getByText("TAHFIZH & ACADEMIC FOUNDATION")).toBeDefined();
    expect(screen.getByText("Madrasah Aliyah Tahfizhul Qur'an As-Surkati")).toBeDefined();
  });

  // 7. Clicking timeline updates milestone card
  it("updates active milestone card when clicking timeline buttons", () => {
    render(<JourneySection stops={mockStops} />);

    // Click University timeline step
    const universityStep = screen.getByTestId("timeline-step-university-uin-malang");
    fireEvent.click(universityStep);

    expect(screen.getByText("COMPUTER SCIENCE")).toBeDefined();
    expect(screen.getByText("Universitas Islam Negeri Maulana Malik Ibrahim Malang")).toBeDefined();
    expect(screen.getByText("2023–Present")).toBeDefined();
  });

  // 8. Previous / Next buttons work
  it("steps through milestones sequentially using Previous and Next controls", () => {
    render(<JourneySection stops={mockStops} />);

    const prevBtn = screen.getByRole("button", { name: /previous/i });
    const nextBtn = screen.getByRole("button", { name: /next/i });

    // Initial state: on step 1 (Origin) -> Previous should be disabled
    expect(prevBtn.getAttribute("disabled")).not.toBeNull();
    expect(nextBtn.getAttribute("disabled")).toBeNull();

    // Click Next -> advances to step 2 (Salatiga)
    fireEvent.click(nextBtn);
    expect(screen.getByText("TAHFIZH & ACADEMIC FOUNDATION")).toBeDefined();
    expect(prevBtn.getAttribute("disabled")).toBeNull();

    // Click Next -> advances to step 3 (University)
    fireEvent.click(nextBtn);
    expect(screen.getByText("COMPUTER SCIENCE")).toBeDefined();

    // Click Next -> advances to step 4 (Current Base)
    fireEvent.click(nextBtn);
    expect(screen.getByText("CURRENT BASE")).toBeDefined();
    expect(nextBtn.getAttribute("disabled")).not.toBeNull();

    // Click Previous -> steps back to step 3
    fireEvent.click(prevBtn);
    expect(screen.getByText("COMPUTER SCIENCE")).toBeDefined();
  });

  // 9. Keyboard interaction works
  it("navigates milestones using keyboard arrow keys", () => {
    render(<JourneySection stops={mockStops} />);

    expect(screen.getByText("ORIGIN")).toBeDefined();

    // Press ArrowRight
    fireEvent.keyDown(window, { key: "ArrowRight" });
    expect(screen.getByText("TAHFIZH & ACADEMIC FOUNDATION")).toBeDefined();

    // Press ArrowLeft
    fireEvent.keyDown(window, { key: "ArrowLeft" });
    expect(screen.getByText("ORIGIN")).toBeDefined();
  });

  // 10. Reduced-motion behavior works
  it("honors prefers-reduced-motion without breaking map rendering", () => {
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

    render(<JourneySection stops={mockStops} />);
    const corridorPath = screen.getByTestId("journey-corridor-path");
    expect(corridorPath).toBeDefined();
    expect(corridorPath.getAttribute("style")).toContain("stroke-dasharray: none");
  });

  // 11. Privacy values do not render
  it("strictly asserts private values (birth year 2004, full birth date, phone, address) do not render", () => {
    const { container } = render(<JourneySection stops={mockStops} />);
    const html = container.innerHTML;

    // No birth year 2004
    expect(html).not.toContain("2004");
    // No private words
    expect(html).not.toContain("birth_date");
    expect(html).not.toContain("telephone");
    expect(html).not.toContain("phone");
    expect(html).not.toContain("street");
    expect(html).not.toContain("residence");
  });

  // 12. No TODO markers render
  it("strictly asserts no TODO markers leak into the rendered HTML", () => {
    const { container } = render(<JourneySection stops={mockStops} />);
    expect(container.innerHTML).not.toContain("TODO_USER");
    expect(container.innerHTML).not.toContain("TODO_");
  });
});
