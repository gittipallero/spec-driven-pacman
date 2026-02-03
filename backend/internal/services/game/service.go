// Package game provides the game service for orchestrating game sessions.
package game

import (
	"errors"
	"sync"

	"github.com/spec-driven-pacman/backend/internal/domain/game"
	"github.com/spec-driven-pacman/backend/internal/domain/maze"
	"github.com/spec-driven-pacman/backend/pkg/types"
)

var (
	// ErrSessionNotFound is returned when a session doesn't exist.
	ErrSessionNotFound = errors.New("session not found")
	// ErrGameNotPlaying is returned when trying to pause a non-playing game.
	ErrGameNotPlaying = errors.New("game is not playing")
	// ErrGameNotPaused is returned when trying to resume a non-paused game.
	ErrGameNotPaused = errors.New("game is not paused")
)

// Service manages game sessions.
type Service struct {
	sessions map[string]*game.GameSession
	mu       sync.RWMutex
}

// NewService creates a new game service.
func NewService() *Service {
	return &Service{
		sessions: make(map[string]*game.GameSession),
	}
}

// CreateGame creates a new game session.
func (s *Service) CreateGame() (*game.GameSession, error) {
	// Load the classic maze
	gameMaze := maze.LoadClassicMaze()
	
	// Create new session
	session := game.NewGameSession(gameMaze)

	// Store session
	s.mu.Lock()
	s.sessions[session.ID] = session
	s.mu.Unlock()

	return session, nil
}

// GetSession retrieves a game session by ID.
func (s *Service) GetSession(id string) (*game.GameSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[id]
	if !exists {
		return nil, ErrSessionNotFound
	}
	return session, nil
}

// HandleInput processes a direction input for a session.
func (s *Service) HandleInput(sessionID string, dir types.Direction) error {
	session, err := s.GetSession(sessionID)
	if err != nil {
		return err
	}

	session.ProcessInput(dir)
	return nil
}

// PauseGame pauses a game session.
func (s *Service) PauseGame(sessionID string) error {
	session, err := s.GetSession(sessionID)
	if err != nil {
		return err
	}

	if session.GetStatus() != types.GameStatusPlaying {
		return ErrGameNotPlaying
	}

	session.Pause()
	return nil
}

// ResumeGame resumes a paused game session.
func (s *Service) ResumeGame(sessionID string) error {
	session, err := s.GetSession(sessionID)
	if err != nil {
		return err
	}

	if session.GetStatus() != types.GameStatusPaused {
		return ErrGameNotPaused
	}

	session.Resume()
	return nil
}

// TogglePause toggles the pause state of a game session.
func (s *Service) TogglePause(sessionID string) error {
	session, err := s.GetSession(sessionID)
	if err != nil {
		return err
	}

	session.TogglePause()
	return nil
}

// EndGame ends and removes a game session.
func (s *Service) EndGame(sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.sessions[sessionID]; !exists {
		return ErrSessionNotFound
	}

	delete(s.sessions, sessionID)
	return nil
}

// GetActiveSessions returns the number of active sessions.
func (s *Service) GetActiveSessions() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.sessions)
}

