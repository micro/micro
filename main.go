package main

import (
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
