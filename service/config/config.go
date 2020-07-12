package config

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	proto "github.com/micro/go-micro/v2/config/source/service/proto"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/micro/v2/client/cli/namespace"
	"github.com/micro/micro/v2/client/cli/util"
	"github.com/micro/micro/v2/cmd"
	"github.com/micro/micro/v2/internal/client"
	"github.com/micro/micro/v2/internal/helper"
	"github.com/micro/micro/v2/service/config/handler"
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

	srvOpts = append(srvOpts, micro.Name(Name))

	service := service.New(srvOpts...)

	h := &handler.Config{
		Store: *cmd.DefaultCmd.Options().Store,
	}

	proto.RegisterConfigHandler(service.Server(), h)
	micro.RegisterSubscriber(handler.WatchTopic, service.Server(), handler.Watcher)

	if err := service.Run(); err != nil {
		log.Fatalf("config Run the service error: ", err)
	}
}

func init() {
	// register the commands
	cmd.Commands = append(app.Commands, Commands()...)
}

func setConfig(ctx *cli.Context) error {
	pb := proto.NewConfigService("go.micro.config", client.New(ctx))

	args := ctx.Args()

	if args.Len() == 0 {
		fmt.Println("Required usage: micro config set key val")
		os.Exit(1)
	}

	// key val
	key := args.Get(0)
	val := args.Get(1)

	ns, err := namespace.Get(util.GetEnv(ctx).Name)
	if err != nil {
		return err
	}

	// TODO: allow the specifying of a config.Key. This will be service name
	// The actuall key-val set is a path e.g micro/accounts/key
	_, err = pb.Update(context.TODO(), &proto.UpdateRequest{
		Change: &proto.Change{
			// the current namespace
			Namespace: ns,
			// actual key for the value
			Path: key,
			// The value
			ChangeSet: &proto.ChangeSet{
				Data:      string(val),
				Format:    "json",
				Source:    "cli",
				Timestamp: time.Now().Unix(),
			},
		},
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return nil
}

func getConfig(ctx *cli.Context) error {
	pb := proto.NewConfigService("go.micro.config", client.New(ctx))

	args := ctx.Args()

	if args.Len() == 0 {
		fmt.Println("Required usage: micro config get key")
		os.Exit(1)
	}

	// key val
	key := args.Get(0)

	if len(key) == 0 {
		log.Fatal("key cannot be blank")
	}

	ns, err := namespace.Get(util.GetEnv(ctx).Name)
	if err != nil {
		return err
	}

	// TODO: allow the specifying of a config.Key. This will be service name
	// The actuall key-val set is a path e.g micro/accounts/key

	rsp, err := pb.Read(context.TODO(), &proto.ReadRequest{
		// The current namespace,
		Namespace: ns,
		// The actual key for the val
		Path: key,
	})

	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			fmt.Println("not found")
			os.Exit(1)
		}
		fmt.Println(err)
		os.Exit(1)
	}

	if rsp.Change == nil || rsp.Change.ChangeSet == nil {
		fmt.Println("not found")
		os.Exit(1)
	}

	// don't do it
	if v := rsp.Change.ChangeSet.Data; len(v) == 0 || string(v) == "null" {
		fmt.Println("not found")
		os.Exit(1)
	}

	fmt.Println(string(rsp.Change.ChangeSet.Data))

	return nil
}

func delConfig(ctx *cli.Context) error {
	pb := proto.NewConfigService("go.micro.config", client.New(ctx))

	args := ctx.Args()

	if args.Len() == 0 {
		fmt.Println("Required usage: micro config get key")
		os.Exit(1)
	}

	// key val
	key := args.Get(0)

	if len(key) == 0 {
		log.Fatal("key cannot be blank")
	}

	ns, err := namespace.Get(util.GetEnv(ctx).Name)
	if err != nil {
		return err
	}

	// TODO: allow the specifying of a config.Key. This will be service name
	// The actuall key-val set is a path e.g micro/accounts/key

	_, err = pb.Delete(context.TODO(), &proto.DeleteRequest{
		Change: &proto.Change{
			// The current namespace
			Namespace: ns,
			// The actual key for the val
			Path: key,
		},
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return nil
}

func Commands(options ...micro.Option) []*cli.Command {
	command := &cli.Command{
		Name:  "config",
		Usage: "Manage configuration values",
		Subcommands: []*cli.Command{
			{
				Name:   "get",
				Usage:  "Get a value; micro config get key",
				Action: getConfig,
			},
			{
				Name:   "set",
				Usage:  "Set a key-val; micro config set key val",
				Action: setConfig,
			},
			{
				Name:   "del",
				Usage:  "Delete a value; micro config del key",
				Action: delConfig,
			},
		},
		Action: func(ctx *cli.Context) error {
			if err := helper.UnexpectedSubcommand(ctx); err != nil {
				return err
			}
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
