package server

import (
	"github.com/micro/go-micro/v3/server"
	"github.com/micro/go-micro/v3/server/grpc"
)

// DefaultServer for the service
var DefaultServer server.Server = grpc.NewServer()

// Register a handler
func Handle(hdlr server.Handler) error {
	return DefaultServer.Handle(hdlr)
}

// Create a new handler
func NewHandler(hdlr interface{}) server.Handler {
	return DefaultServer.NewHandler(hdlr)
}

// Create a new subscriber
func NewSubscriber(topic string, hdlr interface{}) server.Subscriber {
	return DefaultServer.NewSubscriber(topic, hdlr)
}

// Register a subscriber
func Subscribe(sub server.Subscriber) error {
	return DefaultServer.Subscribe(sub)
}
