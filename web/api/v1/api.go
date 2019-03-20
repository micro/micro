package v1

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/registry"
	"github.com/micro/micro/internal/helper"
	"github.com/micro/micro/web/common"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	staticDir = "/webapp/dist/"
	Namespace = "go.micro.web"
)

// Rsp is the struct of api response
type Rsp struct {
	Code    uint        `json:"code,omitempty"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// API is handler of all API calls.
type API struct {
}

// serviceAPIDetail is the service api detail
type serviceAPIDetail struct {
	Name      string               `json:"name,omitempty"`
	Endpoints []*registry.Endpoint `json:"endpoints,omitempty"`
}

type rpcRequest struct {
	Service  string
	Endpoint string
	Method   string
	Address  string
	timeout  int
	Request  interface{}
}

func (api *API) InitV1Handler(r *mux.Router) {
	r.HandleFunc("/api/v1/services", api.services).Methods("GET")
	r.HandleFunc("/api/v1/service/{name:[a-zA-Z0-9]+}", api.service).Methods("GET")
	r.HandleFunc("/api/v1/service-details", api.serviceDetails).Methods("GET")
	r.HandleFunc("/api/v1/stats/{name:[a-zA-Z0-9]+}", api.stats).Methods("GET")
	r.HandleFunc("/api/v1/web-services", api.webServices).Methods("GET")
	r.HandleFunc("/api/v1/rpc", api.rpc).Methods("POST")
	r.HandleFunc("/api/v1/health", api.health).Methods("GET")
}

func (api *API) webServices(w http.ResponseWriter, r *http.Request) {
	services, err := (*cmd.DefaultOptions().Registry).ListServices()
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}

	webServices := make([]*registry.Service, 0)
	for _, s := range services {

		if strings.Index(s.Name, Namespace) == 0 && len(strings.TrimPrefix(s.Name, Namespace)) > 0 {
			s.Name = strings.Replace(s.Name, Namespace+".", "", 1)
			webServices = append(webServices, s)
		}
	}

	sort.Sort(common.SortedServices{Services: services})

	writeJsonData(w, webServices)

	return
}

func (api *API) services(w http.ResponseWriter, r *http.Request) {

	services, err := (*cmd.DefaultOptions().Registry).ListServices()
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}

	sort.Sort(common.SortedServices{Services: services})

	writeJsonData(w, services)
	return
}

func (api *API) serviceDetails(w http.ResponseWriter, r *http.Request) {
	services, err := (*cmd.DefaultOptions().Registry).ListServices()
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}

	sort.Sort(common.SortedServices{Services: services})

	serviceDetails := make([]*serviceAPIDetail, 0)
	for _, service := range services {
		s, err := (*cmd.DefaultOptions().Registry).GetService(service.Name)
		if err != nil {
			continue
		}
		if len(s) == 0 {
			continue
		}

		serviceDetails = append(serviceDetails, &serviceAPIDetail{
			Name:      service.Name,
			Endpoints: s[0].Endpoints,
		})
	}

	writeJsonData(w, serviceDetails)
	return
}

func (api *API) service(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if len(vars["name"]) > 0 {
		s, err := (*cmd.DefaultOptions().Registry).GetService(vars["name"])
		if err != nil {
			http.Error(w, "Error occurred:"+err.Error(), 500)
			return
		}

		if len(s) == 0 {
			writeError(w, "Service Is Not found")
			return
		}

		writeJsonData(w, s)
		return
	}

	return
}

func (api *API) rpc(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	rpcReq := &rpcRequest{}

	d := json.NewDecoder(r.Body)
	d.UseNumber()

	if err := d.Decode(&rpcReq); err != nil {
		writeError(w, err.Error())
		return
	}

	if len(rpcReq.Endpoint) == 0 {
		rpcReq.Endpoint = rpcReq.Method
	}

	rpcReq.timeout, _ = strconv.Atoi(r.Header.Get("Timeout"))

	rpc(w, helper.RequestToContext(r), rpcReq)
}

func (api *API) health(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	rpcReq := &rpcRequest{
		Service:  r.URL.Query().Get("service"),
		Endpoint: "Debug.Health",
		Request:  "{}",
		Address:  r.URL.Query().Get("address"),
	}

	rpc(w, helper.RequestToContext(r), rpcReq)
}

