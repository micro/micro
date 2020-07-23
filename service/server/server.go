package server

import (
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-micro/v2/server/grpc"
)

// DefaultServer for the service
var DefaultServer server.Server = grpc.NewServer()
