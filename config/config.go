package config

import (
	"sync"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/config/reader/json"
	proto "github.com/micro/go-micro/config/source/mucp/proto"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/micro/config/db"
	_ "github.com/micro/micro/config/db/cockroach"
	_ "github.com/micro/micro/config/db/mysql"
)

var (
	Name       = "go.micro.config"
	WatchTopic = "go.micro.config.events"

	PathSplitter = "/"
	Database     = "cockroach"
	// we now support json only
	reader   = json.NewReader()
	mtx      sync.RWMutex
	watchers = make(map[string][]*watcher)
)

func Run(c *cli.Context, srvOpts ...micro.Option) {
	if len(c.GlobalString("server_name")) > 0 {
		Name = c.GlobalString("server_name")
	}

	if len(c.String("watch_topic")) > 0 {
		WatchTopic = c.String("watch_topic")
	}

	if len(c.String("database")) > 0 {
		Database = c.String("database")
	}

	srvOpts = append(srvOpts, micro.Name(Name))

	service := micro.NewService(srvOpts...)
	proto.RegisterSourceHandler(service.Server(), new(Config))

	_ = service.Server().Subscribe(service.Server().NewSubscriber(WatchTopic, Watcher))

	if err := db.Init(
		db.WithDBName(Database),
		db.WithUrl(c.String("database_url")),
	)
		err != nil {
		log.Fatalf("micro config init database error: %s", err)
	}

	if err := service.Run(); err != nil {
		log.Fatalf("micro config Run the service error: ", err)
	}
}

func Commands(options ...micro.Option) []cli.Command {
	command := cli.Command{
		Name:  "config",
		Usage: "Run the config server",
		Action: func(c *cli.Context) {
			Run(c, options...)
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "namespace",
				Usage:  "Set the namespace used by the Config Service e.g. go.micro.srv.config",
				EnvVar: "MICRO_CONFIG_NAMESPACE",
			},
			cli.StringFlag{
				Name:   "database_url",
				EnvVar: "MICRO_CONFIG_DATABASE_URL",
				Usage:  "The database URL e.g root:123@(127.0.0.1:3306)/config?charset=utf8&parseTime=true&loc=Asia%2FShanghai",
			},
			cli.StringFlag{
				Name:   "database",
				EnvVar: "MICRO_CONFIG_DATABASE",
				Usage:  "The database e.g mysql(default), postgresql, but now we only support mysql and cockroach(pg).",
			},
			cli.StringFlag{
				Name:   "watch_topic",
				EnvVar: "MICRO_CONFIG_WATCH_TOPIC",
				Usage:  "watch the change event.",
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

	return []cli.Command{command}
}
