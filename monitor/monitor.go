// Package monitor provides a monitoring service
package monitor

import (
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/monitor"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/micro/monitor/handler"
	pb "github.com/micro/micro/monitor/proto"
)

var (
	Name = "go.micro.monitor"
)

func run(ctx *cli.Context, opts ...micro.Option) {
	log.Name("monitor")

	// create a new monitor
	m := monitor.NewMonitor()
	if err := m.Run(); err != nil {
		log.Fatalf("Failed to start monitoring: %v", err)
	}

	defer m.Stop()

	// check if we need to monitor a service
	serviceName := ctx.String("service")

	if len(serviceName) > 0 {
		if err := m.Watch(serviceName); err != nil {
			log.Fatalf("Failed to monitor %s: %v\n", serviceName, err)
		}
		log.Logf("Monitoring service: %s\n", serviceName)
	}

	// new service
	service := micro.NewService(
		micro.Name(Name),
	)

	// register monitoring handler
	pb.RegisterMonitorHandler(service.Server(), &handler.Monitor{Monitor: m})

	// run service
	service.Run()
}

func Commands(options ...micro.Option) []cli.Command {
	command := cli.Command{
		Name:  "monitor",
		Usage: "Run the monitoring service",
		Action: func(ctx *cli.Context) {
			run(ctx, options...)
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "service",
				Usage:  "Name of the micro service to monitor",
				EnvVar: "MICRO_SERVICE_NAME",
			},
		},
	}

	for _, p := range Plugins() {
		if cmds := p.Commands(); len(cmds) > 0 {
			command.Subcommands = append(command.Subcommands, cmds...)
		}

		if flags := p.Flags(); len(flags) > 0 {
			command.Flags = append(command.Flags, flags...)
		}
	}

	return []cli.Command{command}
}
