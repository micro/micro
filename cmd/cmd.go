package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/micro/go-micro/v3/broker"
	"github.com/micro/go-micro/v3/client"
	"github.com/micro/go-micro/v3/config"
	"github.com/micro/go-micro/v3/runtime"
	"github.com/micro/go-micro/v3/server"
	"github.com/micro/go-micro/v3/store"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v3/auth"
	"github.com/micro/go-micro/v3/cmd"
	"github.com/micro/go-micro/v3/logger"
	"github.com/micro/go-micro/v3/registry"
	"github.com/micro/micro/v2/client/cli/util"
	"github.com/micro/micro/v2/internal/helper"
	_ "github.com/micro/micro/v2/internal/usage"
	"github.com/micro/micro/v2/internal/wrapper"
	"github.com/micro/micro/v2/plugin"
	"github.com/micro/micro/v2/profile"

	authCli "github.com/micro/micro/v2/service/auth/client"
	brokerCli "github.com/micro/micro/v2/service/broker/client"
	configCli "github.com/micro/micro/v2/service/config/client"
	registryCli "github.com/micro/micro/v2/service/registry/client"
	storeCli "github.com/micro/micro/v2/service/store/client"

	muauth "github.com/micro/micro/v2/service/auth"
	mubroker "github.com/micro/micro/v2/service/broker"
	muclient "github.com/micro/micro/v2/service/client"
	muconfig "github.com/micro/micro/v2/service/config"
	muregistry "github.com/micro/micro/v2/service/registry"
	muruntime "github.com/micro/micro/v2/service/runtime"
	muserver "github.com/micro/micro/v2/service/server"
	mustore "github.com/micro/micro/v2/service/store"
)

type command struct {
	opts cmd.Options
	app  *cli.App
}

var (
	DefaultCmd cmd.Cmd = New()

	// name of the binary
	name = "micro"
	// description of the binary
	description = "A framework for cloud native development\n\n	 Use `micro [command] --help` to see command specific help."
	// defaultFlags which are used on all commands
	defaultFlags = []cli.Flag{
		&cli.StringFlag{
			Name:    "profile",
			Usage:   "Set the micro profile: e.g. local or platform",
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
		&cli.StringFlag{
			Name:    "update_url",
			Usage:   "Set the url to retrieve system updates from",
			EnvVars: []string{"MICRO_UPDATE_URL"},
			Value:   "https://micro.mu/update",
		},
		&cli.BoolFlag{
			Name:    "report_usage",
			Usage:   "Report usage statistics",
			EnvVars: []string{"MICRO_REPORT_USAGE"},
			Value:   true,
		},
		&cli.StringFlag{
			Name:    "env",
			Aliases: []string{"e"},
			Usage:   "Override environment",
			EnvVars: []string{"MICRO_ENV"},
		},
	}
)

func init() {
	rand.Seed(time.Now().Unix())
}

func New(opts ...cmd.Option) cmd.Cmd {
	options := cmd.Options{}
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
	cmd.app.Before = cmd.Before

	// run a custom action, this allows us to run a service
	// after parsing the cli flags and setting up micro
	if action := actionFromContext(options.Context); action != nil {
		cmd.app.Action = action
	}

	return cmd
}

func (c *command) App() *cli.App {
	return c.app
}

func (c *command) Options() cmd.Options {
	return c.opts
}

func (c *command) Before(ctx *cli.Context) error {
	// set the proxy address. TODO: Refactor to be a client option.
	util.SetProxyAddress(ctx)

	// initialize plugins
	for _, p := range plugin.Plugins() {
		if err := p.Init(ctx); err != nil {
			return err
		}
	}

	// default the profile for the server
	prof := ctx.String("profile")
	a := ctx.Args().First()
	if len(prof) == 0 && (a == "service" || a == "server") {
		prof = "local"
	}

	// apply the profile
	if profile, ok := profile.Profiles[prof]; ok {
		profile()
		logger.Infof("Configuring micro with the %v profile", prof)
	} else if len(prof) > 0 {
		logger.Fatalf("Unknown profile: %v", prof)
	}

	// wrap the client
	muclient.DefaultClient = wrapper.AuthClient(muclient.DefaultClient)
	muclient.DefaultClient = wrapper.CacheClient(muclient.DefaultClient)
	muclient.DefaultClient = wrapper.TraceCall(muclient.DefaultClient)

	// wrap the server
	muserver.DefaultServer.Init(
		server.WrapHandler(wrapper.AuthHandler()),
		server.WrapHandler(wrapper.TraceHandler()),
		server.WrapHandler(wrapper.HandlerStats()),
	)

	// setup auth
	authOpts := []auth.Option{authCli.WithClient(muclient.DefaultClient)}
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
	if len(ctx.String("auth_public_key")) > 0 {
		authOpts = append(authOpts, auth.PublicKey(ctx.String("auth_public_key")))
	}
	if len(ctx.String("auth_private_key")) > 0 {
		authOpts = append(authOpts, auth.PrivateKey(ctx.String("auth_private_key")))
	}
	muauth.DefaultAuth.Init(authOpts...)

	// setup registry
	registryOpts := []registry.Option{registryCli.WithClient(muclient.DefaultClient)}

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
	if err := muregistry.DefaultRegistry.Init(registryOpts...); err != nil {
		logger.Fatalf("Error configuring registry: %v", err)
	}

	// Setup broker options.
	brokerOpts := []broker.Option{brokerCli.WithClient(muclient.DefaultClient)}
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
	if err := mubroker.DefaultBroker.Init(brokerOpts...); err != nil {
		logger.Fatalf("Error configuring broker: %v", err)
	}

	// Setup store options
	storeOpts := []store.Option{storeCli.WithClient(muclient.DefaultClient)}
	if len(ctx.String("store_address")) > 0 {
		storeOpts = append(storeOpts, store.Nodes(strings.Split(ctx.String("store_address"), ",")...))
	}
	if len(ctx.String("namespace")) > 0 {
		storeOpts = append(storeOpts, store.Database(ctx.String("namespace")))
	}
	if err := mustore.DefaultStore.Init(storeOpts...); err != nil {
		logger.Fatalf("Error configuring store: %v", err)
	}

	// Set runtime client
	if err := muruntime.DefaultRuntime.Init(runtime.WithClient(muclient.DefaultClient)); err != nil {
		logger.Fatalf("Error configuring runtime: %v", err)
	}

	// set the registry in the client and server
	muclient.DefaultClient.Init(client.Registry(muregistry.DefaultRegistry))
	muserver.DefaultServer.Init(server.Registry(muregistry.DefaultRegistry))

	// set the credentials from the CLI. If a service is run, it'll override
	// these when it's started.
	if err := util.SetAuthToken(ctx); err != nil {
		return err
	}

	// Setup config. Do this after auth is configured since it'll load the config
	// from the service immediately.
	conf, err := config.NewConfig(config.WithSource(configCli.NewSource()))
	if err != nil {
		logger.Fatalf("Error configuring config: %v", err)
	}
	muconfig.DefaultConfig = conf

	return nil
}

func (c *command) Init(opts ...cmd.Option) error {
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

	return nil
}

func (c *command) Run() error {
	return c.app.Run(os.Args)
}

func (c *command) String() string {
	return "micro"
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
		// firstly check to see if the service "go.micro.config"
		// exists within the current namespace, then it would
		// execute the Config.Set RPC, setting the flags in the
		// request.
		if srv, err := lookupService(c); err != nil {
			fmt.Printf("Error querying registry for service: %v", err)
			os.Exit(1)
		} else if srv != nil && c.Args().Len() == 1 {
			fmt.Println(formatServiceUsage(srv, c.Args().First()))
			os.Exit(1)
		} else if srv != nil {
			if err := callService(srv, c); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			os.Exit(0)
		}

	}

	return helper.MissingCommand(c)
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
