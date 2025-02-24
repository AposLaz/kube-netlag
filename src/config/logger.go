package config

import (
	"fmt"
	"time"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
)

// Logger prints a log message with the given level and message to the console, along with the current timestamp.
// The message is formatted according to the fmt.Printf syntax, and the level is used to determine the color of the text.
// Allowed log levels are "INFO", "WARN", "ERROR", and "DEBUG", and the corresponding colors are green, yellow, red, and cyan, respectively.
func Logger(level, message string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02T15:04:05")

	// Apply color based on log level
	var color string
	switch level {
	case "INFO":
		color = Reset
	case "WARN":
		color = Yellow
	case "ERROR":
		color = Red
	case "DEBUG":
		color = Cyan
	default:
		color = Reset
	}

	logPrint := fmt.Sprintf("%s[%s] [%s] %s%s\n", color, timestamp, level, fmt.Sprintf(message, args...), Reset)
	fmt.Print(logPrint)
}
