package cmd

import (
	"github.com/micro/micro/v3/service/auth"
	authSrv "github.com/micro/micro/v3/service/auth/client"
	"github.com/micro/micro/v3/service/client"
	grpcCli "github.com/micro/micro/v3/service/client/grpc"
	"github.com/micro/micro/v3/service/network"
	mucpNet "github.com/micro/micro/v3/service/network/mucp"
	"github.com/micro/micro/v3/service/server"
	grpcSvr "github.com/micro/micro/v3/service/server/grpc"
)

// setupDefaults sets the default auth, broker etc implementations incase they werent configured by
// a profile. The default implementations are always the RPC implementations.
func setupDefaults() {
	if client.DefaultClient == nil {
		client.DefaultClient = grpcCli.NewClient()
	}
	if server.DefaultServer == nil {
		server.DefaultServer = grpcSvr.NewServer()
	}
	if network.DefaultNetwork == nil {
		network.DefaultNetwork = mucpNet.NewNetwork()
	}

	// setup rpc implementations after the client is configured
	if auth.DefaultAuth == nil {
		auth.DefaultAuth = authSrv.NewAuth()
	}
}
