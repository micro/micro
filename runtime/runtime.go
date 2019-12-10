// Package runtime is the micro runtime
package runtime

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"text/tabwriter"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/config/cmd"
	"github.com/micro/go-micro/runtime"
	rs "github.com/micro/go-micro/runtime/service"
	pb "github.com/micro/go-micro/runtime/service/proto"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/micro/runtime/handler"
	"github.com/micro/micro/runtime/notifier"
)

const (
	// RunUsage message for the run command
	RunUsage = "Required usage: micro run service --name example --version latest --source go/package/import/path"
	// KillUsage message for the kill command
	KillUsage = "Require usage: micro kill service --name example (optional: --version latest)"
	// Getusage message for micro get command
	GetUsage = "Require usage: micro ps service --name example (optional: --version latest)"
)

var (
	// Name of the runtime
	Name = "go.micro.runtime"
	// Address of the runtime
	Address = ":8088"
)

func defaultEnv() []string {
	var env []string
	for _, evar := range os.Environ() {
		if strings.HasPrefix(evar, "MICRO_") {
			env = append(env, evar)
		}
	}

	return env
}

func runService(ctx *cli.Context, srvOpts ...micro.Option) {
	log.Name("runtime")

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	if len(ctx.Args()) == 0 || ctx.Args()[0] != "service" {
		log.Fatal(RunUsage)
	}

	// get the args
	name := ctx.String("name")
	version := ctx.String("version")
	source := ctx.String("source")
	env := ctx.StringSlice("env")
	local := ctx.Bool("local")

	var r runtime.Runtime
	var exec []string

	// must specify service name
	if len(name) == 0 {
		if len(source) > 0 {
			name = filepath.Base(source)
		} else {
			cwd, _ := os.Getwd()
			name = filepath.Base(cwd)
		}
	}

	// local usage specified
	switch local {
	case true:
		r = *cmd.DefaultCmd.Options().Runtime
		// NOTE: When in local mode, we consider source to be
		// the filesystem path to the source of the service
		if len(source) > 0 {
			if err := os.Chdir(source); err != nil {
				log.Fatalf("Failed to change directory to %s: %v", source, err)
			}
		}
		exec = []string{"go", "run", "."}

		// specify the runtime notifier to update wiht local file changes
		if err := r.Init(runtime.WithNotifier(notifier.New(name, version, source))); err != nil {
			log.Fatalf("Could not start notifier: %v", err)
		}
	default:
		// new service runtime
		r = rs.NewRuntime()
		// NOTE: we consider source in default mode
		// to be the canonical Go module import path
		// if source is empty, we bail as this can
		// lead to a potential K8s API object creation DDOS
		if len(source) == 0 {
			log.Fatal(RunUsage)
		}
		exec = []string{"go", "run", source}
	}

	// start the local runtime
	if err := r.Start(); err != nil {
		log.Fatalf("Could not start: %v", err)
	}

	service := &runtime.Service{
		Name:     name,
		Source:   source,
		Version:  version,
		Metadata: make(map[string]string),
	}

	// default environment
	environment := defaultEnv()
	// add environment variable passed in via cli
	for _, evar := range env {
		for _, e := range strings.Split(evar, ",") {
			if len(e) > 0 {
				environment = append(environment, strings.TrimSpace(e))
			}
		}
	}

	// runtime based on environment we run the service in
	// TODO: how will this work with runtime service
	opts := []runtime.CreateOption{
		runtime.WithCommand(exec...),
		runtime.WithOutput(os.Stdout),
		runtime.WithEnv(environment),
	}

	log.Debugf("Attempting to start service: %v", service)

	// run the service
	if err := r.Create(service, opts...); err != nil {
		log.Fatal(err)
	}

	// if in local mode register signal handlers
	if local {
		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

		// wait for shutdown
		<-shutdown

		log.Logf("Stopping service")

		// delete service from runtime
		if err := r.Delete(service); err != nil {
			log.Fatal(err)
		}

		if err := r.Stop(); err != nil {
			log.Fatal(err)
		}
	}
}

func killService(ctx *cli.Context, srvOpts ...micro.Option) {
	// we expect `micro run service`
	if len(ctx.Args()) == 0 || ctx.Args()[0] != "service" {
		log.Fatal(KillUsage)
	}

	// get the args
	name := ctx.String("name")
	version := ctx.String("version")
	local := ctx.Bool("local")

	if len(name) == 0 {
		log.Fatal(KillUsage)
	}

	var r runtime.Runtime
	switch local {
	case true:
		r = *cmd.DefaultCmd.Options().Runtime
	default:
		r = rs.NewRuntime()
	}

	service := &runtime.Service{
		Name:    name,
		Version: version,
	}

	if err := r.Delete(service); err != nil {
		log.Fatal(err)
	}
}

func getService(ctx *cli.Context, srvOpts ...micro.Option) {
	// get the args
	name := ctx.String("name")
	version := ctx.String("version")
	local := ctx.Bool("local")
	runType := ctx.Bool("runtime")

	var r runtime.Runtime
	switch local {
	case true:
		r = *cmd.DefaultCmd.Options().Runtime
	default:
		r = rs.NewRuntime()
	}

	var list bool

	if len(ctx.Args()) == 0 || ctx.Args()[0] != "service" {
		list = true
	}

	var services []*runtime.Service
	var err error

	// return a list of services
	switch list {
	case true:
		// return the runtiem services
		if runType {
			services, err = r.Read(runtime.ReadType("runtime"))
		} else {
			// list all running services
			services, err = r.List()
		}
	// return one service
	default:
		// check if service name was passed in
		if len(name) == 0 {
			log.Fatal(GetUsage)
		}

		// get service with name and version
		opts := []runtime.ReadOption{
			runtime.ReadService(name),
			runtime.ReadVersion(version),
		}

		// return the runtime services
		if runType {
			opts = append(opts, runtime.ReadType("runtime"))
		}

		// read the service
		services, err = r.Read(opts...)
	}

	// check the error
	if err != nil {
		log.Fatal(err)
	}

	// make sure we return UNKNOWN when empty string is supplied
	parse := func(m string) string {
		if len(m) == 0 {
			return "n/a"
		}
		return m
	}

	// don't do anything if there's no services
	if len(services) == 0 {
		return
	}

	sort.Slice(services, func(i, j int) bool { return services[i].Name < services[j].Name })

	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)
	fmt.Fprintln(writer, "NAME\tVERSION\tSOURCE\tSTATUS\tBUILD\tMETADATA")
	for _, service := range services {
		fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\n",
			service.Name,
			parse(service.Version),
			parse(service.Source),
			parse(service.Metadata["status"]),
			parse(service.Metadata["build"]),
			fmt.Sprintf("owner=%s,group=%s", parse(service.Metadata["owner"]), parse(service.Metadata["group"])))
	}
	writer.Flush()
}

func Run(ctx *cli.Context, srvOpts ...micro.Option) {
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

	// create runtime
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

// Flags is shared flags so we don't have to continually re-add
func Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "Set the name of the service to run",
		},
		cli.StringFlag{
			Name:  "version",
			Usage: "Set the version of the service to run",
			Value: "latest",
		},
		cli.StringFlag{
			Name:  "source",
			Usage: "Set the source url of the service e.g /path/to/source",
		},
		cli.BoolFlag{
			Name:  "local",
			Usage: "Set to run the service from local path",
		},
		cli.StringSliceFlag{
			Name:  "env",
			Usage: "Set the environment variables e.g. foo=bar",
		},
		cli.BoolFlag{
			Name:  "runtime",
			Usage: "Return the runtime services",
		},
	}
}

func Commands(options ...micro.Option) []cli.Command {
	command := []cli.Command{
		{
			Name:  "runtime",
			Usage: "Run the micro runtime",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "address",
					Usage:  "Set the registry http address e.g 0.0.0.0:8088",
					EnvVar: "MICRO_SERVER_ADDRESS",
				},
			},
			Action: func(ctx *cli.Context) {
				Run(ctx, options...)
			},
		},
		{
			// In future we'll also have `micro run [x]` hence `micro run service` requiring "service"
			Name:  "run",
			Usage: "Run a service e.g micro run service version",
			Flags: Flags(),
			Action: func(ctx *cli.Context) {
				runService(ctx, options...)
			},
		},
		{
			Name:  "kill",
			Usage: "Kill removes a running service e.g micro kill service",
			Flags: Flags(),
			Action: func(ctx *cli.Context) {
				killService(ctx, options...)
			},
		},
		{
			Name:  "ps",
			Usage: "Ps returns status of a running service or lists all running services e.g. micro ps",
			Flags: Flags(),
			Action: func(ctx *cli.Context) {
				getService(ctx, options...)
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
