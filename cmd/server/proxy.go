package server

import (
	bmem "github.com/micro/micro/v5/service/broker/memory"
	"github.com/micro/micro/v5/service/client"
	log "github.com/micro/micro/v5/service/logger"
	"github.com/micro/micro/v5/service/registry/noop"
	"github.com/micro/micro/v5/service/server"
	sgrpc "github.com/micro/micro/v5/service/server/grpc"
	"github.com/micro/micro/v5/util/proxy"
	"github.com/micro/micro/v5/util/proxy/grpc"
	"github.com/urfave/cli/v2"
)

func runProxy(ctx *cli.Context, wait chan bool) error {
	// set the context
	popts := []proxy.Option{
		proxy.WithClient(client.DefaultClient),
	}

	serverOpts := []server.Option{
		server.Address(Address),
		server.Registry(noop.NewRegistry()),
		server.Broker(bmem.NewBroker()),
	}

	// default to the grpc proxy
	p := grpc.NewProxy(popts...)

	// wrap the proxy using the proxy's authHandler
	authOpt := server.WrapHandler(authHandler())
	serverOpts = append(serverOpts, authOpt)
	serverOpts = append(serverOpts, server.WithRouter(p))

	// create a new grpc server
	srv := sgrpc.NewServer(serverOpts...)

	// Start the proxy server
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}

	<-wait

	// Stop the server
	if err := srv.Stop(); err != nil {
		log.Fatal(err)
	}

	return nil
}
