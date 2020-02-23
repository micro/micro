package tunnel

import (
	"os"
	"strings"
	"time"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	cmucp "github.com/micro/go-micro/v2/client/mucp"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/proxy"
	"github.com/micro/go-micro/v2/proxy/mucp"
	"github.com/micro/go-micro/v2/registry/memory"
	"github.com/micro/go-micro/v2/router"
	"github.com/micro/go-micro/v2/server"
	smucp "github.com/micro/go-micro/v2/server/mucp"
	tun "github.com/micro/go-micro/v2/tunnel"
	"github.com/micro/go-micro/v2/tunnel/transport"
	"github.com/micro/go-micro/v2/util/mux"
)

var (
	// Name of the tunnel service
	Name = "go.micro.tunnel"
	// Address is the tunnel address
	Address = ":8083"
	// Tunnel is the name of the tunnel
	Tunnel = "tun:0"
	// The tunnel token
	Token = "micro"
)

// run runs the micro server
func run(ctx *cli.Context, srvOpts ...micro.Option) {
	log.Info("tunnel")

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	if len(ctx.String("server_name")) > 0 {
		Name = ctx.String("server_name")
	}
	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}
	if len(ctx.String("token")) > 0 {
		Token = ctx.String("token")
	}
	if len(ctx.String("id")) > 0 {
		Tunnel = ctx.String("id")
		// We need host:port for the Endpoint value in the proxy
		parts := strings.Split(Tunnel, ":")
		if len(parts) == 1 {
			Tunnel = Tunnel + ":0"
		}
	}
	var nodes []string
	if len(ctx.String("server")) > 0 {
		nodes = strings.Split(ctx.String("server"), ",")
	}

	// Initialise service
	service := micro.NewService(
		micro.Name(Name),
		micro.RegisterTTL(time.Duration(ctx.Int("register_ttl"))*time.Second),
		micro.RegisterInterval(time.Duration(ctx.Int("register_interval"))*time.Second),
	)

	// local tunnel router
	r := router.NewRouter(
		router.Id(service.Server().Options().Id),
		router.Registry(service.Client().Options().Registry),
	)

	// start the router
	if err := r.Start(); err != nil {
		log.Errorf("Tunnel error starting router: %s", err)
		os.Exit(1)
	}

	// create a tunnel
	t := tun.NewTunnel(
		tun.Address(Address),
		tun.Nodes(nodes...),
		tun.Token(Token),
	)

	// start the tunnel
	if err := t.Connect(); err != nil {
		log.Errorf("Tunnel error connecting: %v", err)
	}

	log.Infof("Tunnel [%s] listening on %s", Tunnel, Address)

	// create tunnel client with tunnel transport
	tunTransport := transport.NewTransport(
		transport.WithTunnel(t),
	)

	// local server client talks to tunnel
	localSrvClient := cmucp.NewClient(
		client.Transport(tunTransport),
	)

	// local proxy
	localProxy := mucp.NewProxy(
		proxy.WithClient(localSrvClient),
		proxy.WithEndpoint(Tunnel),
	)

	// create new muxer
	muxer := mux.New(Name, localProxy)

	// init server
	service.Server().Init(
		server.WithRouter(muxer),
	)

	// local transport client
	service.Client().Init(
		client.Transport(service.Options().Transport),
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
		server.Address(Tunnel),
		server.Transport(tunTransport),
		server.WithRouter(tunProxy),
		server.Registry(memRegistry),
	)

	if err := tunSrv.Start(); err != nil {
		log.Errorf("Tunnel error starting tunnel server: %v", err)
		os.Exit(1)
	}

	if err := service.Run(); err != nil {
		log.Errorf("Tunnel %s failed: %v", Name, err)
	}

	// stop the router
	if err := r.Stop(); err != nil {
		log.Errorf("Tunnel error stopping tunnel router: %v", err)
	}

	// stop the server
	if err := tunSrv.Stop(); err != nil {
		log.Errorf("Tunnel error stopping tunnel server: %v", err)
	}

	if err := t.Close(); err != nil {
		log.Errorf("Tunnel error stopping tunnel: %v", err)
	}
}

func Commands(options ...micro.Option) []*cli.Command {
	command := &cli.Command{
		Name:  "tunnel",
		Usage: "Run the micro network tunnel",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "address",
				Usage:   "Set the micro tunnel address :8083",
				EnvVars: []string{"MICRO_TUNNEL_ADDRESS"},
			},
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
		},
		Action: func(ctx *cli.Context) error {
			run(ctx, options...)
			return nil
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

	return []*cli.Command{command}
}
