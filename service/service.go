package service

import (
	"context"
	"os"
	"os/signal"
	"runtime"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/client"
	debug "github.com/micro/go-micro/v2/debug/service/handler"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/model"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-micro/v2/store"
	signalutil "github.com/micro/go-micro/v2/util/signal"
	"github.com/micro/micro/v2/cmd"
	muclient "github.com/micro/micro/v2/service/client"
	mudebug "github.com/micro/micro/v2/service/debug"
	mumodel "github.com/micro/micro/v2/service/model"
	muserver "github.com/micro/micro/v2/service/server"
)

// Service is a Micro Service which honours the go-micro/service interface
type Service struct {
	opts Options
}

// New returns a new Micro Service
func New(opts ...Option) *Service {
	return &Service{opts: newOptions(opts...)}
}

// Name of the service
func (s *Service) Name() string {
	return muserver.DefaultServer.Options().Name
}

// Init initialises options. Additionally it calls cmd.Init
// which parses command line flags. cmd.Init is only called
// on first Init.
func (s *Service) Init(opts ...Option) {
	for _, o := range opts {
		o(&s.opts)
	}
}

func (s *Service) Options() Options {
	return s.opts
}

func (s *Service) Client() client.Client {
	return muclient.DefaultClient
}

func (s *Service) Server() server.Server {
	return muserver.DefaultServer
}

func (s *Service) Model() model.Model {
	return mumodel.DefaultModel
}

func (s *Service) String() string {
	return "micro"
}

func (s *Service) Start() error {
	// set the store to use the service name as the table
	store.DefaultStore.Init(store.Table(s.Name()))

	for _, fn := range s.opts.BeforeStart {
		if err := fn(); err != nil {
			return err
		}
	}

	if err := muserver.DefaultServer.Start(); err != nil {
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

	if err := muserver.DefaultServer.Stop(); err != nil {
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
	// run the app wrapped by the cmd package so it
	// initializes micro before running the service.
	return cmd.New(cmd.Action(s.run)).Run()
}

func (s *Service) run(ctx *cli.Context) error {
	// register the debug handler
	muserver.DefaultServer.Handle(
		muserver.DefaultServer.NewHandler(
			debug.NewHandler(muclient.DefaultClient),
			server.InternalHandler(true),
		),
	)

	// setup service auth credentials
	if err := setupAuth(); err != nil {
		return err
	}

	// start the profiler
	if mudebug.DefaultProfiler != nil {
		// to view mutex contention
		runtime.SetMutexProfileFraction(5)
		// to view blocking profile
		runtime.SetBlockProfileRate(1)

		if err := mudebug.DefaultProfiler.Start(); err != nil {
			return err
		}

		defer mudebug.DefaultProfiler.Stop()
	}

	if logger.V(logger.InfoLevel, logger.DefaultLogger) {
		logger.Infof("Starting [service] %s", s.Name())
	}

	if err := s.Start(); err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	if s.opts.Signal {
		signal.Notify(ch, signalutil.Shutdown()...)
	}

	// wait on kill signal
	<-ch
	return s.Stop()
}

// RegisterHandler is syntactic sugar for registering a handler
func RegisterHandler(h interface{}, opts ...server.HandlerOption) error {
	return muserver.DefaultServer.Handle(muserver.DefaultServer.NewHandler(h, opts...))
}

// RegisterSubscriber is syntactic sugar for registering a subscriber
func RegisterSubscriber(topic string, h interface{}, opts ...server.SubscriberOption) error {
	return muserver.DefaultServer.Subscribe(muserver.DefaultServer.NewSubscriber(topic, h, opts...))
}

// Event is an object messages are published to
type Event struct {
	topic string
}

// Publish a message to an event
func (e *Event) Publish(ctx context.Context, msg interface{}, opts ...client.PublishOption) error {
	return muclient.DefaultClient.Publish(ctx, muclient.DefaultClient.NewMessage(e.topic, msg), opts...)
}

// NewEvent creates a new event publisher
func NewEvent(topic string) *Event {
	return &Event{topic}
}
