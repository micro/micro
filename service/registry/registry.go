// Package registry is the micro registry
package registry

import (
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/micro/v2/service/registry/client"
)

var (
	// DefaultRegistry implementation
	DefaultRegistry registry.Registry = client.NewRegistry()
)
