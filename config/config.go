package config

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"time"

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
	// Default key
	Namespace = "global"
)

func setPlatform(ctx *cli.Context) {
	if ctx.Bool("platform") {
		os.Setenv("MICRO_PROXY", "service")
		os.Setenv("MICRO_PROXY_ADDRESS", "proxy.micro.mu:443")
	}
}

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

func setConfig(ctx *cli.Context) error {
	setPlatform(ctx)

	pb := proto.NewConfigService("go.micro.config", *cmd.DefaultCmd.Options().Client)

	args := ctx.Args()

	if args.Len() == 0 {
		log.Fatal("Required usage; micro config set key val")
	}

	// key val
	key := args.Get(0)
	val := args.Get(1)

	// TODO: allow the specifiying of a config.Key. This will be service name
	// The actuall key-val set is a path e.g micro/accounts/key

	_, err := pb.Update(context.TODO(), &proto.UpdateRequest{
		Change: &proto.Change{
			// global key
			Namespace: Namespace,
			// actual key for the value
			Path: key,
			// The value
			ChangeSet: &proto.ChangeSet{
				Data:      []byte(val),
				Format:    "json",
				Source:    "cli",
				Timestamp: time.Now().Unix(),
			},
		},
	})
	if err != nil {
		log.Fatalf("Error setting key-val: %v", err)
	}

	return nil
}

func getConfig(ctx *cli.Context) error {
	setPlatform(ctx)

	pb := proto.NewConfigService("go.micro.config", *cmd.DefaultCmd.Options().Client)

	args := ctx.Args()

	if args.Len() == 0 {
		log.Fatal("Required usage; micro config get key")
	}

	// key val
	key := args.Get(0)

	if len(key) == 0 {
		log.Fatal("key cannot be blank")
	}

	// TODO: allow the specifiying of a config.Key. This will be service name
	// The actuall key-val set is a path e.g micro/accounts/key

	rsp, err := pb.Read(context.TODO(), &proto.ReadRequest{
		// The global key,
		Namespace: Namespace,
		// The actual key for the val
		Path: key,
	})
	if err != nil {
		log.Fatalf("Error reading key-val: %v", err)
	}

	if rsp.Change == nil || rsp.Change.ChangeSet == nil {
		return nil
	}

	// don't do it
	if len(rsp.Change.ChangeSet.Data) == 0 {
		return nil
	}

	v, err := base64.StdEncoding.DecodeString(string(rsp.Change.ChangeSet.Data))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Value:", string(v))

	return nil
}

func delConfig(ctx *cli.Context) error {
	setPlatform(ctx)

	pb := proto.NewConfigService("go.micro.config", *cmd.DefaultCmd.Options().Client)

	args := ctx.Args()

	if args.Len() == 0 {
		log.Fatal("Required usage; micro config get key")
	}

	// key val
	key := args.Get(0)

	if len(key) == 0 {
		log.Fatal("key cannot be blank")
	}

	// TODO: allow the specifiying of a config.Key. This will be service name
	// The actuall key-val set is a path e.g micro/accounts/key

	_, err := pb.Delete(context.TODO(), &proto.DeleteRequest{
		Change: &proto.Change{
			// The global key,
			Namespace: Namespace,
			// The actual key for the val
			Path: key,
		},
	})
	if err != nil {
		log.Fatalf("Error deleting key-val: %v", err)
	}

	return nil
}

func Commands(options ...micro.Option) []*cli.Command {
	command := &cli.Command{
		Name:  "config",
		Usage: "Run the config server",
		Subcommands: []*cli.Command{
			{
				Name:   "set",
				Usage:  "Set a key-val; micro config set key val",
				Action: setConfig,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "platform",
						Usage: "Call through to the platform",
					},
				},
			},
			{
				Name:   "get",
				Usage:  "Get a value; micro config get key",
				Action: getConfig,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "platform",
						Usage: "Call through to the platform",
					},
				},
			},
			{
				Name:   "del",
				Usage:  "Delete a value; micro config del key",
				Action: delConfig,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "platform",
						Usage: "Call through to the platform",
					},
				},
			},
		},
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
