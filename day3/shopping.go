package main

import (
	"fmt"
	"time"
)

// This function will now run concurrently for each request.
func handleBuyRequest(requestID string) {
	fmt.Printf("Request %s started.\n", requestID)

	// Create a channel to signal when all tasks are done.
	done := make(chan bool)

	// Launch each task as a separate goroutine.
	go processPayment(requestID, done)
	go updateInventory(requestID, done)
	go sendConfirmationEmail(requestID, done)
	go logTransaction(requestID, done)

	// Wait for all four tasks to signal they are done.
	for i := 0; i < 4; i++ {
		<-done
	}

	fmt.Printf("All tasks for request %s finished. Request complete.\n", requestID)
}

func processPayment(requestID string, done chan bool) {
	fmt.Printf("Processing payment for %s...\n", requestID)
	time.Sleep(2 * time.Second) // Simulate waiting for a bank API
	fmt.Printf("Payment for %s complete.\n", requestID)
	done <- true // Signal that this task is done.
}

func updateInventory(requestID string, done chan bool) {
	fmt.Printf("Updating inventory for %s...\n", requestID)
	time.Sleep(1 * time.Second) // Simulate waiting for a database
	fmt.Printf("Inventory for %s updated.\n", requestID)
	done <- true // Signal that this task is done.
}

func sendConfirmationEmail(requestID string, done chan bool) {
	fmt.Printf("Sending confirmation email for %s...\n", requestID)
	time.Sleep(3 * time.Second) // Simulate waiting for an email service
	fmt.Printf("Confirmation email for %s sent.\n", requestID)
	done <- true // Signal that this task is done.
}

func logTransaction(requestID string, done chan bool) {
	fmt.Printf("Logging transaction for %s...\n", requestID)
	time.Sleep(1 * time.Second) // Simulate writing to a log file
	fmt.Printf("Transaction for %s logged.\n", requestID)
	done <- true // Signal that this task is done.
}

func main() {
	handleBuyRequest("user-1")
	// If you run this with multiple requests, you'll see how they all interleave.
	// For example:
	// handleBuyRequest("user-2")
}
