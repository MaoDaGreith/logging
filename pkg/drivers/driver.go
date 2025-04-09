package drivers

import (
	"github.com/MaoDaGreith/logging/pkg/core"
)

// DriverConstructor is a function type that creates a new driver instance from options
type DriverConstructor func(options map[string]interface{}) (core.Driver, error)

// registry holds all registered driver constructors
var registry = make(map[string]DriverConstructor)

// Register adds a driver constructor to the registry
func Register(name string, constructor DriverConstructor) {
	registry[name] = constructor
}

// Create instantiates a driver by name with the given options
func Create(name string, options map[string]interface{}) (core.Driver, error) {
	if constructor, ok := registry[name]; ok {
		return constructor(options)
	}

	return nil, core.ErrDriverNotFound
}
