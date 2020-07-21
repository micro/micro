package service

import (
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/cmd"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/debug/trace"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-micro/v2/store"

	svcAuth "github.com/micro/go-micro/v2/auth/service"
	brokerSrv "github.com/micro/go-micro/v2/broker/service"
	gcli "github.com/micro/go-micro/v2/client/grpc"
	memTracer "github.com/micro/go-micro/v2/debug/trace/memory"
	regSrv "github.com/micro/go-micro/v2/registry/service"
	srvRouter "github.com/micro/go-micro/v2/router/service"
	srvRuntime "github.com/micro/go-micro/v2/runtime/service"
	gsrv "github.com/micro/go-micro/v2/server/grpc"
	memStore "github.com/micro/go-micro/v2/store/memory"
	svcStore "github.com/micro/go-micro/v2/store/service"
)

func init() {
	// set defaults
	client.DefaultClient = gcli.NewClient()
	server.DefaultServer = gsrv.NewServer()
	store.DefaultStore = memStore.NewStore()
	trace.DefaultTracer = memTracer.NewTracer()

	// tmp: import services
	cmd.DefaultAuths["service"] = svcAuth.NewAuth
	cmd.DefaultBrokers["service"] = brokerSrv.NewBroker
	cmd.DefaultConfigs["service"] = config.NewConfig
	cmd.DefaultRegistries["service"] = regSrv.NewRegistry
	cmd.DefaultRuntimes["service"] = srvRuntime.NewRuntime
	cmd.DefaultRouters["service"] = srvRouter.NewRouter
	cmd.DefaultStores["service"] = svcStore.NewStore
}
