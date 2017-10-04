// Package run is a micro service runtime
package run

import (
	"fmt"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"

	"github.com/micro/go-run"
	"github.com/micro/go-run/runtime/go"
	proto "github.com/micro/micro/run/proto"

	"golang.org/x/net/context"
)

var (
	Name = "go.micro.run"
)

func manage(r run.Runtime, url string, re, u bool) error {
	// get the source
	log.Logf("fetching %s\n", url)
	src, err := r.Fetch(url, run.Update(u))
	if err != nil {
		return err
	}

	// build the binary
	log.Logf("building %s\n", url)
	bin, err := r.Build(src)
	if err != nil {
		return err
	}

	for {
		// execute the binary
		log.Logf("executing %s\n", url)
		proc, err := r.Exec(bin)
		if err != nil {
			return err
		}

		// wait till exit
		log.Logf("running %s\n", url)

		// bail if not restarting
		if !re {
			return r.Wait(proc)
		}

		// log error since we manage the cycle
		if err := r.Wait(proc); err != nil {
			log.Logf("exited with err %v\n", err)
		}

		// cruft log
		log.Logf("restarting %s\n", url)
	}
}

func runc(ctx *cli.Context) {
	if len(ctx.GlobalString("server_name")) > 0 {
		Name = ctx.GlobalString("server_name")
	}

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	// Initialise Server
	service := micro.NewService(
		micro.Name(Name),
		micro.RegisterTTL(
			time.Duration(ctx.GlobalInt("register_ttl"))*time.Second,
		),
		micro.RegisterInterval(
			time.Duration(ctx.GlobalInt("register_interval"))*time.Second,
		),
	)

	r := gorun.NewRuntime()
	m := newManager(r)
	prun := micro.NewPublisher(Name+".run", service.Client())
	pstop := micro.NewPublisher(Name+".stop", service.Client())

	// micro run github.com/my/service
	// args: github.com/my/service
	if len(ctx.Args()) > 0 {
		// Notes to self:
		// 1. Done: look for flag to loop run cycle infinitely
		// 2. Done: look for flag to defer execution to go.micro.run service
		// 3. look for daemonize flag
		// 4. Done: look for flag to defer update

		ki := ctx.Bool("k")
		st := ctx.Bool("s")
		re := ctx.Bool("r")
		up := ctx.Bool("u")
		xe := ctx.Bool("x")

		// kill a service
		if ki {
			// call runtime manager service
			cl := proto.NewServiceClient(Name, client.DefaultClient)
			_, err := cl.Stop(context.TODO(), &proto.StopRequest{
				Url: ctx.Args().First(),
			})
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		// get status
		if st {
			// call runtime manager service
			cl := proto.NewServiceClient(Name, client.DefaultClient)
			rsp, err := cl.Status(context.TODO(), &proto.StatusRequest{
				Url: ctx.Args().First(),
			})
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(rsp.Info)
			return
		}

		// defer to service
		if xe {
			// call runtime manager service
			cl := proto.NewServiceClient(Name, client.DefaultClient)
			_, err := cl.Run(context.TODO(), &proto.RunRequest{
				Url:     ctx.Args().First(),
				Restart: re,
				Update:  up,
			})
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		// TODO: should return stream?
		go m.Run(ctx.Args().First(), re, up)

		// publish run event
		go prun.Publish(context.TODO(), &proto.RunRequest{
			Url:     ctx.Args().First(),
			Restart: re,
			Update:  up,
		})
	}

	// handlers
	proto.RegisterRuntimeHandler(service.Server(), &runtimeHandler{r})
	proto.RegisterServiceHandler(service.Server(), &serviceHandler{m, prun, pstop})

	// subscribers
	sb := &subHandler{m}
	// run commands
	micro.RegisterSubscriber(Name+".run", service.Server(), sb.Run)
	// stop commands
	micro.RegisterSubscriber(Name+".stop", service.Server(), sb.Stop)

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func Commands() []cli.Command {
	command := cli.Command{
		Name:   "run",
		Usage:  "Run the micro runtime",
		Action: runc,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "k",
				Usage: "Kill service",
			},
			cli.BoolFlag{
				Name:  "r",
				Usage: "Restart if dies. Default: false",
			},
			cli.BoolFlag{
				Name:  "u",
				Usage: "Update the source. Default: false",
			},
			cli.BoolFlag{
				Name:  "x",
				Usage: "Defer run to service. Default: false",
			},
			cli.BoolFlag{
				Name:  "s",
				Usage: "Get service status",
			},
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
