// Package runtime is the micro runtime
package runtime

import (
	"github.com/micro/go-micro/v3/runtime"
	"github.com/micro/micro/v3/service/runtime/client"
)

var (
	// DefaultRuntime implementation
	DefaultRuntime runtime.Runtime = client.NewRuntime()
)

// Create registers a service
func Create(srv *runtime.Service, opts ...runtime.CreateOption) error {
	return DefaultRuntime.Create(srv, opts...)
}

// Read returns the service
func Read(opts ...runtime.ReadOption) ([]*runtime.Service, error) {
	return DefaultRuntime.Read(opts...)
}

// Update the service in place
func Update(srv *runtime.Service, opts ...runtime.UpdateOption) error {
	return DefaultRuntime.Update(srv, opts...)
}

// Delete a service
func Delete(srv *runtime.Service, opts ...runtime.DeleteOption) error {
	return DefaultRuntime.Delete(srv, opts...)
}

// Logs returns the logs for a service
func Logs(srv *runtime.Service, opts ...runtime.LogsOption) (runtime.LogStream, error) {
	return DefaultRuntime.Logs(srv, opts...)
}
