// Package health is a healthchecking sidecar
package health

import (
	"fmt"
	"net/http"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	proto "github.com/micro/go-micro/debug/service/proto"
	"github.com/micro/go-micro/util/log"
	mcli "github.com/micro/micro/cli"
	qcli "github.com/micro/micro/internal/command/cli"
	"golang.org/x/net/context"
)

var (
	healthAddress = ":8088"
	serverAddress string
	serverName    string
)

func Run(ctx *cli.Context) {
	log.Name("health")

	// just check service health
	if len(ctx.Args()) > 0 {
		mcli.Print(qcli.QueryHealth)(ctx)
		return
	}

	serverName = ctx.String("check_service")
	serverAddress = ctx.String("check_address")

	if addr := ctx.String("address"); len(addr) > 0 {
		healthAddress = addr
	}

	if len(healthAddress) == 0 {
		log.Fatal("health address not set")
	}
	if len(serverName) == 0 {
		log.Fatal("service name not set")
	}
	if len(serverAddress) == 0 {
		log.Fatal("service address not set")
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		req := client.NewRequest(serverName, "Debug.Health", &proto.HealthRequest{})
		rsp := &proto.HealthResponse{}

		err := client.Call(context.TODO(), req, rsp, client.WithAddress(serverAddress))
		if err != nil || rsp.Status != "ok" {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "NOT_HEALTHY")
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	log.Infof("Health check running at %s/health", healthAddress)
	log.Infof("Health check defined for %s at %s", serverName, serverAddress)

	if err := http.ListenAndServe(healthAddress, nil); err != nil {
		log.Fatal(err)
	}
}

func Commands(options ...micro.Option) []cli.Command {
	command := cli.Command{
		Name:  "health",
		Usage: "Check the health of a service",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "address",
				Usage:  "Set the address exposed for the http server e.g :8088",
				EnvVar: "MICRO_HEALTH_ADDRESS",
			},
			cli.StringFlag{
				Name:   "check_service",
				Usage:  "Name of the service to query",
				EnvVar: "MICRO_HEALTH_CHECK_SERVICE",
			},
			cli.StringFlag{
				Name:   "check_address",
				Usage:  "Set the service address to query",
				EnvVar: "MICRO_HEALTH_CHECK_ADDRESS",
			},
		},
		Action: func(ctx *cli.Context) {
			Run(ctx)
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
