// Package monitor provides a monitoring service
package monitor

import (
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/monitor"
	"github.com/micro/micro/v2/monitor/handler"
	pb "github.com/micro/micro/v2/monitor/proto"
)

var (
	Name = "go.micro.monitor"
)

func run(ctx *cli.Context, opts ...micro.Option) {
	log.Info("monitor")

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
		log.Infof("Monitoring service: %s\n", serviceName)
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

func Commands(options ...micro.Option) []*cli.Command {
	command := &cli.Command{
		Name:  "monitor",
		Usage: "Run the monitoring service",
		Action: func(ctx *cli.Context) error {
			run(ctx, options...)
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "service",
				Usage:   "Name of the micro service to monitor",
				EnvVars: []string{"MICRO_SERVICE_NAME"},
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

	return []*cli.Command{command}
}
