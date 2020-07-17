// Package runtime is the micro runtime
package runtime

import (
	"os"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/cmd"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/runtime"
	pb "github.com/micro/go-micro/v2/runtime/service/proto"
	"github.com/micro/micro/v2/service/runtime/handler"
	"github.com/micro/micro/v2/service/runtime/manager"
	"github.com/micro/micro/v2/service/runtime/profile"
)

var (
	// Name of the runtime
	Name = "go.micro.runtime"
	// Address of the runtime
	Address = ":8088"
	// Flags specific to the runtime service
	Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "profile",
			Usage:   "Set the runtime profile to use for services e.g local, kubernetes, platform",
			EnvVars: []string{"MICRO_RUNTIME_PROFILE"},
		},
		&cli.StringFlag{
			Name:    "source",
			Usage:   "Set the runtime source, e.g. micro/services",
			EnvVars: []string{"MICRO_RUNTIME_SOURCE"},
		},
		&cli.IntFlag{
			Name:    "retries",
			Usage:   "Set the max retries per service",
			EnvVars: []string{"MICRO_RUNTIME_RETRIES"},
		},
	}
)

// Run the runtime service
func Run(ctx *cli.Context, srvOpts ...micro.Option) {
	log.Init(log.WithFields(map[string]interface{}{"service": "runtime"}))

	// Get the profile
	var prof []string
	switch ctx.String("profile") {
	case "local":
		prof = profile.Local()
	case "server":
		prof = profile.Server()
	case "kubernetes":
		prof = profile.Kubernetes()
	case "platform":
		prof = profile.Platform()
	default:
		log.Fatal("Missing runtime profile")
	}

	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}

	if len(ctx.String("server_name")) > 0 {
		Name = ctx.String("server_name")
	}

	if len(Address) > 0 {
		srvOpts = append(srvOpts, micro.Address(Address))
	}

	// create runtime
	muRuntime := *cmd.DefaultCmd.Options().Runtime
	if ctx.IsSet("source") {
		muRuntime.Init(runtime.WithSource(ctx.String("source")))
	}

	// append name
	srvOpts = append(srvOpts, micro.Name(Name))

	// new service
	service := micro.NewService(srvOpts...)

	// create a new runtime manager
	manager := manager.New(muRuntime,
		manager.Auth(service.Options().Auth),
		manager.Store(service.Options().Store),
		manager.Profile(prof),
		manager.CacheStore(service.Options().Store),
	)

	// start the manager
	if err := manager.Start(); err != nil {
		log.Errorf("failed to start: %s", err)
		os.Exit(1)
	}

	// register the runtime handler
	pb.RegisterRuntimeHandler(service.Server(), &handler.Runtime{
		// Client to publish events
		Client: micro.NewEvent("go.micro.runtime.events", service.Client()),
		// using the micro runtime
		Runtime: manager,
	})

	// start runtime service
	if err := service.Run(); err != nil {
		log.Errorf("error running service: %v", err)
	}

	// stop the manager
	if err := manager.Stop(); err != nil {
		log.Errorf("failed to stop: %s", err)
		os.Exit(1)
	}
}

// flags is shared flags so we don't have to continually re-add
func flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "source",
			Usage: "Set the source url of the service e.g github.com/micro/services",
		},
		&cli.StringFlag{
			Name:  "image",
			Usage: "Set the image to use for the container",
		},
		&cli.StringFlag{
			Name:  "command",
			Usage: "Command to exec",
		},
		&cli.StringFlag{
			Name:  "args",
			Usage: "Command args",
		},
		&cli.StringFlag{
			Name:  "type",
			Usage: "The type of service operate on",
		},
		&cli.StringSliceFlag{
			Name:  "env_vars",
			Usage: "Set the environment variables e.g. foo=bar",
		},
	}
}

func Commands(options ...micro.Option) []*cli.Command {
	command := []*cli.Command{
		{
			// In future we'll also have `micro run [x]` hence `micro run service` requiring "service"
			Name:  "run",
			Usage: RunUsage,
			Description: `Examples:
			micro run github.com/micro/examples/helloworld
			micro run .  # deploy local folder to your local micro server
			micro run ../path/to/folder # deploy local folder to your local micro server
			micro run helloworld # deploy latest version, translates to micro run github.com/micro/services/helloworld
			micro run helloworld@9342934e6180 # deploy certain version
			micro run helloworld@branchname	# deploy certain branch`,
			Flags: flags(),
			Action: func(ctx *cli.Context) error {
				runService(ctx, options...)
				return nil
			},
		},
		{
			Name:  "update",
			Usage: UpdateUsage,
			Description: `Examples:
			micro update github.com/micro/examples/helloworld
			micro update .  # deploy local folder to your local micro server
			micro update ../path/to/folder # deploy local folder to your local micro server
			micro update helloworld # deploy master branch, translates to micro update github.com/micro/services/helloworld
			micro update helloworld@branchname	# deploy certain branch`,
			Flags: flags(),
			Action: func(ctx *cli.Context) error {
				updateService(ctx, options...)
				return nil
			},
		},
		{
			Name:  "kill",
			Usage: KillUsage,
			Flags: flags(),
			Description: `Examples:
			micro kill github.com/micro/examples/helloworld
			micro kill .  # kill service deployed from local folder
			micro kill ../path/to/folder # kill service deployed from local folder
			micro kill helloworld # kill serviced deployed from master branch, translates to micro kill github.com/micro/services/helloworld
			micro kill helloworld@branchname	# kill service deployed from certain branch`,
			Action: func(ctx *cli.Context) error {
				killService(ctx, options...)
				return nil
			},
		},
		{
			Name:  "status",
			Usage: GetUsage,
			Flags: flags(),
			Action: func(ctx *cli.Context) error {
				getService(ctx, options...)
				return nil
			},
		},
		{
			Name:  "logs",
			Usage: "Get logs for a service",
			Flags: logFlags(),
			Action: func(ctx *cli.Context) error {
				getLogs(ctx, options...)
				return nil
			},
		},
	}

	return command
}
