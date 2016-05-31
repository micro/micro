package car

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/registry"
	"github.com/micro/micro/internal/handler"
	"github.com/micro/micro/internal/helper"
	"github.com/micro/micro/internal/server"
	"github.com/micro/micro/internal/stats"
	"github.com/pborman/uuid"
)

var (
	Address = ":8081"
)

type Sidecar struct {
	name    string
	address string
	hcUrl   string
}

var (
	BrokerPath   = "/broker"
	HealthPath   = "/health"
	RegistryPath = "/registry"
	RPCPath      = "/rpc"
	CORS         = map[string]bool{"*": true}
)

func run(ctx *cli.Context, car *Sidecar) {
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

	r := http.NewServeMux()

	// new server
	srv := server.NewServer(Address)
	srv.Init(opts...)

	// register handlers
	if car != nil {
		log.Printf("Registering Health handler at %s", HealthPath)
		r.Handle(HealthPath, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if c, err := car.hc(); err != nil {
				http.Error(w, err.Error(), c)
				return
			}
		}))
	}

	log.Printf("Registering Registry handler at %s", RegistryPath)
	r.Handle(RegistryPath, http.HandlerFunc(handler.Registry))

	log.Printf("Registering RPC handler at %s", RPCPath)
	r.Handle(RPCPath, http.HandlerFunc(handler.RPC))

	log.Printf("Registering Broker handler at %s", BrokerPath)
	r.Handle(BrokerPath, http.HandlerFunc(handler.Broker))

	var h http.Handler = r

	if ctx.GlobalBool("enable_stats") {
		st := stats.New()
		r.Handle("/stats", http.HandlerFunc(st.StatsHandler))
		h = st.ServeHTTP(r)
		st.Start()
		defer st.Stop()
	}

	// reverse wrap handler
	plugins := Plugins()
	for i := len(plugins); i > 0; i-- {
		h = plugins[i-1].Handle(h)
	}

	srv.Handle("/", h)

	// Initialise Server
	service := micro.NewService(
		micro.Name("go.micro.sidecar"),
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

func (s *Sidecar) hc() (int, error) {
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

func (s *Sidecar) hcLoop(service *registry.Service, exitCh chan bool) {
	tick := time.NewTicker(time.Second * 30)
	registered := true

	for {
		select {
		case <-tick.C:
			_, err := s.hc()
			if err != nil && registered {
				log.Printf("Healthcheck error. Deregistering %v", service.Nodes[0].Id)
				(*cmd.DefaultOptions().Registry).Deregister(service)
				registered = false
			} else if err == nil && !registered {
				log.Printf("Healthcheck success. Registering %v", service.Nodes[0].Id)
				(*cmd.DefaultOptions().Registry).Register(service)
				registered = true
			}
		case <-exitCh:
			return
		}
	}
}

// run healthchecker
func (s *Sidecar) run(exit chan bool) {
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

	log.Printf("Registering %s", node.Id)
	(*cmd.DefaultOptions().Registry).Register(service)

	if len(s.hcUrl) == 0 {
		return
	}

	log.Print("Starting sidecar healthchecker")
	go s.hcLoop(service, exit)
	<-exit
}

func New(name, address, hcUrl string) *Sidecar {
	return &Sidecar{
		name:    name,
		address: address,
		hcUrl:   hcUrl,
	}
}

func Commands() []cli.Command {
	command := cli.Command{
		Name:  "sidecar",
		Usage: "Run the micro sidecar",
		Flags: []cli.Flag{
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
			car := New(name, address, hcUrl)
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
