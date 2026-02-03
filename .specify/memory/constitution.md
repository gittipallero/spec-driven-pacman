<!--
=== SYNC IMPACT REPORT ===
Version change: N/A → 1.0.0 (Initial creation)

Added Principles:
- I. Security-First Development (NEW)
- II. Test-Driven Development (NEW)
- III. Clean Architecture (NEW)
- IV. API-First Design (NEW)
- V. Performance & Game Loop Integrity (NEW)

Added Sections:
- Technology Stack
- Security Requirements
- Development Workflow
- Governance

Templates reviewed:
- ✅ plan-template.md - Compatible with constitution principles
- ✅ spec-template.md - Compatible with constitution requirements
- ✅ tasks-template.md - Compatible with task categorization

Follow-up TODOs: None
==================
-->

# Spec-Driven-Pacman Constitution

## Core Principles

### I. Security-First Development

All code MUST be written with security as a primary concern. Security is non-negotiable
and takes precedence over feature velocity.

- Input validation MUST occur at all system boundaries (API endpoints, WebSocket messages, user inputs)
- All user-supplied data MUST be sanitized before processing or storage
- Authentication and authorization MUST be implemented for all protected endpoints
- Secrets and credentials MUST NOT be hardcoded; use environment variables or secure vaults
- Dependencies MUST be regularly audited for known vulnerabilities
- CORS, CSP, and other security headers MUST be properly configured
- Game state manipulation attempts MUST be detected and logged server-side

**Rationale**: Arcade games are targets for cheating and exploitation. Server-authoritative
game logic and secure communication protect game integrity and user trust.

### II. Test-Driven Development

Tests MUST be written before implementation code. The Red-Green-Refactor cycle is mandatory
for all feature development.

- Unit tests MUST cover all business logic and game mechanics
- Integration tests MUST verify API contracts between React client and Go backend
- Contract tests MUST validate WebSocket message formats and game state synchronization
- Test coverage MUST be maintained above 80% for critical game logic paths
- Tests MUST be deterministic and not rely on timing or external state
- Game simulation tests MUST verify core mechanics (movement, collision, scoring)

**Rationale**: Quality code requires verification. TDD ensures correctness of game mechanics
and prevents regressions in complex state management.

### III. Clean Architecture

Code MUST follow clean architecture principles with clear separation of concerns.

- Backend: Domain logic MUST be isolated from HTTP/WebSocket handlers
- Frontend: UI components MUST be separated from game logic and state management
- Shared types/contracts MUST be defined in dedicated locations
- Dependencies MUST point inward (handlers → services → domain)
- Game engine logic MUST be framework-agnostic where possible
- No circular dependencies between packages/modules

**Rationale**: Clean architecture enables independent testing, easier maintenance,
and the ability to swap infrastructure without affecting game logic.

### IV. API-First Design

All communication between frontend and backend MUST be defined through explicit contracts
before implementation.

- REST API endpoints MUST be documented with OpenAPI/Swagger specifications
- WebSocket message formats MUST be defined with JSON schemas
- Breaking API changes MUST increment major version and require migration plan
- Client and server MUST validate messages against defined schemas
- Error responses MUST follow consistent format with actionable messages

**Rationale**: React and Go services evolve independently. Explicit contracts prevent
integration failures and enable parallel development.

### V. Performance & Game Loop Integrity

Game performance MUST meet arcade-quality standards. The game loop is sacred.

- Client MUST maintain 60 FPS minimum under normal conditions
- Server game tick rate MUST be consistent (recommended: 20-60 ticks/second)
- Network latency compensation MUST be implemented for smooth gameplay
- Memory leaks MUST be prevented; resources MUST be properly cleaned up
- Asset loading MUST NOT block the game loop
- State updates MUST be batched efficiently to minimize re-renders

**Rationale**: Arcade games require responsive, smooth gameplay. Performance issues
directly impact user experience and game feel.

## Technology Stack

The following technology choices are binding for this project:

**Frontend (Client)**:

- Language: TypeScript (strict mode enabled)
- Framework: React 18+ with functional components and hooks
- State Management: React Context or Zustand for game state
- Build Tool: Vite
- Testing: Vitest + React Testing Library

**Backend (Server)**:

- Language: Go 1.21+
- HTTP Framework: Chi or Gin (lightweight, performant)
- WebSocket: gorilla/websocket or nhooyr/websocket
- Testing: Go standard testing + testify
- Linting: golangci-lint with strict configuration

**Communication**:

- REST API for non-realtime operations (scores, user management)
- WebSocket for realtime game state synchronization
- JSON for all data serialization

**Infrastructure**:

- Containerization: Docker with multi-stage builds
- CI/CD: GitHub Actions or equivalent

## Security Requirements

The following security controls are MANDATORY:

1. **Authentication**: JWT-based authentication for user sessions
2. **Authorization**: Role-based access control for admin functions
3. **Rate Limiting**: All endpoints MUST be rate-limited to prevent abuse
4. **Input Validation**: All inputs validated server-side (never trust the client)
5. **Game State Authority**: Server MUST be the single source of truth for game state
6. **Anti-Cheat**: Client inputs MUST be validated against possible game states
7. **Secure Transport**: HTTPS/WSS required in production
8. **Dependency Scanning**: Automated vulnerability scanning in CI pipeline
9. **Logging**: Security events MUST be logged (auth failures, suspicious patterns)
10. **No Sensitive Data Exposure**: Stack traces and internal errors hidden from clients

## Development Workflow

All development MUST follow this workflow:

1. **Specification First**: Features start with a specification document
2. **Plan Creation**: Implementation plan created from specification
3. **Contract Definition**: API contracts defined before implementation
4. **Test Writing**: Tests written based on specification and contracts
5. **Implementation**: Code written to make tests pass
6. **Review**: Code review required for all changes
7. **Integration**: Continuous integration validates all tests pass
8. **Documentation**: User-facing changes documented

Code reviews MUST verify:

- Adherence to constitution principles
- Test coverage for new code
- Security considerations addressed
- Performance impact assessed for game-critical paths

## Governance

This constitution is the authoritative source for project standards. All code, reviews,
and architectural decisions MUST comply with these principles.

**Amendment Process**:

1. Propose amendment with rationale in writing
2. Review impact on existing code and documentation
3. Update constitution with new version number
4. Update dependent templates if affected
5. Document migration path for existing code if needed

**Versioning Policy**:

- MAJOR: Principle removal or incompatible redefinition
- MINOR: New principle or section added
- PATCH: Clarifications and non-semantic changes

**Compliance**:

- All pull requests MUST include constitution compliance verification
- Violations MUST be justified in writing with complexity tracking
- Repeated violations trigger architecture review

**Guidance**: Use `.specify/templates/` for runtime development templates and workflows.

**Version**: 1.0.0 | **Ratified**: 2026-02-03 | **Last Amended**: 2026-02-03
