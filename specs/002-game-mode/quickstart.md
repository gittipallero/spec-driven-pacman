# Quickstart: Pac-Man Game Mode

**Feature**: 002-game-mode | **Date**: 2026-02-03

## Prerequisites

- **Node.js**: v18+ (for frontend)
- **Go**: v1.21+ (for backend)
- **Git**: For version control

## Project Setup

### 1. Clone and Navigate

```bash
cd spec-driven-pacman
```

### 2. Backend Setup

```bash
cd backend

# Initialize Go module (if not exists)
go mod init github.com/your-org/pacman-backend

# Install dependencies
go mod tidy

# Run tests
go test ./...

# Start server (development)
go run cmd/server/main.go
```

**Backend runs on**: `http://localhost:8080`

### 3. Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Install new dependencies for game mode
npm install zustand

# Run tests
npm test

# Start development server
npm run dev
```

**Frontend runs on**: `http://localhost:5173`

## Development Workflow

### TDD Cycle (Per Constitution)

1. **Write test** for new functionality
2. **Run test** - verify it fails (Red)
3. **Implement** minimal code to pass
4. **Run test** - verify it passes (Green)
5. **Refactor** if needed
6. **Repeat**

### Backend Development

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./internal/domain/game/...

# Run with coverage
go test -cover ./...

# Run with race detection
go test -race ./...

# Lint
golangci-lint run
```

### Frontend Development

```bash
# Run tests in watch mode
npm test -- --watch

# Run specific test file
npm test -- src/components/Game/Game.test.tsx

# Run with coverage
npm test -- --coverage

# Lint
npm run lint
```

## Key Files to Create

### Backend Structure

```text
backend/
├── cmd/server/main.go           # Entry point
├── internal/
│   ├── domain/
│   │   ├── game/
│   │   │   ├── game.go          # Game session logic
│   │   │   ├── game_test.go
│   │   │   ├── state.go         # Game state types
│   │   │   └── rules.go         # Game rules/scoring
│   │   ├── maze/
│   │   │   ├── maze.go          # Maze layout
│   │   │   ├── maze_test.go
│   │   │   └── classic.go       # Classic maze data
│   │   └── ghost/
│   │       ├── ghost.go         # Ghost entity
│   │       ├── ghost_test.go
│   │       └── ai.go            # Ghost AI behaviors
│   ├── handlers/
│   │   ├── http/
│   │   │   ├── game.go          # REST handlers
│   │   │   └── game_test.go
│   │   └── websocket/
│   │       ├── handler.go       # WebSocket handler
│   │       └── handler_test.go
│   └── services/
│       └── game/
│           ├── service.go       # Game orchestration
│           └── service_test.go
└── go.mod
```

### Frontend Structure

```text
frontend/src/
├── components/
│   ├── Game/
│   │   ├── Game.tsx             # Main game container
│   │   ├── Game.test.tsx
│   │   └── Game.module.css
│   ├── Maze/
│   │   ├── Maze.tsx             # Canvas maze renderer
│   │   ├── Maze.test.tsx
│   │   └── Maze.module.css
│   ├── PacMan/
│   │   ├── PacMan.tsx           # Pac-Man sprite
│   │   └── PacMan.test.tsx
│   ├── Ghost/
│   │   ├── Ghost.tsx            # Ghost sprite
│   │   └── Ghost.test.tsx
│   ├── HUD/
│   │   ├── HUD.tsx              # Score/lives display
│   │   ├── HUD.test.tsx
│   │   └── HUD.module.css
│   └── GameOver/
│       ├── GameOver.tsx         # End screen
│       ├── GameOver.test.tsx
│       └── GameOver.module.css
├── hooks/
│   ├── useGameState.ts          # Game state subscription
│   ├── useGameState.test.ts
│   ├── useKeyboard.ts           # Keyboard input
│   └── useKeyboard.test.ts
├── services/
│   ├── api.ts                   # REST client
│   ├── api.test.ts
│   ├── websocket.ts             # WebSocket client
│   └── websocket.test.ts
└── stores/
    ├── gameStore.ts             # Zustand store
    └── gameStore.test.ts
```

## API Contracts

### Start Game (REST)

```bash
curl -X POST http://localhost:8080/api/game/start
```

Response:
```json
{
  "sessionId": "uuid",
  "websocketUrl": "/ws/game/uuid",
  "state": { ... },
  "maze": { ... }
}
```

### WebSocket Connection

```javascript
const ws = new WebSocket('ws://localhost:8080/ws/game/{sessionId}');

// Send input
ws.send(JSON.stringify({ type: 'input', direction: 'up' }));

// Receive state updates
ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  if (message.type === 'state') {
    // Update game state
  }
};
```

## Testing Strategy

### Unit Tests (Critical Path)

| Component | Coverage Target | Focus Areas |
| --------- | --------------- | ----------- |
| Game Logic | >90% | Movement, collision, scoring |
| Ghost AI | >90% | Targeting, state transitions |
| Maze | >80% | Tile queries, pathfinding |
| WebSocket | >80% | Message handling, errors |

### Integration Tests

| Test | Description |
| ---- | ----------- |
| Game Start | REST → WebSocket → Initial state |
| Input Flow | Key press → WS message → State update |
| Ghost Collision | Movement → Collision → Life lost |
| Level Complete | All dots → Level transition |

## Environment Variables

### Backend

```bash
# .env (create in backend/)
PORT=8080
LOG_LEVEL=debug
GAME_TICK_RATE=20
```

### Frontend

```bash
# .env (create in frontend/)
VITE_API_URL=http://localhost:8080
VITE_WS_URL=ws://localhost:8080
```

## Common Commands

| Task | Command |
| ---- | ------- |
| Start backend | `cd backend && go run cmd/server/main.go` |
| Start frontend | `cd frontend && npm run dev` |
| Test backend | `cd backend && go test ./...` |
| Test frontend | `cd frontend && npm test` |
| Lint backend | `cd backend && golangci-lint run` |
| Lint frontend | `cd frontend && npm run lint` |
| Build frontend | `cd frontend && npm run build` |

## Next Steps

1. Create task breakdown with `/speckit.tasks`
2. Implement backend game loop (P1)
3. Implement frontend rendering (P1)
4. Connect via WebSocket (P1)
5. Add ghost AI (P1)
6. Add scoring/lives (P2)

