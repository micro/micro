package v1

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/registry"
	"github.com/micro/micro/web/common"
	"net/http"
	"sort"
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

// ServiceAPIDetail is the service api detail
type ServiceAPIDetail struct {
	Name      string               `json:"name,omitempty"`
	Endpoints []*registry.Endpoint `json:"endpoints,omitempty"`
}

func (api *API) InitV1Handler(r *mux.Router) {
	r.HandleFunc("/v1/services", api.services)
	r.HandleFunc("/v1/service/{name}", api.service)
	r.HandleFunc("/v1/service-details", api.serviceDetails)
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

	serviceDetails := make([]*ServiceAPIDetail, 0)
	for _, service := range services {
		s, err := (*cmd.DefaultOptions().Registry).GetService(service.Name)
		if err != nil {
			continue
		}
		if len(s) == 0 {
			continue
		}

		serviceDetails = append(serviceDetails, &ServiceAPIDetail{
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
			writeError(w, "Service Are Not found")
			return
		}

		writeJsonData(w, s)
		return
	}

	return
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
