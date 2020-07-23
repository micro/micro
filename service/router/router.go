package router

import (
	"github.com/micro/go-micro/v2/router"
	"github.com/micro/go-micro/v2/router/registry"
	muregistry "github.com/micro/micro/v2/service/registry"
)

var (
	// DefaultRouter implementation
	DefaultRouter router.Router = registry.NewRouter(
		router.Registry(muregistry.DefaultRegistry),
	)
)
