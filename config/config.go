package config

import (
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/config/cmd"
	proto "github.com/micro/go-micro/v2/config/source/service/proto"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/micro/v2/config/db"
	"github.com/micro/micro/v2/config/handler"

	// TODO: decruft to just use the store
	_ "github.com/micro/micro/v2/config/db/cockroach"
	_ "github.com/micro/micro/v2/config/db/etcd"
	_ "github.com/micro/micro/v2/config/db/memory"
	_ "github.com/micro/micro/v2/config/db/store"
)

var (
	// Service name
	Name = "go.micro.config"
	// Default database store
	Database = "store"
)

func Run(c *cli.Context, srvOpts ...micro.Option) {
	if len(c.String("server_name")) > 0 {
		Name = c.String("server_name")
	}

	if len(c.String("watch_topic")) > 0 {
		handler.WatchTopic = c.String("watch_topic")
	}

	if len(c.String("database")) > 0 {
		Database = c.String("database")
	}

	srvOpts = append(srvOpts, micro.Name(Name))

	service := micro.NewService(srvOpts...)

	proto.RegisterConfigHandler(service.Server(), new(handler.Handler))
	micro.RegisterSubscriber(handler.WatchTopic, service.Server(), handler.Watcher)

	if err := db.Init(
		db.WithDatabase(Database),
		db.WithUrl(c.String("database_url")),
		db.WithStore(*cmd.DefaultCmd.Options().Store),
	); err != nil {
		log.Fatalf("config init database error: %s", err)
	}

	if err := service.Run(); err != nil {
		log.Fatalf("config Run the service error: ", err)
	}
}

func Commands(options ...micro.Option) []*cli.Command {
	command := &cli.Command{
		Name:  "config",
		Usage: "Run the config server",
		Action: func(ctx *cli.Context) error {
			Run(ctx, options...)
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "namespace",
				EnvVars: []string{"MICRO_CONFIG_NAMESPACE"},
				Usage:   "Set the namespace used by the Config Service e.g. go.micro.srv.config",
			},
			&cli.StringFlag{
				Name:    "database_url",
				EnvVars: []string{"MICRO_CONFIG_DATABASE_URL"},
				Usage:   "The database URL e.g root:123@(127.0.0.1:3306)/config?charset=utf8&parseTime=true&loc=Asia%2FShanghai",
			},
			&cli.StringFlag{
				Name:    "database",
				EnvVars: []string{"MICRO_CONFIG_DATABASE"},
				Usage:   "The database e.g mysql(default), postgresql, but now we only support mysql and cockroach(pg).",
			},
			&cli.StringFlag{
				Name:    "watch_topic",
				EnvVars: []string{"MICRO_CONFIG_WATCH_TOPIC"},
				Usage:   "watch the change event.",
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
