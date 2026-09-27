import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import React from "react";
import GithubProjectCard from "../features/github-projects/GithubProjectCard";
import GithubProjectsSection from "../features/github-projects/GithubProjectsSection";
import { getGithubProjects } from "../features/github-projects/github-projects.api";
import type { GithubProjectItem } from "../features/github-projects/github-projects.types";
import WorkPage from "../app/work/page";
import { PortfolioUIProvider } from "../context/PortfolioUIContext";

const mockProject1: GithubProjectItem = {
  id: "uuid-1",
  name: "EduTrace",
  full_name: "Nachsyas/EduTrace",
  description: "Educational tracking and analytics platform",
  html_url: "https://github.com/Nachsyas/EduTrace",
  homepage: "https://edutrace.example.com",
  language: "TypeScript",
  stars: 12,
  forks: 3,
  topics: ["nextjs", "ai", "education"],
  synced_at: "2026-09-14T00:00:00Z",
};

const mockProject2: GithubProjectItem = {
  id: "uuid-2",
  name: "GoStream",
  full_name: "Nachsyas/GoStream",
  description: "High-throughput telemetry streaming in pure Go",
  html_url: "https://github.com/Nachsyas/GoStream",
  homepage: "",
  language: "Go",
  stars: 5,
  forks: 0,
  topics: ["golang", "concurrency"],
  synced_at: "2026-09-20T12:00:00Z",
};

describe("GitHub Project Explorer (Phase 2A)", () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe("GithubProjectCard", () => {
    it("renders repository details, language, metrics, topics, and external link", () => {
      render(<GithubProjectCard project={mockProject1} />);

      // Repository Name
      expect(screen.getByRole("heading", { name: "EduTrace" })).toBeDefined();

      // Description
      expect(screen.getByText("Educational tracking and analytics platform")).toBeDefined();

      // Language badge
      expect(screen.getByText("TypeScript")).toBeDefined();

      // Topics with hashtag
      expect(screen.getByText("#nextjs")).toBeDefined();
      expect(screen.getByText("#ai")).toBeDefined();
      expect(screen.getByText("#education")).toBeDefined();

      // Stars & Forks
      expect(screen.getByText("12")).toBeDefined();
      expect(screen.getByText("3")).toBeDefined();

      // Synced date formatted
      expect(screen.getByText(/Synced: Sep 14, 2026/i)).toBeDefined();

      // External link
      const link = screen.getByRole("link", { name: /View EduTrace repository on GitHub/i });
      expect(link.getAttribute("href")).toBe("https://github.com/Nachsyas/EduTrace");
      expect(link.getAttribute("target")).toBe("_blank");
      expect(link.getAttribute("rel")).toContain("noopener");
      expect(link.getAttribute("rel")).toContain("noreferrer");
    });

    it("renders fallback text gracefully when description or topics are empty", () => {
      const bareProject: GithubProjectItem = {
        id: "uuid-3",
        name: "MinimalRepo",
        full_name: "Nachsyas/MinimalRepo",
        description: "",
        html_url: "https://github.com/Nachsyas/MinimalRepo",
        homepage: "",
        language: "",
        stars: 0,
        forks: 0,
        topics: [],
        synced_at: "",
      };

      render(<GithubProjectCard project={bareProject} />);

      expect(screen.getByRole("heading", { name: "MinimalRepo" })).toBeDefined();
      expect(screen.getByText("Open source repository synchronized from verified GitHub profile.")).toBeDefined();
      expect(screen.getByText("Repository")).toBeDefined();
    });
  });

  describe("GithubProjectsSection", () => {
    it("renders initial projects directly without network loading", () => {
      render(<GithubProjectsSection initialProjects={[mockProject1, mockProject2]} />);

      expect(screen.getByRole("heading", { name: /GitHub Repositories/i })).toBeDefined();
      expect(screen.getByRole("heading", { name: "EduTrace" })).toBeDefined();
      expect(screen.getByRole("heading", { name: "GoStream" })).toBeDefined();
      expect(screen.queryByTestId("github-projects-loading")).toBeNull();
      expect(screen.queryByTestId("github-projects-error")).toBeNull();
    });

    it("fetches and renders projects successfully from API", async () => {
      global.fetch = vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => ({
          status: "success",
          data: [mockProject1, mockProject2],
        }),
      } as Response);

      render(<GithubProjectsSection />);

      // Shows loading skeleton initially
      expect(screen.getByTestId("github-projects-loading")).toBeDefined();

      // Wait for projects to render
      await waitFor(() => {
        expect(screen.getByTestId("github-projects-grid")).toBeDefined();
      });

      expect(screen.getByRole("heading", { name: "EduTrace" })).toBeDefined();
      expect(screen.getByRole("heading", { name: "GoStream" })).toBeDefined();
      expect(screen.queryByTestId("github-projects-loading")).toBeNull();
    });

    it("displays non-blocking error notice when API fails", async () => {
      global.fetch = vi.fn().mockRejectedValue(new Error("Network Error"));

      render(<GithubProjectsSection />);

      await waitFor(() => {
        expect(screen.getByTestId("github-projects-error")).toBeDefined();
      });

      expect(screen.getByText(/GitHub projects are temporarily unavailable\./i)).toBeDefined();
      expect(screen.getByText(/Canonical projects and case studies remain fully accessible above\./i)).toBeDefined();
      expect(screen.queryByTestId("github-projects-loading")).toBeNull();
    });

    it("displays empty state when API returns empty array", async () => {
      global.fetch = vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => ({
          status: "success",
          data: [],
        }),
      } as Response);

      render(<GithubProjectsSection />);

      await waitFor(() => {
        expect(screen.getByText(/No synchronized repositories found at this moment\./i)).toBeDefined();
      });
    });
  });

  describe("API Client: getGithubProjects", () => {
    it("fetches projects with page and limit parameters", async () => {
      const mockFetch = vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => ({
          status: "success",
          data: [mockProject1],
        }),
      });
      global.fetch = mockFetch;

      const data = await getGithubProjects(2, 10);
      expect(data).toHaveLength(1);
      expect(data[0].name).toBe("EduTrace");
      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining("/api/v1/github/projects?page=2&limit=10"),
        expect.objectContaining({ method: "GET" })
      );
    });

    it("handles non-ok response with server error message", async () => {
      global.fetch = vi.fn().mockResolvedValue({
        ok: false,
        status: 500,
        json: async () => ({
          error: {
            message: "database query failure",
          },
        }),
      });

      await expect(getGithubProjects(1, 20)).rejects.toThrow("database query failure");
    });
  });

  describe("Unified /work Integration & Fault-Tolerance", () => {
    it("renders canonical projects AND handles GitHub API failure gracefully", async () => {
      // Simulate backend down
      global.fetch = vi.fn().mockRejectedValue(new Error("Connection refused"));

      render(
        <PortfolioUIProvider>
          <WorkPage />
        </PortfolioUIProvider>
      );

      // Canonical projects MUST remain visible
      expect(screen.getByRole("heading", { name: "EduTrace" })).toBeDefined();
      expect(screen.getByRole("heading", { name: "GDGOC E-Commerce" })).toBeDefined();

      // Wait for GitHub section error notice
      await waitFor(() => {
        expect(screen.getByTestId("github-projects-error")).toBeDefined();
      });

      // Canonical work is NOT replaced or removed
      expect(screen.getByText("WORK // SELECTED PROJECTS")).toBeDefined();
      expect(screen.getByRole("heading", { name: /Selected Systems & Engineering Projects/i })).toBeDefined();
    });

    it("renders canonical projects AND synchronized GitHub repositories together", async () => {
      global.fetch = vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => ({
          status: "success",
          data: [mockProject2],
        }),
      });

      render(
        <PortfolioUIProvider>
          <WorkPage />
        </PortfolioUIProvider>
      );

      // Canonical projects present
      expect(screen.getByRole("heading", { name: "GDGOC E-Commerce" })).toBeDefined();

      // Wait for GitHub project to appear
      await waitFor(() => {
        expect(screen.getByRole("heading", { name: "GoStream" })).toBeDefined();
      });

      expect(screen.getByText("#golang")).toBeDefined();
    });
  });
});
