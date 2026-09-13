import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import SmokePage from "../app/page";

describe("Phase 0 Smoke Test Page", () => {
  it("renders the portfolio owner's name and role from validated data", () => {
    render(<SmokePage />);
    expect(screen.getByText("Nachsyas Arham Mumtaz Nashohi")).toBeDefined();
    expect(screen.getByText("Software Engineer")).toBeDefined();
    expect(screen.getByText(/PHASE 0 — FOUNDATION READY/i)).toBeDefined();
    expect(screen.getByText(/Zod Verified/i)).toBeDefined();
  });
});
