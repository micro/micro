package api

import (
	"net/http"
	"time"

	log "github.com/golang/glog"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/micro/internal/handler"
)

// API interface represents our API server.
// It should provide the facility to server HTTP requests,
// forward on to the appropriate services and return a
// response.
type API interface {
	Address() string
	Init() error
	Handle(path string, handler http.Handler)
	Start() error
	Stop() error
}

var (
	Address      = ":8080"
	RPCPath      = "/rpc"
	APIPath      = "/"
	Namespace    = "go.micro.api"
	HeaderPrefix = "X-Micro-"
)

func run(ctx *cli.Context) {
	// Init API
	api := New(Address)
	api.Init()

	log.Infof("Registering RPC Handler at %s", RPCPath)
	api.Handle(RPCPath, http.HandlerFunc(handler.RPC))
	log.Infof("Registering API Handler at %s", APIPath)
	api.Handle(APIPath, http.HandlerFunc(restHandler))

	// Initialise Server
	service := micro.NewService(
		micro.Name("go.micro.api"),
		micro.RegisterTTL(
			time.Duration(ctx.GlobalInt("register_ttl"))*time.Second,
		),
		micro.RegisterInterval(
			time.Duration(ctx.GlobalInt("register_interval"))*time.Second,
		),
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
			Name:   "api",
			Usage:  "Run the micro API",
			Action: run,
		},
	}
}
