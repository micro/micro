package server

import (
	pb "github.com/micro/micro/v5/proto/config"
	"github.com/micro/micro/v5/service"
	"github.com/micro/micro/v5/service/config/handler"
	"github.com/micro/micro/v5/service/logger"
	"github.com/micro/micro/v5/service/store"
	"github.com/urfave/cli/v2"
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
