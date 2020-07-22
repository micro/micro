package service

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"sync"

	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/cmd"
	debug "github.com/micro/go-micro/v2/debug/service/handler"
	"github.com/micro/go-micro/v2/debug/trace"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/model"
	"github.com/micro/go-micro/v2/plugin"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-micro/v2/store"
	signalutil "github.com/micro/go-micro/v2/util/signal"

	// defaults
	grpcClient "github.com/micro/go-micro/v2/client/grpc"
	"github.com/micro/go-micro/v2/debug/stats"
	memTracer "github.com/micro/go-micro/v2/debug/trace/memory"
	"github.com/micro/go-micro/v2/model/mud"
	grpcServer "github.com/micro/go-micro/v2/server/grpc"
)

var (
	// DefaultClient for the service
	DefaultClient client.Client = grpcClient.NewClient()
	// DefaultServer for the service
	DefaultServer server.Server = grpcServer.NewServer()
	// DefaultModel for the service
	DefaultModel model.Model = mud.NewModel()

	// debugging interfaces
	defaultTracer trace.Tracer = memTracer.NewTracer()
	defaultStats  stats.Stats  = stats.NewStats()
)

func init() {
	// wrap the client
	DefaultClient = cacheClient(DefaultClient)
	DefaultClient = authClient(DefaultClient)
	DefaultClient = fromService(DefaultClient)
	DefaultClient = traceCall(DefaultClient)

	// wrap the server
	DefaultServer.Init(
		server.WrapHandler(handlerStats()),
		server.WrapHandler(traceHandler()),
		server.WrapHandler(authHandler()),
	)
}

// Service is a Micro Service which honours the go-micro/service interface
type Service struct {
	opts Options
	once sync.Once
}

// New returns a new Micro Service
func New(opts ...Option) *Service {
	return &Service{opts: newOptions(opts...)}
}

// Name of the service
func (s *Service) Name() string {
	return DefaultServer.Options().Name
}

// Init initialises options. Additionally it calls cmd.Init
// which parses command line flags. cmd.Init is only called
// on first Init.
func (s *Service) Init(opts ...Option) {
	// process options
	for _, o := range opts {
		o(&s.opts)
	}

	s.once.Do(func() {
		// setup the plugins
		for _, p := range strings.Split(os.Getenv("MICRO_PLUGIN"), ",") {
			if len(p) == 0 {
				continue
			}

			// load the plugin
			c, err := plugin.Load(p)
			if err != nil {
				logger.Fatal(err)
			}

			// initialise the plugin
			if err := plugin.Init(c); err != nil {
				logger.Fatal(err)
			}
		}

		// set cmd name
		if len(s.opts.Cmd.App().Name) == 0 {
			s.opts.Cmd.App().Name = s.Server().Options().Name
		}

		// Initialise the command options
		s.opts.Cmd.Init(
			cmd.Client(&DefaultClient),
			cmd.Server(&DefaultServer),
		)

		// run the command line
		// TODO: move to service.Run
		if err := s.opts.Cmd.Run(); err != nil {
			logger.Fatal(err)
		}

		// Explicitly set the table name to the service name
		name := s.Server().Options().Name
		store.DefaultStore.Init(store.Table(name))
	})
}

func (s *Service) Options() Options {
	return s.opts
}

func (s *Service) Client() client.Client {
	return DefaultClient
}

func (s *Service) Server() server.Server {
	return DefaultServer
}

func (s *Service) Model() model.Model {
	return DefaultModel
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

	if err := DefaultServer.Start(); err != nil {
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

	if err := DefaultServer.Stop(); err != nil {
		return err
	}

	for _, fn := range s.opts.AfterStop {
		if err := fn(); err != nil {
			gerr = err
		}
	}

	return gerr
}

func (s *Service) Run() error {
	// register the debug handler
	DefaultServer.Handle(
		DefaultServer.NewHandler(
			debug.NewHandler(DefaultClient),
			server.InternalHandler(true),
		),
	)

	// start the profiler
	// if s.opts.Profile != nil {
	// 	// to view mutex contention
	// 	rtime.SetMutexProfileFraction(5)
	// 	// to view blocking profile
	// 	rtime.SetBlockProfileRate(1)

	// 	if err := s.opts.Profile.Start(); err != nil {
	// 		return err
	// 	}
	// 	defer s.opts.Profile.Stop()
	// }

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
	return DefaultServer.Handle(DefaultServer.NewHandler(h, opts...))
}

// RegisterSubscriber is syntactic sugar for registering a subscriber
func RegisterSubscriber(topic string, h interface{}, opts ...server.SubscriberOption) error {
	return DefaultServer.Subscribe(DefaultServer.NewSubscriber(topic, h, opts...))
}

// Event is an object messages are published to
type Event struct {
	topic string
}

// Publish a message to an event
func (e *Event) Publish(ctx context.Context, msg interface{}, opts ...client.PublishOption) error {
	return DefaultClient.Publish(ctx, DefaultClient.NewMessage(e.topic, msg), opts...)
}

// NewEvent creates a new event publisher
func NewEvent(topic string) *Event {
	return &Event{topic}
}
