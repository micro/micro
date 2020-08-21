// Package runtime is the micro runtime
package runtime

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/micro/cli/v2"
	golog "github.com/micro/go-micro/v3/logger"
	goruntime "github.com/micro/go-micro/v3/runtime"
	"github.com/micro/go-micro/v3/runtime/local/git"
	"github.com/micro/go-micro/v3/util/file"
	"github.com/micro/micro/v3/client/cli/namespace"
	"github.com/micro/micro/v3/client/cli/util"
	cliutil "github.com/micro/micro/v3/client/cli/util"
	"github.com/micro/micro/v3/internal/config"
	muclient "github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/context"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/runtime"
	"github.com/micro/micro/v3/service/runtime/server"
	"google.golang.org/grpc/status"
)

const (
	// RunUsage message for the run command
	RunUsage = "Run a service: micro run [source]"
	// KillUsage message for the kill command
	KillUsage = "Kill a service: micro kill [source]"
	// UpdateUsage message for the update command
	UpdateUsage = "Update a service: micro update [source]"
	// GetUsage message for micro get command
	GetUsage = "Get the status of services"
	// ServicesUsage message for micro services command
	ServicesUsage = "micro services"
	// CannotWatch message for the run command
	CannotWatch = "Cannot watch filesystem on this runtime"
)

var (
	// DefaultRetries which should be attempted when starting a service
	DefaultRetries = 3
	// DefaultImage which should be run
	DefaultImage = "micro/cells:go"
	// Git orgs we currently support for credentials
	GitOrgs = []string{"github", "bitbucket", "gitlab"}
)

const (
	credentialsKey = "GIT_CREDENTIALS"
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

func sourceExists(source *git.Source) error {
	ref := source.Ref
	if ref == "" || ref == "latest" {
		ref = "master"
	}

	sourceExistsAt := func(url string, source *git.Source) error {
		req, _ := http.NewRequest("GET", url, nil)

		// add the git credentials if set
		if creds, ok := getGitCredentials(source.Repo); ok {
			req.Header.Set("Authorization", "token "+creds)
		}

		client := new(http.Client)
		resp, err := client.Do(req)

		// @todo gracefully degrade?
		if err != nil {
			return err
		}
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			return fmt.Errorf("service at %v@%v not found", source.Repo, ref)
		}
		return nil
	}

	if strings.Contains(source.Repo, "github") {
		// Github specific existence checs
		repo := strings.ReplaceAll(source.Repo, "github.com/", "")
		url := fmt.Sprintf("https://api.github.com/repos/%v/contents/%v?ref=%v", repo, source.Folder, ref)
		return sourceExistsAt(url, source)
	} else if strings.Contains(source.Repo, "gitlab") {
		// Gitlab specific existence checks

		// @todo better check for gitlab
		url := fmt.Sprintf("https://%v", source.Repo)
		return sourceExistsAt(url, source)
	}
	return nil
}

func appendSourceBase(ctx *cli.Context, workDir, source string) string {
	isLocal, _ := git.IsLocal(workDir, source)
	// @todo add list of supported hosts here or do this check better
	if !isLocal && !strings.Contains(source, ".com") && !strings.Contains(source, ".org") && !strings.Contains(source, ".net") {
		baseURL, _ := config.Get("git", util.GetEnv(ctx).Name, "baseurl")
		if len(baseURL) == 0 {
			baseURL, _ = config.Get("git", "baseurl")
		}
		if len(baseURL) == 0 {
			return path.Join("github.com/micro/services", source)
		}
		return path.Join(baseURL, source)
	}
	return source
}

func runService(ctx *cli.Context) error {
	// we need some args to run
	if ctx.Args().Len() == 0 {
		fmt.Println(RunUsage)
		return nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	source, err := git.ParseSourceLocal(wd, appendSourceBase(ctx, wd, ctx.Args().Get(0)))
	if err != nil {
		return err
	}
	var newSource string
	if source.Local {
		if cliutil.IsPlatform(ctx) {
			fmt.Println("Local sources are not yet supported on m3o. It's coming soon though!")
			os.Exit(1)
		}
		newSource, err = upload(ctx, source)
		if err != nil {
			return err
		}
	} else {
		err := sourceExists(source)
		if err != nil {
			return err
		}
	}

	typ := ctx.String("type")
	command := strings.TrimSpace(ctx.String("command"))
	args := strings.TrimSpace(ctx.String("args"))

	runtimeSource := source.RuntimeSource()
	if source.Local {
		runtimeSource = newSource
	}

	var retries = DefaultRetries
	if ctx.IsSet("retries") {
		retries = ctx.Int("retries")
	}

	var image = DefaultImage
	if ctx.IsSet("image") {
		image = ctx.String("image")
	}

	// when using the micro/cells:go image, we pass the source as the argument
	args = runtimeSource
	if len(source.Ref) > 0 {
		args += "@" + source.Ref
	}

	// specify the options
	opts := []goruntime.CreateOption{
		goruntime.WithOutput(os.Stdout),
		goruntime.WithRetries(retries),
		goruntime.CreateImage(image),
		goruntime.CreateType(typ),
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
		opts = append(opts, goruntime.WithEnv(environment))
	}

	if len(command) > 0 {
		opts = append(opts, goruntime.WithCommand(strings.Split(command, " ")...))
	}

	if len(args) > 0 {
		opts = append(opts, goruntime.WithArgs(strings.Split(args, " ")...))
	}

	// determine the namespace
	ns, err := namespace.Get(util.GetEnv(ctx).Name)
	if err != nil {
		return err
	}
	opts = append(opts, goruntime.CreateNamespace(ns))
	gitCreds, ok := getGitCredentials(source.Repo)
	if ok {
		opts = append(opts, goruntime.WithSecret(credentialsKey, gitCreds))
	}

	// run the service
	service := &goruntime.Service{
		Name:     source.RuntimeName(),
		Source:   runtimeSource,
		Version:  source.Ref,
		Metadata: make(map[string]string),
	}

	if err := runtime.Create(service, opts...); err != nil {
		return err
	}

	if runtime.DefaultRuntime.String() == "local" {
		// we need to wait
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)
		<-ch
		// delete the service
		return runtime.Delete(service)
	}

	return nil
}

func getGitCredentials(repo string) (string, bool) {
	repo = strings.Split(repo, "/")[0]

	for _, org := range GitOrgs {
		if !strings.Contains(repo, org) {
			continue
		}

		// check the creds for the org
		creds, err := config.Get("git", "credentials", org)
		if err == nil && len(creds) > 0 {
			return creds, true
		}
	}

	return "", false
}

func killService(ctx *cli.Context) error {
	// we need some args to run
	if ctx.Args().Len() == 0 {
		fmt.Println(RunUsage)
		return nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	source, err := git.ParseSourceLocal(wd, appendSourceBase(ctx, wd, ctx.Args().Get(0)))
	if err != nil {
		return err
	}
	service := &goruntime.Service{
		Name:    source.RuntimeName(),
		Source:  source.RuntimeSource(),
		Version: source.Ref,
	}

	// determine the namespace
	ns, err := namespace.Get(util.GetEnv(ctx).Name)
	if err != nil {
		return err
	}

	if err := runtime.Delete(service, goruntime.DeleteNamespace(ns)); err != nil {
		return err
	}

	return nil
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
		file := filepath.Join(path, f.Name())
		b, err := ioutil.ReadFile(file)
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
	uploadedFileName := filepath.Base(source.Folder) + ".tar.gz"
	path := filepath.Join(os.TempDir(), uploadedFileName)

	var err error
	if len(source.LocalRepoRoot) > 0 {
		// @todo currently this uploads the whole repo all the time to support local dependencies
		// in parents (ie service path is `repo/a/b/c` and it depends on `repo/a/b`).
		// Optimise this by only uploading things that are needed.
		err = server.Compress(source.LocalRepoRoot, path)
	} else {
		err = server.Compress(source.FullPath, path)
	}

	if err != nil {
		return "", err
	}
	cli := muclient.DefaultClient
	err = file.New("server", cli, file.WithContext(context.DefaultContext)).Upload(uploadedFileName, path)
	if err != nil {
		return "", err
	}
	// ie. if relative folder path to repo root is `test/service/example`
	// file name becomes `example.tar.gz/test/service`
	parts := strings.Split(source.Folder, "/")
	if len(parts) == 1 {
		return uploadedFileName, nil
	}
	allButLastDir := parts[0 : len(parts)-1]
	return filepath.Join(append([]string{uploadedFileName}, allButLastDir...)...), nil
}

func updateService(ctx *cli.Context) error {
	// we need some args to run
	if ctx.Args().Len() == 0 {
		fmt.Println(RunUsage)
		return nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	source, err := git.ParseSourceLocal(wd, appendSourceBase(ctx, wd, ctx.Args().Get(0)))
	if err != nil {
		return err
	}
	var newSource string
	if source.Local {
		newSource, err = upload(ctx, source)
		if err != nil {
			return err
		}
	} else {
		err := sourceExists(source)
		if err != nil {
			return err
		}
	}

	runtimeSource := source.RuntimeSource()
	if source.Local {
		runtimeSource = newSource
	}
	service := &goruntime.Service{
		Name:    source.RuntimeName(),
		Source:  runtimeSource,
		Version: source.Ref,
	}

	// determine the namespace
	ns, err := namespace.Get(util.GetEnv(ctx).Name)
	if err != nil {
		return err
	}

	opts := []goruntime.UpdateOption{goruntime.UpdateNamespace(ns)}
	gitCreds, ok := getGitCredentials(source.Repo)
	if ok {
		opts = append(opts, goruntime.UpdateSecret(credentialsKey, gitCreds))
	}
	return runtime.Update(service, goruntime.UpdateNamespace(ns))
}

func getService(ctx *cli.Context) error {
	name := ""
	version := "latest"
	typ := ctx.String("type")

	if ctx.Args().Len() > 0 {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		source, err := git.ParseSourceLocal(wd, ctx.Args().Get(0))
		if err != nil {
			return err
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

	var services []*goruntime.Service
	var readOpts []goruntime.ReadOption

	// return a list of services
	switch list {
	case true:
		// return specific type listing
		if len(typ) > 0 {
			readOpts = append(readOpts, goruntime.ReadType(typ))
		}
	// return one service
	default:
		// check if service name was passed in
		if len(name) == 0 {
			fmt.Println(GetUsage)
			return nil
		}

		// get service with name and version
		readOpts = []goruntime.ReadOption{
			goruntime.ReadService(name),
			goruntime.ReadVersion(version),
		}

		// return the runtime services
		if len(typ) > 0 {
			readOpts = append(readOpts, goruntime.ReadType(typ))
		}
	}

	// determine the namespace
	ns, err := namespace.Get(util.GetEnv(ctx).Name)
	if err != nil {
		return err
	}
	readOpts = append(readOpts, goruntime.ReadNamespace(ns))

	// read the service
	services, err = runtime.Read(readOpts...)
	if err != nil {
		return err
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
		return nil
	}

	sort.Slice(services, func(i, j int) bool { return services[i].Name < services[j].Name })

	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)
	fmt.Fprintln(writer, "NAME\tVERSION\tSOURCE\tSTATUS\tBUILD\tUPDATED\tMETADATA")
	for _, service := range services {
		status := parse(service.Metadata["status"])

		// cut the commit down to first 7 characters
		build := parse(service.Metadata["build"])
		if len(build) > 7 {
			build = build[:7]
		}

		// if there is an error, display this in metadata (there is no error field)
		metadata := fmt.Sprintf("owner=%s, group=%s", parse(service.Metadata["owner"]), parse(service.Metadata["group"]))
		if status == "error" {
			metadata = fmt.Sprintf("%v, error=%v", metadata, parse(service.Metadata["error"]))
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
			metadata)
	}
	writer.Flush()
	return nil
}

const (
	// logUsage message for logs command
	logUsage = "Required usage: micro log example"
)

func getLogs(ctx *cli.Context) error {
	logger.DefaultLogger.Init(golog.WithFields(map[string]interface{}{"service": "runtime"}))
	if ctx.Args().Len() == 0 {
		fmt.Println("Service name is required")
		return nil
	}

	name := ctx.Args().Get(0)

	// must specify service name
	if len(name) == 0 {
		fmt.Println(logUsage)
		return nil
	}

	// get the args
	options := []goruntime.LogsOption{}

	count := ctx.Int("lines")
	if count > 0 {
		options = append(options, goruntime.LogsCount(int64(count)))
	} else {
		options = append(options, goruntime.LogsCount(int64(15)))
	}

	follow := ctx.Bool("follow")

	if follow {
		options = append(options, goruntime.LogsStream(follow))
	}

	// @todo reintroduce since
	//since := ctx.String("since")
	//var readSince time.Time
	//d, err := time.ParseDuration(since)
	//if err == nil {
	//	readSince = time.Now().Add(-d)
	//}

	// determine the namespace
	ns, err := namespace.Get(util.GetEnv(ctx).Name)
	if err != nil {
		return err
	}
	options = append(options, goruntime.LogsNamespace(ns))

	logs, err := runtime.Logs(&goruntime.Service{Name: name}, options...)

	if err != nil {
		return err
	}

	output := ctx.String("output")
	for {
		select {
		case record, ok := <-logs.Chan():
			if !ok {
				if err := logs.Error(); err != nil {
					fmt.Printf("Error reading logs: %s\n", status.Convert(err).Message())
					os.Exit(1)
				}
				return nil
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
