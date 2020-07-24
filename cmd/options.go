package cmd

import (
	"context"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/cmd"
)

type serviceKey struct{}

// Action for cmd to execute
func Action(a cli.ActionFunc) cmd.Option {
	return func(o *cmd.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, serviceKey{}, a)
	}
}

func actionFromContext(ctx context.Context) cli.ActionFunc {
	if ctx == nil {
		return nil
	}

	a, _ := ctx.Value(serviceKey{}).(cli.ActionFunc)
	return a
}
