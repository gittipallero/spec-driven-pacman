# API Contracts: Arcade Main Page

**Feature**: 001-main-page  
**Date**: 2026-02-03

## Overview

This feature does not require any API contracts. The main page is a static frontend component with no backend communication.

## Why No Contracts?

The main page feature:

- Displays static content (title, button)
- Has no user authentication
- Does not fetch or persist data
- Does not communicate with any backend service

## Future Contracts

When game features are implemented, the following contracts will be defined:

| Contract               | Method    | Purpose                      |
| ---------------------- | --------- | ---------------------------- |
| `POST /api/game/start` | REST      | Initialize new game session  |
| `GET /api/scores/high` | REST      | Fetch high score leaderboard |
| `WS /ws/game`          | WebSocket | Real-time game state sync    |

These will be defined in their respective feature specifications.

## Navigation Contract (Internal)

While not an API contract, the main page establishes the following navigation interface:

### Start Game Action

```typescript
// When user clicks "Start Game", navigate to game route
// This is a frontend-only routing action

interface NavigationAction {
  type: "NAVIGATE";
  destination: "/game"; // Future game route
}
```

For MVP, the Start Game button will navigate to a placeholder route until game functionality is implemented.
