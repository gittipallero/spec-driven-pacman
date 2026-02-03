// Package http provides HTTP handlers and middleware for the game API.
package http

import (
	"log"
	"net/http"
	"time"
)

// Logger is a middleware that logs all HTTP requests with timing information.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the response writer to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)
		log.Printf("%s %s %d %v", r.Method, r.URL.Path, wrapped.statusCode, duration)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// RateLimiter returns a rate limiting middleware.
// maxRequests is the max requests per second per IP.
type rateLimiter struct {
	requests map[string][]time.Time
	maxPerSec int
}

// NewRateLimiter creates a new rate limiter middleware.
func NewRateLimiter(maxRequestsPerSecond int) func(http.Handler) http.Handler {
	rl := &rateLimiter{
		requests:  make(map[string][]time.Time),
		maxPerSec: maxRequestsPerSecond,
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr

			// Clean old requests
			now := time.Now()
			cutoff := now.Add(-time.Second)
			var recent []time.Time
			for _, t := range rl.requests[ip] {
				if t.After(cutoff) {
					recent = append(recent, t)
				}
			}
			rl.requests[ip] = recent

			// Check limit
			if len(rl.requests[ip]) >= rl.maxPerSec {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			// Add this request
			rl.requests[ip] = append(rl.requests[ip], now)

			next.ServeHTTP(w, r)
		})
	}
}

// ValidateSessionID is a middleware that validates session ID format.
func ValidateSessionID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Session ID validation is done in handlers
		next.ServeHTTP(w, r)
	})
}

// Recovery is a middleware that recovers from panics.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

