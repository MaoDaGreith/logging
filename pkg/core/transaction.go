package core

import (
	"time"
)

// Transaction represents a group of related log entries
type Transaction interface {
	// Debug logs a message at Debug level
	Debug(msg string, attrs ...Attributes) error

	// Info logs a message at Info level
	Info(msg string, attrs ...Attributes) error

	// Warning logs a message at Warning level
	Warning(msg string, attrs ...Attributes) error

	// Error logs a message at Error level
	Error(msg string, attrs ...Attributes) error

	// Log logs a message at the specified level
	Log(level Level, msg string, attrs ...Attributes) error

	// ID returns the transaction ID
	ID() string
}

// transaction implements the Transaction interface
type transaction struct {
	id     string
	logger *logger // Use concrete type to avoid circular dependency issues
}

// newTransaction creates a new transaction with the specified ID and logger
func newTransaction(id string, logger *logger) Transaction {
	return &transaction{
		id:     id,
		logger: logger,
	}
}

// Debug logs a message at Debug level
func (t *transaction) Debug(msg string, attrs ...Attributes) error {
	return t.Log(Debug, msg, attrs...)
}

// Info logs a message at Info level
func (t *transaction) Info(msg string, attrs ...Attributes) error {
	return t.Log(Info, msg, attrs...)
}

// Warning logs a message at Warning level
func (t *transaction) Warning(msg string, attrs ...Attributes) error {
	return t.Log(Warning, msg, attrs...)
}

// Error logs a message at Error level
func (t *transaction) Error(msg string, attrs ...Attributes) error {
	return t.Log(Error, msg, attrs...)
}

// Log logs a message at the specified level
func (t *transaction) Log(level Level, msg string, attrs ...Attributes) error {
	entry := &LogEntry{
		Timestamp:     time.Now(),
		Level:         level,
		Message:       msg,
		TransactionID: t.id,
	}

	if len(attrs) > 0 {
		entry.Attrs = attrs[0]
	}

	var lastErr error
	for _, driver := range t.logger.drivers {
		if err := driver.Log(entry); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

// ID returns the transaction ID
func (t *transaction) ID() string {
	return t.id
}
