package config

import (
	"github.com/micro/go-micro/v3/config"
)

// DefaultConfig implementation. Setup in the cmd package, this will
// be refactored following the updated config interface.
var DefaultConfig config.Config

type (
	// Value is an alias for reader.Value
	Value = config.Value
)

// Get a value at the path
func Get(path ...string) Value {
	return DefaultConfig.Get(path...)
}

// Set the value at a path
func Set(val interface{}, path ...string) {
	DefaultConfig.Set(val)
}

// Delete the value at a path
func Delete(path string) {
	DefaultConfig.Delete(path)
}
