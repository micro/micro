package service

import (
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/cmd"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/debug/trace"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-micro/v2/store"

	// clients
	gcli "github.com/micro/go-micro/v2/client/grpc"

	gsrv "github.com/micro/go-micro/v2/server/grpc"

	// brokers

	brokerSrv "github.com/micro/go-micro/v2/broker/service"

	// registries

	regSrv "github.com/micro/go-micro/v2/registry/service"

	// routers

	srvRouter "github.com/micro/go-micro/v2/router/service"

	// runtimes

	srvRuntime "github.com/micro/go-micro/v2/runtime/service"

	// selectors

	// transports

	// stores
	memStore "github.com/micro/go-micro/v2/store/memory"
	svcStore "github.com/micro/go-micro/v2/store/service"

	// tracers
	// jTracer "github.com/micro/go-micro/v2/debug/trace/jaeger"
	memTracer "github.com/micro/go-micro/v2/debug/trace/memory"

	// auth

	svcAuth "github.com/micro/go-micro/v2/auth/service"
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
