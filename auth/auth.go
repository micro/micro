package auth

import (
	"fmt"
	"os"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"
	srvAuth "github.com/micro/go-micro/v2/auth/service"
	pb "github.com/micro/go-micro/v2/auth/service/proto"
	"github.com/micro/go-micro/v2/config/cmd"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/util/config"
	"github.com/micro/micro/v2/auth/api"
	"github.com/micro/micro/v2/auth/handler"
)

var (
	// Name of the service
	Name = "go.micro.auth"
	// Address of the service
	Address = ":8010"
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

	// setup service
	srvOpts = append(srvOpts, micro.Name(Name))
	service := micro.NewService(srvOpts...)

	// run service
	pb.RegisterAuthHandler(service.Server(), handler.New())
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
	acc, err := authFromContext(ctx).Verify(token)
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
					Action: func(ctx *cli.Context) error {
						api.Run(ctx, srvOpts...)
						return nil
					},
				},
			}),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "address",
					Usage:   "Set the auth http address e.g 0.0.0.0:8010",
					EnvVars: []string{"MICRO_SERVER_ADDRESS"},
				},
			},
		},
		&cli.Command{
			Name:  "login",
			Usage: "Login using a token",
			Action: func(ctx *cli.Context) error {
				login(ctx)
				return nil
			},
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "platform",
					Usage: "Connect to the platform",
					Value: false,
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
