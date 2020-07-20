package config

import (
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	proto "github.com/micro/go-micro/v2/config/source/service/proto"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/micro/v2/service/config/handler"
)

const (
	name = "go.micro.config"
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
func Run(c *cli.Context, srvOpts ...micro.Option) {
	srvOpts = append([]micro.Option{
		micro.Name(name),
	}, srvOpts...)

	if len(c.String("watch_topic")) > 0 {
		handler.WatchTopic = c.String("watch_topic")
	}

	srv := micro.NewService(srvOpts...)

	h := &handler.Config{Store: srv.Options().Store}
	proto.RegisterConfigHandler(srv.Server(), h)
	micro.RegisterSubscriber(handler.WatchTopic, srv.Server(), handler.Watcher)

	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
