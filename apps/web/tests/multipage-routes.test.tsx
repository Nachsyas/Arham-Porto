import { describe, it, expect } from "vitest";
import { render, screen, fireEvent, renderHook, act } from "@testing-library/react";
import React from "react";
import { PortfolioUIProvider, usePortfolioUI } from "../context/PortfolioUIContext";
import { NavigationProgressProvider, useNavigationProgress } from "../context/NavigationProgressContext";
import NavLink, { isModifiedClick } from "../components/motion/NavLink";
import { getSkills, getProjects, getEvidence, getProfile, getJourney } from "arham-porto-data";
import SkillsExplorer from "../features/skills/SkillsExplorer";
import HomeView from "../features/home/HomeView";
import ContactView from "../features/contact/ContactView";
import WorkPage from "../app/work/page";
import SkillsPage from "../app/skills/page";
import JourneyPage from "../app/journey/page";
import ContactPage from "../app/contact/page";
import Navbar from "../components/Navbar";

describe("Multi-Page Experience & Routing Architecture", () => {
  describe("PortfolioUIProvider & Context Hook", () => {
    it("manages openAskArham, closeAskArham, openQuickReview, and closeQuickReview states correctly", () => {
      const wrapper = ({ children }: { children: React.ReactNode }) => (
        <PortfolioUIProvider>{children}</PortfolioUIProvider>
      );

      const { result } = renderHook(() => usePortfolioUI(), { wrapper });

      // Initial state
      expect(result.current.isAskArhamOpen).toBe(false);
      expect(result.current.isQuickReviewOpen).toBe(false);

      // openAskArham
      act(() => {
        result.current.openAskArham();
      });
      expect(result.current.isAskArhamOpen).toBe(true);

      // closeAskArham
      act(() => {
        result.current.closeAskArham();
      });
      expect(result.current.isAskArhamOpen).toBe(false);

      // openQuickReview
      act(() => {
        result.current.openQuickReview();
      });
      expect(result.current.isQuickReviewOpen).toBe(true);

      // closeQuickReview
      act(() => {
        result.current.closeQuickReview();
      });
      expect(result.current.isQuickReviewOpen).toBe(false);
    });

    it("throws an error when usePortfolioUI is invoked outside PortfolioUIProvider", () => {
      // Suppress error logging during this intentional test
      const originalError = console.error;
      console.error = () => {};

      expect(() => {
        renderHook(() => usePortfolioUI());
      }).toThrow("usePortfolioUI must be used within a PortfolioUIProvider");

      console.error = originalError;
    });
  });

  describe("Navbar Route Navigation", () => {
    it("renders all canonical multi-page route destinations", () => {
      render(
        <PortfolioUIProvider>
          <Navbar />
        </PortfolioUIProvider>
      );

      expect(screen.getByRole("link", { name: /^Home$/i }).getAttribute("href")).toBe("/");
      expect(screen.getByRole("link", { name: /^Work$/i }).getAttribute("href")).toBe("/work");
      expect(screen.getByRole("link", { name: /^Skills$/i }).getAttribute("href")).toBe("/skills");
      expect(screen.getByRole("link", { name: /^Journey$/i }).getAttribute("href")).toBe("/journey");
      expect(screen.getByRole("link", { name: /^Contact$/i }).getAttribute("href")).toBe("/contact");
    });
  });

  describe("Dedicated Route Pages", () => {
    it("renders /work with PageHeader and project cards", () => {
      render(
        <PortfolioUIProvider>
          <WorkPage />
        </PortfolioUIProvider>
      );

      expect(screen.getByText("WORK // SELECTED PROJECTS")).toBeDefined();
      expect(screen.getByRole("heading", { name: /Selected Systems & Engineering Projects/i })).toBeDefined();
      expect(screen.getByRole("heading", { name: "EduTrace" })).toBeDefined();
      expect(screen.getByRole("heading", { name: "GDGOC E-Commerce" })).toBeDefined();
    });

    it("renders /skills with PageHeader and verified evidence without arbitrary percentage numbers", () => {
      const { container } = render(
        <PortfolioUIProvider>
          <SkillsPage />
        </PortfolioUIProvider>
      );

      expect(screen.getByText("SKILLS // VERIFIED CAPABILITIES")).toBeDefined();
      expect(screen.getByRole("heading", { name: /Capabilities Backed by Verified Evidence/i })).toBeDefined();

      // Zero arbitrary percentage ratings
      const text = container.textContent || "";
      expect(text).not.toMatch(/\d+%/);
    });

    it("renders /journey with PageHeader, interactive map, and academic history", () => {
      render(
        <PortfolioUIProvider>
          <JourneyPage />
        </PortfolioUIProvider>
      );

      expect(screen.getByText("JOURNEY // GEOGRAPHIC & ACADEMIC")).toBeDefined();
      expect(screen.getByRole("heading", { name: /Geographic & Academic Evolution/i })).toBeDefined();
      expect(screen.getByLabelText("Stylized geographic map showing academic journey across Central Java and East Java, Indonesia")).toBeDefined();
      expect(screen.getByRole("heading", { name: /Academic & Formative Records/i })).toBeDefined();
    });

    it("renders /contact with PageHeader, verified GitHub channel, and Ask Arham CTA", () => {
      render(
        <PortfolioUIProvider>
          <ContactPage />
        </PortfolioUIProvider>
      );

      expect(screen.getByText("CONTACT // VERIFIED CHANNELS")).toBeDefined();
      expect(screen.getByRole("heading", { name: /Get in Touch & Technical Inquiries/i })).toBeDefined();
      expect(screen.getByRole("link", { name: /github\.com\/Nachsyas/i }).getAttribute("href")).toBe("https://github.com/Nachsyas");
      expect(screen.getByRole("button", { name: /Launch Ask Arham Assistant/i })).toBeDefined();
    });
  });

  describe("Evidence Integrity, Claim Precision, and Route Semantics Regressions", () => {
    const skills = getSkills();
    const projects = getProjects();
    const evidence = getEvidence();
    const profile = getProfile();
    const journeyStops = getJourney();

    it("A: ai-rag with evidenceIds=[] does NOT show 'Verified' or 'Evidence Grounded'", () => {
      render(
        <PortfolioUIProvider>
          <SkillsExplorer skills={skills} projects={projects} evidence={evidence} initialSkillId="ai-rag" />
        </PortfolioUIProvider>
      );

      // Verify that "ai-rag" button does NOT say "Verified"
      const ragButton = screen.getByRole("button", { name: /Retrieval-Augmented Generation/i });
      expect(ragButton.textContent).not.toMatch(/Verified/i);
      expect(ragButton.textContent).toContain("Portfolio Skill");

      // Verify detail panel does NOT say "Evidence Grounded"
      expect(screen.queryByText(/Evidence Grounded/i)).toBeNull();
      expect(screen.getByText("No linked portfolio evidence yet")).toBeDefined();

      // Verify claim heading is neutral overview, not verified claim
      expect(screen.queryByText(/VERIFIED CAPABILITY CLAIM:/i)).toBeNull();
      expect(screen.getByText("CAPABILITY OVERVIEW:")).toBeDefined();
    });

    it("B: ai-rag does NOT automatically display Maritime AI as evidence", () => {
      render(
        <PortfolioUIProvider>
          <SkillsExplorer skills={skills} projects={projects} evidence={evidence} />
        </PortfolioUIProvider>
      );

      const ragButton = screen.getByRole("button", { name: /Retrieval-Augmented Generation/i });
      fireEvent.click(ragButton);

      // Maritime AI should NOT be listed under demonstrated projects
      expect(screen.queryByText(/Maritime AI/i)).toBeNull();
    });

    it("C: empty evidence state is shown correctly for skills without linked evidence", () => {
      render(
        <PortfolioUIProvider>
          <SkillsExplorer skills={skills} projects={projects} evidence={evidence} />
        </PortfolioUIProvider>
      );

      const ragButton = screen.getByRole("button", { name: /Retrieval-Augmented Generation/i });
      fireEvent.click(ragButton);

      // Evidence list empty state
      expect(
        screen.getByText("No linked portfolio evidence is currently available for this skill.")
      ).toBeDefined();

      // Project list empty state
      expect(
        screen.getByText("No public projects are currently linked to this skill via canonical evidence.")
      ).toBeDefined();
    });

    it("D: Home does not claim all capabilities have evidence", () => {
      const { container } = render(
        <PortfolioUIProvider>
          <NavigationProgressProvider>
            <HomeView
              profile={profile}
              featuredProjects={projects.filter((p) => p.featured)}
              skills={skills}
              journeyStops={journeyStops}
            />
          </NavigationProgressProvider>
        </PortfolioUIProvider>
      );

      // Must not claim every capability has evidence
      expect(container.textContent).not.toMatch(
        /Every claimed technical capability in this portfolio is grounded in verifiable repository commits/i
      );
      expect(
        screen.getByText(
          /Technical evidence is linked where verified repository evidence is available\./i
        )
      ).toBeDefined();
    });

    it("E: Home does not claim candidate availability when profile.availability is null", () => {
      expect(profile.availability).toBeNull();

      const { container } = render(
        <PortfolioUIProvider>
          <NavigationProgressProvider>
            <HomeView
              profile={profile}
              featuredProjects={projects.filter((p) => p.featured)}
              skills={skills}
              journeyStops={journeyStops}
            />
          </NavigationProgressProvider>
        </PortfolioUIProvider>
      );

      // Must not infer availability
      expect(container.textContent).not.toMatch(
        /Open for software engineering opportunities, distributed backend collaboration/i
      );
      expect(
        screen.getByText(
          /Explore verified engineering work, public repositories, or ask the grounded portfolio reviewer\./i
        )
      ).toBeDefined();
    });

    it("F: Contact does not contain 'hybrid vector search' or 'microservices'", () => {
      const { container } = render(
        <PortfolioUIProvider>
          <ContactView profile={profile} />
        </PortfolioUIProvider>
      );

      expect(container.textContent).not.toMatch(/hybrid vector search/i);
      expect(container.textContent).not.toMatch(/microservices/i);
      expect(
        screen.getByText(
          /Powered by semantic vector retrieval over approved portfolio sources with strict citation anchoring/i
        )
      ).toBeDefined();
      expect(
        screen.getByText(/Review verified commit history, Go backend services, Next\.js web applications/i)
      ).toBeDefined();
    });

    it("G: Journey does not invent 'boarding school'", () => {
      const { container: homeContainer } = render(
        <PortfolioUIProvider>
          <NavigationProgressProvider>
            <HomeView
              profile={profile}
              featuredProjects={projects.filter((p) => p.featured)}
              skills={skills}
              journeyStops={journeyStops}
            />
          </NavigationProgressProvider>
        </PortfolioUIProvider>
      );

      expect(homeContainer.textContent).not.toMatch(/boarding school/i);
      expect(homeContainer.textContent).toContain("Tahfizh & Academic Foundation");
      expect(homeContainer.textContent).toContain("Karanganyar → Jakarta → Salatiga → Malang");

      const { container: journeyContainer } = render(
        <PortfolioUIProvider>
          <JourneyPage />
        </PortfolioUIProvider>
      );

      expect(journeyContainer.textContent).not.toMatch(/boarding school/i);
      expect(journeyContainer.textContent).toContain("Karanganyar → Jakarta → Salatiga → Malang");
    });

    it("H: Route navigation and navigation progress context function normally", () => {
      const wrapper = ({ children }: { children: React.ReactNode }) => (
        <PortfolioUIProvider>
          <NavigationProgressProvider>{children}</NavigationProgressProvider>
        </PortfolioUIProvider>
      );

      const { result } = renderHook(() => useNavigationProgress(), { wrapper });

      expect(result.current.isNavigating).toBe(false);

      act(() => {
        result.current.startProgress();
      });
      expect(result.current.isNavigating).toBe(true);

      act(() => {
        result.current.completeProgress();
      });
      expect(result.current.isNavigating).toBe(false);

      // Test isModifiedClick modifier safety
      const normalClick = { metaKey: false, ctrlKey: false, shiftKey: false, altKey: false, button: 0 } as React.MouseEvent;
      expect(isModifiedClick(normalClick)).toBe(false);

      const cmdClick = { metaKey: true, ctrlKey: false, shiftKey: false, altKey: false, button: 0 } as React.MouseEvent;
      expect(isModifiedClick(cmdClick)).toBe(true);

      const middleClick = { metaKey: false, ctrlKey: false, shiftKey: false, altKey: false, button: 1 } as React.MouseEvent;
      expect(isModifiedClick(middleClick)).toBe(true);

      // Test NavLink component interactions
      function NavLinkHarness() {
        const { isNavigating } = useNavigationProgress();
        return (
          <div>
            <span data-testid="nav-status">{isNavigating ? "NAVIGATING" : "IDLE"}</span>
            <NavLink href="/work">Go to Work</NavLink>
            <NavLink href="#overview">Anchor Link</NavLink>
            <NavLink href="https://github.com" target="_blank">External Link</NavLink>
          </div>
        );
      }

      render(
        <PortfolioUIProvider>
          <NavigationProgressProvider>
            <NavLinkHarness />
          </NavigationProgressProvider>
        </PortfolioUIProvider>
      );

      expect(screen.getByTestId("nav-status").textContent).toBe("IDLE");

      // Suppress jsdom navigation error for internal link test
      const originalConsoleError = console.error;
      console.error = (...args: unknown[]) => {
        if (typeof args[0] === "string" && args[0].includes("Not implemented: navigation")) return;
        if (args[0] instanceof Error && args[0].message.includes("Not implemented: navigation")) return;
        originalConsoleError(...args);
      };

      try {
        // Clicking hash anchor should not trigger navigation progress
        fireEvent.click(screen.getByText("Anchor Link"));
        expect(screen.getByTestId("nav-status").textContent).toBe("IDLE");

        // Clicking external target=_blank link should not trigger navigation progress
        fireEvent.click(screen.getByText("External Link"));
        expect(screen.getByTestId("nav-status").textContent).toBe("IDLE");

        // Clicking internal route link triggers navigation progress
        fireEvent.click(screen.getByText("Go to Work"));
        expect(screen.getByTestId("nav-status").textContent).toBe("NAVIGATING");
      } finally {
        console.error = originalConsoleError;
      }
    });
  });
});
