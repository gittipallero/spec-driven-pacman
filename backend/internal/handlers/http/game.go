// Package http provides HTTP handlers for the game API.
package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/spec-driven-pacman/backend/internal/services/game"
)

// GameHandler handles game-related HTTP requests.
type GameHandler struct {
	service *game.Service
}

// NewGameHandler creates a new game handler.
func NewGameHandler(service *game.Service) *GameHandler {
	return &GameHandler{service: service}
}

// StartGame handles POST /api/game/start
func (h *GameHandler) StartGame(w http.ResponseWriter, r *http.Request) {
	session, err := h.service.CreateGame()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	wsURL := "/ws/game/" + session.ID
	response := session.ToStartResponse(wsURL)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetGameState handles GET /api/game/{sessionId}
func (h *GameHandler) GetGameState(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	
	session, err := h.service.GetSession(sessionID)
	if err != nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session.ToJSON())
}

// PauseGame handles POST /api/game/{sessionId}/pause
func (h *GameHandler) PauseGame(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	
	err := h.service.PauseGame(sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	session, _ := h.service.GetSession(sessionID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session.ToJSON())
}

// ResumeGame handles POST /api/game/{sessionId}/resume
func (h *GameHandler) ResumeGame(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	
	err := h.service.ResumeGame(sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	session, _ := h.service.GetSession(sessionID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session.ToJSON())
}

// EndGame handles DELETE /api/game/{sessionId}
func (h *GameHandler) EndGame(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	
	err := h.service.EndGame(sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeError(w http.ResponseWriter, code string, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Code:    code,
		Message: message,
	})
}

