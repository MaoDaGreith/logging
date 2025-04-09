package drivers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/MaoDaGreith/logging/pkg/core"
)

// JSONFileDriverName is the name to use in configuration
const JSONFileDriverName = "json_file"

func init() {
	Register(JSONFileDriverName, NewJSONFileDriver)
}

// JSONFileDriver outputs logs to a JSON file
type JSONFileDriver struct {
	filePath string
	file     *os.File
	encoder  *json.Encoder
	minLevel core.Level
	mu       sync.Mutex
}

// JSONLogEntry represents a log entry in JSON format
type JSONLogEntry struct {
	Timestamp     string            `json:"timestamp"`
	Level         string            `json:"level"`
	Message       string            `json:"message"`
	Attributes    map[string]string `json:"attributes,omitempty"`
	TransactionID string            `json:"transaction_id,omitempty"`
}

// NewJSONFileDriver creates a new JSON file driver from a map of options
func NewJSONFileDriver(options map[string]interface{}) (core.Driver, error) {
	filePath, ok := options["file_path"].(string)
	if !ok || filePath == "" {
		return nil, fmt.Errorf("file_path is required")
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	driver := &JSONFileDriver{
		filePath: filePath,
		file:     file,
		encoder:  json.NewEncoder(file),
		minLevel: core.Debug,
	}

	if levelStr, ok := options["min_level"].(string); ok {
		if level, err := core.ParseLevel(levelStr); err == nil {
			driver.minLevel = level
		}
	}

	return driver, nil
}

// Log writes a log entry to the JSON file
func (d *JSONFileDriver) Log(entry *core.LogEntry) error {
	if entry.Level < d.minLevel {
		return nil
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if d.file == nil {
		return fmt.Errorf("driver is closed")
	}

	jsonEntry := struct {
		Timestamp     string            `json:"timestamp"`
		Level         string            `json:"level"`
		Message       string            `json:"message"`
		Attributes    map[string]string `json:"attributes,omitempty"`
		TransactionID string            `json:"transaction_id,omitempty"`
	}{
		Timestamp:     entry.Timestamp.Format("2006-01-02T15:04:05.000Z07:00"),
		Level:         entry.Level.String(),
		Message:       entry.Message,
		Attributes:    entry.Attrs,
		TransactionID: entry.TransactionID,
	}

	return d.encoder.Encode(jsonEntry)
}

// Close closes the file
func (d *JSONFileDriver) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.file == nil {
		return nil
	}

	err := d.file.Close()
	d.file = nil
	return err
}
