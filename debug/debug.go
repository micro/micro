// Package debug implements metrics/logging/introspection/... of go-micro services
package debug

import (
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/debug/log"
	"github.com/micro/go-micro/v2/debug/log/kubernetes"
	dservice "github.com/micro/go-micro/v2/debug/service"
	ulog "github.com/micro/go-micro/v2/logger"
	"github.com/micro/micro/v2/cli/util"
	logHandler "github.com/micro/micro/v2/debug/log/handler"
	pblog "github.com/micro/micro/v2/debug/log/proto"
	statshandler "github.com/micro/micro/v2/debug/stats/handler"
	pbstats "github.com/micro/micro/v2/debug/stats/proto"
	tracehandler "github.com/micro/micro/v2/debug/trace/handler"
	pbtrace "github.com/micro/micro/v2/debug/trace/proto"
	"github.com/micro/micro/v2/debug/web"
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

	// TODO: parse out --log_source
	// if kubernetes then .. go-micro/debug/log/kubernetes.New

	// default log initialiser
	newLog := func(service string) log.Log {
		// service log calls the actual service for the log
		return dservice.NewLog(
			// log with service name
			log.Name(service),
		)
	}

	source := ctx.String("log")
	switch source {
	case "service":
		newLog = func(service string) log.Log {
			return dservice.NewLog(
				log.Name(service),
			)
		}
	case "kubernetes":
		newLog = func(service string) log.Log {
			return kubernetes.NewLog(
				log.Name(service),
			)
		}
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

	// log handler
	lgHandler := &logHandler.Log{
		// create the log map
		Logs: make(map[string]log.Log),
		// Create the new func
		New: newLog,
	}

	// Register the stats handler
	pbstats.RegisterStatsHandler(service.Server(), statsHandler)
	// register trace handler
	pbtrace.RegisterTraceHandler(service.Server(), traceHandler)
	// Register the logs handler
	pblog.RegisterLogHandler(service.Server(), lgHandler)

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
				&cli.StringFlag{
					Name:    "log",
					Usage:   "Specify the log source to use e.g service, kubernetes",
					EnvVars: []string{"MICRO_DEBUG_LOG"},
					Value:   "service",
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
			Subcommands: []*cli.Command{
				&cli.Command{
					Name:  "web",
					Usage: "Start the debug web dashboard",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:    "netdata_url",
							Usage:   "The Full URL to the netdata server",
							EnvVars: []string{"MICRO_NETDATA_URL"},
							Value:   "http://localhost:19999",
						},
					},
					Action: func(c *cli.Context) error {
						web.Run(c)
						return nil
					},
				},
			},
		},
		{
			Name:  "log",
			Usage: "Get logs for a service",
			Flags: logFlags(),
			Action: func(ctx *cli.Context) error {
				getLog(ctx, options...)
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
