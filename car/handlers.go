package car

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	log "github.com/golang/glog"
	"github.com/gorilla/websocket"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/registry"

	"golang.org/x/net/context"
)

const (
	pingTime      = (readDeadline * 9) / 10
	readLimit     = 16384
	readDeadline  = 60 * time.Second
	writeDeadline = 10 * time.Second
)

type conn struct {
	topic string
	ws    *websocket.Conn
}

var (
	once sync.Once

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
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
	err = registry.Register(service)
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
	err = registry.Deregister(service)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func getService(w http.ResponseWriter, r *http.Request) {
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
	if s == nil || len(s) == 0 || len(s[0].Name) == 0 {
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

func rpcHandler(w http.ResponseWriter, r *http.Request) {
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

	var response map[string]interface{}
	req := client.NewJsonRequest(service, method, request)
	err := client.Call(context.Background(), req, &response)
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

func brokerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	r.ParseForm()
	topic := r.Form.Get("topic")
	if len(topic) == 0 {
		http.Error(w, "Topic not specified", 400)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err.Error())
		return
	}

	once.Do(func() {
		broker.Init()
		broker.Connect()
	})

	c := &conn{
		topic: topic,
		ws:    ws,
	}

	go c.writeLoop()
	c.readLoop()
}

func (c *conn) readLoop() {
	defer func() {
		c.ws.Close()
	}()

	c.ws.SetReadLimit(readLimit)
	c.ws.SetReadDeadline(time.Now().Add(readDeadline))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(readDeadline))
		return nil
	})

	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			return
		}
		broker.Publish(c.topic, &broker.Message{Body: message})
	}
}

func (c *conn) write(mType int, data []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeDeadline))
	return c.ws.WriteMessage(mType, data)
}

func (c *conn) writeLoop() {
	ticker := time.NewTicker(pingTime)

	subscriber, err := broker.Subscribe(c.topic, func(p broker.Publication) error {
		b, err := json.Marshal(p.Message())
		if err != nil {
			return nil
		}
		return c.write(websocket.TextMessage, b)
	})

	defer func() {
		subscriber.Unsubscribe()
		ticker.Stop()
		c.ws.Close()
	}()

	if err != nil {
		log.Error(err.Error())
		return
	}

	for _ = range ticker.C {
		if err := c.write(websocket.PingMessage, []byte{}); err != nil {
			return
		}
	}
}
