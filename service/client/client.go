package client

import (
	"context"

	"github.com/micro/go-micro/v3/client"
	"github.com/micro/go-micro/v3/client/grpc"
)

// DefaultClient for the service
var DefaultClient client.Client = grpc.NewClient()

type (
	// Message is an alias for client.Message
	Message = client.Message
	// Request is an alias for client.Request
	Request = client.Request
	// xStream is an alias for client.Stream
	xStream = client.Stream
)

// NewMessage returns a message which can be published
func NewMessage(topic string, msg interface{}, opts ...client.MessageOption) Message {
	return DefaultClient.NewMessage(topic, msg, opts...)
}

// NewRequest returns a request can which be executed using Call or Stream
func NewRequest(service, endpoint string, req interface{}, reqOpts ...client.RequestOption) Request {
	return DefaultClient.NewRequest(service, endpoint, req, reqOpts...)
}

// Call performs a request
func Call(ctx context.Context, req Request, rsp interface{}, opts ...client.CallOption) error {
	return DefaultClient.Call(ctx, req, rsp, opts...)
}

// Stream performs a streaming request
func Stream(ctx context.Context, req Request, opts ...client.CallOption) (xStream, error) {
	return DefaultClient.Stream(ctx, req, opts...)
}

// Publish a message
func Publish(ctx context.Context, msg Message, opts ...client.PublishOption) error {
	return DefaultClient.Publish(ctx, msg, opts...)
}
