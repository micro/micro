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
	r.HandleFunc("/api/v1/service/{name}", api.service).Methods("GET")
	r.HandleFunc("/api/v1/service-details", api.serviceDetails).Methods("GET")
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

func rpc(w http.ResponseWriter, ctx context.Context, rpcReq *rpcRequest) {

	if len(rpcReq.Service) == 0 {
		writeError(w, "Service Is Not found")
	}

	if len(rpcReq.Endpoint) == 0 {
		writeError(w, "Endpoint Is Not found")
	}

	// decode rpc request param body
	if req, ok := rpcReq.Request.(string); ok {
		d := json.NewDecoder(strings.NewReader(req))
		d.UseNumber()

		if err := d.Decode(&rpcReq.Request); err != nil {
			writeError(w, "error decoding request string: "+err.Error())
			return
		}
	}

	// create request/response
	var response json.RawMessage
	var err error
	req := (*cmd.DefaultOptions().Client).NewRequest(rpcReq.Service, rpcReq.Endpoint, rpcReq.Request, client.WithContentType("application/json"))

	var opts []client.CallOption

	// set timeout
	if rpcReq.timeout > 0 {
		opts = append(opts, client.WithRequestTimeout(time.Duration(rpcReq.timeout)*time.Second))
	}

	// remote call
	if len(rpcReq.Address) > 0 {
		opts = append(opts, client.WithAddress(rpcReq.Address))
	}

	// remote call
	err = (*cmd.DefaultOptions().Client).Call(ctx, req, &response, opts...)
	if err != nil {
		ce := errors.Parse(err.Error())
		switch ce.Code {
		case 0:
			// assuming it's totally screwed
			ce.Code = 500
			ce.Id = "go.micro.rpc"
			ce.Status = http.StatusText(500)
			ce.Detail = "error during request: " + ce.Detail
			w.WriteHeader(500)
		default:
			w.WriteHeader(int(ce.Code))
		}
		w.Write([]byte(ce.Error()))
		return
	}

	writeJsonData(w, response)
}

func writeJsonData(w http.ResponseWriter, data interface{}) {

	rsp := &Rsp{
		Data:    data,
		Success: true,
	}

	b, err := json.Marshal(rsp)
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func writeError(w http.ResponseWriter, msg string) {

	rsp := &Rsp{
		Error:   msg,
		Success: false,
	}

	b, err := json.Marshal(rsp)
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
