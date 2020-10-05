// Package runtime is the micro runtime
package runtime

import (
	"github.com/micro/go-micro/v3/runtime"
	"github.com/micro/micro/v3/service/runtime/client"
)

var (
	// DefaultRuntime implementation
	DefaultRuntime runtime.Runtime = client.NewRuntime()
	// ErrAlreadyExists is an alias for runtime.ErrAlreadyExists
	ErrAlreadyExists = runtime.ErrAlreadyExists
	// ErrNotFound is an alias for runtime.ErrNotFound
	ErrNotFound = runtime.ErrNotFound
)

type (
	// Logs is an alias for runtime.Logs
	Logs = runtime.Logs
	// Resource is an alias for runtime.Resource
	Resource = runtime.Resource
	// Service is an alias for runtime.Service
	Service = runtime.Service
	// ServiceStatus is an alias for runtime.ServiceStatus
	ServiceStatus = runtime.ServiceStatus
)

const (
	// Unknown indicates the status of the service is not known
	Unknown = runtime.Unknown
	// Pending is the initial status of a service
	Pending = runtime.Pending
	// Building is the status when the service is being built
	Building = runtime.Building
	// Starting is the status when the service has been started but is not yet ready to accept traffic
	Starting = runtime.Starting
	// Running is the status when the service is active and accepting traffic
	Running = runtime.Running
	// Stopping is the status when a service is stopping
	Stopping = runtime.Stopping
	// Stopped is the status when a service has been stopped or has completed
	Stopped = runtime.Stopped
	// Error is the status when an error occured, this could be a build error or a run error. The error
	// details can be found within the service's metadata
	Error = runtime.Error
)

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

// Log returns the logs for a resource (service)
func Log(resource Resource, opts ...LogsOption) (Logs, error) {
	return DefaultRuntime.Logs(resource, opts...)
}
