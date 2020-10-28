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
// Original source: github.com/micro/micro/v3/router/watcher.go

package router

import (
	"errors"
	"time"
)

var (
	// ErrWatcherStopped is returned when routing table watcher has been stopped
	ErrWatcherStopped = errors.New("watcher stopped")
)

// EventType defines routing table event
type EventType int

const (
	// Create is emitted when a new route has been created
	Create EventType = iota
	// Delete is emitted when an existing route has been deleted
	Delete
	// Update is emitted when an existing route has been updated
	Update
)

// String returns human readable event type
func (t EventType) String() string {
	switch t {
	case Create:
		return "create"
	case Delete:
		return "delete"
	case Update:
		return "update"
	default:
		return "unknown"
	}
}

// Event is returned by a call to Next on the watcher.
type Event struct {
	// Unique id of the event
	Id string
	// Type defines type of event
	Type EventType
	// Timestamp is event timestamp
	Timestamp time.Time
	// Route is table route
	Route Route
}

// Watcher defines routing table watcher interface
// Watcher returns updates to the routing table
type Watcher interface {
	// Next is a blocking call that returns watch result
	Next() (*Event, error)
	// Chan returns event channel
	Chan() (<-chan *Event, error)
	// Stop stops watcher
	Stop()
}

// WatchOption is used to define what routes to watch in the table
type WatchOption func(*WatchOptions)

// WatchOptions are table watcher options
// TODO: expand the options to watch based on other criteria
type WatchOptions struct {
	// Service allows to watch specific service routes
	Service string
}

// WatchService sets what service routes to watch
// Service is the microservice name
func WatchService(s string) WatchOption {
	return func(o *WatchOptions) {
		o.Service = s
	}
}
