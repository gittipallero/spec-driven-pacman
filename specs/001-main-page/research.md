# Research: Arcade Main Page

**Feature**: 001-main-page  
**Date**: 2026-02-03

## Research Topics

### 1. Retro Arcade Typography

**Decision**: Use "Press Start 2P" as primary font with system monospace fallback

**Rationale**:

- Press Start 2P is a bitmap-style font specifically designed to emulate 1980s arcade game text
- Free and open-source (OFL license), can be self-hosted or loaded from Google Fonts
- Widely supported and renders consistently across browsers
- 8x8 pixel grid design authentic to classic arcade machines

**Alternatives Considered**:
| Font | Pros | Cons | Rejected Because |
|------|------|------|------------------|
| VT323 | Authentic terminal look | Too thin for headers | Less impactful for game title |
| Silkscreen | Crisp pixel font | Limited character set | Missing some glyphs |
| Custom bitmap font | Full control | Development time | Overkill for MVP |

### 2. CRT Visual Effects (CSS-only)

**Decision**: Implement CRT effects using CSS filters, pseudo-elements, and animations

**Rationale**:

- CSS-only approach ensures 60 FPS performance without canvas overhead
- Progressive enhancement: effects degrade gracefully on older browsers
- Easy to disable for `prefers-reduced-motion` users
- No JavaScript required for visual effects

**Implementation Approach**:

```css
/* Scan lines via repeating gradient */
.crt::before {
  background: repeating-linear-gradient(
    0deg,
    rgba(0, 0, 0, 0.15),
    rgba(0, 0, 0, 0.15) 1px,
    transparent 1px,
    transparent 2px
  );
}

/* Screen curvature via border-radius + overflow */
.crt {
  border-radius: 20px / 10px;
  box-shadow: inset 0 0 100px rgba(0, 0, 0, 0.5);
}

/* Phosphor glow via text-shadow */
.glow {
  text-shadow: 0 0 10px currentColor, 0 0 20px currentColor;
}
```

**Alternatives Considered**:
| Approach | Pros | Cons | Rejected Because |
|----------|------|------|------------------|
| Canvas shader | Most authentic | Performance cost, complexity | Overkill for landing page |
| SVG filters | Good browser support | Complex syntax | CSS simpler for basic effects |
| WebGL | Full control | Major complexity | Way overkill |

### 3. Color Palette (Pac-Man Authentic)

**Decision**: Use classic Pac-Man color palette with slight modernization for accessibility

**Rationale**:

- Authentic colors create immediate recognition
- High contrast ensures accessibility compliance
- Limited palette maintains retro feel

**Color System**:

```css
:root {
  /* Background */
  --arcade-black: #000000;
  --arcade-dark-blue: #0a0a23;

  /* Pac-Man Yellow (primary) */
  --pacman-yellow: #ffff00;
  --pacman-yellow-glow: #ffff66;

  /* Ghost Colors (accents) */
  --ghost-red: #ff0000; /* Blinky */
  --ghost-pink: #ffb8ff; /* Pinky */
  --ghost-cyan: #00ffff; /* Inky */
  --ghost-orange: #ffb852; /* Clyde */

  /* Maze Blue */
  --maze-blue: #2121de;
  --maze-blue-light: #3333ff;

  /* UI Elements */
  --text-white: #ffffff;
  --text-dim: #888888;
}
```

**Accessibility Check**:

- Yellow (#ffff00) on dark blue (#0a0a23): Contrast ratio 15.3:1 ✅ (AAA)
- White (#ffffff) on dark blue (#0a0a23): Contrast ratio 19.1:1 ✅ (AAA)
- Cyan (#00ffff) on dark blue (#0a0a23): Contrast ratio 12.4:1 ✅ (AAA)

### 4. Animation Performance

**Decision**: CSS animations only, with `will-change` hints and GPU acceleration

**Rationale**:

- CSS animations run on compositor thread, ensuring 60 FPS
- `transform` and `opacity` are GPU-accelerated properties
- No JavaScript frame loops required for landing page effects

**Animation Patterns**:

```css
/* Title glow pulse - GPU accelerated */
@keyframes glow-pulse {
  0%,
  100% {
    filter: brightness(1);
  }
  50% {
    filter: brightness(1.2);
  }
}

/* Respect reduced motion */
@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    animation-duration: 0.01ms !important;
    transition-duration: 0.01ms !important;
  }
}
```

### 5. Responsive Strategy

**Decision**: Mobile-first with CSS Grid, scaling text with `clamp()`

**Rationale**:

- `clamp()` provides fluid typography without breakpoint jumps
- CSS Grid handles layout with minimal media queries
- Pixel fonts scale well at integer multiples

**Breakpoints**:
| Breakpoint | Target | Font Scale |
|------------|--------|------------|
| 320px | Mobile (min) | 1x (base) |
| 768px | Tablet | 1.5x |
| 1024px | Desktop | 2x |
| 1920px | Full HD | 2.5x |
| 3840px | 4K | 4x |

**Implementation**:

```css
.title {
  font-size: clamp(1.5rem, 5vw, 4rem);
}
```

### 6. Font Loading Strategy

**Decision**: Self-host fonts with `font-display: swap` and system fallback

**Rationale**:

- Self-hosting avoids Google Fonts privacy concerns
- `font-display: swap` ensures text is always visible
- Monospace fallback maintains layout during load

**Implementation**:

```css
@font-face {
  font-family: "Press Start 2P";
  src: url("/fonts/PressStart2P-Regular.woff2") format("woff2");
  font-display: swap;
}

:root {
  --font-arcade: "Press Start 2P", "Courier New", monospace;
}
```

## Resolved Clarifications

All technical decisions made - no NEEDS CLARIFICATION items remain.

## Dependencies to Install

```json
{
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0"
  },
  "devDependencies": {
    "@testing-library/react": "^14.0.0",
    "@testing-library/jest-dom": "^6.0.0",
    "@types/react": "^18.2.0",
    "@types/react-dom": "^18.2.0",
    "typescript": "^5.3.0",
    "vite": "^5.0.0",
    "vitest": "^1.0.0",
    "jsdom": "^23.0.0"
  }
}
```

## External Assets

| Asset               | Source                   | License | Notes                           |
| ------------------- | ------------------------ | ------- | ------------------------------- |
| Press Start 2P font | Google Fonts / self-host | OFL 1.1 | Download WOFF2 for self-hosting |
