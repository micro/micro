package tunnel

import (
	"os"
	"strings"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/config/options"
	"github.com/micro/go-micro/proxy"
	"github.com/micro/go-micro/proxy/mucp"
	"github.com/micro/go-micro/registry/memory"
	"github.com/micro/go-micro/router"
	"github.com/micro/go-micro/server"
	tun "github.com/micro/go-micro/tunnel"
	"github.com/micro/go-micro/tunnel/transport"
	"github.com/micro/go-micro/util/log"
)

var (
	// Name of the tunnel service
	Name = "go.micro.tunnel"
	// Address is the tunnel address
	Address = ":9095"
	// Tunnel is the name of the tunnel
	Tunnel = ":9096"
)

// run runs the micro server
func run(ctx *cli.Context, srvOpts ...micro.Option) {
	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	if len(ctx.GlobalString("server_name")) > 0 {
		Name = ctx.GlobalString("server_name")
	}
	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}
	if len(ctx.String("name")) > 0 {
		Tunnel = ctx.String("name")
	}
	var nodes []string
	if len(ctx.String("nodes")) > 0 {
		nodes = strings.Split(ctx.String("nodes"), ",")
	}

	// Initialise service
	service := micro.NewService(
		micro.Name(Name),
		micro.RegisterTTL(time.Duration(ctx.GlobalInt("register_ttl"))*time.Second),
		micro.RegisterInterval(time.Duration(ctx.GlobalInt("register_interval"))*time.Second),
	)

	// local tunnel router
	r := router.NewRouter(
		router.Id(service.Server().Options().Id),
		router.Registry(service.Client().Options().Registry),
	)

	// create a tunnel
	t := tun.NewTunnel(
		tun.Address(Address),
		tun.Nodes(nodes...),
	)

	// create tunnel client with tunnel transport
	tunTransport := transport.NewTransport(
		transport.WithTunnel(t),
	)

	// local server client talks to tunnel
	localSrvClient := client.NewClient(
		client.Transport(tunTransport),
	)

	// local proxy
	localProxy := mucp.NewProxy(
		proxy.WithRouter(r),
		proxy.WithClient(localSrvClient),
		proxy.WithEndpoint(Tunnel),
	)

	// init server
	service.Server().Init(
		server.WithRouter(localProxy),
	)

	// local transport client
	tunSrvClient := client.NewClient(
		client.Transport(service.Options().Transport),
	)

	// local proxy
	tunProxy := mucp.NewProxy(
		proxy.WithClient(tunSrvClient),
	)

	// create memory registry
	memRegistry := memory.NewRegistry()

	// local server
	tunSrv := server.NewServer(
		server.Address(Tunnel),
		server.Transport(tunTransport),
		server.WithRouter(tunProxy),
		server.Registry(memRegistry),
	)

	if err := tunSrv.Start(); err != nil {
		log.Logf("Tunnel error starting tunnel server: %v", err)
		os.Exit(1)
	}

	if err := service.Run(); err != nil {
		log.Log("Tunnel %s failed: %v", Name, err)
	}

	// stop the router
	if err := r.Stop(); err != nil {
		log.Logf("Tunnel error stopping tunnel router: %v", err)
	}

	// stop the server
	if err := tunSrv.Stop(); err != nil {
		log.Logf("Tunnel error stopping tunnel server: %v", err)
	}

	if err := t.Connect(); err != nil {
		log.Logf("Tunnel error stopping tunnel: %v", err)
	}

	log.Logf("Tunnel stopped")
}

func Commands(options ...micro.Option) []cli.Command {
	command := cli.Command{
		Name:  "tunnel",
		Usage: "Run the micro network tunnel",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "address",
				Usage:  "Set the micro tunnel address :9095",
				EnvVar: "MICRO_TUNNEL_ADDRESS",
			},
			cli.StringFlag{
				Name:   "name",
				Usage:  "Name of the tunnel used as the internal dial/listen address",
				EnvVar: "MICRO_TUNNEL_NAME",
			},
			cli.StringFlag{
				Name:   "nodes",
				Usage:  "Set the micro tunnel nodes",
				EnvVar: "MICRO_TUNNEL_NODES",
			},
		},
		Action: func(ctx *cli.Context) {
			run(ctx, options...)
		},
	}

	for _, p := range Plugins() {
		if cmds := p.Commands(); len(cmds) > 0 {
			command.Subcommands = append(command.Subcommands, cmds...)
		}

		if flags := p.Flags(); len(flags) > 0 {
			command.Flags = append(command.Flags, flags...)
		}
	}

	return []cli.Command{command}
}
