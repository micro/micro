// Package car is a language agnostic rpc proxy
package car

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/micro/cli"
	"github.com/micro/go-api"
	"github.com/micro/go-api/router"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/registry"
	"github.com/micro/micro/internal/handler"
	"github.com/micro/micro/internal/helper"
	"github.com/micro/micro/internal/server"
	"github.com/micro/micro/internal/stats"
	"github.com/micro/micro/plugin"
	"github.com/pborman/uuid"
)

type sidecar struct {
	name    string
	address string
	hcUrl   string
}

type srv struct {
	*mux.Router
}

var (
	Name         = "go.micro.sidecar"
	Address      = ":8081"
	Handler      = "rpc"
	RootPath     = "/"
	ProxyPath    = "/{service:[a-zA-Z0-9]+}"
	BrokerPath   = "/broker"
	HealthPath   = "/health"
	RegistryPath = "/registry"
	RPCPath      = "/rpc"
	CORS         = map[string]bool{"*": true}
	Namespace    = "go.micro.srv"
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

func newSidecar(name, address, hcUrl string) *sidecar {
	return &sidecar{
		name:    name,
		address: address,
		hcUrl:   hcUrl,
	}
}

func run(ctx *cli.Context, car *sidecar) {
	if len(ctx.GlobalString("server_name")) > 0 {
		Name = ctx.GlobalString("server_name")
	}
	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}
	if len(ctx.String("cors")) > 0 {
		origins := make(map[string]bool)
		for _, origin := range strings.Split(ctx.String("cors"), ",") {
			origins[origin] = true
		}
		CORS = origins
	}
	if len(ctx.String("handler")) > 0 {
		Handler = ctx.String("handler")
	}
	if len(ctx.String("namespace")) > 0 {
		Namespace = ctx.String("namespace")
	}

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

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

	// register handlers
	if car != nil {
		log.Logf("Registering Health handler at %s", HealthPath)
		r.Handle(HealthPath, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if c, err := car.hc(); err != nil {
				http.Error(w, err.Error(), c)
				return
			}
		}))
	}

	log.Logf("Registering Registry handler at %s", RegistryPath)
	r.Handle(RegistryPath, http.HandlerFunc(handler.Registry))

	log.Logf("Registering RPC handler at %s", RPCPath)
	r.Handle(RPCPath, http.HandlerFunc(handler.RPC))

	log.Logf("Registering Broker handler at %s", BrokerPath)
	r.Handle(BrokerPath, http.HandlerFunc(handler.Broker))

	switch Handler {
	case "proxy":
		log.Logf("Registering Proxy Handler at %s", ProxyPath)
		rt := router.NewRouter(router.WithNamespace(Namespace), router.WithHandler(api.Proxy))
		r.PathPrefix(ProxyPath).Handler(handler.Proxy(rt, nil, false))
	// rpc
	default:
		log.Logf("Registering Root Handler at %s", RootPath)
		rt := router.NewRouter(router.WithNamespace(Namespace), router.WithHandler(api.Rpc))
		r.PathPrefix(RootPath).Handler(handler.RPCX(rt, nil))
	}

	// reverse wrap handler
	plugins := append(Plugins(), plugin.Plugins()...)
	for i := len(plugins); i > 0; i-- {
		h = plugins[i-1].Handler()(h)
	}

	srv.Handle("/", h)

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

func (s *sidecar) hc() (int, error) {
	if len(s.hcUrl) == 0 {
		return 200, nil
	}
	rsp, err := http.Get(s.hcUrl)
	if err != nil {
		return 500, err
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != 200 {
		return rsp.StatusCode, fmt.Errorf("Non 200 response: %d", rsp.StatusCode)
	}
	return 200, nil
}

func (s *sidecar) hcLoop(service *registry.Service, exitCh chan bool) {
	tick := time.NewTicker(time.Second * 30)
	registered := true

	for {
		select {
		case <-tick.C:
			_, err := s.hc()
			if err != nil && registered {
				log.Logf("Healthcheck error. Deregistering %v", service.Nodes[0].Id)
				(*cmd.DefaultOptions().Registry).Deregister(service)
				registered = false
			} else if err == nil && !registered {
				log.Logf("Healthcheck success. Registering %v", service.Nodes[0].Id)
				(*cmd.DefaultOptions().Registry).Register(service)
				registered = true
			}
		case <-exitCh:
			return
		}
	}
}

// run healthchecker
func (s *sidecar) run(exit chan bool) {
	parts := strings.Split(s.address, ":")
	host := strings.Join(parts[:len(parts)-1], ":")
	port, _ := strconv.Atoi(parts[len(parts)-1])

	id := s.name + "-" + uuid.NewUUID().String()
	node := &registry.Node{
		Id:      id,
		Address: host,
		Port:    port,
	}

	service := &registry.Service{
		Name:  s.name,
		Nodes: []*registry.Node{node},
	}

	log.Logf("Registering %s", node.Id)
	(*cmd.DefaultOptions().Registry).Register(service)

	if len(s.hcUrl) == 0 {
		return
	}

	log.Log("Starting sidecar healthchecker")
	go s.hcLoop(service, exit)
	<-exit
}

func Commands() []cli.Command {
	command := cli.Command{
		Name:  "sidecar",
		Usage: "Run the micro sidecar",
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
			cli.StringFlag{
				Name:   "handler",
				Usage:  "Specify the request handler to be used for mapping HTTP requests to services; {proxy, rpc}",
				EnvVar: "MICRO_SIDECAR_HANDLER",
			},
			cli.StringFlag{
				Name:   "namespace",
				Usage:  "Set the namespace used by the sidecar e.g. com.example.srv",
				EnvVar: "MICRO_SIDECAR_NAMESPACE",
			},
			cli.StringFlag{
				Name:  "server_name",
				Usage: "Server name of the app",
			},
			cli.StringFlag{
				Name:  "server_address",
				Usage: "Server address and port of the app",
			},
			cli.StringFlag{
				Name:  "healthcheck_url",
				Usage: "URL to check health of the app",
			},
		},
		Action: func(c *cli.Context) {
			name := c.String("server_name")
			address := c.String("server_address")
			hcUrl := c.String("healthcheck_url")

			if len(name) == 0 && len(address) == 0 {
				run(c, nil)
				return
			}

			if len(name) == 0 {
				fmt.Println("Require server name")
				return
			}

			if len(address) == 0 {
				fmt.Println("Require server address")
				return
			}

			// exit chan
			exit := make(chan bool)

			// start the healthchecker
			car := newSidecar(name, address, hcUrl)
			go car.run(exit)

			// run the server
			run(c, car)

			// kill healthchecker
			close(exit)
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
