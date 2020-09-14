package server

import (
	"github.com/micro/go-micro/v3/store"
	pb "github.com/micro/micro/v3/proto/config"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
	mustore "github.com/micro/micro/v3/service/store"
	"github.com/urfave/cli/v2"
)

const (
	name    = "config"
	address = ":8001"
)

var (
	// Flags specific to the config service
	Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "watch_topic",
			EnvVars: []string{"MICRO_CONFIG_WATCH_TOPIC"},
			Usage:   "watch the change event.",
		},
	}
)

// Run micro config
func Run(c *cli.Context) error {
	if len(c.String("watch_topic")) > 0 {
		watchTopic = c.String("watch_topic")
	}

	srv := service.New(
		service.Name(name),
		service.Address(address),
	)

	mustore.DefaultStore.Init(store.Table("config"))

	// register the handler
	pb.RegisterConfigHandler(srv.Server(), new(Config))
	// register the subscriber
	srv.Subscribe(watchTopic, new(watcher))

	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
	return nil
}
