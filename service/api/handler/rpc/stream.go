// Copyright 2020 Asim Aslam
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Original source: github.com/micro/go-micro/v3/api/handler/rpc/stream.go

package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	pbapi "github.com/micro/micro/v3/proto/api"
	"github.com/micro/micro/v3/service/api"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	raw "github.com/micro/micro/v3/util/codec/bytes"
	"github.com/micro/micro/v3/util/router"
)

const (
	// Time allowed to write a message to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = 15 * time.Second

	// Maximum message size allowed from client.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func serveStream(ctx context.Context, w http.ResponseWriter, r *http.Request, service *api.Service, c client.Client) {
	// serve as websocket if thats the case
	if isWebSocket(r) {
		serveWebsocket(ctx, w, r, service, c)
		return
	}

	ct := r.Header.Get("Content-Type")
	// Strip charset from Content-Type (like `application/json; charset=UTF-8`)
	if idx := strings.IndexRune(ct, ';'); idx >= 0 {
		ct = ct[:idx]
	}

	payload, err := api.RequestPayload(r)
	if err != nil {
		if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
			logger.Error(err)
		}
		return
	}
	if len(payload) == 0 {
		// make it valid json
		payload = []byte("{}")
	}

	var request interface{}
	if !bytes.Equal(payload, []byte(`{}`)) {
		switch ct {
		case "application/json", "":
			m := json.RawMessage(payload)
			request = &m
		default:
			request = &raw.Frame{Data: payload}
		}
	}

	// we always need to set content type for message
	if ct == "" {
		ct = "application/json"
	}
	req := c.NewRequest(
		service.Name,
		service.Endpoint.Name,
		request,
		client.WithContentType(ct),
		client.StreamingRequest(),
	)

	w.Header().Set("Content-Type", ct)

	// create custom router
	callOpt := client.WithRouter(router.New(service.Services))

	// create a new stream
	stream, err := c.Stream(ctx, req, callOpt)
	if err != nil {
		if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
			logger.Error(err)
		}
		merr, ok := err.(*errors.Error)
		if ok {
			w.WriteHeader(int(merr.Code))
			w.Write([]byte(merr.Error()))
		}
		return
	}
	defer stream.Close()

	// send request even if nil because it triggers the call in case server expects no input
	// without this, we establish a connection but don't kick off the stream of communication
	if err = stream.Send(request); err != nil {
		if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
			logger.Error(err)
		}
		merr, ok := err.(*errors.Error)
		if ok {
			w.WriteHeader(int(merr.Code))
			w.Write([]byte(merr.Error()))
		} else {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		}
		return
	}

	rsp := stream.Response()

	// receive from stream and send to client
	for {
		select {
		case <-ctx.Done():
			return
		case <-stream.Context().Done():
			return
		default:
			// read backend response body
			buf, err := rsp.Read()
			if err != nil {
				// clean exit
				if err == io.EOF {
					return
				}
				// wants to avoid import  grpc/status.Status
				if strings.Contains(err.Error(), "context canceled") {
					return
				}
				if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
					logger.Error(err)
				}
				merr, ok := err.(*errors.Error)
				if ok {
					w.WriteHeader(int(merr.Code))
					w.Write([]byte(merr.Error()))
				}
				return
			}
			var bufOut string
			var apiRsp pbapi.Response
			if err := json.Unmarshal(buf, &apiRsp); err == nil && apiRsp.StatusCode > 0 {
				// bit of a hack. If the response is actually an api response we want to set the headers and status code
				for _, v := range apiRsp.Header {
					for _, s := range v.Values {
						w.Header().Add(v.Key, s)
					}
				}
				w.WriteHeader(int(apiRsp.StatusCode))
				bufOut = apiRsp.Body
			} else {
				bufOut = string(buf)
			}

			// send the buffer
			_, err = fmt.Fprint(w, bufOut)
			if err != nil {
				if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
					logger.Error(err)
				}
			}

			// flush it
			flusher, ok := w.(http.Flusher)
			if ok {
				flusher.Flush()
			}
		}
	}
}

type stream struct {
	// message type requested (binary or text)
	messageType int
	// request context
	ctx context.Context
	// the websocket connection.
	conn *websocket.Conn
	// the downstream connection.
	stream client.Stream
}

func (s *stream) processWSReadsAndWrites() {
	defer func() {
		s.conn.Close()
	}()

	msgs := make(chan []byte)

	stopCtx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(3)
	go s.rspToBufLoop(cancel, &wg, stopCtx, msgs)
	go s.bufToClientLoop(cancel, &wg, stopCtx, msgs)
	go s.clientToServerLoop(cancel, &wg, stopCtx)
	wg.Wait()
}

func (s *stream) clientToServerLoop(cancel context.CancelFunc, wg *sync.WaitGroup, stopCtx context.Context) {
	defer func() {
		s.stream.Close()
		cancel()
		wg.Done()
	}()
	s.conn.SetReadLimit(maxMessageSize)
	s.conn.SetReadDeadline(time.Now().Add(pongWait))
	s.conn.SetPongHandler(func(string) error { s.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		select {
		case <-stopCtx.Done():
			return
		default:
		}

		_, msg, err := s.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
					logger.Error(err)
				}
			}
			return
		}

		var request interface{}
		switch s.messageType {
		case websocket.TextMessage:
			m := json.RawMessage(msg)
			request = &m
		default:
			request = &raw.Frame{Data: msg}
		}

		if err := s.stream.Send(request); err != nil {
			if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
				logger.Error(err)
			}
			return
		}
	}

}

func (s *stream) rspToBufLoop(cancel context.CancelFunc, wg *sync.WaitGroup, stopCtx context.Context, msgs chan []byte) {
	defer func() {
		cancel()
		wg.Done()
	}()
	rsp := s.stream.Response()
	for {
		select {
		case <-stopCtx.Done():
			return
		default:
		}
		bytes, err := rsp.Read()
		if err != nil {
			if err == io.EOF {
				// clean exit
				return
			}
			// write error then close the connection
			b, _ := json.Marshal(err)
			s.conn.WriteMessage(s.messageType, b)
			s.conn.WriteMessage(websocket.CloseAbnormalClosure, []byte{})
			return
		}
		select {
		case <-stopCtx.Done():
			return
		case msgs <- bytes:
		}

	}

}

func (s *stream) bufToClientLoop(cancel context.CancelFunc, wg *sync.WaitGroup, stopCtx context.Context, msgs chan []byte) {
	defer func() {
		s.conn.Close()
		cancel()
		wg.Done()

	}()
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-stopCtx.Done():
			return
		case <-s.ctx.Done():
			return
		case <-s.stream.Context().Done():
			s.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		case <-ticker.C:
			s.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := s.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case msg := <-msgs:
			// read response body
			s.conn.SetWriteDeadline(time.Now().Add(writeWait))
			w, err := s.conn.NextWriter(s.messageType)
			if err != nil {
				return
			}
			if _, err := w.Write(msg); err != nil {
				return
			}
			if err := w.Close(); err != nil {
				return
			}
		}
	}

}

// serveWebsocket will stream rpc back over websockets assuming json
func serveWebsocket(ctx context.Context, w http.ResponseWriter, r *http.Request, service *api.Service, c client.Client) {
	var rspHdr http.Header
	// we use Sec-Websocket-Protocol to pass auth headers so just accept anything here
	if prots := r.Header.Values("Sec-WebSocket-Protocol"); len(prots) > 0 {
		rspHdr = http.Header{}
		for _, p := range prots {
			rspHdr.Add("Sec-WebSocket-Protocol", p)
		}
	}

	conn, err := upgrader.Upgrade(w, r, rspHdr)
	if err != nil {
		if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
			logger.Error(err)
		}
		return
	}

	// determine the content type
	ct := r.Header.Get("Content-Type")
	// strip charset from Content-Type (like `application/json; charset=UTF-8`)
	if idx := strings.IndexRune(ct, ';'); idx >= 0 {
		ct = ct[:idx]
	}
	if len(ct) == 0 {
		ct = "application/json"
	}

	// create stream
	req := c.NewRequest(service.Name, service.Endpoint.Name, nil, client.WithContentType(ct), client.StreamingRequest())
	str, err := c.Stream(ctx, req, client.WithRouter(router.New(service.Services)))
	if err != nil {
		if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
			logger.Error(err)
		}
		return
	}

	// determine the message type
	msgType := websocket.BinaryMessage
	if ct == "application/json" {
		msgType = websocket.TextMessage
	}

	s := stream{ctx: ctx, conn: conn, stream: str, messageType: msgType}
	s.processWSReadsAndWrites()
}

func isStream(r *http.Request, srv *api.Service) bool {
	// check if the endpoint supports streaming
	for _, service := range srv.Services {
		for _, ep := range service.Endpoints {
			// skip if it doesn't match the name
			if ep.Name != srv.Endpoint.Name {
				continue
			}
			// matched if the name
			if v := ep.Metadata["stream"]; v == "true" {
				return true
			}
		}
	}

	return false
}

func isWebSocket(r *http.Request) bool {
	contains := func(key, val string) bool {
		vv := strings.Split(r.Header.Get(key), ",")
		for _, v := range vv {
			if val == strings.ToLower(strings.TrimSpace(v)) {
				return true
			}
		}
		return false
	}

	if contains("Connection", "upgrade") && contains("Upgrade", "websocket") {
		return true
	}

	return false
}
