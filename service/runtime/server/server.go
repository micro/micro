package server

import (
	"os"

	goruntime "github.com/micro/go-micro/v3/runtime"
	"github.com/micro/go-micro/v3/runtime/builder"
	pb "github.com/micro/micro/v3/proto/runtime"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/runtime"
	builderSrv "github.com/micro/micro/v3/service/runtime/builder/client"
	"github.com/micro/micro/v3/service/runtime/manager"
	"github.com/urfave/cli/v2"
)

var (
	// name of the runtime
	name = "runtime"
	// address of the runtime
	address = ":8088"
	// builder to use
	build builder.Builder

	// Flags specific to the runtime service
	Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "source",
			Usage:   "Set the runtime source, e.g. micro/services",
			EnvVars: []string{"MICRO_RUNTIME_SOURCE"},
		},
		&cli.StringFlag{
			Name:    "builder",
			Usage:   "Set the builder, e.g service",
			EnvVars: []string{"MICRO_RUNTIME_BUILDER"},
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

	// configure the builder which is used to precompile source
	switch ctx.String("builder") {
	case "":
		// no builder is enabled, runtime will run source code directly
	case "service":
		build = builderSrv.NewBuilder()
	default:
		logger.Fatalf("Unknown builder: %v", ctx.String("builder"))
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

	// create a new runtime manager
	manager := manager.New(manager.Builder(build))

	// start the manager
	if err := manager.Start(); err != nil {
		log.Errorf("failed to start: %s", err)
		os.Exit(1)
	}

	// register the runtime handler
	pb.RegisterRuntimeHandler(srv.Server(), &Runtime{
		Runtime: manager,
	})

	// register the source handler
	pb.RegisterSourceHandler(srv.Server(), &Source{})

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
