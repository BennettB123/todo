package main

import "fmt"

type Logger struct {
	debug bool
}

func (logger Logger) LogError(message string) {
	fmt.Printf("[Error]: %s\n", message)
}

func (logger Logger) LogDebug(message string) {
	if !logger.debug {
		return
	}

	fmt.Printf("[Debug]: %s\n", message)
}
