import { describe, it, expect, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { MainPage } from "../../src/components/MainPage";

describe("MainPage Integration", () => {
  it("displays all main page elements together", () => {
    render(<MainPage />);

    // Title is visible
    expect(screen.getByText("Spec-Driven-Pacman")).toBeInTheDocument();

    // Start button is visible
    expect(screen.getByText("START GAME")).toBeInTheDocument();
  });

  it("has correct visual hierarchy", () => {
    render(<MainPage />);

    // Title should be h1
    const title = screen.getByRole("heading", { level: 1 });
    expect(title).toBeInTheDocument();

    // Button should be interactive
    const button = screen.getByRole("button");
    expect(button).toBeInTheDocument();
  });

  it("button click triggers navigation callback", () => {
    const handleStart = vi.fn();
    render(<MainPage onStartGame={handleStart} />);

    const button = screen.getByRole("button");
    fireEvent.click(button);

    expect(handleStart).toHaveBeenCalledTimes(1);
  });

  it("renders with dark background (CRT effect)", () => {
    const { container } = render(<MainPage />);
    const crtContainer = container.querySelector('[class*="crtContainer"]');
    expect(crtContainer).toBeInTheDocument();
  });

  it("is accessible - all interactive elements can be focused", () => {
    render(<MainPage />);

    const button = screen.getByRole("button");
    button.focus();
    expect(document.activeElement).toBe(button);
  });
});

