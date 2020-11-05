// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Original source: github.com/micro/go-micro/v3/config/config.go

// Package config is an interface for dynamic configuration.
package config

import (
	"time"
)

// DefaultConfig implementation
var DefaultConfig Config

// Config is an interface abstraction for dynamic configuration
type Config interface {
	Get(path string, options ...Option) (Value, error)
	Set(path string, val interface{}, options ...Option) error
	Delete(path string, options ...Option) error
}

// Value represents a value of any type
type Value interface {
	Exists() bool
	Bool(def bool) bool
	Int(def int) int
	String(def string) string
	Float64(def float64) float64
	Duration(def time.Duration) time.Duration
	StringSlice(def []string) []string
	StringMap(def map[string]string) map[string]string
	Scan(val interface{}) error
	Bytes() []byte
}

type Options struct {
	Secret bool
}

type Option func(o *Options)

func Secret(b bool) Option {
	return func(o *Options) {
		o.Secret = b
	}
}

type Secrets interface {
	Config
}

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
