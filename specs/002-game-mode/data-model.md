# Data Model: Pac-Man Game Mode

**Feature**: 002-game-mode | **Date**: 2026-02-03

## Overview

This document defines the data entities, their attributes, relationships, and state transitions for the Pac-Man game mode.

---

## Entities

### 1. Game Session

Represents a single game session from start to game over.

| Field | Type | Description |
| ----- | ---- | ----------- |
| id | string | Unique session identifier (UUID) |
| status | GameStatus | Current game status |
| score | integer | Current score (0+) |
| lives | integer | Remaining lives (0-3) |
| level | integer | Current level (1+) |
| vulnerabilityTimer | integer | Remaining ghost vulnerability time (ms) |
| createdAt | timestamp | Session start time |

**GameStatus Enum**:
- `waiting` - Game initialized, not started
- `playing` - Active gameplay
- `paused` - Game paused
- `gameOver` - No lives remaining
- `levelComplete` - All dots collected

---

### 2. Pac-Man

The player-controlled character.

| Field | Type | Description |
| ----- | ---- | ----------- |
| x | integer | Grid X position (0-27) |
| y | integer | Grid Y position (0-30) |
| direction | Direction | Current movement direction |
| nextDirection | Direction | Queued direction for next valid turn |
| state | PacManState | Animation/visual state |

**Direction Enum**: `up`, `down`, `left`, `right`, `none`

**PacManState Enum**:
- `moving` - Normal movement
- `dying` - Death animation playing
- `idle` - Not moving (waiting for input)

---

### 3. Ghost

Enemy characters that pursue Pac-Man.

| Field | Type | Description |
| ----- | ---- | ----------- |
| name | GhostName | Ghost identity |
| x | integer | Grid X position (0-27) |
| y | integer | Grid Y position (0-30) |
| direction | Direction | Current movement direction |
| state | GhostState | Current behavior state |
| mode | GhostMode | AI behavior mode |

**GhostName Enum**: `blinky`, `pinky`, `inky`, `clyde`

**GhostState Enum**:
- `normal` - Standard chasing behavior
- `vulnerable` - Can be eaten by Pac-Man
- `eaten` - Returning to ghost house
- `respawning` - Exiting ghost house

**GhostMode Enum**:
- `chase` - Actively targeting Pac-Man
- `scatter` - Moving to corner territory
- `frightened` - Random movement (when vulnerable)

---

### 4. Maze

The game board layout.

| Field | Type | Description |
| ----- | ---- | ----------- |
| width | integer | Grid width (28) |
| height | integer | Grid height (31) |
| tiles | integer[][] | 2D array of tile types |
| dotCount | integer | Total dots in maze |
| dotsRemaining | integer | Uncollected dots |

**Tile Types**:
```
0 = EMPTY      - Walkable corridor
1 = WALL       - Impassable
2 = DOT        - Regular dot (10 points)
3 = POWER      - Power pellet (50 points)
4 = GHOST_HOUSE - Ghost spawn area
5 = TUNNEL     - Wraps to opposite side
```

---

### 5. Dot

Collectible item in the maze.

| Field | Type | Description |
| ----- | ---- | ----------- |
| x | integer | Grid X position |
| y | integer | Grid Y position |
| type | DotType | Regular dot or power pellet |
| collected | boolean | Whether collected |

**DotType Enum**: `dot`, `powerPellet`

**Point Values**:
- Regular dot: 10 points
- Power pellet: 50 points

---

## Relationships

```text
┌─────────────────┐
│  Game Session   │
├─────────────────┤
│ Contains 1      │────────┐
│ Pac-Man         │        │
├─────────────────┤        ▼
│ Contains 4      │    ┌─────────┐
│ Ghosts          │───▶│ Pac-Man │
├─────────────────┤    └─────────┘
│ Contains 1      │        
│ Maze            │───▶┌─────────┐
├─────────────────┤    │  Ghost  │ x4
│ Has many        │    └─────────┘
│ Dots            │        
└─────────────────┘───▶┌─────────┐
                       │  Maze   │
                       └────┬────┘
                            │
                            ▼
                       ┌─────────┐
                       │  Dot    │ x240+
                       └─────────┘
```

---

## State Transitions

### Game Session State Machine

```text
                    ┌──────────┐
                    │ waiting  │
                    └────┬─────┘
                         │ start game
                         ▼
                    ┌──────────┐
       ┌───────────│ playing  │◄──────────┐
       │           └────┬─────┘           │
       │ pause          │                 │ resume
       ▼                │                 │
  ┌──────────┐          │            ┌────┴─────┐
  │  paused  │──────────┼───────────▶│  paused  │
  └──────────┘          │            └──────────┘
                        │
         ┌──────────────┼──────────────┐
         │ lives = 0    │ dots = 0     │
         ▼              │              ▼
    ┌──────────┐        │      ┌─────────────────┐
    │ gameOver │        │      │ levelComplete   │
    └──────────┘        │      └────────┬────────┘
                        │               │ next level
                        │               ▼
                        └──────────────────
```

### Ghost State Machine

```text
                    ┌──────────┐
                    │  normal  │◄─────────────────┐
                    └────┬─────┘                  │
                         │ power pellet eaten     │
                         ▼                        │
                    ┌────────────┐                │
                    │ vulnerable │────────────────┤
                    └────┬───────┘ timer expires  │
                         │ eaten by Pac-Man       │
                         ▼                        │
                    ┌──────────┐                  │
                    │  eaten   │                  │
                    └────┬─────┘                  │
                         │ reached ghost house    │
                         ▼                        │
                    ┌────────────┐                │
                    │ respawning │────────────────┘
                    └────────────┘ exit complete
```

### Pac-Man State Machine

```text
    ┌──────────┐
    │   idle   │◄──────────────────────┐
    └────┬─────┘                       │
         │ direction input             │ hit wall / no input
         ▼                             │
    ┌──────────┐                       │
    │  moving  │───────────────────────┘
    └────┬─────┘
         │ collision with ghost (normal)
         ▼
    ┌──────────┐
    │  dying   │
    └────┬─────┘
         │ animation complete
         ▼
    [respawn at start position or game over]
```

---

## Validation Rules

### Game Session
- `score` must be >= 0
- `lives` must be 0-3
- `level` must be >= 1
- `vulnerabilityTimer` must be >= 0

### Position
- `x` must be 0-27 (maze width - 1)
- `y` must be 0-30 (maze height - 1)
- Position must not be on a WALL tile

### Movement
- Pac-Man can only move in cardinal directions
- Movement blocked if target tile is WALL
- Tunnel tiles wrap to opposite side (x: 0 ↔ 27)

### Collision
- Pac-Man collides with ghost if same tile position
- Pac-Man collects dot if on same tile position

### Scoring
- Dot: +10 points
- Power pellet: +50 points
- Ghost eaten (combo): 200, 400, 800, 1600 points

---

## Initial State

When a new game session starts:

```text
Game Session:
  - status: "playing"
  - score: 0
  - lives: 3
  - level: 1
  - vulnerabilityTimer: 0

Pac-Man:
  - position: (13, 23) [center bottom]
  - direction: "left"
  - state: "idle"

Ghosts:
  - Blinky: (13, 11), "normal", "scatter"
  - Pinky: (13, 14), "normal", "scatter"  [in ghost house]
  - Inky: (11, 14), "normal", "scatter"   [in ghost house]
  - Clyde: (15, 14), "normal", "scatter"  [in ghost house]

Maze:
  - Standard 28x31 layout
  - 240 dots + 4 power pellets
```

