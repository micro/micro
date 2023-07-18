package server

import (
	"github.com/urfave/cli/v2"
	pb "micro.dev/v4/proto/config"
	"micro.dev/v4/service"
	"micro.dev/v4/service/config/handler"
	"micro.dev/v4/service/logger"
	"micro.dev/v4/service/store"
)

const (
	address = ":8001"
)

var (
	// Flags specific to the config service
	Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "config_secret",
			EnvVars: []string{"MICRO_CONFIG_SECRET_KEY"},
		},
	}
)

// Run micro config
func Run(c *cli.Context) error {
	srv := service.New(
		service.Name("config"),
		service.Address(address),
	)

	store.DefaultStore.Init(store.Table("config"))

	// register the handler
	pb.RegisterConfigHandler(srv.Server(), handler.NewConfig(c.String("config_secret")))

	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
	return nil
}
