# Tasks: Arcade Main Page

**Input**: Design documents from `/specs/001-main-page/`  
**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, data-model.md ✅, quickstart.md ✅

**Tests**: Included per Constitution Principle II (Test-Driven Development)

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Web app**: `frontend/src/` for source code
- All paths relative to repository root

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Create frontend project structure with Vite + React + TypeScript in frontend/
- [x] T002 Install dependencies (react, react-dom) and dev dependencies (vitest, @testing-library/react, @testing-library/jest-dom, jsdom) in frontend/package.json
- [x] T003 [P] Configure TypeScript strict mode in frontend/tsconfig.json
- [x] T004 [P] Configure Vitest with jsdom environment in frontend/vitest.config.ts
- [x] T005 [P] Create test setup file in frontend/tests/setup.ts
- [x] T006 [P] Download and add Press Start 2P font to frontend/public/fonts/PressStart2P-Regular.woff2

**Checkpoint**: Project scaffolding complete, ready to run `npm run dev` and `npm test`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core CSS infrastructure that ALL user stories depend on

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T007 Create CSS custom properties (colors, typography) in frontend/src/styles/variables.css
- [x] T008 [P] Create global styles with font-face declarations in frontend/src/styles/global.css
- [x] T009 [P] Create keyframe animations (glow-pulse) with prefers-reduced-motion support in frontend/src/styles/animations.css
- [x] T010 Create useReducedMotion accessibility hook in frontend/src/hooks/useReducedMotion.ts
- [x] T011 Update frontend/src/main.tsx to import global styles

**Checkpoint**: Foundation ready - CSS variables, fonts, and animations available for all components

---

## Phase 3: User Story 1 - View Game Title (Priority: P1) 🎯 MVP

**Goal**: Display "Spec-Driven-Pacman" title with arcade-style font and glow effects

**Independent Test**: Load main page and verify title "Spec-Driven-Pacman" is visible with retro arcade styling

### Tests for User Story 1

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T012 [P] [US1] Create ArcadeTitle component test in frontend/src/components/ArcadeTitle/ArcadeTitle.test.tsx

### Implementation for User Story 1

- [x] T013 [US1] Create ArcadeTitle component in frontend/src/components/ArcadeTitle/ArcadeTitle.tsx
- [x] T014 [US1] Create ArcadeTitle styles with glow effect in frontend/src/components/ArcadeTitle/ArcadeTitle.module.css
- [x] T015 [US1] Create component index export in frontend/src/components/ArcadeTitle/index.ts

**Checkpoint**: ArcadeTitle component renders "Spec-Driven-Pacman" with arcade font and glow - User Story 1 complete and independently testable

---

## Phase 4: User Story 2 - Retro Arcade Atmosphere (Priority: P2)

**Goal**: Wrap content in CRT-style visual effects (scan lines, screen curvature, dark background)

**Independent Test**: Load main page and verify CRT effects are visible (scan lines overlay, dark background, screen glow)

### Tests for User Story 2

- [x] T016 [P] [US2] Create CRTEffect component test in frontend/src/components/CRTEffect/CRTEffect.test.tsx

### Implementation for User Story 2

- [x] T017 [US2] Create CRTEffect wrapper component in frontend/src/components/CRTEffect/CRTEffect.tsx
- [x] T018 [US2] Create CRTEffect styles (scan lines, curvature, glow) in frontend/src/components/CRTEffect/CRTEffect.module.css
- [x] T019 [US2] Create component index export in frontend/src/components/CRTEffect/index.ts

**Checkpoint**: CRTEffect component wraps children with retro CRT visual effects - User Story 2 complete and independently testable

---

## Phase 5: User Story 3 - Navigate to Game (Priority: P3)

**Goal**: Display "START GAME" button with arcade styling and hover effects

**Independent Test**: Load main page, verify Start Game button is visible and responds to click/hover

### Tests for User Story 3

- [x] T020 [P] [US3] Create StartButton component test in frontend/src/components/StartButton/StartButton.test.tsx

### Implementation for User Story 3

- [x] T021 [US3] Create StartButton component in frontend/src/components/StartButton/StartButton.tsx
- [x] T022 [US3] Create StartButton styles (arcade theme, hover glow) in frontend/src/components/StartButton/StartButton.module.css
- [x] T023 [US3] Create component index export in frontend/src/components/StartButton/index.ts

**Checkpoint**: StartButton component renders with arcade styling and responds to interaction - User Story 3 complete

---

## Phase 6: Integration - MainPage Composition

**Goal**: Compose all components into the final MainPage

### Tests for MainPage

- [x] T024 [P] Create MainPage component test in frontend/src/components/MainPage/MainPage.test.tsx
- [x] T025 [P] Create integration test for full page in frontend/tests/integration/main-page.test.tsx

### Implementation for MainPage

- [x] T026 Create MainPage component (composes ArcadeTitle, CRTEffect, StartButton) in frontend/src/components/MainPage/MainPage.tsx
- [x] T027 Create MainPage layout styles in frontend/src/components/MainPage/MainPage.module.css
- [x] T028 Create component index export in frontend/src/components/MainPage/index.ts
- [x] T029 Update App.tsx to render MainPage in frontend/src/App.tsx
- [x] T030 Create App test in frontend/src/App.test.tsx

**Checkpoint**: Complete main page renders with all components - Full feature testable

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Responsive design, accessibility verification, performance validation

### Responsive Design

- [x] T031 [P] Add responsive styles (320px - 3840px) to frontend/src/components/ArcadeTitle/ArcadeTitle.module.css
- [x] T032 [P] Add responsive styles to frontend/src/components/MainPage/MainPage.module.css
- [x] T033 [P] Add responsive styles to frontend/src/components/StartButton/StartButton.module.css

### Accessibility & Performance

- [x] T034 Verify color contrast meets WCAG AA (run automated accessibility test)
- [x] T035 Verify prefers-reduced-motion disables animations
- [x] T036 Run Lighthouse performance audit (target: >90 performance score)

### Quality & Documentation

- [x] T037 Run all tests and verify 100% pass rate
- [x] T038 Run quickstart.md verification checklist

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-5)**: All depend on Foundational phase completion
  - User stories can proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3)
- **Integration (Phase 6)**: Depends on ALL user stories being complete
- **Polish (Phase 7)**: Depends on Integration being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Independent of US1
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) - Independent of US1/US2

### Within Each User Story

- Tests MUST be written and FAIL before implementation
- Component implementation before styles (or parallel)
- Index export after component is complete

### Parallel Opportunities

- T003, T004, T005, T006 can run in parallel (Setup phase)
- T008, T009 can run in parallel (Foundational phase)
- Once Foundational completes: US1, US2, US3 can all start in parallel
- T012, T016, T020 (all test files) can be written in parallel
- T024, T025 (integration tests) can run in parallel
- T031, T032, T033 (responsive styles) can run in parallel

---

## Parallel Example: All User Story Tests

```bash
# Write all component tests in parallel (TDD - Red phase):
Task T012: "Create ArcadeTitle component test in frontend/src/components/ArcadeTitle/ArcadeTitle.test.tsx"
Task T016: "Create CRTEffect component test in frontend/src/components/CRTEffect/CRTEffect.test.tsx"
Task T020: "Create StartButton component test in frontend/src/components/StartButton/StartButton.test.tsx"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1 (ArcadeTitle)
4. **STOP and VALIDATE**: Title visible with arcade styling
5. Can demo/deploy MVP with just the title

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 (ArcadeTitle) → Test independently → **MVP ready!**
3. Add User Story 2 (CRTEffect) → Enhanced visuals
4. Add User Story 3 (StartButton) → Navigation ready
5. Integration → Full main page
6. Polish → Production ready

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (ArcadeTitle)
   - Developer B: User Story 2 (CRTEffect)
   - Developer C: User Story 3 (StartButton)
3. Reconvene for Integration phase

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing (TDD Red-Green-Refactor)
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Constitution Principle II requires TDD - all tests included
