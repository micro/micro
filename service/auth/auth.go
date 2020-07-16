package auth

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"
	srvAuth "github.com/micro/go-micro/v2/auth/service"
	pb "github.com/micro/go-micro/v2/auth/service/proto"
	"github.com/micro/go-micro/v2/auth/token"
	"github.com/micro/go-micro/v2/auth/token/jwt"
	"github.com/micro/go-micro/v2/cmd"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/micro/v2/client/cli/namespace"
	clitoken "github.com/micro/micro/v2/client/cli/token"
	cliutil "github.com/micro/micro/v2/client/cli/util"
	"github.com/micro/micro/v2/internal/client"
	authHandler "github.com/micro/micro/v2/service/auth/handler/auth"
	rulesHandler "github.com/micro/micro/v2/service/auth/handler/rules"
	"golang.org/x/crypto/ssh/terminal"

	// imported specifically for signup
	platform "github.com/micro/micro/v2/platform/cli"
)

var (
	// Name of the service
	Name = "go.micro.auth"
	// Address of the service
	Address = ":8010"
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
	// assuming --otp go to platform.Signup
	if isOTP := ctx.Bool("otp"); isOTP {
		platform.Signup(ctx)
		return
	}

	// otherwise assume username/password login

	// get the environment
	env := cliutil.GetEnv(ctx)
	// get the email address
	email := ctx.String("email")

	// email is blank
	if len(email) == 0 {
		fmt.Print("Enter email address: ")
		// read out the email from prompt if blank
		reader := bufio.NewReader(os.Stdin)
		email, _ = reader.ReadString('\n')
		email = strings.TrimSpace(email)
	}

	authSrv := authFromContext(ctx)
	ns, err := namespace.Get(env.Name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	password := ctx.String("password")
	if len(password) == 0 {
		pw, err := getPassword()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		password = strings.TrimSpace(pw)
		fmt.Println()
	}
	tok, err := authSrv.Token(auth.WithCredentials(email, password), auth.WithTokenIssuer(ns))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	clitoken.Save(env.Name, tok)

	fmt.Println("Successfully logged in.")
}

// taken from https://stackoverflow.com/questions/2137357/getpasswd-functionality-in-go
func getPassword() (string, error) {
	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	password := string(bytePassword)
	return strings.TrimSpace(password), nil
}

//Commands for auth
func Commands(srvOpts ...micro.Option) []*cli.Command {
	commands := []*cli.Command{
		{
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
						{
							Name:  "account",
							Usage: "Delete an auth account",
							Flags: RuleFlags,
							Action: func(ctx *cli.Context) error {
								deleteAccount(ctx)
								return nil
							},
						},
					}),
				},
			}),
		},
		{
			Name:        "login",
			Usage:       `Interactive login flow.`,
			Description: "Run 'micro login' for micro servers or 'micro login --otp' for the Micro Platform.",
			Action: func(ctx *cli.Context) error {
				login(ctx)
				return nil
			},
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "otp",
					Usage: "Login/signup with a One Time Password.",
				},
				&cli.StringFlag{
					Name:  "password",
					Usage: "Password to use for login. If not provided, will be asked for during login. Useful for automated scripts",
				},
				&cli.StringFlag{
					Name:  "email",
					Usage: "Email address to use for login",
				},
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
