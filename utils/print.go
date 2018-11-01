package utils

import (
	"fmt"
	"log"
)

type Level int

const (
	NORMAL Level = iota
	INFO
	WARNING
	ERROR
)

// Print uses log to print a simple string.
func Print(message string, level Level) {
	switch level {
	case ERROR:
		printError(&message)
	case INFO:
		info(&message)
	default:
		log.Println(message)
	}
}

func printError(message *string) {
	log.Printf("\x1b[0;91m%s\x1b[0m", *message)
}

func info(message *string) {
	log.Println(fmt.Sprintf("\x1b[0;32m%s\x1b[0m", *message))
}
