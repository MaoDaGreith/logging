package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/MaoDaGreith/logging/pkg/core"
	"github.com/MaoDaGreith/logging/pkg/drivers"
	"gopkg.in/yaml.v3"
)

// Config represents the logger configuration
type Config struct {
	Logger       Logger
	DefaultLevel string         `json:"default_level"`
	Drivers      []DriverConfig `json:"drivers"`
}

// DriverConfig represents a single driver configuration
type DriverConfig struct {
	Type     string                 `json:"type"`
	MinLevel string                 `json:"min_level"`
	Options  map[string]interface{} `json:"options,omitempty"`
}

// Logger is the main interface for logging
type Logger interface {
	Debug(msg string, attrs ...core.Attributes) error
	Info(msg string, attrs ...core.Attributes) error
	Warning(msg string, attrs ...core.Attributes) error
	Error(msg string, attrs ...core.Attributes) error
	Log(level core.Level, msg string, attrs ...core.Attributes) error
	NewTransaction(txID string) core.Transaction
	Close() error
}

// LoadFromFile loads a configuration from a file
func LoadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	logger, err := config.CreateLogger()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	config.Logger = logger
	return &config, nil
}

// LoadDefault loads the default configuration from a standard location
func LoadDefault() (*Config, error) {
	if path := os.Getenv("LOGGING_CONFIG_PATH"); path != "" {
		return LoadFromFile(path)
	}

	locations := []string{
		"config/logging.json",
		"/etc/logging/config.json",
	}

	for _, loc := range locations {
		if _, err := os.Stat(loc); err == nil {
			return LoadFromFile(loc)
		}
	}

	return &Config{
		DefaultLevel: "info",
		Drivers: []DriverConfig{
			{
				Type:     "console",
				MinLevel: "debug",
			},
		},
	}, nil
}

// CreateLogger creates a logger from a configuration
func (c *Config) CreateLogger() (Logger, error) {
	driverInstances := make([]core.Driver, 0, len(c.Drivers))

	for _, driverConfig := range c.Drivers {
		if driverConfig.Options == nil {
			driverConfig.Options = make(map[string]interface{})
		}

		if _, ok := driverConfig.Options["min_level"]; !ok && driverConfig.MinLevel != "" {
			driverConfig.Options["min_level"] = driverConfig.MinLevel
		}

		driver, err := drivers.Create(driverConfig.Type, driverConfig.Options)
		if err != nil {
			return nil, fmt.Errorf("failed to create driver '%s': %w", driverConfig.Type, err)
		}

		driverInstances = append(driverInstances, driver)
	}

	return core.NewLogger(driverInstances...), nil
}

// SaveToFile saves the configuration to a file
func (c *Config) SaveToFile(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(c); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	return nil
}
