package service

import (
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v3/cmd"
	incmd "github.com/micro/micro/v3/internal/cmd"
)

var (
	defaultCmd   cmd.Cmd = newCmd()
	defaultFlags         = []cli.Flag{
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
	cmd.app.Flags = defaultFlags
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
	if n := ctx.String("service_name"); len(n) > 0 {
		Name(n)(&defaultService.opts)
	}
	if v := ctx.String("service_version"); len(v) > 0 {
		Version(v)(&defaultService.opts)
	}

	for _, fn := range incmd.InitFuncs {
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
