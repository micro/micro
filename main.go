package main

import (
	"strings"

	ccli "github.com/micro/cli"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/micro/api"
	"github.com/micro/micro/car"
	"github.com/micro/micro/cli"
	"github.com/micro/micro/web"
)

func setup(app *ccli.App) {
	app.Flags = append(app.Flags,
		ccli.StringFlag{
			Name:   "proxy_address",
			Usage:  "Proxy requests via the HTTP address specified",
			EnvVar: "MICRO_PROXY_ADDRESS",
		},
		ccli.BoolFlag{
			Name:   "enable_tls",
			Usage:  "Enable TLS",
			EnvVar: "MICRO_ENABLE_TLS",
		},
		ccli.StringFlag{
			Name:   "tls_cert_file",
			Usage:  "TLS Certificate file",
			EnvVar: "MICRO_TLS_CERT_File",
		},
		ccli.StringFlag{
			Name:   "tls_key_file",
			Usage:  "TLS Key file",
			EnvVar: "MICRO_TLS_KEY_File",
		},
		ccli.StringFlag{
			Name:   "api_address",
			Usage:  "Set the api address e.g 0.0.0.0:8080",
			EnvVar: "MICRO_API_ADDRESS",
		},
		ccli.StringFlag{
			Name:   "sidecar_address",
			Usage:  "Set the sidecar address e.g 0.0.0.0:8081",
			EnvVar: "MICRO_SIDECAR_ADDRESS",
		},
		ccli.StringFlag{
			Name:   "web_address",
			Usage:  "Set the web UI address e.g 0.0.0.0:8082",
			EnvVar: "MICRO_WEB_ADDRESS",
		},
		ccli.IntFlag{
			Name:   "register_ttl",
			EnvVar: "MICRO_REGISTER_TTL",
			Usage:  "Register TTL in seconds",
		},
		ccli.IntFlag{
			Name:   "register_interval",
			EnvVar: "MICRO_REGISTER_INTERVAL",
			Usage:  "Register interval in seconds",
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
			Name:   "sidecar_cors",
			Usage:  "Comma separated whitelist of allowed origins for CORS",
			EnvVar: "MICRO_SIDECAR_CORS",
		},
	)

	before := app.Before

	app.Before = func(ctx *ccli.Context) error {
		if len(ctx.String("api_address")) > 0 {
			api.Address = ctx.String("api_address")
		}
		if len(ctx.String("sidecar_address")) > 0 {
			car.Address = ctx.String("sidecar_address")
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
		if len(ctx.String("sidecar_cors")) > 0 {
			car.CORS = fn(ctx.String("sidecar_cors"))
		}
		if len(ctx.String("web_cors")) > 0 {
			web.CORS = fn(ctx.String("web_cors"))
		}

		return before(ctx)
	}
}

func main() {
	app := cmd.App()
	app.Commands = append(app.Commands, api.Commands()...)
	app.Commands = append(app.Commands, cli.Commands()...)
	app.Commands = append(app.Commands, car.Commands()...)
	app.Commands = append(app.Commands, web.Commands()...)
	app.Action = func(context *ccli.Context) { ccli.ShowAppHelp(context) }

	setup(app)

	cmd.Init(
		cmd.Name("micro"),
		cmd.Description("A microservices toolkit"),
		cmd.Version("latest"),
	)
}
