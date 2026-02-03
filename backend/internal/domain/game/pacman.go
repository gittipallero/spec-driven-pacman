// Package game provides game domain logic for Pac-Man.
package game

import (
	"github.com/spec-driven-pacman/backend/pkg/types"
)

// PacManState represents the visual/animation state of Pac-Man.
type PacManState string

const (
	PacManStateIdle   PacManState = "idle"
	PacManStateMoving PacManState = "moving"
	PacManStateDying  PacManState = "dying"
)

// PacMan represents the player-controlled character.
type PacMan struct {
	Position      types.Position  `json:"position"`
	Direction     types.Direction `json:"direction"`
	NextDirection types.Direction `json:"nextDirection"`
	State         PacManState     `json:"state"`
}

// NewPacMan creates a new Pac-Man at the specified starting position.
func NewPacMan(startPos types.Position) *PacMan {
	return &PacMan{
		Position:      startPos,
		Direction:     types.DirectionLeft,
		NextDirection: types.DirectionNone,
		State:         PacManStateIdle,
	}
}

// SetNextDirection queues a direction change for the next valid intersection.
func (p *PacMan) SetNextDirection(dir types.Direction) {
	p.NextDirection = dir
}

// CanMove checks if Pac-Man can move in the given direction.
// The maze parameter provides wall collision checking.
func (p *PacMan) CanMove(dir types.Direction, isWalkable func(types.Position) bool) bool {
	newPos := p.Position.Move(dir)
	return isWalkable(newPos)
}

// Move updates Pac-Man's position if the move is valid.
// Returns true if the move was successful.
func (p *PacMan) Move(isWalkable func(types.Position) bool) bool {
	// Try to apply queued direction first
	if p.NextDirection != types.DirectionNone && p.CanMove(p.NextDirection, isWalkable) {
		p.Direction = p.NextDirection
		p.NextDirection = types.DirectionNone
	}

	// Move in current direction if possible
	if p.CanMove(p.Direction, isWalkable) {
		p.Position = p.Position.Move(p.Direction)
		p.State = PacManStateMoving
		return true
	}

	p.State = PacManStateIdle
	return false
}

// SetDying sets Pac-Man to the dying state.
func (p *PacMan) SetDying() {
	p.State = PacManStateDying
}

// Reset resets Pac-Man to the starting position and state.
func (p *PacMan) Reset(startPos types.Position) {
	p.Position = startPos
	p.Direction = types.DirectionLeft
	p.NextDirection = types.DirectionNone
	p.State = PacManStateIdle
}

// ToJSON returns a JSON-serializable representation.
type PacManJSON struct {
	X         int             `json:"x"`
	Y         int             `json:"y"`
	Direction types.Direction `json:"direction"`
	State     PacManState     `json:"state"`
}

// ToJSON converts PacMan to a JSON-friendly struct.
func (p *PacMan) ToJSON() PacManJSON {
	return PacManJSON{
		X:         p.Position.X,
		Y:         p.Position.Y,
		Direction: p.Direction,
		State:     p.State,
	}
}

