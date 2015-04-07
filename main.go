package main

import (
	"os"

	"github.com/asim/go-micro/cmd"
	"github.com/asim/micro/api"
	"github.com/asim/micro/cli"
	ccli "github.com/codegangsta/cli"
)

func main() {
	cmd.Init()

	app := ccli.NewApp()
	app.Name = "micro"
	app.Usage = "A microservices toolchain"
	app.Version = "0.0.1"
	app.Commands = append(app.Commands, api.Commands()...)
	app.Commands = append(app.Commands, cli.Commands()...)
	app.Run(os.Args)
}
