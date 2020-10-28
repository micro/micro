// Package server is the micro server which runs the whole system
package server

import (
	"os"
	"strings"

	"github.com/micro/micro/v3/cmd"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/auth"
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
			&cli.StringFlag{
				Name:    "image",
				Usage:   "Set the micro server image",
				EnvVars: []string{"MICRO_SERVER_IMAGE"},
				Value:   "micro/micro:latest",
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

		envvars = append(envvars, val)
	}

	// start the services
	for _, service := range services {
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

		// set the proxy addres, default to the network running locally
		if service != "network" {
			proxy := context.String("proxy_address")
			if len(proxy) == 0 {
				proxy = "127.0.0.1:8443"
			}
			env = append(env, "MICRO_PROXY="+proxy)
		}

		// for kubernetes we want to provide a port and instruct the service to bind to it. we don't do
		// this locally because the services are not isolated and the ports will conflict
		var port string
		if runtime.DefaultRuntime.String() == "kubernetes" {
			switch service {
			case "api":
				// run the api on :443, the standard port for HTTPs
				port = "443"
				env = append(env, "MICRO_API_ADDRESS=:443")
				// pass :8080 for the internal service address, since this is the default port used for the
				// static (k8s) router. Because the http api will register on :443 it won't conflict
				env = append(env, "MICRO_SERVICE_ADDRESS=:8080")
			case "proxy":
				// run the proxy on :443, the standard port for HTTPs
				port = "443"
				env = append(env, "MICRO_PROXY_ADDRESS=:443")
				// pass :8080 for the internal service address, since this is the default port used for the
				// static (k8s) router. Because the grpc proxy will register on :443 it won't conflict
				env = append(env, "MICRO_SERVICE_ADDRESS=:8080")
			case "network":
				port = "8443"
				env = append(env, "MICRO_SERVICE_ADDRESS=:8443")
			default:
				port = "8080"
				env = append(env, "MICRO_SERVICE_ADDRESS=:8080")
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
		args := []runtime.CreateOption{
			runtime.WithCommand(os.Args[0]),
			runtime.WithArgs(cmdArgs...),
			runtime.WithEnv(env),
			runtime.WithPort(port),
			runtime.WithRetries(10),
			runtime.WithServiceAccount("micro"),
			runtime.WithVolume("store-pvc", "/store"),
			runtime.CreateImage(context.String("image")),
			runtime.CreateNamespace("micro"),
			runtime.WithSecret("MICRO_AUTH_PUBLIC_KEY", auth.DefaultAuth.Options().PublicKey),
			runtime.WithSecret("MICRO_AUTH_PRIVATE_KEY", auth.DefaultAuth.Options().PrivateKey),
		}

		// NOTE: we use Version right now to check for the latest release
		muService := &runtime.Service{Name: service, Version: "latest"}
		if err := runtime.Create(muService, args...); err != nil {
			log.Errorf("Failed to create runtime environment: %v", err)
			return err
		}
	}

	// server is deployed as a pod in k8s, meaning it should exit once the services have been created.
	if runtime.DefaultRuntime.String() == "kubernetes" {
		return nil
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

	// start the server
	if err := srv.Run(); err != nil {
		log.Fatalf("Error running server: %v", err)
	}

	log.Info("Stopped server")
	return nil
}
