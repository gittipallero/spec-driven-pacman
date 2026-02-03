package game

import (
	"testing"

	"github.com/spec-driven-pacman/backend/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestNewPacMan(t *testing.T) {
	startPos := types.Position{X: 13, Y: 23}
	p := NewPacMan(startPos)

	assert.Equal(t, startPos, p.Position)
	assert.Equal(t, types.DirectionLeft, p.Direction)
	assert.Equal(t, types.DirectionNone, p.NextDirection)
	assert.Equal(t, PacManStateIdle, p.State)
}

func TestPacManMove(t *testing.T) {
	// Helper to create a walkability checker
	allWalkable := func(pos types.Position) bool {
		return true
	}
	
	blockedRight := func(pos types.Position) bool {
		return pos.X < 14 // Can't move right past x=13
	}

	t.Run("moves in current direction when valid", func(t *testing.T) {
		p := NewPacMan(types.Position{X: 13, Y: 23})
		p.Direction = types.DirectionRight
		
		moved := p.Move(allWalkable)
		
		assert.True(t, moved)
		assert.Equal(t, 14, p.Position.X)
		assert.Equal(t, 23, p.Position.Y)
		assert.Equal(t, PacManStateMoving, p.State)
	})

	t.Run("does not move when blocked", func(t *testing.T) {
		p := NewPacMan(types.Position{X: 13, Y: 23})
		p.Direction = types.DirectionRight
		
		moved := p.Move(blockedRight)
		
		assert.False(t, moved)
		assert.Equal(t, 13, p.Position.X)
		assert.Equal(t, PacManStateIdle, p.State)
	})

	t.Run("applies queued direction when valid", func(t *testing.T) {
		p := NewPacMan(types.Position{X: 13, Y: 23})
		p.Direction = types.DirectionRight
		p.SetNextDirection(types.DirectionUp)
		
		moved := p.Move(allWalkable)
		
		assert.True(t, moved)
		assert.Equal(t, types.DirectionUp, p.Direction)
		assert.Equal(t, types.DirectionNone, p.NextDirection)
		assert.Equal(t, 22, p.Position.Y) // Moved up
	})

	t.Run("ignores queued direction when invalid but continues current", func(t *testing.T) {
		blockedUp := func(pos types.Position) bool {
			return pos.Y >= 23 // Can't move up past y=23
		}
		
		p := NewPacMan(types.Position{X: 13, Y: 23})
		p.Direction = types.DirectionRight
		p.SetNextDirection(types.DirectionUp) // Invalid - blocked
		
		moved := p.Move(blockedUp)
		
		assert.True(t, moved)
		assert.Equal(t, types.DirectionRight, p.Direction) // Kept current
		assert.Equal(t, types.DirectionUp, p.NextDirection) // Still queued
		assert.Equal(t, 14, p.Position.X) // Moved right
	})
}

func TestPacManSetNextDirection(t *testing.T) {
	p := NewPacMan(types.Position{X: 13, Y: 23})
	
	p.SetNextDirection(types.DirectionUp)
	
	assert.Equal(t, types.DirectionUp, p.NextDirection)
}

func TestPacManCanMove(t *testing.T) {
	p := NewPacMan(types.Position{X: 13, Y: 23})
	
	t.Run("returns true when walkable", func(t *testing.T) {
		canMove := p.CanMove(types.DirectionRight, func(pos types.Position) bool {
			return true
		})
		assert.True(t, canMove)
	})

	t.Run("returns false when not walkable", func(t *testing.T) {
		canMove := p.CanMove(types.DirectionRight, func(pos types.Position) bool {
			return false
		})
		assert.False(t, canMove)
	})
}

func TestPacManSetDying(t *testing.T) {
	p := NewPacMan(types.Position{X: 13, Y: 23})
	
	p.SetDying()
	
	assert.Equal(t, PacManStateDying, p.State)
}

func TestPacManReset(t *testing.T) {
	p := NewPacMan(types.Position{X: 13, Y: 23})
	p.Position = types.Position{X: 5, Y: 5}
	p.Direction = types.DirectionUp
	p.NextDirection = types.DirectionDown
	p.State = PacManStateDying
	
	newStartPos := types.Position{X: 10, Y: 10}
	p.Reset(newStartPos)
	
	assert.Equal(t, newStartPos, p.Position)
	assert.Equal(t, types.DirectionLeft, p.Direction)
	assert.Equal(t, types.DirectionNone, p.NextDirection)
	assert.Equal(t, PacManStateIdle, p.State)
}

func TestPacManToJSON(t *testing.T) {
	p := NewPacMan(types.Position{X: 13, Y: 23})
	p.Direction = types.DirectionUp
	p.State = PacManStateMoving
	
	json := p.ToJSON()
	
	assert.Equal(t, 13, json.X)
	assert.Equal(t, 23, json.Y)
	assert.Equal(t, types.DirectionUp, json.Direction)
	assert.Equal(t, PacManStateMoving, json.State)
}

