package config

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/micro/cli/v2"
	proto "github.com/micro/micro/v2/service/config/proto"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/micro/v2/client/cli/namespace"
	"github.com/micro/micro/v2/client/cli/util"
	"github.com/micro/micro/v2/cmd"
	"github.com/micro/micro/v2/internal/client"
	cliconfig "github.com/micro/micro/v2/internal/config"
	"github.com/micro/micro/v2/internal/helper"
)

var (
	subcommandFlags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "local",
			Usage: "Connect to local user micro config file and not to micro server config",
		},
	}
)

func setConfig(ctx *cli.Context) error {
	args := ctx.Args()
	// key val
	key := args.Get(0)
	val := args.Get(1)

	if ctx.Bool("local") {
		return cliconfig.Set(val, strings.Split(key, ".")...)
	}
	cli, err := client.New(ctx)
	if err != nil {
		return err
	}
	pb := proto.NewConfigService("go.micro.config", cli)

	if args.Len() == 0 {
		return fmt.Errorf("Required usage: micro config set key val")
	}

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
	return err
}

func getConfig(ctx *cli.Context) error {
	args := ctx.Args()

	if args.Len() == 0 {
		return fmt.Errorf("Required usage: micro config get key")
	}
	// key val
	key := args.Get(0)
	if len(key) == 0 {
		return fmt.Errorf("key cannot be blank")
	}

	if ctx.Bool("local") {
		val, err := cliconfig.Get(strings.Split(key, ".")...)
		if err == nil {
			fmt.Println(val)
		}
		return err
	}
	cli, err := client.New(ctx)
	if err != nil {
		return err
	}
	pb := proto.NewConfigService("go.micro.config", cli)

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
		return err
	}

	if rsp.Change == nil || rsp.Change.ChangeSet == nil {
		return fmt.Errorf("not found")
	}

	// don't do it
	if v := rsp.Change.ChangeSet.Data; len(v) == 0 || string(v) == "null" {
		return fmt.Errorf("not found")
	}

	fmt.Println(string(rsp.Change.ChangeSet.Data))
	return nil
}

func delConfig(ctx *cli.Context) error {
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

	cli, err := client.New(ctx)
	if err != nil {
		return err
	}
	pb := proto.NewConfigService("go.micro.config", cli)

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
	return err
}

func init() {
	cmd.Register(
		&cli.Command{
			Name:   "config",
			Usage:  "Manage configuration values",
			Action: helper.UnexpectedSubcommand,
			Subcommands: []*cli.Command{
				{
					Name:   "get",
					Usage:  "Get a value; micro config get key",
					Action: getConfig,
					Flags:  subcommandFlags,
				},
				{
					Name:   "set",
					Usage:  "Set a key-val; micro config set key val",
					Action: setConfig,
					Flags:  subcommandFlags,
				},
				{
					Name:   "del",
					Usage:  "Delete a value; micro config del key",
					Action: delConfig,
					Flags:  subcommandFlags,
				},
			},
		},
	)
}
