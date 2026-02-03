package game

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/spec-driven-pacman/backend/internal/domain/ghost"
	"github.com/spec-driven-pacman/backend/internal/domain/maze"
	"github.com/spec-driven-pacman/backend/pkg/types"
)

const (
	// InitialLives is the starting number of lives.
	InitialLives = 3
	// VulnerabilityDuration is how long ghosts stay vulnerable (ms).
	VulnerabilityDuration = 6000
	// GhostComboBase is the base score for eating a ghost.
	GhostComboBase = 200
)

// PacManStartPos is Pac-Man's starting position.
var PacManStartPos = types.Position{X: 13, Y: 23}

// GameSession represents a single game session.
type GameSession struct {
	ID                  string            `json:"id"`
	Status              types.GameStatus  `json:"status"`
	Score               int               `json:"score"`
	Lives               int               `json:"lives"`
	Level               int               `json:"level"`
	VulnerabilityTimer  int               `json:"vulnerabilityTimer"` // Remaining ms
	GhostsEatenCombo    int               `json:"-"`                  // Tracks combo multiplier
	
	PacMan              *PacMan           `json:"-"`
	Ghosts              []*ghost.Ghost    `json:"-"`
	Maze                *maze.Maze        `json:"-"`
	OriginalMaze        *maze.Maze        `json:"-"` // For resetting dots
	
	mu                  sync.RWMutex
	lastTick            time.Time
}

// NewGameSession creates a new game session with initial state.
func NewGameSession(gameMaze *maze.Maze) *GameSession {
	// Create a copy of the maze for resetting
	originalMaze := maze.NewMaze(gameMaze.Width, gameMaze.Height)
	for y := 0; y < gameMaze.Height; y++ {
		for x := 0; x < gameMaze.Width; x++ {
			originalMaze.Tiles[y][x] = gameMaze.Tiles[y][x]
		}
	}
	originalMaze.DotsRemaining = gameMaze.DotsRemaining
	originalMaze.TotalDots = gameMaze.TotalDots

	return &GameSession{
		ID:                 uuid.New().String(),
		Status:             types.GameStatusPlaying,
		Score:              0,
		Lives:              InitialLives,
		Level:              1,
		VulnerabilityTimer: 0,
		GhostsEatenCombo:   0,
		PacMan:             NewPacMan(PacManStartPos),
		Ghosts:             ghost.CreateAllGhosts(),
		Maze:               gameMaze,
		OriginalMaze:       originalMaze,
		lastTick:           time.Now(),
	}
}

// ProcessInput handles a direction input from the player.
func (s *GameSession) ProcessInput(dir types.Direction) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Status != types.GameStatusPlaying {
		return
	}

	if dir.IsValid() && dir != types.DirectionNone {
		s.PacMan.SetNextDirection(dir)
	}
}

// Pause pauses the game.
func (s *GameSession) Pause() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Status == types.GameStatusPlaying {
		s.Status = types.GameStatusPaused
	}
}

// Resume resumes a paused game.
func (s *GameSession) Resume() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Status == types.GameStatusPaused {
		s.Status = types.GameStatusPlaying
		s.lastTick = time.Now() // Reset tick timer to avoid jumps
	}
}

// TogglePause toggles between playing and paused states.
func (s *GameSession) TogglePause() {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch s.Status {
	case types.GameStatusPlaying:
		s.Status = types.GameStatusPaused
	case types.GameStatusPaused:
		s.Status = types.GameStatusPlaying
		s.lastTick = time.Now()
	}
}

// Tick processes one game tick.
// deltaMs is the time since last tick in milliseconds.
func (s *GameSession) Tick(deltaMs int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Status != types.GameStatusPlaying {
		return
	}

	// Move Pac-Man
	isWalkable := func(pos types.Position) bool {
		return s.Maze.IsWalkable(pos)
	}
	s.PacMan.Move(isWalkable)

	// Handle tunnel wrapping for Pac-Man
	s.PacMan.Position = s.Maze.WrapPosition(s.PacMan.Position)

	// Check dot collection
	s.checkDotCollection()

	// Update vulnerability timer
	if s.VulnerabilityTimer > 0 {
		s.VulnerabilityTimer -= deltaMs
		if s.VulnerabilityTimer <= 0 {
			s.VulnerabilityTimer = 0
			s.endVulnerability()
		}
	}

	// Check for level completion
	if s.Maze.DotsRemaining == 0 {
		s.completeLevel()
		return
	}

	// Check ghost collisions
	s.checkGhostCollisions()
}

// checkDotCollection handles dot and power pellet collection.
func (s *GameSession) checkDotCollection() {
	points := s.Maze.CollectDot(s.PacMan.Position)
	if points > 0 {
		s.Score += points
		
		// Check if power pellet
		if points == maze.PowerPelletPoints {
			s.startVulnerability()
		}
	}
}

// startVulnerability triggers ghost vulnerability mode.
func (s *GameSession) startVulnerability() {
	s.VulnerabilityTimer = VulnerabilityDuration
	s.GhostsEatenCombo = 0
	
	for _, g := range s.Ghosts {
		if g.State == ghost.GhostStateNormal {
			g.SetVulnerable()
		}
	}
}

// endVulnerability ends ghost vulnerability mode.
func (s *GameSession) endVulnerability() {
	for _, g := range s.Ghosts {
		if g.State == ghost.GhostStateVulnerable {
			g.SetNormal()
		}
	}
	s.GhostsEatenCombo = 0
}

// checkGhostCollisions handles collisions between Pac-Man and ghosts.
func (s *GameSession) checkGhostCollisions() {
	for _, g := range s.Ghosts {
		if !g.Position.Equals(s.PacMan.Position) {
			continue
		}

		switch g.State {
		case ghost.GhostStateNormal:
			s.handlePacManDeath()
			return
		case ghost.GhostStateVulnerable:
			s.eatGhost(g)
		}
	}
}

// handlePacManDeath handles Pac-Man being caught by a ghost.
func (s *GameSession) handlePacManDeath() {
	s.Lives--
	s.PacMan.SetDying()

	if s.Lives <= 0 {
		s.Status = types.GameStatusGameOver
	} else {
		// Reset positions after a short delay (handled by caller)
		s.resetPositions()
	}
}

// eatGhost handles Pac-Man eating a vulnerable ghost.
func (s *GameSession) eatGhost(g *ghost.Ghost) {
	g.SetEaten()
	
	// Calculate combo score: 200, 400, 800, 1600
	s.GhostsEatenCombo++
	comboMultiplier := 1 << (s.GhostsEatenCombo - 1) // 1, 2, 4, 8
	points := GhostComboBase * comboMultiplier
	s.Score += points
}

// resetPositions resets Pac-Man and ghost positions without clearing dots.
func (s *GameSession) resetPositions() {
	s.PacMan.Reset(PacManStartPos)
	for _, g := range s.Ghosts {
		g.Reset()
	}
	s.VulnerabilityTimer = 0
	s.GhostsEatenCombo = 0
}

// completeLevel handles completing a level.
func (s *GameSession) completeLevel() {
	s.Status = types.GameStatusLevelComplete
	s.Level++
	
	// Reset the maze with all dots
	s.Maze.ResetDots(s.OriginalMaze)
	s.resetPositions()
}

// GetStatus returns the current game status thread-safely.
func (s *GameSession) GetStatus() types.GameStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Status
}

// AddScore adds points to the score.
func (s *GameSession) AddScore(points int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Score += points
}

// GameStateJSON is the JSON-serializable game state.
type GameStateJSON struct {
	Status             types.GameStatus `json:"status"`
	Score              int              `json:"score"`
	Lives              int              `json:"lives"`
	Level              int              `json:"level"`
	PacMan             PacManJSON       `json:"pacman"`
	Ghosts             []ghost.GhostJSON `json:"ghosts"`
	DotsRemaining      int              `json:"dotsRemaining"`
	VulnerabilityTimer int              `json:"vulnerabilityTimer"`
}

// ToJSON converts the session to a JSON-friendly struct.
func (s *GameSession) ToJSON() GameStateJSON {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ghosts := make([]ghost.GhostJSON, len(s.Ghosts))
	for i, g := range s.Ghosts {
		ghosts[i] = g.ToJSON()
	}

	return GameStateJSON{
		Status:             s.Status,
		Score:              s.Score,
		Lives:              s.Lives,
		Level:              s.Level,
		PacMan:             s.PacMan.ToJSON(),
		Ghosts:             ghosts,
		DotsRemaining:      s.Maze.DotsRemaining,
		VulnerabilityTimer: s.VulnerabilityTimer,
	}
}

// GameStartResponseJSON is the response for starting a new game.
type GameStartResponseJSON struct {
	SessionID    string          `json:"sessionId"`
	WebsocketURL string          `json:"websocketUrl"`
	State        GameStateJSON   `json:"state"`
	Maze         maze.MazeJSON   `json:"maze"`
}

// ToStartResponse creates a game start response.
func (s *GameSession) ToStartResponse(wsURL string) GameStartResponseJSON {
	return GameStartResponseJSON{
		SessionID:    s.ID,
		WebsocketURL: wsURL,
		State:        s.ToJSON(),
		Maze:         s.Maze.ToJSON(),
	}
}

