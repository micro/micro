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
	// Node is an alias for registry.Node
	Node = registry.Node
	// Service is an alias for registry.Service
	Service = registry.Service
	// Watcher is an alias for registry.Watcher
	Watcher = registry.Watcher
)

// GetService from the registry
func GetService(service string) ([]*Service, error) {
	return DefaultRegistry.GetService(service)
}

// ListServices in the registry
func ListServices() ([]*Service, error) {
	return DefaultRegistry.ListServices()
}

// Watch the registry for updates
func Watch() (Watcher, error) {
	return DefaultRegistry.Watch()
}
