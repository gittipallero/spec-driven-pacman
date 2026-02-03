// Package ghost provides ghost entity and AI behavior.
package ghost

import (
	"github.com/spec-driven-pacman/backend/pkg/types"
)

// GhostName identifies each of the four ghosts.
type GhostName string

const (
	GhostBlinky GhostName = "blinky" // Red - targets Pac-Man directly
	GhostPinky  GhostName = "pinky"  // Pink - targets ahead of Pac-Man
	GhostInky   GhostName = "inky"   // Cyan - complex targeting
	GhostClyde  GhostName = "clyde"  // Orange - shy, scatters when close
)

// GhostState represents the current behavior state of a ghost.
type GhostState string

const (
	GhostStateNormal     GhostState = "normal"
	GhostStateVulnerable GhostState = "vulnerable"
	GhostStateEaten      GhostState = "eaten"
	GhostStateRespawning GhostState = "respawning"
)

// GhostMode represents the AI behavior mode.
type GhostMode string

const (
	GhostModeChase      GhostMode = "chase"
	GhostModeScatter    GhostMode = "scatter"
	GhostModeFrightened GhostMode = "frightened"
)

// Ghost represents an enemy character.
type Ghost struct {
	Name         GhostName       `json:"name"`
	Position     types.Position  `json:"position"`
	Direction    types.Direction `json:"direction"`
	State        GhostState      `json:"state"`
	Mode         GhostMode       `json:"mode"`
	StartPos     types.Position  `json:"-"` // Starting position for respawn
	ScatterTarget types.Position `json:"-"` // Corner target for scatter mode
}

// NewGhost creates a new ghost with the given name and starting position.
func NewGhost(name GhostName, startPos types.Position, scatterTarget types.Position) *Ghost {
	return &Ghost{
		Name:          name,
		Position:      startPos,
		Direction:     types.DirectionUp,
		State:         GhostStateNormal,
		Mode:          GhostModeScatter,
		StartPos:      startPos,
		ScatterTarget: scatterTarget,
	}
}

// SetVulnerable changes the ghost to vulnerable state.
func (g *Ghost) SetVulnerable() {
	if g.State == GhostStateNormal {
		g.State = GhostStateVulnerable
		g.Mode = GhostModeFrightened
		// Reverse direction when becoming vulnerable
		g.Direction = g.Direction.Opposite()
	}
}

// SetEaten marks the ghost as eaten, triggering return to ghost house.
func (g *Ghost) SetEaten() {
	g.State = GhostStateEaten
}

// SetRespawning sets the ghost to respawning state (exiting ghost house).
func (g *Ghost) SetRespawning() {
	g.State = GhostStateRespawning
}

// SetNormal returns the ghost to normal chase/scatter behavior.
func (g *Ghost) SetNormal() {
	g.State = GhostStateNormal
	g.Mode = GhostModeChase
}

// Reset returns the ghost to its starting position and state.
func (g *Ghost) Reset() {
	g.Position = g.StartPos
	g.Direction = types.DirectionUp
	g.State = GhostStateNormal
	g.Mode = GhostModeScatter
}

// IsAtHome checks if the ghost is at its starting position.
func (g *Ghost) IsAtHome() bool {
	return g.Position.Equals(g.StartPos)
}

// ReverseDirection reverses the ghost's current direction.
func (g *Ghost) ReverseDirection() {
	g.Direction = g.Direction.Opposite()
}

// GhostJSON is the JSON-serializable representation of a ghost.
type GhostJSON struct {
	Name      GhostName       `json:"name"`
	X         int             `json:"x"`
	Y         int             `json:"y"`
	Direction types.Direction `json:"direction"`
	State     GhostState      `json:"state"`
	Mode      GhostMode       `json:"mode"`
}

// ToJSON converts Ghost to a JSON-friendly struct.
func (g *Ghost) ToJSON() GhostJSON {
	return GhostJSON{
		Name:      g.Name,
		X:         g.Position.X,
		Y:         g.Position.Y,
		Direction: g.Direction,
		State:     g.State,
		Mode:      g.Mode,
	}
}

// CreateAllGhosts creates the four standard Pac-Man ghosts at their starting positions.
func CreateAllGhosts() []*Ghost {
	return []*Ghost{
		NewGhost(GhostBlinky, types.Position{X: 13, Y: 11}, types.Position{X: 25, Y: 0}),   // Top-right scatter
		NewGhost(GhostPinky, types.Position{X: 13, Y: 14}, types.Position{X: 2, Y: 0}),    // Top-left scatter
		NewGhost(GhostInky, types.Position{X: 11, Y: 14}, types.Position{X: 27, Y: 30}),   // Bottom-right scatter
		NewGhost(GhostClyde, types.Position{X: 15, Y: 14}, types.Position{X: 0, Y: 30}),   // Bottom-left scatter
	}
}

