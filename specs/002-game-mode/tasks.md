# Tasks: Pac-Man Game Mode

**Input**: Design documents from `/specs/002-game-mode/`
**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, data-model.md ✅, contracts/ ✅

**Tests**: Per Constitution Principle II (Test-Driven Development), tests are written FIRST for all critical game logic. Tests MUST fail before implementation.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

### Backend Setup

- [x] T001 [P] Initialize Go module in `backend/` with `go mod init`
- [x] T002 [P] Add dependencies: chi, gorilla/websocket, testify, uuid in `backend/go.mod`
- [x] T003 [P] Create directory structure per plan.md:
  ```
  backend/cmd/server/
  backend/internal/domain/{game,maze,ghost}/
  backend/internal/handlers/{http,websocket}/
  backend/internal/services/game/
  backend/pkg/types/
  ```
- [x] T004 [P] Configure golangci-lint with `.golangci.yml`

### Frontend Setup

- [x] T005 [P] Install Zustand dependency: `npm install zustand`
- [x] T006 [P] Create directory structure per plan.md:
  ```
  frontend/src/components/{Game,Maze,PacMan,Ghost,HUD,GameOver}/
  frontend/src/hooks/
  frontend/src/services/
  frontend/src/stores/
  ```
- [x] T007 [P] Add environment variables to `frontend/.env` for API/WS URLs (blocked by gitignore, using vite.config.ts defaults)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

### Backend Domain Types

- [x] T008 [P] Create shared types in `backend/pkg/types/types.go`:
  - Direction enum (up, down, left, right, none)
  - GameStatus enum (waiting, playing, paused, gameOver, levelComplete)
  - Position struct (x, y)
- [x] T009 [P] Create Pac-Man entity in `backend/internal/domain/game/pacman.go`:
  - PacMan struct with position, direction, nextDirection, state
  - PacManState enum (idle, moving, dying)
- [x] T010 [P] Create Ghost entity in `backend/internal/domain/ghost/ghost.go`:
  - Ghost struct with name, position, direction, state, mode
  - GhostName enum (blinky, pinky, inky, clyde)
  - GhostState enum (normal, vulnerable, eaten, respawning)
  - GhostMode enum (chase, scatter, frightened)
- [x] T011 [P] Create Maze entity in `backend/internal/domain/maze/maze.go`:
  - Maze struct with width, height, tiles
  - TileType constants (empty, wall, dot, power, ghostHouse, tunnel)
- [x] T012 Create GameSession entity in `backend/internal/domain/game/session.go`:
  - GameSession struct with id, status, score, lives, level, vulnerabilityTimer
  - Contains PacMan, Ghosts, Maze references

### Backend HTTP Server

- [x] T013 Create main entry point in `backend/cmd/server/main.go`:
  - Initialize Chi router
  - Configure CORS middleware
  - Register routes
  - Start HTTP server on configurable port
- [x] T014 Create REST handler structure in `backend/internal/handlers/http/game.go`:
  - Mount routes: POST /api/game/start, GET /api/game/{sessionId}, POST /api/game/{sessionId}/pause, POST /api/game/{sessionId}/resume
- [x] T015 Create WebSocket handler structure in `backend/internal/handlers/websocket/handler.go`:
  - WebSocket upgrade endpoint: /ws/game/{sessionId}
  - Connection management skeleton
  - Handle pause/resume messages from client

### Frontend Core Infrastructure

- [x] T016 [P] Create Zustand game store in `frontend/src/stores/gameStore.ts`:
  - State: sessionId, gameStatus, maze, pacman, ghosts, score, lives, dotsRemaining
  - Actions: setGameState, updateFromServer, reset, togglePause
- [x] T017 [P] Create REST API client in `frontend/src/services/api.ts`:
  - startGame(): POST /api/game/start
  - getGameState(sessionId): GET /api/game/{sessionId}
  - pauseGame(sessionId), resumeGame(sessionId), endGame(sessionId)
- [x] T018 Create WebSocket client in `frontend/src/services/websocket.ts`:
  - connect(sessionId): establish WebSocket connection
  - sendInput(direction): send player input
  - sendPause(): send pause command
  - sendResume(): send resume command
  - onStateUpdate(callback): handle state messages
  - disconnect(): clean up connection
  - Auto-reconnect with exponential backoff

**Checkpoint**: Foundation ready - user story implementation can now begin

---

## Phase 3: User Story 1 - Start Game and Move Pac-Man (Priority: P1) 🎯 MVP

**Goal**: Player clicks START GAME, maze appears, Pac-Man moves with arrow keys, Escape pauses game

**Independent Test**: Click Start → game screen with maze → arrow keys move Pac-Man → Escape toggles pause

### Tests for User Story 1 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T019 [P] [US1] Unit test for Maze loading in `backend/internal/domain/maze/maze_test.go`:
  - Test classic maze has correct dimensions (28x31)
  - Test tile types are correctly placed
- [x] T020 [P] [US1] Unit test for Pac-Man movement in `backend/internal/domain/game/pacman_test.go`:
  - Test move in valid direction
  - Test blocked by wall
  - Test direction queuing
- [x] T021 [P] [US1] Unit test for GameSession initialization in `backend/internal/domain/game/session_test.go`:
  - Test NewGameSession creates valid initial state
  - Test Pac-Man starts at correct position
  - Test pause/resume toggles game status
- [x] T022 [P] [US1] Integration test for /api/game/start in `backend/tests/integration/game_start_test.go`:
  - Test returns 201 with sessionId, websocketUrl, state, maze
- [x] T023 [P] [US1] Component test for Maze rendering in `frontend/src/components/Maze/Maze.test.tsx`:
  - Test renders correct number of tiles
  - Test wall vs corridor styling
- [x] T024 [P] [US1] Hook test for useKeyboard in `frontend/src/hooks/useKeyboard.test.ts`:
  - Test arrow key capture
  - Test Escape key triggers pause callback
  - Test callback invoked with correct direction

### Implementation for User Story 1

- [x] T025 [US1] Implement classic maze data in `backend/internal/domain/maze/classic.go`:
  - Define 28x31 tile array with standard Pac-Man layout
  - Include all walls, dots, power pellets, ghost house, tunnels
- [x] T026 [US1] Implement maze loading in `backend/internal/domain/maze/maze.go`:
  - LoadClassicMaze() returns maze with all tiles
  - GetTile(x, y) returns tile type
  - IsWalkable(x, y) checks if position is valid
- [x] T027 [US1] Implement Pac-Man movement logic in `backend/internal/domain/game/pacman.go`:
  - Move(direction, maze) updates position if valid
  - SetNextDirection(direction) queues direction change
  - CanMove(direction, maze) validates against walls
- [x] T028 [US1] Implement GameSession creation and pause in `backend/internal/domain/game/session.go`:
  - NewGameSession() initializes with maze, Pac-Man at (13,23), 3 lives, score 0
  - ProcessInput(direction) applies input to Pac-Man
  - Pause() sets status to paused, stops game tick
  - Resume() sets status to playing, resumes game tick
- [x] T029 [US1] Implement game service in `backend/internal/services/game/service.go`:
  - CreateGame() returns new session with ID
  - GetSession(id) retrieves session
  - HandleInput(sessionId, direction) processes player input
  - PauseGame(sessionId) pauses the game
  - ResumeGame(sessionId) resumes the game
- [x] T030 [US1] Implement POST /api/game/start handler in `backend/internal/handlers/http/game.go`:
  - Create game via service
  - Return GameStartResponse with sessionId, wsUrl, state, maze
- [x] T031 [US1] Implement pause/resume handlers in `backend/internal/handlers/http/game.go`:
  - POST /api/game/{sessionId}/pause → call service.PauseGame()
  - POST /api/game/{sessionId}/resume → call service.ResumeGame()
- [x] T032 [US1] Implement WebSocket message handling in `backend/internal/handlers/websocket/handler.go`:
  - Parse "input" messages for direction
  - Parse "pause" messages to pause game
  - Parse "resume" messages to resume game
  - Call appropriate service methods
  - Broadcast state update
- [x] T033 [US1] Implement useKeyboard hook in `frontend/src/hooks/useKeyboard.ts`:
  - Listen for keydown events (ArrowUp, ArrowDown, ArrowLeft, ArrowRight)
  - Listen for Escape key to toggle pause
  - Call provided callbacks (onDirection, onPause)
  - Handle cleanup on unmount
- [x] T034 [US1] Implement Maze component in `frontend/src/components/Maze/Maze.tsx`:
  - Canvas-based rendering of maze tiles
  - Wall styling (blue), corridor (black), dots (white)
  - Props: maze data from store
- [x] T035 [US1] Implement PacMan component in `frontend/src/components/PacMan/PacMan.tsx`:
  - Render Pac-Man sprite at position
  - Animate based on direction
  - Props: x, y, direction, state
- [x] T036 [US1] Implement PauseOverlay component in `frontend/src/components/Game/PauseOverlay.tsx`:
  - Display "PAUSED" text when game status is paused
  - Show "Press ESC to resume" hint
  - Semi-transparent overlay styling
- [x] T037 [US1] Implement Game container in `frontend/src/components/Game/Game.tsx`:
  - Call api.startGame() on mount
  - Connect to WebSocket
  - Render Maze + PacMan + PauseOverlay
  - Wire useKeyboard to send inputs and handle Escape for pause
  - Toggle pause state on Escape key press
- [x] T038 [US1] Integrate StartButton click to navigate to Game in `frontend/src/App.tsx`:
  - Add routing or conditional render
  - StartButton onClick → show Game component

**Checkpoint**: User Story 1 complete - can start game, move Pac-Man, and pause with Escape

---

## Phase 4: User Story 2 - Collect Dots and Power Pellets (Priority: P1)

**Goal**: Pac-Man eats dots for points, power pellets trigger ghost vulnerability

**Independent Test**: Move Pac-Man over dot → dot disappears, score increases

### Tests for User Story 2 ⚠️

- [x] T039 [P] [US2] Unit test for dot collection in `backend/internal/domain/game/session_test.go`:
  - Test dot removed from maze when Pac-Man overlaps
  - Test score increases by 10 for dot
  - Test score increases by 50 for power pellet
- [x] T040 [P] [US2] Unit test for vulnerability timer in `backend/internal/domain/game/session_test.go`:
  - Test power pellet sets vulnerability timer
  - Test timer decrements each tick
  - Test ghosts return to normal when timer expires

### Implementation for User Story 2

- [x] T041 [US2] Implement dot tracking in `backend/internal/domain/maze/maze.go`:
  - dotsRemaining counter
  - CollectDot(x, y) removes dot, returns points
  - IsDot(x, y), IsPowerPellet(x, y) checks
- [x] T042 [US2] Implement scoring in `backend/internal/domain/game/session.go`:
  - CheckDotCollection() on each tick
  - AddScore(points) updates session score
  - Handle power pellet → set vulnerabilityTimer to 6000ms
- [x] T043 [US2] Implement vulnerability timer tick in `backend/internal/domain/game/session.go`:
  - Tick() decrements timer by tick interval
  - When timer reaches 0, trigger ghost state change
- [x] T044 [US2] Update Maze component for dot rendering in `frontend/src/components/Maze/Maze.tsx`:
  - Render dots as small white circles
  - Render power pellets as larger flashing circles
  - Update when dotsRemaining changes
- [x] T045 [US2] Update gameStore to track collected dots in `frontend/src/stores/gameStore.ts`:
  - Handle collectedDots array from server state
  - Animate dot collection

**Checkpoint**: User Story 2 complete - dots and power pellets work

---

## Phase 5: User Story 3 - Ghost Encounters (Priority: P1)

**Goal**: Four ghosts chase Pac-Man with AI, collision loses life or earns points

**Independent Test**: Ghost touches Pac-Man → lose life; vulnerable ghost touched → points

### Tests for User Story 3 ⚠️

- [x] T046 [P] [US3] Unit test for ghost AI in `backend/internal/domain/ghost/ai_test.go`:
  - Test Blinky targets Pac-Man position
  - Test Pinky targets 4 tiles ahead
  - Test scatter mode moves to corner
- [x] T047 [P] [US3] Unit test for ghost collision in `backend/internal/domain/game/session_test.go`:
  - Test normal ghost collision → lose life
  - Test vulnerable ghost collision → ghost eaten, score +200/400/800/1600
  - Test eaten ghost returns to ghost house
- [x] T048 [P] [US3] Unit test for ghost state transitions in `backend/internal/domain/ghost/ghost_test.go`:
  - Test normal → vulnerable on power pellet
  - Test vulnerable → normal on timer expire
  - Test eaten → respawning → normal cycle

### Implementation for User Story 3

- [x] T049 [US3] Implement ghost AI targeting in `backend/internal/domain/ghost/ai.go`:
  - BlinkyTarget(pacmanPos) returns Pac-Man position
  - PinkyTarget(pacmanPos, pacmanDir) returns 4 tiles ahead
  - InkyTarget(pacmanPos, blinkyPos) returns complex target
  - ClydeTarget(pacmanPos, clydePos) returns conditional target
- [x] T050 [US3] Implement ghost movement in `backend/internal/domain/ghost/ghost.go`:
  - Move(target, maze) moves toward target avoiding walls
  - ChooseDirection(target, maze) picks best valid direction
  - ReverseDirection() for mode changes
- [x] T051 [US3] Implement ghost state machine in `backend/internal/domain/ghost/ghost.go`:
  - SetVulnerable() changes state, enables eating
  - SetEaten() triggers return to ghost house
  - SetRespawning() exits ghost house
  - SetNormal() resumes chase
- [x] T052 [US3] Implement collision detection in `backend/internal/domain/game/session.go`:
  - CheckGhostCollisions() each tick
  - HandleNormalCollision() → lose life, reset positions
  - HandleVulnerableCollision(ghost) → eat ghost, add combo score
- [x] T053 [US3] Implement game loop tick in `backend/internal/services/game/service.go`:
  - RunGameLoop() at 20 TPS
  - Each tick: move ghosts, check collisions, update timers
  - Skip tick processing when game is paused
  - Broadcast state to connected clients
- [x] T054 [US3] Implement Ghost component in `frontend/src/components/Ghost/Ghost.tsx`:
  - Render ghost sprite at position
  - Color based on name (blinky=red, pinky=pink, inky=cyan, clyde=orange)
  - Blue color when vulnerable
  - Eyes only when eaten
- [x] T055 [US3] Update Game component to render ghosts in `frontend/src/components/Game/Game.tsx`:
  - Map over ghosts array from store
  - Render Ghost component for each

**Checkpoint**: User Story 3 complete - full gameplay loop works

---

## Phase 6: User Story 4 - Score and Lives Display (Priority: P2)

**Goal**: Show current score and remaining lives during gameplay

**Independent Test**: Score updates on dot collection; lives decrease on death

### Tests for User Story 4 ⚠️

- [x] T056 [P] [US4] Component test for HUD in `frontend/src/components/HUD/HUD.test.tsx`:
  - Test score displays correctly
  - Test lives display correct count
  - Test updates on prop change

### Implementation for User Story 4

- [x] T057 [US4] Implement HUD component in `frontend/src/components/HUD/HUD.tsx`:
  - Display score with retro font styling
  - Display lives as Pac-Man icons
  - Display level number
  - Props: score, lives, level from store
- [x] T058 [US4] Style HUD in `frontend/src/components/HUD/HUD.module.css`:
  - Position at top of game screen
  - Retro arcade styling matching MainPage theme
- [x] T059 [US4] Integrate HUD into Game component in `frontend/src/components/Game/Game.tsx`:
  - Render HUD above maze
  - Connect to gameStore for values

**Checkpoint**: User Story 4 complete - score and lives visible

---

## Phase 7: User Story 5 - Level Completion and Game Over (Priority: P2)

**Goal**: Level complete when all dots eaten; game over when no lives remain

**Independent Test**: Collect all dots → level complete; lose all lives → game over

### Tests for User Story 5 ⚠️

- [x] T060 [P] [US5] Unit test for level completion in `backend/internal/domain/game/session_test.go`:
  - Test level complete when dotsRemaining = 0
  - Test level increments
  - Test positions reset with new dots
- [x] T061 [P] [US5] Unit test for game over in `backend/internal/domain/game/session_test.go`:
  - Test game over when lives = 0
  - Test final score preserved
- [x] T062 [P] [US5] Component test for GameOver in `frontend/src/components/GameOver/GameOver.test.tsx`:
  - Test shows final score
  - Test shows restart button
  - Test shows main menu button

### Implementation for User Story 5

- [x] T063 [US5] Implement level completion in `backend/internal/domain/game/session.go`:
  - CheckLevelComplete() when dotsRemaining = 0
  - IncrementLevel() resets maze, positions, increments level
  - Status → levelComplete briefly, then playing
- [x] T064 [US5] Implement game over in `backend/internal/domain/game/session.go`:
  - CheckGameOver() when lives = 0
  - Status → gameOver
  - Stop game loop for session
- [x] T065 [US5] Implement GameOver component in `frontend/src/components/GameOver/GameOver.tsx`:
  - Display "GAME OVER" or "LEVEL COMPLETE" based on status
  - Show final score
  - Restart button → call api.startGame(), reset store
  - Main Menu button → navigate back to MainPage
- [x] T066 [US5] Style GameOver in `frontend/src/components/GameOver/GameOver.module.css`:
  - Overlay styling
  - Match arcade aesthetic from MainPage
- [x] T067 [US5] Integrate GameOver into Game component in `frontend/src/components/Game/Game.tsx`:
  - Show GameOver when gameStatus is gameOver or levelComplete
  - Handle button callbacks

**Checkpoint**: User Story 5 complete - full game lifecycle works

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

### Security Hardening (Constitution Principle I)

- [x] T068 Input validation audit in `backend/internal/handlers/`:
  - Validate direction is valid enum value
  - Validate sessionId format
  - Rate limit input messages (max 60/sec)
- [x] T069 Add CORS configuration in `backend/cmd/server/main.go`:
  - Configure allowed origins
  - Configure allowed methods and headers
- [x] T070 Add request logging middleware in `backend/internal/handlers/http/middleware.go`:
  - Log all requests with timing
  - Log WebSocket connections

### Network Resilience

- [x] T071 Implement reconnection UI in `frontend/src/components/Game/Game.tsx`:
  - Show "Reconnecting..." overlay on disconnect
  - Hide overlay on successful reconnect
  - Show "Connection Lost" after 5 failures
- [x] T072 Implement input queuing in `frontend/src/services/websocket.ts`:
  - Queue inputs during disconnect
  - Send queued inputs on reconnect

### Performance (Constitution Principle V)

- [x] T073 Implement client-side interpolation in `frontend/src/hooks/useGameState.ts`:
  - Interpolate positions between server ticks
  - Smooth 60 FPS rendering from 20 TPS server
- [x] T074 Optimize canvas rendering in `frontend/src/components/Maze/Maze.tsx`:
  - Use requestAnimationFrame
  - Only redraw changed tiles
  - Pre-render static maze layer

### Quality & Documentation

- [x] T075 [P] Add error boundaries in `frontend/src/components/Game/Game.tsx`:
  - Catch rendering errors
  - Show user-friendly error message
- [x] T076 [P] Run quickstart.md validation:
  - Verify all setup commands work
  - Update any outdated instructions

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-7)**: All depend on Foundational phase completion
  - US1 (Phase 3): Must complete first - establishes core game loop and pause
  - US2 (Phase 4): Can start after US1 tests pass
  - US3 (Phase 5): Can start after US1 tests pass
  - US4 (Phase 6): Can start after US1 tests pass
  - US5 (Phase 7): Depends on US2 (dots) and US3 (lives) completion
- **Polish (Phase 8)**: Depends on all user stories being complete

### Within Each User Story

1. Tests MUST be written and FAIL before implementation
2. Backend domain logic before handlers
3. Frontend components before integration
4. Integration and wiring last

### Parallel Opportunities

- All Setup tasks (T001-T007) can run in parallel
- All Foundational type definitions (T008-T012) can run in parallel
- All tests for a user story can run in parallel
- US2, US3, US4 can start in parallel after US1 core is complete

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Game starts, Pac-Man moves, Escape pauses
5. Demo if ready

### Full P1 Delivery

1. Complete Setup + Foundational
2. Complete US1 → Validate (including pause functionality)
3. Complete US2 + US3 in parallel → Validate full gameplay
4. Demo playable game

### Complete Feature

1. All above
2. Add US4 (HUD) + US5 (Game Over)
3. Complete Phase 8 (Polish)
4. Final validation against spec.md acceptance scenarios

---

## Task Summary

| Phase                     | Tasks     | Focus                                   |
| ------------------------- | --------- | --------------------------------------- |
| 1. Setup                  | T001-T007 | Project structure                       |
| 2. Foundational           | T008-T018 | Core infrastructure                     |
| 3. US1 - Movement & Pause | T019-T038 | Start game, move Pac-Man, Escape pauses |
| 4. US2 - Dots             | T039-T045 | Collect dots, scoring                   |
| 5. US3 - Ghosts           | T046-T055 | Ghost AI, collisions                    |
| 6. US4 - HUD              | T056-T059 | Score/lives display                     |
| 7. US5 - End States       | T060-T067 | Level complete, game over               |
| 8. Polish                 | T068-T076 | Security, performance, quality          |

**Total**: 76 tasks across 8 phases

---

## Notes

- [P] tasks = different files, no dependencies
- [US#] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing
- Commit after each task or logical group
- Constitution requires TDD - never skip writing tests first
- **Escape key toggles pause** - implemented in US1 (no HUD pause button per clarification)
