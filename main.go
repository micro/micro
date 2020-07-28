package main

//go:generate ./scripts/generate.sh

import (
	"fmt"
	"os"

	"github.com/micro/micro/v3/cmd"

	// internal packages
	_ "github.com/micro/micro/v3/internal/usage"

	// load packages so they can register commands
	_ "github.com/micro/micro/v3/client/api"
	_ "github.com/micro/micro/v3/client/bot"
	_ "github.com/micro/micro/v3/client/cli"
	_ "github.com/micro/micro/v3/client/cli/new"
	_ "github.com/micro/micro/v3/client/proxy"
	_ "github.com/micro/micro/v3/client/web"
	_ "github.com/micro/micro/v3/platform/cli"
	_ "github.com/micro/micro/v3/server"
	_ "github.com/micro/micro/v3/service/auth/cli"
	_ "github.com/micro/micro/v3/service/cli"
	_ "github.com/micro/micro/v3/service/config/cli"
	_ "github.com/micro/micro/v3/service/debug/cli"
	_ "github.com/micro/micro/v3/service/network/cli"
	_ "github.com/micro/micro/v3/service/runtime/cli"
	_ "github.com/micro/micro/v3/service/store/cli"
)

func main() {
	if err := cmd.DefaultCmd.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
