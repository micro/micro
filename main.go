package main

//go:generate ./.github/generate.sh

import (
	"fmt"
	"os"

	"github.com/micro/micro/v2/cmd"

	// load packages so they can register commands
	_ "github.com/micro/micro/v2/client/api"
	_ "github.com/micro/micro/v2/client/bot"
	_ "github.com/micro/micro/v2/client/cli"
	_ "github.com/micro/micro/v2/client/cli/new"
	_ "github.com/micro/micro/v2/client/proxy"
	_ "github.com/micro/micro/v2/client/web"
	_ "github.com/micro/micro/v2/platform/cli"
	_ "github.com/micro/micro/v2/server"
	_ "github.com/micro/micro/v2/service/auth/cli"
	_ "github.com/micro/micro/v2/service/cli"
	_ "github.com/micro/micro/v2/service/config/cli"
	_ "github.com/micro/micro/v2/service/debug/cli"
	_ "github.com/micro/micro/v2/service/network/cli"
	_ "github.com/micro/micro/v2/service/runtime/cli"
	_ "github.com/micro/micro/v2/service/store/cli"
)

func main() {
	if err := cmd.DefaultCmd.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
