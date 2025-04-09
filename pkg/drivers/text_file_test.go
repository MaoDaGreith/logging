package drivers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/MaoDaGreith/logging/pkg/core"
)

func TestTextFileDriver(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "text-driver-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name    string
		level   core.Level
		message string
		attrs   core.Attributes
		txID    string
	}{
		{
			name:    "debug level",
			level:   core.Debug,
			message: "debug message",
		},
		{
			name:    "info level",
			level:   core.Info,
			message: "info message",
		},
		{
			name:    "warning level",
			level:   core.Warning,
			message: "warning message",
		},
		{
			name:    "error level",
			level:   core.Error,
			message: "error message",
		},
		{
			name:    "with attributes",
			level:   core.Info,
			message: "message with attrs",
			attrs:   core.Attributes{"key": "value"},
		},
		{
			name:    "with transaction",
			level:   core.Info,
			message: "tx message",
			txID:    "tx-123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new file for each test
			filePath := filepath.Join(tempDir, tt.name+".log")

			// Create driver
			driver, err := NewTextFileDriver(map[string]interface{}{
				"file_path":   filePath,
				"min_level":   "debug",
				"time_format": time.RFC3339,
			})
			if err != nil {
				t.Fatalf("Failed to create driver: %v", err)
			}
			defer driver.Close()

			// Create log entry
			entry := &core.LogEntry{
				Timestamp:     time.Now(),
				Level:         tt.level,
				Message:       tt.message,
				Attrs:         tt.attrs,
				TransactionID: tt.txID,
			}

			// Log the entry
			err = driver.Log(entry)
			if err != nil {
				t.Errorf("Log() error = %v", err)
				return
			}

			// Close the driver to ensure file is written
			driver.Close()

			// Read the file
			data, err := os.ReadFile(filePath)
			if err != nil {
				t.Errorf("Failed to read file: %v", err)
				return
			}

			output := string(data)

			// Verify output format
			if !strings.Contains(output, tt.level.String()) {
				t.Errorf("Output does not contain level %s", tt.level)
			}

			if !strings.Contains(output, tt.message) {
				t.Errorf("Output does not contain message %q", tt.message)
			}

			if len(tt.attrs) > 0 {
				for k, v := range tt.attrs {
					if !strings.Contains(output, k+"="+v) {
						t.Errorf("Output does not contain attribute %s=%s", k, v)
					}
				}
			}

			if tt.txID != "" && !strings.Contains(output, tt.txID) {
				t.Errorf("Output does not contain transaction ID %s", tt.txID)
			}
		})
	}
}

func TestTextFileDriverMinLevel(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "text-driver-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test file path
	filePath := filepath.Join(tempDir, "min_level_test.log")

	// Create driver with min level Info
	driver, err := NewTextFileDriver(map[string]interface{}{
		"file_path": filePath,
		"min_level": "info",
	})
	if err != nil {
		t.Fatalf("Failed to create driver: %v", err)
	}
	defer driver.Close()

	// Test debug level (should not write)
	entry := &core.LogEntry{
		Timestamp: time.Now(),
		Level:     core.Debug,
		Message:   "debug message",
	}

	err = driver.Log(entry)
	if err != nil {
		t.Errorf("Log() error = %v", err)
	}

	// Close and reopen to ensure file is written
	driver.Close()

	// Check file is empty
	data, err := os.ReadFile(filePath)
	if err != nil && !os.IsNotExist(err) {
		t.Errorf("Failed to read file: %v", err)
		return
	}

	if len(data) > 0 {
		t.Error("Expected file to be empty for debug level with min level Info")
	}

	// Create new driver for info level test
	driver, err = NewTextFileDriver(map[string]interface{}{
		"file_path": filePath,
		"min_level": "info",
	})
	if err != nil {
		t.Fatalf("Failed to create driver: %v", err)
	}
	defer driver.Close()

	// Test info level (should write)
	entry.Level = core.Info
	err = driver.Log(entry)
	if err != nil {
		t.Errorf("Log() error = %v", err)
	}

	// Close to ensure file is written
	driver.Close()

	// Check file has content
	data, err = os.ReadFile(filePath)
	if err != nil {
		t.Errorf("Failed to read file: %v", err)
		return
	}

	if len(data) == 0 {
		t.Error("Expected file to have content for info level with min level Info")
	}

	output := string(data)
	if !strings.Contains(output, "INFO") {
		t.Error("Output does not contain INFO level")
	}
}

func TestTextFileDriverClose(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "text-driver-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test file path
	filePath := filepath.Join(tempDir, "close_test.log")

	// Create driver
	driver, err := NewTextFileDriver(map[string]interface{}{
		"file_path": filePath,
	})
	if err != nil {
		t.Fatalf("Failed to create driver: %v", err)
	}

	// Close the driver
	err = driver.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}

	// Try to write after close
	entry := &core.LogEntry{
		Timestamp: time.Now(),
		Level:     core.Info,
		Message:   "test message",
	}

	err = driver.Log(entry)
	if err == nil {
		t.Error("Expected error when writing to closed driver")
	}
}
