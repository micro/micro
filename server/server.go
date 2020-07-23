package server

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/micro/cli/v2"
	log "github.com/micro/go-micro/v2/logger"
	gorun "github.com/micro/go-micro/v2/runtime"
	"github.com/micro/go-micro/v2/util/file"
	"github.com/micro/micro/v2/client/cli/util"
	"github.com/micro/micro/v2/cmd"
	"github.com/micro/micro/v2/internal/update"
	"github.com/micro/micro/v2/service"
	"github.com/micro/micro/v2/service/client"
	muruntime "github.com/micro/micro/v2/service/runtime"
)

var (
	// list of services managed
	services = []string{
		// runtime services
		"config", // ????
		"auth",   // :8010
		// "network",  // :8085
		"runtime",  // :8088
		"registry", // :8000
		"broker",   // :8001
		"store",    // :8002
		"debug",    // :????
		"proxy",    // :8081
		"api",      // :8080
		"web",      // :8082
	}
)

var (
	// Name of the server microservice
	Name = "go.micro.server"
	// Address is the router microservice bind address
	Address = ":10001"
)

// upload is used for file uploads to the server
func upload(ctx *cli.Context, args []string) ([]byte, error) {
	if ctx.Args().Len() == 0 {
		return nil, errors.New("Required filename to upload")
	}

	filename := ctx.Args().Get(0)
	localfile := ctx.Args().Get(1)

	fileClient := file.New("go.micro.server", client.DefaultClient)
	return nil, fileClient.Upload(filename, localfile)
}

func init() {
	command := &cli.Command{
		Name:  "server",
		Usage: "Run the micro server",
		Description: `Launching the micro server ('micro server') will enable one to connect to it by
		setting the appropriate Micro environment (see 'micro env' && 'micro env --help') commands.`,
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
				Name:    "runtime_profile",
				Usage:   "Set the micro profile: server or platform",
				EnvVars: []string{"MICRO_RUNTIME_PROFILE"},
				Value:   "server",
			},
			&cli.BoolFlag{
				Name:    "auto_update",
				Usage:   "Enable automatic updates",
				EnvVars: []string{"MICRO_AUTO_UPDATE"},
			},
			&cli.StringFlag{
				Name:    "update_url",
				Usage:   "Set the url to retrieve system updates from",
				EnvVars: []string{"MICRO_UPDATE_URL"},
				Value:   update.DefaultURL,
			},
		},
		Action: func(ctx *cli.Context) error {
			Run(ctx)
			return nil
		},
		Subcommands: []*cli.Command{{
			Name:  "file",
			Usage: "Move files between your local machine and the server",
			Subcommands: []*cli.Command{
				{
					Name:   "upload",
					Action: util.Print(upload),
				},
			},
		}},
	}

	for _, p := range Plugins() {
		if cmds := p.Commands(); len(cmds) > 0 {
			command.Subcommands = append(command.Subcommands, cmds...)
		}

		if flags := p.Flags(); len(flags) > 0 {
			command.Flags = append(command.Flags, flags...)
		}
	}

	cmd.Register(command)
}

// Run runs the entire platform
func Run(context *cli.Context) error {
	if context.Args().Len() > 0 {
		cli.ShowSubcommandHelp(context)
		os.Exit(1)
	}

	// get the network flag
	peer := context.Bool("peer")

	// pass the env to the server
	env := os.Environ()

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
	muRuntime := muruntime.DefaultRuntime

	// Use default update notifier
	if context.Bool("auto_update") {
		updateURL := context.String("update_url")
		if len(updateURL) == 0 {
			updateURL = update.DefaultURL
		}

		options := []gorun.Option{
			gorun.WithScheduler(update.NewScheduler(updateURL, fmt.Sprintf("%d", time.Now().Unix()))),
		}
		muRuntime.Init(options...)
	}

	for _, service := range services {
		name := service

		if namespace := context.String("namespace"); len(namespace) > 0 {
			name = fmt.Sprintf("%s.%s", namespace, service)
		}

		log.Infof("Registering %s", name)
		// @todo this is a hack
		envs := env
		cmdArgs := []string{}

		switch service {
		case "proxy", "web", "api", "bot", "cli":
			envs = append(envs, "MICRO_AUTH=service")
			envs = append(envs, "MICRO_REGISTRY=service")
		default:
			// run server as "micro service [cmd]"
			cmdArgs = append(cmdArgs, "service")
			// pass the profile for the server
			envs = append(envs, "MICRO_PROFILE="+context.String("runtime_profile"))
		}

		// we want to pass through the global args so go up one level in the context lineage
		if len(context.Lineage()) > 1 {
			globCtx := context.Lineage()[1]
			for _, f := range globCtx.FlagNames() {
				cmdArgs = append(cmdArgs, "--"+f, context.String(f))
			}
		}
		cmdArgs = append(cmdArgs, service)

		// runtime based on environment we run the service in
		args := []gorun.CreateOption{
			gorun.WithCommand(os.Args[0]),
			gorun.WithArgs(cmdArgs...),
			gorun.WithEnv(envs),
			gorun.WithOutput(os.Stdout),
			gorun.WithRetries(10),
			gorun.CreateImage("micro/micro"),
		}

		// NOTE: we use Version right now to check for the latest release
		muService := &gorun.Service{Name: name, Version: fmt.Sprintf("%d", time.Now().Unix())}
		if err := muRuntime.Create(muService, args...); err != nil {
			log.Errorf("Failed to create runtime environment: %v", err)
			return err
		}
	}

	log.Info("Starting service runtime")

	// start the runtime
	if err := muRuntime.Start(); err != nil {
		log.Fatal(err)
		return err
	}

	log.Info("Service runtime started")

	// TODO: should we launch the console?
	// start the console
	// cli.Init(context)

	server := service.New(
		service.Name(Name),
		service.Address(Address),
	)

	// @todo make this configurable
	uploadDir := filepath.Join(os.TempDir(), "micro", "uploads")
	os.MkdirAll(uploadDir, 0777)
	file.RegisterHandler(server.Server(), uploadDir)

	// start the server
	if err := server.Run(); err != nil {
		log.Fatalf("Error running server: %v", err)
	}

	log.Info("Stopping service runtime")

	// stop all the things
	if err := muRuntime.Stop(); err != nil {
		log.Fatal(err)
		return err
	}

	log.Info("Service runtime shutdown")

	// exit success
	os.Exit(0)
	return nil
}
