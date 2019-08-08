package tunnel

import (
	"sync"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/util/log"
)

var (
	// Name of the router microservice
	Name = "go.micro.tunnel"
	// Address is the tunnel microservice bind address
	Address = ":8084"
	// Tunnel is the tunnel bind address
	Tunnel = ":9096"
	// Network is the network id
	Network = "local"
)

// run runs the micro server
func run(ctx *cli.Context, srvOpts ...micro.Option) {
	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	if len(ctx.GlobalString("server_name")) > 0 {
		Name = ctx.GlobalString("server_name")
	}
	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}
	if len(ctx.String("network_address")) > 0 {
		Network = ctx.String("network")
	}
	if len(ctx.String("tunnel_address")) > 0 {
		Tunnel = ctx.String("tunnel")
	}

	// Initialise service
	service := micro.NewService(
		micro.Name(Name),
		micro.Address(Address),
		micro.RegisterTTL(time.Duration(ctx.GlobalInt("register_ttl"))*time.Second),
		micro.RegisterInterval(time.Duration(ctx.GlobalInt("register_interval"))*time.Second),
	)

	var wg sync.WaitGroup

	errChan := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		errChan <- service.Run()
	}()

	log.Log("[tunnel] attempting to stop the tunnel")

	wg.Wait()

	log.Logf("[tunnel] successfully stopped")

}

func Commands(options ...micro.Option) []cli.Command {
	command := cli.Command{
		Name:  "tunnel",
		Usage: "Run the micro network tunnel",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "tunnel_address",
				Usage:  "Set the micro tunnel address :9096",
				EnvVar: "MICRO_TUNNEL_ADDRESS",
			},
			cli.StringFlag{
				Name:   "network_address",
				Usage:  "Set the micro network address: local",
				EnvVar: "MICRO_NETWORK_ADDRESS",
			},
		},
		Action: func(ctx *cli.Context) {
			run(ctx, options...)
		},
	}

	for _, p := range Plugins() {
		if cmds := p.Commands(); len(cmds) > 0 {
			command.Subcommands = append(command.Subcommands, cmds...)
		}

		if flags := p.Flags(); len(flags) > 0 {
			command.Flags = append(command.Flags, flags...)
		}
	}

	return []cli.Command{command}
}
