package config

import (
	"github.com/micro/go-micro/v3/config"
)

// DefaultConfig implementation. Setup in the cmd package, this will
// be refactored following the updated config interface.
var DefaultConfig config.Config

type (
	// Value is an alias for reader.Value
	Value   = config.Value
	Options = config.Options
	Option  = config.Option
)

// Get a value at the path
func Get(path string, options ...Option) (Value, error) {
	return DefaultConfig.Get(path, options...)
}

// Set the value at a path
func Set(path string, val interface{}, options ...Option) error {
	return DefaultConfig.Set(path, val, options...)
}

// Delete the value at a path
func Delete(path string, options ...Option) error {
	return DefaultConfig.Delete(path, options...)
}

var Secret = config.Secret
