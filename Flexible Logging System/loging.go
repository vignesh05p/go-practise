// You have multiple logging strategies (console, file, remote). How do you design a flexible logger in Go using interfaces so that your backend code can switch logging strategies easily?
package main

import "fmt"

type logger interface {
	log(message string) error
}

// Console logger implementation

type consolelogger struct{}

func (c consolelogger) log(message string) error {
	fmt.Println("Console log:", message)
	return nil
}
