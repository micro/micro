package client

import (
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/grpc"
	"github.com/micro/micro/v2/service/registry"
)

// DefaultClient for the service
var DefaultClient client.Client = grpc.NewClient(
	client.Registry(registry.DefaultRegistry),
)
