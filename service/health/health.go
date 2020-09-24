// Package health is a healthchecking sidecar
package health

import (
	"fmt"
	"net/http"

	goclient "github.com/micro/go-micro/v3/client"
	"github.com/micro/micro/v3/client/cli/util"
	qcli "github.com/micro/micro/v3/internal/command"
	proto "github.com/micro/micro/v3/proto/debug"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/logger"
	"github.com/urfave/cli/v2"
	"golang.org/x/net/context"
)

const healthAddress = ":8088"

// Flags specific to the health service
var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:    "check_service",
		Usage:   "Name of the service to query",
		EnvVars: []string{"MICRO_HEALTH_CHECK_SERVICE"},
	},
	&cli.StringFlag{
		Name:    "check_address",
		Usage:   "Set the service address to query",
		EnvVars: []string{"MICRO_HEALTH_CHECK_ADDRESS"},
	},
}

// Run micro health
func Run(ctx *cli.Context) error {
	// just check service health
	if ctx.Args().Len() > 0 {
		util.Print(qcli.QueryHealth)(ctx)
		return nil
	}

	serverName := ctx.String("check_service")
	serverAddress := ctx.String("check_address")

	if len(serverName) == 0 {
		logger.Fatal("service name not set")
	}
	if len(serverAddress) == 0 {
		logger.Fatal("service address not set")
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		req := client.NewRequest(serverName, "Debug.Health", &proto.HealthRequest{})
		rsp := &proto.HealthResponse{}

		err := client.DefaultClient.Call(context.TODO(), req, rsp, goclient.WithAddress(serverAddress))
		if err != nil || rsp.Status != "ok" {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "NOT_HEALTHY")
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	logger.Infof("Health check running at %s/health", healthAddress)
	logger.Infof("Health check defined for %s at %s", serverName, serverAddress)

	if err := http.ListenAndServe(healthAddress, nil); err != nil {
		logger.Fatal(err)
	}
	return nil
}
