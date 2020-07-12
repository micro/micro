// Package service provides a micro service
package service

import (
	"os"
	"os/signal"
	rtime "runtime"
	"strings"
	"sync"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/debug/service/handler"
	"github.com/micro/go-micro/v2/debug/stats"
	"github.com/micro/go-micro/v2/debug/trace"
	"github.com/micro/go-micro/v2/logger"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/plugin"
	"github.com/micro/go-micro/v2/proxy"
	"github.com/micro/go-micro/v2/proxy/grpc"
	"github.com/micro/go-micro/v2/proxy/http"
	"github.com/micro/go-micro/v2/proxy/mucp"
	"github.com/micro/go-micro/v2/runtime"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-micro/v2/store"
	signalutil "github.com/micro/go-micro/v2/util/signal"
	"github.com/micro/go-micro/v2/util/wrapper"
	"github.com/micro/micro/v2/cmd"
	"github.com/micro/micro/v2/service/handler/exec"
	"github.com/micro/micro/v2/service/handler/file"

	// used to setup defaults
	gcli "github.com/micro/go-micro/v2/client/grpc"
	memTrace "github.com/micro/go-micro/v2/debug/trace/memory"
	gsrv "github.com/micro/go-micro/v2/server/grpc"
	memoryStore "github.com/micro/go-micro/v2/store/memory"
)

// Service is a Micro Service which honours the go-micro/service interface
type Service struct {
	opts Options

	once sync.Once
}

func init() {
	// default client
	client.DefaultClient = gcli.NewClient()
	// default server
	server.DefaultServer = gsrv.NewServer()
	// default store
	store.DefaultStore = memoryStore.NewStore()
	// set default trace
	trace.DefaultTracer = memTrace.NewTracer()
	// register the commands
	cmd.Commands = append(app.Commands, Commands()...)
}

// New returns a new Micro Service
func New(opts ...Option) *Service {
	service := new(service)
	options := newOptions(opts...)

	// service name
	serviceName := options.Server.Options().Name

	// we pass functions to the wrappers since the values can change during initialisation
	authFn := func() auth.Auth { return options.Server.Options().Auth }
	cacheFn := func() *client.Cache { return options.Client.Options().Cache }

	// wrap client to inject From-Service header on any calls
	options.Client = wrapper.FromService(serviceName, options.Client)
	options.Client = wrapper.TraceCall(serviceName, trace.DefaultTracer, options.Client)
	options.Client = wrapper.CacheClient(cacheFn, options.Client)
	options.Client = wrapper.AuthClient(authFn, options.Client)

	// wrap the server to provide handler stats
	options.Server.Init(
		server.WrapHandler(wrapper.HandlerStats(stats.DefaultStats)),
		server.WrapHandler(wrapper.TraceHandler(trace.DefaultTracer)),
		server.WrapHandler(wrapper.AuthHandler(authFn)),
	)

	// set opts
	service.opts = options

	return service
}

func (s *service) Name() string {
	return s.opts.Server.Options().Name
}

// Init initialises options. Additionally it calls cmd.Init
// which parses command line flags. cmd.Init is only called
// on first Init.
func (s *service) Init(opts ...Option) {
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

		// Initialise the command flags, overriding new service
		if err := s.opts.Cmd.Init(
			cmd.Auth(&s.opts.Auth),
			cmd.Broker(&s.opts.Broker),
			cmd.Registry(&s.opts.Registry),
			cmd.Runtime(&s.opts.Runtime),
			cmd.Transport(&s.opts.Transport),
			cmd.Client(&s.opts.Client),
			cmd.Config(&s.opts.Config),
			cmd.Server(&s.opts.Server),
			cmd.Store(&s.opts.Store),
			cmd.Profile(&s.opts.Profile),
		); err != nil {
			logger.Fatal(err)
		}

		// Explicitly set the table name to the service name
		name := s.opts.Cmd.App().Name
		s.opts.Store.Init(store.Table(name))
	})
}

func (s *service) Options() Options {
	return s.opts
}

func (s *service) Client() client.Client {
	return s.opts.Client
}

func (s *service) Server() server.Server {
	return s.opts.Server
}

func (s *service) String() string {
	return "micro"
}

func (s *service) Start() error {
	for _, fn := range s.opts.BeforeStart {
		if err := fn(); err != nil {
			return err
		}
	}

	if err := s.opts.Server.Start(); err != nil {
		return err
	}

	for _, fn := range s.opts.AfterStart {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) Stop() error {
	var gerr error

	for _, fn := range s.opts.BeforeStop {
		if err := fn(); err != nil {
			gerr = err
		}
	}

	if err := s.opts.Server.Stop(); err != nil {
		return err
	}

	for _, fn := range s.opts.AfterStop {
		if err := fn(); err != nil {
			gerr = err
		}
	}

	return gerr
}

func (s *service) Run() error {
	// register the debug handler
	s.opts.Server.Handle(
		s.opts.Server.NewHandler(
			handler.NewHandler(s.opts.Client),
			server.InternalHandler(true),
		),
	)

	// start the profiler
	if s.opts.Profile != nil {
		// to view mutex contention
		rtime.SetMutexProfileFraction(5)
		// to view blocking profile
		rtime.SetBlockProfileRate(1)

		if err := s.opts.Profile.Start(); err != nil {
			return err
		}
		defer s.opts.Profile.Stop()
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

	select {
	// wait on kill signal
	case <-ch:
	// wait on context cancel
	case <-s.opts.Context.Done():
	}

	return s.Stop()
}

// Run starts a micro service sidecar to encapsulate any app
func Run(ctx *cli.Context, opts ...micro.Option) {
	log.Init(log.WithFields(map[string]interface{}{"service": "service"}))

	name := ctx.String("name")
	address := ctx.String("address")
	endpoint := ctx.String("endpoint")

	metadata := make(map[string]string)
	for _, md := range ctx.StringSlice("metadata") {
		parts := strings.Split(md, "=")
		if len(parts) < 2 {
			continue
		}

		key := parts[0]
		val := strings.Join(parts[1:], "=")

		// set the key/val
		metadata[key] = val
	}

	if len(metadata) > 0 {
		opts = append(opts, micro.Metadata(metadata))
	}

	if len(name) > 0 {
		opts = append(opts, micro.Name(name))
	}

	if len(address) > 0 {
		opts = append(opts, micro.Address(address))
	}

	if len(endpoint) == 0 {
		endpoint = proxy.DefaultEndpoint
	}

	var p proxy.Proxy

	switch {
	case strings.HasPrefix(endpoint, "grpc"):
		endpoint = strings.TrimPrefix(endpoint, "grpc://")
		p = grpc.NewProxy(proxy.WithEndpoint(endpoint))
	case strings.HasPrefix(endpoint, "http"):
		p = http.NewProxy(proxy.WithEndpoint(endpoint))
	case strings.HasPrefix(endpoint, "file"):
		endpoint = strings.TrimPrefix(endpoint, "file://")
		p = file.NewProxy(proxy.WithEndpoint(endpoint))
	case strings.HasPrefix(endpoint, "exec"):
		endpoint = strings.TrimPrefix(endpoint, "exec://")
		p = exec.NewProxy(proxy.WithEndpoint(endpoint))
	default:
		p = mucp.NewProxy(proxy.WithEndpoint(endpoint))
	}

	// run the service if asked to
	if ctx.Args().Len() > 0 {
		args := []runtime.CreateOption{
			runtime.WithCommand(ctx.Args().Slice()...),
			runtime.WithOutput(os.Stdout),
		}

		// create new local runtime
		r := runtime.NewRuntime()

		// start the runtime
		r.Start()

		// register the service
		r.Create(&runtime.Service{
			Name: name,
		}, args...)

		// stop the runtime
		defer func() {
			r.Delete(&runtime.Service{
				Name: name,
			})
			r.Stop()
		}()
	}

	log.Infof("Service [%s] Serving %s at endpoint %s\n", p.String(), name, endpoint)

	// new service
	service := NewService(opts...)

	// create new muxer
	//	muxer := mux.New(name, p)

	// set the router
	service.Server().Init(
		server.WithRouter(p),
	)

	// run service
	service.Run()
}

// Commands sets up the CLI command for `micro service`
func Commands(options ...micro.Option) []*cli.Command {
	command := &cli.Command{
		Name:  "service",
		Usage: "Run a micro service",
		Action: func(ctx *cli.Context) error {
			Run(ctx, options...)
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Usage:   "Name of the service",
				EnvVars: []string{"MICRO_SERVICE_NAME"},
				Value:   "service",
			},
			&cli.StringFlag{
				Name:    "address",
				Usage:   "Address of the service",
				EnvVars: []string{"MICRO_SERVICE_ADDRESS"},
			},
			&cli.StringFlag{
				Name:    "endpoint",
				Usage:   "The local service endpoint (Defaults to localhost:9090); {http, grpc, file, exec}://path-or-address e.g http://localhost:9090",
				EnvVars: []string{"MICRO_SERVICE_ENDPOINT"},
			},
			&cli.StringSliceFlag{
				Name:    "metadata",
				Usage:   "Add metadata as key-value pairs describing the service e.g owner=john@example.com",
				EnvVars: []string{"MICRO_SERVICE_METADATA"},
			},
		},
	}

	for _, p := range Plugins() {
		if cmds := p.Commands(); len(cmds) > 0 {
			command.Subcommands = append(command.Subcommands, cmds...)
		}

		if flags := p.Flags(); len(flags) > 0 {
			command.Flags = append(command.Flags, flags...)
		}
	}

	return []*cli.Command{command}
}
