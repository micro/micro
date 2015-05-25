package main

import (
	ccli "github.com/codegangsta/cli"
	"github.com/myodc/go-micro/cmd"
	"github.com/myodc/micro/api"
	"github.com/myodc/micro/car"
	"github.com/myodc/micro/cli"
)

func main() {
	app := ccli.NewApp()
	app.Name = "micro"
	app.Usage = "A microservices toolchain"
	app.Version = "0.0.1"
	app.Commands = append(app.Commands, api.Commands()...)
	app.Commands = append(app.Commands, cli.Commands()...)
	app.Commands = append(app.Commands, car.Commands()...)
	app.Flags = cmd.Flags
	app.Before = cmd.Setup
	app.RunAndExitOnError()
}
