package service

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"micro.dev/v4/cmd"
	"micro.dev/v4/service/client"
	"micro.dev/v4/service/logger"
	"micro.dev/v4/service/model"
	"micro.dev/v4/service/profile"
	"micro.dev/v4/service/server"
)

var (
	// errMissingName is returned by service.Run when a service is run
	// prior to it's name being set.
	errMissingName = errors.New("missing service name")
)

// Service is a Micro Service which honours the go-micro/service interface
type Service struct {
	opts Options
}

// New returns a new Micro Service
func New(opts ...Option) *Service {
	// setup micro, this triggers the Before
	// function which parses CLI flags.
	cmd.New(cmd.Service()).Run()

	// setup auth
	profile.SetupAccount(nil)

	// return a new service
	return &Service{opts: newOptions(opts...)}
}

// Name of the service
func (s *Service) Name() string {
	return s.opts.Name
}

// Version of the service
func (s *Service) Version() string {
	return s.opts.Version
}

// Handle registers a handler
func (s *Service) Handle(v interface{}) error {
	return s.Server().Handle(s.Server().NewHandler(v))
}

// Subscribe registers a subscriber
func (s *Service) Subscribe(topic string, v interface{}) error {
	return s.Server().Subscribe(s.Server().NewSubscriber(topic, v))
}

func (s *Service) Init(opts ...Option) {
	for _, o := range opts {
		o(&s.opts)
	}
}

func (s *Service) Options() Options {
	return s.opts
}

func (s *Service) Client() client.Client {
	return client.DefaultClient
}

func (s *Service) Server() server.Server {
	return server.DefaultServer
}

func (s *Service) Model() model.Model {
	return model.DefaultModel
}

func (s *Service) Start() error {
	for _, fn := range s.opts.BeforeStart {
		if err := fn(); err != nil {
			return err
		}
	}

	if err := s.Server().Start(); err != nil {
		return err
	}

	for _, fn := range s.opts.AfterStart {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) Stop() error {
	var gerr error

	for _, fn := range s.opts.BeforeStop {
		if err := fn(); err != nil {
			gerr = err
		}
	}

	if err := server.DefaultServer.Stop(); err != nil {
		return err
	}

	for _, fn := range s.opts.AfterStop {
		if err := fn(); err != nil {
			gerr = err
		}
	}

	return gerr
}

// Run the service
func (s *Service) Run() error {
	// ensure service's have a name, this is injected by the runtime manager
	if len(s.Name()) == 0 {
		return errMissingName
	}

	if logger.V(logger.InfoLevel, logger.DefaultLogger) {
		logger.Infof("Starting [service] %s", s.Name())
	}

	if err := s.Start(); err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	if s.opts.Signal {
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL)
	}

	// wait on kill signal
	<-ch
	return s.Stop()
}

// Handle is syntactic sugar for registering a handler
func Handle(h interface{}, opts ...server.HandlerOption) error {
	return server.DefaultServer.Handle(server.DefaultServer.NewHandler(h, opts...))
}

// Subscribe is syntactic sugar for registering a subscriber
func Subscribe(topic string, h interface{}, opts ...server.SubscriberOption) error {
	return server.DefaultServer.Subscribe(server.DefaultServer.NewSubscriber(topic, h, opts...))
}

// Event is an object messages are published to
type Event struct {
	topic string
}

// Publish a message to an event
func (e *Event) Publish(ctx context.Context, msg interface{}) error {
	return client.Publish(ctx, client.NewMessage(e.topic, msg))
}

// NewEvent creates a new event publisher
func NewEvent(topic string) *Event {
	return &Event{topic}
}
