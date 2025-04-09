package drivers

import (
	"testing"

	"github.com/MaoDaGreith/logging/pkg/core"
)

func TestDriverRegistry(t *testing.T) {
	// Save original registry
	originalRegistry := make(map[string]DriverConstructor)
	for k, v := range registry {
		originalRegistry[k] = v
	}

	// Clear registry for testing
	registry = make(map[string]DriverConstructor)

	// Test registration
	testDriver := func(options map[string]interface{}) (core.Driver, error) {
		return nil, nil
	}

	Register("test", testDriver)

	// Test driver creation
	driver, err := Create("test", nil)
	if err != nil {
		t.Errorf("Create() error = %v", err)
	}
	if driver != nil {
		t.Error("Expected nil driver from test constructor")
	}

	// Test non-existent driver
	driver, err = Create("non-existent", nil)
	if err != core.ErrDriverNotFound {
		t.Errorf("Create() error = %v, want %v", err, core.ErrDriverNotFound)
	}
	if driver != nil {
		t.Error("Expected nil driver for non-existent type")
	}

	// Restore original registry
	registry = originalRegistry
}

func TestConsoleDriverRegistration(t *testing.T) {
	driver, err := Create("console", map[string]interface{}{
		"min_level": "info",
	})
	if err != nil {
		t.Errorf("Create() error = %v", err)
	}
	if driver == nil {
		t.Error("Expected non-nil console driver")
	}
}

func TestJSONFileDriverRegistration(t *testing.T) {
	tempDir := t.TempDir()
	driver, err := Create("json_file", map[string]interface{}{
		"file_path": tempDir + "/test.json",
		"min_level": "info",
	})
	if err != nil {
		t.Errorf("Create() error = %v", err)
	}
	if driver == nil {
		t.Error("Expected non-nil JSON file driver")
	}
}

func TestTextFileDriverRegistration(t *testing.T) {
	tempDir := t.TempDir()
	driver, err := Create("text_file", map[string]interface{}{
		"file_path": tempDir + "/test.log",
		"min_level": "info",
	})
	if err != nil {
		t.Errorf("Create() error = %v", err)
	}
	if driver == nil {
		t.Error("Expected non-nil text file driver")
	}
}

func TestDriverCreationWithInvalidOptions(t *testing.T) {
	tests := []struct {
		name        string
		driverType  string
		options     map[string]interface{}
		expectError bool
	}{
		{
			name:       "console with invalid min level",
			driverType: "console",
			options: map[string]interface{}{
				"min_level": "invalid",
			},
			expectError: false, // Console driver ignores invalid min level
		},
		{
			name:        "json file without file path",
			driverType:  "json_file",
			options:     map[string]interface{}{},
			expectError: true,
		},
		{
			name:        "text file without file path",
			driverType:  "text_file",
			options:     map[string]interface{}{},
			expectError: true,
		},
		{
			name:       "json file with invalid min level",
			driverType: "json_file",
			options: map[string]interface{}{
				"file_path": "test.json",
				"min_level": "invalid",
			},
			expectError: false, // File drivers ignore invalid min level
		},
		{
			name:       "text file with invalid time format",
			driverType: "text_file",
			options: map[string]interface{}{
				"file_path":   "test.log",
				"time_format": 123, // Invalid type
			},
			expectError: false, // Text driver ignores invalid time format
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			driver, err := Create(tt.driverType, tt.options)
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				if driver != nil {
					t.Error("Expected nil driver when error occurs")
				}
			} else {
				if err != nil {
					t.Errorf("Create() error = %v", err)
				}
				if driver == nil {
					t.Error("Expected non-nil driver")
				}
			}
		})
	}
}
