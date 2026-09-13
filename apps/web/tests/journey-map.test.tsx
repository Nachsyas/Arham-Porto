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
    coordinates: [110.9515, -7.5976],
    period: null, // Birth year strictly omitted for privacy
    description: "Early formative origin in Karanganyar, Central Java.",
    image: null,
    public: true,
    todo: [],
  },
  {
    id: "journey-tk",
    category: "tk",
    title: null,
    institution: null,
    city: null,
    region: null,
    country: "Indonesia",
    coordinates: null,
    period: null,
    description: null,
    image: null,
    public: false, // Early education skipped from public Journey
    todo: [],
  },
  {
    id: "residence-jakarta",
    category: "residence",
    title: "Residence",
    institution: null,
    city: "Jakarta",
    region: "DKI Jakarta",
    country: "Indonesia",
    coordinates: [106.8272, -6.1754],
    period: null,
    description: "Formative base and residence in Jakarta, DKI Jakarta.",
    image: null,
    public: true,
    todo: [],
  },
  {
    id: "mi-al-hamid-jakarta",
    category: "sd",
    title: "Primary Education",
    institution: "Madrasah Ibtidaiyah Terpadu Al-Hamid",
    city: "Jakarta Timur",
    region: "DKI Jakarta",
    country: "Indonesia",
    coordinates: [106.9004, -6.225],
    period: null, // Unknown/unconfirmed period: must not be fabricated
    description: "Primary academic foundation at Madrasah Ibtidaiyah Terpadu Al-Hamid, Jakarta Timur.",
    image: null,
    public: true,
    todo: [],
  },
  {
    id: "mtsn30-jakarta",
    category: "smp",
    title: "Lower Secondary Education",
    institution: "Madrasah Tsanawiyah Negeri 30 Jakarta Timur",
    city: "Jakarta Timur",
    region: "DKI Jakarta",
    country: "Indonesia",
    coordinates: [106.9004, -6.225],
    period: null, // Unknown/unconfirmed period: must not be fabricated
    description: "Lower secondary academic foundation at MTsN 30 Jakarta Timur.",
    image: null,
    public: true,
    todo: [],
  },
  {
    id: "ma-assurkati-salatiga",
    category: "sma",
    title: "Tahfizh & Academic Foundation",
    institution: "Madrasah Aliyah Tahfizhul Qur'an As-Surkati",
    city: "Salatiga",
    region: "Jawa Tengah",
    country: "Indonesia",
    coordinates: [110.5084, -7.3305],
    period: "2019–2023",
    description: "Secondary education focusing on Quranic memorization (Tahfizh) and foundational academic studies.",
    image: null,
    public: true,
    todo: [],
  },
  {
    id: "university-uin-malang",
    category: "university",
    title: "Computer Science",
    institution: "Universitas Islam Negeri Maulana Malik Ibrahim Malang",
    city: "Malang",
    region: "Jawa Timur",
    country: "Indonesia",
    coordinates: [112.6081, -7.9525],
    period: "2023–Present",
    description: "Undergraduate studies in Computer Science / Informatics Engineering.",
    image: null,
    public: true,
    todo: [],
  },
  {
    id: "current-base-malang",
    category: "current",
    title: "Current Base",
    institution: null,
    city: "Malang",
    region: "Jawa Timur",
    country: "Indonesia",
    coordinates: [112.6308, -7.9826],
    period: "Present",
    description: "Active software engineering home base in Malang, East Java.",
    image: null,
    public: true,
    todo: [],
  },
];

describe("Phase 3 Interactive Journey Map (Expanded 7 Milestones)", () => {
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

  // 1. Origin renders
  it("renders Origin milestone", () => {
    render(<JourneySection stops={mockStops} />);
    expect(screen.getByText("ORIGIN")).toBeDefined();
    expect(screen.getAllByText(/Karanganyar/i).length).toBeGreaterThan(0);
  });

  // 2. Jakarta Residence renders
  it("renders Jakarta Residence milestone when selected", () => {
    render(<JourneySection stops={mockStops} />);
    const resStep = screen.getByTestId("timeline-step-residence-jakarta");
    fireEvent.click(resStep);
    expect(screen.getByText("RESIDENCE")).toBeDefined();
    expect(screen.getAllByText(/Jakarta/i).length).toBeGreaterThan(0);
  });

  // 3. MI Al-Hamid renders
  it("renders MI Al-Hamid milestone when selected", () => {
    render(<JourneySection stops={mockStops} />);
    const miStep = screen.getByTestId("timeline-step-mi-al-hamid-jakarta");
    fireEvent.click(miStep);
    expect(screen.getByText("PRIMARY EDUCATION")).toBeDefined();
    expect(screen.getByText("Madrasah Ibtidaiyah Terpadu Al-Hamid")).toBeDefined();
  });

  // 4. MTsN 30 Jakarta Timur renders
  it("renders MTsN 30 Jakarta Timur milestone when selected", () => {
    render(<JourneySection stops={mockStops} />);
    const mtsStep = screen.getByTestId("timeline-step-mtsn30-jakarta");
    fireEvent.click(mtsStep);
    expect(screen.getByText("LOWER SECONDARY EDUCATION")).toBeDefined();
    expect(screen.getByText("Madrasah Tsanawiyah Negeri 30 Jakarta Timur")).toBeDefined();
  });

  // 5. MA As-Surkati renders
  it("renders MA As-Surkati milestone when selected", () => {
    render(<JourneySection stops={mockStops} />);
    const maStep = screen.getByTestId("timeline-step-ma-assurkati-salatiga");
    fireEvent.click(maStep);
    expect(screen.getByText("SECONDARY EDUCATION")).toBeDefined();
    expect(screen.getByText("Madrasah Aliyah Tahfizhul Qur'an As-Surkati")).toBeDefined();
  });

  // 6. University renders
  it("renders University milestone when selected", () => {
    render(<JourneySection stops={mockStops} />);
    const uinStep = screen.getByTestId("timeline-step-university-uin-malang");
    fireEvent.click(uinStep);
    expect(screen.getByText("COMPUTER SCIENCE")).toBeDefined();
    expect(screen.getByText("Universitas Islam Negeri Maulana Malik Ibrahim Malang")).toBeDefined();
  });

  // 7. Current Base renders
  it("renders Current Base milestone when selected", () => {
    render(<JourneySection stops={mockStops} />);
    const currentStep = screen.getByTestId("timeline-step-current-base-malang");
    fireEvent.click(currentStep);
    expect(screen.getByText("CURRENT BASE")).toBeDefined();
  });

  // 8. Exactly four geographic map markers exist
  it("renders exactly four unique geographic markers on the map", () => {
    render(<JourneySection stops={mockStops} />);

    expect(screen.getByTestId("journey-marker-karanganyar")).toBeDefined();
    expect(screen.getByTestId("journey-marker-jakarta")).toBeDefined();
    expect(screen.getByTestId("journey-marker-salatiga")).toBeDefined();
    expect(screen.getByTestId("journey-marker-malang")).toBeDefined();

    const markers = screen.getAllByTestId(/^journey-marker-/);
    expect(markers).toHaveLength(4);
  });

  // 9. Residence + MI + MTs share Jakarta geography
  it("ensures Residence, MI, and MTs share the same geographic Jakarta pin", () => {
    render(<JourneySection stops={mockStops} />);
    const jakartaMarker = screen.getByTestId("journey-marker-jakarta");

    // Residence
    fireEvent.click(screen.getByTestId("timeline-step-residence-jakarta"));
    expect(jakartaMarker.getAttribute("aria-pressed")).toBe("true");

    // MI
    fireEvent.click(screen.getByTestId("timeline-step-mi-al-hamid-jakarta"));
    expect(jakartaMarker.getAttribute("aria-pressed")).toBe("true");

    // MTs
    fireEvent.click(screen.getByTestId("timeline-step-mtsn30-jakarta"));
    expect(jakartaMarker.getAttribute("aria-pressed")).toBe("true");
  });

  // 10. University + Current Base share Malang geography
  it("ensures University and Current Base share the same geographic Malang pin", () => {
    render(<JourneySection stops={mockStops} />);
    const malangMarker = screen.getByTestId("journey-marker-malang");

    // University
    fireEvent.click(screen.getByTestId("timeline-step-university-uin-malang"));
    expect(malangMarker.getAttribute("aria-pressed")).toBe("true");

    // Current Base
    fireEvent.click(screen.getByTestId("timeline-step-current-base-malang"));
    expect(malangMarker.getAttribute("aria-pressed")).toBe("true");
  });

  // 11. No Jakarta -> Jakarta route segments
  it("ensures no Jakarta -> Jakarta route segment exists in corridor curve", () => {
    const routePath = getCorridorRoutePath();
    // Jakarta coordinates: x: 172, y: 72
    const jakartaCoordMatches = (routePath.match(/\b172(\.0)? 72(\.0)?\b/g) || []).length;
    expect(jakartaCoordMatches).toBe(1);
  });

  // 12. No Malang -> Malang route segment
  it("ensures no Malang -> Malang route segment exists in corridor curve", () => {
    const routePath = getCorridorRoutePath();
    // Malang coordinates: (603.4, 239.2)
    const malangCoordMatches = (routePath.match(/603\.4 239\.2/g) || []).length;
    expect(malangCoordMatches).toBe(1);
  });

  // 13. Clicking Jakarta marker selects Residence when entering Jakarta
  it("clicking Jakarta marker selects Residence when entering Jakarta from outside", () => {
    render(<JourneySection stops={mockStops} />);
    // Currently on Origin (Karanganyar)
    const jakartaMarker = screen.getByTestId("journey-marker-jakarta");
    fireEvent.click(jakartaMarker);

    expect(screen.getByText("RESIDENCE")).toBeDefined();
  });

  // 14. Timeline can select MI
  it("timeline directly selects MI milestone", () => {
    render(<JourneySection stops={mockStops} />);
    fireEvent.click(screen.getByTestId("timeline-step-mi-al-hamid-jakarta"));
    expect(screen.getByText("Madrasah Ibtidaiyah Terpadu Al-Hamid")).toBeDefined();
  });

  // 15. Timeline can select MTs
  it("timeline directly selects MTs milestone", () => {
    render(<JourneySection stops={mockStops} />);
    fireEvent.click(screen.getByTestId("timeline-step-mtsn30-jakarta"));
    expect(screen.getByText("Madrasah Tsanawiyah Negeri 30 Jakarta Timur")).toBeDefined();
  });

  // 16. Previous / Next traverses all seven milestones
  it("traverses all seven milestones sequentially using Next and Previous", () => {
    render(<JourneySection stops={mockStops} />);
    const prevBtn = screen.getByRole("button", { name: /previous/i });
    const nextBtn = screen.getByRole("button", { name: /next/i });

    expect(prevBtn.getAttribute("disabled")).not.toBeNull();
    expect(screen.getByText("ORIGIN")).toBeDefined();

    // 01 Origin -> 02 Residence
    fireEvent.click(nextBtn);
    expect(screen.getByText("RESIDENCE")).toBeDefined();

    // 02 Residence -> 03 MI
    fireEvent.click(nextBtn);
    expect(screen.getByText("PRIMARY EDUCATION")).toBeDefined();

    // 03 MI -> 04 MTs
    fireEvent.click(nextBtn);
    expect(screen.getByText("LOWER SECONDARY EDUCATION")).toBeDefined();

    // 04 MTs -> 05 MA
    fireEvent.click(nextBtn);
    expect(screen.getByText("SECONDARY EDUCATION")).toBeDefined();

    // 05 MA -> 06 University
    fireEvent.click(nextBtn);
    expect(screen.getByText("COMPUTER SCIENCE")).toBeDefined();

    // 06 University -> 07 Current Base
    fireEvent.click(nextBtn);
    expect(screen.getByText("CURRENT BASE")).toBeDefined();
    expect(nextBtn.getAttribute("disabled")).not.toBeNull();

    // Step back: 07 -> 06
    fireEvent.click(prevBtn);
    expect(screen.getByText("COMPUTER SCIENCE")).toBeDefined();
  });

  // 17. MI and MTs render no fabricated period
  it("does not render any fabricated or placeholder period for MI and MTs", () => {
    render(<JourneySection stops={mockStops} />);

    // Check MI
    fireEvent.click(screen.getByTestId("timeline-step-mi-al-hamid-jakarta"));
    const miArticle = screen.getByRole("article", { name: /Milestone details/i });
    expect(miArticle.textContent).not.toContain("Unknown");
    expect(miArticle.textContent).not.toContain("TBD");
    expect(miArticle.textContent).not.toContain("TODO");

    // Check MTs
    fireEvent.click(screen.getByTestId("timeline-step-mtsn30-jakarta"));
    const mtsArticle = screen.getByRole("article", { name: /Milestone details/i });
    expect(mtsArticle.textContent).not.toContain("Unknown");
    expect(mtsArticle.textContent).not.toContain("TBD");
    expect(mtsArticle.textContent).not.toContain("TODO");
  });

  // 18. Birth year remains hidden
  it("strictly hides birth year (2004) across all rendered milestones", () => {
    const { container } = render(<JourneySection stops={mockStops} />);
    expect(container.innerHTML).not.toContain("2004");
  });

  // 19. Residential address remains absent
  it("strictly excludes private residential address or complex name", () => {
    const { container } = render(<JourneySection stops={mockStops} />);
    const html = container.innerHTML.toLowerCase();
    expect(html).not.toContain("bambu kuning");
    expect(html).not.toContain("cipayung");
    expect(html).not.toContain("rt ");
    expect(html).not.toContain("rw ");
    expect(html).not.toContain("street");
  });

  // 20. TODO strings remain absent
  it("strictly asserts no TODO strings leak into the rendered DOM", () => {
    const { container } = render(<JourneySection stops={mockStops} />);
    expect(container.innerHTML).not.toContain("TODO_USER");
    expect(container.innerHTML).not.toContain("TODO_");
  });

  // 21. Reduced motion remains functional
  it("honors prefers-reduced-motion without breaking rendering", () => {
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

  // 22. Keyboard interaction remains functional
  it("navigates milestones using keyboard arrow keys", () => {
    render(<JourneySection stops={mockStops} />);
    expect(screen.getByText("ORIGIN")).toBeDefined();

    // ArrowRight -> Residence
    fireEvent.keyDown(window, { key: "ArrowRight" });
    expect(screen.getByText("RESIDENCE")).toBeDefined();

    // ArrowLeft -> Origin
    fireEvent.keyDown(window, { key: "ArrowLeft" });
    expect(screen.getByText("ORIGIN")).toBeDefined();
  });
});
