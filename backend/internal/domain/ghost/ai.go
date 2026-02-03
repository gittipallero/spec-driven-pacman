package ghost

import (
	"math"

	"github.com/spec-driven-pacman/backend/pkg/types"
)

// BlinkyTarget returns Blinky's target: directly where Pac-Man is.
func BlinkyTarget(pacmanPos types.Position) types.Position {
	return pacmanPos
}

// PinkyTarget returns Pinky's target: 4 tiles ahead of Pac-Man.
// In the original game, when Pac-Man faces up, the target is also offset 4 tiles left.
func PinkyTarget(pacmanPos types.Position, pacmanDir types.Direction) types.Position {
	target := pacmanPos
	
	switch pacmanDir {
	case types.DirectionUp:
		// Original overflow bug: also moves 4 tiles left
		target.Y -= 4
		target.X -= 4
	case types.DirectionDown:
		target.Y += 4
	case types.DirectionLeft:
		target.X -= 4
	case types.DirectionRight:
		target.X += 4
	}
	
	return target
}

// InkyTarget returns Inky's target: uses both Pac-Man and Blinky's positions.
// Target is the vector from Blinky to 2 tiles ahead of Pac-Man, doubled from Blinky.
func InkyTarget(pacmanPos types.Position, pacmanDir types.Direction, blinkyPos types.Position) types.Position {
	// Get position 2 tiles ahead of Pac-Man
	intermediate := pacmanPos
	switch pacmanDir {
	case types.DirectionUp:
		intermediate.Y -= 2
		intermediate.X -= 2 // Same overflow bug as Pinky
	case types.DirectionDown:
		intermediate.Y += 2
	case types.DirectionLeft:
		intermediate.X -= 2
	case types.DirectionRight:
		intermediate.X += 2
	}
	
	// Vector from Blinky to intermediate
	vecX := intermediate.X - blinkyPos.X
	vecY := intermediate.Y - blinkyPos.Y
	
	// Double the vector and add to Blinky's position
	return types.Position{
		X: blinkyPos.X + (vecX * 2),
		Y: blinkyPos.Y + (vecY * 2),
	}
}

// ClydeTarget returns Clyde's target: Pac-Man if far away, otherwise scatter corner.
func ClydeTarget(pacmanPos types.Position, clydePos types.Position) types.Position {
	distance := calculateDistance(pacmanPos, clydePos)
	
	if distance > 8 {
		return pacmanPos
	}
	
	// Default scatter corner for Clyde (bottom-left)
	return types.Position{X: 0, Y: 30}
}

// ClydeTargetWithScatter returns Clyde's target with a custom scatter corner.
func ClydeTargetWithScatter(pacmanPos types.Position, clydePos types.Position, scatterTarget types.Position) types.Position {
	distance := calculateDistance(pacmanPos, clydePos)
	
	if distance > 8 {
		return pacmanPos
	}
	
	return scatterTarget
}

// calculateDistance calculates the Euclidean distance between two positions.
func calculateDistance(a, b types.Position) float64 {
	dx := float64(a.X - b.X)
	dy := float64(a.Y - b.Y)
	return math.Sqrt(dx*dx + dy*dy)
}

// GetTarget returns the appropriate target for a ghost based on its name and mode.
func GetTarget(ghost *Ghost, pacmanPos types.Position, pacmanDir types.Direction, blinkyPos types.Position) types.Position {
	// In scatter mode, return to scatter corner
	if ghost.Mode == GhostModeScatter {
		return ghost.ScatterTarget
	}
	
	// In frightened mode, targets don't matter (random movement)
	if ghost.Mode == GhostModeFrightened {
		return ghost.Position // Will use random direction selection
	}
	
	// Chase mode - use ghost-specific targeting
	switch ghost.Name {
	case GhostBlinky:
		return BlinkyTarget(pacmanPos)
	case GhostPinky:
		return PinkyTarget(pacmanPos, pacmanDir)
	case GhostInky:
		return InkyTarget(pacmanPos, pacmanDir, blinkyPos)
	case GhostClyde:
		return ClydeTargetWithScatter(pacmanPos, ghost.Position, ghost.ScatterTarget)
	default:
		return pacmanPos
	}
}

// ChooseDirection selects the best direction for a ghost to reach its target.
// Ghosts cannot reverse direction except when changing modes.
func ChooseDirection(ghost *Ghost, target types.Position, isWalkable func(types.Position) bool) types.Direction {
	currentDir := ghost.Direction
	oppositeDir := currentDir.Opposite()
	pos := ghost.Position
	
	// Possible directions to check (not including reverse)
	directions := []types.Direction{
		types.DirectionUp,
		types.DirectionLeft,
		types.DirectionDown,
		types.DirectionRight,
	}
	
	var bestDir types.Direction = types.DirectionNone
	bestDistance := math.MaxFloat64
	
	for _, dir := range directions {
		// Cannot reverse
		if dir == oppositeDir {
			continue
		}
		
		// Check if walkable
		nextPos := pos.Move(dir)
		if !isWalkable(nextPos) {
			continue
		}
		
		// Calculate distance to target
		dist := calculateDistance(nextPos, target)
		
		// Prefer directions with shorter distance
		// Tie-breaker: up > left > down > right (per original game)
		if dist < bestDistance {
			bestDistance = dist
			bestDir = dir
		}
	}
	
	// If no valid direction found, allow reverse (dead end)
	if bestDir == types.DirectionNone {
		nextPos := pos.Move(oppositeDir)
		if isWalkable(nextPos) {
			return oppositeDir
		}
	}
	
	return bestDir
}

