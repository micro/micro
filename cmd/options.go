package cmd

import (
	"context"

	"github.com/micro/cli/v2"

	"github.com/micro/go-micro/v3/cmd"
)

type beforeKey struct{}
type setupOnlyKey struct{}

// Before sets a function to be called before micro is setup
func Before(f cli.BeforeFunc) cmd.Option {
	return func(o *cmd.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, beforeKey{}, f)
	}
}

func beforeFromContext(ctx context.Context, def cli.BeforeFunc) cli.BeforeFunc {
	if ctx == nil {
		return def
	}

	a, ok := ctx.Value(beforeKey{}).(cli.BeforeFunc)
	if !ok {
		return def
	}

	// perform the before func passed in the context before the default
	return func(ctx *cli.Context) error {
		if err := a(ctx); err != nil {
			return err
		}
		return def(ctx)
	}
}

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
