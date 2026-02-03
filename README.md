# Spec-Driven-Pacman

A quality-focused, security-oriented arcade game inspired by Pac-Man, built with React and Go.

## Project Philosophy

This project follows a **specification-driven development** approach where:

- Features start as specifications before any code is written
- Implementation follows explicit plans derived from specifications
- Quality and security are non-negotiable principles

See [Constitution](.specify/memory/constitution.md) for the full set of governing principles.

## Technology Stack

| Layer         | Technology                                                                |
| ------------- | ------------------------------------------------------------------------- |
| Frontend      | React 18+ (TypeScript), Vite                                              |
| Backend       | Go 1.21+                                                                  |
| Communication | REST API + WebSocket                                                      |
| Testing       | Vitest + React Testing Library (Frontend), Go testing + testify (Backend) |

## Project Structure

```
spec-driven-pacman/
├── .specify/              # Specification-driven development framework
│   ├── memory/            # Project constitution and memory
│   └── templates/         # Plan, spec, and task templates
├── backend/               # Go backend (API + game server)
│   ├── cmd/               # Entry points
│   ├── internal/          # Private application code
│   │   ├── domain/        # Game logic and entities
│   │   ├── handlers/      # HTTP and WebSocket handlers
│   │   └── services/      # Business logic services
│   └── tests/             # Backend tests
├── frontend/              # React client
│   ├── src/
│   │   ├── components/    # UI components
│   │   ├── game/          # Game engine and rendering
│   │   ├── hooks/         # Custom React hooks
│   │   └── services/      # API and WebSocket clients
│   └── tests/             # Frontend tests
├── specs/                 # Feature specifications
└── docs/                  # Project documentation
```

## Getting Started

> **Note**: This project is under development. Setup instructions will be added as the project progresses.

### Prerequisites

- Node.js 20+
- Go 1.21+
- Docker (optional, for containerized development)

### Development

```bash
# Frontend
cd frontend
npm install
npm run dev

# Backend
cd backend
go mod download
go run cmd/server/main.go
```

## Core Principles

1. **Security-First Development** - All code written with security as primary concern
2. **Test-Driven Development** - Tests before implementation (Red-Green-Refactor)
3. **Clean Architecture** - Clear separation of concerns
4. **API-First Design** - Contracts defined before implementation
5. **Performance & Game Loop Integrity** - 60 FPS client, consistent server tick rate

## Contributing

All contributions must adhere to the [Constitution](.specify/memory/constitution.md). Pull requests require:

- Compliance with all five core principles
- Test coverage for new code
- Security impact assessment
- Performance consideration for game-critical paths

## License

TBD
