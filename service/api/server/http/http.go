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
// Original source: github.com/micro/go-micro/v3/api/server/http/http.go

// Package http provides a http server with features; acme, cors, etc
package http

import (
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/handlers"
	"github.com/micro/micro/v5/service/api"
	"github.com/micro/micro/v5/service/logger"
)

type httpServer struct {
	mux  *http.ServeMux
	opts api.Options

	mtx     sync.RWMutex
	address string
	exit    chan chan error
}

func NewServer(address string, opts ...api.Option) api.Server {
	var options api.Options
	for _, o := range opts {
		o(&options)
	}

	return &httpServer{
		opts:    options,
		mux:     http.NewServeMux(),
		address: address,
		exit:    make(chan chan error),
	}
}

func (s *httpServer) Address() string {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.address
}

func (s *httpServer) Init(opts ...api.Option) error {
	for _, o := range opts {
		o(&s.opts)
	}
	return nil
}

func (s *httpServer) Handle(path string, handler http.Handler) {
	// TODO: move this stuff out to one place with ServeHTTP

	// apply the wrappers, e.g. auth
	for _, wrapper := range s.opts.Wrappers {
		handler = wrapper(handler)
	}

	// wrap with cors
	if s.opts.EnableCORS {
		handler = CombinedCORSHandler(handler)
	}

	// wrap with logger
	handler = handlers.CombinedLoggingHandler(os.Stdout, handler)

	s.mux.Handle(path, handler)
}

func (s *httpServer) Start() error {
	l, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	if logger.V(logger.InfoLevel, logger.DefaultLogger) {
		logger.Infof("HTTP Server Listening on %s", l.Addr().String())
	}

	s.mtx.Lock()
	s.address = l.Addr().String()
	s.mtx.Unlock()

	go func() {
		if err := http.Serve(l, s.mux); err != nil {
			// temporary fix
			//logger.Fatal(err)
		}
	}()

	go func() {
		ch := <-s.exit
		ch <- l.Close()
	}()

	return nil
}

func (s *httpServer) Stop() error {
	ch := make(chan error)
	s.exit <- ch
	return <-ch
}

func (s *httpServer) String() string {
	return "http"
}
