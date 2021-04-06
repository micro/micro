# Plugins

Plugins are a way of extending the functionality of Micro

## Overview

Plugins enable Micro to be extending and intercepted to provide additional functionality and features. 
This may include logging, metrics, tracing, authentication, etc. The Plugin model requires registering 
a struct that matches a plugin interface. It's then registered and setup when Micro starts.

## Design

Here's the interface design

```go
// Plugin is the interface for plugins to micro. It differs from go-micro in that it's for
// the micro API, Web, Sidecar, CLI. It's a method of building middleware for the HTTP side.
type Plugin interface {
	// Global Flags
	Flags() []cli.Flag
	// Sub-commands
	Commands() []cli.Command
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
        Plugins() map[string]Plugin
        Register(name string, plugin Plugin) error
}

// Handler is the plugin middleware handler which wraps an existing http.Handler passed in.
// Its the responsibility of the Handler to call the next http.Handler in the chain.
type Handler func(http.Handler) http.Handler
```

## How to use it

Here's a simple example of a plugin that adds a flag and then prints the value

### The plugin

Create a plugin.go file in the top level dir

```go
package main

import (
	"log"
	"github.com/urfave/cli/v2"
	"github.com/micro/micro/plugin"
)

func init() {
	plugin.Register(plugin.NewPlugin(
		plugin.WithName("example"),
		plugin.WithFlag(&cli.StringFlag{
			Name:   "example_flag",
			Usage:  "This is an example plugin flag",
			EnvVars: []string{"EXAMPLE_FLAG"},
			Value: "avalue",
		}),
		plugin.WithInit(func(ctx *cli.Context) error {
			log.Println("Got value for example_flag", ctx.String("example_flag"))
			return nil
		}),
	))
}
```

### Building the code

Simply build micro with the plugin

```shell
go build -o micro ./main.go ./plugin.go
```

