package main

import (
	"fmt"
	"sync"
)

func handleMessage(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Handling message %d\n", id)
}

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 10000; i++ {
		wg.Add(1)
		go handleMessage(i, &wg) // each message handled in its own goroutine
	}

	wg.Wait()
	fmt.Println("All messages handled")
}
