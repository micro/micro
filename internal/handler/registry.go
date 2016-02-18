package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/registry"
)

func addService(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer r.Body.Close()

	var service *registry.Service
	err = json.Unmarshal(b, &service)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = (*cmd.DefaultOptions().Registry).Register(service)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func delService(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer r.Body.Close()

	var service *registry.Service
	err = json.Unmarshal(b, &service)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = (*cmd.DefaultOptions().Registry).Deregister(service)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func getService(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	service := r.Form.Get("service")

	var s []*registry.Service
	var err error

	if len(service) == 0 {
		s, err = (*cmd.DefaultOptions().Registry).ListServices()
	} else {
		s, err = (*cmd.DefaultOptions().Registry).GetService(service)
	}

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if s == nil || (len(service) > 0 && (len(s) == 0 || len(s[0].Name) == 0)) {
		http.Error(w, "Service not found", 404)
		return
	}

	b, err := json.Marshal(s)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.Write(b)
}

func Registry(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getService(w, r)
	case "POST":
		addService(w, r)
	case "DELETE":
		delService(w, r)
	}
}
