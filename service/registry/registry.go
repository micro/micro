// Package registry is the micro registry
package registry

import (
	"github.com/micro/go-micro/v3/registry"
	"github.com/micro/micro/v3/service/registry/client"
)

var (
	// DefaultRegistry implementation
	DefaultRegistry registry.Registry = client.NewRegistry()
)

type (
	// Service is an alias for registry.Service
	Service = registry.Service
	// Watcher is an alias for registry.Watcher
	Watcher = registry.Watcher
)

// Register a service
func Register(service *Service, opts ...registry.RegisterOption) error {
	return DefaultRegistry.Register(service, opts...)
}

// Deregister a service
func Deregister(service *Service, opts ...registry.DeregisterOption) error {
	return DefaultRegistry.Deregister(service, opts...)
}

// GetService from the registry
func GetService(service string, opts ...registry.GetOption) ([]*Service, error) {
	return DefaultRegistry.GetService(service, opts...)
}

// ListServices in the registry
func ListServices(opts ...registry.ListOption) ([]*Service, error) {
	return DefaultRegistry.ListServices(opts...)
}

// Watch the registry for updates
func Watch(opts ...registry.WatchOption) (Watcher, error) {
	return DefaultRegistry.Watch(opts...)
}
