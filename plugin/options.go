package plugin

import (
	"github.com/micro/cli/v2"
)

// Options are used as part of a new plugin
type Options struct {
	Name     string
	Flags    []cli.Flag
	Commands []*cli.Command
	Handlers []Handler
	Init     func(*cli.Context) error
}

type Option func(o *Options)

// WithFlag adds flags to a plugin
func WithFlag(flag ...cli.Flag) Option {
	return func(o *Options) {
		o.Flags = append(o.Flags, flag...)
	}
}

// WithCommand adds commands to a plugin
func WithCommand(cmd ...*cli.Command) Option {
	return func(o *Options) {
		o.Commands = append(o.Commands, cmd...)
	}
}

// WithHandler adds middleware handlers to
func WithHandler(h ...Handler) Option {
	return func(o *Options) {
		o.Handlers = append(o.Handlers, h...)
	}
}

// WithName defines the name of the plugin
func WithName(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// WithInit sets the init function
func WithInit(fn func(*cli.Context) error) Option {
	return func(o *Options) {
		o.Init = fn
	}
}
