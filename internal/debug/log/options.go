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
// Original source: github.com/micro/go-micro/v3/debug/log/options.go

package log

import "time"

// Option used by the logger
type Option func(*Options)

// Options are logger options
type Options struct {
	// Name of the log
	Name string
	// Size is the size of ring buffer
	Size int
	// Format specifies the output format
	Format FormatFunc
}

// Name of the log
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// Size sets the size of the ring buffer
func Size(s int) Option {
	return func(o *Options) {
		o.Size = s
	}
}

func Format(f FormatFunc) Option {
	return func(o *Options) {
		o.Format = f
	}
}

// DefaultOptions returns default options
func DefaultOptions() Options {
	return Options{
		Size: DefaultSize,
	}
}

// ReadOptions for querying the logs
type ReadOptions struct {
	// Since what time in past to return the logs
	Since time.Time
	// Count specifies number of logs to return
	Count int
	// Stream requests continuous log stream
	Stream bool
}

// ReadOption used for reading the logs
type ReadOption func(*ReadOptions)

// Since sets the time since which to return the log records
func Since(s time.Time) ReadOption {
	return func(o *ReadOptions) {
		o.Since = s
	}
}

// Count sets the number of log records to return
func Count(c int) ReadOption {
	return func(o *ReadOptions) {
		o.Count = c
	}
}
