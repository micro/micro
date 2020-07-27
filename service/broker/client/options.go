package client

import (
	"context"

	"github.com/micro/go-micro/v3/broker"
	"github.com/micro/go-micro/v3/client"
)

type clientKey struct{}

// WithClient to call broker service
func WithClient(c client.Client) broker.Option {
	return func(o *broker.Options) {
		if o.Context == nil {
			o.Context = context.WithValue(context.Background(), clientKey{}, c)
			return
		}

		o.Context = context.WithValue(o.Context, clientKey{}, c)
	}
}
