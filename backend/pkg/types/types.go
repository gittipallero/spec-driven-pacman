// Package types provides shared types and enums for the Pac-Man game.
package types

// Direction represents the movement direction.
type Direction string

const (
	DirectionUp    Direction = "up"
	DirectionDown  Direction = "down"
	DirectionLeft  Direction = "left"
	DirectionRight Direction = "right"
	DirectionNone  Direction = "none"
)

// IsValid checks if the direction is a valid value.
func (d Direction) IsValid() bool {
	switch d {
	case DirectionUp, DirectionDown, DirectionLeft, DirectionRight, DirectionNone:
		return true
	default:
		return false
	}
}

// Opposite returns the opposite direction.
func (d Direction) Opposite() Direction {
	switch d {
	case DirectionUp:
		return DirectionDown
	case DirectionDown:
		return DirectionUp
	case DirectionLeft:
		return DirectionRight
	case DirectionRight:
		return DirectionLeft
	default:
		return DirectionNone
	}
}

// GameStatus represents the current state of a game session.
type GameStatus string

const (
	GameStatusWaiting       GameStatus = "waiting"
	GameStatusPlaying       GameStatus = "playing"
	GameStatusPaused        GameStatus = "paused"
	GameStatusGameOver      GameStatus = "gameOver"
	GameStatusLevelComplete GameStatus = "levelComplete"
)

// Position represents a coordinate on the maze grid.
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Equals checks if two positions are the same.
func (p Position) Equals(other Position) bool {
	return p.X == other.X && p.Y == other.Y
}

// Move returns a new position moved in the given direction.
func (p Position) Move(dir Direction) Position {
	switch dir {
	case DirectionUp:
		return Position{X: p.X, Y: p.Y - 1}
	case DirectionDown:
		return Position{X: p.X, Y: p.Y + 1}
	case DirectionLeft:
		return Position{X: p.X - 1, Y: p.Y}
	case DirectionRight:
		return Position{X: p.X + 1, Y: p.Y}
	default:
		return p
	}
}

