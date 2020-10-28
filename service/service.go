package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime"

	"github.com/micro/micro/v3/cmd"
	signalutil "github.com/micro/micro/v3/internal/signal"
	"github.com/micro/micro/v3/service/client"
	mudebug "github.com/micro/micro/v3/service/debug"
	debug "github.com/micro/micro/v3/service/debug/handler"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/server"
	"github.com/urfave/cli/v2"
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

// Run the default service and waits for it to exist
func Run() {
	// setup a new service, calling New() will trigger the cmd package
	// to parse the command line and
	srv := New()

	if err := srv.Run(); err == errMissingName {
		fmt.Println("Micro services must be run using \"micro run\"")
		os.Exit(1)
	} else if err != nil {
		logger.Fatalf("Error running %v service: %v", srv.Name(), err)
	}
}

// New returns a new Micro Service
func New(opts ...Option) *Service {
	// before extracts service options from the CLI flags. These
	// aren't set by the cmd package to prevent a circular dependancy.
	// prepend them to the array so options passed by the user to this
	// function are applied after (taking precedence)
	before := func(ctx *cli.Context) error {
		if n := ctx.String("service_name"); len(n) > 0 {
			opts = append([]Option{Name(n)}, opts...)
		}
		if v := ctx.String("service_version"); len(v) > 0 {
			opts = append([]Option{Version(v)}, opts...)
		}

		// service address injected by the runtime takes priority as the service port must match the
		// port the server is running on
		if a := ctx.String("service_address"); len(a) > 0 {
			opts = append(opts, Address(a))
		}
		return nil
	}

	// setup micro, this triggers the Before
	// function which parses CLI flags.
	cmd.New(cmd.SetupOnly(), cmd.Before(before)).Run()

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

// Handler registers a handler
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

func (s *Service) String() string {
	return "micro"
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

	// register the debug handler
	s.Server().Handle(
		s.Server().NewHandler(
			debug.NewHandler(s.Client()),
			server.InternalHandler(true),
		),
	)

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
