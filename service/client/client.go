package client

import (
	"context"

	"github.com/micro/go-micro/v3/client"
	"github.com/micro/go-micro/v3/client/grpc"
)

// DefaultClient for the service
var DefaultClient client.Client = grpc.NewClient()

type (
	// Connection is an alias for client.Stream
	Connection = client.Stream
	// Message is an alias for client.Message
	Message = client.Message
	// Request is an alias for client.Request
	Request = client.Request
)

// NewMessage returns a message which can be published
func NewMessage(topic string, msg interface{}) Message {
	return DefaultClient.NewMessage(topic, msg)
}

// NewRequest returns a request can which be executed using Call or Stream
func NewRequest(service, endpoint string, req interface{}) Request {
	return DefaultClient.NewRequest(service, endpoint, req)
}

// Call performs a request
func Call(ctx context.Context, req Request, rsp interface{}) error {
	return DefaultClient.Call(ctx, req, rsp)
}

// Stream performs a streaming request
func Stream(ctx context.Context, req Request) (Connection, error) {
	return DefaultClient.Stream(ctx, req)
}

// Publish a message
func Publish(ctx context.Context, msg Message) error {
	return DefaultClient.Publish(ctx, msg)
}
