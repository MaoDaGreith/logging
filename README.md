# Logging Library

A simple and flexible logging library for Go applications.

## Features

- Multiple logging drivers (Console, JSON file, Text file)
- Configurable log levels
- Thread-safe logging
- Support for structured logging
- Easy to extend with custom drivers
- Configuration via JSON or YAML files

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

The library can be configured via JSON or YAML configuration files. A sample configuration file (`config.yaml.sample`) is provided in the root directory.

### YAML Configuration Example

```yaml
logging:
  level: info
  timestamp_format: "2006-01-02 15:04:05.000"
  
  drivers:
    - type: console
      options:
        format: text
        output: stdout
        colors: true
    
    - type: json_file
      options:
        file_path: "logs/app.json"
        max_size: 10485760    # 10MB
        max_backups: 5
        max_age: 30
    
    - type: text_file
      options:
        file_path: "logs/app.log"
        max_size: 10485760    # 10MB
        max_backups: 5
        max_age: 30
        format: "[%timestamp%] [%level%] %message%"
```

### JSON Configuration Example

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

The library looks for configuration in these locations:
1. Path specified in `LOGGING_CONFIG_PATH` environment variable
2. `config/logging.json`
3. `/etc/logging/config.json`

If no configuration is found, it falls back to a default configuration with just a console driver.

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

