package main // import micro.dev/v4/cmd/micro

//go:generate ./scripts/generate.sh

import (
	"micro.dev/v4/cmd"

	// load packages so they can register commands
	_ "micro.dev/v4/cmd/client"
	_ "micro.dev/v4/cmd/server"
	_ "micro.dev/v4/cmd/web"
)

func main() {
	cmd.Run()
}
