# Logging/Telemetry Package - Project Summary

## Overview

This project implements a flexible logging and telemetry library for Go applications. The package allows for logging at different severity levels, with metadata attributes, and supports transaction-based logging for grouping related log entries. The library can output logs to multiple destinations through a driver-based system.

## Key Features

1. **Multiple Log Levels**
   - Debug, Info, Warning, Error levels supported
   - Each log entry includes a timestamp

2. **Metadata Support**
   - Logs can include key-value attributes for additional context
   - Attributes are stored as string key-value pairs

3. **Transaction Logging**
   - Logs can be grouped by transaction ID
   - Useful for tracking operations across multiple components

4. **Multiple Output Drivers**
   - Console output with optional colorization
   - JSON file output for structured logging
   - Text file output for human-readable logs
   - Extensible driver system for custom outputs

5. **Configuration System**
   - Load configuration from JSON files
   - Environment-aware configuration loading
   - Driver-specific configuration options

## Architecture

The package is organized into the following components:

1. **Core Package**
   - Defines log levels (Debug, Info, Warning, Error)
   - Implements the main Logger interface
   - Implements Transaction interface for grouping logs
   - Defines the Driver interface for output plugins

2. **Drivers Package**
   - Implementation of Console, JSON file, and Text file drivers
   - Driver registry for extensibility
   - Factory functions for creating drivers from configuration

3. **Config Package**
   - Configuration loading from JSON files
   - Default configuration provision
   - Logger construction from configuration

## Usage Examples

1. **Basic Usage**
   ```go
   logger := core.NewLogger(drivers.NewConsoleDriverWithOptions())
   logger.Debug("Debug message")
   logger.Info("Info message")
   logger.Warning("Warning message", core.Attributes{"source": "example"})
   logger.Error("Error message", core.Attributes{"code": "500"})
   ```

2. **Transaction Logging**
   ```go
   tx := logger.NewTransaction("request-123")
   tx.Info("Processing request")
   tx.Info("Database query executed", core.Attributes{
     "table": "users",
     "duration_ms": "50",
   })
   tx.Info("Request completed")
   ```

3. **Configuration-based Usage**
   ```go
   // Load configuration from file
   cfg, _ := config.LoadFromFile("logging.json")
   
   // Create logger from configuration
   logger, _ := cfg.CreateLogger()
   
   // Use the logger
   logger.Info("Application started")
   ```

## Extensibility

The system can be extended with custom drivers by:

1. Implementing the `core.Driver` interface
2. Registering the driver with a name
3. Creating a constructor function

```go
// Custom driver implementation
type CustomDriver struct {
    // Fields
}

// Implement the Driver interface
func (d *CustomDriver) Log(entry *core.LogEntry) error {
    // Implementation
}

func (d *CustomDriver) Close() error {
    // Implementation
}

// Register the driver
func init() {
    drivers.Register("custom", NewCustomDriver)
}

// Constructor function
func NewCustomDriver(options map[string]interface{}) (core.Driver, error) {
    // Create and return driver
}
```

## Testing

The package includes a comprehensive test suite covering:
- Log level parsing and formatting
- Logger core functionality
- Transaction logging behavior
- Driver interaction

## Future Improvements

Potential enhancements for the future:
1. Additional drivers (Syslog, ELK, Cloud logging services)
2. Structured logging with different formats (e.g., logfmt)
3. Log rotation capabilities
4. Sampling and filtering options for high-volume logs
5. Performance optimizations for high-throughput scenarios 