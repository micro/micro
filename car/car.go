package car

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"github.com/codegangsta/cli"
	log "github.com/golang/glog"
	"github.com/micro/go-micro/registry"
)

var (
	Address = ":8081"
)

type Sidecar struct {
	name    string
	address string
	hcUrl   string
}

func run() {
	sc := New("", "", "")
	sc.serve()
}

func runApp(name, address, hcUrl string) {
	sc := New(name, address, hcUrl)
	go sc.serve()
	sc.run()
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
				registry.Deregister(service)
				registered = false
			} else if err == nil && !registered {
				log.Infof("Healthcheck success. Registering %v", service.Nodes[0].Id)
				registry.Register(service)
				registered = true
			}
		case <-exitCh:
			return
		}
	}
}

func (s *Sidecar) serve() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if c, err := s.hc(); err != nil {
			http.Error(w, err.Error(), c)
			return
		}
	})

	http.HandleFunc("/registry", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getService(w, r)
		case "POST":
			addService(w, r)
		case "DELETE":
			delService(w, r)
		}
	})

	http.HandleFunc("/rpc", rpcHandler)

	http.HandleFunc("/broker", brokerHandler)

	log.Infof("Listening on %s", Address)
	if err := http.ListenAndServe(Address, nil); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func (s *Sidecar) run() {
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
	registry.Register(service)

	if len(s.hcUrl) > 0 {
		log.Info("Starting sidecar healthchecker")
		exitCh := make(chan bool, 1)
		go s.hcLoop(service, exitCh)
		defer func() {
			exitCh <- true
		}()
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	<-ch
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
					run()
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

				runApp(name, address, hcUrl)
			},
		},
	}
}
