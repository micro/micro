package server

import (
	"fmt"
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
		return fmt.Errorf("failed to list local services: %v", err)
	}

	// add services to routing table
	for _, service := range services {
		log.Logf("[server] adding route for local service %v", service)
		// create new micro network route
		r := router.NewRoute(
			router.DestAddr(service.Name),
			router.Gateway(s.router),
			router.Network("local"),
			router.Metric(1),
		)
		// add new route to routing table
		if err := s.router.Table().Add(r); err != nil {
			log.Logf("[server] failed to add route for service: %v", service)
		}
	}

	// start micor network router
	if err := s.router.Start(); err != nil {
		return fmt.Errorf("failed to start router: %v", err)
	}

	log.Logf("[server] router has started: \n%s", s.router)
	log.Logf("[server] initial routing table: \n%s", s.router.Table())

	// get local registry watcher
	w, err := s.service.Client().Options().Registry.Watch()
	if err != nil {
		return fmt.Errorf("failed to create local registry watch: %v", err)
	}

	s.wg.Add(1)
	go s.watch(w)

	return nil
}

// watch local registry
func (s *srv) watch(w registry.Watcher) {
	log.Logf("[server] starting local registry watcher")
	defer s.wg.Done()

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

		log.Logf("[server] watcher action: %s, Service: %v", res.Action, res.Service.Name)
		// create new route
		r := router.NewRoute(
			router.DestAddr(res.Service.Name),
			router.Gateway(s.router),
			router.Network("local"),
			router.Metric(1),
		)

		switch res.Action {
		case "create":
			if len(res.Service.Nodes) > 0 {
				log.Logf("[server] adding route for local service %v", res.Service.Name)
				// add new route to routing table
				if err := s.router.Table().Add(r); err != nil {
					log.Logf("[server] failed to add route for service: %v", res.Service.Name)
				}
			}
		case "delete":
			log.Logf("[server] removing route for local service %v", res.Service.Name)
			// delete route from routing table
			if err := s.router.Table().Remove(r); err != nil {
				log.Logf("[server] failed to remove route for service: %v", res.Service.Name)
			}
		}
	}
}

func (s *srv) stop() error {
	log.Log("[server] attempting to stop server")

	// notify all goroutines to finish
	close(s.exit)

	// wait for all goroutines to finish
	s.wg.Wait()

	log.Log("[server] attempting to stop router")

	// stop the router
	if err := s.router.Stop(); err != nil {
		return fmt.Errorf("failed to stop router: %v", err)
	}

	log.Log("[server] server successfully stopped")

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
		router.GossipAddr(Router),
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
