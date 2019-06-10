package server

import (
	"os"
	"sync"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/router"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-micro/transport"
	"github.com/micro/go-micro/transport/grpc"
	"github.com/micro/go-micro/util/log"
)

var (
	// Name of the server
	Name = "go.micro.server"
	// Address to bind route microservices to
	Address = ":8083"
	// Router address to bind to for router gossip
	Router = ":9094"
	// Network address to bind to
	Network = ":9093"
)

type srv struct {
	exit    chan struct{}
	service micro.Service
	router  router.Router
	network server.Server
	wg      *sync.WaitGroup
}

func newServer(s micro.Service, r router.Router) *srv {
	// NOTE: this will end up being QUIC transport
	t := grpc.NewTransport(transport.Addrs(Network))
	n := server.NewServer(server.Transport(t))

	return &srv{
		exit:    make(chan struct{}),
		service: s,
		router:  r,
		network: n,
		wg:      &sync.WaitGroup{},
	}
}

func (s *srv) start() error {
	log.Log("[server] starting micro server")

	// list all local services
	services, err := s.service.Client().Options().Registry.ListServices()
	if err != nil {
		log.Logf("[server] failed to list local services: %v", err)
		return err
	}

	// add services to routing table
	for _, service := range services {
		log.Logf("[server] adding route for local service %v", service)
		// create new micro network route
		r := router.NewRoute(
			router.DestAddr(service.Name),
			router.Hop(s.router),
			router.Network("local"),
			router.Metric(1),
		)
		// add new route to routing table
		if err := s.router.Table().Add(r); err != nil {
			log.Logf("[server] failed to add route to service: %v", service)
		}
	}

	log.Logf("[server] router has started: \n%s", s.router)
	log.Logf("[server] initial routing table: \n%s", s.router.Table())

	s.wg.Add(1)
	go s.watch()

	return nil
}

func (s *srv) watch() {
	log.Logf("[server] starting local registry watcher")

	defer s.wg.Done()
	w, err := s.service.Client().Options().Registry.Watch()
	if err != nil {
		log.Logf("[server] failed to create registry watch: %v", err)
		return
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		<-s.exit
		log.Logf("[server] stopping local registry watcher")
		w.Stop()
	}()

	// watch for changes to services
	for {
		res, err := w.Next()
		if err == registry.ErrWatcherStopped {
			log.Logf("[server] registry watcher stopped")
			return
		}

		if err != nil {
			log.Logf("[server] error watching registry: %s", err)
			return
		}

		switch res.Action {
		case "create":
			if len(res.Service.Nodes) > 0 {
				log.Logf("Action: %s, Service: %v", res.Action, res.Service.Name)
			}
		case "delete":
			log.Logf("Action: %s, Service: %v", res.Action, res.Service.Name)
		}
	}
}

func (s *srv) stop() error {
	log.Log("[server] attempting to stop")

	// notify all goroutines to finish
	close(s.exit)

	// wait for all goroutines to finish
	s.wg.Wait()

	return nil
}

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
	if len(ctx.String("network")) > 0 {
		Network = ctx.String("network")
	}
	if len(ctx.String("router")) > 0 {
		Router = ctx.String("router")
	}

	// Initialise service
	service := micro.NewService(
		micro.Name(Name),
		micro.Address(Address),
		micro.RegisterTTL(time.Duration(ctx.GlobalInt("register_ttl"))*time.Second),
		micro.RegisterInterval(time.Duration(ctx.GlobalInt("register_interval"))*time.Second),
	)

	// create new router
	r := router.NewRouter(
		router.ID(service.Server().Options().Id),
		router.Address(Address),
		router.GossipAddress(Router),
		router.NetworkAddr(Network),
	)

	// create new server and start it
	s := newServer(service, r)

	if err := s.start(); err != nil {
		log.Logf("error starting server: %v", err)
		os.Exit(1)
	}

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

	// stop the server
	if err := s.stop(); err != nil {
		log.Logf("error stopping server: %v", err)
		os.Exit(1)
	}

	log.Logf("[server] successfully stopped")
}

func Commands(options ...micro.Option) []cli.Command {
	command := cli.Command{
		Name:  "server",
		Usage: "Run the micro network server",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "address",
				Usage:  "Set the micro server address :8083",
				EnvVar: "MICRO_SERVER_ADDRESS",
			},
			cli.StringFlag{
				Name:   "network",
				Usage:  "Set the micro network address :9093",
				EnvVar: "MICRO_NETWORK_ADDRESS",
			},
			cli.StringFlag{
				Name:   "router",
				Usage:  "Set the micro router address :9094",
				EnvVar: "MICRO_ROUTER_ADDRESS",
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
