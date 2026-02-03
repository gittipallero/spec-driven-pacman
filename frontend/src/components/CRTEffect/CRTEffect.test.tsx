import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { CRTEffect } from "./CRTEffect";

describe("CRTEffect", () => {
  it("renders children content", () => {
    render(
      <CRTEffect>
        <span>Test Content</span>
      </CRTEffect>
    );
    expect(screen.getByText("Test Content")).toBeInTheDocument();
  });

  it("wraps content in a container with CRT styling", () => {
    const { container } = render(
      <CRTEffect>
        <span data-testid="child">Content</span>
      </CRTEffect>
    );
    // CSS modules mangle class names, so check for the wrapper structure
    const crtWrapper = container.firstElementChild;
    expect(crtWrapper).toBeInTheDocument();
    expect(crtWrapper?.className).toMatch(/crtContainer/);
  });

  it("includes scan lines overlay element", () => {
    const { container } = render(
      <CRTEffect>
        <span>Content</span>
      </CRTEffect>
    );
    const scanLines = container.querySelector('[class*="scanLines"]');
    expect(scanLines).toBeInTheDocument();
  });

  it("has proper aria attributes for decoration", () => {
    const { container } = render(
      <CRTEffect>
        <span>Content</span>
      </CRTEffect>
    );
    const scanLines = container.querySelector('[class*="scanLines"]');
    expect(scanLines).toHaveAttribute("aria-hidden", "true");
  });
});

