// Package debug allows to debug services
package debug

import (
	"fmt"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/debug/log"
	dbg "github.com/micro/go-micro/debug/service"
)

const (
	// LogsUsage message for logs command
	LogsUsage = "Required usage: micro logs --name example"
)

var (
	// Name of the service
	Name = "go.micro.debug"
	// Address of the service
	Address = ":8089"
)

func getLogs(ctx *cli.Context, srvOpts ...micro.Option) {
	log.Name("debug")

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	// get the args
	name := ctx.String("name")
	since := ctx.String("since")
	count := ctx.Int("count")
	stream := ctx.Bool("stream")

	// must specify service name
	if len(name) == 0 {
		log.Fatal(LogsUsage)
	}

	service := dbg.NewDebug(name)

	var options []log.ReadOption

	d, err := time.ParseDuration(since)
	if err == nil {
		readSince := time.Now().Add(-d)
		options = append(options, log.Since(readSince))
	}

	if count > 0 {
		options = append(options, log.Count(count))
	}

	if stream {
		options = append(options, log.Stream(stream))
	}

	logs, err := service.Logs(options...)
	if err != nil {
		log.Fatal(err)
	}

	for record := range logs {
		fmt.Printf("%v\n", record)
	}
}

func getStats(ctx *cli.Context, srvOpts ...micro.Option) {
	log.Name("debug")

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	// TODO: implement this cruft
}

func run(ctx *cli.Context, srvOpts ...micro.Option) {
	log.Name("debug")

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}

	if len(ctx.GlobalString("server_name")) > 0 {
		Name = ctx.GlobalString("server_name")
	}

	if len(Address) > 0 {
		srvOpts = append(srvOpts, micro.Address(Address))
	}

	// append name
	srvOpts = append(srvOpts, micro.Name(Name))

	// new service
	service := micro.NewService(srvOpts...)

	// TODO: implement debug service for k8s cruft

	// start debug service
	if err := service.Run(); err != nil {
		log.Errorf("error running service: %v", err)
	}

	log.Infof("successfully stopped")
}

// Flags is shared flags so we don't have to continually re-add
func Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "Set the name of the service to debug",
		},
		cli.StringFlag{
			Name:  "version",
			Usage: "Set the version of the service to debug",
			Value: "latest",
		},
		cli.BoolFlag{
			Name:  "stream",
			Usage: "Set to stream logs continuously",
		},
		cli.StringFlag{
			Name:  "since",
			Usage: "Set to the relative time from which to show the logs for e.g. 1h",
		},
		cli.IntFlag{
			Name:  "count",
			Usage: "Set to query the last number of log events",
		},
	}
}

func Commands(options ...micro.Option) []cli.Command {
	command := []cli.Command{
		{
			Name:  "debug",
			Usage: "Run the micro debug service",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "address",
					Usage:  "Set the registry http address e.g 0.0.0.0:8089",
					EnvVar: "MICRO_SERVER_ADDRESS",
				},
			},
			Action: func(ctx *cli.Context) {
				run(ctx, options...)
			},
		},
		{
			Name:  "logs",
			Usage: "Get logs for a service",
			Flags: Flags(),
			Action: func(ctx *cli.Context) {
				getLogs(ctx, options...)
			},
		},
		{
			Name:  "stats",
			Usage: "Get stats for a service",
			Flags: Flags(),
			Action: func(ctx *cli.Context) {
				getStats(ctx, options...)
			},
		},
	}

	for _, p := range Plugins() {
		if cmds := p.Commands(); len(cmds) > 0 {
			command[0].Subcommands = append(command[0].Subcommands, cmds...)
		}

		if flags := p.Flags(); len(flags) > 0 {
			command[0].Flags = append(command[0].Flags, flags...)
		}
	}

	return command
}
