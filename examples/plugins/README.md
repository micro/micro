# Plugins

The micro toolkit supports plugins for the binary itself. These are separate from go-micro plugins.

Plugins can be used to add flags, commands and middleware handlers. An example would be authentication, 
logging, tracing, etc.

## A simple example

Here's a simple example of a plugin that adds a flag and then prints the value

```
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
