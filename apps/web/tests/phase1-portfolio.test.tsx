import { describe, it, expect, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import PortfolioApp from "../features/portfolio/PortfolioApp";
import ProjectCard from "../features/projects/ProjectCard";
import ExperienceSection from "../features/experience/ExperienceSection";
import EducationSection from "../features/education/EducationSection";
import SkillsSection from "../features/skills/SkillsSection";
import type { Profile, Project, Skill } from "arham-porto-schema";

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
  todo: ["TODO_USER_POSITIONING"],
};

const mockProjects: Project[] = [
  {
    id: "edutrace",
    title: "EduTrace",
    slug: "edutrace",
    summary: "Full-stack educational tracking system.",
    problem: "Educational institutions struggle with fragmented records.",
    solution: "A centralized, verified educational tracking platform.",
    role: [],
    contributions: [],
    technologies: ["TypeScript", "Next.js", "React"],
    githubUrl: "https://github.com/Nachsyas/EduTrace",
    demoUrl: null,
    image: null,
    featured: true,
    category: "Full-Stack",
    evidenceIds: [],
  },
  {
    id: "gdgoc-ecommerce",
    title: "GDGOC E-Commerce",
    slug: "gdgoc-ecommerce",
    summary: "Modular backend service for transactions.",
    problem: "Monolithic backends face coupling.",
    solution: "Decoupled Go REST API with clean architecture.",
    role: [],
    contributions: [],
    technologies: ["Go", "PostgreSQL", "Clean Architecture"],
    githubUrl: "https://github.com/Nachsyas/gdgoc-ecommerce",
    demoUrl: null,
    image: null,
    featured: true,
    category: "Backend",
    evidenceIds: [],
  },
];

const mockSkills: Skill[] = [
  {
    id: "backend-go",
    name: "Backend Engineering with Go",
    category: "Backend",
    claim: "Building scalable backend services and REST APIs with modern Go.",
    evidenceIds: [],
  },
  {
    id: "frontend-nextjs",
    name: "Frontend Engineering with Next.js",
    category: "Frontend",
    claim: "Developing accessible user interfaces.",
    evidenceIds: [],
  },
];

describe("Phase 1: Reviewer-First Static Portfolio", () => {
  it("renders the Hero section with full name and primary role", () => {
    render(
      <PortfolioApp
        profile={mockProfile}
        projects={mockProjects}
        skills={mockSkills}
        experiences={[]}
        education={[]}
      />
    );

    expect(screen.getAllByText(/NACHSYAS ARHAM/i).length).toBeGreaterThan(0);
    expect(screen.getAllByText(/MUMTAZ NASHOHI/i).length).toBeGreaterThan(0);
    expect(screen.getAllByText("Software Engineer").length).toBeGreaterThan(0);
    expect(screen.getByText("#Backend")).toBeDefined();
    expect(screen.getByText("#Full-Stack")).toBeDefined();
  });

  it("opens the 60-Second Quick Review drawer on trigger and does not leak TODO strings", () => {
    render(
      <PortfolioApp
        profile={mockProfile}
        projects={mockProjects}
        skills={mockSkills}
        experiences={[]}
        education={[]}
      />
    );

    const reviewButton = screen.getAllByRole("button", { name: /quick review/i })[0];
    fireEvent.click(reviewButton);

    expect(screen.getByRole("dialog")).toBeDefined();
    expect(screen.getByRole("heading", { name: "60-Second Quick Review" })).toBeDefined();
    expect(screen.getByText("Candidate Profile")).toBeDefined();

    // Verify zero TODO_USER strings are rendered to users
    expect(screen.queryByText(/TODO_USER/i)).toBeNull();
  });

  it("filters project cards dynamically upon selecting a category tab", () => {
    render(
      <PortfolioApp
        profile={mockProfile}
        projects={mockProjects}
        skills={mockSkills}
        experiences={[]}
        education={[]}
      />
    );

    // Initial state: all projects rendered
    expect(screen.getByRole("heading", { name: "EduTrace" })).toBeDefined();
    expect(screen.getByRole("heading", { name: "GDGOC E-Commerce" })).toBeDefined();

    // Click Backend filter tab
    const backendTab = screen.getByRole("tab", { name: "Backend" });
    fireEvent.click(backendTab);

    // GDGOC E-Commerce should remain, EduTrace (Full-Stack) should be filtered out
    expect(screen.getByRole("heading", { name: "GDGOC E-Commerce" })).toBeDefined();
    expect(screen.queryByRole("heading", { name: "EduTrace" })).toBeNull();

    // Click All filter tab
    const allTab = screen.getByRole("tab", { name: "All" });
    fireEvent.click(allTab);
    expect(screen.getByRole("heading", { name: "EduTrace" })).toBeDefined();
  });

  it("renders skills grounded in claims without arbitrary percentage numbers", () => {
    const { container } = render(
      <SkillsSection skills={mockSkills} projects={mockProjects} />
    );

    expect(screen.getByText("Backend Engineering with Go")).toBeDefined();
    expect(screen.getByText("Building scalable backend services and REST APIs with modern Go.")).toBeDefined();

    // Assert that no arbitrary percentage numbers (e.g. 95%, 90%) are present
    const content = container.textContent || "";
    expect(content).not.toMatch(/\d+%/);
  });

  it("handles empty Experience gracefully without fabricated claims", () => {
    render(<ExperienceSection experiences={[]} />);

    expect(screen.getByText("Experience Records in Verification")).toBeDefined();
    expect(screen.getByText(/undergoing verification against primary source documentation/i)).toBeDefined();
  });

  it("handles empty Education gracefully without fabricated claims", () => {
    render(<EducationSection education={[]} />);

    expect(screen.getByText("Academic Background Verification")).toBeDefined();
    expect(screen.getByText(/undergoing verification before public listing/i)).toBeDefined();
  });
});
