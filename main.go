package main

import (
	"github.com/micro/micro/v2/cmd"

	// services
	_ "github.com/micro/micro/v2/service"
	_ "github.com/micro/micro/v2/service/auth"
	_ "github.com/micro/micro/v2/service/broker"
	_ "github.com/micro/micro/v2/service/config"
	_ "github.com/micro/micro/v2/service/debug"
	_ "github.com/micro/micro/v2/service/health"
	_ "github.com/micro/micro/v2/service/network"
	_ "github.com/micro/micro/v2/service/registry"
	_ "github.com/micro/micro/v2/service/router"
	_ "github.com/micro/micro/v2/service/runtime"
	_ "github.com/micro/micro/v2/service/store"
	_ "github.com/micro/micro/v2/service/tunnel"
)

func main() {
	cmd.Init()
}
