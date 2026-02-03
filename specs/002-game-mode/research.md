# Research: Pac-Man Game Mode

**Feature**: 002-game-mode | **Date**: 2026-02-03

## Overview

This document consolidates research decisions for implementing the Pac-Man game mode with server-authoritative architecture.

---

## 1. Real-time Communication Protocol

### Decision: WebSocket for game state updates

**Rationale**: 
- WebSocket provides bidirectional, low-latency communication essential for real-time games
- Maintains persistent connection, avoiding HTTP overhead for frequent updates
- Supports server-push model for game tick updates
- gorilla/websocket is mature, performant, and widely used in Go ecosystem

**Alternatives Considered**:

| Option | Pros | Cons | Verdict |
| ------ | ---- | ---- | ------- |
| WebSocket | Low latency, bidirectional, persistent | Requires connection management | ✅ Selected |
| Server-Sent Events (SSE) | Simple, HTTP-based | Unidirectional only, can't send inputs efficiently | ❌ Rejected |
| HTTP Polling | Simple to implement | High latency, wasteful requests | ❌ Rejected |
| WebRTC | Ultra-low latency | Overkill for single-player, complex setup | ❌ Rejected |

---

## 2. Game Loop Architecture

### Decision: Server-side game loop with fixed tick rate (20 ticks/second)

**Rationale**:
- 20 TPS provides ~50ms between updates, sufficient for Pac-Man's grid-based movement
- Server-authoritative prevents cheating and ensures consistent game state
- Client interpolates between server states for smooth 60 FPS rendering
- Lower tick rate reduces server load while maintaining playability

**Implementation Pattern**:
```text
Server: Run game tick every 50ms
  → Calculate positions, collisions, ghost AI
  → Broadcast state to connected clients

Client: Render at 60 FPS
  → Interpolate between last two server states
  → Send inputs immediately (optimistic)
  → Apply server corrections on next state
```

**Alternatives Considered**:

| Option | Pros | Cons | Verdict |
| ------ | ---- | ---- | ------- |
| 20 TPS server tick | Balanced load/responsiveness | Requires interpolation | ✅ Selected |
| 60 TPS server tick | Matches client FPS | High server load, unnecessary | ❌ Rejected |
| Client-authoritative | Simpler networking | Enables cheating, inconsistent state | ❌ Rejected |

---

## 3. Ghost AI Approach

### Decision: Classic Pac-Man AI patterns (Blinky, Pinky, Inky, Clyde)

**Rationale**:
- Well-documented algorithms create authentic Pac-Man experience
- Each ghost has distinct personality (target selection strategy)
- Scatter/Chase mode alternation adds strategic depth
- Relatively simple to implement with grid-based pathfinding

**Ghost Behavior Summary**:
- **Blinky (Red)**: Directly targets Pac-Man's current position
- **Pinky (Pink)**: Targets 4 tiles ahead of Pac-Man
- **Inky (Cyan)**: Complex targeting based on Blinky's position
- **Clyde (Orange)**: Targets Pac-Man when far, scatters when close

**Pathfinding**: Simple grid-based direction selection (no A* needed for authentic behavior)

---

## 4. Maze Representation

### Decision: 2D grid array with tile types

**Rationale**:
- Simple integer/enum grid efficiently represents maze layout
- Easy collision detection via tile lookup
- Straightforward serialization for API transport
- Matches classic Pac-Man's tile-based design (28x31 grid)

**Tile Types**:
```text
0 = Empty (walkable)
1 = Wall
2 = Dot
3 = Power Pellet
4 = Ghost House
5 = Tunnel (wraps to opposite side)
```

**Storage**: Hardcoded in backend (single classic maze for MVP)

---

## 5. State Synchronization Strategy

### Decision: Full state broadcast with delta compression consideration

**Rationale**:
- Game state is relatively small (~2-3KB per tick)
- Full state simplifies client logic (no state reconstruction needed)
- Server remains single source of truth
- Delta compression can be added later if bandwidth becomes concern

**State Payload Structure**:
```text
{
  pacman: { x, y, direction, state }
  ghosts: [{ x, y, direction, state, name }]
  score: number
  lives: number
  level: number
  dots: [{ x, y, collected }]  // or bitmask for efficiency
  powerPellets: [{ x, y, collected }]
  gameStatus: "playing" | "paused" | "gameOver" | "levelComplete"
  vulnerabilityTimer: number (remaining ms)
}
```

---

## 6. Input Handling

### Decision: Immediate input transmission with server validation

**Rationale**:
- Client sends direction input immediately for responsiveness
- Server validates against game rules (wall collision, valid direction)
- Input queuing on client handles network latency
- Server applies input at next tick if valid

**Input Message**:
```text
{ type: "input", direction: "up" | "down" | "left" | "right" }
```

---

## 7. Frontend State Management

### Decision: Zustand for game state

**Rationale**:
- Lightweight (~1KB) compared to Redux
- Simple API with hooks integration
- Supports subscriptions for selective re-renders
- No boilerplate, perfect for game state shape

**Alternatives Considered**:

| Option | Pros | Cons | Verdict |
| ------ | ---- | ---- | ------- |
| Zustand | Simple, lightweight, hooks-native | Less ecosystem than Redux | ✅ Selected |
| Redux Toolkit | Mature, dev tools | Overkill for single game state | ❌ Rejected |
| React Context | Built-in, no deps | Re-render issues at high frequency | ❌ Rejected |
| Jotai | Atomic, minimal | Less familiar pattern | ❌ Rejected |

---

## 8. Connection Handling & Resilience

### Decision: Auto-reconnect with exponential backoff

**Rationale**:
- Network interruptions are common, especially on WiFi
- Automatic reconnection improves user experience
- Exponential backoff prevents server overload
- Game pauses during disconnection, resumes on reconnect

**Reconnection Strategy**:
```text
On disconnect:
  1. Pause game, show "Reconnecting..." overlay
  2. Attempt reconnect with exponential backoff (1s, 2s, 4s, max 30s)
  3. On success: Request current game state, resume
  4. After 5 failures: Show "Connection lost" with manual retry option
```

---

## 9. Rendering Approach

### Decision: Canvas-based rendering with React wrapper

**Rationale**:
- Canvas provides optimal performance for game rendering at 60 FPS
- React wrapper manages component lifecycle and state integration
- Avoids DOM manipulation overhead of SVG/HTML elements
- Easy sprite/animation handling

**Implementation**:
- React component owns canvas element
- useEffect hook handles render loop
- Zustand subscription triggers re-renders on state change
- RequestAnimationFrame for smooth 60 FPS

---

## 10. Testing Strategy

### Decision: Unit tests for game logic, integration tests for API

**Backend Testing**:
- Unit tests for domain logic (movement, collision, ghost AI, scoring)
- Integration tests for WebSocket message handling
- Table-driven tests for edge cases

**Frontend Testing**:
- Unit tests for hooks and stores
- Component tests with React Testing Library
- Integration tests for WebSocket client

**Coverage Target**: >80% on critical paths (game logic, state management)

---

## Summary of Key Decisions

| Area | Decision | Key Rationale |
| ---- | -------- | ------------- |
| Communication | WebSocket | Low latency, bidirectional |
| Game Loop | 20 TPS server-side | Balance of responsiveness and load |
| Ghost AI | Classic patterns | Authentic experience |
| Maze | 2D grid array | Simple, efficient |
| State Sync | Full state broadcast | Simplicity, server authority |
| Input | Immediate send + server validation | Responsiveness + integrity |
| Frontend State | Zustand | Lightweight, hooks-native |
| Reconnection | Auto with exponential backoff | User experience |
| Rendering | Canvas with React wrapper | Performance |
| Testing | Unit + Integration | TDD per constitution |

