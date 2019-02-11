// Package proxy is a cli proxy
package proxy

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/micro/cli"
	"github.com/micro/go-api/server"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"github.com/micro/go-proxy/router/mucp"
	"github.com/micro/micro/internal/handler"
	"github.com/micro/micro/internal/helper"
	"github.com/micro/micro/internal/stats"
	"github.com/micro/micro/plugin"

	ahandler "github.com/micro/go-api/handler"
	abroker "github.com/micro/go-api/handler/broker"
	aregistry "github.com/micro/go-api/handler/registry"
)

type srv struct {
	*mux.Router
}

var (
	// Name of the proxy
	Name = "go.micro.proxy"
	// The http address of the proxy
	Address = ":8081"
	// The backend host to route to
	Backend string
	// The paths for http endpoints
	BrokerPath   = "/broker"
	RegistryPath = "/registry"
	RPCPath      = "/rpc"
)

func run(ctx *cli.Context, srvOpts ...micro.Option) {
	if len(ctx.GlobalString("server_name")) > 0 {
		Name = ctx.GlobalString("server_name")
	}
	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}
	if len(ctx.String("backend")) > 0 {
		Backend = ctx.String("backend")
	}

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	var opts []server.Option

	if ctx.GlobalBool("enable_acme") {
		hosts := helper.ACMEHosts(ctx)
		opts = append(opts, server.EnableACME(true))
		opts = append(opts, server.ACMEHosts(hosts...))
	} else if ctx.GlobalBool("enable_tls") {
		config, err := helper.TLSConfig(ctx)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		opts = append(opts, server.EnableTLS(true))
		opts = append(opts, server.TLSConfig(config))
	}

	r := mux.NewRouter()
	var h http.Handler = r

	if ctx.GlobalBool("enable_stats") {
		st := stats.New()
		r.Handle("/stats", http.HandlerFunc(st.StatsHandler))
		h = st.ServeHTTP(r)
		st.Start()
		defer st.Stop()
	}

	// new server
	srv := server.NewServer(Address)
	srv.Init(opts...)

	// service opts
	srvOpts = append(srvOpts, micro.Name(Name))
	if i := time.Duration(ctx.GlobalInt("register_ttl")); i > 0 {
		srvOpts = append(srvOpts, micro.RegisterTTL(i*time.Second))
	}
	if i := time.Duration(ctx.GlobalInt("register_interval")); i > 0 {
		srvOpts = append(srvOpts, micro.RegisterInterval(i*time.Second))
	}

	// set backend
	if len(Backend) > 0 {
		srvOpts = append(srvOpts, mucp.WithBackend(Backend))
	}

	// Initialise Server
	service := mucp.NewService(srvOpts...)

	log.Logf("Registering Registry handler at %s", RegistryPath)
	r.Handle(RegistryPath, aregistry.NewHandler(ahandler.WithService(service)))

	log.Logf("Registering RPC handler at %s", RPCPath)
	r.Handle(RPCPath, http.HandlerFunc(handler.RPC))

	log.Logf("Registering Broker handler at %s", BrokerPath)
	br := abroker.NewHandler(
		ahandler.WithService(service),
	)
	r.Handle(BrokerPath, br)

	// reverse wrap handler
	plugins := append(Plugins(), plugin.Plugins()...)
	for i := len(plugins); i > 0; i-- {
		h = plugins[i-1].Handler()(h)
	}

	srv.Handle("/", h)

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

	if err := srv.Stop(); err != nil {
		log.Fatal(err)
	}
}

func Commands(options ...micro.Option) []cli.Command {
	command := cli.Command{
		Name:  "proxy",
		Usage: "Run the service proxy",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "address",
				Usage:  "Set the proxy http address e.g 0.0.0.0:8081",
				EnvVar: "MICRO_PROXY_ADDRESS",
			},
			cli.StringFlag{
				Name:   "backend",
				Usage:  "Set the backend to route to e.g greeter or localhost:9090",
				EnvVar: "MICRO_PROXY_BACKEND",
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
