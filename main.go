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
	app.Flags = append(app.Flags, ccli.StringFlag{
		Name:   "api_address",
		Usage:  "Set the api address e.g 0.0.0.0:8080",
		EnvVar: "MICRO_API_ADDRESS",
	})

	before := app.Before

	app.Before = func(ctx *ccli.Context) error {
		if len(ctx.String("api_address")) > 0 {
			api.Address = ctx.String("api_address")
		}
		return before(ctx)
	}
}

func main() {
	micro := cmd.NewCmd(
		cmd.Name("micro"),
		cmd.Description("A microservices toolchain"),
		cmd.Version("latest"),
	)

	app := micro.App()
	app.Commands = append(app.Commands, api.Commands()...)
	app.Commands = append(app.Commands, cli.Commands()...)
	app.Commands = append(app.Commands, car.Commands()...)
	app.Commands = append(app.Commands, web.Commands()...)
	app.Action = func(context *ccli.Context) { ccli.ShowAppHelp(context) }

	setup(app)
	micro.Init()
}
