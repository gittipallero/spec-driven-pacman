# API Contracts: Pac-Man Game Mode

**Feature**: 002-game-mode | **Date**: 2026-02-03

## Overview

This directory contains API contracts for communication between the React frontend and Go backend.

## Files

| File | Description |
| ---- | ----------- |
| `rest-api.yaml` | OpenAPI 3.0 specification for REST endpoints |
| `websocket-messages.json` | JSON Schema for WebSocket message types |

## Communication Pattern

```text
┌──────────────┐                    ┌──────────────┐
│   Frontend   │                    │   Backend    │
│   (React)    │                    │    (Go)      │
└──────┬───────┘                    └──────┬───────┘
       │                                   │
       │  POST /api/game/start             │
       │──────────────────────────────────▶│
       │  { maze, pacman, ghosts, ... }    │
       │◀──────────────────────────────────│
       │                                   │
       │  WebSocket /ws/game/:sessionId    │
       │══════════════════════════════════▶│
       │                                   │
       │  input: { direction: "up" }       │
       │──────────────────────────────────▶│
       │                                   │
       │  state: { pacman, ghosts, ... }   │
       │◀──────────────────────────────────│
       │  (continuous at ~20 TPS)          │
       │                                   │
```

## Endpoints Summary

### REST API

| Method | Endpoint | Description |
| ------ | -------- | ----------- |
| POST | `/api/game/start` | Initialize new game session |
| GET | `/api/game/:sessionId` | Get current game state |
| POST | `/api/game/:sessionId/pause` | Pause game |
| POST | `/api/game/:sessionId/resume` | Resume game |
| DELETE | `/api/game/:sessionId` | End game session |

### WebSocket

| Direction | Message Type | Description |
| --------- | ------------ | ----------- |
| Client → Server | `input` | Player direction input |
| Client → Server | `pause` | Pause request |
| Client → Server | `resume` | Resume request |
| Server → Client | `state` | Game state update (~20/sec) |
| Server → Client | `event` | Game events (death, level complete) |
| Server → Client | `error` | Error notification |

## Validation

- All inputs validated server-side
- Invalid messages return error response
- Rate limiting: max 60 input messages per second per session

