package car

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/golang/glog"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/registry"
	"github.com/micro/micro/internal/handler"
	"github.com/micro/micro/internal/server"
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
)

func run(ctx *cli.Context, car *Sidecar) {
	var opts []server.Option

	if ctx.GlobalBool("enable_tls") {
		cert := ctx.GlobalString("tls_cert_file")
		key := ctx.GlobalString("tls_key_file")

		if len(cert) > 0 && len(key) > 0 {
			certs, err := tls.LoadX509KeyPair(cert, key)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			config := &tls.Config{
				Certificates: []tls.Certificate{certs},
			}
			opts = append(opts, server.EnableTLS(true))
			opts = append(opts, server.TLSConfig(config))
		} else {
			fmt.Println("Enable TLS specified without certificate and key files")
			return
		}
	}

	// new server
	srv := server.NewServer(Address)
	srv.Init(opts...)

	// register handlers
	if car != nil {
		log.Infof("Registering Health handler at %s", HealthPath)
		srv.Handle(HealthPath, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if c, err := car.hc(); err != nil {
				http.Error(w, err.Error(), c)
				return
			}
		}))
	}

	log.Infof("Registering Registry handler at %s", RegistryPath)
	srv.Handle(RegistryPath, http.HandlerFunc(handler.Registry))

	log.Infof("Registering RPC handler at %s", RPCPath)
	srv.Handle(RPCPath, http.HandlerFunc(handler.RPC))

	log.Infof("Registering Broker handler at %s", BrokerPath)
	srv.Handle(BrokerPath, http.HandlerFunc(handler.Broker))

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
				log.Infof("Healthcheck error. Deregistering %v", service.Nodes[0].Id)
				(*cmd.DefaultOptions().Registry).Deregister(service)
				registered = false
			} else if err == nil && !registered {
				log.Infof("Healthcheck success. Registering %v", service.Nodes[0].Id)
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

	log.Infof("Registering %s", node.Id)
	(*cmd.DefaultOptions().Registry).Register(service)

	if len(s.hcUrl) == 0 {
		return
	}

	log.Info("Starting sidecar healthchecker")
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
	return []cli.Command{
		{
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
		},
	}
}
