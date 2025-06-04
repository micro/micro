package main

import (
	"go-micro.dev/v5/cmd"

	_ "github.com/micro/micro/v5/cmd/micro-cli"
	_ "github.com/micro/micro/v5/cmd/micro-run"
)

var version = "5.0.0-dev"

func main() {
	cmd.Init(
		cmd.Name("micro"),
		cmd.Version(version),
	)
}
