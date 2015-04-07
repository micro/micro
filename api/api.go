package api

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/asim/go-micro/client"
	"github.com/asim/go-micro/errors"
	"github.com/asim/go-micro/server"
	"github.com/codegangsta/cli"
	log "github.com/golang/glog"
)

type ApiServer struct {
	mtx     sync.RWMutex
	address string
	exit    chan chan error
}

var (
	Address = ":8080"
	RpcPath = "/rpc"
)

func run() {
	server.Name = "go.micro.api"

	// Initialise Server
	server.Init()

	// Init API
	api := New(Address)
	api.Init()

	// Start API
	if err := api.Start(); err != nil {
		log.Fatal(err)
	}

	// Run server
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

	// Stop API
	if err := api.Stop(); err != nil {
		log.Fatal(err)
	}
}

func (s *ApiServer) Address() string {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.address
}

func (s *ApiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	var service, method string
	var request interface{}

	// response content type
	w.Header().Set("Content-Type", "application/json")

	switch r.Header.Get("Content-Type") {
	case "application/json":
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			e := errors.BadRequest("io.myodc.platform.server", err.Error())
			w.WriteHeader(400)
			w.Write([]byte(e.Error()))
			return
		}

		var body map[string]interface{}
		err = json.Unmarshal(b, &body)
		if err != nil {
			e := errors.BadRequest("io.myodc.platform.server", err.Error())
			w.WriteHeader(400)
			w.Write([]byte(e.Error()))
			return
		}

		service = body["service"].(string)
		method = body["method"].(string)
		request = body["request"]
	default:
		r.ParseForm()
		service = r.Form.Get("service")
		method = r.Form.Get("method")
		json.Unmarshal([]byte(r.Form.Get("request")), &request)
	}

	log.Infof("API Request: /rpc service: %s, method: %s", service, method)
	var response map[string]interface{}
	req := client.NewJsonRequest(service, method, request)
	err := client.Call(req, &response)
	if err != nil {
		log.Errorf("Error calling %s.%s: %v", service, method, err)
		ce := errors.Parse(err.Error())
		switch ce.Code {
		case 0:
			w.WriteHeader(500)
		default:
			w.WriteHeader(int(ce.Code))
		}
		w.Write([]byte(ce.Error()))
		return
	}

	b, _ := json.Marshal(response)
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.Write(b)
}

func (s *ApiServer) Init() error {
	log.Infof("API Rpc handler %s", RpcPath)
	http.Handle(RpcPath, s)
	return nil
}

func (s *ApiServer) Start() error {
	l, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	log.Infof("Listening on %s", l.Addr().String())

	s.mtx.Lock()
	s.address = l.Addr().String()
	s.mtx.Unlock()

	go http.Serve(l, nil)

	go func() {
		ch := <-s.exit
		ch <- l.Close()
	}()

	return nil
}

func (s *ApiServer) Stop() error {
	ch := make(chan error)
	s.exit <- ch
	return <-ch
}

func New(address string) *ApiServer {
	return &ApiServer{
		address: address,
		exit:    make(chan chan error),
	}
}

func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:  "api",
			Usage: "Run the micro API",
			Action: func(c *cli.Context) {
				run()
			},
		},
	}
}
