package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"sort"

	ccli "github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/cmd"
	gostore "github.com/micro/go-micro/v2/store"
	"github.com/micro/micro/v2/plugin"
	"github.com/micro/micro/v2/server"

	// clients
	"github.com/micro/micro/v2/client/api"
	"github.com/micro/micro/v2/client/bot"
	"github.com/micro/micro/v2/client/cli"
	"github.com/micro/micro/v2/client/cli/new"
	"github.com/micro/micro/v2/client/cli/util"
	"github.com/micro/micro/v2/client/proxy"
	"github.com/micro/micro/v2/client/web"

	// load cli packages so they can register commands
	_ "github.com/micro/micro/v2/service/auth/cli"
	_ "github.com/micro/micro/v2/service/cli"
	_ "github.com/micro/micro/v2/service/config/cli"
	_ "github.com/micro/micro/v2/service/debug/cli"
	_ "github.com/micro/micro/v2/service/network/cli"
	_ "github.com/micro/micro/v2/service/runtime/cli"
	_ "github.com/micro/micro/v2/service/store/cli"

	// internals
	"github.com/micro/micro/v2/command"
	inauth "github.com/micro/micro/v2/internal/auth"
	"github.com/micro/micro/v2/internal/helper"
	_ "github.com/micro/micro/v2/internal/plugins"
	"github.com/micro/micro/v2/internal/update"
	_ "github.com/micro/micro/v2/internal/usage"

	// platform related commands
	platform "github.com/micro/micro/v2/platform/cli"
)

var (
	// BuildDate is when the micro binary was built. Injeted via LDFLAGS
	BuildDate string

	// name of the binary
	name = "micro"

	// description of the binary
	description = "A framework for cloud native development\n\n	 Use `micro [command] --help` to see command specific help."
)

type commands []*ccli.Command

func (s commands) Len() int      { return len(s) }
func (s commands) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s commands) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

func setup(app *ccli.App) {
	app.Flags = append(app.Flags,
		&ccli.BoolFlag{
			Name:  "local",
			Usage: "Enable local only development: Defaults to true.",
		},
		&ccli.BoolFlag{
			Name:    "enable_acme",
			Usage:   "Enables ACME support via Let's Encrypt. ACME hosts should also be specified.",
			EnvVars: []string{"MICRO_ENABLE_ACME"},
		},
		&ccli.StringFlag{
			Name:    "acme_hosts",
			Usage:   "Comma separated list of hostnames to manage ACME certs for",
			EnvVars: []string{"MICRO_ACME_HOSTS"},
		},
		&ccli.StringFlag{
			Name:    "acme_provider",
			Usage:   "The provider that will be used to communicate with Let's Encrypt. Valid options: autocert, certmagic",
			EnvVars: []string{"MICRO_ACME_PROVIDER"},
		},
		&ccli.BoolFlag{
			Name:    "enable_tls",
			Usage:   "Enable TLS support. Expects cert and key file to be specified",
			EnvVars: []string{"MICRO_ENABLE_TLS"},
		},
		&ccli.StringFlag{
			Name:    "tls_cert_file",
			Usage:   "Path to the TLS Certificate file",
			EnvVars: []string{"MICRO_TLS_CERT_FILE"},
		},
		&ccli.StringFlag{
			Name:    "tls_key_file",
			Usage:   "Path to the TLS Key file",
			EnvVars: []string{"MICRO_TLS_KEY_FILE"},
		},
		&ccli.StringFlag{
			Name:    "tls_client_ca_file",
			Usage:   "Path to the TLS CA file to verify clients against",
			EnvVars: []string{"MICRO_TLS_CLIENT_CA_FILE"},
		},
		&ccli.StringFlag{
			Name:    "api_address",
			Usage:   "Set the api address e.g 0.0.0.0:8080",
			EnvVars: []string{"MICRO_API_ADDRESS"},
		},
		&ccli.StringFlag{
			Name:    "namespace",
			Usage:   "Set the micro service namespace",
			EnvVars: []string{"MICRO_NAMESPACE"},
			Value:   "micro",
		},
		&ccli.StringFlag{
			Name:    "proxy_address",
			Usage:   "Proxy requests via the HTTP address specified",
			EnvVars: []string{"MICRO_PROXY_ADDRESS"},
		},
		&ccli.StringFlag{
			Name:    "web_address",
			Usage:   "Set the web UI address e.g 0.0.0.0:8082",
			EnvVars: []string{"MICRO_WEB_ADDRESS"},
		},
		&ccli.StringFlag{
			Name:    "network",
			Usage:   "Set the micro network name: local, go.micro",
			EnvVars: []string{"MICRO_NETWORK"},
		},
		&ccli.StringFlag{
			Name:    "network_address",
			Usage:   "Set the micro network address e.g. :9093",
			EnvVars: []string{"MICRO_NETWORK_ADDRESS"},
		},
		&ccli.StringFlag{
			Name:    "gateway_address",
			Usage:   "Set the micro default gateway address e.g. :9094",
			EnvVars: []string{"MICRO_GATEWAY_ADDRESS"},
		},
		&ccli.StringFlag{
			Name:    "tunnel_address",
			Usage:   "Set the micro tunnel address e.g. :8083",
			EnvVars: []string{"MICRO_TUNNEL_ADDRESS"},
		},
		&ccli.StringFlag{
			Name:    "api_handler",
			Usage:   "Specify the request handler to be used for mapping HTTP requests to services; {api, proxy, rpc}",
			EnvVars: []string{"MICRO_API_HANDLER"},
		},
		&ccli.StringFlag{
			Name:    "api_namespace",
			Usage:   "Set the namespace used by the API e.g. com.example.api",
			EnvVars: []string{"MICRO_API_NAMESPACE"},
		},
		&ccli.StringFlag{
			Name:    "web_namespace",
			Usage:   "Set the namespace used by the Web proxy e.g. com.example.web",
			EnvVars: []string{"MICRO_WEB_NAMESPACE"},
		},
		&ccli.StringFlag{
			Name:    "web_url",
			Usage:   "Set the host used for the web dashboard e.g web.example.com",
			EnvVars: []string{"MICRO_WEB_HOST"},
		},
		&ccli.BoolFlag{
			Name:    "enable_stats",
			Usage:   "Enable stats",
			EnvVars: []string{"MICRO_ENABLE_STATS"},
		},
		&ccli.BoolFlag{
			Name:    "auto_update",
			Usage:   "Enable automatic updates",
			EnvVars: []string{"MICRO_AUTO_UPDATE"},
		},
		&ccli.StringFlag{
			Name:    "update_url",
			Usage:   "Set the url to retrieve system updates from",
			EnvVars: []string{"MICRO_UPDATE_URL"},
			Value:   update.DefaultURL,
		},
		&ccli.BoolFlag{
			Name:    "report_usage",
			Usage:   "Report usage statistics",
			EnvVars: []string{"MICRO_REPORT_USAGE"},
			Value:   true,
		},
		&ccli.StringFlag{
			Name:    "env",
			Aliases: []string{"e"},
			Usage:   "Override environment",
			EnvVars: []string{"MICRO_ENV"},
		},
	)

	plugins := plugin.Plugins()

	for _, p := range plugins {
		if flags := p.Flags(); len(flags) > 0 {
			app.Flags = append(app.Flags, flags...)
		}

		if cmds := p.Commands(); len(cmds) > 0 {
			app.Commands = append(app.Commands, cmds...)
		}
	}

	before := app.Before

	app.Before = func(ctx *ccli.Context) error {

		if len(ctx.String("api_handler")) > 0 {
			api.Handler = ctx.String("api_handler")
		}
		if len(ctx.String("api_address")) > 0 {
			api.Address = ctx.String("api_address")
		}
		if len(ctx.String("proxy_address")) > 0 {
			proxy.Address = ctx.String("proxy_address")
		}
		if len(ctx.String("web_address")) > 0 {
			web.Address = ctx.String("web_address")
		}
		if len(ctx.String("api_namespace")) > 0 {
			api.Namespace = ctx.String("api_namespace")
		}
		if len(ctx.String("web_namespace")) > 0 {
			web.Namespace = ctx.String("web_namespace")
		}
		if len(ctx.String("web_host")) > 0 {
			web.Host = ctx.String("web_host")
		}

		for _, p := range plugins {
			if err := p.Init(ctx); err != nil {
				return err
			}
		}

		util.SetupCommand(ctx)
		// now do previous before
		if err := before(ctx); err != nil {
			// DO NOT return this error otherwise the action will fail
			// and help will be printed.
			fmt.Println(err)
			os.Exit(1)
			return err
		}

		var opts []gostore.Option

		// the database is not overriden by flag then set it
		if len(ctx.String("store_database")) == 0 {
			opts = append(opts, gostore.Database(cmd.App().Name))
		}

		// if the table is not overriden by flag then set it
		if len(ctx.String("store_table")) == 0 {
			table := cmd.App().Name

			// if an arg is specified use that as the name
			// so each service has its own table preconfigured
			if name := ctx.Args().First(); len(name) > 0 {
				table = name
			}

			opts = append(opts, gostore.Table(table))
		}

		// TODO: move this entire initialisation elsewhere
		// maybe in service.Run so all things are configured
		if len(opts) > 0 {
			(*cmd.DefaultCmd.Options().Store).Init(opts...)
		}

		// add the system rules if we're using the JWT implementation
		// which doesn't have access to the rules in the auth service
		if (*cmd.DefaultCmd.Options().Auth).String() == "jwt" {
			for _, rule := range inauth.SystemRules {
				if err := (*cmd.DefaultCmd.Options().Auth).Grant(rule); err != nil {
					return err
				}
			}
		}

		return nil
	}
}

// Run executes the command line
func Run(options ...micro.Option) {
	// get the app
	app := cmd.App()

	// register commands
	app.Commands = append(app.Commands, command.List()...)

	// Add the client commmands
	app.Commands = append(app.Commands, api.Commands()...)
	app.Commands = append(app.Commands, web.Commands()...)
	app.Commands = append(app.Commands, proxy.Commands()...)
	app.Commands = append(app.Commands, bot.Commands()...)
	app.Commands = append(app.Commands, cli.Commands()...)

	// Add the service commands
	app.Commands = append(app.Commands, new.Commands()...)
	app.Commands = append(app.Commands, server.Commands(options...)...)
	app.Commands = append(app.Commands, platform.Commands(options...)...)

	sort.Sort(commands(app.Commands))

	// boot micro runtime
	app.Action = func(c *ccli.Context) error {
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

	// add additional flags and setup functions
	setup(app)

	// initialise our options
	cmd.DefaultCmd.Init(
		cmd.Name(name),
		cmd.Description(description),
		cmd.Version(buildVersion()),
	)

	// run the command
	if err := cmd.DefaultCmd.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
