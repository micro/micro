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
// Original source: github.com/micro/go-micro/v3/runtime/runtime.go

// Package runtime is a service runtime manager
package runtime

import (
	"errors"
	"time"
)

var (
	// DefaultRuntime implementation
	DefaultRuntime Runtime

	ErrAlreadyExists   = errors.New("already exists")
	ErrInvalidResource = errors.New("invalid resource")
	ErrNotFound        = errors.New("not found")
)

// Runtime is a service runtime manager
type Runtime interface {
	// Init initializes runtime
	Init(...Option) error
	// Create a resource
	Create(Resource, ...CreateOption) error
	// Read a resource
	Read(...ReadOption) ([]*Service, error)
	// Update a resource
	Update(Resource, ...UpdateOption) error
	// Delete a resource
	Delete(Resource, ...DeleteOption) error
	// Logs returns the logs for a resource
	Logs(Resource, ...LogsOption) (LogStream, error)
	// Start starts the runtime
	Start() error
	// Stop shuts down the runtime
	Stop() error
	// String defines the runtime implementation
	String() string
}

// LogStream returns a log stream
type LogStream interface {
	Error() error
	Chan() chan Log
	Stop() error
}

// Log is a log message
type Log struct {
	Message  string
	Metadata map[string]string
}

// EventType defines schedule event
type EventType int

const (
	// CreateEvent is emitted when a new build has been craeted
	CreateEvent EventType = iota
	// UpdateEvent is emitted when a new update become available
	UpdateEvent
	// DeleteEvent is emitted when a build has been deleted
	DeleteEvent
)

// String returns human readable event type
func (t EventType) String() string {
	switch t {
	case CreateEvent:
		return "create"
	case DeleteEvent:
		return "delete"
	case UpdateEvent:
		return "update"
	default:
		return "unknown"
	}
}

// Event is notification event
type Event struct {
	// ID of the event
	ID string
	// Type is event type
	Type EventType
	// Timestamp is event timestamp
	Timestamp time.Time
	// Service the event relates to
	Service *Service
	// Options to use when processing the event
	Options *CreateOptions
}

// ServiceStatus defines service statuses
type ServiceStatus int

const (
	// Unknown indicates the status of the service is not known
	Unknown ServiceStatus = iota
	// Pending is the initial status of a service
	Pending
	// Building is the status when the service is being built
	Building
	// Starting is the status when the service has been started but is not yet ready to accept traffic
	Starting
	// Running is the status when the service is active and accepting traffic
	Running
	// Stopping is the status when a service is stopping
	Stopping
	// Stopped is the status when a service has been stopped or has completed
	Stopped
	// Error is the status when an error occured, this could be a build error or a run error. The error
	// details can be found within the service's metadata
	Error
)

// Resources which are allocated to a serivce
type Resources struct {
	// CPU is the maximum amount of CPU the service will be allocated (unit millicpu)
	// e.g. 0.25CPU would be passed as 250
	CPU int
	// Mem is the maximum amount of memory the service will be allocated (unit mebibyte)
	// e.g. 128 MiB of memory would be passed as 128
	Mem int
	// Disk is the maximum amount of disk space the service will be allocated (unit mebibyte)
	// e.g. 128 MiB of memory would be passed as 128
	Disk int
}

// Create a resource
func Create(resource Resource, opts ...CreateOption) error {
	return DefaultRuntime.Create(resource, opts...)
}

// Read returns the service
func Read(opts ...ReadOption) ([]*Service, error) {
	return DefaultRuntime.Read(opts...)
}

// Update the resource in place
func Update(resource Resource, opts ...UpdateOption) error {
	return DefaultRuntime.Update(resource, opts...)
}

// Delete a resource
func Delete(resource Resource, opts ...DeleteOption) error {
	return DefaultRuntime.Delete(resource, opts...)
}

// Logs for a resource
func Logs(resource Resource, opts ...LogsOption) (LogStream, error) {
	return DefaultRuntime.Logs(resource, opts...)
}
