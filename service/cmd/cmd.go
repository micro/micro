package cmd

import (
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v3/cmd"
)

var (
	beforeFuncs  []cli.BeforeFunc
	DefaultCmd   cmd.Cmd = newCmd()
	DefaultFlags         = []cli.Flag{
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

type command struct {
	app *cli.App
}

func newCmd() cmd.Cmd {
	cmd := new(command)
	cmd.app = cli.NewApp()
	cmd.app.Flags = DefaultFlags
	cmd.app.Before = cmd.Before
	return cmd
}

func (c *command) App() *cli.App {
	return c.app
}

func (c *command) Options() cmd.Options {
	return cmd.Options{}
}

// After is executed after any subcommand
func (c *command) After(ctx *cli.Context) error {
	return nil
}

// Before is executed before any subcommand
func (c *command) Before(ctx *cli.Context) error {
	for _, fn := range beforeFuncs {
		if err := fn(ctx); err != nil {
			return err
		}
	}
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

// Register a function which is called before an action is executed
func Register(fn cli.BeforeFunc) {
	beforeFuncs = append(beforeFuncs, fn)
}

// // use the proxy address passed as a flag, this is normally
// 	// the micro network
// 	if proxy := ctx.String("proxy_address"); len(proxy) > 0 {
// 		muclient.DefaultClient.Init(client.Proxy(proxy))
// 	}

// 	// wrap the client
// 	muclient.DefaultClient = wrapper.AuthClient(muclient.DefaultClient)
// 	muclient.DefaultClient = wrapper.CacheClient(muclient.DefaultClient)
// 	muclient.DefaultClient = wrapper.TraceCall(muclient.DefaultClient)
// 	muclient.DefaultClient = wrapper.FromService(muclient.DefaultClient)

// 	// wrap the server
// 	muserver.DefaultServer.Init(
// 		server.WrapHandler(wrapper.AuthHandler()),
// 		server.WrapHandler(wrapper.TraceHandler()),
// 		server.WrapHandler(wrapper.HandlerStats()),
// 	)

// 	// setup auth
// 	authOpts := []auth.Option{}
// 	if len(ctx.String("namespace")) > 0 {
// 		authOpts = append(authOpts, auth.Issuer(ctx.String("namespace")))
// 	}
// 	if len(ctx.String("auth_id")) > 0 || len(ctx.String("auth_secret")) > 0 {
// 		authOpts = append(authOpts, auth.Credentials(
// 			ctx.String("auth_id"), ctx.String("auth_secret"),
// 		))
// 	}
// 	muauth.DefaultAuth.Init(authOpts...)

// 	// Setup store options
// 	storeDB := store.Database(ctx.String("namespace"))
// 	if err := mustore.DefaultStore.Init(storeDB); err != nil {
// 		logger.Fatalf("Error configuring store: %v", err)
// 	}

// 	// set the registry in the client and server
// 	muclient.DefaultClient.Init(client.Registry(muregistry.DefaultRegistry))
// 	muserver.DefaultServer.Init(server.Registry(muregistry.DefaultRegistry))

// 	// setup auth credentials and refresh token periodically
// 	if err := setupAuthForService(); err != nil {
// 		logger.Fatalf("Error setting up auth: %v", err)
// 	}
// 	go refreshAuthToken(c.exit)

// 	// Setup config. Do this after auth is configured since it'll load the config
// 	// from the service immediately
// 	conf, err := config.NewConfig(config.WithSource(configCli.NewSource()))
// 	if err != nil {
// 		logger.Fatalf("Error configuring config: %v", err)
// 	}
// 	muconfig.DefaultConfig = conf
// 	return nil

// before extracts service options from the CLI flags. These
// aren't set by the cmd package to prevent a circular dependancy.
// prepend them to the array so options passed by the user to this
// function are applied after (taking precedence)
// before := func(ctx *cli.Context) error {
// 	var opts []service.Option
// 	if n := ctx.String("service_name"); len(n) > 0 {
// 		opts = append(opts, Name(n))
// 	}
// 	if v := ctx.String("service_version"); len(v) > 0 {
// 		opts = append(opts, Version(v))
// 	}
// 	return nil
// }
