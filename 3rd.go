package main

import (
	"fmt"
	"slices"
)

func main() {

	var num1 = []int{1, 2, 3, 4, 5}
	var num2 = []int{1, 2, 3, 4, 5}

	fmt.Print(slices.Equal(num1, num2))

}
