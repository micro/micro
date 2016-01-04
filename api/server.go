package api

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"sync"

	log "github.com/golang/glog"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/errors"

	"golang.org/x/net/context"
)

type ApiServer struct {
	mtx     sync.RWMutex
	address string
	exit    chan chan error
}

func newApiServer(address string) API {
	return &ApiServer{
		address: address,
		exit:    make(chan chan error),
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
			e := errors.BadRequest("go.micro.api", err.Error())
			w.WriteHeader(400)
			w.Write([]byte(e.Error()))
			return
		}

		var body map[string]interface{}
		err = json.Unmarshal(b, &body)
		if err != nil {
			e := errors.BadRequest("go.micro.api", err.Error())
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
	req := (*cmd.DefaultOptions().Client).NewJsonRequest(service, method, request)
	err := (*cmd.DefaultOptions().Client).Call(context.Background(), req, &response)
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
	log.Infof("API RPC handler %s", RpcPath)
	http.Handle(RpcPath, s)
	http.HandleFunc(HttpPath, restHandler)
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
