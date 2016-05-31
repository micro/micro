# Plugins

Plugins are a way of integrating external code into the Micro toolkit. This is completely separate to go-micro plugins. 
Using plugins here allows you to add additional flags, commandsand HTTP handlers to the toolkit. 

## How it works

There is a global plugin manager under micro/plugin which consists of plugins that will be used across the entire toolkit. 
Plugins can be registered by calling `plugin.Register`. Each component (api, web, sidecar, cli, bot) has a separate 
plugin manager used to register plugins which should only be added as part of that component. They can be used in 
the same way by called `api.Register`, `web.Register`, etc.

Here's the interface

```go
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
        Plugins() map[string]Plugin
        Register(name string, plugin Plugin) error
}
```
