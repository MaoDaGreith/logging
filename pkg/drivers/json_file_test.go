package drivers

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/MaoDaGreith/logging/pkg/core"
)

func TestJSONFileDriver(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "json-driver-test")
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
			filePath := filepath.Join(tempDir, tt.name+".json")

			// Create driver
			driver, err := NewJSONFileDriver(map[string]interface{}{
				"file_path": filePath,
				"min_level": "debug",
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

			// Decode the JSON
			var logEntry struct {
				Timestamp     string            `json:"timestamp"`
				Level         string            `json:"level"`
				Message       string            `json:"message"`
				Attributes    map[string]string `json:"attributes,omitempty"`
				TransactionID string            `json:"transaction_id,omitempty"`
			}

			if err := json.Unmarshal(data, &logEntry); err != nil {
				t.Errorf("Failed to decode JSON: %v", err)
				return
			}

			// Verify fields
			if logEntry.Level != tt.level.String() {
				t.Errorf("Level = %v, want %v", logEntry.Level, tt.level.String())
			}

			if logEntry.Message != tt.message {
				t.Errorf("Message = %v, want %v", logEntry.Message, tt.message)
			}

			if tt.txID != "" && logEntry.TransactionID != tt.txID {
				t.Errorf("TransactionID = %v, want %v", logEntry.TransactionID, tt.txID)
			}

			if len(tt.attrs) > 0 {
				if len(logEntry.Attributes) != len(tt.attrs) {
					t.Errorf("Attributes length = %v, want %v", len(logEntry.Attributes), len(tt.attrs))
				}
				for k, v := range tt.attrs {
					if logEntry.Attributes[k] != v {
						t.Errorf("Attribute[%v] = %v, want %v", k, logEntry.Attributes[k], v)
					}
				}
			}
		})
	}
}

func TestJSONFileDriverMinLevel(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "json-driver-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test file path
	filePath := filepath.Join(tempDir, "min_level_test.json")

	// Create driver with min level Info
	driver, err := NewJSONFileDriver(map[string]interface{}{
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
	driver, err = NewJSONFileDriver(map[string]interface{}{
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

	var logEntry struct {
		Level   string `json:"level"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(data, &logEntry); err != nil {
		t.Errorf("Failed to decode JSON: %v", err)
		return
	}

	if logEntry.Level != "INFO" {
		t.Errorf("Level = %v, want INFO", logEntry.Level)
	}
}

func TestJSONFileDriverClose(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "json-driver-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test file path
	filePath := filepath.Join(tempDir, "close_test.json")

	// Create driver
	driver, err := NewJSONFileDriver(map[string]interface{}{
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
