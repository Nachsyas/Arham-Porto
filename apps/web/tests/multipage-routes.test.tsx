import { describe, it, expect } from "vitest";
import { render, screen, fireEvent, renderHook, act } from "@testing-library/react";
import React from "react";
import { PortfolioUIProvider, usePortfolioUI } from "../context/PortfolioUIContext";
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
});
