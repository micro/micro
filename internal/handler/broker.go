package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/cmd"
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
		(*cmd.DefaultOptions().Broker).Publish(c.topic, &broker.Message{Body: message})
	}
}

func (c *conn) write(mType int, data []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeDeadline))
	return c.ws.WriteMessage(mType, data)
}

func (c *conn) writeLoop() {
	ticker := time.NewTicker(pingTime)

	subscriber, err := (*cmd.DefaultOptions().Broker).Subscribe(c.topic, func(p broker.Publication) error {
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
		log.Print(err.Error())
		return
	}

	for _ = range ticker.C {
		if err := c.write(websocket.PingMessage, []byte{}); err != nil {
			return
		}
	}
}

func Broker(w http.ResponseWriter, r *http.Request) {
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
		log.Print(err.Error())
		return
	}

	once.Do(func() {
		(*cmd.DefaultOptions().Broker).Init()
		(*cmd.DefaultOptions().Broker).Connect()
	})

	c := &conn{
		topic: topic,
		ws:    ws,
	}

	go c.writeLoop()
	c.readLoop()
}
