import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import CaseStudyPage, { generateMetadata, generateStaticParams } from "../app/projects/[slug]/page";

describe("Dynamic Case Study Route (/projects/[slug])", () => {
  it("pre-renders all static parameters for approved projects", async () => {
    const params = await generateStaticParams();
    const slugs = params.map((p) => p.slug);
    expect(slugs).toContain("edutrace");
    expect(slugs).toContain("gdgoc-ecommerce");
    expect(slugs).toContain("maritime-ai-dashboard");
    expect(slugs).toContain("smart-kitchen");
  });

  it("generates correct metadata for a valid project slug", async () => {
    const meta = await generateMetadata({
      params: Promise.resolve({ slug: "edutrace" }),
    });
    expect(meta.title).toContain("EduTrace");
  });

  it("renders verified case study sections for EduTrace", async () => {
    const pageComponent = await CaseStudyPage({
      params: Promise.resolve({ slug: "edutrace" }),
    });

    render(pageComponent);

    expect(screen.getByRole("heading", { name: "EduTrace" })).toBeDefined();
    expect(screen.getByText(/Decentralized academic record ledger/i)).toBeDefined();
    expect(screen.getByText("01 // The Problem & Context")).toBeDefined();
    expect(screen.getByText("02 // Architectural Solution")).toBeDefined();
    expect(screen.getByText("03 // Technology Stack")).toBeDefined();
    expect(screen.getByText("Back to Selected Work")).toBeDefined();
  });
});
