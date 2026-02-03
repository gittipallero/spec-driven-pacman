import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { MainPage } from "./MainPage";

describe("MainPage", () => {
  it("renders the game title", () => {
    render(<MainPage />);
    expect(screen.getByText("Spec-Driven-Pacman")).toBeInTheDocument();
  });

  it("renders the start game button", () => {
    render(<MainPage />);
    expect(screen.getByRole("button")).toBeInTheDocument();
    expect(screen.getByText("START GAME")).toBeInTheDocument();
  });

  it("renders with CRT effect wrapper", () => {
    const { container } = render(<MainPage />);
    const crtContainer = container.querySelector('[class*="crtContainer"]');
    expect(crtContainer).toBeInTheDocument();
  });

  it("has proper heading hierarchy", () => {
    render(<MainPage />);
    const heading = screen.getByRole("heading", { level: 1 });
    expect(heading).toHaveTextContent("Spec-Driven-Pacman");
  });

  it("renders scan lines overlay for CRT effect", () => {
    const { container } = render(<MainPage />);
    const scanLines = container.querySelector('[class*="scanLines"]');
    expect(scanLines).toBeInTheDocument();
  });
});

