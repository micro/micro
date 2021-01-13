package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"runtime/debug"
	"sort"
	"strings"
	"time"

	"github.com/micro/micro/v3/client/cli/util"
	uconf "github.com/micro/micro/v3/internal/config"
	"github.com/micro/micro/v3/internal/helper"
	"github.com/micro/micro/v3/internal/network"
	"github.com/micro/micro/v3/internal/report"
	_ "github.com/micro/micro/v3/internal/usage"
	"github.com/micro/micro/v3/internal/user"
	"github.com/micro/micro/v3/internal/wrapper"
	"github.com/micro/micro/v3/plugin"
	"github.com/micro/micro/v3/profile"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/broker"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/config"
	configCli "github.com/micro/micro/v3/service/config/client"
	storeConf "github.com/micro/micro/v3/service/config/store"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/registry"
	"github.com/micro/micro/v3/service/server"
	"github.com/micro/micro/v3/service/store"
	"github.com/urfave/cli/v2"

	muruntime "github.com/micro/micro/v3/service/runtime"
)

type Cmd interface {
	// Init initialises options
	// Note: Use Run to parse command line
	Init(opts ...Option) error
	// Options set within this command
	Options() Options
	// The cli app within this cmd
	App() *cli.App
	// Run executes the command
	Run() error
	// Implementation
	String() string
}

type command struct {
	opts Options
	app  *cli.App

	// before is a function which should
	// be called in Before if not nil
	before cli.ActionFunc

	// indicates whether this is a service
	service bool
}

var (
	DefaultCmd Cmd = New()

	// name of the binary
	name = "micro"
	// description of the binary
	description = "A framework for cloud native development\n\n	 Use `micro [command] --help` to see command specific help."
	// defaultFlags which are used on all commands
	defaultFlags = []cli.Flag{
		&cli.StringFlag{
			Name:    "c",
			Usage:   "Set the config file: Defaults to ~/.micro/config.json",
			EnvVars: []string{"MICRO_CONFIG_FILE"},
		},
		&cli.StringFlag{
			Name:    "env",
			Aliases: []string{"e"},
			Usage:   "Set the environment to operate in",
			EnvVars: []string{"MICRO_ENV"},
		},
		&cli.StringFlag{
			Name:    "profile",
			Usage:   "Set the micro server profile: e.g. local or kubernetes",
			EnvVars: []string{"MICRO_PROFILE"},
		},
		&cli.StringFlag{
			Name:    "namespace",
			EnvVars: []string{"MICRO_NAMESPACE"},
			Usage:   "Namespace the service is operating in",
			Value:   "micro",
		},
		&cli.StringFlag{
			Name:    "auth_address",
			EnvVars: []string{"MICRO_AUTH_ADDRESS"},
			Usage:   "Comma-separated list of auth addresses",
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
			Name:    "auth_public_key",
			EnvVars: []string{"MICRO_AUTH_PUBLIC_KEY"},
			Usage:   "Public key for JWT auth (base64 encoded PEM)",
		},
		&cli.StringFlag{
			Name:    "auth_private_key",
			EnvVars: []string{"MICRO_AUTH_PRIVATE_KEY"},
			Usage:   "Private key for JWT auth (base64 encoded PEM)",
		},
		&cli.StringFlag{
			Name:    "registry_address",
			EnvVars: []string{"MICRO_REGISTRY_ADDRESS"},
			Usage:   "Comma-separated list of registry addresses",
		},
		&cli.StringFlag{
			Name:    "registry_tls_ca",
			Usage:   "Certificate authority for TLS with registry",
			EnvVars: []string{"MICRO_REGISTRY_TLS_CA"},
		},
		&cli.StringFlag{
			Name:    "registry_tls_cert",
			Usage:   "Client cert for TLS with registry",
			EnvVars: []string{"MICRO_REGISTRY_TLS_CERT"},
		},
		&cli.StringFlag{
			Name:    "registry_tls_key",
			Usage:   "Client key for TLS with registry",
			EnvVars: []string{"MICRO_REGISTRY_TLS_KEY"},
		},
		&cli.StringFlag{
			Name:    "broker_address",
			EnvVars: []string{"MICRO_BROKER_ADDRESS"},
			Usage:   "Comma-separated list of broker addresses",
		},
		&cli.StringFlag{
			Name:    "events_tls_ca",
			Usage:   "Certificate authority for TLS with events",
			EnvVars: []string{"MICRO_EVENTS_TLS_CA"},
		},
		&cli.StringFlag{
			Name:    "events_tls_cert",
			Usage:   "Client cert for TLS with events",
			EnvVars: []string{"MICRO_EVENTS_TLS_CERT"},
		},
		&cli.StringFlag{
			Name:    "events_tls_key",
			Usage:   "Client key for TLS with events",
			EnvVars: []string{"MICRO_EVENTS_TLS_KEY"},
		},
		&cli.StringFlag{
			Name:    "broker_tls_ca",
			Usage:   "Certificate authority for TLS with broker",
			EnvVars: []string{"MICRO_BROKER_TLS_CA"},
		},
		&cli.StringFlag{
			Name:    "broker_tls_cert",
			Usage:   "Client cert for TLS with broker",
			EnvVars: []string{"MICRO_BROKER_TLS_CERT"},
		},
		&cli.StringFlag{
			Name:    "broker_tls_key",
			Usage:   "Client key for TLS with broker",
			EnvVars: []string{"MICRO_BROKER_TLS_KEY"},
		},
		&cli.StringFlag{
			Name:    "store_address",
			EnvVars: []string{"MICRO_STORE_ADDRESS"},
			Usage:   "Comma-separated list of store addresses",
		},
		&cli.StringFlag{
			Name:    "proxy_address",
			Usage:   "Proxy requests via the HTTP address specified",
			EnvVars: []string{"MICRO_PROXY"},
		},
		&cli.BoolFlag{
			Name:    "report_usage",
			Usage:   "Report usage statistics",
			EnvVars: []string{"MICRO_REPORT_USAGE"},
			Value:   true,
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
		&cli.StringFlag{
			Name:    "service_address",
			Usage:   "Address to run the service on",
			EnvVars: []string{"MICRO_SERVICE_ADDRESS"},
		},
		&cli.BoolFlag{
			Name:    "prompt_update",
			Usage:   "Provide an update prompt when a new binary is available. Enabled for release binaries only.",
			Value:   true,
			EnvVars: []string{"MICRO_PROMPT_UPDATE"},
		},
		&cli.StringFlag{
			Name:    "config_secret_key",
			Usage:   "Key to use when encoding/decoding secret config values. Will be generated and saved to file if not provided.",
			Value:   "",
			EnvVars: []string{"MICRO_CONFIG_SECRET_KEY"},
		},
	}
)

func init() {
	rand.Seed(time.Now().Unix())

	// configure defaults for all packages
	setupDefaults()
}

func action(c *cli.Context) error {
	if c.Args().Len() > 0 {
		// if an executable is available with the name of
		// the command, execute it with the arguments from
		// index 1 on.
		v, err := exec.LookPath("micro-" + c.Args().First())
		if err == nil {
			ce := exec.Command(v, c.Args().Slice()[1:]...)
			ce.Stdout = os.Stdout
			ce.Stderr = os.Stderr
			return ce.Run()
		}

		// lookup the service, e.g. "micro config set" would
		// firstly check to see if the service, e.g. config
		// exists within the current namespace, then it would
		// execute the Config.Set RPC, setting the flags in the
		// request.
		if srv, ns, err := lookupService(c); err != nil {
			return util.CliError(err)
		} else if srv != nil && shouldRenderHelp(c) {
			return cli.Exit(formatServiceUsage(srv, c), 0)
		} else if srv != nil {
			err := callService(srv, ns, c)
			return util.CliError(err)
		}

		// srv == nil
		return helper.UnexpectedCommand(c)

	}

	return helper.MissingCommand(c)
}

func New(opts ...Option) *command {
	options := Options{}
	for _, o := range opts {
		o(&options)
	}

	cmd := new(command)
	cmd.opts = options
	cmd.app = cli.NewApp()
	cmd.app.Name = name
	cmd.app.Version = buildVersion()
	cmd.app.Usage = description
	cmd.app.Flags = defaultFlags
	cmd.app.Action = action
	cmd.app.Before = beforeFromContext(options.Context, cmd.Before)

	// if this option has been set, we're running a service
	// and no action needs to be performed. The CMD package
	// is just being used to parse flags and configure micro.
	if setupOnlyFromContext(options.Context) {
		cmd.service = true
		cmd.app.Action = func(ctx *cli.Context) error { return nil }
	}

	//flags to add
	if len(options.Flags) > 0 {
		cmd.app.Flags = append(cmd.app.Flags, options.Flags...)
	}
	//action to replace
	if options.Action != nil {
		cmd.app.Action = options.Action
	}
	// cmd to add to use registry

	return cmd
}

func (c *command) App() *cli.App {
	return c.app
}

func (c *command) Options() Options {
	return c.opts
}

// Before is executed before any subcommand
func (c *command) Before(ctx *cli.Context) error {
	if v := ctx.Args().First(); len(v) > 0 {
		switch v {
		case "service", "server":
			// do nothing
		default:
			// check for the latest release
			// TODO: write a local file to detect
			// when we last checked so we don't do it often
			updated, err := confirmAndSelfUpdate(ctx)
			if err != nil {
				return err
			}
			// if updated we expect to re-execute the command
			// TODO: maybe require relogin or update of the
			// config...
			if updated {
				// considering nil actually continues
				// we need to os.Exit(0)
				os.Exit(0)
				return nil
			}
		}
	}

	// set the config file if specified
	if cf := ctx.String("c"); len(cf) > 0 {
		uconf.SetConfig(cf)
	}

	// initialize plugins
	for _, p := range plugin.Plugins() {
		if err := p.Init(ctx); err != nil {
			return err
		}
	}

	// default the profile for the server
	prof := ctx.String("profile")

	// if no profile is set then set one
	if len(prof) == 0 {
		switch ctx.Args().First() {
		case "service", "server":
			prof = "local"
		default:
			prof = "client"
		}
	}

	// apply the profile
	if profile, err := profile.Load(prof); err != nil {
		logger.Fatal(err)
	} else {
		// load the profile
		profile.Setup(ctx)
	}

	// set the proxy address
	var proxy string
	if c.service || ctx.IsSet("proxy_address") {
		// use the proxy address passed as a flag, this is normally
		// the micro network
		proxy = ctx.String("proxy_address")
	} else {
		// for CLI, use the external proxy which is loaded from the
		// local config
		var err error
		proxy, err = util.CLIProxyAddress(ctx)
		if err != nil {
			return err
		}
	}
	if len(proxy) > 0 {
		client.DefaultClient.Init(client.Proxy(proxy))
	}

	// use the internal network lookup
	client.DefaultClient.Init(
		client.Lookup(network.Lookup),
	)

	// wrap the client
	client.DefaultClient = wrapper.AuthClient(client.DefaultClient)
	client.DefaultClient = wrapper.CacheClient(client.DefaultClient)
	client.DefaultClient = wrapper.TraceCall(client.DefaultClient)
	client.DefaultClient = wrapper.FromService(client.DefaultClient)
	client.DefaultClient = wrapper.LogClient(client.DefaultClient)

	// wrap the server
	server.DefaultServer.Init(
		server.WrapHandler(wrapper.AuthHandler()),
		server.WrapHandler(wrapper.TraceHandler()),
		server.WrapHandler(wrapper.HandlerStats()),
		server.WrapHandler(wrapper.LogHandler()),
		server.WrapHandler(wrapper.MetricsHandler()),
	)

	// setup auth
	authOpts := []auth.Option{}
	if len(ctx.String("namespace")) > 0 {
		authOpts = append(authOpts, auth.Issuer(ctx.String("namespace")))
	}
	if len(ctx.String("auth_address")) > 0 {
		authOpts = append(authOpts, auth.Addrs(ctx.String("auth_address")))
	}
	if len(ctx.String("auth_id")) > 0 || len(ctx.String("auth_secret")) > 0 {
		authOpts = append(authOpts, auth.Credentials(
			ctx.String("auth_id"), ctx.String("auth_secret"),
		))
	}

	// load the jwt private and public keys, in the case of the server we want to generate them if not
	// present. The server will inject these creds into the core services, if the services generated
	// the credentials themselves then they wouldn't match
	if len(ctx.String("auth_public_key")) > 0 || len(ctx.String("auth_private_key")) > 0 {
		authOpts = append(authOpts, auth.PublicKey(ctx.String("auth_public_key")))
		authOpts = append(authOpts, auth.PrivateKey(ctx.String("auth_private_key")))
	} else if ctx.Args().First() == "server" || ctx.Args().First() == "service" {
		privKey, pubKey, err := user.GetJWTCerts()
		if err != nil {
			logger.Fatalf("Error getting keys: %v", err)
		}
		authOpts = append(authOpts, auth.PublicKey(string(pubKey)), auth.PrivateKey(string(privKey)))
	}

	auth.DefaultAuth.Init(authOpts...)

	// setup auth credentials, use local credentials for the CLI and injected creds
	// for the service.
	var err error
	if c.service {
		err = setupAuthForService()
	} else {
		err = setupAuthForCLI(ctx)
	}
	if err != nil {
		logger.Fatalf("Error setting up auth: %v", err)
	}
	go refreshAuthToken()

	// initialize the server with the namespace so it knows which domain to register in
	server.DefaultServer.Init(server.Namespace(ctx.String("namespace")))

	// setup registry
	registryOpts := []registry.Option{}

	// Parse registry TLS certs
	if len(ctx.String("registry_tls_cert")) > 0 || len(ctx.String("registry_tls_key")) > 0 {
		cert, err := tls.LoadX509KeyPair(ctx.String("registry_tls_cert"), ctx.String("registry_tls_key"))
		if err != nil {
			logger.Fatalf("Error loading registry tls cert: %v", err)
		}

		// load custom certificate authority
		caCertPool := x509.NewCertPool()
		if len(ctx.String("registry_tls_ca")) > 0 {
			crt, err := ioutil.ReadFile(ctx.String("registry_tls_ca"))
			if err != nil {
				logger.Fatalf("Error loading registry tls certificate authority: %v", err)
			}
			caCertPool.AppendCertsFromPEM(crt)
		}

		cfg := &tls.Config{Certificates: []tls.Certificate{cert}, RootCAs: caCertPool}
		registryOpts = append(registryOpts, registry.TLSConfig(cfg))
	}
	if len(ctx.String("registry_address")) > 0 {
		addresses := strings.Split(ctx.String("registry_address"), ",")
		registryOpts = append(registryOpts, registry.Addrs(addresses...))
	}
	if err := registry.DefaultRegistry.Init(registryOpts...); err != nil {
		logger.Fatalf("Error configuring registry: %v", err)
	}

	// Setup broker options.
	brokerOpts := []broker.Option{}
	if len(ctx.String("broker_address")) > 0 {
		brokerOpts = append(brokerOpts, broker.Addrs(ctx.String("broker_address")))
	}

	// Parse broker TLS certs
	if len(ctx.String("broker_tls_cert")) > 0 || len(ctx.String("broker_tls_key")) > 0 {
		cert, err := tls.LoadX509KeyPair(ctx.String("broker_tls_cert"), ctx.String("broker_tls_key"))
		if err != nil {
			logger.Fatalf("Error loading broker TLS cert: %v", err)
		}

		// load custom certificate authority
		caCertPool := x509.NewCertPool()
		if len(ctx.String("broker_tls_ca")) > 0 {
			crt, err := ioutil.ReadFile(ctx.String("broker_tls_ca"))
			if err != nil {
				logger.Fatalf("Error loading broker TLS certificate authority: %v", err)
			}
			caCertPool.AppendCertsFromPEM(crt)
		}

		cfg := &tls.Config{Certificates: []tls.Certificate{cert}, RootCAs: caCertPool}
		brokerOpts = append(brokerOpts, broker.TLSConfig(cfg))
	}
	if err := broker.DefaultBroker.Init(brokerOpts...); err != nil {
		logger.Fatalf("Error configuring broker: %v", err)
	}
	if err := broker.DefaultBroker.Connect(); err != nil {
		logger.Fatalf("Error connecting to broker: %v", err)
	}

	// Setup runtime. This is a temporary fix to trigger the runtime to recreate
	// its client now the client has been replaced with a wrapped one.
	if err := muruntime.DefaultRuntime.Init(); err != nil {
		logger.Fatalf("Error configuring runtime: %v", err)
	}

	// Setup store options
	storeOpts := []store.StoreOption{}
	if len(ctx.String("store_address")) > 0 {
		storeOpts = append(storeOpts, store.Nodes(strings.Split(ctx.String("store_address"), ",")...))
	}
	if len(ctx.String("namespace")) > 0 {
		storeOpts = append(storeOpts, store.Database(ctx.String("namespace")))
	}
	if len(ctx.String("service_name")) > 0 {
		storeOpts = append(storeOpts, store.Table(ctx.String("service_name")))
	}
	if err := store.DefaultStore.Init(storeOpts...); err != nil {
		logger.Fatalf("Error configuring store: %v", err)
	}

	// set the registry and broker in the client and server
	client.DefaultClient.Init(
		client.Broker(broker.DefaultBroker),
		client.Registry(registry.DefaultRegistry),
	)
	server.DefaultServer.Init(
		server.Broker(broker.DefaultBroker),
		server.Registry(registry.DefaultRegistry),
	)

	// Setup config. Do this after auth is configured since it'll load the config
	// from the service immediately. We only do this if the action is nil, indicating
	// a service is being run
	if c.service && config.DefaultConfig == nil {
		config.DefaultConfig = configCli.NewConfig(ctx.String("namespace"))
	} else if config.DefaultConfig == nil {
		config.DefaultConfig, _ = storeConf.NewConfig(store.DefaultStore, ctx.String("namespace"))
	}

	return nil
}

func (c *command) Init(opts ...Option) error {
	for _, o := range opts {
		o(&c.opts)
	}
	if len(c.opts.Name) > 0 {
		c.app.Name = c.opts.Name
	}
	if len(c.opts.Version) > 0 {
		c.app.Version = c.opts.Version
	}
	c.app.HideVersion = len(c.opts.Version) == 0
	c.app.Usage = c.opts.Description

	//allow user's flags to add
	if len(c.opts.Flags) > 0 {
		c.app.Flags = append(c.app.Flags, c.opts.Flags...)
	}
	//action to replace
	if c.opts.Action != nil {
		c.app.Action = c.opts.Action
	}

	return nil
}

func (c *command) Run() error {
	defer func() {
		if r := recover(); r != nil {
			report.Errorf(nil, fmt.Sprintf("panic: %v", string(debug.Stack())))
			panic(r)
		}
	}()
	return c.app.Run(os.Args)
}

func (c *command) String() string {
	return "micro"
}

// Register CLI commands
func Register(cmds ...*cli.Command) {
	app := DefaultCmd.App()
	app.Commands = append(app.Commands, cmds...)

	// sort the commands so they're listed in order on the cli
	// todo: move this to micro/cli so it's only run when the
	// commands are printed during "help"
	sort.Slice(app.Commands, func(i, j int) bool {
		return app.Commands[i].Name < app.Commands[j].Name
	})
}

// Run the default command
func Run() {
	if err := DefaultCmd.Run(); err != nil {
		fmt.Println(formatErr(err))
		os.Exit(1)
	}
}
