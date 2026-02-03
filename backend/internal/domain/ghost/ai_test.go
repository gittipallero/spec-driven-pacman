package ghost

import (
	"testing"

	"github.com/spec-driven-pacman/backend/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestBlinkyTarget(t *testing.T) {
	t.Run("targets Pac-Man position directly", func(t *testing.T) {
		pacmanPos := types.Position{X: 10, Y: 15}
		
		target := BlinkyTarget(pacmanPos)
		
		assert.Equal(t, pacmanPos, target)
	})
}

func TestPinkyTarget(t *testing.T) {
	t.Run("targets 4 tiles ahead of Pac-Man moving right", func(t *testing.T) {
		pacmanPos := types.Position{X: 10, Y: 15}
		pacmanDir := types.DirectionRight
		
		target := PinkyTarget(pacmanPos, pacmanDir)
		
		assert.Equal(t, types.Position{X: 14, Y: 15}, target)
	})

	t.Run("targets 4 tiles ahead of Pac-Man moving up", func(t *testing.T) {
		pacmanPos := types.Position{X: 10, Y: 15}
		pacmanDir := types.DirectionUp
		
		target := PinkyTarget(pacmanPos, pacmanDir)
		
		// Classic bug: also moves 4 tiles left when moving up
		assert.Equal(t, types.Position{X: 6, Y: 11}, target)
	})

	t.Run("targets 4 tiles ahead of Pac-Man moving down", func(t *testing.T) {
		pacmanPos := types.Position{X: 10, Y: 15}
		pacmanDir := types.DirectionDown
		
		target := PinkyTarget(pacmanPos, pacmanDir)
		
		assert.Equal(t, types.Position{X: 10, Y: 19}, target)
	})

	t.Run("targets 4 tiles ahead of Pac-Man moving left", func(t *testing.T) {
		pacmanPos := types.Position{X: 10, Y: 15}
		pacmanDir := types.DirectionLeft
		
		target := PinkyTarget(pacmanPos, pacmanDir)
		
		assert.Equal(t, types.Position{X: 6, Y: 15}, target)
	})
}

func TestInkyTarget(t *testing.T) {
	t.Run("targets based on Blinky and Pac-Man positions", func(t *testing.T) {
		pacmanPos := types.Position{X: 10, Y: 15}
		pacmanDir := types.DirectionRight
		blinkyPos := types.Position{X: 8, Y: 13}
		
		target := InkyTarget(pacmanPos, pacmanDir, blinkyPos)
		
		// Inky's target is the vector from Blinky to 2 tiles ahead of Pac-Man, doubled
		// 2 ahead of Pac-Man: (12, 15)
		// Vector from Blinky: (12-8, 15-13) = (4, 2)
		// Doubled and added to Blinky: (8+8, 13+4) = (16, 17)
		assert.Equal(t, types.Position{X: 16, Y: 17}, target)
	})
}

func TestClydeTarget(t *testing.T) {
	t.Run("targets Pac-Man when far away", func(t *testing.T) {
		pacmanPos := types.Position{X: 10, Y: 15}
		clydePos := types.Position{X: 20, Y: 25} // More than 8 tiles away
		
		target := ClydeTarget(pacmanPos, clydePos)
		
		assert.Equal(t, pacmanPos, target)
	})

	t.Run("targets scatter corner when close to Pac-Man", func(t *testing.T) {
		pacmanPos := types.Position{X: 10, Y: 15}
		clydePos := types.Position{X: 12, Y: 17} // Within 8 tiles
		scatterTarget := types.Position{X: 0, Y: 30} // Clyde's scatter corner
		
		target := ClydeTargetWithScatter(pacmanPos, clydePos, scatterTarget)
		
		assert.Equal(t, scatterTarget, target)
	})
}

func TestGhostScatterMode(t *testing.T) {
	t.Run("ghosts have scatter targets in corners", func(t *testing.T) {
		ghosts := CreateAllGhosts()
		
		// Blinky scatters to top-right
		assert.Equal(t, types.Position{X: 25, Y: 0}, ghosts[0].ScatterTarget)
		
		// Pinky scatters to top-left
		assert.Equal(t, types.Position{X: 2, Y: 0}, ghosts[1].ScatterTarget)
		
		// Inky scatters to bottom-right
		assert.Equal(t, types.Position{X: 27, Y: 30}, ghosts[2].ScatterTarget)
		
		// Clyde scatters to bottom-left
		assert.Equal(t, types.Position{X: 0, Y: 30}, ghosts[3].ScatterTarget)
	})
}

