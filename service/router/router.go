package router

import (
	"github.com/micro/go-micro/v3/router"
	"github.com/micro/micro/v2/service/router/client"
)

var (
	// DefaultRouter implementation
	DefaultRouter router.Router = client.NewRouter()
)
