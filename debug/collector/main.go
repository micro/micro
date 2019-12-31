package main

import (
	"os"
	"path"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/util/log"
	plugin "github.com/micro/micro/debug/collector/micro"

	"github.com/netdata/go-orchestrator"
	"github.com/netdata/go-orchestrator/cli"
	"github.com/netdata/go-orchestrator/pkg/multipath"
)

var (
	cd, _         = os.Getwd()
	netdataConfig = multipath.New(
		os.Getenv("NETDATA_USER_CONFIG_DIR"),
		os.Getenv("NETDATA_STOCK_CONFIG_DIR"),
		path.Join(cd, "/../../../../etc/netdata"),
		path.Join(cd, "/../../../../usr/lib/netdata/conf.d"),
	)
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.debug.collector"),
		micro.Version("latest"),
	)

	if len(os.Args) > 1 {
		os.Args = append(os.Args[:1], os.Args[2:]...)
	}

	// Initialise service
	service.Init()

	go func() {
		log.Fatal(service.Run())
	}()

	// register the new plugin
	plugin.New(service.Client()).Register()

	netdata := orchestrator.New()
	netdata.Name = "micro.d"
	netdata.Option = &cli.Option{
		UpdateEvery: 1,
		Debug:       true,
		Module:      "all",
		ConfigDir:   netdataConfig,
		Version:     false,
	}
	netdata.ConfigPath = netdataConfig

	if !netdata.Setup() {
		log.Fatal("Netdata failed to Setup()")
	}

	netdata.Serve()
}
