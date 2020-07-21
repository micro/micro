package main

import (
	"github.com/micro/micro/v2/cmd"

	// load cli packages so they can register commands
	_ "github.com/micro/micro/v2/service/auth/cli"
	_ "github.com/micro/micro/v2/service/cli"
	_ "github.com/micro/micro/v2/service/config/cli"
	_ "github.com/micro/micro/v2/service/debug/cli"
	_ "github.com/micro/micro/v2/service/network/cli"
	_ "github.com/micro/micro/v2/service/runtime/cli"
	_ "github.com/micro/micro/v2/service/store/cli"
)

func main() {
	cmd.Run()
}
