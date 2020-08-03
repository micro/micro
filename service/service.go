package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime"

	"github.com/micro/go-micro/v3/client"
	debug "github.com/micro/go-micro/v3/debug/service/handler"
	"github.com/micro/go-micro/v3/server"
	"github.com/micro/go-micro/v3/store"
	signalutil "github.com/micro/go-micro/v3/util/signal"
	muclient "github.com/micro/micro/v3/service/client"
	mudebug "github.com/micro/micro/v3/service/debug"
	"github.com/micro/micro/v3/service/logger"
	muserver "github.com/micro/micro/v3/service/server"
)

var (
	// defaultService sets up a new service
	defaultService *Service
	// errMissingName is returned by service.Run when a service is run
	// prior to it's name being set.
	errMissingName = errors.New("missing service name")
)

// Service is a Micro Service which honours the go-micro/service interface
type Service struct {
	opts Options
}

// New returns a new service
func New(opts ...Option) *Service {
	s := &Service{opts: newOptions(opts...)}
	s.Options().Cmd.Run()
	return s
}

// Init the default service
func Init(opts ...Option) {
	if defaultService == nil {
		setupDefaultService()
		return
	}

	for _, o := range opts {
		o(&defaultService.opts)
	}
}

// Run the default service and waits for it to exist
func Run() {
	if defaultService == nil {
		setupDefaultService()
	}

	if err := defaultService.Run(); err == errMissingName {
		fmt.Println("Micro services must be run using \"micro run\"")
		os.Exit(1)
	} else if err != nil {
		logger.Fatalf("Error running %v service: %v", defaultService.Name(), err)
	}
}

// Handle registers a handler
func Handle(h interface{}, opts ...server.HandlerOption) error {
	return muserver.DefaultServer.Handle(muserver.DefaultServer.NewHandler(h, opts...))
}

// Subscribe to a topic
func Subscribe(topic string, h interface{}, opts ...server.SubscriberOption) error {
	return muserver.DefaultServer.Subscribe(muserver.DefaultServer.NewSubscriber(topic, h, opts...))
}

// Name of the service
func (s *Service) Name() string {
	return s.opts.Name
}

// Version of the service
func (s *Service) Version() string {
	return s.opts.Version
}

// Options for the service
func (s *Service) Options() Options {
	return s.opts
}

// Start the service
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

// Stop the service
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
	// ensure service's have a name, this is injected by the runtime manager
	if len(s.Name()) == 0 {
		return errMissingName
	}

	// register the debug handler
	muserver.DefaultServer.Handle(
		muserver.DefaultServer.NewHandler(
			debug.NewHandler(muclient.DefaultClient),
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

// Event is an object messages are published to
type Event struct {
	topic string
}

// Publish a message to an event
func (e *Event) Publish(ctx context.Context, msg interface{}, opts ...client.PublishOption) error {
	return muclient.Publish(ctx, muclient.NewMessage(e.topic, msg), opts...)
}

// NewEvent creates a new event publisher
func NewEvent(topic string) *Event {
	return &Event{topic}
}

// setupDefaultService sets up the defaultService variable. We don't do
// this in init because it will result in micro always being configured
// as if a service was being run
func setupDefaultService() {
	defaultService = &Service{opts: newOptions()}
	defaultService.Options().Cmd.Run()
}
