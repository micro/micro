// Package registry is the micro registry
package registry

import (
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/service"
)

var (
	// DefaultRegistry implementation
	DefaultRegistry registry.Registry = service.NewRegistry()
)
