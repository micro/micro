package cmd

import (
	"github.com/micro/micro/v3/service/auth"
	authSrv "github.com/micro/micro/v3/service/auth/client"
	"github.com/micro/micro/v3/service/broker"
	brokerSrv "github.com/micro/micro/v3/service/broker/client"
	"github.com/micro/micro/v3/service/client"
	grpcCli "github.com/micro/micro/v3/service/client/grpc"
	"github.com/micro/micro/v3/service/model"
	"github.com/micro/micro/v3/service/model/mud"
	"github.com/micro/micro/v3/service/events"
	eventsSrv "github.com/micro/micro/v3/service/events/client"
	"github.com/micro/micro/v3/service/network"
	mucpNet "github.com/micro/micro/v3/service/network/mucp"
	"github.com/micro/micro/v3/service/server"
	grpcSvr "github.com/micro/micro/v3/service/server/grpc"
)

// setupDefaults sets the default auth, broker etc implementations incase they arent configured by
// a profile. The default implementations are always the RPC implementations.
func setupDefaults() {
	client.DefaultClient = grpcCli.NewClient()
	server.DefaultServer = grpcSvr.NewServer()
	network.DefaultNetwork = mucpNet.NewNetwork()
	model.DefaultModel = mud.NewModel()

	// setup rpc implementations after the client is configured
	auth.DefaultAuth = authSrv.NewAuth()
	broker.DefaultBroker = brokerSrv.NewBroker()
	events.DefaultStream = eventsSrv.NewStream()
	events.DefaultStore = eventsSrv.NewStore()
}
