package build

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/micro/cli/v2"
	log "github.com/micro/go-micro/v2/logger"
	goplugin "github.com/micro/go-micro/v2/plugin"
	"github.com/micro/micro/v2/plugin"
)

func build(ctx *cli.Context) error {
	name := ctx.String("name")
	path := ctx.String("path")
	newfn := ctx.String("func")
	typ := ctx.String("type")
	out := ctx.String("output")

	if len(name) == 0 {
		// TODO return err
		fmt.Println("specify --name of plugin")
		os.Exit(1)
	}

	if len(typ) == 0 {
		// TODO return err
		fmt.Println("specify --type of plugin")
		os.Exit(1)
	}

	// set the path
	if len(path) == 0 {
		// github.com/micro/go-plugins/broker/rabbitmq
		// github.com/micro/go-plugins/micro/basic_auth
		path = filepath.Join("github.com/micro/go-plugins", typ, name)
	}

	// set the newfn
	if len(newfn) == 0 {
		if typ == "micro" {
			newfn = "NewPlugin"
		} else {
			newfn = "New" + strings.Title(typ)
		}
	}

	if len(out) == 0 {
		out = "./"
	}

	// create a .so file
	if !strings.HasSuffix(out, ".so") {
		out = filepath.Join(out, name+".so")
	}

	if err := goplugin.Build(out, &goplugin.Config{
		Name:    name,
		Type:    typ,
		Path:    path,
		NewFunc: newfn,
	}); err != nil {
		// TODO return err
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Plugin %s generated at %s\n", name, out)
	return nil
}

func pluginCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:   "build",
			Usage:  "Build a micro plugin",
			Action: build,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "name",
					Usage: "Name of the plugin e.g rabbitmq",
				},
				&cli.StringFlag{
					Name:  "type",
					Usage: "Type of the plugin e.g broker",
				},
				&cli.StringFlag{
					Name:  "path",
					Usage: "Import path of the plugin",
				},
				&cli.StringFlag{
					Name:  "func",
					Usage: "New plugin function creator name e.g NewBroker",
				},
				&cli.StringFlag{
					Name:  "output, o",
					Usage: "Output dir or file for the plugin",
				},
			},
		},
	}
}

// Commands returns license commands
func Commands() []*cli.Command {
	return []*cli.Command{{
		Name:        "plugin",
		Usage:       "Plugin commands",
		Subcommands: pluginCommands(),
	}}
}

// returns a micro plugin which loads plugins
func Flags() plugin.Plugin {
	return plugin.NewPlugin(
		plugin.WithName("plugin"),
		plugin.WithFlag(
			&cli.StringSliceFlag{
				Name:    "plugin",
				EnvVars: []string{"MICRO_PLUGIN"},
				Usage:   "Comma separated list of plugins e.g broker/rabbitmq, registry/etcd, micro/basic_auth, /path/to/plugin.so",
			},
		),
		plugin.WithInit(func(ctx *cli.Context) error {
			plugins := ctx.StringSlice("plugin")
			if len(plugins) == 0 {
				return nil
			}

			for _, p := range plugins {
				if err := load(p); err != nil {
					log.Errorf("Error loading plugin %s: %v", p, err)
					return err
				}
				log.Infof("Loaded plugin %s\n", p)
			}

			return nil
		}),
	)
}
