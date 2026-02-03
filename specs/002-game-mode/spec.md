# Feature Specification: Pac-Man Game Mode

**Feature Branch**: `002-game-mode`  
**Created**: 2026-02-03  
**Status**: Draft  
**Input**: User description: "Game mode. After user has clicked StartButton game starts. A predefined pacman map and commonly known rules apply here"

## Architecture Constraints

This feature follows a **server-authoritative architecture**:

### Backend Responsibilities (Game Server)

- **All game logic**: Movement validation, collision detection, scoring calculations
- **Ghost AI**: All ghost behavior, pathfinding, and state management
- **Game state management**: Lives, score, level progression, win/lose conditions
- **Map storage**: The maze layout is stored and served by the backend
- **Game tick/loop**: The backend runs the game simulation

### Frontend Responsibilities (Game Client)

- **Rendering**: Display the game state received from the backend
- **Input capture**: Capture keyboard inputs from the player
- **Input transmission**: Send player inputs to the backend
- **State visualization**: Show score, lives, game status based on backend data
- **No game logic**: The frontend does NOT calculate movement, collisions, or any game rules

### Communication Pattern

- Frontend requests game start → Backend initializes game and returns initial state (including map)
- Frontend sends player input (direction) → Backend processes and returns updated game state
- Backend continuously sends game state updates → Frontend renders each update

## Clarifications

### Session 2026-02-03

- Q: How should players trigger the game pause? → A: Escape key (no HUD pause button)

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Start Game and Move Pac-Man (Priority: P1)

As a player, I want to start the game by clicking the Start button and control Pac-Man using keyboard inputs, so that I can navigate through the maze.

**Why this priority**: This is the core interaction - without the ability to start the game and control Pac-Man, no gameplay is possible. This delivers the fundamental playable experience.

**Independent Test**: Can be fully tested by clicking Start, verifying the game screen appears with a maze and Pac-Man, and confirming arrow key inputs move Pac-Man in the corresponding direction. Delivers the core interactive experience.

**Acceptance Scenarios**:

1. **Given** the main page is displayed with the Start button, **When** the player clicks "START GAME", **Then** the game screen loads showing the predefined Pac-Man maze with Pac-Man positioned at the starting location
2. **Given** the game is running, **When** the player presses an arrow key (Up/Down/Left/Right), **Then** Pac-Man moves in that direction if the path is not blocked by a wall
3. **Given** the game is running and Pac-Man is moving, **When** the player presses a different arrow key, **Then** Pac-Man changes direction at the next valid intersection
4. **Given** Pac-Man is moving toward a wall, **When** Pac-Man reaches the wall, **Then** Pac-Man stops and waits for a new valid direction input

---

### User Story 2 - Collect Dots and Power Pellets (Priority: P1)

As a player, I want to eat dots and power pellets as I navigate the maze, so that I can earn points and gain temporary ghost-eating ability.

**Why this priority**: Collecting dots is the primary objective of Pac-Man and the core scoring mechanism. Without this, there is no game goal.

**Independent Test**: Can be tested by moving Pac-Man over dots and power pellets, verifying they disappear, score increases, and power pellet triggers the vulnerability state.

**Acceptance Scenarios**:

1. **Given** the game is running and Pac-Man moves over a dot, **When** Pac-Man's position overlaps the dot, **Then** the dot disappears and the score increases by a fixed amount
2. **Given** the game is running and Pac-Man moves over a power pellet, **When** Pac-Man's position overlaps the power pellet, **Then** the power pellet disappears, the score increases, and ghosts enter vulnerable state for a limited time
3. **Given** ghosts are in vulnerable state, **When** the vulnerability timer expires, **Then** ghosts return to their normal dangerous state

---

### User Story 3 - Ghost Encounters (Priority: P1)

As a player, I want ghosts to chase me through the maze, so that I have a challenge to overcome while collecting dots.

**Why this priority**: Ghosts are essential to the Pac-Man gameplay challenge. Without them, there would be no risk or tension.

**Independent Test**: Can be tested by observing ghost movement patterns and verifying collision behavior (losing a life when touched by normal ghost, earning points when eating a vulnerable ghost).

**Acceptance Scenarios**:

1. **Given** the game is running, **When** the game starts, **Then** four ghosts appear at their designated starting positions in the maze
2. **Given** ghosts are in normal state, **When** a ghost touches Pac-Man, **Then** Pac-Man loses one life and the game resets positions (if lives remain) or shows game over (if no lives remain)
3. **Given** ghosts are in vulnerable state (after power pellet), **When** Pac-Man touches a vulnerable ghost, **Then** the ghost is eaten, the score increases significantly, and the ghost returns to the starting area
4. **Given** a ghost was eaten, **When** the ghost reaches the starting area, **Then** the ghost regenerates and resumes chasing Pac-Man in normal state

---

### User Story 4 - Score and Lives Display (Priority: P2)

As a player, I want to see my current score and remaining lives displayed during gameplay, so that I can track my progress and know how many attempts I have left.

**Why this priority**: While the game is playable without visible score/lives, displaying them is essential for a complete game experience and player feedback.

**Independent Test**: Can be tested by verifying the UI shows score (updating when dots are collected) and lives (updating when Pac-Man dies).

**Acceptance Scenarios**:

1. **Given** the game starts, **When** the game screen loads, **Then** the score displays as zero and the lives display shows the starting number of lives (3 lives)
2. **Given** the game is running, **When** Pac-Man collects a dot or power pellet, **Then** the displayed score updates immediately to reflect the new total
3. **Given** Pac-Man has lives remaining, **When** Pac-Man loses a life, **Then** the lives display decreases by one

---

### User Story 5 - Level Completion and Game Over (Priority: P2)

As a player, I want the level to complete when I eat all dots and the game to end when I lose all lives, so that I have clear win/lose conditions.

**Why this priority**: Win/lose conditions provide closure to the gameplay session, but the core gameplay loop can function for initial testing without them.

**Independent Test**: Can be tested by collecting all dots (level complete) or losing all lives (game over) and verifying the appropriate end state is displayed.

**Acceptance Scenarios**:

1. **Given** the game is running, **When** Pac-Man collects the last remaining dot in the maze, **Then** a level complete message is displayed and the level restarts with all dots replenished
2. **Given** Pac-Man has one life remaining, **When** a ghost catches Pac-Man, **Then** a "Game Over" screen is displayed with the final score
3. **Given** the game over screen is displayed, **When** the player views the screen, **Then** an option to return to the main menu or restart is available

---

### Edge Cases

**Input Handling**:

- What happens when the player presses multiple arrow keys simultaneously? → The most recent key press takes priority
- How does the system handle rapid key presses? → Direction changes queue for the next valid intersection
- What happens if the player presses an arrow key toward a wall? → Pac-Man continues in current direction until a valid turn is available
- What happens if the browser window loses focus during gameplay? → Game pauses automatically
- How does the player pause the game? → Pressing the Escape key toggles between playing and paused states (no HUD pause button)

**Game Logic**:

- What happens if all ghosts are eaten during one power pellet? → Each ghost eaten increases in point value (ghost combo scoring)

**Network/Communication**:

- What happens if the connection to the backend is lost during gameplay? → Game pauses and displays a "Reconnecting..." message; game resumes from last known state when connection restores
- What happens if player input cannot be sent to the backend? → Input is queued locally and sent when connection restores; visual feedback shows pending state
- What happens if the backend cannot be reached when starting a game? → Display an error message and return to the main menu

## Requirements _(mandatory)_

### Functional Requirements

**Game Initialization**:

- **FR-001**: Frontend MUST request game start from the backend when player clicks "START GAME"
- **FR-002**: Backend MUST initialize the game state and return the maze layout, Pac-Man position, ghost positions, and initial score/lives
- **FR-003**: Frontend MUST render the maze and all game entities based on data received from the backend

**Player Input & Movement**:

- **FR-004**: Frontend MUST capture arrow key inputs (Up, Down, Left, Right) and send them to the backend
- **FR-004a**: Frontend MUST capture Escape key to toggle pause state and send pause/resume command to backend
- **FR-005**: Backend MUST validate movement inputs and update Pac-Man's position according to game rules
- **FR-006**: Backend MUST prevent Pac-Man from moving through walls
- **FR-007**: Backend MUST send updated game state to frontend after processing each input

**Collectibles & Scoring**:

- **FR-008**: Backend MUST detect when Pac-Man's position overlaps a dot and update the score
- **FR-009**: Backend MUST detect when Pac-Man's position overlaps a power pellet, update the score, and trigger ghost vulnerability
- **FR-010**: Backend MUST track and manage the vulnerability timer for ghosts

**Ghost Behavior**:

- **FR-011**: Backend MUST control all ghost movement and AI behavior
- **FR-012**: Backend MUST detect collisions between Pac-Man and ghosts
- **FR-013**: Backend MUST handle ghost state transitions (normal → vulnerable → eaten → respawning)

**Game State & Display**:

- **FR-014**: Frontend MUST display the current score received from backend
- **FR-015**: Frontend MUST display remaining lives received from backend
- **FR-016**: Backend MUST determine when Pac-Man loses a life and update the game state
- **FR-017**: Backend MUST determine "Game Over" condition (no lives remaining)
- **FR-018**: Backend MUST determine level completion (all dots collected)
- **FR-019**: Frontend MUST display "Game Over" or "Level Complete" screens based on backend state

**Game Flow**:

- **FR-020**: Backend MUST reset positions after Pac-Man loses a life (if lives remain)
- **FR-021**: Frontend MUST provide option to return to main menu or restart (triggering appropriate backend calls)

### Key Entities

- **Pac-Man**: The player-controlled character; has position, direction, and animation state
- **Ghost**: Enemy characters (4 total); have position, movement behavior, and state (normal, vulnerable, eaten)
- **Maze**: The game board; contains walls, corridors, dot positions, power pellet positions, and starting positions
- **Dot**: Collectible item worth points; has position and collected state
- **Power Pellet**: Special collectible that triggers ghost vulnerability; has position and collected state
- **Game State**: Tracks score, lives remaining, level, and game status (playing, paused, game over, level complete)

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: Players can start gameplay within 2 seconds of clicking the Start button (including map fetch from backend)
- **SC-002**: Pac-Man responds to keyboard input within 150 milliseconds (including round-trip to backend), providing responsive controls
- **SC-003**: Ghost movement appears smooth and consistent at 60 frames per second on standard hardware
- **SC-004**: Players can complete a full game session (win or lose) without encountering errors or crashes
- **SC-005**: Score updates are displayed within 200 milliseconds of collecting items (accounting for backend processing)
- **SC-006**: The game correctly tracks and displays all 240+ dots in the standard maze layout
- **SC-007**: Players can pause and resume gameplay by pressing the Escape key without losing game state
- **SC-008**: The game runs consistently across modern desktop browsers (Chrome, Firefox, Safari, Edge)
- **SC-009**: The game gracefully handles temporary network interruptions without losing player progress

## Assumptions

**Architecture**:

- Server-authoritative architecture: all game logic runs on the backend
- The frontend is a thin client that only renders and captures input
- The frontend must use as much reusable components as possible
- The maze layout is stored on and fetched from the backend
- Real-time communication between frontend and backend for game state updates

**Game Rules**:

- The game uses a classic Pac-Man maze layout (standard 28x31 tile grid)
- The game starts with 3 lives (industry standard)
- Standard scoring is used: 10 points per dot, 50 points per power pellet, 200-1600 points for eating ghosts (combo scaling)
- Ghost vulnerability duration is approximately 6-8 seconds (standard timing)
- The four ghosts follow their traditional chase/scatter behavior patterns

**Technical**:

- Arrow keys are the primary control method for movement
- Escape key toggles pause (no HUD pause button)
- The game targets 60 FPS for smooth rendering on the frontend
- Backend game tick rate supports smooth gameplay experience
- Single-player mode only (no multiplayer specified)
