package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
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

	var opts []registry.RegisterOption

	// parse ttl
	if ttl := r.Form.Get("ttl"); len(ttl) > 0 {
		d, err := time.ParseDuration(ttl)
		if err == nil {
			opts = append(opts, registry.RegisterTTL(d))
		}
	}

	var service *registry.Service
	err = json.Unmarshal(b, &service)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = (*cmd.DefaultOptions().Registry).Register(service, opts...)
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
		//
		upgrade := r.Header.Get("Upgrade")
		connect := r.Header.Get("Connection")

		// watch if websockets
		if upgrade == "websocket" && connect == "Upgrade" {
			rw, err := (*cmd.DefaultOptions().Registry).Watch()
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			watch(rw, w, r)
			return
		}

		// otherwise list services
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

func ping(ws *websocket.Conn, exit chan bool) {
	ticker := time.NewTicker(pingTime)

	for {
		select {
		case <-ticker.C:
			ws.SetWriteDeadline(time.Now().Add(writeDeadline))
			err := ws.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				return
			}
		case <-exit:
			return
		}
	}
}

func watch(rw registry.Watcher, w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we need an exit chan
	exit := make(chan bool)

	defer func() {
		close(exit)
	}()

	// ping the socket
	go ping(ws, exit)

	for {
		// get next result
		r, err := rw.Next()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// write to client
		ws.SetWriteDeadline(time.Now().Add(writeDeadline))
		if err := ws.WriteJSON(r); err != nil {
			return
		}
	}
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
