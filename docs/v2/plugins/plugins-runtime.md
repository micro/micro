---
title: Runtime Plugins
keywords: plugins
tags: [plugins]
sidebar: home_sidebar
permalink: /plugins-runtime
summary: 
---

Plugins are a way of integrating external code into the Micro toolkit. This is completely separate to go-micro plugins. 
Using plugins here allows you to add additional flags, commands and HTTP handlers to the toolkit. 

## How it works

There is a global plugin manager under micro/plugin which consists of plugins that will be used across the entire toolkit. 
Plugins can be registered by calling `plugin.Register`. Each component (api, web, proxy, cli, bot) has a separate 
plugin manager used to register plugins which should only be added as part of that component. They can be used in 
the same way by called `api.Register`, `web.Register`, etc.

Here's the interface

```go
// Plugin is the interface for plugins to micro. It differs from go-micro in that it's for
// the micro API, Web, Proxy, CLI. It's a method of building middleware for the HTTP side.
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
	"github.com/micro/cli"
	"github.com/micro/micro/plugin"
)

func init() {
	plugin.Register(plugin.NewPlugin(
		plugin.WithName("example"),
		plugin.WithFlag(cli.StringFlag{
			Name:   "example_flag",
			Usage:  "This is an example plugin flag",
			EnvVar: "EXAMPLE_FLAG",
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

## Repository

The plugins for the toolkit can be found in [github.com/micro/go-plugins/micro](https://github.com/micro/go-plugins/tree/master/micro).

