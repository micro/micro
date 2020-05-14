package auth

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"
	srvAuth "github.com/micro/go-micro/v2/auth/service"
	pb "github.com/micro/go-micro/v2/auth/service/proto"
	"github.com/micro/go-micro/v2/auth/token"
	"github.com/micro/go-micro/v2/auth/token/jwt"
	"github.com/micro/go-micro/v2/config/cmd"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/util/config"
	"github.com/micro/micro/v2/auth/api"
	accountsHandler "github.com/micro/micro/v2/auth/handler/accounts"
	authHandler "github.com/micro/micro/v2/auth/handler/auth"
	rulesHandler "github.com/micro/micro/v2/auth/handler/rules"
	cliutil "github.com/micro/micro/v2/cli/util"
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
			Name:     "resource",
			Usage:    "The resource to amend in the format namespace:type:name:endpoint, e.g. micro:service:go.micro.auth:*",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "access",
			Usage: "The access level, must be granted or denied",
			Value: "granted",
		},
		&cli.IntFlag{
			Name:  "priority",
			Usage: "The priority level, default is 0, the greater the number the higher the priority",
			Value: 0,
		},
	}
	// AccountFlags are provided to the create account command
	AccountFlags = []cli.Flag{
		&cli.StringFlag{
			Name:     "id",
			Usage:    "The account id",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "secret",
			Usage:    "The account secret (password)",
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:  "roles",
			Usage: "Comma seperated list of roles to give the account",
		},
	}
)

// run the auth service
func Run(ctx *cli.Context, srvOpts ...micro.Option) {
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

	st := *cmd.DefaultCmd.Options().Store

	// set the handlers store
	authH.Init(auth.Store(st))
	ruleH.Init(auth.Store(st))
	accountH.Init(auth.Store(st))

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
	if cliutil.IsLocal() {
		return *cmd.DefaultCmd.Options().Auth
	}
	return srvAuth.NewAuth()
}

// login using a token
func login(ctx *cli.Context) {
	// check for the token flag
	if tok := ctx.String("token"); len(tok) > 0 {
		_, err := authFromContext(ctx).Inspect(tok)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := config.Set(tok, "micro", "auth", "token"); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("You have been logged in")
		return
	}

	if ctx.Args().Len() != 2 {
		fmt.Println("Usage: `micro login {id} {secret} OR micro login --token {token}`")
		os.Exit(1)
	}
	id := ctx.Args().Get(0)
	secret := ctx.Args().Get(1)

	// Execute the request
	tok, err := authFromContext(ctx).Token(auth.WithCredentials(id, secret), auth.WithExpiry(time.Hour*24))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	envName := cliutil.GetEnv().Name
	// Store the access token in micro config
	if err := config.Set(tok.AccessToken, "micro", "envs", envName, "auth", "token"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Store the refresh token in micro config
	if err := config.Set(tok.RefreshToken, "micro", "envs", envName, "auth", "refresh-token"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Inform the user
	fmt.Println("You have been logged in")
}

// whoami returns info about the logged in user
func whoami(ctx *cli.Context) {
	// Get the token from micro config
	envName := cliutil.GetEnv().Name
	tok, err := config.Get("micro", "envs", envName, "auth", "token")
	if err != nil {
		fmt.Println("You are not logged in")
		os.Exit(1)
	}

	// Inspect the token
	acc, err := authFromContext(ctx).Inspect(tok)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("ID: %v\n", acc.ID)
	fmt.Printf("Roles: %v\n", strings.Join(acc.Roles, ", "))
}

//Commands for auth
func Commands(srvOpts ...micro.Option) []*cli.Command {
	commands := []*cli.Command{
		{
			Name:  "auth",
			Usage: "Run the auth service",
			Action: func(ctx *cli.Context) error {
				Run(ctx)
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
					Name:  "list",
					Usage: "List auth resources",
					Subcommands: append([]*cli.Command{
						{
							Name:  "rules",
							Usage: "List auth rules",
							Action: func(ctx *cli.Context) error {
								listRules(ctx)
								return nil
							},
						},
						{
							Name:  "accounts",
							Usage: "List auth accounts",
							Action: func(ctx *cli.Context) error {
								listAccounts(ctx)
								return nil
							},
						},
					}),
				},
				{
					Name:  "create",
					Usage: "Create an auth resource",
					Subcommands: append([]*cli.Command{
						{
							Name:  "rule",
							Usage: "Create an auth rule",
							Flags: append(RuleFlags),
							Action: func(ctx *cli.Context) error {
								createRule(ctx)
								return nil
							},
						},
						{
							Name:  "account",
							Usage: "Create an auth account",
							Flags: append(AccountFlags),
							Action: func(ctx *cli.Context) error {
								createAccount(ctx)
								return nil
							},
						},
					}),
				},
				{
					Name:  "delete",
					Usage: "Delete a auth resource",
					Subcommands: append([]*cli.Command{
						{
							Name:  "rule",
							Usage: "Delete an auth rule",
							Flags: RuleFlags,
							Action: func(ctx *cli.Context) error {
								deleteRule(ctx)
								return nil
							},
						},
					}),
				},
			}),
		},
		{
			Name:  "login",
			Usage: "Login using a token",
			Action: func(ctx *cli.Context) error {
				login(ctx)
				return nil
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "token",
					Usage: "The token to set",
				},
			},
		},
		{
			Name:  "whoami",
			Usage: "Account information",
			Action: func(ctx *cli.Context) error {
				whoami(ctx)
				return nil
			},
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
