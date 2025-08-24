package main

import (
	"fmt"
	"net/http"
	"time"
)

// Logging middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		fmt.Println("Request started:", r.Method, r.URL.Path)

		next.ServeHTTP(w, r) // call the next handler

		fmt.Println("Request completed in:", time.Since(start))
	})
}

// Main handler
func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Go!"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler)

	// Wrap with middleware
	wrapped := loggingMiddleware(mux)

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", wrapped)
}
