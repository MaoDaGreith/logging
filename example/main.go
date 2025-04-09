package main

import (
	"fmt"
	"os"
	"time"

	"github.com/MaoDaGreith/logging/pkg/config"
	"github.com/MaoDaGreith/logging/pkg/core"
	"github.com/MaoDaGreith/logging/pkg/drivers"
)

func main() {
	// Example 1: Create a logger with console driver directly
	fmt.Println("Example 1: Console Driver")
	consoleLogger := core.NewLogger(drivers.NewConsoleDriverWithOptions())

	consoleLogger.Debug("This is a debug message")
	consoleLogger.Info("This is an info message")
	consoleLogger.Warning("This is a warning message", core.Attributes{
		"source": "example",
	})
	consoleLogger.Error("This is an error message", core.Attributes{
		"code": "500",
	})

	// Example 2: Transaction logging
	fmt.Println("\nExample 2: Transaction Logging")
	tx := consoleLogger.NewTransaction("request-123")
	tx.Info("Processing request")
	tx.Debug("Validating input")
	tx.Info("Executing database query", core.Attributes{
		"table": "users",
		"query": "SELECT * FROM users WHERE id = ?",
	})
	tx.Warning("Slow query detected", core.Attributes{
		"duration_ms": "150",
	})
	tx.Info("Request completed", core.Attributes{
		"status": "success",
	})

	// Example 3: Multiple drivers
	fmt.Println("\nExample 3: Multiple Drivers")

	// Create a temporary JSON log file
	jsonFilePath := "example_log.json"
	defer os.Remove(jsonFilePath) // Clean up after example

	jsonDriver, err := drivers.NewJSONFileDriver(map[string]interface{}{
		"file_path": jsonFilePath,
		"min_level": "warning",
	})
	if err != nil {
		fmt.Printf("Failed to create JSON driver: %v\n", err)
		return
	}

	// Create a temporary text log file
	textFilePath := "example_log.txt"
	defer os.Remove(textFilePath) // Clean up after example

	textDriver, err := drivers.NewTextFileDriver(map[string]interface{}{
		"file_path": textFilePath,
	})
	if err != nil {
		fmt.Printf("Failed to create text driver: %v\n", err)
		return
	}

	// Create a logger with multiple drivers
	multiLogger := core.NewLogger(
		drivers.NewConsoleDriverWithOptions(drivers.WithMinLevel(core.Info)),
		jsonDriver,
		textDriver,
	)

	multiLogger.Debug("This debug message only goes to text file")      // Not to console or JSON
	multiLogger.Info("This info message goes to console and text file") // Not to JSON
	multiLogger.Warning("This warning goes to all drivers")
	multiLogger.Error("This error goes to all drivers")

	// Close the logger to ensure file drivers flush and close
	multiLogger.Close()

	fmt.Printf("\nLogs written to %s and %s\n", jsonFilePath, textFilePath)

	// Example 4: Configuration-based logger
	fmt.Println("\nExample 4: Configuration-based Logger")

	// Create a configuration
	cfg := &config.Config{
		DefaultLevel: "info",
		Drivers: []config.DriverConfig{
			{
				Type:     "console",
				MinLevel: "info",
			},
			{
				Type:     "text_file",
				MinLevel: "debug",
				Options: map[string]interface{}{
					"file_path": "config_example.log",
				},
			},
		},
	}

	// Save configuration to file for demonstration
	configPath := "example_config.json"
	if err := cfg.SaveToFile(configPath); err != nil {
		fmt.Printf("Failed to save config: %v\n", err)
		return
	}
	defer os.Remove(configPath)           // Clean up after example
	defer os.Remove("config_example.log") // Clean up after example

	// Load configuration from file
	loadedCfg, err := config.LoadFromFile(configPath)
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}

	// Create logger from configuration
	configLogger, err := loadedCfg.CreateLogger()
	if err != nil {
		fmt.Printf("Failed to create logger from config: %v\n", err)
		return
	}

	configLogger.Debug("This debug message only goes to text file") // Not to console
	configLogger.Info("This info message goes to both drivers")
	configLogger.Warning("This warning goes to both drivers")
	configLogger.Error("This error goes to both drivers")

	// Create a transaction with the config logger
	configTx := configLogger.NewTransaction("config-tx-123")
	configTx.Info("Transaction from config-based logger")

	// Close the logger
	configLogger.Close()

	fmt.Printf("\nConfig saved to %s and logs written to config_example.log\n", configPath)

	// Wait a moment to ensure files are written
	time.Sleep(100 * time.Millisecond)
}
