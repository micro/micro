// Package debug implements metrics/logging/introspection/... of go-micro services
package debug

import (
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	ulog "github.com/micro/go-micro/v2/logger"
	cliutil "github.com/micro/micro/v2/cli/util"
	statshandler "github.com/micro/micro/v2/debug/stats/handler"
	pbstats "github.com/micro/micro/v2/debug/stats/proto"
	tracehandler "github.com/micro/micro/v2/debug/trace/handler"
	pbtrace "github.com/micro/micro/v2/debug/trace/proto"
)

var (
	// Name of the service
	Name = "go.micro.debug"
	// Address of the service
	Address = ":8089"
)

func Run(ctx *cli.Context, srvOpts ...micro.Option) {
	ulog.Init(ulog.WithFields(map[string]interface{}{"service": "debug"}))

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}

	if len(ctx.String("server_name")) > 0 {
		Name = ctx.String("server_name")
	}

	if len(Address) > 0 {
		srvOpts = append(srvOpts, micro.Address(Address))
	}

	// append name
	srvOpts = append(srvOpts, micro.Name(Name))

	// new service
	service := micro.NewService(srvOpts...)

	done := make(chan bool)
	defer func() {
		close(done)
	}()

	// create a service cache
	c := newCache(done)

	// stats handler
	statsHandler, err := statshandler.New(done, ctx.Int("window"), c.services)
	if err != nil {
		ulog.Fatal(err)
	}

	// stats handler
	traceHandler, err := tracehandler.New(done, ctx.Int("window"), c.services)
	if err != nil {
		ulog.Fatal(err)
	}

	// Register the stats handler
	pbstats.RegisterStatsHandler(service.Server(), statsHandler)
	// register trace handler
	pbtrace.RegisterTraceHandler(service.Server(), traceHandler)

	// TODO: implement debug service for k8s cruft

	// start debug service
	if err := service.Run(); err != nil {
		ulog.Fatal(err)
	}
}

// Commands populates the debug commands
func Commands(options ...micro.Option) []*cli.Command {
	cliutil.SetupCommand()
	command := []*cli.Command{
		{
			Name:  "debug",
			Usage: "Run the micro debug service",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "address",
					Usage:   "Set the registry http address e.g 0.0.0.0:8089",
					EnvVars: []string{"MICRO_SERVER_ADDRESS"},
				},
				&cli.IntFlag{
					Name:    "window",
					Usage:   "Specifies how many seconds of stats snapshots to retain in memory",
					EnvVars: []string{"MICRO_DEBUG_WINDOW"},
					Value:   360,
				},
			},
			Action: func(ctx *cli.Context) error {
				Run(ctx, options...)
				return nil
			},
		},
		{
			Name:  "trace",
			Usage: "Get tracing info from a service",
			Action: func(ctx *cli.Context) error {
				getTrace(ctx, options...)
				return nil
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
