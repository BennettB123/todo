package main

import "fmt"

const (
	Debug = "DEBUG"
	Error = "ERROR"
)

type Logger struct {
	debug bool
}

func (logger Logger) Log(logLevel string, message string) {
	if logLevel == Debug && !logger.debug {
		return
	}

	fmt.Printf("[%s]: %s\n", logLevel, message)
}

func (logger Logger) LogError(message string, err error) {
	logger.Log(Error, message)
	logger.Log(Error, fmt.Sprintf("\t%s", err))
}
