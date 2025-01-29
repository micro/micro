package server

import (
	"os"

	pb "github.com/micro/micro/v5/proto/runtime"
	"github.com/micro/micro/v5/service"
	log "github.com/micro/micro/v5/service/logger"
	"github.com/micro/micro/v5/service/runtime/handler"
	"github.com/micro/micro/v5/service/runtime/manager"
	"github.com/urfave/cli/v2"
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
