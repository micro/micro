package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"time"
	"unicode"

	"github.com/urfave/cli/v2"
	clitoken "micro.dev/v4/cmd/client/token"
	"micro.dev/v4/cmd/client/util"
	"micro.dev/v4/service/auth"
	"micro.dev/v4/service/broker"
	"micro.dev/v4/service/client"
	"micro.dev/v4/service/config"
	configCli "micro.dev/v4/service/config/client"
	storeConf "micro.dev/v4/service/config/store"
	"micro.dev/v4/service/errors"
	"micro.dev/v4/service/logger"
	"micro.dev/v4/service/profile"
	"micro.dev/v4/service/registry"
	"micro.dev/v4/service/runtime"
	"micro.dev/v4/service/server"
	"micro.dev/v4/service/store"
	uauth "micro.dev/v4/util/auth"
	uconf "micro.dev/v4/util/config"
	"micro.dev/v4/util/helper"
	"micro.dev/v4/util/namespace"
)

type Cmd interface {
	// Options set within this command
	Options() Options
	// The cli app within this cmd
	App() *cli.App
	// Run executes the command
	Run() error
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
	description = "API first development platform"
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
			Name:    "namespace",
			EnvVars: []string{"MICRO_NAMESPACE"},
			Usage:   "Namespace the service is operating in",
			Value:   "micro",
		},
		&cli.StringFlag{
			Name:    "client_id",
			EnvVars: []string{"MICRO_CLIENT_ID"},
			Usage:   "Account ID used for client authentication",
		},
		&cli.StringFlag{
			Name:    "client_secret",
			EnvVars: []string{"MICRO_CLIENT_SECRET"},
			Usage:   "Account secret used for client authentication",
		},
		&cli.StringFlag{
			Name:    "public_key",
			EnvVars: []string{"MICRO_PUBLIC_KEY"},
			Usage:   "Public key for JWT auth (base64 encoded PEM)",
		},
		&cli.StringFlag{
			Name:    "private_key",
			EnvVars: []string{"MICRO_PRIVATE_KEY"},
			Usage:   "Private key for JWT auth (base64 encoded PEM)",
		},
		&cli.StringFlag{
			Name:    "name",
			Usage:   "Set the service name",
			EnvVars: []string{"MICRO_SERVICE_NAME"},
		},
		&cli.StringFlag{
			Name:    "profile",
			Usage:   "Set the micro server profile: e.g. local or kubernetes",
			EnvVars: []string{"MICRO_SERVICE_PROFILE"},
		},
		&cli.StringFlag{
			Name:    "network",
			Usage:   "Service network address",
			EnvVars: []string{"MICRO_SERVICE_NETWORK"},
		},
	}
)

func init() {
	rand.Seed(time.Now().Unix())
}

func formatErr(err error) string {
	switch v := err.(type) {
	case *errors.Error:
		return upcaseInitial(v.Detail)
	default:
		return upcaseInitial(err.Error())
	}
}

func upcaseInitial(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func action(c *cli.Context) error {
	if c.Args().Len() == 0 {
		return helper.MissingCommand(c)
	}

	// lookup the service, e.g. "micro config set" would
	// firstly check to see if the service, e.g. config
	// exists within the current namespace, then it would
	// execute the Config.Set RPC, setting the flags in the
	// request.
	if srv, ns, err := util.LookupService(c); err != nil {
		return util.CliError(err)
	} else if srv != nil && util.ShouldRenderHelp(c) {
		return cli.Exit(util.FormatServiceUsage(srv, c), 0)
	} else if srv != nil {
		err := util.CallService(srv, ns, c)
		return util.CliError(err)
	}

	// srv == nil
	return helper.UnexpectedCommand(c)
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
	if serviceFromContext(options.Context) {
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

	return cmd
}

// setupAuth handles exchanging refresh tokens to access tokens
// The structure of the local micro userconfig file is the following:
// micro.auth.[envName].token: temporary access token
// micro.auth.[envName].refresh-token: long lived refresh token
// micro.auth.[envName].expiry: expiration time of the access token, seconds since Unix epoch.
func (c *command) setupAuth(ctx *cli.Context) error {
	if c.service || ctx.Args().First() == "server" {
		return nil
	}

	env, err := util.GetEnv(ctx)
	if err != nil {
		return err
	}
	ns, err := namespace.Get(env.Name)
	if err != nil {
		return err
	}

	tok, err := clitoken.Get(ctx)
	if err != nil {
		return err
	}

	// If there is no refresh token, do not try to refresh it
	if len(tok.RefreshToken) == 0 {
		return nil
	}

	// setup auth token
	// profile.SetupJWT()

	// Check if token is valid
	if time.Now().Before(tok.Expiry.Add(time.Minute * -1)) {
		auth.DefaultAuth.Init(
			auth.ClientToken(tok),
			auth.Issuer(ns),
		)
		return nil
	}

	// Get new access token from refresh token if it's close to expiry
	tok, err = auth.DefaultAuth.Token(
		auth.WithToken(tok.RefreshToken),
		auth.WithTokenIssuer(ns),
		auth.WithExpiry(time.Hour*24),
	)
	if err != nil {
		return nil
	}

	// Save the token to user config file
	auth.DefaultAuth.Init(
		auth.ClientToken(tok),
		auth.Issuer(ns),
	)

	return clitoken.Save(ctx, tok)
}

func (c *command) App() *cli.App {
	return c.app
}

func (c *command) Options() Options {
	return c.opts
}

// Before is executed before any subcommand
func (c *command) Before(ctx *cli.Context) error {
	// set the config file if specified
	if cf := ctx.String("c"); len(cf) > 0 {
		uconf.SetConfig(cf)
	}

	command := ctx.Args().First()

	// certain commands don't require loading
	if command == "env" {
		return nil
	}

	// default the profile for the server
	prof := ctx.String("profile")

	// if no profile is set then set one
	if command == "server" {
		prof = "server"
	} else if len(prof) == 0 {
		prof = "client"
	}

	// apply the profile
	if profile, err := profile.Load(prof); err != nil {
		logger.Fatal(err)
	} else {
		// load the profile
		profile.Setup(ctx)
	}

	// set the proxy address
	var netAddress string
	if c.service || ctx.IsSet("network") {
		// use the proxy address passed as a flag, this is normally
		// the micro network
		netAddress = ctx.String("network")
	} else if command != "server" {
		var err error
		netAddress, err = util.CLIProxyAddress(ctx)
		if err != nil {
			return err
		}
	}
	if len(netAddress) > 0 {
		client.DefaultClient.Init(client.Network(netAddress))
	}

	authOpts := []auth.Option{}
	if len(ctx.String("namespace")) > 0 {
		authOpts = append(authOpts, auth.Issuer(ctx.String("namespace")))
	}
	if len(ctx.String("client_id")) > 0 || len(ctx.String("client_secret")) > 0 {
		authOpts = append(authOpts, auth.Credentials(
			ctx.String("client_id"), ctx.String("client_secret"),
		))
	}

	if len(ctx.String("public_key")) > 0 || len(ctx.String("private_key")) > 0 {
		authOpts = append(authOpts, auth.PublicKey(ctx.String("public_key")))
		authOpts = append(authOpts, auth.PrivateKey(ctx.String("private_key")))
	}

	// setup auth
	auth.DefaultAuth.Init(authOpts...)

	if err := c.setupAuth(ctx); err != nil {
		logger.Fatalf("Error setting up auth: %v", err)
	}

	// refresh the auth token
	go uauth.RefreshToken()

	// initialize the server with the namespace so it knows which domain to register in
	server.DefaultServer.Init(server.Namespace(ctx.String("namespace")))

	if err := registry.DefaultRegistry.Init(); err != nil {
		logger.Fatalf("Error configuring registry: %v", err)
	}

	if err := broker.DefaultBroker.Connect(); err != nil {
		logger.Fatalf("Error connecting to broker: %v", err)
	}

	// Setup runtime. This is a temporary fix to trigger the runtime to recreate
	// its client now the client has been replaced with a wrapped one.
	if err := runtime.DefaultRuntime.Init(); err != nil {
		logger.Fatalf("Error configuring runtime: %v", err)
	}

	// Setup store options
	storeOpts := []store.StoreOption{}

	if len(ctx.String("namespace")) > 0 {
		storeOpts = append(storeOpts, store.Database(ctx.String("namespace")))
	}

	if len(ctx.String("name")) > 0 {
		storeOpts = append(storeOpts, store.Table(ctx.String("name")))
	}

	if err := store.DefaultStore.Init(storeOpts...); err != nil {
		logger.Fatalf("Error configuring store: %v", err)
	}

	// Setup config. Do this after auth is configured since it'll load the config
	// from the service immediately. We only do this if the action is nil, indicating
	// a service is being run
	if c.service && config.DefaultConfig == nil {
		config.DefaultConfig = configCli.NewConfig()
	} else if config.DefaultConfig == nil {
		config.DefaultConfig, _ = storeConf.NewConfig(store.DefaultStore, ctx.String("namespace"))
	}

	return nil
}

func (c *command) Run() error {
	return c.app.Run(os.Args)
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
