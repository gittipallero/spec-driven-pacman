// Package main is the entry point for the Pac-Man game server.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	
	httphandlers "github.com/spec-driven-pacman/backend/internal/handlers/http"
	wshandlers "github.com/spec-driven-pacman/backend/internal/handlers/websocket"
	"github.com/spec-driven-pacman/backend/internal/services/game"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize services
	gameService := game.NewService()

	// Initialize handlers
	gameHandler := httphandlers.NewGameHandler(gameService)
	wsHandler := wshandlers.NewHandler(gameService)

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(httphandlers.NewRateLimiter(60)) // 60 requests per second
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Routes
	r.Route("/api/game", func(r chi.Router) {
		r.Post("/start", gameHandler.StartGame)
		r.Get("/{sessionId}", gameHandler.GetGameState)
		r.Post("/{sessionId}/pause", gameHandler.PauseGame)
		r.Post("/{sessionId}/resume", gameHandler.ResumeGame)
		r.Delete("/{sessionId}", gameHandler.EndGame)
	})

	// WebSocket route
	r.Get("/ws/game/{sessionId}", wsHandler.HandleConnection)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

