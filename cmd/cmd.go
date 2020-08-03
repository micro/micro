package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/micro/go-micro/v3/client"
	"github.com/micro/go-micro/v3/server"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v3/cmd"
	"github.com/micro/micro/v3/client/cli/util"
	incmd "github.com/micro/micro/v3/internal/cmd"
	uconf "github.com/micro/micro/v3/internal/config"
	"github.com/micro/micro/v3/internal/helper"
	_ "github.com/micro/micro/v3/internal/usage"
	"github.com/micro/micro/v3/plugin"
	"github.com/micro/micro/v3/profile"
	"github.com/micro/micro/v3/service/logger"

	muclient "github.com/micro/micro/v3/service/client"
	muregistry "github.com/micro/micro/v3/service/registry"
	muserver "github.com/micro/micro/v3/service/server"
)

type command struct {
	opts cmd.Options
	app  *cli.App

	sync.Mutex
	// exit is a channel which is closed
	// on exit for anything that requires
	// cleanup
	exit chan bool
}

var (
	DefaultCmd cmd.Cmd = New()

	// name of the binary
	name = "micro"
	// description of the binary
	description = "A framework for cloud native development\n\n	 Use `micro [command] --help` to see command specific help."
	// defaultFlags which are used on all commands
	defaultFlags = []cli.Flag{
		&cli.StringFlag{
			Name:    "c",
			Usage:   "Set the config file: Defaults to ~/.micro/config.json",
			EnvVars: []string{"MICRO_CONFIG_FILE"},
		},
		&cli.StringFlag{
			Name:    "env",
			Aliases: []string{"e"},
			Usage:   "Set the environment to operate in",
			EnvVars: []string{"MICRO_ENV"},
		},
		&cli.StringFlag{
			Name:    "profile",
			Usage:   "Set the micro server profile: e.g. local or kubernetes",
			EnvVars: []string{"MICRO_PROFILE"},
		},
		&cli.StringFlag{
			Name:    "namespace",
			EnvVars: []string{"MICRO_NAMESPACE"},
			Usage:   "Namespace the service is operating in",
			Value:   "micro",
		},
		&cli.StringFlag{
			Name:    "auth_id",
			EnvVars: []string{"MICRO_AUTH_ID"},
			Usage:   "Account ID used for client authentication",
		},
		&cli.StringFlag{
			Name:    "auth_secret",
			EnvVars: []string{"MICRO_AUTH_SECRET"},
			Usage:   "Account secret used for client authentication",
		},
		&cli.StringFlag{
			Name:    "auth_public_key",
			EnvVars: []string{"MICRO_AUTH_PUBLIC_KEY"},
			Usage:   "Public key for JWT auth (base64 encoded PEM)",
		},
		&cli.StringFlag{
			Name:    "auth_private_key",
			EnvVars: []string{"MICRO_AUTH_PRIVATE_KEY"},
			Usage:   "Private key for JWT auth (base64 encoded PEM)",
		},
		&cli.StringFlag{
			Name:    "registry_address",
			EnvVars: []string{"MICRO_REGISTRY_ADDRESS"},
			Usage:   "Comma-separated list of registry addresses",
		},
		&cli.StringFlag{
			Name:    "registry_tls_ca",
			Usage:   "Certificate authority for TLS with registry",
			EnvVars: []string{"MICRO_REGISTRY_TLS_CA"},
		},
		&cli.StringFlag{
			Name:    "registry_tls_cert",
			Usage:   "Client cert for TLS with registry",
			EnvVars: []string{"MICRO_REGISTRY_TLS_CERT"},
		},
		&cli.StringFlag{
			Name:    "registry_tls_key",
			Usage:   "Client key for TLS with registry",
			EnvVars: []string{"MICRO_REGISTRY_TLS_KEY"},
		},
		&cli.StringFlag{
			Name:    "broker_address",
			EnvVars: []string{"MICRO_BROKER_ADDRESS"},
			Usage:   "Comma-separated list of broker addresses",
		},
		&cli.StringFlag{
			Name:    "broker_tls_ca",
			Usage:   "Certificate authority for TLS with broker",
			EnvVars: []string{"MICRO_BROKER_TLS_CA"},
		},
		&cli.StringFlag{
			Name:    "broker_tls_cert",
			Usage:   "Client cert for TLS with broker",
			EnvVars: []string{"MICRO_BROKER_TLS_CERT"},
		},
		&cli.StringFlag{
			Name:    "broker_tls_key",
			Usage:   "Client key for TLS with broker",
			EnvVars: []string{"MICRO_BROKER_TLS_KEY"},
		},
		&cli.StringFlag{
			Name:    "store_address",
			EnvVars: []string{"MICRO_STORE_ADDRESS"},
			Usage:   "Comma-separated list of store addresses",
		},
		&cli.StringFlag{
			Name:    "proxy_address",
			Usage:   "Proxy requests via the HTTP address specified",
			EnvVars: []string{"MICRO_PROXY"},
		},
		&cli.BoolFlag{
			Name:    "report_usage",
			Usage:   "Report usage statistics",
			EnvVars: []string{"MICRO_REPORT_USAGE"},
			Value:   true,
		},
		&cli.StringFlag{
			Name:    "service_name",
			Usage:   "Name of the micro service",
			EnvVars: []string{"MICRO_SERVICE_NAME"},
		},
		&cli.StringFlag{
			Name:    "service_version",
			Usage:   "Version of the micro service",
			EnvVars: []string{"MICRO_SERVICE_VERSION"},
		},
	}
)

func init() {
	rand.Seed(time.Now().Unix())
}

func New(opts ...cmd.Option) cmd.Cmd {
	options := cmd.Options{}
	for _, o := range opts {
		o(&options)
	}

	cmd := new(command)
	cmd.exit = make(chan bool)
	cmd.opts = options
	cmd.app = cli.NewApp()
	cmd.app.Name = name
	cmd.app.Version = buildVersion()
	cmd.app.Usage = description
	cmd.app.Flags = defaultFlags
	cmd.app.Action = action
	cmd.app.Before = beforeFromContext(options.Context, cmd.Before)
	cmd.app.After = cmd.After
	return cmd
}

func (c *command) App() *cli.App {
	return c.app
}

func (c *command) Options() cmd.Options {
	return c.opts
}

// After is executed after any subcommand
func (c *command) After(ctx *cli.Context) error {
	c.Lock()
	defer c.Unlock()

	select {
	case <-c.exit:
		return nil
	default:
		close(c.exit)
	}

	return nil
}

// Before is executed before any subcommand
func (c *command) Before(ctx *cli.Context) error {
	// set the config file if specified
	if cf := ctx.String("c"); len(cf) > 0 {
		uconf.SetConfig(cf)
	}

	// initialize plugins
	for _, p := range plugin.Plugins() {
		if err := p.Init(ctx); err != nil {
			return err
		}
	}

	// default the profile for the server
	prof := ctx.String("profile")

	// if no profile is set then set one
	if len(prof) == 0 {
		switch ctx.Args().First() {
		case "service", "server":
			prof = "local"
		default:
			prof = "client"
		}
	}

	// apply the profile
	if profile, err := profile.Load(prof); err != nil {
		logger.Fatal(err)
	} else {
		// load the profile
		profile.Setup(ctx)
	}

	// use the external proxy which is loaded from the local config
	if proxy := util.CLIProxyAddress(ctx); len(proxy) > 0 {
		muclient.DefaultClient.Init(client.Proxy(proxy))
	}

	// set the registry in the client and server
	muclient.DefaultClient.Init(client.Registry(muregistry.DefaultRegistry))
	muserver.DefaultServer.Init(server.Registry(muregistry.DefaultRegistry))

	// inject auth credentials from the local storage for CLI auth
	if err := setupAuthForCLI(ctx); err != nil {
		logger.Fatalf("Error setting up auth: %v", err)
	}

	// call the init funcs registered by other modules
	// in micro such as auth, registry etc
	for _, fn := range incmd.InitFuncs {
		if err := fn(ctx); err != nil {
			return err
		}
	}

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

func action(c *cli.Context) error {
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
			cmdStr := strings.Join(c.Args().Slice(), " ")
			fmt.Printf("Error querying registry for service %v: %v", cmdStr, err)
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

	return helper.MissingCommand(c)
}

// Register CLI commands
func Register(cmds ...*cli.Command) {
	app := DefaultCmd.App()
	app.Commands = append(app.Commands, cmds...)

	// sort the commands so they're listed in order on the cli
	// todo: move this to micro/cli so it's only run when the
	// commands are printed during "help"
	sort.Slice(app.Commands, func(i, j int) bool {
		return app.Commands[i].Name < app.Commands[j].Name
	})
}
