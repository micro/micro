package main

import (
	"embed"
	"go-micro.dev/v5/cmd"

	_ "github.com/micro/micro/v5/cmd/micro-cli"
	_ "github.com/micro/micro/v5/cmd/micro-run"
	"github.com/micro/micro/v5/cmd/micro-server"
)

//go:embed html/*
var htmlFS embed.FS

var version = "5.0.0-dev"

func init() {
	server.HTML = htmlFS
}

func main() {
	cmd.Init(
		cmd.Name("micro"),
		cmd.Version(version),
	)
}
