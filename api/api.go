package api

import (
	log "github.com/golang/glog"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
)

// API interface represents our API server.
// It should provide the facility to server HTTP requests,
// forward on to the appropriate services and return a
// response.
type API interface {
	Address() string
	Init() error
	Start() error
	Stop() error
}

var (
	Address      = ":8080"
	RpcPath      = "/rpc"
	HttpPath     = "/"
	Namespace    = "go.micro.api"
	HeaderPrefix = "X-Micro-"
)

func run() {
	// Init API
	api := New(Address)
	api.Init()

	// Initialise Server
	service := micro.NewService(
		micro.Name("go.micro.api"),
	)

	// Start API
	if err := api.Start(); err != nil {
		log.Fatal(err)
	}

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

	// Stop API
	if err := api.Stop(); err != nil {
		log.Fatal(err)
	}
}

func New(address string) API {
	return newApiServer(address)
}

func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:  "api",
			Usage: "Run the micro API",
			Action: func(c *cli.Context) {
				run()
			},
		},
	}
}
