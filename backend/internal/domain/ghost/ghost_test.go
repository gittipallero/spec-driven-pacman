package ghost

import (
	"testing"

	"github.com/spec-driven-pacman/backend/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestNewGhost(t *testing.T) {
	startPos := types.Position{X: 13, Y: 14}
	scatterTarget := types.Position{X: 25, Y: 0}
	
	ghost := NewGhost(GhostBlinky, startPos, scatterTarget)
	
	assert.Equal(t, GhostBlinky, ghost.Name)
	assert.Equal(t, startPos, ghost.Position)
	assert.Equal(t, startPos, ghost.StartPos)
	assert.Equal(t, scatterTarget, ghost.ScatterTarget)
	assert.Equal(t, GhostStateNormal, ghost.State)
	assert.Equal(t, GhostModeScatter, ghost.Mode)
	assert.Equal(t, types.DirectionUp, ghost.Direction)
}

func TestGhostStateTransitions(t *testing.T) {
	t.Run("normal to vulnerable", func(t *testing.T) {
		ghost := NewGhost(GhostBlinky, types.Position{X: 13, Y: 14}, types.Position{X: 25, Y: 0})
		ghost.State = GhostStateNormal
		ghost.Direction = types.DirectionRight
		
		ghost.SetVulnerable()
		
		assert.Equal(t, GhostStateVulnerable, ghost.State)
		assert.Equal(t, GhostModeFrightened, ghost.Mode)
		assert.Equal(t, types.DirectionLeft, ghost.Direction) // Reversed
	})

	t.Run("vulnerable to eaten", func(t *testing.T) {
		ghost := NewGhost(GhostBlinky, types.Position{X: 13, Y: 14}, types.Position{X: 25, Y: 0})
		ghost.SetVulnerable()
		
		ghost.SetEaten()
		
		assert.Equal(t, GhostStateEaten, ghost.State)
	})

	t.Run("eaten to respawning", func(t *testing.T) {
		ghost := NewGhost(GhostBlinky, types.Position{X: 13, Y: 14}, types.Position{X: 25, Y: 0})
		ghost.SetEaten()
		
		ghost.SetRespawning()
		
		assert.Equal(t, GhostStateRespawning, ghost.State)
	})

	t.Run("respawning to normal", func(t *testing.T) {
		ghost := NewGhost(GhostBlinky, types.Position{X: 13, Y: 14}, types.Position{X: 25, Y: 0})
		ghost.SetRespawning()
		
		ghost.SetNormal()
		
		assert.Equal(t, GhostStateNormal, ghost.State)
		assert.Equal(t, GhostModeChase, ghost.Mode)
	})

	t.Run("vulnerable ignored when not normal", func(t *testing.T) {
		ghost := NewGhost(GhostBlinky, types.Position{X: 13, Y: 14}, types.Position{X: 25, Y: 0})
		ghost.SetEaten()
		originalState := ghost.State
		
		ghost.SetVulnerable()
		
		assert.Equal(t, originalState, ghost.State) // Should not change
	})
}

func TestGhostReset(t *testing.T) {
	ghost := NewGhost(GhostBlinky, types.Position{X: 13, Y: 14}, types.Position{X: 25, Y: 0})
	ghost.Position = types.Position{X: 5, Y: 5}
	ghost.State = GhostStateEaten
	ghost.Mode = GhostModeChase
	ghost.Direction = types.DirectionDown
	
	ghost.Reset()
	
	assert.Equal(t, types.Position{X: 13, Y: 14}, ghost.Position)
	assert.Equal(t, GhostStateNormal, ghost.State)
	assert.Equal(t, GhostModeScatter, ghost.Mode)
	assert.Equal(t, types.DirectionUp, ghost.Direction)
}

func TestGhostIsAtHome(t *testing.T) {
	startPos := types.Position{X: 13, Y: 14}
	ghost := NewGhost(GhostBlinky, startPos, types.Position{X: 25, Y: 0})
	
	t.Run("returns true at start position", func(t *testing.T) {
		assert.True(t, ghost.IsAtHome())
	})

	t.Run("returns false away from start", func(t *testing.T) {
		ghost.Position = types.Position{X: 10, Y: 10}
		assert.False(t, ghost.IsAtHome())
	})
}

func TestGhostReverseDirection(t *testing.T) {
	ghost := NewGhost(GhostBlinky, types.Position{X: 13, Y: 14}, types.Position{X: 25, Y: 0})
	
	ghost.Direction = types.DirectionRight
	ghost.ReverseDirection()
	assert.Equal(t, types.DirectionLeft, ghost.Direction)
	
	ghost.ReverseDirection()
	assert.Equal(t, types.DirectionRight, ghost.Direction)
	
	ghost.Direction = types.DirectionUp
	ghost.ReverseDirection()
	assert.Equal(t, types.DirectionDown, ghost.Direction)
}

func TestCreateAllGhosts(t *testing.T) {
	ghosts := CreateAllGhosts()
	
	assert.Len(t, ghosts, 4)
	
	// Check each ghost
	names := make(map[GhostName]bool)
	for _, g := range ghosts {
		names[g.Name] = true
		assert.Equal(t, GhostStateNormal, g.State)
		assert.Equal(t, GhostModeScatter, g.Mode)
	}
	
	assert.True(t, names[GhostBlinky])
	assert.True(t, names[GhostPinky])
	assert.True(t, names[GhostInky])
	assert.True(t, names[GhostClyde])
}

func TestGhostToJSON(t *testing.T) {
	ghost := NewGhost(GhostBlinky, types.Position{X: 13, Y: 14}, types.Position{X: 25, Y: 0})
	ghost.Direction = types.DirectionRight
	ghost.State = GhostStateVulnerable
	ghost.Mode = GhostModeFrightened
	
	json := ghost.ToJSON()
	
	assert.Equal(t, GhostBlinky, json.Name)
	assert.Equal(t, 13, json.X)
	assert.Equal(t, 14, json.Y)
	assert.Equal(t, types.DirectionRight, json.Direction)
	assert.Equal(t, GhostStateVulnerable, json.State)
	assert.Equal(t, GhostModeFrightened, json.Mode)
}

