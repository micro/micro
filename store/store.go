package store

import (
	"strings"
	"time"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/store"
	pb "github.com/micro/go-micro/v2/store/service/proto"
	mcli "github.com/micro/micro/v2/cli"
	"github.com/micro/micro/v2/store/handler"

	"github.com/micro/go-micro/v2/store/cockroach"
	"github.com/micro/go-micro/v2/store/memory"
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
	log.Init(log.WithFields(map[string]interface{}{"service": "store"}))

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	if len(ctx.String("server_name")) > 0 {
		Name = ctx.String("server_name")
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
	if len(ctx.String("database")) > 0 {
		Namespace = ctx.String("database")
	}
	if len(ctx.String("table")) > 0 {
		Prefix = ctx.String("table")
	}

	// Initialise service
	service := micro.NewService(
		micro.Name(Name),
		micro.RegisterTTL(time.Duration(ctx.Int("register_ttl"))*time.Second),
		micro.RegisterInterval(time.Duration(ctx.Int("register_interval"))*time.Second),
	)

	opts := []store.Option{store.Nodes(Nodes...)}
	if len(Namespace) > 0 {
		opts = append(opts, store.Database(Namespace))
	}
	if len(Prefix) > 0 {
		opts = append(opts, store.Table(Prefix))
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
		storeHandler.New = func(namespace string, prefix string) (store.Store, error) {
			// return a new memory store
			return memory.NewStore(
				store.Database(namespace),
				store.Table(prefix),
			), nil
		}
	case "cockroach":
		// set the default store
		storeHandler.Default = cockroach.NewStore(opts...)
		// set the new store initialiser
		storeHandler.New = func(namespace string, prefix string) (store.Store, error) {
			storeDB := cockroach.NewStore(
				store.Nodes(Nodes...),
				store.Database(namespace),
				store.Table(prefix),
			)
			if err := storeDB.Init(); err != nil {
				return nil, err
			}
			return storeDB, nil
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
func Commands(options ...micro.Option) []*cli.Command {
	command := &cli.Command{
		Name:  "store",
		Usage: "Run the micro store service",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "address",
				Usage:   "Set the micro tunnel address :8002",
				EnvVars: []string{"MICRO_SERVER_ADDRESS"},
			},
			&cli.StringFlag{
				Name:    "backend",
				Usage:   "Set the backend for the micro store",
				EnvVars: []string{"MICRO_STORE_BACKEND"},
				Value:   "memory",
			},
			&cli.StringFlag{
				Name:    "nodes",
				Usage:   "Comma separated list of Nodes to pass to the store backend",
				EnvVars: []string{"MICRO_STORE_NODES"},
			},
			&cli.StringFlag{
				Name:    "database",
				Usage:   "Database option to pass to the store backend",
				EnvVars: []string{"MICRO_STORE_DATABASE"},
			},
			&cli.StringFlag{
				Name:    "table",
				Usage:   "Table option to pass to the store backend",
				EnvVars: []string{"MICRO_STORE_TABLE"},
			},
		},
		Action: func(ctx *cli.Context) error {
			run(ctx, options...)
			return nil
		},
		Subcommands: mcli.StoreCommands(),
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
