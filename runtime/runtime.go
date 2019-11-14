// Package runtime is the micro runtime
package runtime

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/config/cmd"
	"github.com/micro/go-micro/runtime"
	"github.com/micro/go-micro/runtime/service/handler"
	pb "github.com/micro/go-micro/runtime/service/proto"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/micro/runtime/notifier"
)

var (
	// Name of the runtime
	Name = "go.micro.runtime"
	// Address of the runtime
	Address = ":8088"
)

func runService(ctx *cli.Context, srvOpts ...micro.Option) {
	log.Name("runtime")

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	if len(ctx.Args()) == 0 || ctx.Args()[0] != "service" {
		log.Fatal("Require usage: micro run service --name example --version latest (optional: --source /path/to/source)")
	}

	// get the args
	name := ctx.String("name")
	version := ctx.String("version")
	source := ctx.String("source")

	if len(name) == 0 {
		log.Fatal("Require usage: micro run service --name example --version latest")
	}

	// get the default runtime
	r := runtime.DefaultRuntime

	// specifier the notifier
	r.Init(runtime.WithNotifier(notifier.New(name, version, source)))

	// start the rutime
	r.Start()
	defer r.Stop()

	// change to the directory of the source
	// TODO: in future
	if len(source) > 0 {
		dir := filepath.Dir(source)
		if err := os.Chdir(dir); err != nil {
			log.Fatalf("Could not change to directory %s: %v", dir, err)
		}
	}

	log.Logf("Starting service: %s version: %s", name, version)

	service := &runtime.Service{
		Name:    name,
		Version: fmt.Sprintf("%d", time.Now().Unix()),
		Exec:    "go run main.go",
	}

	// runtime based on environment we run the service in
	args := []runtime.CreateOption{
		runtime.WithOutput(os.Stdout),
	}

	// run the service
	if err := r.Create(service, args...); err != nil {
		log.Fatal(err)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// wait for shutdown
	<-shutdown

	log.Logf("Stopping service")

	if err := r.Delete(service); err != nil {
		log.Fatal(err)
	}
}

func run(ctx *cli.Context, srvOpts ...micro.Option) {
	log.Name("runtime")

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}

	if len(ctx.GlobalString("server_name")) > 0 {
		Name = ctx.GlobalString("server_name")
	}

	if len(Address) > 0 {
		srvOpts = append(srvOpts, micro.Address(Address))
	}

	muRuntime := *cmd.DefaultCmd.Options().Runtime

	// start the runtime
	if err := muRuntime.Start(); err != nil {
		log.Logf("failed to start: %s", err)
		os.Exit(1)
	}

	// append name
	srvOpts = append(srvOpts, micro.Name(Name))

	// new service
	service := micro.NewService(srvOpts...)

	// register the runtime handler
	pb.RegisterRuntimeHandler(service.Server(), &handler.Runtime{
		// using the mdns runtime
		Runtime: muRuntime,
	})

	// start runtime service
	if err := service.Run(); err != nil {
		log.Logf("error running service: %v", err)
	}

	// stop the runtime
	if err := muRuntime.Stop(); err != nil {
		log.Logf("failed to stop: %s", err)
		os.Exit(1)
	}

	log.Logf("successfully stopped")
}

func Commands(options ...micro.Option) []cli.Command {
	command := []cli.Command{
		{
			Name:  "runtime",
			Usage: "Run the micro runtime",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "address",
					Usage:  "Set the registry http address e.g 0.0.0.0:8000",
					EnvVar: "MICRO_SERVER_ADDRESS",
				},
			},
			Action: func(ctx *cli.Context) {
				run(ctx, options...)
			},
		},
		{
			Name:  "run",
			Usage: "Run a service e.g micro run service version",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name",
					Usage: "Set the name of the service to run",
					Value: "service",
				},
				cli.StringFlag{
					Name:  "version",
					Usage: "Set the version of the service to run",
					Value: "latest",
				},
				cli.StringFlag{
					Name:  "source",
					Usage: "Set the source location of the service e.g /path/to/source",
					Value: ".",
				},
				// TODO: change to BoolFlag
				cli.BoolTFlag{
					Name:  "local",
					Usage: "Set to run the service local",
				},
			},
			Action: func(ctx *cli.Context) {
				runService(ctx, options...)
			},
		},
	}

	for _, p := range Plugins() {
		if cmds := p.Commands(); len(cmds) > 0 {
			command[0].Subcommands = append(command[0].Subcommands, cmds...)
		}

		if flags := p.Flags(); len(flags) > 0 {
			command[0].Flags = append(command[0].Flags, flags...)
		}
	}

	return command
}
