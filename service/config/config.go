package config

import (
	"github.com/micro/go-micro/v3/config"
	"github.com/micro/go-micro/v3/config/reader"
)

// DefaultConfig implementation. Setup in the cmd package, this will
// be refactored following the updated config interface.
var DefaultConfig config.Config

// Bytes representation of config
func Bytes() []byte {
	return DefaultConfig.Bytes()
}

// Get a value at the path
func Get(path ...string) reader.Value {
	return DefaultConfig.Get(path...)
}

// Set the value at a path
func Set(val interface{}, path ...string) {
	DefaultConfig.Set(val)
}

// Delete the value at a path
func Delete(path ...string) {
	DefaultConfig.Del(path...)
}

// Map represesentation of config
func Map() map[string]interface{} {
	return DefaultConfig.Map()
}

// Scan config into the value provided
func Scan(v interface{}) error {
	return DefaultConfig.Scan(v)
}
