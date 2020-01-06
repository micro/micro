package store

import (
	"strings"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/store"
	pb "github.com/micro/go-micro/store/service/proto"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/micro/store/handler"

	"github.com/micro/go-micro/store/cockroach"
	"github.com/micro/go-micro/store/memory"
)

var (
	// Name of the tunnel service
	Name = "go.micro.store"
	// Address is the tunnel address
	Address = ":8002"
	// Backend is the implementation of the store
	Backend = "memory"
	// Nodes is passed to the underlying backend
	Nodes = []string{"localhost"}
	// Namespace is passed to the underlying backend if set.
	Namespace = ""
	// Prefix is passed to the underlying backend if set.
	Prefix = ""
)

// run runs the micro server
func run(ctx *cli.Context, srvOpts ...micro.Option) {
	log.Name("store")

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	if len(ctx.GlobalString("server_name")) > 0 {
		Name = ctx.GlobalString("server_name")
	}
	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}
	if len(ctx.String("backend")) > 0 {
		Backend = ctx.String("backend")
	}
	if len(ctx.String("nodes")) > 0 {
		Nodes = strings.Split(ctx.String("nodes"), ",")
	}
	if len(ctx.String("namespace")) > 0 {
		Namespace = ctx.String("namespace")
	}

	// Initialise service
	service := micro.NewService(
		micro.Name(Name),
		micro.RegisterTTL(time.Duration(ctx.GlobalInt("register_ttl"))*time.Second),
		micro.RegisterInterval(time.Duration(ctx.GlobalInt("register_interval"))*time.Second),
	)

	opts := []store.Option{store.Nodes(Nodes...)}
	if len(Namespace) > 0 {
		opts = append(opts, store.Namespace(Namespace))
	}
	if len(Prefix) > 0 {
		opts = append(opts, store.Prefix(Prefix))
	}

	// the store handler
	storeHandler := &handler.Store{
		Stores: make(map[string]store.Store),
	}

	switch Backend {
	case "memory":
		// set the default store
		storeHandler.Default = memory.NewStore(opts...)
		// set the new store initialiser
		storeHandler.New = func(namespace string, prefix string) store.Store {
			// return a new memory store
			return memory.NewStore(
				store.Namespace(namespace),
				store.Prefix(prefix),
			)
		}
	case "cockroach":
		// set the default store
		storeHandler.Default = cockroach.NewStore(opts...)
		// set the new store initialiser
		storeHandler.New = func(namespace string, prefix string) store.Store {
			return cockroach.NewStore(
				store.Nodes(Nodes...),
				store.Namespace(namespace),
				store.Prefix(prefix),
			)
		}
	default:
		log.Fatalf("%s is not an implemented store", Backend)
	}

	pb.RegisterStoreHandler(service.Server(), storeHandler)

	// start the service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

// Commands is the cli interface for the store service
func Commands(options ...micro.Option) []cli.Command {
	command := cli.Command{
		Name:  "store",
		Usage: "Run the micro store service",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "address",
				Usage:  "Set the micro tunnel address :8002",
				EnvVar: "MICRO_SERVER_ADDRESS",
			},
			cli.StringFlag{
				Name:   "backend",
				Usage:  "Set the backend for the micro store",
				EnvVar: "MICRO_STORE_BACKEND",
				Value:  "memory",
			},
			cli.StringFlag{
				Name:   "nodes",
				Usage:  "Comma separated list of Nodes to pass to the store backend",
				EnvVar: "MICRO_STORE_NODES",
			},
			cli.StringFlag{
				Name:   "namespace",
				Usage:  "Namespace to pass to the store backend",
				EnvVar: "MICRO_STORE_NAMESPACE",
			},
			cli.StringFlag{
				Name:   "prefix",
				Usage:  "Key prefix to pass to the store backend",
				EnvVar: "MICRO_STORE_PREFIX",
			},
		},
		Action: func(ctx *cli.Context) {
			run(ctx, options...)
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
