import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import React from "react";
import {
  AskArhamLauncher,
  AskArhamPanel,
  SuggestedQuestions,
  AskArhamAnswer,
  EvidenceList,
  SourceList,
  SafeActionButtons,
} from "../features/ask-arham";
import type { GroundedResponse, PublicEvidenceItem, SourceCitation, SafeAction } from "../features/ask-arham/ask-arham.types";

const mockRouterPush = vi.fn();
vi.mock("next/navigation", () => ({
  useRouter: () => ({
    push: mockRouterPush,
  }),
}));

describe("Ask Arham AI UI Components", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe("AskArhamLauncher", () => {
    it("renders launcher button with accessible label and Alt+A badge", () => {
      const onClick = vi.fn();
      render(<AskArhamLauncher isOpen={false} onClick={onClick} />);

      const button = screen.getByRole("button", { name: /Open Ask Arham AI Reviewer Copilot/i });
      expect(button).toBeDefined();
      expect(screen.getByText("Ask Arham AI")).toBeDefined();
      expect(screen.getByText("Alt+A")).toBeDefined();

      fireEvent.click(button);
      expect(onClick).toHaveBeenCalledTimes(1);
    });

    it("handles Alt+A keyboard shortcut to toggle panel", () => {
      const onClick = vi.fn();
      render(<AskArhamLauncher isOpen={false} onClick={onClick} />);

      fireEvent.keyDown(window, { key: "a", altKey: true });
      expect(onClick).toHaveBeenCalledTimes(1);
    });
  });

  describe("SuggestedQuestions", () => {
    it("renders recruiter-focused questions and triggers onSelect", () => {
      const onSelect = vi.fn();
      render(<SuggestedQuestions onSelect={onSelect} disabled={false} />);

      expect(screen.getByText("Suggested Reviewer Questions")).toBeDefined();
      const questionBtn = screen.getByText("What evidence shows Go experience?");
      expect(questionBtn).toBeDefined();

      fireEvent.click(questionBtn);
      expect(onSelect).toHaveBeenCalledWith("What evidence shows Go experience?");
    });
  });

  describe("EvidenceList (Zero Similarity Scores Rule)", () => {
    const mockEvidence: PublicEvidenceItem[] = [
      {
        id: "E1",
        title: "EduTrace - Clean Architecture",
        kind: "github",
        repository: "Nachsyas/EduTrace",
        excerpt: "Domain layer contains core business logic without external frameworks.",
      },
      {
        id: "E2",
        title: "Software Engineer Profile",
        kind: "portfolio",
        excerpt: "Backend engineering expertise in Go and PostgreSQL.",
      },
    ];

    it("renders evidence items without any similarity scores (Correction 17 & 48)", () => {
      render(<EvidenceList evidence={mockEvidence} />);

      // Accordion toggle button
      const toggle = screen.getByRole("button", { name: /Verified Evidence/i });
      expect(toggle).toBeDefined();

      // Expand accordion
      fireEvent.click(toggle);

      expect(screen.getByText(/EduTrace - Clean Architecture/i)).toBeDefined();
      expect(screen.getByText(/Software Engineer Profile/i)).toBeDefined();
      expect(screen.getByText(/Domain layer contains core business logic without external frameworks./i)).toBeDefined();

      // Verify ZERO similarity scores are present
      expect(screen.queryByText(/similarity/i)).toBeNull();
      expect(screen.queryByText(/score/i)).toBeNull();
      expect(screen.queryByText(/0\.\d+/)).toBeNull();
    });
  });

  describe("SourceList (GitHub & Canonical Portfolio Citations)", () => {
    const mockSources: SourceCitation[] = [
      {
        id: "E1",
        kind: "github",
        label: "EduTrace Repository",
        url: "https://github.com/Nachsyas/EduTrace/blob/main/README.md",
        repository: "Nachsyas/EduTrace",
      },
      {
        id: "E2",
        kind: "portfolio",
        label: "Academic Journey & Formative Base",
        // Note: URL is optional for canonical sources (Correction 14 & 15)
      },
    ];

    it("renders immutable GitHub links and safe canonical labels without local file paths", () => {
      render(<SourceList sources={mockSources} />);

      expect(screen.getByText("Certified Sources")).toBeDefined();

      // GitHub source has an anchor link
      const githubLink = screen.getByRole("link", { name: /EduTrace Repository/i });
      expect(githubLink.getAttribute("href")).toBe("https://github.com/Nachsyas/EduTrace/blob/main/README.md");
      expect(githubLink.getAttribute("target")).toBe("_blank");
      expect(githubLink.getAttribute("rel")).toContain("noreferrer");

      // Canonical portfolio citation has no raw local paths
      expect(screen.getByText("Academic Journey & Formative Base")).toBeDefined();
      expect(screen.queryByText(/data\/journey\/journey\.json/i)).toBeNull();
      expect(screen.queryByText(/\/Users\/user/i)).toBeNull();
    });
  });

  describe("SafeActionButtons", () => {
    const mockActions: SafeAction[] = [
      {
        id: "view-project-edutrace",
        label: "View EduTrace Case Study",
      },
      {
        id: "go-to-skills",
        label: "Explore Skills Matrix",
      },
    ];

    it("renders verified action buttons and invokes router push or scroll", () => {
      const onActionTriggered = vi.fn();
      render(<SafeActionButtons actions={mockActions} onActionTriggered={onActionTriggered} />);

      expect(screen.getByText("Contextual Actions")).toBeDefined();

      const btn1 = screen.getByRole("button", { name: /View EduTrace Case Study/i });
      fireEvent.click(btn1);
      expect(onActionTriggered).toHaveBeenCalled();
      expect(mockRouterPush).toHaveBeenCalledWith("/projects/edutrace");

      const btn2 = screen.getByRole("button", { name: /Explore Skills Matrix/i });
      fireEvent.click(btn2);
      expect(onActionTriggered).toHaveBeenCalled();
      expect(mockRouterPush).toHaveBeenCalledWith("/#skills");
    });
  });

  describe("AskArhamAnswer", () => {
    const mockResponse: GroundedResponse = {
      status: "supported",
      answer: "Nachsyas Arham Mumtaz Nashohi is a Software Engineer specializing in Go.",
      segments: [
        {
          text: "Nachsyas Arham Mumtaz Nashohi is a Software Engineer specializing in Go.",
          evidence_ids: ["E1"],
        },
      ],
      evidence: [
        {
          id: "E1",
          title: "Profile Information",
          kind: "portfolio",
          excerpt: "Software Engineer with expertise in Go, Clean Architecture, and PostgreSQL.",
        },
      ],
      sources: [
        {
          id: "E1",
          kind: "portfolio",
          label: "Portfolio Profile",
        },
      ],
      actions: [
        {
          id: "go-to-skills",
          label: "Explore Skills",
        },
      ],
    };

    it("renders status badge, grounded answer segments, sources, and actions", () => {
      render(<AskArhamAnswer response={mockResponse} onActionTriggered={vi.fn()} />);

      expect(screen.getByText("Verified Evidence")).toBeDefined();
      expect(screen.getByText("Nachsyas Arham Mumtaz Nashohi is a Software Engineer specializing in Go.")).toBeDefined();
      expect(screen.getByText("[E1]")).toBeDefined();
      expect(screen.getByText("Portfolio Profile")).toBeDefined();
      expect(screen.getByText("Explore Skills")).toBeDefined();
    });

    it("renders insufficient evidence badge appropriately", () => {
      const insufficientResponse: GroundedResponse = {
        status: "insufficient_evidence",
        answer: "I couldn't find verified evidence of production Kubernetes deployment in approved sources.",
        segments: [
          {
            text: "I couldn't find verified evidence of production Kubernetes deployment in approved sources.",
            evidence_ids: [],
          },
        ],
        evidence: [],
        sources: [],
        actions: [],
      };

      render(<AskArhamAnswer response={insufficientResponse} onActionTriggered={vi.fn()} />);
      expect(screen.getByText("Insufficient Evidence")).toBeDefined();
      expect(screen.getByText(/I couldn't find verified evidence/i)).toBeDefined();
    });
  });

  describe("AskArhamPanel & Accessibility", () => {
    it("renders panel with restrained aria-live announcer and handles Escape key", () => {
      const onClose = vi.fn();
      render(
        <AskArhamPanel
          isOpen={true}
          onClose={onClose}
        />
      );

      // Panel header
      expect(screen.getByRole("dialog")).toBeDefined();
      expect(screen.getByRole("heading", { name: /Ask Arham AI/i })).toBeDefined();

      // aria-live announcer
      const liveRegions = document.querySelectorAll("[aria-live]");
      expect(liveRegions.length).toBeGreaterThan(0);
      expect(liveRegions[0].getAttribute("aria-live")).toBe("polite");

      // Press Escape key
      fireEvent.keyDown(window, { key: "Escape" });
      expect(onClose).toHaveBeenCalledTimes(1);
    });
  });

  describe("XSS Safety Boundary", () => {
    it("renders untrusted HTML in answers and evidence as inert text", () => {
      const xssPayload = "<script>alert('pwned')</script><img src=x onerror=alert(1)>";
      const xssResponse: GroundedResponse = {
        status: "supported",
        answer: xssPayload,
        segments: [
          {
            text: xssPayload,
            evidence_ids: ["E1"],
          },
        ],
        evidence: [
          {
            id: "E1",
            title: "Untrusted Evidence",
            kind: "github",
            excerpt: xssPayload,
          },
        ],
        sources: [],
        actions: [],
      };

      const { container } = render(
        <AskArhamAnswer response={xssResponse} onActionTriggered={vi.fn()} />
      );

      // Verify no script tags are created in the DOM
      expect(container.querySelectorAll("script")).toHaveLength(0);
      expect(container.querySelectorAll("img[src='x']")).toHaveLength(0);
      // Plain text shows the escaped markup characters
      expect(screen.getAllByText((content) => content.includes("<script>")).length).toBeGreaterThan(0);
    });
  });
});
