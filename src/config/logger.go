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