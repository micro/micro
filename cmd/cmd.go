package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"time"

	"github.com/micro/cli/v2"
	ccli "github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/cmd"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/micro/v2/client/cli/util"
	"github.com/micro/micro/v2/internal/helper"
	_ "github.com/micro/micro/v2/internal/usage"
	"github.com/micro/micro/v2/plugin"
	"github.com/micro/micro/v2/profile"
)

type command struct {
	opts cmd.Options
	app  *cli.App
}

var (
	DefaultCmd cmd.Cmd = newCmd()

	// name of the binary
	name = "micro"
	// description of the binary
	description = "A framework for cloud native development\n\n	 Use `micro [command] --help` to see command specific help."
	// defaultFlags which are used on all commands
	defaultFlags = []ccli.Flag{
		&ccli.StringFlag{
			Name:    "profile",
			Usage:   "Set the micro profile: local, test or platform",
			EnvVars: []string{"MICRO_PROFILE"},
			Value:   "local",
		},
		// &ccli.StringFlag{
		// 	Name:    "api_address",
		// 	Usage:   "Set the api address e.g 0.0.0.0:8080",
		// 	EnvVars: []string{"MICRO_API_ADDRESS"},
		// },
		// &ccli.StringFlag{
		// 	Name:    "namespace",
		// 	Usage:   "Set the micro service namespace",
		// 	EnvVars: []string{"MICRO_NAMESPACE"},
		// 	Value:   "micro",
		// },
		// &ccli.StringFlag{
		// 	Name:    "proxy_address",
		// 	Usage:   "Proxy requests via the HTTP address specified",
		// 	EnvVars: []string{"MICRO_PROXY_ADDRESS"},
		// },
		// &ccli.StringFlag{
		// 	Name:    "web_address",
		// 	Usage:   "Set the web UI address e.g 0.0.0.0:8082",
		// 	EnvVars: []string{"MICRO_WEB_ADDRESS"},
		// },
		// &ccli.StringFlag{
		// 	Name:    "network",
		// 	Usage:   "Set the micro network name: local, go.micro",
		// 	EnvVars: []string{"MICRO_NETWORK"},
		// },
		// &ccli.StringFlag{
		// 	Name:    "network_address",
		// 	Usage:   "Set the micro network address e.g. :9093",
		// 	EnvVars: []string{"MICRO_NETWORK_ADDRESS"},
		// },
		// &ccli.StringFlag{
		// 	Name:    "gateway_address",
		// 	Usage:   "Set the micro default gateway address e.g. :9094",
		// 	EnvVars: []string{"MICRO_GATEWAY_ADDRESS"},
		// },
		// &ccli.StringFlag{
		// 	Name:    "tunnel_address",
		// 	Usage:   "Set the micro tunnel address e.g. :8083",
		// 	EnvVars: []string{"MICRO_TUNNEL_ADDRESS"},
		// },
		// &ccli.StringFlag{
		// 	Name:    "api_handler",
		// 	Usage:   "Specify the request handler to be used for mapping HTTP requests to services; {api, proxy, rpc}",
		// 	EnvVars: []string{"MICRO_API_HANDLER"},
		// },
		// &ccli.StringFlag{
		// 	Name:    "api_namespace",
		// 	Usage:   "Set the namespace used by the API e.g. com.example.api",
		// 	EnvVars: []string{"MICRO_API_NAMESPACE"},
		// },
		// &ccli.StringFlag{
		// 	Name:    "web_namespace",
		// 	Usage:   "Set the namespace used by the Web proxy e.g. com.example.web",
		// 	EnvVars: []string{"MICRO_WEB_NAMESPACE"},
		// },
		// &ccli.StringFlag{
		// 	Name:    "web_url",
		// 	Usage:   "Set the host used for the web dashboard e.g web.example.com",
		// 	EnvVars: []string{"MICRO_WEB_HOST"},
		// },
		// &ccli.BoolFlag{
		// 	Name:    "enable_stats",
		// 	Usage:   "Enable stats",
		// 	EnvVars: []string{"MICRO_ENABLE_STATS"},
		// },
		// &ccli.BoolFlag{
		// 	Name:    "auto_update",
		// 	Usage:   "Enable automatic updates",
		// 	EnvVars: []string{"MICRO_AUTO_UPDATE"},
		// },
		// &ccli.StringFlag{
		// 	Name:    "update_url",
		// 	Usage:   "Set the url to retrieve system updates from",
		// 	EnvVars: []string{"MICRO_UPDATE_URL"},
		// 	Value:   update.DefaultURL,
		// },
		// &ccli.BoolFlag{
		// 	Name:    "report_usage",
		// 	Usage:   "Report usage statistics",
		// 	EnvVars: []string{"MICRO_REPORT_USAGE"},
		// 	Value:   true,
		// },
		// &ccli.StringFlag{
		// 	Name:    "env",
		// 	Aliases: []string{"e"},
		// 	Usage:   "Override environment",
		// 	EnvVars: []string{"MICRO_ENV"},
		// },
	}
)

func init() {
	rand.Seed(time.Now().Unix())
}

func newCmd(opts ...cmd.Option) cmd.Cmd {
	options := cmd.Options{}
	for _, o := range opts {
		o(&options)
	}

	cmd := new(command)
	cmd.opts = options
	cmd.app = cli.NewApp()
	cmd.app.Name = name
	cmd.app.Version = buildVersion()
	cmd.app.Usage = description
	cmd.app.Flags = defaultFlags
	cmd.app.Action = action
	cmd.app.Before = before

	return cmd
}

func (c *command) App() *cli.App {
	return c.app
}

func (c *command) Options() cmd.Options {
	return c.opts
}

func (c *command) Before(ctx *cli.Context) error {
	return nil
}

func (c *command) Init(opts ...cmd.Option) error {
	for _, o := range opts {
		o(&c.opts)
	}
	if len(c.opts.Name) > 0 {
		c.app.Name = c.opts.Name
	}
	if len(c.opts.Version) > 0 {
		c.app.Version = c.opts.Version
	}
	c.app.HideVersion = len(c.opts.Version) == 0
	c.app.Usage = c.opts.Description

	return nil
}

func (c *command) Run() error {
	return c.app.Run(os.Args)
}

func (c *command) String() string {
	return "micro"
}

func before(ctx *ccli.Context) error {
	for _, p := range plugin.Plugins() {
		if err := p.Init(ctx); err != nil {
			return err
		}
	}

	// configure the server
	if ctx.Args().First() == "service" || ctx.Args().First() == "server" {
		switch ctx.String("profile") {
		case "local":
			profile.Local()
		case "test":
			profile.Test()
		case "platform":
			profile.Platform()
		default:
			logger.Fatalf("Unknown profile: %v", ctx.String("profile"))
		}
	}

	// set the proxy address
	util.SetupCommand(ctx)

	return nil
}

func action(c *ccli.Context) error {
	if c.Args().Len() > 0 {
		// if an executable is available with the name of
		// the command, execute it with the arguments from
		// index 1 on.
		v, err := exec.LookPath("micro-" + c.Args().First())
		if err == nil {
			ce := exec.Command(v, c.Args().Slice()[1:]...)
			ce.Stdout = os.Stdout
			ce.Stderr = os.Stderr
			return ce.Run()
		}

		// lookup the service, e.g. "micro config set" would
		// firstly check to see if the service "go.micro.config"
		// exists within the current namespace, then it would
		// execute the Config.Set RPC, setting the flags in the
		// request.
		if srv, err := lookupService(c); err != nil {
			fmt.Printf("Error querying registry for service: %v", err)
			os.Exit(1)
		} else if srv != nil && c.Args().Len() == 1 {
			fmt.Println(formatServiceUsage(srv, c.Args().First()))
			os.Exit(1)
		} else if srv != nil {
			if err := callService(srv, c); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			os.Exit(0)
		}

	}
	fmt.Println(helper.MissingCommand(c))
	os.Exit(1)
	return nil
}

// Register CLI commands
func Register(cmds ...*ccli.Command) {
	app := DefaultCmd.App()
	app.Commands = append(app.Commands, cmds...)

	// sort the commands so they're listed in order on the cli
	// todo: move this to micro/cli so it's only run when the
	// commands are printed during "help"
	sort.Slice(app.Commands, func(i, j int) bool {
		return app.Commands[i].Name < app.Commands[j].Name
	})
}
