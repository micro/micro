// Package proxy is a cli proxy
package proxy

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/micro/cli"
	"github.com/micro/go-api/server"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
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
	Name         = "go.micro.proxy"
	Address      = ":8081"
	BrokerPath   = "/broker"
	RegistryPath = "/registry"
	RPCPath      = "/rpc"
	CORS         = map[string]bool{"*": true}
)

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
	s := &srv{r}

	var h http.Handler = s

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

	// Initialise Server
	service := micro.NewService(
		micro.Name(Name),
		micro.RegisterTTL(
			time.Duration(ctx.GlobalInt("register_ttl"))*time.Second,
		),
		micro.RegisterInterval(
			time.Duration(ctx.GlobalInt("register_interval"))*time.Second,
		),
	)

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

func Commands() []cli.Command {
	command := cli.Command{
		Name:  "proxy",
		Usage: "Run the micro proxy",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "address",
				Usage:  "Set the sidecar address e.g 0.0.0.0:8081",
				EnvVar: "MICRO_SIDECAR_ADDRESS",
			},
			cli.StringFlag{
				Name:   "cors",
				Usage:  "Comma separated whitelist of allowed origins for CORS",
				EnvVar: "MICRO_SIDECAR_CORS",
			},
		},
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
