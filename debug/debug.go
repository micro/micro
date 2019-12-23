// Package debug implements metrics/logging/introspection/... of go-micro services
package debug

import (
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/debug/log"
	"github.com/micro/go-micro/debug/log/kubernetes"
	dservice "github.com/micro/go-micro/debug/service"
	ulog "github.com/micro/go-micro/util/log"
	logHandler "github.com/micro/micro/debug/log/handler"
	pblog "github.com/micro/micro/debug/log/proto"
	"github.com/micro/micro/debug/stats"
	statshandler "github.com/micro/micro/debug/stats/handler"
	pbstats "github.com/micro/micro/debug/stats/proto"
	"github.com/micro/micro/debug/web"
)

var (
	// Name of the service
	Name = "go.micro.debug"
	// Address of the service
	Address = ":8089"
)

func Run(ctx *cli.Context, srvOpts ...micro.Option) {
	ulog.Name("debug")

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

	// stats handler
	statsHandler, err := statshandler.New(done)
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

	// Register the logs handler
	pblog.RegisterLogHandler(service.Server(), lgHandler)

	// TODO: implement debug service for k8s cruft

	// start debug service
	if err := service.Run(); err != nil {
		ulog.Fatal(err)
	}
}

// Commands populates the debug commands
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
				cli.StringFlag{
					Name:   "log",
					Usage:  "Specify the log source to use e.g service, kubernetes",
					EnvVar: "MICRO_DEBUG_LOG",
					Value:  "service",
				},
			},
			Action: func(ctx *cli.Context) {
				Run(ctx, options...)
			},
			Subcommands: []cli.Command{
				cli.Command{
					Name:  "web",
					Usage: "Start the debug web dashboard",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "netdata_url",
							Usage:  "The Full URL to the netdata server",
							EnvVar: "MICRO_NETDATA_URL",
							Value:  "http://localhost:19999",
						},
					},
					Action: func(c *cli.Context) {
						web.Run(c)
					},
				},
				cli.Command{
					Name:  "stats",
					Usage: "Start the debug stats scraper",
					Action: func(c *cli.Context) {
						stats.Run(c)
					},
				},
			},
		},
		{
			Name:  "log",
			Usage: "Get logs for a service",
			Flags: logFlags(),
			Action: func(ctx *cli.Context) {
				getLog(ctx, options...)
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
