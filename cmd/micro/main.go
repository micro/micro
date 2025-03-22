package main // import github.com/micro/micro/v5/cmd/micro

//go:generate ./scripts/generate.sh

import (
	"github.com/micro/micro/v5/cmd"

	// load packages so they can register commands
	_ "github.com/micro/micro/v5/cmd/cli"
	_ "github.com/micro/micro/v5/cmd/server"
)

func main() {
	cmd.Run()
}
