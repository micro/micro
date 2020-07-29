package server

import (
	"os"
	"strings"
	"time"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v3/client"
	cmucp "github.com/micro/go-micro/v3/client/mucp"
	log "github.com/micro/go-micro/v3/logger"
	"github.com/micro/go-micro/v3/proxy"
	"github.com/micro/go-micro/v3/proxy/mucp"
	"github.com/micro/go-micro/v3/registry/memory"
	"github.com/micro/go-micro/v3/router"
	regRouter "github.com/micro/go-micro/v3/router/registry"
	"github.com/micro/go-micro/v3/server"
	smucp "github.com/micro/go-micro/v3/server/mucp"
	"github.com/micro/go-micro/v3/transport"
	tun "github.com/micro/go-micro/v3/tunnel"
	tuntransport "github.com/micro/go-micro/v3/tunnel/transport"
	"github.com/micro/go-micro/v3/util/mux"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/registry"
	mutunnel "github.com/micro/micro/v3/service/tunnel"
)

var (
	// name of the tunnel service
	name = "go.micro.tunnel"
	// address is the tunnel address
	address = ":8083"
	// tunnel is the name of the tunnel
	tunnel = "tun:0"
	// the tunnel token
	token = "micro"

	// Flags specific to the tunnel service
	Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "id",
			Usage:   "Id of the tunnel used as the internal dial/listen address.",
			EnvVars: []string{"MICRO_TUNNEL_ID"},
		},
		&cli.StringFlag{
			Name:    "server",
			Usage:   "Set the micro tunnel server address. This can be a comma separated list.",
			EnvVars: []string{"MICRO_TUNNEL_SERVER"},
		},
		&cli.StringFlag{
			Name:    "token",
			Usage:   "Set the micro tunnel token for authentication",
			EnvVars: []string{"MICRO_TUNNEL_TOKEN"},
		},
	}
)

// Run micro tunnel
func Run(ctx *cli.Context) error {
	if len(ctx.String("server_name")) > 0 {
		name = ctx.String("server_name")
	}
	if len(ctx.String("address")) > 0 {
		address = ctx.String("address")
	}
	if len(ctx.String("tunnel_address")) > 0 {
		address = ctx.String("tunnel_address")
	}
	if len(ctx.String("token")) > 0 {
		token = ctx.String("token")
	}
	if len(ctx.String("id")) > 0 {
		tunnel = ctx.String("id")
		// We need host:port for the Endpoint value in the proxy
		parts := strings.Split(tunnel, ":")
		if len(parts) == 1 {
			tunnel = tunnel + ":0"
		}
	}
	var nodes []string
	if len(ctx.String("server")) > 0 {
		nodes = strings.Split(ctx.String("server"), ",")
	}

	// Initialise service
	service := service.New(
		service.Name(name),
		service.RegisterTTL(time.Duration(ctx.Int("register_ttl"))*time.Second),
		service.RegisterInterval(time.Duration(ctx.Int("register_interval"))*time.Second),
	)

	// local tunnel router
	r := regRouter.NewRouter(
		router.Id(service.Server().Options().Id),
		router.Registry(registry.DefaultRegistry),
	)

	// create a tunnel
	t := tun.NewTunnel(
		tun.Address(address),
		tun.Nodes(nodes...),
		tun.Token(token),
	)

	// start the tunnel
	if err := t.Connect(); err != nil {
		log.Errorf("Tunnel error connecting: %v", err)
	}

	log.Infof("Tunnel [%s] listening on %s", tunnel, address)

	// create tunnel client with tunnel transport
	tunTransport := transport.NewTransport(
		tuntransport.WithTunnel(t),
	)

	// local server client talks to tunnel
	localSrvClient := cmucp.NewClient(
		client.Transport(tunTransport),
	)

	// local proxy
	localProxy := mucp.NewProxy(
		proxy.WithClient(localSrvClient),
		proxy.WithEndpoint(tunnel),
	)

	// create new muxer
	muxer := mux.New(name, localProxy)

	// init server
	service.Server().Init(
		server.WithRouter(muxer),
	)

	// local transport client
	service.Client().Init(
		client.Transport(mutunnel.DefaultTransport),
	)

	// local proxy
	tunProxy := mucp.NewProxy(
		proxy.WithRouter(r),
		proxy.WithClient(service.Client()),
	)

	// create memory registry
	memRegistry := memory.NewRegistry()

	// local server
	tunSrv := smucp.NewServer(
		server.Address(tunnel),
		server.Transport(tunTransport),
		server.WithRouter(tunProxy),
		server.Registry(memRegistry),
	)

	if err := tunSrv.Start(); err != nil {
		log.Errorf("Tunnel error starting tunnel server: %v", err)
		os.Exit(1)
	}

	if err := service.Run(); err != nil {
		log.Errorf("Tunnel %s failed: %v", name, err)
	}

	// stop the router
	if err := r.Close(); err != nil {
		log.Errorf("Tunnel error closing tunnel router: %v", err)
	}

	// stop the server
	if err := tunSrv.Stop(); err != nil {
		log.Errorf("Tunnel error stopping tunnel server: %v", err)
	}

	if err := t.Close(); err != nil {
		log.Errorf("Tunnel error stopping tunnel: %v", err)
	}

	return nil
}
