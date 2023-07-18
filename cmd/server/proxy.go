package server

import (
	"github.com/urfave/cli/v2"
	bmem "micro.dev/v4/service/broker/memory"
	"micro.dev/v4/service/client"
	log "micro.dev/v4/service/logger"
	"micro.dev/v4/service/registry/noop"
	"micro.dev/v4/service/server"
	sgrpc "micro.dev/v4/service/server/grpc"
	"micro.dev/v4/util/proxy"
	"micro.dev/v4/util/proxy/grpc"
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
