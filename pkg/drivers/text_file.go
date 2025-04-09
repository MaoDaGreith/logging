package drivers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/MaoDaGreith/logging/pkg/core"
)

// TextFileDriverName is the name to use in configuration
const TextFileDriverName = "text_file"

func init() {
	Register(TextFileDriverName, NewTextFileDriver)
}

// TextFileDriver outputs logs to a text file
type TextFileDriver struct {
	filePath string
	file     *os.File
	minLevel core.Level
	mu       sync.Mutex
}

// NewTextFileDriver creates a new text file driver from a map of options
func NewTextFileDriver(options map[string]interface{}) (core.Driver, error) {
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

	driver := &TextFileDriver{
		filePath: filePath,
		file:     file,
		minLevel: core.Debug,
	}

	if levelStr, ok := options["min_level"].(string); ok {
		if level, err := core.ParseLevel(levelStr); err == nil {
			driver.minLevel = level
		}
	}

	return driver, nil
}

// Log writes a log entry to the text file
func (d *TextFileDriver) Log(entry *core.LogEntry) error {
	if entry.Level < d.minLevel {
		return nil
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if d.file == nil {
		return fmt.Errorf("driver is closed")
	}

	var builder strings.Builder
	builder.WriteString(entry.Timestamp.Format("2006-01-02T15:04:05.000Z07:00"))
	builder.WriteString(" [")
	builder.WriteString(entry.Level.String())
	builder.WriteString("] ")
	builder.WriteString(entry.Message)

	if len(entry.Attrs) > 0 {
		builder.WriteString(" {")
		first := true
		for k, v := range entry.Attrs {
			if !first {
				builder.WriteString(", ")
			}
			first = false
			builder.WriteString(k)
			builder.WriteString("=")
			builder.WriteString(v)
		}
		builder.WriteString("}")
	}

	if entry.TransactionID != "" {
		builder.WriteString(" (txn: ")
		builder.WriteString(entry.TransactionID)
		builder.WriteString(")")
	}

	builder.WriteString("\n")

	_, err := d.file.WriteString(builder.String())
	return err
}

// Close closes the file
func (d *TextFileDriver) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.file == nil {
		return nil
	}

	err := d.file.Close()
	d.file = nil
	return err
}
