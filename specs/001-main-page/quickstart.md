# Quickstart: Arcade Main Page

**Feature**: 001-main-page  
**Date**: 2026-02-03

## Prerequisites

- Node.js 20+ (LTS recommended)
- npm 10+ or pnpm 8+

## Setup

### 1. Create Frontend Project

```bash
# From repository root
cd frontend

# Initialize Vite + React + TypeScript project
npm create vite@latest . -- --template react-ts

# Install dependencies
npm install

# Install testing dependencies
npm install -D vitest @testing-library/react @testing-library/jest-dom jsdom
```

### 2. Configure TypeScript (strict mode)

Update `tsconfig.json`:

```json
{
  "compilerOptions": {
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true,
    "noUncheckedIndexedAccess": true
  }
}
```

### 3. Configure Vitest

Create `vitest.config.ts`:

```typescript
import { defineConfig } from "vitest/config";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  test: {
    environment: "jsdom",
    globals: true,
    setupFiles: ["./tests/setup.ts"],
  },
});
```

Create `tests/setup.ts`:

```typescript
import "@testing-library/jest-dom";
```

### 4. Add Font Assets

Download Press Start 2P font:

```bash
mkdir -p public/fonts
# Download from Google Fonts or:
# https://fonts.google.com/specimen/Press+Start+2P
# Place PressStart2P-Regular.woff2 in public/fonts/
```

### 5. Create CSS Variables

Create `src/styles/variables.css`:

```css
:root {
  /* Backgrounds */
  --arcade-black: #000000;
  --arcade-dark-blue: #0a0a23;

  /* Pac-Man Colors */
  --pacman-yellow: #ffff00;
  --pacman-yellow-glow: #ffff66;

  /* Ghost Colors */
  --ghost-red: #ff0000;
  --ghost-pink: #ffb8ff;
  --ghost-cyan: #00ffff;
  --ghost-orange: #ffb852;

  /* Maze */
  --maze-blue: #2121de;

  /* Text */
  --text-white: #ffffff;
  --text-dim: #888888;

  /* Typography */
  --font-arcade: "Press Start 2P", "Courier New", monospace;
}
```

## Development

### Run Development Server

```bash
npm run dev
```

Open http://localhost:5173 in your browser.

### Run Tests

```bash
# Run all tests
npm test

# Run tests in watch mode
npm test -- --watch

# Run tests with coverage
npm test -- --coverage
```

### Build for Production

```bash
npm run build
```

Output will be in `dist/` directory.

## Verification Checklist

After setup, verify:

- [ ] `npm run dev` starts without errors
- [ ] http://localhost:5173 shows default Vite page
- [ ] `npm test` runs without errors
- [ ] TypeScript strict mode catches type errors
- [ ] Font files are accessible at `/fonts/PressStart2P-Regular.woff2`

## Component Development Order

Implement components in this order (TDD):

1. **CSS Variables & Global Styles** - Foundation
2. **useReducedMotion hook** - Accessibility foundation
3. **ArcadeTitle component** - Core visual element (P1)
4. **CRTEffect component** - Visual wrapper (P2)
5. **StartButton component** - Navigation element (P3)
6. **MainPage component** - Composition of all components

For each component:

1. Write test first (Red)
2. Implement to pass test (Green)
3. Refactor if needed

## Troubleshooting

### Font not loading

- Check file exists at `public/fonts/PressStart2P-Regular.woff2`
- Verify `@font-face` declaration in CSS
- Check browser dev tools Network tab for 404

### Tests failing

- Ensure `vitest.config.ts` has `environment: 'jsdom'`
- Verify `tests/setup.ts` imports `@testing-library/jest-dom`

### TypeScript errors

- Run `npm run build` to see all type errors
- Ensure `strict: true` in `tsconfig.json`
