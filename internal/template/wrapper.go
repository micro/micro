package template

var (
	WrapperAPI = `package client

import (
	"context"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/server"
	{{dehyphen .Alias}} "path/to/service/proto/{{.Alias}}"
)

type {{dehyphen .Alias}}Key struct {}

// FromContext retrieves the client from the Context
func {{title .Alias}}FromContext(ctx context.Context) ({{dehyphen .Alias}}.{{title .Alias}}Service, bool) {
	c, ok := ctx.Value({{dehyphen .Alias}}Key{}).({{dehyphen .Alias}}.{{title .Alias}}Service)
	return c, ok
}

// Client returns a wrapper for the {{title .Alias}}Client
func {{title .Alias}}Wrapper(service micro.Service) server.HandlerWrapper {
	client := {{dehyphen .Alias}}.New{{title .Alias}}Service("go.micro.service.template", service.Client())

	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			ctx = context.WithValue(ctx, {{dehyphen .Alias}}Key{}, client)
			return fn(ctx, req, rsp)
		}
	}
}
`
)
