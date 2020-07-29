package server

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/micro/cli/v2"
	log "github.com/micro/go-micro/v3/logger"
	gorun "github.com/micro/go-micro/v3/runtime"
	"github.com/micro/go-micro/v3/util/file"
	"github.com/micro/micro/v3/client/cli/util"
	"github.com/micro/micro/v3/cmd"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/runtime"
)

var (
	// list of services managed
	services = []string{
		// runtime services
		"config",   // ????
		"auth",     // :8010
		"network",  // :8085
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

	// TODO: reimplement peering of servers e.g --peer=node1,node2,node3
	// peers are configured as network nodes to cluster between

	log.Info("Loading core services")
	for _, service := range services {
		name := service

		if namespace := context.String("namespace"); len(namespace) > 0 {
			name = fmt.Sprintf("%s.%s", namespace, service)
		}

		log.Infof("Registering %s", name)
		// @todo this is a hack
		env := []string{}
		cmdArgs := []string{}

		switch service {
		case "proxy", "web", "api", "bot", "cli":
			// pull the values we care about from environment
			for _, val := range os.Environ() {
				// only process MICRO_ values
				if !strings.HasPrefix(val, "MICRO_") {
					continue
				}
				// override any profile value because clients
				// talk to services, these may be started
				// differently in future as a `micro client`
				if strings.HasPrefix(val, "MICRO_PROFILE=") {
					val = "MICRO_PROFILE=client"
				}
				env = append(env, val)
			}
		default:
			// run server as "micro service [cmd]"
			cmdArgs = append(cmdArgs, "service")

			// pull the values we care about from environment
			for _, val := range os.Environ() {
				// only process MICRO_ values
				if !strings.HasPrefix(val, "MICRO_") {
					continue
				}
				env = append(env, val)
			}
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
			gorun.WithEnv(env),
			gorun.WithOutput(os.Stdout),
			gorun.WithRetries(10),
			gorun.CreateImage("micro/micro"),
		}

		// NOTE: we use Version right now to check for the latest release
		muService := &gorun.Service{Name: name, Version: fmt.Sprintf("%d", time.Now().Unix())}
		if err := runtime.Create(muService, args...); err != nil {
			log.Errorf("Failed to create runtime environment: %v", err)
			return err
		}
	}

	log.Info("Starting service runtime")

	// start the runtime
	if err := runtime.DefaultRuntime.Start(); err != nil {
		log.Fatal(err)
		return err
	}
	defer runtime.DefaultRuntime.Stop()
	log.Info("Service runtime started")

	// TODO: should we launch the console?
	// start the console
	// cli.Init(context)

	srv := service.New(
		service.Name(Name),
		service.Address(Address),
	)

	// @todo make this configurable
	uploadDir := filepath.Join(os.TempDir(), "micro", "uploads")
	os.MkdirAll(uploadDir, 0777)
	file.RegisterHandler(srv.Server(), uploadDir)

	// start the server
	if err := srv.Run(); err != nil {
		log.Fatalf("Error running server: %v", err)
	}

	return nil
}
