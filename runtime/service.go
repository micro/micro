// Package runtime is the micro runtime
package runtime

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/config/cmd"
	"github.com/micro/go-micro/v2/runtime"
	srvRuntime "github.com/micro/go-micro/v2/runtime/service"
)

const (
	// RunUsage message for the run command
	RunUsage = "Required usage: micro run [service] [version] [--source github.com/micro/services]"
	// KillUsage message for the kill command
	KillUsage = "Require usage: micro kill [service] [version]"
	// UpdateUsage message for the update command
	UpdateUsage = "Require usage: micro update [service] [version]"
	// GetUsage message for micro get command
	GetUsage = "Require usage: micro ps [service] [version]"
	// ServicesUsage message for micro services command
	ServicesUsage = "Require usage: micro services"
	// CannotWatch message for the run command
	CannotWatch = "Cannot watch filesystem on this runtime"
)

var (
	// DefaultRetries which should be attempted when starting a service
	DefaultRetries = 3
	// Image to specify if none is specified
	Image = "docker.pkg.github.com/micro/services"
	// Source where we get services from
	Source = "github.com/micro/services"
)

// timeAgo returns the time passed
func timeAgo(v string) string {
	if len(v) == 0 {
		return "unknown"
	}
	t, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return v
	}
	return fmt.Sprintf("%v ago", time.Since(t).Truncate(time.Second))
}

func runtimeFromContext(ctx *cli.Context) runtime.Runtime {
	if ctx.Bool("platform") {
		os.Setenv("MICRO_PROXY", "service")
		os.Setenv("MICRO_PROXY_ADDRESS", "proxy.micro.mu:443")
		return srvRuntime.NewRuntime()
	}

	return *cmd.DefaultCmd.Options().Runtime
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
	version := "latest"
	source := ctx.String("source")
	typ := ctx.String("type")
	image := ctx.String("image")
	command := strings.TrimSpace(ctx.String("command"))
	args := strings.TrimSpace(ctx.String("args"))

	// load the runtime
	r := runtimeFromContext(ctx)

	// set the version (arg 2, optional)
	if ctx.Args().Len() > 1 {
		version = ctx.Args().Get(1)
	}

	// add environment variable passed in via cli
	var environment []string
	for _, evar := range ctx.StringSlice("env") {
		for _, e := range strings.Split(evar, ",") {
			if len(e) > 0 {
				environment = append(environment, strings.TrimSpace(e))
			}
		}
	}

	var retries = DefaultRetries
	if ctx.IsSet("retries") {
		retries = ctx.Int("retries")
	}

	// set the image from our images if its the platform
	if ctx.Bool("platform") && len(image) == 0 {
		if strings.HasPrefix(name, ".") || strings.HasPrefix(name, "/") {
			fmt.Println(RunUsage)
			return
		}

		formattedName := strings.ReplaceAll(name, "/", "-")
		image = fmt.Sprintf("%v/%v", Image, formattedName)
	}

	// check the source is set
	if ctx.Bool("platform") && len(source) == 0 {
		source = Source
	}

	// specify the options
	opts := []runtime.CreateOption{
		runtime.WithOutput(os.Stdout),
		runtime.WithRetries(retries),
		runtime.CreateImage(image),
		runtime.CreateType(typ),
	}

	if len(environment) > 0 {
		opts = append(opts, runtime.WithEnv(environment))
	}

	if len(command) > 0 {
		opts = append(opts, runtime.WithCommand(strings.Split(command, " ")...))
	}

	if len(args) > 0 {
		opts = append(opts, runtime.WithArgs(strings.Split(args, " ")...))
	}

	// don't pass through dotted names unless
	if strings.HasPrefix(name, ".") || strings.HasPrefix(name, "/") {
		if r.String() != "local" {
			fmt.Println(RunUsage)
			return
		}
		path, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			return
		}
		name = filepath.Base(path)
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

	if r.String() == "local" {
		// we need to wait
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)
		<-ch
		// delete the service
		r.Delete(service)
	}
}

func killService(ctx *cli.Context, srvOpts ...micro.Option) {
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

	service := &runtime.Service{
		Name:    name,
		Version: version,
	}

	if err := runtimeFromContext(ctx).Delete(service); err != nil {
		fmt.Println(err)
		return
	}
}

func updateService(ctx *cli.Context, srvOpts ...micro.Option) {
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

	service := &runtime.Service{
		Name:    name,
		Version: version,
	}

	if err := runtimeFromContext(ctx).Update(service); err != nil {
		fmt.Println(err)
		return
	}
}

func getService(ctx *cli.Context, srvOpts ...micro.Option) {
	name := ctx.Args().Get(0)
	version := "latest"
	typ := ctx.String("type")
	r := runtimeFromContext(ctx)

	if strings.HasPrefix(name, ".") || strings.HasPrefix(name, "/") {
		fmt.Println(GetUsage)
		return
	}

	// set version as second arg
	if ctx.Args().Len() > 1 {
		version = ctx.Args().Get(1)
	}

	// should we list sevices
	var list bool

	// zero args so list all
	if ctx.Args().Len() == 0 {
		list = true
	}

	var err error
	var services []*runtime.Service

	// return a list of services
	switch list {
	case true:
		// return specific type listing
		if len(typ) > 0 {
			services, err = r.Read(runtime.ReadType(typ))
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
		if len(typ) > 0 {
			opts = append(opts, runtime.ReadType(typ))
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
	fmt.Fprintln(writer, "NAME\tVERSION\tSOURCE\tSTATUS\tBUILD\tUPDATED\tMETADATA")
	for _, service := range services {
		status := parse(service.Metadata["status"])
		if status == "error" {
			status = service.Metadata["error"]
		}

		// cut the commit down to first 7 characters
		build := parse(service.Metadata["build"])
		if len(build) > 7 {
			build = build[:7]
		}

		// parse when the service was started
		updated := parse(timeAgo(service.Metadata["started"]))

		fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			service.Name,
			parse(service.Version),
			parse(service.Source),
			status,
			build,
			updated,
			fmt.Sprintf("owner=%s,group=%s", parse(service.Metadata["owner"]), parse(service.Metadata["group"])))
	}
	writer.Flush()
}
