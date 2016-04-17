package plugin

import (
	"net/http"

	"github.com/micro/cli"
)

// Plugin is the interface for plugins to micro.
// It differs from go-micro in that it's for
// the micro API, Web, Sidecar, CLI.
// It's a method of building middleware for the
// HTTP side.
type Plugin interface {
	// Global Flags
	Flags() []cli.Flag
	// Sub commands
	Commands() []cli.Command
	// Init called when command line args are parsed.
	// The initialised cli.Context is passed in.
	Init(*cli.Context) error
	// Handle is the middleware handler for
	// HTTP requests. We pass in the existing
	// handler so it can be wrapped to create
	// a call chain.
	Handle(http.Handler) http.Handler
}
