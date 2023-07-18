package server

import (
	"os"

	"github.com/urfave/cli/v2"
	pb "micro.dev/v4/proto/runtime"
	"micro.dev/v4/service"
	log "micro.dev/v4/service/logger"
	"micro.dev/v4/service/runtime/handler"
	"micro.dev/v4/service/runtime/manager"
)

var (
	// address of the runtime
	address = ":8088"
)

// Run the runtime service
func Run(ctx *cli.Context) error {
	// new service
	srv := service.New(
		service.Name("runtime"),
		service.Address(address),
	)

	// create a new runtime manager
	manager := manager.New()

	// start the manager
	if err := manager.Start(); err != nil {
		log.Errorf("failed to start: %s", err)
		os.Exit(1)
	}

	// register the handlers
	pb.RegisterRuntimeHandler(srv.Server(), &handler.Runtime{Runtime: manager})
	pb.RegisterBuildHandler(srv.Server(), new(handler.Build))
	pb.RegisterSourceHandler(srv.Server(), new(handler.Source))

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
