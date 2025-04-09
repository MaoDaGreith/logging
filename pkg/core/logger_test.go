package core

import (
	"errors"
	"testing"
)

// MockDriver is a mock implementation of the Driver interface for testing
type MockDriver struct {
	Logs        []*LogEntry
	ShouldError bool
	Closed      bool
}

// Log records the log entry and returns an error if configured to do so
func (d *MockDriver) Log(entry *LogEntry) error {
	d.Logs = append(d.Logs, entry)
	if d.ShouldError {
		return errors.New("mock driver error")
	}
	return nil
}

// Close marks the driver as closed
func (d *MockDriver) Close() error {
	d.Closed = true
	if d.ShouldError {
		return errors.New("mock driver close error")
	}
	return nil
}

func TestLoggerLevels(t *testing.T) {
	mockDriver := &MockDriver{}
	logger := NewLogger(mockDriver)

	// Test each log level
	tests := []struct {
		name     string
		logFunc  func(msg string, attrs ...Attributes) error
		expected Level
	}{
		{"Debug", logger.Debug, Debug},
		{"Info", logger.Info, Info},
		{"Warning", logger.Warning, Warning},
		{"Error", logger.Error, Error},
	}

	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			msg := "test message"
			err := test.logFunc(msg)

			// Check no error
			if err != nil {
				t.Errorf("logger.%s() error = %v", test.name, err)
			}

			// Check log entry was created correctly
			if len(mockDriver.Logs) != i+1 {
				t.Fatalf("Expected %d logs, got %d", i+1, len(mockDriver.Logs))
			}

			log := mockDriver.Logs[i]
			if log.Level != test.expected {
				t.Errorf("Log level = %v, want %v", log.Level, test.expected)
			}
			if log.Message != msg {
				t.Errorf("Log message = %q, want %q", log.Message, msg)
			}
		})
	}
}

func TestLoggerWithAttributes(t *testing.T) {
	mockDriver := &MockDriver{}
	logger := NewLogger(mockDriver)

	attrs := Attributes{
		"key1": "value1",
		"key2": "value2",
	}

	err := logger.Info("test with attributes", attrs)

	// Check no error
	if err != nil {
		t.Errorf("logger.Info() error = %v", err)
	}

	// Check log entry was created correctly
	if len(mockDriver.Logs) != 1 {
		t.Fatalf("Expected 1 log, got %d", len(mockDriver.Logs))
	}

	log := mockDriver.Logs[0]
	if log.Level != Info {
		t.Errorf("Log level = %v, want %v", log.Level, Info)
	}
	if log.Message != "test with attributes" {
		t.Errorf("Log message = %q, want %q", log.Message, "test with attributes")
	}

	// Check attributes
	for k, v := range attrs {
		if log.Attrs[k] != v {
			t.Errorf("Log attribute %q = %q, want %q", k, log.Attrs[k], v)
		}
	}
}

func TestLoggerDriverError(t *testing.T) {
	mockDriver := &MockDriver{ShouldError: true}
	logger := NewLogger(mockDriver)

	err := logger.Info("test error")

	// Check error
	if err == nil {
		t.Errorf("Expected error but got nil")
	}
}

func TestLoggerMultipleDrivers(t *testing.T) {
	mockDriver1 := &MockDriver{}
	mockDriver2 := &MockDriver{}
	logger := NewLogger(mockDriver1, mockDriver2)

	msg := "test multiple drivers"
	err := logger.Info(msg)

	// Check no error
	if err != nil {
		t.Errorf("logger.Info() error = %v", err)
	}

	// Check both drivers received the log
	for i, driver := range []*MockDriver{mockDriver1, mockDriver2} {
		if len(driver.Logs) != 1 {
			t.Fatalf("Driver %d: Expected 1 log, got %d", i+1, len(driver.Logs))
		}

		log := driver.Logs[0]
		if log.Level != Info {
			t.Errorf("Driver %d: Log level = %v, want %v", i+1, log.Level, Info)
		}
		if log.Message != msg {
			t.Errorf("Driver %d: Log message = %q, want %q", i+1, log.Message, msg)
		}
	}
}

func TestLoggerClose(t *testing.T) {
	mockDriver1 := &MockDriver{}
	mockDriver2 := &MockDriver{ShouldError: true}
	logger := NewLogger(mockDriver1, mockDriver2)

	err := logger.Close()

	// Check error reported from second driver
	if err == nil {
		t.Errorf("Expected error but got nil")
	}

	// Check both drivers were closed
	if !mockDriver1.Closed {
		t.Errorf("First driver not closed")
	}
	if !mockDriver2.Closed {
		t.Errorf("Second driver not closed")
	}
}
