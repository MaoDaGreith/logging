package drivers

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/MaoDaGreith/logging/pkg/core"
)

// ConsoleDriverName is the name to use in configuration
const ConsoleDriverName = "console"

func init() {
	Register(ConsoleDriverName, NewConsoleDriver)
}

// ConsoleDriver outputs logs to the console (stdout/stderr)
type ConsoleDriver struct {
	stdout     io.Writer
	stderr     io.Writer
	minLevel   core.Level
	timeFormat string
	colorized  bool
}

// ConsoleDriverOption represents an option for the console driver
type ConsoleDriverOption func(*ConsoleDriver)

// WithMinLevel sets the minimum log level to output
func WithMinLevel(level core.Level) ConsoleDriverOption {
	return func(d *ConsoleDriver) {
		d.minLevel = level
	}
}

// WithTimeFormat sets the time format for logs
func WithTimeFormat(format string) ConsoleDriverOption {
	return func(d *ConsoleDriver) {
		d.timeFormat = format
	}
}

// WithColorized enables/disables colorized output
func WithColorized(colorized bool) ConsoleDriverOption {
	return func(d *ConsoleDriver) {
		d.colorized = colorized
	}
}

// WithStdout sets the output writer for non-error logs
func WithStdout(w io.Writer) ConsoleDriverOption {
	return func(d *ConsoleDriver) {
		d.stdout = w
	}
}

// WithStderr sets the output writer for error logs
func WithStderr(w io.Writer) ConsoleDriverOption {
	return func(d *ConsoleDriver) {
		d.stderr = w
	}
}

// NewConsoleDriverWithOptions creates a new console driver with options
func NewConsoleDriverWithOptions(options ...ConsoleDriverOption) *ConsoleDriver {
	driver := &ConsoleDriver{
		stdout:     os.Stdout,
		stderr:     os.Stderr,
		minLevel:   core.Debug,
		timeFormat: time.RFC3339,
		colorized:  true,
	}

	for _, option := range options {
		option(driver)
	}

	return driver
}

// NewConsoleDriver creates a new console driver from a map of options
func NewConsoleDriver(options map[string]interface{}) (core.Driver, error) {
	driver := &ConsoleDriver{
		stdout:     os.Stdout,
		stderr:     os.Stderr,
		minLevel:   core.Debug,
		timeFormat: time.RFC3339,
		colorized:  true,
	}

	if levelStr, ok := options["min_level"].(string); ok {
		if level, err := core.ParseLevel(levelStr); err == nil {
			driver.minLevel = level
		}
	}

	if format, ok := options["time_format"].(string); ok {
		driver.timeFormat = format
	}

	if colorized, ok := options["colorized"].(bool); ok {
		driver.colorized = colorized
	}

	return driver, nil
}

// Log writes a log entry to the console
func (d *ConsoleDriver) Log(entry *core.LogEntry) error {
	if entry.Level < d.minLevel {
		return nil
	}

	// Format the log entry
	formatted := d.format(entry)

	// Write to the appropriate output
	var out io.Writer
	if entry.Level >= core.Error {
		out = d.stderr
	} else {
		out = d.stdout
	}

	_, err := fmt.Fprintln(out, formatted)
	return err
}

// Close is a no-op for the console driver
func (d *ConsoleDriver) Close() error {
	return nil
}

// format formats a log entry as a string
func (d *ConsoleDriver) format(entry *core.LogEntry) string {
	// Format timestamp
	timestamp := entry.Timestamp.Format(d.timeFormat)

	// Format log level
	levelStr := entry.Level.String()
	if d.colorized {
		levelStr = d.colorizeLevel(entry.Level, levelStr)
	}

	// Format message
	message := entry.Message

	// Format attributes
	attrsStr := ""
	if len(entry.Attrs) > 0 {
		attrs := make([]string, 0, len(entry.Attrs))
		for k, v := range entry.Attrs {
			attrs = append(attrs, fmt.Sprintf("%s=%s", k, v))
		}
		attrsStr = fmt.Sprintf(" [%s]", strings.Join(attrs, ", "))
	}

	// Format transaction ID
	txID := ""
	if entry.TransactionID != "" {
		txID = fmt.Sprintf(" (tx: %s)", entry.TransactionID)
	}

	return fmt.Sprintf("%s [%s]%s%s %s", timestamp, levelStr, txID, attrsStr, message)
}

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
)

// colorizeLevel adds ANSI color codes to log level
func (d *ConsoleDriver) colorizeLevel(level core.Level, levelStr string) string {
	switch level {
	case core.Debug:
		return colorBlue + levelStr + colorReset
	case core.Info:
		return colorGreen + levelStr + colorReset
	case core.Warning:
		return colorYellow + levelStr + colorReset
	case core.Error:
		return colorRed + levelStr + colorReset
	default:
		return levelStr
	}
}
