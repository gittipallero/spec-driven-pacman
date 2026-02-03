import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import App from "./App";

describe("App", () => {
  it("renders the main page", () => {
    render(<App />);
    expect(screen.getByText("Spec-Driven-Pacman")).toBeInTheDocument();
  });

  it("renders the start game button", () => {
    render(<App />);
    expect(screen.getByText("START GAME")).toBeInTheDocument();
  });
});

