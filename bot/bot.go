package bot

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/micro/cli"
	"github.com/micro/micro/bot/plugin"
	_ "github.com/micro/micro/bot/plugin/slack"
)

func run(ctx *cli.Context) {
	if len(ctx.String("plugins")) == 0 {
		fmt.Println("No plugins specified")
		os.Exit(1)
	}

	plugins := strings.Split(ctx.String("plugins"), ",")
	if len(plugins) == 0 {
		fmt.Println("No plugins specified")
		os.Exit(1)
	}

	for _, p := range plugins {
		if _, ok := plugin.Plugins[p]; !ok {
			fmt.Printf("Plugin %s not found\n", p)
			os.Exit(1)
		}
	}

	for _, p := range plugins {
		pg := plugin.Plugins[p]

		fmt.Println("Starting plugin", p)

		if err := pg.Init(ctx); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := pg.Start(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	fmt.Println(<-ch)

	for _, p := range plugins {
		pg := plugin.Plugins[p]
		fmt.Println("Stopping plugin", pg)
		if err := pg.Stop(); err != nil {
			fmt.Println(err)
		}
	}
}

func Commands() []cli.Command {
	flags := []cli.Flag{
		cli.StringFlag{
			Name:  "plugins",
			Usage: "Plugins to load on startup",
		},
	}

	// setup plugin flags
	for _, plugin := range plugin.Plugins {
		flags = append(flags, plugin.Flags()...)
	}

	return []cli.Command{
		{
			Name:   "bot",
			Usage:  "Run the micro bot",
			Flags:  flags,
			Action: run,
		},
	}
}
