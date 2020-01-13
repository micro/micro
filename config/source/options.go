package mucp

import (
	"context"

	"github.com/micro/go-micro/config/source"
)

type serviceNameKey struct{}
type idKey struct{}

func ServiceName(a string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, serviceNameKey{}, a)
	}
}

// Id sets the key prefix to use
func Id(p string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, idKey{}, p)
	}
}
