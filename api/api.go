package api

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	log "github.com/golang/glog"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/micro/internal/handler"
	"github.com/micro/micro/internal/server"
)

var (
	Address      = ":8080"
	RPCPath      = "/rpc"
	APIPath      = "/"
	Namespace    = "go.micro.api"
	HeaderPrefix = "X-Micro-"
)

func run(ctx *cli.Context) {
	// Init API
	var opts []server.Option

	if ctx.GlobalBool("enable_tls") {
		cert := ctx.GlobalString("tls_cert_file")
		key := ctx.GlobalString("tls_key_file")

		if len(cert) > 0 && len(key) > 0 {
			certs, err := tls.LoadX509KeyPair(cert, key)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			config := &tls.Config{
				Certificates: []tls.Certificate{certs},
			}
			opts = append(opts, server.EnableTLS(true))
			opts = append(opts, server.TLSConfig(config))
		} else {
			fmt.Println("Enable TLS specified without certificate and key files")
			return
		}
	}

	api := server.NewServer(Address)
	api.Init(opts...)

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

func New(address string) server.Server {
	return server.NewServer(address)
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
