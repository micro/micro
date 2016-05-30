package template

var (
	WrapperAPI = `package client

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
	example "github.com/micro/micro/examples/template/srv/proto/example"

	"golang.org/x/net/context"
)

type exampleKey struct {}

// FromContext retrieves the client from the Context
func ExampleFromContext(ctx context.Context) (example.ExampleClient, bool) {
	c, ok := ctx.Value(exampleKey{}).(example.ExampleClient)
	return c, ok
}

// Client returns a wrapper for the ExampleClient
func ExampleWrapper(service micro.Service) server.HandlerWrapper {
	client := example.NewExampleClient("go.micro.srv.template", service.Client())

	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			ctx = context.WithValue(ctx, exampleKey{}, client)
			return fn(ctx, req, rsp)
		}
	}
}
`
)
