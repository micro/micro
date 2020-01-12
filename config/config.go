package config

import (
	"context"
	"sync"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	cr "github.com/micro/go-micro/config/reader"
	"github.com/micro/go-micro/config/reader/json"
	"github.com/micro/go-micro/config/source"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/micro/config/db"
	proto "github.com/micro/micro/config/proto/config"
)

var (
	Namespace   = "go.micro.srv.config"
	DatabaseURL = "root:123@(127.0.0.1:3306)/config?charset=utf8&parseTime=true&loc=Asia%2FShanghai"
	DB          = "mysql"

	PathSplitter = "/"
	WatchTopic   = "micro.config.watch"

	// we now support json only
	reader   = json.NewReader()
	mtx      sync.RWMutex
	watchers = make(map[string][]*watcher)
)

func Merge(ch ...*source.ChangeSet) (*source.ChangeSet, error) {
	return reader.Merge(ch...)
}

func Values(ch *source.ChangeSet) (cr.Values, error) {
	return reader.Values(ch)
}

// Watch created by a client RPC request
func Watch(id string) (*watcher, error) {
	mtx.Lock()
	w := &watcher{
		id:   id,
		exit: make(chan bool),
		next: make(chan *proto.WatchResponse),
	}
	watchers[id] = append(watchers[id], w)
	mtx.Unlock()
	return w, nil
}

// Used as a subscriber between config services for events
func Watcher(ctx context.Context, ch *proto.WatchResponse) error {
	mtx.RLock()
	for _, sub := range watchers[ch.Id] {
		select {
		case sub.next <- ch:
		case <-time.After(time.Millisecond * 100):
		}
	}
	mtx.RUnlock()
	return nil
}

// Publish a change
func Publish(ctx context.Context, ch *proto.WatchResponse) error {
	req := client.NewMessage(WatchTopic, ch)
	return client.Publish(ctx, req)
}

func run(c *cli.Context, srvOpts ...micro.Option) {
	service := micro.NewService(srvOpts...)
	service.Init()

	if len(c.String("database_url")) > 0 {
		DatabaseURL = c.String("database_url")
	}

	if len(c.String("database")) > 0 {
		DB = c.String("database")
	}

	proto.RegisterConfigHandler(service.Server(), new(Config))

	_ = service.Server().Subscribe(service.Server().NewSubscriber(WatchTopic, Watcher))

	if err := db.Init(DB); err != nil {
		log.Fatalf("micro config init database error: %s", err)
	}

	if err := service.Run(); err != nil {
		log.Fatalf("micro config run the service error: ", err)
	}
}

func Commands(options ...micro.Option) []cli.Command {
	command := cli.Command{
		Name:  "config",
		Usage: "Run the config server",
		Action: func(c *cli.Context) {
			run(c, options...)
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "namespace",
				Usage:  "Set the namespace used by the Config Service e.g. go.micro.srv.config",
				EnvVar: "MICRO_API_NAMESPACE",
			},
			cli.StringFlag{
				Name:   "database_url",
				EnvVar: "MICRO_CONFIG_DATABASE_URL",
				Usage:  "The database URL e.g root:123@(127.0.0.1:3306)/config?charset=utf8&parseTime=true&loc=Asia%2FShanghai",
			},
			cli.StringFlag{
				Name:   "database",
				EnvVar: "MICRO_CONFIG_DATABASE",
				Usage:  "The database e.g mysql(default), postgresql, but now we only support mysql.",
				Value:  "mysql",
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
