package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/micro/internal/handler"
	"github.com/micro/micro/internal/helper"
	"github.com/micro/micro/internal/server"
	"github.com/micro/micro/internal/stats"
)

var (
	Address      = ":8080"
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

	switch ctx.GlobalString("api_handler") {
	case "proxy":
		log.Printf("Registering API Proxy Handler at %s", ProxyPath)
		r.PathPrefix(ProxyPath).Handler(handler.Proxy(Namespace, false))
	default:
		log.Printf("Registering API Handler at %s", APIPath)
		r.PathPrefix(APIPath).HandlerFunc(apiHandler)
	}

	// reverse wrap handler
	plugins := Plugins()
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

func New(address string) server.Server {
	return server.NewServer(address)
}

func Commands() []cli.Command {
	command := cli.Command{
		Name:   "api",
		Usage:  "Run the micro API",
		Action: run,
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
