// Package websocket provides WebSocket handlers for real-time game communication.
package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/spec-driven-pacman/backend/internal/services/game"
	"github.com/spec-driven-pacman/backend/pkg/types"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from localhost for development
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Handler manages WebSocket connections for game sessions.
type Handler struct {
	service     *game.Service
	connections map[string]map[*websocket.Conn]bool
	mu          sync.RWMutex
}

// NewHandler creates a new WebSocket handler.
func NewHandler(service *game.Service) *Handler {
	h := &Handler{
		service:     service,
		connections: make(map[string]map[*websocket.Conn]bool),
	}
	return h
}

// HandleConnection handles WebSocket connection upgrades.
func (h *Handler) HandleConnection(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")

	// Verify session exists
	_, err := h.service.GetSession(sessionID)
	if err != nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	// Upgrade to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	// Register connection
	h.addConnection(sessionID, conn)
	defer h.removeConnection(sessionID, conn)

	log.Printf("WebSocket connected for session %s", sessionID)

	// Start sending game state updates
	done := make(chan struct{})
	go h.sendStateUpdates(sessionID, conn, done)

	// Handle incoming messages
	h.handleMessages(sessionID, conn)
	close(done)
}

// addConnection registers a connection for a session.
func (h *Handler) addConnection(sessionID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.connections[sessionID] == nil {
		h.connections[sessionID] = make(map[*websocket.Conn]bool)
	}
	h.connections[sessionID][conn] = true
}

// removeConnection removes a connection from a session.
func (h *Handler) removeConnection(sessionID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.connections[sessionID] != nil {
		delete(h.connections[sessionID], conn)
		if len(h.connections[sessionID]) == 0 {
			delete(h.connections, sessionID)
		}
	}
	conn.Close()
}

// ClientMessage represents a message from the client.
type ClientMessage struct {
	Type      string          `json:"type"`
	Direction types.Direction `json:"direction,omitempty"`
	Timestamp int64           `json:"timestamp,omitempty"`
}

// ServerMessage represents a message to the client.
type ServerMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

// handleMessages processes incoming WebSocket messages.
func (h *Handler) handleMessages(sessionID string, conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var msg ClientMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			h.sendError(conn, "INVALID_MESSAGE", "Invalid message format")
			continue
		}

		switch msg.Type {
		case "input":
			if !msg.Direction.IsValid() {
				h.sendError(conn, "INVALID_DIRECTION", "Direction must be one of: up, down, left, right")
				continue
			}
			h.service.HandleInput(sessionID, msg.Direction)

		case "pause":
			h.service.PauseGame(sessionID)

		case "resume":
			h.service.ResumeGame(sessionID)

		default:
			h.sendError(conn, "INVALID_MESSAGE", "Unknown message type")
		}
	}
}

// sendStateUpdates sends periodic game state updates to the client.
func (h *Handler) sendStateUpdates(sessionID string, conn *websocket.Conn, done chan struct{}) {
	ticker := time.NewTicker(50 * time.Millisecond) // 20 TPS
	defer ticker.Stop()

	tick := 0
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			session, err := h.service.GetSession(sessionID)
			if err != nil {
				return
			}

			// Process game tick if playing
			if session.GetStatus() == types.GameStatusPlaying {
				session.Tick(50) // 50ms per tick
			}

			tick++
			state := session.ToJSON()
			
			msg := struct {
				Type               string      `json:"type"`
				Tick               int         `json:"tick"`
				Status             string      `json:"status"`
				Score              int         `json:"score"`
				Lives              int         `json:"lives"`
				Level              int         `json:"level"`
				PacMan             interface{} `json:"pacman"`
				Ghosts             interface{} `json:"ghosts"`
				DotsRemaining      int         `json:"dotsRemaining"`
				VulnerabilityTimer int         `json:"vulnerabilityTimer"`
			}{
				Type:               "state",
				Tick:               tick,
				Status:             string(state.Status),
				Score:              state.Score,
				Lives:              state.Lives,
				Level:              state.Level,
				PacMan:             state.PacMan,
				Ghosts:             state.Ghosts,
				DotsRemaining:      state.DotsRemaining,
				VulnerabilityTimer: state.VulnerabilityTimer,
			}

			if err := conn.WriteJSON(msg); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}
		}
	}
}

// sendError sends an error message to the client.
func (h *Handler) sendError(conn *websocket.Conn, code, message string) {
	msg := struct {
		Type    string `json:"type"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}{
		Type:    "error",
		Code:    code,
		Message: message,
	}
	conn.WriteJSON(msg)
}

// BroadcastToSession sends a message to all connections in a session.
func (h *Handler) BroadcastToSession(sessionID string, msg interface{}) {
	h.mu.RLock()
	conns := h.connections[sessionID]
	h.mu.RUnlock()

	for conn := range conns {
		if err := conn.WriteJSON(msg); err != nil {
			log.Printf("Broadcast error: %v", err)
		}
	}
}

