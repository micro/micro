package cmd

import (
	"strings"

	ccli "github.com/micro/cli"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/micro/api"
	"github.com/micro/micro/bot"
	"github.com/micro/micro/cli"
	"github.com/micro/micro/new"
	"github.com/micro/micro/plugin"
	"github.com/micro/micro/proxy"
	"github.com/micro/micro/web"
)

var (
	name        = "micro"
	description = "A cloud-native toolkit"
	version     = "0.9.0"
)

func setup(app *ccli.App) {
	app.Flags = append(app.Flags,
		ccli.BoolFlag{
			Name:   "enable_acme",
			Usage:  "Enables ACME support via Let's Encrypt. ACME hosts should also be specified.",
			EnvVar: "MICRO_ENABLE_ACME",
		},
		ccli.StringFlag{
			Name:   "acme_hosts",
			Usage:  "Comma separated list of hostnames to manage ACME certs for",
			EnvVar: "MICRO_ACME_HOSTS",
		},
		ccli.BoolFlag{
			Name:   "enable_tls",
			Usage:  "Enable TLS support. Expects cert and key file to be specified",
			EnvVar: "MICRO_ENABLE_TLS",
		},
		ccli.StringFlag{
			Name:   "tls_cert_file",
			Usage:  "Path to the TLS Certificate file",
			EnvVar: "MICRO_TLS_CERT_FILE",
		},
		ccli.StringFlag{
			Name:   "tls_key_file",
			Usage:  "Path to the TLS Key file",
			EnvVar: "MICRO_TLS_KEY_FILE",
		},
		ccli.StringFlag{
			Name:   "tls_client_ca_file",
			Usage:  "Path to the TLS CA file to verify clients against",
			EnvVar: "MICRO_TLS_CLIENT_CA_FILE",
		},
		ccli.StringFlag{
			Name:   "api_address",
			Usage:  "Set the api address e.g 0.0.0.0:8080",
			EnvVar: "MICRO_API_ADDRESS",
		},
		ccli.StringFlag{
			Name:   "proxy_address",
			Usage:  "Proxy requests via the HTTP address specified",
			EnvVar: "MICRO_PROXY_ADDRESS",
		},
		ccli.StringFlag{
			Name:   "web_address",
			Usage:  "Set the web UI address e.g 0.0.0.0:8082",
			EnvVar: "MICRO_WEB_ADDRESS",
		},
		ccli.StringFlag{
			Name:   "api_handler",
			Usage:  "Specify the request handler to be used for mapping HTTP requests to services; {api, proxy, rpc}",
			EnvVar: "MICRO_API_HANDLER",
		},
		ccli.StringFlag{
			Name:   "api_namespace",
			Usage:  "Set the namespace used by the API e.g. com.example.api",
			EnvVar: "MICRO_API_NAMESPACE",
		},
		ccli.StringFlag{
			Name:   "web_namespace",
			Usage:  "Set the namespace used by the Web proxy e.g. com.example.web",
			EnvVar: "MICRO_WEB_NAMESPACE",
		},
		ccli.StringFlag{
			Name:   "api_cors",
			Usage:  "Comma separated whitelist of allowed origins for CORS",
			EnvVar: "MICRO_API_CORS",
		},
		ccli.StringFlag{
			Name:   "web_cors",
			Usage:  "Comma separated whitelist of allowed origins for CORS",
			EnvVar: "MICRO_WEB_CORS",
		},
		ccli.StringFlag{
			Name:   "proxy_cors",
			Usage:  "Comma separated whitelist of allowed origins for CORS",
			EnvVar: "MICRO_PROXY_CORS",
		},
		ccli.BoolFlag{
			Name:   "enable_stats",
			Usage:  "Enable stats",
			EnvVar: "MICRO_ENABLE_STATS",
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

		// origin comma separated string to map
		fn := func(s string) map[string]bool {
			origins := make(map[string]bool)
			for _, origin := range strings.Split(s, ",") {
				origins[origin] = true
			}
			return origins
		}

		if len(ctx.String("api_cors")) > 0 {
			api.CORS = fn(ctx.String("api_cors"))
		}
		if len(ctx.String("proxy_cors")) > 0 {
			proxy.CORS = fn(ctx.String("proxy_cors"))
		}
		if len(ctx.String("web_cors")) > 0 {
			web.CORS = fn(ctx.String("web_cors"))
		}

		for _, p := range plugins {
			p.Init(ctx)
		}

		return before(ctx)
	}
}

func Init() {
	app := cmd.App()
	app.Commands = append(app.Commands, api.Commands()...)
	app.Commands = append(app.Commands, bot.Commands()...)
	app.Commands = append(app.Commands, cli.Commands()...)
	app.Commands = append(app.Commands, proxy.Commands()...)
	app.Commands = append(app.Commands, new.Commands()...)
	app.Commands = append(app.Commands, web.Commands()...)
	app.Action = func(context *ccli.Context) { ccli.ShowAppHelp(context) }

	setup(app)

	cmd.Init(
		cmd.Name(name),
		cmd.Description(description),
		cmd.Version(version),
	)
}
