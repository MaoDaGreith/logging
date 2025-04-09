# Logging Library

A simple and flexible logging library for Go applications.

## Features

- Multiple logging drivers (Console, JSON file, Text file)
- Configurable log levels
- Thread-safe logging
- Support for structured logging
- Easy to extend with custom drivers

## Installation

```bash
go get github.com/MaoDaGreith/logging
```

## Usage

```go
package main

import (
    "github.com/MaoDaGreith/logging/pkg/core"
    "github.com/MaoDaGreith/logging/pkg/drivers"
)

func main() {
    // Create a new logger with console driver
    logger := core.NewLogger(drivers.NewConsoleDriver())
    
    // Basic logging
    logger.Debug("Debug message")
    logger.Info("Info message")
    logger.Warning("Warning message")
    logger.Error("Error message")
    
    // Logging with attributes
    logger.Info("User login", core.Attributes{
        "userId": "123",
        "source": "web",
    })
    
    // Transaction logging
    tx := logger.NewTransaction("request-123")
    tx.Info("Processing request")
    tx.Info("Database query executed", core.Attributes{
        "table": "users",
        "duration_ms": "50",
    })
    tx.Info("Request completed")
}
```

## Configuration

The library can be configured via a JSON configuration file:

```json
{
  "default_level": "info",
  "drivers": [
    {
      "type": "console",
      "min_level": "debug"
    },
    {
      "type": "json_file",
      "min_level": "warning",
      "options": {
        "file_path": "/var/log/app.json"
      }
    }
  ]
}
```

## Extending with Custom Drivers

To create a custom driver, implement the `Driver` interface and register it with the logger:

```go
type CustomDriver struct {
    // Your driver fields
}

func (d *CustomDriver) Log(entry *core.LogEntry) error {
    // Implement your logging logic here
    return nil
}

// Then use it with the logger
logger := core.NewLogger(&CustomDriver{})
```

