package service

import (
	"sync"

	"github.com/micro/go-micro/v3/client"
	"github.com/micro/go-micro/v3/cmd"
	"github.com/micro/go-micro/v3/config"
	"github.com/micro/go-micro/v3/server"
	"github.com/micro/go-micro/v3/store"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v3/auth"
	"github.com/micro/micro/v3/internal/wrapper"
	"github.com/micro/micro/v3/service/logger"

	muauth "github.com/micro/micro/v3/service/auth"
	muclient "github.com/micro/micro/v3/service/client"
	muconfig "github.com/micro/micro/v3/service/config"
	configCli "github.com/micro/micro/v3/service/config/client"
	muregistry "github.com/micro/micro/v3/service/registry"
	muserver "github.com/micro/micro/v3/service/server"
	mustore "github.com/micro/micro/v3/service/store"
)

type command struct {
	app *cli.App

	sync.Mutex
	// exit is a channel which is closed
	// on exit for anything that requires
	// cleanup
	exit chan bool
}

var (
	defaultCmd cmd.Cmd

	defaultFlags = []cli.Flag{
		&cli.StringFlag{
			Name:    "namespace",
			EnvVars: []string{"MICRO_NAMESPACE"},
			Usage:   "Namespace the service is operating in",
			Value:   "micro",
		},
		&cli.StringFlag{
			Name:    "auth_id",
			EnvVars: []string{"MICRO_AUTH_ID"},
			Usage:   "Account ID used for client authentication",
		},
		&cli.StringFlag{
			Name:    "auth_secret",
			EnvVars: []string{"MICRO_AUTH_SECRET"},
			Usage:   "Account secret used for client authentication",
		},
		&cli.StringFlag{
			Name:    "proxy_address",
			Usage:   "Proxy requests via the HTTP address specified",
			EnvVars: []string{"MICRO_PROXY"},
		},
		&cli.StringFlag{
			Name:    "service_name",
			Usage:   "Name of the micro service",
			EnvVars: []string{"MICRO_SERVICE_NAME"},
		},
		&cli.StringFlag{
			Name:    "service_version",
			Usage:   "Version of the micro service",
			EnvVars: []string{"MICRO_SERVICE_VERSION"},
		},
	}
)

func newCmd() cmd.Cmd {
	cmd := new(command)
	cmd.exit = make(chan bool)
	cmd.app = cli.NewApp()
	cmd.app.Usage = description
	cmd.app.Flags = defaultFlags
	cmd.app.Before = cmd.Before
	return cmd
}

func (c *command) App() *cli.App {
	return c.app
}

func (c *command) Options() cmd.Options {
	return nil
}

// After is executed after any subcommand
func (c *command) After(ctx *cli.Context) error {
	return nil
}

// Before is executed before any subcommand
func (c *command) Before(ctx *cli.Context) error {
	// use the proxy address passed as a flag, this is normally
	// the micro network
	if proxy := ctx.String("proxy_address"); len(proxy) > 0 {
		muclient.DefaultClient.Init(client.Proxy(proxy))
	}

	// wrap the client
	muclient.DefaultClient = wrapper.AuthClient(muclient.DefaultClient)
	muclient.DefaultClient = wrapper.CacheClient(muclient.DefaultClient)
	muclient.DefaultClient = wrapper.TraceCall(muclient.DefaultClient)
	muclient.DefaultClient = wrapper.FromService(muclient.DefaultClient)

	// wrap the server
	muserver.DefaultServer.Init(
		server.WrapHandler(wrapper.AuthHandler()),
		server.WrapHandler(wrapper.TraceHandler()),
		server.WrapHandler(wrapper.HandlerStats()),
	)

	// setup auth
	authOpts := []auth.Option{}
	if len(ctx.String("namespace")) > 0 {
		authOpts = append(authOpts, auth.Issuer(ctx.String("namespace")))
	}
	if len(ctx.String("auth_id")) > 0 || len(ctx.String("auth_secret")) > 0 {
		authOpts = append(authOpts, auth.Credentials(
			ctx.String("auth_id"), ctx.String("auth_secret"),
		))
	}
	muauth.DefaultAuth.Init(authOpts...)

	// Setup store options
	storeDB := store.Database(ctx.String("namespace"))
	if err := mustore.DefaultStore.Init(storeDB); err != nil {
		logger.Fatalf("Error configuring store: %v", err)
	}

	// set the registry in the client and server
	muclient.DefaultClient.Init(client.Registry(muregistry.DefaultRegistry))
	muserver.DefaultServer.Init(server.Registry(muregistry.DefaultRegistry))

	// setup auth credentials and refresh token periodically
	if err := setupAuthForService(); err != nil {
		logger.Fatalf("Error setting up auth: %v", err)
	}
	go refreshAuthToken(c.exit)

	// Setup config. Do this after auth is configured since it'll load the config
	// from the service immediately
	conf, err := config.NewConfig(config.WithSource(configCli.NewSource()))
	if err != nil {
		logger.Fatalf("Error configuring config: %v", err)
	}
	muconfig.DefaultConfig = conf
	return nil
}

func (c *command) Init(opts ...cmd.Option) error {
	return nil
}

func (c *command) Run() error {
	return nil
}

func (c *command) String() string {
	return "service"
}
