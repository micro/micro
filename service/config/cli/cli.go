package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	goclient "github.com/micro/go-micro/v3/client"
	"github.com/micro/micro/v3/client/cli/namespace"
	"github.com/micro/micro/v3/client/cli/util"
	"github.com/micro/micro/v3/cmd"
	"github.com/micro/micro/v3/internal/helper"
	proto "github.com/micro/micro/v3/proto/config"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/context"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/urfave/cli/v2"
)

func setConfig(ctx *cli.Context) error {
	args := ctx.Args()
	// key val
	key := args.Get(0)
	val := args.Get(1)

	pb := proto.NewConfigService("config", client.DefaultClient)

	if args.Len() == 0 {
		return fmt.Errorf("Required usage: micro config set key val")
	}

	ns, err := namespace.Get(util.GetEnv(ctx).Name)
	if err != nil {
		return err
	}

	v, _ := json.Marshal(parseValue(val))

	// TODO: allow the specifying of a config.Key. This will be service name
	// The actuall key-val set is a path e.g micro/accounts/key
	_, err = pb.Set(context.DefaultContext, &proto.SetRequest{
		// the current namespace
		Namespace: ns,
		// actual key for the value
		Path: key,
		// The value
		Value: &proto.Value{
			Data: string(v),
			//Format: "json",
		},
		Options: &proto.Options{
			Secret: ctx.Bool("secret"),
		},
	}, goclient.WithAuthToken())
	return err
}

func parseValue(s string) interface{} {
	b, err := strconv.ParseBool(s)
	if err == nil {
		return b
	}
	f, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return f
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		return i
	}
	return s
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

	ns, err := namespace.Get(util.GetEnv(ctx).Name)
	if err != nil {
		return err
	}

	// TODO: allow the specifying of a config.Key. This will be service name
	// The actuall key-val set is a path e.g micro/accounts/key
	pb := proto.NewConfigService("config", client.DefaultClient)
	rsp, err := pb.Get(context.DefaultContext, &proto.GetRequest{
		// The current namespace,
		Namespace: ns,
		// The actual key for the val
		Path: key,
		Options: &proto.Options{
			Secret: ctx.Bool("secret"),
		},
	}, goclient.WithAuthToken())
	if err != nil {
		return err
	}

	if v := rsp.Value.Data; len(v) == 0 || strings.TrimSpace(string(v)) == "null" {
		return fmt.Errorf("not found")
	}

	if strings.HasPrefix(rsp.Value.Data, "\"") && strings.HasSuffix(rsp.Value.Data, "\"") {
		fmt.Println(rsp.Value.Data[1 : len(rsp.Value.Data)-1])
		return nil
	}
	fmt.Println(string(rsp.Value.Data))
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

	ns, err := namespace.Get(util.GetEnv(ctx).Name)
	if err != nil {
		return err
	}

	// TODO: allow the specifying of a config.Key. This will be service name
	// The actuall key-val set is a path e.g micro/accounts/key
	pb := proto.NewConfigService("config", client.DefaultClient)
	_, err = pb.Delete(context.DefaultContext, &proto.DeleteRequest{
		// The current namespace
		Namespace: ns,
		// The actual key for the val
		Path: key,
	}, goclient.WithAuthToken())
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
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:    "secret",
							Aliases: []string{"s"},
							Usage:   "Set it as a secret value",
						},
					},
				},
				{
					Name:   "set",
					Usage:  "Set a key-val; micro config set key val",
					Action: setConfig,
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:    "secret",
							Aliases: []string{"s"},
							Usage:   "Set it as a secret value",
						},
					},
				},
				{
					Name:   "del",
					Usage:  "Delete a value; micro config del key",
					Action: delConfig,
				},
			},
		},
	)
}
