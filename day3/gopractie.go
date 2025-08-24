package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	users := []string{"Alice", "Bob", "Charlie"}
	var wg sync.WaitGroup

	for i := 0; i < len(users); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(200 * time.Millisecond) // pretend sending email takes time
			fmt.Println("Email sent to:", users[i]) // BUG ❌
		}()
	}

	wg.Wait()
}
AC