// Package server is the micro server which runs the whole system
package server

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/micro/go-micro/v3/util/file"
	"github.com/micro/micro/v3/client/cli/util"
	"github.com/micro/micro/v3/cmd"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/context"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/runtime"
	"github.com/urfave/cli/v2"
)

var (
	// list of services managed
	services = []string{
		"network",  // :8443
		"runtime",  // :8088
		"registry", // :8000
		"config",   // :8001
		"store",    // :8002
		"broker",   // :8003
		"events",   // :unset
		"auth",     // :8010
		"proxy",    // :8081
		"api",      // :8080
	}
)

var (
	// Name of the server microservice
	Name = "server"
	// Address is the server address
	Address = ":10001"
)

// upload is used for file uploads to the server
func upload(ctx *cli.Context, args []string) ([]byte, error) {
	if ctx.Args().Len() == 0 {
		return nil, errors.New("Required filename to upload")
	}

	filename := ctx.Args().Get(0)
	localfile := ctx.Args().Get(1)

	fileClient := file.New("server", client.DefaultClient, file.WithContext(context.DefaultContext))
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

	log.Info("Starting server")

	// parse the env vars
	var envvars []string
	for _, val := range os.Environ() {
		comps := strings.Split(val, "=")
		if len(comps) != 2 {
			continue
		}

		// only process MICRO_ values
		if !strings.HasPrefix(comps[0], "MICRO_") {
			continue
		}

		// skip the profile and proxy, that's set below since it can be service specific
		if comps[0] == "MICRO_PROFILE" || comps[0] == "MICRO_PROXY" {
			continue
		}

		// don't pass keys (e.g. MICRO_AUTH_PRIVATE_KEY) as an env var, these should be mounted
		// as secrets
		if strings.HasSuffix(comps[0], "_KEY") {
			continue
		}

		envvars = append(envvars, val)
	}

	// bind to port 8080, this is what the k8s service will use (enabling routing) and the readiness
	// checks will use. this logic exists in the runtime manager but since that's not used by the
	// server, the logic had to be duplicated here
	if runtime.DefaultRuntime.String() == "kubernetes" {
		envvars = append(envvars, "MICRO_SERVICE_ADDRESS=:8080")
	}

	// start the services
	for _, service := range services {
		// set the proxy addres, default to the network running locally
		proxy := context.String("proxy_address")
		if len(proxy) == 0 {
			proxy = "127.0.0.1:8443"
		}

		log.Infof("Registering %s", service)

		// all things run by the server are `micro service [name]`
		cmdArgs := []string{"service"}

		// override the profile for api & proxy
		env := envvars
		if service == "proxy" || service == "api" {
			env = append(env, "MICRO_PROFILE=client")
		} else {
			env = append(env, "MICRO_PROFILE="+context.String("profile"))
		}

		// inject the proxy address for all services but the network, as we don't want
		// that calling itself
		if len(proxy) > 0 && service != "network" {
			env = append(env, "MICRO_PROXY="+proxy)
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
		args := []runtime.CreateOption{
			runtime.WithCommand(os.Args[0]),
			runtime.WithArgs(cmdArgs...),
			runtime.WithEnv(env),
			runtime.WithRetries(10),
			runtime.WithVolume("store-pvc", "/store"),
			runtime.CreateImage("localhost:5000/micro"),
		}

		// NOTE: we use Version right now to check for the latest release
		muService := &runtime.Service{Name: service, Version: "latest"}
		if err := runtime.Create(muService, args...); err != nil {
			log.Errorf("Failed to create runtime environment: %v", err)
			return err
		}
	}

	log.Info("Starting server runtime")

	// start the runtime
	if err := runtime.DefaultRuntime.Start(); err != nil {
		log.Fatal(err)
		return err
	}
	defer runtime.DefaultRuntime.Stop()

	// internal server
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

	log.Info("Stopped server")

	return nil
}
