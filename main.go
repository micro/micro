package main

import (
	"github.com/asim/go-micro/cmd"
	"github.com/asim/micro/api"
	"github.com/asim/micro/cli"
	"github.com/asim/micro/sic"
	ccli "github.com/codegangsta/cli"
)

func main() {
	app := ccli.NewApp()
	app.Name = "micro"
	app.Usage = "A microservices toolchain"
	app.Version = "0.0.1"
	app.Commands = append(app.Commands, api.Commands()...)
	app.Commands = append(app.Commands, cli.Commands()...)
	app.Commands = append(app.Commands, sic.Commands()...)
	app.Flags = cmd.Flags
	app.Before = cmd.Setup
	app.RunAndExitOnError()
}
