package server

import (
	"os"

	"github.com/micro/cli/v2"
	goruntime "github.com/micro/go-micro/v3/runtime"
	"github.com/micro/micro/v3/service"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/runtime"
	"github.com/micro/micro/v3/service/runtime/manager"
	pb "github.com/micro/micro/v3/service/runtime/proto"
)

var (
	// name of the runtime
	name = "runtime"
	// address of the runtime
	address = ":8088"

	// Flags specific to the runtime service
	Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "source",
			Usage:   "Set the runtime source, e.g. micro/services",
			EnvVars: []string{"MICRO_RUNTIME_SOURCE"},
		},
		&cli.IntFlag{
			Name:    "retries",
			Usage:   "Set the max retries per service",
			EnvVars: []string{"MICRO_RUNTIME_RETRIES"},
		},
		&cli.IntFlag{
			Name:    "limit_cpu",
			Usage:   "Set the CPU limit for services created by the runtime, e.g. 500 for 0.5 CPU",
			EnvVars: []string{"MICRO_RUNTIME_LIMIT_CPU"},
		},
		&cli.IntFlag{
			Name:    "limit_memory",
			Usage:   "Set the memory limit for services created by the runtime, e.g. 512 for 512 MiB",
			EnvVars: []string{"MICRO_RUNTIME_LIMIT_MEMORY"},
		},
		&cli.IntFlag{
			Name:    "limit_disk",
			Usage:   "Set the disk limit for services created by the runtime, e.g. 512 for 512 MiB",
			EnvVars: []string{"MICRO_RUNTIME_LIMIT_DISK"},
		},
	}
)

// Run the runtime service
func Run(ctx *cli.Context) error {
	if len(ctx.String("address")) > 0 {
		address = ctx.String("address")
	}

	if len(ctx.String("server_name")) > 0 {
		name = ctx.String("server_name")
	}

	var srvOpts []service.Option
	if len(address) > 0 {
		srvOpts = append(srvOpts, service.Address(address))
	}

	// create runtime
	if ctx.IsSet("source") {
		runtime.DefaultRuntime.Init(goruntime.WithSource(ctx.String("source")))
	}

	// append name
	srvOpts = append(srvOpts, service.Name(name))

	// new service
	srv := service.New(srvOpts...)

	limits := &goruntime.Resources{
		CPU:  ctx.Int("limit_cpu"),
		Mem:  ctx.Int("limit_memory"),
		Disk: ctx.Int("limit_disk"),
	}

	// create a new runtime manager
	manager := manager.New(manager.WithCreateOption(goruntime.ResourceLimits(limits)))

	// start the manager
	if err := manager.Start(); err != nil {
		log.Errorf("failed to start: %s", err)
		os.Exit(1)
	}

	// register the runtime handler
	pb.RegisterRuntimeHandler(srv.Server(), &Runtime{
		Event:   service.NewEvent("runtime.events"),
		Runtime: manager,
	})

	// start runtime service
	if err := srv.Run(); err != nil {
		log.Errorf("error running service: %v", err)
	}

	// stop the manager
	if err := manager.Stop(); err != nil {
		log.Errorf("failed to stop: %s", err)
		os.Exit(1)
	}

	return nil
}
