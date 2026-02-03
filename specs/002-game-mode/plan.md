# Implementation Plan: Pac-Man Game Mode

**Branch**: `002-game-mode` | **Date**: 2026-02-03 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/002-game-mode/spec.md`

## Summary

Implement the core Pac-Man gameplay with a **server-authoritative architecture**. The Go backend handles all game logic (movement, collision, ghost AI, scoring) and serves the maze layout. The React frontend is a thin client that renders game state and captures keyboard input. Communication uses WebSocket for real-time game state updates and REST for game initialization.

## Technical Context

**Language/Version**:

- Frontend: TypeScript (strict mode) with React 18+
- Backend: Go 1.21+

**Primary Dependencies**:

- Frontend: React 18, Vite, Zustand (game state)
- Backend: Chi (HTTP), gorilla/websocket (realtime)

**Storage**: In-memory game state (no persistence required for MVP)

**Testing**:

- Frontend: Vitest + React Testing Library
- Backend: Go standard testing + testify

**Target Platform**: Modern desktop browsers (Chrome, Firefox, Safari, Edge)

**Project Type**: Web application (frontend + backend)

**Performance Goals**:

- Client: 60 FPS minimum
- Server: 20-60 ticks/second game loop
- Input latency: <150ms round-trip

**Constraints**:

- WebSocket connection required for gameplay
- Server-authoritative (no client-side game logic)

**Scale/Scope**: Single-player, single game session at a time

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

| Principle               | Requirement                                                          | Status |
| ----------------------- | -------------------------------------------------------------------- | ------ |
| I. Security-First       | Input validation at all boundaries, server-authoritative game state  | ✅     |
| II. Test-Driven         | Tests written before implementation, >80% coverage on critical paths | ✅     |
| III. Clean Architecture | Domain logic isolated, no circular dependencies                      | ✅     |
| IV. API-First           | Contracts defined before implementation (OpenAPI/JSON schemas)       | ✅     |
| V. Performance          | 60 FPS client target, consistent server tick rate                    | ✅     |

**Notes**:

- Security: Server-authoritative architecture ensures game state integrity; all inputs validated server-side
- TDD: Tests will be written before implementation per spec workflow
- Clean Architecture: Backend domain logic isolated from handlers; frontend components separated from state
- API-First: REST and WebSocket contracts defined in Phase 1
- Performance: 60 FPS client target documented; server tick rate configurable

## Project Structure

### Documentation (this feature)

```text
specs/002-game-mode/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
└── tasks.md             # Phase 2 output (created by /speckit.tasks)
```

### Source Code (repository root)

```text
backend/
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── internal/
│   ├── domain/
│   │   ├── game/             # Game state, rules, entities
│   │   ├── maze/             # Maze layout, pathfinding
│   │   └── ghost/            # Ghost AI behavior
│   ├── handlers/
│   │   ├── http/             # REST handlers
│   │   └── websocket/        # WebSocket handlers
│   └── services/
│       └── game/             # Game service orchestration
├── pkg/
│   └── types/                # Shared types/contracts
└── tests/
    ├── unit/
    └── integration/

frontend/
├── src/
│   ├── components/
│   │   ├── Game/             # Main game container
│   │   ├── Maze/             # Maze rendering
│   │   ├── PacMan/           # Pac-Man sprite
│   │   ├── Ghost/            # Ghost sprites
│   │   ├── HUD/              # Score, lives display
│   │   └── GameOver/         # End screens
│   ├── hooks/
│   │   ├── useGameState.ts   # Game state subscription
│   │   └── useKeyboard.ts    # Input capture
│   ├── services/
│   │   ├── api.ts            # REST client
│   │   └── websocket.ts      # WebSocket client
│   └── stores/
│       └── gameStore.ts      # Zustand game state
└── tests/
    ├── unit/
    └── integration/
```

**Structure Decision**: Web application with separate frontend (React/TypeScript) and backend (Go) projects. Backend follows clean architecture with domain/handlers/services separation. Frontend uses component-based architecture with centralized state management.

## Complexity Tracking

> No constitution violations identified. Design follows all principles.

| Violation | Why Needed | Simpler Alternative Rejected Because |
| --------- | ---------- | ------------------------------------ |
| None      | N/A        | N/A                                  |
