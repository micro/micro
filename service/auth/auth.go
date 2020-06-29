package auth

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"
	srvAuth "github.com/micro/go-micro/v2/auth/service"
	pb "github.com/micro/go-micro/v2/auth/service/proto"
	"github.com/micro/go-micro/v2/auth/token"
	"github.com/micro/go-micro/v2/auth/token/jwt"
	"github.com/micro/go-micro/v2/config/cmd"
	"github.com/micro/go-micro/v2/errors"
	log "github.com/micro/go-micro/v2/logger"
	cliutil "github.com/micro/micro/v2/client/cli/util"
	"github.com/micro/micro/v2/internal/client"
	"github.com/micro/micro/v2/internal/config"
	"github.com/micro/micro/v2/internal/helper"
	"github.com/micro/micro/v2/service/auth/api"
	authHandler "github.com/micro/micro/v2/service/auth/handler/auth"
	rulesHandler "github.com/micro/micro/v2/service/auth/handler/rules"
	signupproto "github.com/micro/services/signup/proto/signup"
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
			Name:  "scope",
			Usage: "The scope to amend, e.g. 'user' or '*', leave blank to make public",
		},
		&cli.StringFlag{
			Name:  "resource",
			Usage: "The resource to amend in the format type:name:endpoint, e.g. service:go.micro.auth:*",
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
			Name:  "secret",
			Usage: "The account secret (password)",
		},
		&cli.StringSliceFlag{
			Name:  "scopes",
			Usage: "Comma seperated list of scopes to give the account",
		},
	}
)

// run the auth service
func Run(ctx *cli.Context, srvOpts ...micro.Option) {
	log.Init(log.WithFields(map[string]interface{}{"service": "auth"}))

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
	ruleH := &rulesHandler.Rules{}
	authH := &authHandler.Auth{}

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

	// setup service
	srvOpts = append(srvOpts, micro.Name(Name))
	service := micro.NewService(srvOpts...)

	// register handlers
	pb.RegisterAuthHandler(service.Server(), authH)
	pb.RegisterRulesHandler(service.Server(), ruleH)
	pb.RegisterAccountsHandler(service.Server(), authH)

	// run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func authFromContext(ctx *cli.Context) auth.Auth {
	if cliutil.IsLocal(ctx) {
		return *cmd.DefaultCmd.Options().Auth
	}
	return srvAuth.NewAuth(
		auth.WithClient(client.New(ctx)),
	)
}

// login flow.
// For documentation of the flow please refer to https://github.com/micro/development/pull/223
func login(ctx *cli.Context) {
	env := cliutil.GetEnv(ctx)
	if tok := ctx.String("token"); len(tok) > 0 {
		_, err := authFromContext(ctx).Inspect(tok)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err := config.Set(tok, "micro", "auth", env.Name, "token"); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("You have been logged in")
		return
	}

	email := ""
	if ctx.Args().Len() > 0 {
		email = ctx.Args().First()
	}
	reader := bufio.NewReader(os.Stdin)
	if len(email) == 0 {
		fmt.Print("Please enter your email address: ")
		email, _ = reader.ReadString('\n')
		email = strings.TrimSpace(email)
	}
	signupService := signupproto.NewSignupService("go.micro.service.signup", client.New(ctx))

	_, err := signupService.SendVerificationEmail(context.TODO(), &signupproto.SendVerificationEmailRequest{
		Email: email,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Print("We have sent you an email with a one time password. Please paste it here: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	rsp, err := signupService.Verify(context.TODO(), &signupproto.VerifyRequest{
		Email: email,
		Token: password,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Already registered users can just get logged in.
	tok := rsp.AuthToken
	if rsp.AuthToken != nil {
		if err := config.Set(tok.AccessToken, "micro", "auth", env.Name, "token"); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// Store the refresh token in micro config
		if err := config.Set(tok.RefreshToken, "micro", "auth", env.Name, "refresh-token"); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Successfully logged in.")
		return
	}

	fmt.Print("Please go to https://m3o.com/subscribe.html and paste the acquired payment method id here: ")
	paymentMethodID, _ := reader.ReadString('\n')
	paymentMethodID = strings.TrimSpace(paymentMethodID)

	signupRsp, err := signupService.CompleteSignup(context.TODO(), &signupproto.CompleteSignupRequest{
		Email:           email,
		Token:           password,
		PaymentMethodID: paymentMethodID,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tok = signupRsp.AuthToken
	if err := config.Set(tok.AccessToken, "micro", "auth", env.Name, "token"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Store the refresh token in micro config
	if err := config.Set(tok.RefreshToken, "micro", "auth", env.Name, "refresh-token"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Successfully logged in.")
	// @todo save the namespace from the last call and use that.
}

// whoami returns info about the logged in user
func whoami(ctx *cli.Context) {
	// Get the token from micro config
	env, _ := config.Get("env")
	tok, err := config.Get("micro", "auth", env, "token")
	if err != nil {
		fmt.Println("You are not logged in")
		os.Exit(1)
	}

	// Inspect the token
	acc, err := authFromContext(ctx).Inspect(tok)
	if verr, ok := err.(*errors.Error); ok {
		fmt.Println("Error: " + verr.Detail)
		return
	} else if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("ID: %v; Scopes: %v\n", acc.ID, strings.Join(acc.Scopes, ", "))
}

//Commands for auth
func Commands(srvOpts ...micro.Option) []*cli.Command {
	commands := []*cli.Command{
		{
			Name:  "auth",
			Usage: "Manage authentication related resources",
			Action: func(ctx *cli.Context) error {
				if err := helper.UnexpectedSubcommand(ctx); err != nil {
					return err
				}
				Run(ctx)
				return nil
			},
			Subcommands: append([]*cli.Command{
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
			}),
		},
		{
			Name:  "login",
			Usage: "Interactive login flow. Just type `micro login` or `micro login [email address]`",
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
