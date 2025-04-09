package drivers

import (
	"bytes"
	"regexp"
	"testing"
	"time"

	"github.com/MaoDaGreith/logging/pkg/core"
)

func TestConsoleDriver(t *testing.T) {
	// Create a buffer to capture output
	var stdout, stderr bytes.Buffer

	// Create a console driver with custom writers
	driver := NewConsoleDriverWithOptions(
		WithStdout(&stdout),
		WithStderr(&stderr),
		WithMinLevel(core.Debug),
		WithTimeFormat(time.RFC3339),
		WithColorized(false),
	)

	// Test cases
	tests := []struct {
		name      string
		level     core.Level
		message   string
		attrs     core.Attributes
		txID      string
		expectOut bool
	}{
		{
			name:      "debug level",
			level:     core.Debug,
			message:   "debug message",
			expectOut: true,
		},
		{
			name:      "info level",
			level:     core.Info,
			message:   "info message",
			expectOut: true,
		},
		{
			name:      "warning level",
			level:     core.Warning,
			message:   "warning message",
			expectOut: true,
		},
		{
			name:      "error level",
			level:     core.Error,
			message:   "error message",
			expectOut: true,
		},
		{
			name:      "with attributes",
			level:     core.Info,
			message:   "message with attrs",
			attrs:     core.Attributes{"key": "value"},
			expectOut: true,
		},
		{
			name:      "with transaction",
			level:     core.Info,
			message:   "tx message",
			txID:      "tx-123",
			expectOut: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear buffers
			stdout.Reset()
			stderr.Reset()

			// Create log entry
			entry := &core.LogEntry{
				Timestamp:     time.Now(),
				Level:         tt.level,
				Message:       tt.message,
				Attrs:         tt.attrs,
				TransactionID: tt.txID,
			}

			// Log the entry
			err := driver.Log(entry)
			if err != nil {
				t.Errorf("Log() error = %v", err)
				return
			}

			// Check output
			var output string
			if tt.level >= core.Error {
				output = stderr.String()
			} else {
				output = stdout.String()
			}

			if tt.expectOut && output == "" {
				t.Error("Expected output but got none")
			} else if !tt.expectOut && output != "" {
				t.Error("Expected no output but got some")
			}

			// Verify output format
			if tt.expectOut {
				// Check timestamp format
				if !timeRegex.MatchString(output) {
					t.Error("Output does not contain valid timestamp")
				}

				// Check level
				if !levelRegex.MatchString(output) {
					t.Error("Output does not contain valid level")
				}

				// Check message
				if !messageRegex.MatchString(output) {
					t.Error("Output does not contain message")
				}

				// Check attributes if present
				if len(tt.attrs) > 0 {
					if !attrsRegex.MatchString(output) {
						t.Error("Output does not contain attributes")
					}
				}

				// Check transaction ID if present
				if tt.txID != "" {
					if !txIDRegex.MatchString(output) {
						t.Error("Output does not contain transaction ID")
					}
				}
			}
		})
	}
}

// Regular expressions for output validation
var (
	timeRegex    = regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}[-+]\d{2}:\d{2}`)
	levelRegex   = regexp.MustCompile(`\[(DEBUG|INFO|WARNING|ERROR)\]`)
	messageRegex = regexp.MustCompile(`message`)
	attrsRegex   = regexp.MustCompile(`\[key=value\]`)
	txIDRegex    = regexp.MustCompile(`\(tx: tx-123\)`)
)

func TestConsoleDriverMinLevel(t *testing.T) {
	var stdout bytes.Buffer
	driver := NewConsoleDriverWithOptions(
		WithStdout(&stdout),
		WithMinLevel(core.Info),
	)

	// Test debug level (should not output)
	entry := &core.LogEntry{
		Timestamp: time.Now(),
		Level:     core.Debug,
		Message:   "debug message",
	}

	err := driver.Log(entry)
	if err != nil {
		t.Errorf("Log() error = %v", err)
	}

	if stdout.String() != "" {
		t.Error("Expected no output for debug level with min level Info")
	}

	// Test info level (should output)
	entry.Level = core.Info
	err = driver.Log(entry)
	if err != nil {
		t.Errorf("Log() error = %v", err)
	}

	if stdout.String() == "" {
		t.Error("Expected output for info level with min level Info")
	}
}

func TestConsoleDriverClose(t *testing.T) {
	driver := NewConsoleDriverWithOptions()
	err := driver.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}
}
