package client

import (
	"context"

	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/store"
)

type clientKey struct{}

// WithClient to call broker service
func WithClient(c client.Client) store.Option {
	return func(o *store.Options) {
		if o.Context == nil {
			o.Context = context.WithValue(context.Background(), clientKey{}, c)
			return
		}

		o.Context = context.WithValue(o.Context, clientKey{}, c)
	}
}
