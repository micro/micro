package service

import (
	"time"

	goclient "github.com/micro/go-micro/v3/client"
	gocmd "github.com/micro/go-micro/v3/cmd"
	goserver "github.com/micro/go-micro/v3/server"

	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/server"
)

// Options for micro service
type Options struct {
	Cmd gocmd.Cmd

	Name    string
	Version string

	// Before and After funcs
	BeforeStart []func() error
	BeforeStop  []func() error
	AfterStart  []func() error
	AfterStop   []func() error

	Signal bool
}

func newOptions(opts ...Option) Options {
	opt := Options{
		Cmd:    defaultCmd,
		Signal: true,
	}

	for _, o := range opts {
		o(&opt)
	}
	return opt
}

type Option func(o *Options)

// HandleSignal toggles automatic installation of the signal handler that
// traps TERM, INT, and QUIT.  Users of this feature to disable the signal
// handler, should control liveness of the service through the context.
func HandleSignal(b bool) Option {
	return func(o *Options) {
		o.Signal = b
	}
}

// Address sets the address of the server
func Address(addr string) Option {
	return func(o *Options) {
		server.DefaultServer.Init(goserver.Address(addr))
	}
}

// Name of the service
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
		server.DefaultServer.Init(goserver.Name(n))
	}
}

// Version of the service
func Version(v string) Option {
	return func(o *Options) {
		o.Version = v
		server.DefaultServer.Init(goserver.Version(v))
	}
}

// Metadata associated with the service
func Metadata(md map[string]string) Option {
	return func(o *Options) {
		server.DefaultServer.Init(goserver.Metadata(md))
	}
}

// RegisterTTL specifies the TTL to use when registering the service
func RegisterTTL(t time.Duration) Option {
	return func(o *Options) {
		server.DefaultServer.Init(goserver.RegisterTTL(t))
	}
}

// RegisterInterval specifies the interval on which to re-register
func RegisterInterval(t time.Duration) Option {
	return func(o *Options) {
		server.DefaultServer.Init(goserver.RegisterInterval(t))
	}
}

// WrapClient is a convenience method for wrapping a Client with
// some middleware component. A list of wrappers can be provided.
// Wrappers are applied in reverse order so the last is executed first.
func WrapClient(w ...goclient.Wrapper) Option {
	return func(o *Options) {
		// apply in reverse
		for i := len(w); i > 0; i-- {
			client.DefaultClient = w[i-1](client.DefaultClient)
		}
	}
}

// WrapCall is a convenience method for wrapping a Client CallFunc
func WrapCall(w ...goclient.CallWrapper) Option {
	return func(o *Options) {
		client.DefaultClient.Init(goclient.WrapCall(w...))
	}
}

// WrapHandler adds a handler Wrapper to a list of options passed into the server
func WrapHandler(w ...goserver.HandlerWrapper) Option {
	return func(o *Options) {
		var wrappers []goserver.Option

		for _, wrap := range w {
			wrappers = append(wrappers, goserver.WrapHandler(wrap))
		}

		// Init once
		server.DefaultServer.Init(wrappers...)
	}
}

// WrapSubscriber adds a subscriber Wrapper to a list of options passed into the server
func WrapSubscriber(w ...goserver.SubscriberWrapper) Option {
	return func(o *Options) {
		var wrappers []goserver.Option

		for _, wrap := range w {
			wrappers = append(wrappers, goserver.WrapSubscriber(wrap))
		}

		// Init once
		server.DefaultServer.Init(wrappers...)
	}
}

// Before and Afters

// BeforeStart run funcs before service starts
func BeforeStart(fn func() error) Option {
	return func(o *Options) {
		o.BeforeStart = append(o.BeforeStart, fn)
	}
}

// BeforeStop run funcs before service stops
func BeforeStop(fn func() error) Option {
	return func(o *Options) {
		o.BeforeStop = append(o.BeforeStop, fn)
	}
}

// AfterStart run funcs after service starts
func AfterStart(fn func() error) Option {
	return func(o *Options) {
		o.AfterStart = append(o.AfterStart, fn)
	}
}

// AfterStop run funcs after service stops
func AfterStop(fn func() error) Option {
	return func(o *Options) {
		o.AfterStop = append(o.AfterStop, fn)
	}
}
