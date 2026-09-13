import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import SmokePage from "../app/page";

describe("Phase 1 Production Portfolio Page", () => {
  it("renders the portfolio owner's name, role, and key navigation links", () => {
    render(<SmokePage />);
    expect(screen.getAllByText(/NACHSYAS ARHAM/i).length).toBeGreaterThan(0);
    expect(screen.getAllByText("Software Engineer").length).toBeGreaterThan(0);
    expect(screen.getAllByText(/Quick Review/i).length).toBeGreaterThan(0);
    expect(screen.getByText("Selected Engineering Work")).toBeDefined();
  });
});
