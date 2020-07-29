package cmd

import (
	"context"

	"github.com/micro/go-micro/v3/cmd"
)

type setupOnlyKey struct{}

// SetupOnly for cmd to execute
func SetupOnly() cmd.Option {
	return func(o *cmd.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, setupOnlyKey{}, true)
	}
}

func setupOnlyFromContext(ctx context.Context) bool {
	if ctx == nil {
		return false
	}

	a, _ := ctx.Value(setupOnlyKey{}).(bool)
	return a
}
