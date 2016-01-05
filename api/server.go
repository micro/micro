package api

import (
	"net"
	"net/http"
	"sync"

	log "github.com/golang/glog"
)

type ApiServer struct {
	mux *http.ServeMux

	mtx     sync.RWMutex
	address string
	exit    chan chan error
}

func newApiServer(address string) API {
	return &ApiServer{
		mux:     http.NewServeMux(),
		address: address,
		exit:    make(chan chan error),
	}
}

func (s *ApiServer) Address() string {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.address
}

func (s *ApiServer) Init() error {
	return nil
}

func (s *ApiServer) Handle(path string, handler http.Handler) {
	s.mux.Handle(path, handler)
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

	go http.Serve(l, s.mux)

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
