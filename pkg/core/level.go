package core

import (
	"fmt"
	"strings"
)

// Level represents the severity level of a log message
type Level int

const (
	// Debug is used for development-time messages and debugging
	Debug Level = iota
	// Info is used for general information about system operation
	Info
	// Warning is used for non-critical issues that might need attention
	Warning
	// Error is used for error conditions
	Error
)

// String returns the string representation of the log level
func (l Level) String() string {
	switch l {
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warning:
		return "WARNING"
	case Error:
		return "ERROR"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", l)
	}
}

// ParseLevel converts a string to a Level
func ParseLevel(levelStr string) (Level, error) {
	switch strings.ToUpper(levelStr) {
	case "DEBUG":
		return Debug, nil
	case "INFO":
		return Info, nil
	case "WARNING", "WARN":
		return Warning, nil
	case "ERROR", "ERR":
		return Error, nil
	default:
		return Info, fmt.Errorf("unknown log level: %s", levelStr)
	}
}
