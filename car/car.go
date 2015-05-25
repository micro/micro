package car

import (
	"encoding/json"
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
	"github.com/myodc/go-micro/registry"
)

var (
	Address = ":8081"
)

type Sidecar struct {
	name    string
	address string
	hcUrl   string
}

type Service struct {
	Name  string  `json:"name"`
	Nodes []*Node `json:"nodes"`
}

type Node struct {
	Id      string `json:"id"`
	Address string `json:"address"`
	Port    int    `json:"port"`
}

func run(name, address, hcUrl string) {
	sc := New(name, address, hcUrl)
	if err := sc.Run(); err != nil {
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

func (s *Sidecar) hcLoop(service registry.Service, exitCh chan bool) {
	tick := time.NewTicker(time.Second * 30)
	registered := true

	for {
		select {
		case <-tick.C:
			_, err := s.hc()
			if err != nil && registered {
				log.Infof("Healthcheck error. Deregistering %v", service.Nodes()[0].Id())
				registry.Deregister(service)
				registered = false
			} else if err == nil && !registered {
				log.Infof("Healthcheck success. Registering %v", service.Nodes()[0].Id())
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
		r.ParseForm()
		service := r.Form.Get("service")
		if len(service) == 0 {
			http.Error(w, "Require service", 400)
			return
		}
		s, err := registry.GetService(service)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if s == nil {
			http.Error(w, "Service not found", 404)
			return
		}
		srv := &Service{
			Name: s.Name(),
		}
		for _, node := range s.Nodes() {
			srv.Nodes = append(srv.Nodes, &Node{
				Id:      node.Id(),
				Address: node.Address(),
				Port:    node.Port(),
			})
		}
		b, err := json.Marshal(srv)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(b)))
		w.Write(b)
	})

	http.HandleFunc("/rpc", rpcHandler)

	http.HandleFunc("/broker", brokerHandler)

	if err := http.ListenAndServe(Address, nil); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func (s *Sidecar) Run() error {
	parts := strings.Split(s.address, ":")
	host := strings.Join(parts[:len(parts)-1], ":")
	port, _ := strconv.Atoi(parts[len(parts)-1])

	id := s.name + "-" + uuid.NewUUID().String()
	node := registry.NewNode(id, host, port)
	service := registry.NewService(s.name, node)

	log.Infof("Registering %s", node.Id())
	registry.Register(service)

	if len(s.hcUrl) > 0 {
		log.Info("Starting sidecar healthchecker")
		exitCh := make(chan bool, 1)
		go s.hcLoop(service, exitCh)
		defer func() {
			exitCh <- true
		}()
	}

	go s.serve()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	log.Infof("Received signal %s", <-ch)

	log.Infof("Deregistering %s", node.Id())
	registry.Deregister(service)
	return nil
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

				if len(name) == 0 {
					fmt.Println("Require server name")
					return
				}

				if len(address) == 0 {
					fmt.Println("Require server address")
					return
				}

				run(name, address, hcUrl)
			},
		},
	}
}
