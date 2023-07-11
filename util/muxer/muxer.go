// Package muxer provides proxy muxing
package muxer

import (
	"context"

	"micro.dev/v4/service/server"
	"micro.dev/v4/util/proxy"
)

// Server is a proxy muxer that incudes the use of the DefaultHandler
type Server struct {
	// name of service
	Name string
	// Proxy handler
	Proxy proxy.Proxy
	// The default handler
	Handler Handler
}

type Handler interface {
	proxy.Proxy
	NewHandler(interface{}, ...server.HandlerOption) server.Handler
	Handle(server.Handler) error
}

func (s *Server) ProcessMessage(ctx context.Context, msg server.Message) error {
	return s.Proxy.ProcessMessage(ctx, msg)
}

func (s *Server) ServeRequest(ctx context.Context, req server.Request, rsp server.Response) error {
	return s.Proxy.ServeRequest(ctx, req, rsp)
}

func New(name string, p proxy.Proxy) *Server {
	return &Server{
		Name:  name,
		Proxy: p,
	}
}
