package plugin

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/micro/cli"
)

// Plugin is the interface for plugins to micro. It differs from go-micro in that it's for
// the micro API, Web, Sidecar, CLI. It's a method of building middleware for the HTTP side.
type Plugin interface {
	// Global Flags
	Flags() []cli.Flag
	// Sub-commands
	Commands() []cli.Command
	// Init called when command line args are parsed.
	// The initialised cli.Context is passed in.
	Init(*cli.Context) error
	// Handle is the middleware handler for HTTP requests. We pass in
	// the existing handler so it can be wrapped to create a call chain.
	Handle(http.Handler) http.Handler
	// Name of the plugin
	String() string
}

// Manager is the plugin manager which stores plugins and allows them to be retrieved.
// This is used by all the components of micro.
type Manager interface {
	Plugins() []Plugin
	Register(Plugin) error
}

type manager struct {
	sync.Mutex
	plugins    []Plugin
	registered map[string]bool
}

var (
	// global plugin manager
	defaultManager = newManager()
)

func newManager() *manager {
	return &manager{
		registered: make(map[string]bool),
	}
}

func (m *manager) Plugins() []Plugin {
	m.Lock()
	defer m.Unlock()
	return m.plugins
}

func (m *manager) Register(plugin Plugin) error {
	m.Lock()
	defer m.Unlock()

	name := plugin.String()

	if m.registered[name] {
		return fmt.Errorf("Plugin with name %s already registered", name)
	}

	m.registered[name] = true
	m.plugins = append(m.plugins, plugin)
	return nil
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
