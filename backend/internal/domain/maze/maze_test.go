package maze

import (
	"testing"

	"github.com/spec-driven-pacman/backend/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadClassicMaze(t *testing.T) {
	t.Run("has correct dimensions", func(t *testing.T) {
		m := LoadClassicMaze()
		
		assert.Equal(t, MazeWidth, m.Width)
		assert.Equal(t, MazeHeight, m.Height)
		assert.Equal(t, 28, m.Width)
		assert.Equal(t, 31, m.Height)
	})

	t.Run("has correct number of rows and columns", func(t *testing.T) {
		m := LoadClassicMaze()
		
		require.Len(t, m.Tiles, MazeHeight)
		for y, row := range m.Tiles {
			assert.Len(t, row, MazeWidth, "row %d should have %d columns", y, MazeWidth)
		}
	})

	t.Run("has walls on edges", func(t *testing.T) {
		m := LoadClassicMaze()
		
		// Top edge (row 0) should be walls
		for x := 0; x < m.Width; x++ {
			assert.Equal(t, TileWall, m.GetTile(types.Position{X: x, Y: 0}), 
				"top edge at x=%d should be wall", x)
		}
		
		// Bottom edge (row 30) should be walls
		for x := 0; x < m.Width; x++ {
			assert.Equal(t, TileWall, m.GetTile(types.Position{X: x, Y: 30}), 
				"bottom edge at x=%d should be wall", x)
		}
	})

	t.Run("has ghost house in center", func(t *testing.T) {
		m := LoadClassicMaze()
		
		// Ghost house is around rows 12-16, columns 10-17
		ghostHousePos := types.Position{X: 13, Y: 14}
		assert.Equal(t, TileGhostHouse, m.GetTile(ghostHousePos), 
			"ghost house should be at center")
	})

	t.Run("has power pellets in corners", func(t *testing.T) {
		m := LoadClassicMaze()
		
		// Power pellets at known positions
		powerPelletPositions := []types.Position{
			{X: 1, Y: 3},   // Top left
			{X: 26, Y: 3},  // Top right
			{X: 1, Y: 23},  // Bottom left
			{X: 26, Y: 23}, // Bottom right
		}
		
		for _, pos := range powerPelletPositions {
			assert.Equal(t, TilePower, m.GetTile(pos), 
				"power pellet should be at (%d, %d)", pos.X, pos.Y)
		}
	})

	t.Run("counts dots correctly", func(t *testing.T) {
		m := LoadClassicMaze()
		
		assert.Greater(t, m.DotsRemaining, 0, "maze should have dots")
		assert.Equal(t, m.TotalDots, m.DotsRemaining, "all dots should remain initially")
	})

	t.Run("has tunnel tiles on sides", func(t *testing.T) {
		m := LoadClassicMaze()
		
		// Tunnel at row 14
		leftTunnel := types.Position{X: 0, Y: 14}
		rightTunnel := types.Position{X: 27, Y: 14}
		
		assert.Equal(t, TileTunnel, m.GetTile(leftTunnel), "left tunnel should exist")
		assert.Equal(t, TileTunnel, m.GetTile(rightTunnel), "right tunnel should exist")
	})
}

func TestMazeIsWalkable(t *testing.T) {
	m := LoadClassicMaze()

	t.Run("walls are not walkable", func(t *testing.T) {
		wallPos := types.Position{X: 0, Y: 0}
		assert.False(t, m.IsWalkable(wallPos), "wall should not be walkable")
	})

	t.Run("corridors are walkable", func(t *testing.T) {
		corridorPos := types.Position{X: 1, Y: 1}
		assert.True(t, m.IsWalkable(corridorPos), "corridor should be walkable")
	})

	t.Run("dots are walkable", func(t *testing.T) {
		dotPos := types.Position{X: 1, Y: 1} // Has a dot initially
		assert.True(t, m.IsWalkable(dotPos), "dot tile should be walkable")
	})

	t.Run("power pellets are walkable", func(t *testing.T) {
		powerPos := types.Position{X: 1, Y: 3}
		assert.True(t, m.IsWalkable(powerPos), "power pellet tile should be walkable")
	})

	t.Run("ghost house is not walkable for pacman", func(t *testing.T) {
		ghostHousePos := types.Position{X: 13, Y: 14}
		assert.False(t, m.IsWalkable(ghostHousePos), "ghost house should not be walkable for pacman")
	})

	t.Run("out of bounds is not walkable", func(t *testing.T) {
		outOfBounds := types.Position{X: 0, Y: -1}
		assert.False(t, m.IsWalkable(outOfBounds), "out of bounds should not be walkable")
	})
}

func TestMazeWrapPosition(t *testing.T) {
	m := LoadClassicMaze()

	t.Run("wraps left to right", func(t *testing.T) {
		pos := types.Position{X: -1, Y: 14}
		wrapped := m.WrapPosition(pos)
		
		assert.Equal(t, 27, wrapped.X, "should wrap to right side")
		assert.Equal(t, 14, wrapped.Y, "Y should not change")
	})

	t.Run("wraps right to left", func(t *testing.T) {
		pos := types.Position{X: 28, Y: 14}
		wrapped := m.WrapPosition(pos)
		
		assert.Equal(t, 0, wrapped.X, "should wrap to left side")
		assert.Equal(t, 14, wrapped.Y, "Y should not change")
	})

	t.Run("no wrap for normal positions", func(t *testing.T) {
		pos := types.Position{X: 14, Y: 14}
		wrapped := m.WrapPosition(pos)
		
		assert.Equal(t, pos, wrapped, "should not change normal position")
	})
}

func TestMazeCollectDot(t *testing.T) {
	t.Run("collecting dot returns points and decrements counter", func(t *testing.T) {
		m := LoadClassicMaze()
		initialDots := m.DotsRemaining
		
		dotPos := types.Position{X: 1, Y: 1} // Known dot position
		points := m.CollectDot(dotPos)
		
		assert.Equal(t, DotPoints, points, "should return dot points")
		assert.Equal(t, initialDots-1, m.DotsRemaining, "dots should decrement")
		assert.Equal(t, TileEmpty, m.GetTile(dotPos), "tile should become empty")
	})

	t.Run("collecting power pellet returns more points", func(t *testing.T) {
		m := LoadClassicMaze()
		initialDots := m.DotsRemaining
		
		powerPos := types.Position{X: 1, Y: 3} // Known power pellet position
		points := m.CollectDot(powerPos)
		
		assert.Equal(t, PowerPelletPoints, points, "should return power pellet points")
		assert.Equal(t, initialDots-1, m.DotsRemaining, "dots should decrement")
		assert.Equal(t, TileEmpty, m.GetTile(powerPos), "tile should become empty")
	})

	t.Run("collecting from empty tile returns zero", func(t *testing.T) {
		m := LoadClassicMaze()
		initialDots := m.DotsRemaining
		
		// Collect once
		dotPos := types.Position{X: 1, Y: 1}
		m.CollectDot(dotPos)
		
		// Collect again
		points := m.CollectDot(dotPos)
		
		assert.Equal(t, 0, points, "should return 0 for empty tile")
		assert.Equal(t, initialDots-1, m.DotsRemaining, "dots should not change")
	})

	t.Run("collecting from wall returns zero", func(t *testing.T) {
		m := LoadClassicMaze()
		
		wallPos := types.Position{X: 0, Y: 0}
		points := m.CollectDot(wallPos)
		
		assert.Equal(t, 0, points, "should return 0 for wall")
	})
}

