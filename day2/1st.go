package main

import "fmt"

// makeCounter returns a function that increments and returns a counter value.
func makeCounter() func() int {
	count := 0 // variable captured by the closure

	return func() int {
		count++
		return count
	}
}

func main() {
	counter1 := makeCounter()
	fmt.Println(counter1()) // 1
	fmt.Println(counter1()) // 2
	fmt.Println(counter1()) // 3

	counter2 := makeCounter()
	fmt.Println(counter2()) // 1 (separate state from counter1)
}
