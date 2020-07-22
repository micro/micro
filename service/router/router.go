package router

import (
	"github.com/micro/go-micro/v2/router"
	"github.com/micro/go-micro/v2/router/service"
)

var (
	// DefaultRouter implementation
	DefaultRouter router.Router = service.NewRouter()
)
