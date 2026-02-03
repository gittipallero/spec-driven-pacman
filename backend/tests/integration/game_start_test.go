package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	httphandlers "github.com/spec-driven-pacman/backend/internal/handlers/http"
	"github.com/spec-driven-pacman/backend/internal/services/game"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRouter() *chi.Mux {
	gameService := game.NewService()
	gameHandler := httphandlers.NewGameHandler(gameService)

	r := chi.NewRouter()
	r.Route("/api/game", func(r chi.Router) {
		r.Post("/start", gameHandler.StartGame)
		r.Get("/{sessionId}", gameHandler.GetGameState)
		r.Post("/{sessionId}/pause", gameHandler.PauseGame)
		r.Post("/{sessionId}/resume", gameHandler.ResumeGame)
		r.Delete("/{sessionId}", gameHandler.EndGame)
	})

	return r
}

func TestStartGame(t *testing.T) {
	router := setupRouter()

	t.Run("returns 201 with valid game response", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/game/start", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Check sessionId exists and is valid
		sessionId, ok := response["sessionId"].(string)
		assert.True(t, ok, "sessionId should be a string")
		assert.NotEmpty(t, sessionId, "sessionId should not be empty")
		assert.Len(t, sessionId, 36, "sessionId should be a UUID")

		// Check websocketUrl exists
		wsUrl, ok := response["websocketUrl"].(string)
		assert.True(t, ok, "websocketUrl should be a string")
		assert.Contains(t, wsUrl, sessionId, "websocketUrl should contain sessionId")

		// Check state exists
		state, ok := response["state"].(map[string]interface{})
		require.True(t, ok, "state should be an object")
		assert.Equal(t, "playing", state["status"])
		assert.Equal(t, float64(0), state["score"])
		assert.Equal(t, float64(3), state["lives"])
		assert.Equal(t, float64(1), state["level"])

		// Check maze exists
		maze, ok := response["maze"].(map[string]interface{})
		require.True(t, ok, "maze should be an object")
		assert.Equal(t, float64(28), maze["width"])
		assert.Equal(t, float64(31), maze["height"])
		assert.NotNil(t, maze["tiles"])
	})
}

func TestGetGameState(t *testing.T) {
	router := setupRouter()

	t.Run("returns 404 for non-existent session", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/game/non-existent-id", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("returns game state for existing session", func(t *testing.T) {
		// First create a game
		startReq := httptest.NewRequest("POST", "/api/game/start", nil)
		startW := httptest.NewRecorder()
		router.ServeHTTP(startW, startReq)

		var startResponse map[string]interface{}
		json.Unmarshal(startW.Body.Bytes(), &startResponse)
		sessionId := startResponse["sessionId"].(string)

		// Then get its state
		req := httptest.NewRequest("GET", "/api/game/"+sessionId, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var state map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &state)
		require.NoError(t, err)

		assert.Equal(t, "playing", state["status"])
	})
}

func TestPauseResume(t *testing.T) {
	router := setupRouter()

	// Create a game first
	startReq := httptest.NewRequest("POST", "/api/game/start", nil)
	startW := httptest.NewRecorder()
	router.ServeHTTP(startW, startReq)

	var startResponse map[string]interface{}
	json.Unmarshal(startW.Body.Bytes(), &startResponse)
	sessionId := startResponse["sessionId"].(string)

	t.Run("pause changes status to paused", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/game/"+sessionId+"/pause", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var state map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &state)
		assert.Equal(t, "paused", state["status"])
	})

	t.Run("resume changes status back to playing", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/game/"+sessionId+"/resume", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var state map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &state)
		assert.Equal(t, "playing", state["status"])
	})
}

func TestEndGame(t *testing.T) {
	router := setupRouter()

	// Create a game first
	startReq := httptest.NewRequest("POST", "/api/game/start", nil)
	startW := httptest.NewRecorder()
	router.ServeHTTP(startW, startReq)

	var startResponse map[string]interface{}
	json.Unmarshal(startW.Body.Bytes(), &startResponse)
	sessionId := startResponse["sessionId"].(string)

	t.Run("delete removes the session", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/game/"+sessionId, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)

		// Verify it's gone
		getReq := httptest.NewRequest("GET", "/api/game/"+sessionId, nil)
		getW := httptest.NewRecorder()
		router.ServeHTTP(getW, getReq)
		assert.Equal(t, http.StatusNotFound, getW.Code)
	})
}

