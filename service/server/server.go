package server

import (
	"github.com/micro/go-micro/v3/server"
	"github.com/micro/go-micro/v3/server/grpc"
	"github.com/micro/micro/v2/service/registry"
)

// DefaultServer for the service
var DefaultServer server.Server = grpc.NewServer(
	server.Registry(registry.DefaultRegistry),
)
