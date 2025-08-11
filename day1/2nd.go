package main

import "fmt"

func printType(i interface{}) {
	switch v := i.(type) {
	case string:
		fmt.Println("String:", v)
	case int:
		fmt.Println("Integer:", v)
	case float64:
		fmt.Println("Float:", v)
	default:
		fmt.Println("Unknown type")
	}
}

func main() {
	printType("Hello")
	printType(42)
	printType(3.14)
	printType(true)
}
