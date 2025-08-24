// /http example?
package main

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello! This is an HTTP server.")
}

func main() {
	http.HandleFunc("/", helloHandler) // Register handler for "/"
	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil) // Start server
}
