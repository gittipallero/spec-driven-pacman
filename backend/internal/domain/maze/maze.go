// Package maze provides the maze layout and tile management.
package maze

import (
	"github.com/spec-driven-pacman/backend/pkg/types"
)

// TileType represents the type of tile at a position.
type TileType int

const (
	TileEmpty      TileType = 0 // Walkable corridor
	TileWall       TileType = 1 // Impassable wall
	TileDot        TileType = 2 // Regular dot (10 points)
	TilePower      TileType = 3 // Power pellet (50 points)
	TileGhostHouse TileType = 4 // Ghost spawn area
	TileTunnel     TileType = 5 // Wraps to opposite side
)

const (
	// MazeWidth is the standard Pac-Man maze width.
	MazeWidth = 28
	// MazeHeight is the standard Pac-Man maze height.
	MazeHeight = 31
	// DotPoints is the score for collecting a regular dot.
	DotPoints = 10
	// PowerPelletPoints is the score for collecting a power pellet.
	PowerPelletPoints = 50
)

// Maze represents the game board.
type Maze struct {
	Width         int        `json:"width"`
	Height        int        `json:"height"`
	Tiles         [][]TileType `json:"tiles"`
	DotsRemaining int        `json:"dotsRemaining"`
	TotalDots     int        `json:"-"`
}

// NewMaze creates a new empty maze with the given dimensions.
func NewMaze(width, height int) *Maze {
	tiles := make([][]TileType, height)
	for y := range tiles {
		tiles[y] = make([]TileType, width)
	}
	return &Maze{
		Width:  width,
		Height: height,
		Tiles:  tiles,
	}
}

// GetTile returns the tile type at the given position.
func (m *Maze) GetTile(pos types.Position) TileType {
	// Handle tunnel wrapping
	x := pos.X
	if x < 0 {
		x = m.Width - 1
	} else if x >= m.Width {
		x = 0
	}

	if pos.Y < 0 || pos.Y >= m.Height {
		return TileWall
	}

	return m.Tiles[pos.Y][x]
}

// SetTile sets the tile type at the given position.
func (m *Maze) SetTile(pos types.Position, tile TileType) {
	if pos.X >= 0 && pos.X < m.Width && pos.Y >= 0 && pos.Y < m.Height {
		m.Tiles[pos.Y][pos.X] = tile
	}
}

// IsWalkable checks if the position can be walked on.
func (m *Maze) IsWalkable(pos types.Position) bool {
	tile := m.GetTile(pos)
	return tile != TileWall && tile != TileGhostHouse
}

// IsWalkableForGhost checks if a ghost can walk on the position.
func (m *Maze) IsWalkableForGhost(pos types.Position) bool {
	tile := m.GetTile(pos)
	return tile != TileWall
}

// IsDot checks if there is a dot at the position.
func (m *Maze) IsDot(pos types.Position) bool {
	return m.GetTile(pos) == TileDot
}

// IsPowerPellet checks if there is a power pellet at the position.
func (m *Maze) IsPowerPellet(pos types.Position) bool {
	return m.GetTile(pos) == TilePower
}

// IsTunnel checks if the position is a tunnel tile.
func (m *Maze) IsTunnel(pos types.Position) bool {
	return m.GetTile(pos) == TileTunnel
}

// CollectDot collects the dot at the given position.
// Returns the points earned (0 if no dot was there).
func (m *Maze) CollectDot(pos types.Position) int {
	tile := m.GetTile(pos)
	
	switch tile {
	case TileDot:
		m.SetTile(pos, TileEmpty)
		m.DotsRemaining--
		return DotPoints
	case TilePower:
		m.SetTile(pos, TileEmpty)
		m.DotsRemaining--
		return PowerPelletPoints
	default:
		return 0
	}
}

// CountDots counts all dots and power pellets in the maze.
func (m *Maze) CountDots() int {
	count := 0
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			tile := m.Tiles[y][x]
			if tile == TileDot || tile == TilePower {
				count++
			}
		}
	}
	return count
}

// ResetDots resets all dots in the maze from the original layout.
func (m *Maze) ResetDots(original *Maze) {
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			origTile := original.Tiles[y][x]
			if origTile == TileDot || origTile == TilePower {
				m.Tiles[y][x] = origTile
			}
		}
	}
	m.DotsRemaining = m.CountDots()
}

// WrapPosition wraps a position through the tunnel if needed.
func (m *Maze) WrapPosition(pos types.Position) types.Position {
	x := pos.X
	if x < 0 {
		x = m.Width - 1
	} else if x >= m.Width {
		x = 0
	}
	return types.Position{X: x, Y: pos.Y}
}

// MazeJSON is the JSON-serializable representation of a maze.
type MazeJSON struct {
	Width         int          `json:"width"`
	Height        int          `json:"height"`
	Tiles         [][]int      `json:"tiles"`
	DotsRemaining int          `json:"dotsRemaining"`
}

// ToJSON converts Maze to a JSON-friendly struct.
func (m *Maze) ToJSON() MazeJSON {
	tiles := make([][]int, m.Height)
	for y := range tiles {
		tiles[y] = make([]int, m.Width)
		for x := range tiles[y] {
			tiles[y][x] = int(m.Tiles[y][x])
		}
	}
	return MazeJSON{
		Width:         m.Width,
		Height:        m.Height,
		Tiles:         tiles,
		DotsRemaining: m.DotsRemaining,
	}
}

