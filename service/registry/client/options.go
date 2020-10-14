package client

import (
	"context"

	"github.com/micro/micro/v3/service/registry"
	"github.com/micro/micro/v3/service/client"
)

type clientKey struct{}

// WithClient sets the RPC client
func WithClient(c client.Client) registry.Option {
	return func(o *registry.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}

		o.Context = context.WithValue(o.Context, clientKey{}, c)
	}
}
