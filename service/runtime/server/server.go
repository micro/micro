package server

import (
	"os"

	"github.com/micro/cli/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/runtime"
	"github.com/micro/micro/v2/service"
	"github.com/micro/micro/v2/service/auth"
	muruntime "github.com/micro/micro/v2/service/runtime"
	"github.com/micro/micro/v2/service/runtime/manager"
	pb "github.com/micro/micro/v2/service/runtime/proto"
	"github.com/micro/micro/v2/service/store"
)

var (
	// name of the runtime
	name = "go.micro.runtime"
	// address of the runtime
	address = ":8088"

	// Flags specific to the runtime service
	Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "profile",
			Usage:   "Set the runtime profile to use for services e.g local, kubernetes, platform",
			EnvVars: []string{"MICRO_RUNTIME_PROFILE"},
		},
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
	muRuntime := muruntime.DefaultRuntime
	if ctx.IsSet("source") {
		muRuntime.Init(runtime.WithSource(ctx.String("source")))
	}

	// append name
	srvOpts = append(srvOpts, service.Name(name))

	// new service
	srv := service.New(srvOpts...)

	// create a new runtime manager
	manager := manager.New(muRuntime,
		manager.Auth(auth.DefaultAuth),
		manager.Store(store.DefaultStore),
		// manager.Profile(prof),
		manager.CacheStore(store.DefaultStore),
	)

	// start the manager
	if err := manager.Start(); err != nil {
		log.Errorf("failed to start: %s", err)
		os.Exit(1)
	}

	// register the runtime handler
	pb.RegisterRuntimeHandler(srv.Server(), &Runtime{
		Event:   service.NewEvent("go.micro.runtime.events"),
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
