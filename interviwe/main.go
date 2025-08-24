package main

import (
	"fmt"
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Request Method: %s, URL: %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello! You reached the main handler.")
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", LoggingMiddleware(http.HandlerFunc(mainHandler)))

	fmt.Println("Server running at :8080")
	http.ListenAndServe(":8080", mux)
}
