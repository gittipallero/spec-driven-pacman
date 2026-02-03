import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { ArcadeTitle } from "./ArcadeTitle";

describe("ArcadeTitle", () => {
  it("renders the game title text", () => {
    render(<ArcadeTitle />);
    expect(screen.getByText("Spec-Driven-Pacman")).toBeInTheDocument();
  });

  it("renders as an h1 heading element", () => {
    render(<ArcadeTitle />);
    const heading = screen.getByRole("heading", { level: 1 });
    expect(heading).toBeInTheDocument();
    expect(heading).toHaveTextContent("Spec-Driven-Pacman");
  });

  it("applies arcade styling class", () => {
    render(<ArcadeTitle />);
    const heading = screen.getByRole("heading", { level: 1 });
    expect(heading.className).toContain("title");
  });

  it("allows custom title text via prop", () => {
    render(<ArcadeTitle title="Custom Title" />);
    expect(screen.getByText("Custom Title")).toBeInTheDocument();
  });
});

