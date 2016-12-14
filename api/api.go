// Package api is an API Gateway
package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/micro/internal/handler"
	"github.com/micro/micro/internal/helper"
	"github.com/micro/micro/internal/server"
	"github.com/micro/micro/internal/stats"
	"github.com/micro/micro/plugin"
)

var (
	Address      = ":8080"
	Handler      = "api"
	RPCPath      = "/rpc"
	APIPath      = "/"
	ProxyPath    = "/{service:[a-zA-Z0-9]+}"
	Namespace    = "go.micro.api"
	HeaderPrefix = "X-Micro-"
	CORS         = map[string]bool{"*": true}
)

type srv struct {
	*mux.Router
}

func (s *srv) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); CORS[origin] {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	} else if len(origin) > 0 && CORS["*"] {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == "OPTIONS" {
		return
	}

	s.Router.ServeHTTP(w, r)
}

func run(ctx *cli.Context) {
	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}
	if len(ctx.String("handler")) > 0 {
		Handler = ctx.String("handler")
	}
	if len(ctx.String("namespace")) > 0 {
		Namespace = ctx.String("namespace")
	}
	if len(ctx.String("cors")) > 0 {
		origins := make(map[string]bool)
		for _, origin := range strings.Split(ctx.String("cors"), ",") {
			origins[origin] = true
		}
		CORS = origins
	}

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	// Init API
	var opts []server.Option

	if ctx.GlobalBool("enable_tls") {
		config, err := helper.TLSConfig(ctx)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		opts = append(opts, server.EnableTLS(true))
		opts = append(opts, server.TLSConfig(config))
	}

	// create the router
	r := mux.NewRouter()
	s := &srv{r}
	var h http.Handler = s

	if ctx.GlobalBool("enable_stats") {
		st := stats.New()
		r.HandleFunc("/stats", st.StatsHandler)
		h = st.ServeHTTP(s)
		st.Start()
		defer st.Stop()
	}

	log.Printf("Registering RPC Handler at %s", RPCPath)
	r.HandleFunc(RPCPath, handler.RPC)

	switch Handler {
	case "rpc":
		log.Printf("Registering API RPC Handler at %s", APIPath)
		r.PathPrefix(APIPath).Handler(handler.RPCX(Namespace))
	case "proxy":
		log.Printf("Registering API Proxy Handler at %s", ProxyPath)
		r.PathPrefix(ProxyPath).Handler(handler.Proxy(Namespace, false))
	default:
		log.Printf("Registering API Default Handler at %s", APIPath)
		r.PathPrefix(APIPath).Handler(handler.API(Namespace))
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

func Commands() []cli.Command {
	command := cli.Command{
		Name:   "api",
		Usage:  "Run the micro API",
		Action: run,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "address",
				Usage:  "Set the api address e.g 0.0.0.0:8080",
				EnvVar: "MICRO_API_ADDRESS",
			},
			cli.StringFlag{
				Name:   "handler",
				Usage:  "Specify the request handler to be used for mapping HTTP requests to services. e.g api, proxy, rpc",
				EnvVar: "MICRO_API_HANDLER",
			},
			cli.StringFlag{
				Name:   "namespace",
				Usage:  "Set the namespace used by the API e.g. com.example.api",
				EnvVar: "MICRO_API_NAMESPACE",
			},
			cli.StringFlag{
				Name:   "cors",
				Usage:  "Comma separated whitelist of allowed origins for CORS",
				EnvVar: "MICRO_API_CORS",
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
