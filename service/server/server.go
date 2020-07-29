package server

import (
	"github.com/micro/go-micro/v3/server"
	"github.com/micro/go-micro/v3/server/grpc"
	"github.com/micro/micro/v3/service/registry"
)

// DefaultServer for the service
var DefaultServer server.Server = grpc.NewServer(
	server.Registry(registry.DefaultRegistry),
)

// Register a handler
func Handle(hdlr server.Handler) error {
	return DefaultServer.Handle(hdlr)
}

// Create a new handler
func NewHandler(hdlr interface{}, opts ...server.HandlerOption) server.Handler {
	return DefaultServer.NewHandler(hdlr, opts...)
}

// Create a new subscriber
func NewSubscriber(topic string, hdlr interface{}, opts ...server.SubscriberOption) server.Subscriber {
	return DefaultServer.NewSubscriber(topic, hdlr, opts...)
}

// Register a subscriber
func Subscribe(sub server.Subscriber) error {
	return DefaultServer.Subscribe(sub)
}
