package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

/*
  Middleware chain (outer -> inner):
  recoveryMiddleware(
      loggingMiddleware(
          rateLimitMiddleware(
              corsMiddleware(
                  authMiddleware(finalHandler)
              )
          )
      )
  )

  Recovery should wrap everything to catch panics from any layer.
  CORS is placed before auth so preflight (OPTIONS) requests don't require auth.
*/

// -------------------- Recovery Middleware --------------------
func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				// Log the panic and return 500
				log.Printf("[PANIC RECOVER] %v - %s %s\n", rec, r.Method, r.URL.Path)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// -------------------- Logging Middleware --------------------
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("[REQUEST START] %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("[REQUEST END] %s %s (%s)", r.Method, r.URL.Path, time.Since(start))
	})
}

// -------------------- Simple Auth Middleware --------------------
// Checks for a header X-API-KEY == "secret123" (change for production)
func authMiddleware(next http.Handler) http.Handler {
	const validKey = "secret123" // change in real app, do not hardcode
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for safe public endpoints if needed, e.g., health
		if r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}

		key := r.Header.Get("X-API-KEY")
		if key == "" {
			http.Error(w, "missing API key", http.StatusUnauthorized)
			return
		}
		if key != validKey {
			http.Error(w, "invalid API key", http.StatusForbidden)
			return
		}
		// Auth OK
		next.ServeHTTP(w, r)
	})
}

// -------------------- CORS Middleware --------------------
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // restrict in prod
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-API-KEY")

		// Handle preflight (OPTIONS) requests directly
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// -------------------- Rate Limiter --------------------
// Simple token-bucket per client IP
type clientBucket struct {
	tokens     float64
	lastRefill time.Time
	ratePerSec float64 // tokens added per second
	maxTokens  float64
}

type rateLimiter struct {
	clients map[string]*clientBucket
	mu      sync.Mutex
	// cleanup ticker to remove stale entries
	cleanupTicker *time.Ticker
}

func newRateLimiter() *rateLimiter {
	rl := &rateLimiter{
		clients:       make(map[string]*clientBucket),
		cleanupTicker: time.NewTicker(5 * time.Minute),
	}
	// Start cleanup goroutine
	go func() {
		for range rl.cleanupTicker.C {
			rl.cleanupStale(10 * time.Minute)
		}
	}()
	return rl
}

func (rl *rateLimiter) getBucket(key string) *clientBucket {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	cb, ok := rl.clients[key]
	if !ok {
		cb = &clientBucket{
			tokens:     5.0, // start full
			lastRefill: time.Now(),
			ratePerSec: 5.0 / 60.0, // 5 tokens per minute => tokens/second
			maxTokens:  5.0,
		}
		rl.clients[key] = cb
	}
	return cb
}

func (rl *rateLimiter) allow(key string) bool {
	cb := rl.getBucket(key)

	// refill tokens based on time passed
	now := time.Now()
	elapsed := now.Sub(cb.lastRefill).Seconds()
	if elapsed > 0 {
		cb.tokens += elapsed * cb.ratePerSec
		if cb.tokens > cb.maxTokens {
			cb.tokens = cb.maxTokens
		}
		cb.lastRefill = now
	}

	if cb.tokens >= 1.0 {
		cb.tokens -= 1.0
		return true
	}
	return false
}

func (rl *rateLimiter) cleanupStale(age time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()
	for k, v := range rl.clients {
		if now.Sub(v.lastRefill) > age {
			delete(rl.clients, k)
		}
	}
}

func (rl *rateLimiter) stop() {
	rl.cleanupTicker.Stop()
}

// Helper to extract client IP (checks X-Forwarded-For then RemoteAddr)
func clientIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// X-Forwarded-For can be comma separated list
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func rateLimitMiddleware(rl *rateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := clientIP(r)
			if !rl.allow(ip) {
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// -------------------- Example Handlers --------------------
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// simulate some work
	time.Sleep(100 * time.Millisecond)
	fmt.Fprintf(w, "Hello! You hit %s\n", r.URL.Path)
}

func panicHandler(w http.ResponseWriter, r *http.Request) {
	// Force a panic to demo recovery middleware
	panic("simulated panic for demo")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

// -------------------- Main & Middleware Chaining --------------------
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler)
	mux.HandleFunc("/panic", panicHandler)
	mux.HandleFunc("/health", healthHandler)

	// create rate limiter
	rl := newRateLimiter()
	defer rl.stop()

	// chain middlewares:
	// recovery -> logging -> rateLimit -> cors -> auth -> mux
	handler := recoveryMiddleware(
		loggingMiddleware(
			rateLimitMiddleware(rl)(
				corsMiddleware(
					authMiddleware(mux),
				),
			),
		),
	)

	addr := ":8080"
	log.Printf("Server starting on %s\n", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
