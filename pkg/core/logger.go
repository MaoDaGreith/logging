package core

import (
	"errors"
	"time"
)

// Errors
var (
	ErrDriverNotFound = errors.New("driver not found")
	ErrInvalidLevel   = errors.New("invalid log level")
)

// Attributes represents additional metadata for log entries
type Attributes map[string]string

// LogEntry represents a single log record
type LogEntry struct {
	// Timestamp when the log entry was created
	Timestamp time.Time

	// Level indicates the severity of the log entry
	Level Level

	// Message is the log message
	Message string

	// Attrs contains additional metadata about the log entry
	Attrs Attributes

	// TransactionID is an optional identifier for grouping related logs
	TransactionID string
}

// Driver defines the interface for log drivers
// This is defined here to avoid circular imports
type Driver interface {
	// Log processes a log entry
	Log(entry *LogEntry) error

	// Close allows the driver to close any resources when the logger is done
	Close() error
}

// logger implements the Logger interface
type logger struct {
	drivers []Driver
}

// NewLogger creates a new logger with the specified drivers
func NewLogger(drivers ...Driver) *logger {
	return &logger{
		drivers: drivers,
	}
}

// Debug logs a message at Debug level
func (l *logger) Debug(msg string, attrs ...Attributes) error {
	return l.Log(Debug, msg, attrs...)
}

// Info logs a message at Info level
func (l *logger) Info(msg string, attrs ...Attributes) error {
	return l.Log(Info, msg, attrs...)
}

// Warning logs a message at Warning level
func (l *logger) Warning(msg string, attrs ...Attributes) error {
	return l.Log(Warning, msg, attrs...)
}

// Error logs a message at Error level
func (l *logger) Error(msg string, attrs ...Attributes) error {
	return l.Log(Error, msg, attrs...)
}

// Log logs a message at the specified level
func (l *logger) Log(level Level, msg string, attrs ...Attributes) error {
	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   msg,
	}

	if len(attrs) > 0 {
		entry.Attrs = attrs[0]
	}

	var lastErr error
	for _, driver := range l.drivers {
		if err := driver.Log(entry); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

// NewTransaction creates a new transaction with the specified ID
func (l *logger) NewTransaction(txID string) Transaction {
	return newTransaction(txID, l)
}

// Close closes all drivers
func (l *logger) Close() error {
	var lastErr error
	for _, driver := range l.drivers {
		if err := driver.Close(); err != nil {
			lastErr = err
		}
	}

	return lastErr
}
