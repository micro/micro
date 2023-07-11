package server

import (
	"context"
	"net/http"

	"micro.dev/v4/service/api"
	"micro.dev/v4/service/api/handler"
	"micro.dev/v4/service/api/handler/rpc"
	"micro.dev/v4/service/api/router"
	"micro.dev/v4/service/client"
	"micro.dev/v4/service/errors"
)

type metaHandler struct {
	c client.Client
	r router.Router
}

var (
	// built in handlers
	handlers = map[string]handler.Handler{
		"rpc": rpc.NewHandler(),
	}
)

// Register a handler
func Register(handler string, hd handler.Handler) {
	handlers[handler] = hd
}

// serverContext
type serverContext struct {
	context.Context
	domain  string
	client  client.Client
	service *api.Service
}

func (c *serverContext) Service() *api.Service {
	return c.service
}

func (c *serverContext) Client() client.Client {
	return c.client
}

func (c *serverContext) Domain() string {
	return c.domain
}

func (m *metaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	service, err := m.r.Route(r)
	if err != nil {
		er := errors.InternalServerError("micro", err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte(er.Error()))
		return
	}

	// inject service into context
	ctx := r.Context()
	// create a new server context
	srvContext := &serverContext{
		Context: ctx,
		domain:  service.Domain,
		client:  m.c,
		service: service,
	}
	// clone request with new context
	req := r.Clone(srvContext)

	// serve request
	handlers["rpc"].ServeHTTP(w, req)
}

// Meta is a http.Handler that routes based on endpoint metadata
func Meta(c client.Client, r router.Router) http.Handler {
	return &metaHandler{
		c: c,
		r: r,
	}
}
