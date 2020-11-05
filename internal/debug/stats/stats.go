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
// Original source: github.com/micro/go-micro/v3/debug/stats/stats.go

// Package stats provides runtime stats
package stats

// Stats provides stats interface
type Stats interface {
	// Read stat snapshot
	Read() ([]*Stat, error)
	// Write a stat snapshot
	Write(*Stat) error
	// Record a request
	Record(error) error
}

// A runtime stat
type Stat struct {
	// Timestamp of recording
	Timestamp int64
	// Start time as unix timestamp
	Started int64
	// Uptime in seconds
	Uptime int64
	// Memory usage in bytes
	Memory uint64
	// Threads aka go routines
	Threads uint64
	// Garbage collection in nanoseconds
	GC uint64
	// Total requests
	Requests uint64
	// Total errors
	Errors uint64
}
