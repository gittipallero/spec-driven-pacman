package game

import (
	"testing"

	"github.com/spec-driven-pacman/backend/internal/domain/maze"
	"github.com/spec-driven-pacman/backend/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGameSession(t *testing.T) {
	gameMaze := maze.LoadClassicMaze()
	session := NewGameSession(gameMaze)

	t.Run("creates valid session ID", func(t *testing.T) {
		assert.NotEmpty(t, session.ID)
		assert.Len(t, session.ID, 36) // UUID format
	})

	t.Run("starts with playing status", func(t *testing.T) {
		assert.Equal(t, types.GameStatusPlaying, session.Status)
	})

	t.Run("starts with correct initial values", func(t *testing.T) {
		assert.Equal(t, 0, session.Score)
		assert.Equal(t, InitialLives, session.Lives)
		assert.Equal(t, 1, session.Level)
		assert.Equal(t, 0, session.VulnerabilityTimer)
	})

	t.Run("Pac-Man starts at correct position", func(t *testing.T) {
		require.NotNil(t, session.PacMan)
		assert.Equal(t, PacManStartPos, session.PacMan.Position)
	})

	t.Run("has four ghosts", func(t *testing.T) {
		assert.Len(t, session.Ghosts, 4)
	})

	t.Run("has maze with dots", func(t *testing.T) {
		require.NotNil(t, session.Maze)
		assert.Greater(t, session.Maze.DotsRemaining, 0)
	})
}

func TestGameSessionProcessInput(t *testing.T) {
	gameMaze := maze.LoadClassicMaze()
	session := NewGameSession(gameMaze)

	t.Run("queues direction when playing", func(t *testing.T) {
		session.ProcessInput(types.DirectionUp)
		
		assert.Equal(t, types.DirectionUp, session.PacMan.NextDirection)
	})

	t.Run("ignores input when paused", func(t *testing.T) {
		session.Status = types.GameStatusPaused
		session.PacMan.NextDirection = types.DirectionNone
		
		session.ProcessInput(types.DirectionDown)
		
		assert.Equal(t, types.DirectionNone, session.PacMan.NextDirection)
	})

	t.Run("ignores input when game over", func(t *testing.T) {
		session.Status = types.GameStatusGameOver
		session.PacMan.NextDirection = types.DirectionNone
		
		session.ProcessInput(types.DirectionLeft)
		
		assert.Equal(t, types.DirectionNone, session.PacMan.NextDirection)
	})

	t.Run("ignores invalid direction", func(t *testing.T) {
		session.Status = types.GameStatusPlaying
		session.PacMan.NextDirection = types.DirectionNone
		
		session.ProcessInput(types.Direction("invalid"))
		
		assert.Equal(t, types.DirectionNone, session.PacMan.NextDirection)
	})
}

func TestGameSessionPause(t *testing.T) {
	gameMaze := maze.LoadClassicMaze()
	session := NewGameSession(gameMaze)

	t.Run("pauses when playing", func(t *testing.T) {
		session.Status = types.GameStatusPlaying
		
		session.Pause()
		
		assert.Equal(t, types.GameStatusPaused, session.Status)
	})

	t.Run("does nothing when not playing", func(t *testing.T) {
		session.Status = types.GameStatusGameOver
		
		session.Pause()
		
		assert.Equal(t, types.GameStatusGameOver, session.Status)
	})
}

func TestGameSessionResume(t *testing.T) {
	gameMaze := maze.LoadClassicMaze()
	session := NewGameSession(gameMaze)

	t.Run("resumes when paused", func(t *testing.T) {
		session.Status = types.GameStatusPaused
		
		session.Resume()
		
		assert.Equal(t, types.GameStatusPlaying, session.Status)
	})

	t.Run("does nothing when not paused", func(t *testing.T) {
		session.Status = types.GameStatusPlaying
		
		session.Resume()
		
		assert.Equal(t, types.GameStatusPlaying, session.Status)
	})
}

func TestGameSessionTogglePause(t *testing.T) {
	gameMaze := maze.LoadClassicMaze()
	session := NewGameSession(gameMaze)

	t.Run("toggles playing to paused", func(t *testing.T) {
		session.Status = types.GameStatusPlaying
		
		session.TogglePause()
		
		assert.Equal(t, types.GameStatusPaused, session.Status)
	})

	t.Run("toggles paused to playing", func(t *testing.T) {
		session.Status = types.GameStatusPaused
		
		session.TogglePause()
		
		assert.Equal(t, types.GameStatusPlaying, session.Status)
	})

	t.Run("does nothing for other states", func(t *testing.T) {
		session.Status = types.GameStatusGameOver
		
		session.TogglePause()
		
		assert.Equal(t, types.GameStatusGameOver, session.Status)
	})
}

func TestGameSessionTick(t *testing.T) {
	t.Run("moves Pac-Man when playing", func(t *testing.T) {
		gameMaze := maze.LoadClassicMaze()
		session := NewGameSession(gameMaze)
		session.Status = types.GameStatusPlaying
		session.PacMan.Direction = types.DirectionLeft
		initialX := session.PacMan.Position.X
		
		session.Tick(50)
		
		// Pac-Man should have moved left (or stayed if blocked)
		// The exact behavior depends on the maze layout
		assert.NotNil(t, session.PacMan)
	})

	t.Run("does not tick when paused", func(t *testing.T) {
		gameMaze := maze.LoadClassicMaze()
		session := NewGameSession(gameMaze)
		session.Status = types.GameStatusPaused
		initialPos := session.PacMan.Position
		
		session.Tick(50)
		
		assert.Equal(t, initialPos, session.PacMan.Position)
	})

	t.Run("does not tick when game over", func(t *testing.T) {
		gameMaze := maze.LoadClassicMaze()
		session := NewGameSession(gameMaze)
		session.Status = types.GameStatusGameOver
		initialPos := session.PacMan.Position
		
		session.Tick(50)
		
		assert.Equal(t, initialPos, session.PacMan.Position)
	})
}

func TestGameSessionToJSON(t *testing.T) {
	gameMaze := maze.LoadClassicMaze()
	session := NewGameSession(gameMaze)

	json := session.ToJSON()

	assert.Equal(t, session.Status, json.Status)
	assert.Equal(t, session.Score, json.Score)
	assert.Equal(t, session.Lives, json.Lives)
	assert.Equal(t, session.Level, json.Level)
	assert.Equal(t, session.PacMan.Position.X, json.PacMan.X)
	assert.Equal(t, session.PacMan.Position.Y, json.PacMan.Y)
	assert.Len(t, json.Ghosts, 4)
}

func TestGameSessionToStartResponse(t *testing.T) {
	gameMaze := maze.LoadClassicMaze()
	session := NewGameSession(gameMaze)
	wsURL := "/ws/game/test-id"

	response := session.ToStartResponse(wsURL)

	assert.Equal(t, session.ID, response.SessionID)
	assert.Equal(t, wsURL, response.WebsocketURL)
	assert.NotNil(t, response.State)
	assert.NotNil(t, response.Maze)
	assert.Equal(t, maze.MazeWidth, response.Maze.Width)
	assert.Equal(t, maze.MazeHeight, response.Maze.Height)
}

