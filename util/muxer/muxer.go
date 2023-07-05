// Package muxer provides proxy muxing
package muxer

import (
	"context"

	"github.com/micro/micro/v3/service/proxy"
	"github.com/micro/micro/v3/service/server"
	"github.com/micro/micro/v3/service/server/mucp"
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
	if msg.Topic() == s.Name {
		return s.Handler.ProcessMessage(ctx, msg)
	}
	return s.Proxy.ProcessMessage(ctx, msg)
}

func (s *Server) ServeRequest(ctx context.Context, req server.Request, rsp server.Response) error {
	if req.Service() == s.Name {
		return s.Handler.ServeRequest(ctx, req, rsp)
	}
	return s.Proxy.ServeRequest(ctx, req, rsp)
}

func New(name string, p proxy.Proxy) *Server {
	r := mucp.DefaultRouter

	return &Server{
		Name:    name,
		Proxy:   p,
		Handler: r,
	}
}
