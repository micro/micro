package auth

import (
	"fmt"
	"os"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"
	srvAuth "github.com/micro/go-micro/v2/auth/service"
	pb "github.com/micro/go-micro/v2/auth/service/proto"
	"github.com/micro/go-micro/v2/auth/token"
	"github.com/micro/go-micro/v2/auth/token/jwt"
	"github.com/micro/go-micro/v2/config/cmd"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/store"
	"github.com/micro/go-micro/v2/util/config"
	"github.com/micro/micro/v2/auth/api"
	accountsHandler "github.com/micro/micro/v2/auth/handler/accounts"
	authHandler "github.com/micro/micro/v2/auth/handler/auth"
	rulesHandler "github.com/micro/micro/v2/auth/handler/rules"
	"github.com/micro/micro/v2/auth/web"
)

var (
	// Name of the service
	Name = "go.micro.auth"
	// Address of the service
	Address = ":8010"
	// ServiceFlags are provided to commands which run micro services
	ServiceFlags = []cli.Flag{
		&cli.StringFlag{
			Name:    "address",
			Usage:   "Set the auth http address e.g 0.0.0.0:8010",
			EnvVars: []string{"MICRO_SERVER_ADDRESS"},
		},
		&cli.StringFlag{
			Name:    "auth_provider",
			EnvVars: []string{"MICRO_AUTH_PROVIDER"},
			Usage:   "Auth provider enables account generation",
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
	}
	// RuleFlags are provided to commands which create or delete rules
	RuleFlags = []cli.Flag{
		&cli.StringFlag{
			Name:     "role",
			Usage:    "The role to amend, e.g. 'user' or '*' to represent all",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "resource_type",
			Usage:    "The type of resouce to amend, e.g. service",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "resource_name",
			Usage:    "The name of the resouce to amend, e.g. go.micro.auth",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "resource_endpoint",
			Usage: "The endpoint of the resouce to amend, e.g. Auth.Generate",
			Value: "*",
		},
		&cli.StringFlag{
			Name:     "access",
			Usage:    "The access level, must be granted or denied",
			Required: true,
		},
	}
	// AccountFlags are provided to the create account command
	AccountFlags = []cli.Flag{
		&cli.StringFlag{
			Name:     "id",
			Usage:    "The account id",
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:     "roles",
			Usage:    "Comma seperated list of roles to give the account",
			Required: true,
		},
	}
	// PlatformFlag connects via proxy
	PlatformFlag = &cli.BoolFlag{
		Name:  "platform",
		Usage: "Connect to the platform",
		Value: false,
	}
)

// run the auth service
func run(ctx *cli.Context, srvOpts ...micro.Option) {
	log.Init(log.WithFields(map[string]interface{}{"service": "auth"}))

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}
	if len(Address) > 0 {
		srvOpts = append(srvOpts, micro.Address(Address))
	}

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	// set store namespace
	store.DefaultStore.Init(store.Namespace(Name))

	// setup the handlers
	authH := &authHandler.Auth{}
	ruleH := &rulesHandler.Rules{}
	accountH := &accountsHandler.Accounts{}

	// setup the auth handler to use JWTs
	pubKey := ctx.String("auth_public_key")
	privKey := ctx.String("auth_private_key")
	if len(pubKey) > 0 || len(privKey) > 0 {
		authH.TokenProvider = jwt.NewTokenProvider(
			token.WithPublicKey(pubKey),
			token.WithPrivateKey(privKey),
		)
	}

	// set the handlers store
	authH.Init(auth.Store(store.DefaultStore))
	ruleH.Init(auth.Store(store.DefaultStore))
	accountH.Init(auth.Store(store.DefaultStore))

	// setup service
	srvOpts = append(srvOpts, micro.Name(Name))
	service := micro.NewService(srvOpts...)

	// register handlers
	pb.RegisterAuthHandler(service.Server(), authH)
	pb.RegisterRulesHandler(service.Server(), ruleH)
	pb.RegisterAccountsHandler(service.Server(), accountH)

	// run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func authFromContext(ctx *cli.Context) auth.Auth {
	if ctx.Bool("platform") {
		os.Setenv("MICRO_PROXY", "service")
		os.Setenv("MICRO_PROXY_ADDRESS", "proxy.micro.mu:443")
		return srvAuth.NewAuth()
	}

	return *cmd.DefaultCmd.Options().Auth
}

// login using a token
func login(ctx *cli.Context) {
	if ctx.Args().Len() != 1 {
		fmt.Println("Usage: `micro login [token]`")
		os.Exit(1)
	}
	token := ctx.Args().First()

	// Execute the request
	acc, err := authFromContext(ctx).Inspect(token)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else if acc == nil {
		fmt.Printf("[%v] did not generate an account\n", authFromContext(ctx).String())
		os.Exit(1)
	}

	// Store the token in micro config
	if err := config.Set("token", token); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Inform the user
	fmt.Println("You have been logged in")
}

func Commands(srvOpts ...micro.Option) []*cli.Command {
	commands := []*cli.Command{
		&cli.Command{
			Name:  "auth",
			Usage: "Run the auth service",
			Action: func(ctx *cli.Context) error {
				run(ctx)
				return nil
			},
			Subcommands: append([]*cli.Command{
				{
					Name:        "api",
					Usage:       "Run the auth api",
					Description: "Run the auth api",
					Flags:       ServiceFlags,
					Action: func(ctx *cli.Context) error {
						api.Run(ctx, srvOpts...)
						return nil
					},
				},
				{
					Name:        "web",
					Usage:       "Run the auth web",
					Description: "Run the auth web",
					Flags:       append(ServiceFlags, PlatformFlag),
					Action: func(ctx *cli.Context) error {
						web.Run(ctx, srvOpts...)
						return nil
					},
				},
				&cli.Command{
					Name:  "list",
					Usage: "List auth resources",
					Subcommands: append([]*cli.Command{
						{
							Name:  "rules",
							Usage: "List auth rules",
							Flags: []cli.Flag{PlatformFlag},
							Action: func(ctx *cli.Context) error {
								listRules(ctx)
								return nil
							},
						},
						{
							Name:  "accounts",
							Usage: "List auth accounts",
							Flags: []cli.Flag{PlatformFlag},
							Action: func(ctx *cli.Context) error {
								listAccounts(ctx)
								return nil
							},
						},
					}),
				},
				&cli.Command{
					Name:  "create",
					Usage: "Create an auth resource",
					Subcommands: append([]*cli.Command{
						{
							Name:  "rule",
							Usage: "Create an auth rule",
							Flags: append(RuleFlags, PlatformFlag),
							Action: func(ctx *cli.Context) error {
								createRule(ctx)
								return nil
							},
						},
						{
							Name:  "account",
							Usage: "Create an auth account",
							Flags: append(AccountFlags, PlatformFlag),
							Action: func(ctx *cli.Context) error {
								createAccount(ctx)
								return nil
							},
						},
					}),
				},
				&cli.Command{
					Name:  "delete",
					Usage: "Delete a auth resource",
					Subcommands: append([]*cli.Command{
						{
							Name:  "rule",
							Usage: "Delete an auth rule",
							Flags: append(RuleFlags, PlatformFlag),
							Action: func(ctx *cli.Context) error {
								deleteRule(ctx)
								return nil
							},
						},
					}),
				},
			}),
		},
		&cli.Command{
			Name:  "login",
			Usage: "Login using a token",
			Action: func(ctx *cli.Context) error {
				login(ctx)
				return nil
			},
			Flags: []cli.Flag{PlatformFlag},
		},
	}

	for _, c := range commands {
		for _, p := range Plugins() {
			if cmds := p.Commands(); len(cmds) > 0 {
				c.Subcommands = append(c.Subcommands, cmds...)
			}

			if flags := p.Flags(); len(flags) > 0 {
				c.Flags = append(c.Flags, flags...)
			}
		}
	}

	return commands
}
