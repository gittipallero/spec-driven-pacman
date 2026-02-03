import { describe, it, expect, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { StartButton } from "./StartButton";

describe("StartButton", () => {
  it("renders the start game text", () => {
    render(<StartButton />);
    expect(screen.getByText("START GAME")).toBeInTheDocument();
  });

  it("renders as a button element", () => {
    render(<StartButton />);
    const button = screen.getByRole("button");
    expect(button).toBeInTheDocument();
    expect(button).toHaveTextContent("START GAME");
  });

  it("applies arcade styling class", () => {
    render(<StartButton />);
    const button = screen.getByRole("button");
    expect(button.className).toContain("button");
  });

  it("calls onClick handler when clicked", () => {
    const handleClick = vi.fn();
    render(<StartButton onClick={handleClick} />);

    const button = screen.getByRole("button");
    fireEvent.click(button);

    expect(handleClick).toHaveBeenCalledTimes(1);
  });

  it("allows custom button text via prop", () => {
    render(<StartButton label="PLAY NOW" />);
    expect(screen.getByText("PLAY NOW")).toBeInTheDocument();
  });

  it("is accessible with proper focus styles", () => {
    render(<StartButton />);
    const button = screen.getByRole("button");
    button.focus();
    expect(document.activeElement).toBe(button);
  });
});

