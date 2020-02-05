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

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/config/cmd"
	"github.com/micro/go-micro/v2/runtime"
	rs "github.com/micro/go-micro/v2/runtime/service"
	"github.com/micro/micro/v2/runtime/scheduler"
)

const (
	// RunUsage message for the run command
	RunUsage = "Required usage: micro run [service] [version] [--source github.com/micro/services --watch]"
	// KillUsage message for the kill command
	KillUsage = "Require usage: micro kill [service] [version]"
	// GetUsage message for micro get command
	GetUsage = "Require usage: micro ps [service] [version]"
	// CannotWatch message for the run command
	CannotWatch = "Cannot watch filesystem on this runtime"
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
	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	// we need some args to run
	if ctx.Args().Len() == 0 {
		fmt.Println(RunUsage)
		return
	}

	// set and validate the name (arg 1)
	name := ctx.Args().Get(0)
	if strings.HasPrefix(name, ".") || strings.HasPrefix(name, "/") {
		fmt.Println(RunUsage)
		return
	}

	// set the version (arg 2, optional)
	version := "latest"
	if ctx.Args().Len() > 1 {
		version = ctx.Args().Get(1)
	}

	// load the runtime. The default runtime is ignored if running on the platform
	var r runtime.Runtime
	if ctx.Bool("platform") {
		r = rs.NewRuntime()
		// TODO @BEN: Proxy to platform
	} else {
		r = *cmd.DefaultCmd.Options().Runtime
	}

	source := ctx.String("source")
	exec := []string{"go", "run", source}

	// Determine the filepath
	fp := filepath.Join(os.Getenv("GOPATH"), "src", source, name)

	// Find the filepath or `go run` will pull from git by default
	if r.String() == "local" && os.Chdir(fp) == nil {
		exec = []string{"go", "run", "."}

		// watch the filesystem for changes
		sched := scheduler.New(name, version, fp)
		if err := r.Init(runtime.WithScheduler(sched)); err != nil {
			fmt.Printf("Could not start scheduler: %v", err)
			return
		}
	}

	// start the runtimes
	if err := r.Start(); err != nil {
		fmt.Printf("Could not start: %v", err)
		return
	}

	// add environment variable passed in via cli
	environment := defaultEnv()
	for _, evar := range ctx.StringSlice("env") {
		for _, e := range strings.Split(evar, ",") {
			if len(e) > 0 {
				environment = append(environment, strings.TrimSpace(e))
			}
		}
	}

	// specify the options
	opts := []runtime.CreateOption{
		runtime.WithCommand(exec...),
		runtime.WithOutput(os.Stdout),
		runtime.WithEnv(environment),
	}

	// run the service
	service := &runtime.Service{
		Name:     name,
		Source:   source,
		Version:  version,
		Metadata: make(map[string]string),
	}
	if err := r.Create(service, opts...); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Started service %v: %v\n", name, source)

	// if local	 then register signal handlers
	if r.String() == "local" {
		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

		// wait for shutdown
		<-shutdown

		// delete service from runtime
		if err := r.Delete(service); err != nil {
			fmt.Println(err)
			return
		}

		if err := r.Stop(); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func killService(ctx *cli.Context, srvOpts ...micro.Option) {
	// get the args
	name := ctx.String("name")
	version := ctx.String("version")
	local := ctx.Bool("local")

	if ctx.Args().Len() > 0 {
		// set name to first arg
		name = ctx.Args().Get(0)
		if ctx.Args().Len() > 1 {
			version = ctx.Args().Get(1)
		}
	}

	if len(name) == 0 {
		fmt.Println(KillUsage)
		return
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
		fmt.Println(err)
		return
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

	// zero args so list all
	if ctx.Args().Len() == 0 {
		list = true
	} else {
		// set name as first arg
		name = ctx.Args().Get(0)
		// set version as second arg
		if ctx.Args().Len() > 1 {
			version = ctx.Args().Get(1)
		}
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
			fmt.Println(GetUsage)
			return
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
		fmt.Println(err)
		return
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
		status := parse(service.Metadata["status"])
		if status == "error" {
			status = service.Metadata["error"]
		}

		fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\n",
			service.Name,
			parse(service.Version),
			parse(service.Source),
			status,
			parse(service.Metadata["build"]),
			fmt.Sprintf("owner=%s,group=%s", parse(service.Metadata["owner"]), parse(service.Metadata["group"])))
	}
	writer.Flush()
}
