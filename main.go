package main

import (
	ccli "github.com/micro/cli"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/micro/api"
	"github.com/micro/micro/car"
	"github.com/micro/micro/cli"
	"github.com/micro/micro/web"
)

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
	micro.Init()
}
