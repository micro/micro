package main

import (
	"code.google.com/p/go.net/context"

	log "github.com/golang/glog"
	"github.com/myodc/go-micro/cmd"
	"github.com/myodc/go-micro/server"
	hello "github.com/myodc/micro/examples/greeter/server/proto/hello"
)

type Say struct{}

func (s *Say) Hello(ctx context.Context, req *hello.Request, rsp *hello.Response) error {
	log.Info("Received Say.Hello request")

	rsp.Msg = server.Config().Id() + ": Hello " + req.Name

	return nil
}

func main() {
	// optionally setup command line usage
	cmd.Init()

	// Initialise Server
	server.Init(
		server.Name("go.micro.srv.greeter"),
	)

	// Register Handlers
	server.Register(
		server.NewReceiver(
			new(Say),
		),
	)

	// Run server
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
