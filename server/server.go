package server

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/config/cmd"
	log "github.com/micro/go-micro/v2/logger"
	gorun "github.com/micro/go-micro/v2/runtime"
	handler "github.com/micro/go-micro/v2/util/file"
	"github.com/micro/micro/v2/internal/platform"
	"github.com/micro/micro/v2/internal/update"
)

var (
	// list of services managed
	services = []string{
		// runtime services
		"config",   // ????
		"network",  // :8085
		"runtime",  // :8088
		"registry", // :8000
		"broker",   // :8001
		"store",    // :8002
		"router",   // :8084
		"debug",    // :????
		"proxy",    // :8081
		"api",      // :8080
		"auth",     // :8010
		"web",      // :8082
		"bot",      // :????
		"init",     // no port, manage self
	}
)

var (
	// Name of the server microservice
	Name = "go.micro.server"
	// Address is the router microservice bind address
	Address = ":10001"
)

func Commands(options ...micro.Option) []*cli.Command {
	command := &cli.Command{
		Name:  "server",
		Usage: "Run the micro server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "address",
				Usage:   "Set the micro server address :10001",
				EnvVars: []string{"MICRO_SERVER_ADDRESS"},
			},
			&cli.BoolFlag{
				Name:  "peer",
				Usage: "Peer with the global network to share services",
			},
			&cli.StringFlag{
				Name:    "profile",
				Usage:   "Set the runtime profile to use for services e.g local, kubernetes, platform",
				EnvVars: []string{"MICRO_RUNTIME_PROFILE"},
			},
		},
		Action: func(ctx *cli.Context) error {
			Run(ctx)
			return nil
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

// Run runs the entire platform
func Run(context *cli.Context) error {
	log.Init(log.WithFields(map[string]interface{}{"service": "micro"}))

	if context.Args().Len() > 0 {
		cli.ShowSubcommandHelp(context)
		os.Exit(1)
	}
	// set default profile
	if len(context.String("profile")) == 0 {
		context.Set("profile", "server")
	}

	// get the network flag
	peer := context.Bool("peer")

	// pass through the environment
	// TODO: perhaps don't do this
	env := os.Environ()
	env = append(env, "MICRO_STORE=file")
	env = append(env, "MICRO_RUNTIME_PROFILE="+context.String("profile"))

	// connect to the network if specified
	if peer {
		log.Info("Setting global network")

		if v := os.Getenv("MICRO_NETWORK_NODES"); len(v) == 0 {
			// set the resolver to use https://micro.mu/network
			env = append(env, "MICRO_NETWORK_NODES=network.micro.mu")
			log.Info("Setting default network micro.mu")
		}
		if v := os.Getenv("MICRO_NETWORK_TOKEN"); len(v) == 0 {
			// set the network token
			env = append(env, "MICRO_NETWORK_TOKEN=micro.mu")
			log.Info("Setting default network token")
		}
	}

	log.Info("Loading core services")

	// create new micro runtime
	muRuntime := cmd.DefaultCmd.Options().Runtime

	// Use default update notifier
	if context.Bool("auto_update") {
		updateURL := context.String("update_url")
		if len(updateURL) == 0 {
			updateURL = update.DefaultURL
		}

		options := []gorun.Option{
			gorun.WithScheduler(update.NewScheduler(updateURL, platform.Version)),
		}
		(*muRuntime).Init(options...)
	}

	for _, service := range services {
		name := service

		if namespace := context.String("namespace"); len(namespace) > 0 {
			name = fmt.Sprintf("%s.%s", namespace, service)
		}

		log.Infof("Registering %s", name)

		// runtime based on environment we run the service in
		args := []gorun.CreateOption{
			gorun.WithCommand(os.Args[0]),
			gorun.WithArgs(service),
			gorun.WithEnv(env),
			gorun.WithOutput(os.Stdout),
		}

		// NOTE: we use Version right now to check for the latest release
		muService := &gorun.Service{Name: name, Version: platform.Version}
		if err := (*muRuntime).Create(muService, args...); err != nil {
			log.Errorf("Failed to create runtime enviroment: %v", err)
			return err
		}
	}

	log.Info("Starting service runtime")

	// start the runtime
	if err := (*muRuntime).Start(); err != nil {
		log.Fatal(err)
		return err
	}

	log.Info("Service runtime started")

	// TODO: should we launch the console?
	// start the console
	// cli.Init(context)

	server := micro.NewService(
		micro.Name(Name),
		micro.Address(Address),
	)

	// @todo make this configurable
	uploadDir := filepath.Join(os.TempDir(), "micro", "uploads")
	os.MkdirAll(uploadDir, 0777)
	handler.RegisterHandler(server.Server(), uploadDir)
	// start the server
	server.Run()

	log.Info("Stopping service runtime")

	// stop all the things
	if err := (*muRuntime).Stop(); err != nil {
		log.Fatal(err)
		return err
	}

	log.Info("Service runtime shutdown")

	// exit success
	os.Exit(0)
	return nil
}
