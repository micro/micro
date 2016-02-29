package server

import (
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"sync"

	log "github.com/golang/glog"
	"github.com/gorilla/handlers"
	"github.com/micro/micro/internal/stats"
)

type Server interface {
	Address() string
	Init(opts ...Option) error
	Handle(path string, handler http.Handler)
	Start() error
	Stop() error
}

type server struct {
	mux  *http.ServeMux
	opts Options

	mtx     sync.RWMutex
	address string
	exit    chan chan error
}

func NewServer(address string) Server {
	return &server{
		opts:    Options{},
		mux:     http.NewServeMux(),
		address: address,
		exit:    make(chan chan error),
	}
}

func (s *server) Address() string {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.address
}

func (s *server) Init(opts ...Option) error {
	for _, o := range opts {
		o(&s.opts)
	}
	return nil
}

func (s *server) Handle(path string, handler http.Handler) {
	s.mux.Handle(path, handlers.CombinedLoggingHandler(os.Stderr, handler))
}

func (s *server) Start() error {
	var l net.Listener
	var err error

	if s.opts.EnableTLS && s.opts.TLSConfig != nil {
		l, err = tls.Listen("tcp", s.address, s.opts.TLSConfig)
	} else {
		l, err = net.Listen("tcp", s.address)
	}
	if err != nil {
		return err
	}

	var h http.Handler
	exit := make(chan bool)

	// enable stats
	if len(s.opts.Stats) > 0 {
		st := stats.New(s.opts.Stats, s.mux)
		h = st
		go func() {
			st.Start()
			<-exit
			st.Stop()
		}()
	} else {
		h = s.mux
	}

	log.Infof("Listening on %s", l.Addr().String())

	s.mtx.Lock()
	s.address = l.Addr().String()
	s.mtx.Unlock()

	go http.Serve(l, h)

	go func() {
		close(exit)
		ch := <-s.exit
		ch <- l.Close()
	}()

	return nil
}

func (s *server) Stop() error {
	ch := make(chan error)
	s.exit <- ch
	return <-ch
}
