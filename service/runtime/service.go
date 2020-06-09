// Package runtime is the micro runtime
package runtime

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/runtime"
	"github.com/micro/go-micro/v2/runtime/local/git"
	srvRuntime "github.com/micro/go-micro/v2/runtime/service"
	"github.com/micro/go-micro/v2/util/file"
	cliutil "github.com/micro/micro/v2/client/cli/util"
	"github.com/micro/micro/v2/internal/client"
	"github.com/micro/micro/v2/service/runtime/handler"
)

const (
	// RunUsage message for the run command
	RunUsage = "Required usage: micro run [source]"
	// KillUsage message for the kill command
	KillUsage = "Require usage: micro kill [source]"
	// UpdateUsage message for the update command
	UpdateUsage = "Require usage: micro update [source]"
	// GetUsage message for micro get command
	GetUsage = "Require usage: micro status [service] [version]"
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
	if cliutil.IsLocal(ctx) {
		return *cmd.DefaultCmd.Options().Runtime
	}

	return srvRuntime.NewRuntime(runtime.WithClient(client.New(ctx)))
}

// exists returns whether the given file or directory exists
func dirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
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

	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	source, err := git.ParseSourceLocal(wd, ctx.Args().Get(0))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var newSource string
	if source.Local {
		newSource, err = upload(ctx, source)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	typ := ctx.String("type")
	image := ctx.String("image")
	command := strings.TrimSpace(ctx.String("command"))
	args := strings.TrimSpace(ctx.String("args"))

	// load the runtime
	r := runtimeFromContext(ctx)

	var retries = DefaultRetries
	if ctx.IsSet("retries") {
		retries = ctx.Int("retries")
	}

	// set the image if not specified
	if len(image) == 0 {
		formattedName := strings.ReplaceAll(source.Folder, "/", "-")
		// eg. docker.pkg.github.com/micro/services/users-api
		image = fmt.Sprintf("%v/%v", Image, formattedName)
	}

	// specify the options
	opts := []runtime.CreateOption{
		runtime.WithOutput(os.Stdout),
		runtime.WithRetries(retries),
		runtime.CreateImage(image),
		runtime.CreateType(typ),
	}

	// add environment variable passed in via cli
	var environment []string
	for _, evar := range ctx.StringSlice("env_vars") {
		for _, e := range strings.Split(evar, ",") {
			if len(e) > 0 {
				environment = append(environment, strings.TrimSpace(e))
			}
		}
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

	runtimeSource := source.RuntimeSource()
	if source.Local {
		runtimeSource = newSource
	}
	// run the service
	service := &runtime.Service{
		Name:     source.RuntimeName(),
		Source:   runtimeSource,
		Version:  source.Ref,
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

	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	source, err := git.ParseSourceLocal(wd, ctx.Args().Get(0))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	service := &runtime.Service{
		Name:    source.RuntimeName(),
		Source:  source.RuntimeSource(),
		Version: source.Ref,
	}

	if err := runtimeFromContext(ctx).Delete(service); err != nil {
		fmt.Println(err)
		return
	}
}

func grepMain(path string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".go") {
			continue
		}
		b, err := ioutil.ReadFile(f.Name())
		if err != nil {
			continue
		}
		if strings.Contains(string(b), "package main") {
			return nil
		}
	}
	return fmt.Errorf("Directory does not contain a main package")
}

func upload(ctx *cli.Context, source *git.Source) (string, error) {
	if err := grepMain(source.FullPath); err != nil {
		return "", err
	}
	uploadedFileName := strings.ReplaceAll(source.Folder, string(filepath.Separator), "-") + ".tar.gz"
	path := filepath.Join(os.TempDir(), uploadedFileName)
	err := handler.Compress(source.FullPath, path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = file.New("go.micro.server", client.New(ctx)).Upload(uploadedFileName, path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return uploadedFileName, nil
}

func updateService(ctx *cli.Context, srvOpts ...micro.Option) {
	// we need some args to run
	if ctx.Args().Len() == 0 {
		fmt.Println(RunUsage)
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	source, err := git.ParseSourceLocal(wd, ctx.Args().Get(0))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var newSource string
	if source.Local {
		newSource, err = upload(ctx, source)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	runtimeSource := source.RuntimeSource()
	if source.Local {
		runtimeSource = newSource
	}
	service := &runtime.Service{
		Name:    source.RuntimeName(),
		Source:  runtimeSource,
		Version: source.Ref,
	}

	if err := runtimeFromContext(ctx).Update(service); err != nil {
		fmt.Println(err)
		return
	}
}

func getService(ctx *cli.Context, srvOpts ...micro.Option) {
	name := ""
	version := "latest"
	typ := ctx.String("type")
	r := runtimeFromContext(ctx)

	if ctx.Args().Len() > 0 {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		source, err := git.ParseSourceLocal(wd, ctx.Args().Get(0))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		name = source.RuntimeName()
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

	var services []*runtime.Service
	var readOpts []runtime.ReadOption

	// return a list of services
	switch list {
	case true:
		// return specific type listing
		if len(typ) > 0 {
			readOpts = append(readOpts, runtime.ReadType(typ))
		}
	// return one service
	default:
		// check if service name was passed in
		if len(name) == 0 {
			fmt.Println(GetUsage)
			return
		}

		// get service with name and version
		readOpts = []runtime.ReadOption{
			runtime.ReadService(name),
			runtime.ReadVersion(version),
		}

		// return the runtime services
		if len(typ) > 0 {
			readOpts = append(readOpts, runtime.ReadType(typ))
		}

	}

	// read the service
	services, err := r.Read(readOpts...)
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
			strings.ToLower(status),
			build,
			updated,
			fmt.Sprintf("owner=%s,group=%s", parse(service.Metadata["owner"]), parse(service.Metadata["group"])))
	}
	writer.Flush()
}

const (
	// logUsage message for logs command
	logUsage = "Required usage: micro log example"
)

func getLogs(ctx *cli.Context, srvOpts ...micro.Option) {
	log.Init(log.WithFields(map[string]interface{}{"service": "runtime"}))
	if ctx.Args().Len() == 0 {
		fmt.Println("Service name is required")
		return
	}

	name := ctx.Args().Get(0)

	// must specify service name
	if len(name) == 0 {
		fmt.Println(logUsage)
		return
	}

	// get the args
	options := []runtime.LogsOption{}

	count := ctx.Int("lines")
	if count > 0 {
		options = append(options, runtime.LogsCount(int64(count)))
	} else {
		options = append(options, runtime.LogsCount(int64(15)))
	}

	follow := ctx.Bool("follow")

	if follow {
		options = append(options, runtime.LogsStream(follow))
	}

	r := runtimeFromContext(ctx)

	// @todo reintroduce since
	//since := ctx.String("since")
	//var readSince time.Time
	//d, err := time.ParseDuration(since)
	//if err == nil {
	//	readSince = time.Now().Add(-d)
	//}

	logs, err := r.Logs(&runtime.Service{Name: name}, options...)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	output := ctx.String("output")
	for {
		select {
		case record, ok := <-logs.Chan():
			if !ok {
				return
			}
			switch output {
			case "json":
				b, _ := json.Marshal(record)
				fmt.Printf("%v\n", string(b))
			default:
				fmt.Printf("%v\n", record.Message)

			}
		}
	}
}

// logFlags is shared flags so we don't have to continually re-add
func logFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "version",
			Usage: "Set the version of the service to debug",
		},
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "Set the output format e.g json, text",
		},
		&cli.BoolFlag{
			Name:    "follow",
			Aliases: []string{"f"},
			Usage:   "Set to stream logs continuously (default: true)",
		},
		&cli.StringFlag{
			Name:  "since",
			Usage: "Set to the relative time from which to show the logs for e.g. 1h",
		},
		&cli.IntFlag{
			Name:    "lines",
			Aliases: []string{"n"},
			Usage:   "Set to query the last number of log events",
		},
	}
}
