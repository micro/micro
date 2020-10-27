package cmd

import (
	"context"

	"github.com/urfave/cli/v2"
)

type Option func(o *Options)

type Options struct {
	// Name of the application
	Name string
	// Description of the application
	Description string
	// Version of the application
	Version string
	// Action to execute when Run is called and there is no subcommand
	Action func(*cli.Context) error
	// TODO replace with built in command definition
	Commands []*cli.Command
	// TODO replace with built in flags definition
	Flags []cli.Flag
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// Command line Name
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// Command line Description
func Description(d string) Option {
	return func(o *Options) {
		o.Description = d
	}
}

// Command line Version
func Version(v string) Option {
	return func(o *Options) {
		o.Version = v
	}
}

// Commands to add
func Commands(c ...*cli.Command) Option {
	return func(o *Options) {
		o.Commands = c
	}
}

// Flags to add
func Flags(f ...cli.Flag) Option {
	return func(o *Options) {
		o.Flags = f
	}
}

// Action to execute
func Action(a func(*cli.Context) error) Option {
	return func(o *Options) {
		o.Action = a
	}
}

type beforeKey struct{}
type setupOnlyKey struct{}

// Before sets a function to be called before micro is setup
func Before(f cli.BeforeFunc) Option {
	return func(o *Options) {
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

// SetupOnly for to execute
func SetupOnly() Option {
	return func(o *Options) {
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
