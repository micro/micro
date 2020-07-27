package client

import (
	"github.com/micro/go-micro/v3/client"
	"github.com/micro/go-micro/v3/client/grpc"
)

// DefaultClient for the service
var DefaultClient client.Client = grpc.NewClient()
