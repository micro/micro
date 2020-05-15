package store

import (
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/config/cmd"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/store"
	pb "github.com/micro/go-micro/v2/store/service/proto"
	mcli "github.com/micro/micro/v2/client/cli"
	"github.com/micro/micro/v2/service/store/handler"
	"github.com/pkg/errors"
)

var (
	// Name of the store service
	Name = "go.micro.store"
	// Address is the store address
	Address = ":8002"
)

// run runs the micro server
func Run(ctx *cli.Context, srvOpts ...micro.Option) {
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

	// Initialise service
	service := micro.NewService(
		micro.Name(Name),
	)

	// the store handler
	storeHandler := &handler.Store{
		Default: *cmd.DefaultOptions().Store,
		Stores:  make(map[string]bool),
	}

	table := "store"
	if v := ctx.String("store_table"); len(v) > 0 {
		table = v
	}

	// set to store table
	storeHandler.Default.Init(
		store.Table(table),
	)

	backend := storeHandler.Default.String()
	options := storeHandler.Default.Options()

	log.Infof("Initialising the [%s] store with opts: %+v", backend, options)

	// set the new store initialiser
	storeHandler.New = func(database string, table string) (store.Store, error) {
		// Record the new database and table in the internal store
		if err := storeHandler.Default.Write(&store.Record{
			Key:   "databases/" + database,
			Value: []byte{},
		}, store.WriteTo("micro", "internal")); err != nil {
			return nil, errors.Wrap(err, "micro store couldn't store new database in internal table")
		}
		if err := storeHandler.Default.Write(&store.Record{
			Key:   "tables/" + database + "/" + table,
			Value: []byte{},
		}, store.WriteTo("micro", "internal")); err != nil {
			return nil, errors.Wrap(err, "micro store couldn't store new table in internal table")
		}

		return storeHandler.Default, nil
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
		},
		Action: func(ctx *cli.Context) error {
			Run(ctx, options...)
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
