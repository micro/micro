package plugin

import (
	"net/http"

	"github.com/micro/cli/v2"
)

// Plugin is the interface for plugins to micro. It differs from go-micro in that it's for
// the micro API, Web, Sidecar, CLI. It's a method of building middleware for the HTTP side.
type Plugin interface {
	// Global Flags
	Flags() []cli.Flag
	// Sub-commands
	Commands() []*cli.Command
	// Handle is the middleware handler for HTTP requests. We pass in
	// the existing handler so it can be wrapped to create a call chain.
	Handler() Handler
	// Init called when command line args are parsed.
	// The initialised cli.Context is passed in.
	Init(*cli.Context) error
	// Name of the plugin
	String() string
}

// Manager is the plugin manager which stores plugins and allows them to be retrieved.
// This is used by all the components of micro.
type Manager interface {
	Plugins() []Plugin
	Register(Plugin) error
}

// Handler is the plugin middleware handler which wraps an existing http.Handler passed in.
// Its the responsibility of the Handler to call the next http.Handler in the chain.
type Handler func(http.Handler) http.Handler

type plugin struct {
	opts    Options
	init    func(ctx *cli.Context) error
	handler Handler
}

func (p *plugin) Flags() []cli.Flag {
	return p.opts.Flags
}

func (p *plugin) Commands() []*cli.Command {
	return p.opts.Commands
}

func (p *plugin) Handler() Handler {
	return p.handler
}

func (p *plugin) Init(ctx *cli.Context) error {
	return p.opts.Init(ctx)
}

func (p *plugin) String() string {
	return p.opts.Name
}

func newPlugin(opts ...Option) Plugin {
	options := Options{
		Name: "default",
		Init: func(ctx *cli.Context) error { return nil },
	}

	for _, o := range opts {
		o(&options)
	}

	handler := func(hdlr http.Handler) http.Handler {
		for _, h := range options.Handlers {
			hdlr = h(hdlr)
		}
		return hdlr
	}

	return &plugin{
		opts:    options,
		handler: handler,
	}
}

// Plugins lists the global plugins
func Plugins() []Plugin {
	return defaultManager.Plugins()
}

// Register registers a global plugins
func Register(plugin Plugin) error {
	return defaultManager.Register(plugin)
}

// NewManager creates a new plugin manager
func NewManager() Manager {
	return newManager()
}

// NewPlugin makes it easy to create a new plugin
func NewPlugin(opts ...Option) Plugin {
	return newPlugin(opts...)
}
