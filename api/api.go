// Package api is an API Gateway
package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/micro/cli"
	ahandler "github.com/micro/go-api/handler"
	aapi "github.com/micro/go-api/handler/api"
	"github.com/micro/go-api/handler/event"
	ahttp "github.com/micro/go-api/handler/http"
	arpc "github.com/micro/go-api/handler/rpc"
	"github.com/micro/go-api/handler/web"
	"github.com/micro/go-api/resolver"
	"github.com/micro/go-api/resolver/grpc"
	"github.com/micro/go-api/resolver/host"
	rrmicro "github.com/micro/go-api/resolver/micro"
	"github.com/micro/go-api/resolver/path"
	"github.com/micro/go-api/router"
	"github.com/micro/go-api/server"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"github.com/micro/micro/internal/handler"
	"github.com/micro/micro/internal/helper"
	"github.com/micro/micro/internal/stats"
	"github.com/micro/micro/plugin"
)

var (
	Name         = "go.micro.api"
	Address      = ":8080"
	Handler      = "meta"
	Resolver     = "micro"
	RPCPath      = "/rpc"
	APIPath      = "/"
	ProxyPath    = "/{service:[a-zA-Z0-9]+}"
	Namespace    = "go.micro.api"
	HeaderPrefix = "X-Micro-"
)

func run(ctx *cli.Context, srvOpts ...micro.Option) {
	if len(ctx.GlobalString("server_name")) > 0 {
		Name = ctx.GlobalString("server_name")
	}
	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}
	if len(ctx.String("handler")) > 0 {
		Handler = ctx.String("handler")
	}
	if len(ctx.String("namespace")) > 0 {
		Namespace = ctx.String("namespace")
	}
	if len(ctx.String("resolver")) > 0 {
		Resolver = ctx.String("resolver")
	}

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	// Init API
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

	// create the router
	var h http.Handler
	r := mux.NewRouter()
	h = r

	if ctx.GlobalBool("enable_stats") {
		st := stats.New()
		r.HandleFunc("/stats", st.StatsHandler)
		h = st.ServeHTTP(r)
		st.Start()
		defer st.Stop()
	}

	srvOpts = append(srvOpts, micro.Name(Name))
	if i := time.Duration(c.GlobalInt("register_ttl")); i > 0 {
		srvOpts = append(srvOpts, micro.RegisterTTL(i*time.Second))
	}
	if i := time.Duration(c.GlobalInt("register_interval")); i > 0 {
		srvOpts = append(srvOpts, micro.RegisterInterval(i*time.Second))
	}

	// initialise service
	service := micro.NewService(srvOpts...)

	// register rpc handler
	log.Logf("Registering RPC Handler at %s", RPCPath)
	r.HandleFunc(RPCPath, handler.RPC)

	// resolver options
	ropts := []resolver.Option{
		resolver.WithNamespace(Namespace),
		resolver.WithHandler(Handler),
	}

	// default resolver
	rr := rrmicro.NewResolver(ropts...)

	switch Resolver {
	case "host":
		rr = host.NewResolver(ropts...)
	case "path":
		rr = path.NewResolver(ropts...)
	case "grpc":
		rr = grpc.NewResolver(ropts...)
	}

	switch Handler {
	case "rpc":
		log.Logf("Registering API RPC Handler at %s", APIPath)
		rt := router.NewRouter(
			router.WithNamespace(Namespace),
			router.WithHandler(arpc.Handler),
			router.WithResolver(rr),
			router.WithRegistry(service.Options().Registry),
		)
		rp := arpc.NewHandler(
			ahandler.WithNamespace(Namespace),
			ahandler.WithRouter(rt),
			ahandler.WithService(service),
		)
		r.PathPrefix(APIPath).Handler(rp)
	case "api":
		log.Logf("Registering API Request Handler at %s", APIPath)
		rt := router.NewRouter(
			router.WithNamespace(Namespace),
			router.WithHandler(aapi.Handler),
			router.WithResolver(rr),
			router.WithRegistry(service.Options().Registry),
		)
		ap := aapi.NewHandler(
			ahandler.WithNamespace(Namespace),
			ahandler.WithRouter(rt),
			ahandler.WithService(service),
		)
		r.PathPrefix(APIPath).Handler(ap)
	case "event":
		log.Logf("Registering API Event Handler at %s", APIPath)
		rt := router.NewRouter(
			router.WithNamespace(Namespace),
			router.WithHandler(event.Handler),
			router.WithResolver(rr),
			router.WithRegistry(service.Options().Registry),
		)
		ev := event.NewHandler(
			ahandler.WithNamespace(Namespace),
			ahandler.WithRouter(rt),
			ahandler.WithService(service),
		)
		r.PathPrefix(APIPath).Handler(ev)
	case "http", "proxy":
		log.Logf("Registering API HTTP Handler at %s", ProxyPath)
		rt := router.NewRouter(
			router.WithNamespace(Namespace),
			router.WithHandler(ahttp.Handler),
			router.WithResolver(rr),
			router.WithRegistry(service.Options().Registry),
		)
		ht := ahttp.NewHandler(
			ahandler.WithNamespace(Namespace),
			ahandler.WithRouter(rt),
			ahandler.WithService(service),
		)
		r.PathPrefix(ProxyPath).Handler(ht)
	case "web":
		log.Logf("Registering API Web Handler at %s", APIPath)
		rt := router.NewRouter(
			router.WithNamespace(Namespace),
			router.WithHandler(web.Handler),
			router.WithResolver(rr),
			router.WithRegistry(service.Options().Registry),
		)
		w := web.NewHandler(
			ahandler.WithNamespace(Namespace),
			ahandler.WithRouter(rt),
			ahandler.WithService(service),
		)
		r.PathPrefix(APIPath).Handler(w)
	default:
		log.Logf("Registering API Default Handler at %s", APIPath)
		rt := router.NewRouter(
			router.WithNamespace(Namespace),
			router.WithRegistry(service.Options().Registry),
		)
		r.PathPrefix(APIPath).Handler(handler.Meta(service, rt))
	}

	// reverse wrap handler
	plugins := append(Plugins(), plugin.Plugins()...)
	for i := len(plugins); i > 0; i-- {
		h = plugins[i-1].Handler()(h)
	}

	// create the server
	api := server.NewServer(Address)
	api.Init(opts...)
	api.Handle("/", h)

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

func Commands(options ...micro.Option) []cli.Command {
	command := cli.Command{
		Name:   "api",
		Usage:  "Run the micro API",
		Action: func(ctx *cli.Context) {
			run(ctx, options...)
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "address",
				Usage:  "Set the api address e.g 0.0.0.0:8080",
				EnvVar: "MICRO_API_ADDRESS",
			},
			cli.StringFlag{
				Name:   "handler",
				Usage:  "Specify the request handler to be used for mapping HTTP requests to services; {api, event, http, rpc}",
				EnvVar: "MICRO_API_HANDLER",
			},
			cli.StringFlag{
				Name:   "namespace",
				Usage:  "Set the namespace used by the API e.g. com.example.api",
				EnvVar: "MICRO_API_NAMESPACE",
			},
			cli.StringFlag{
				Name:   "resolver",
				Usage:  "Set the hostname resolver used by the API {host, path, grpc}",
				EnvVar: "MICRO_API_RESOLVER",
			},
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
