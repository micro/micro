package main

import (
	"fmt"
	"net/http"

	"github.com/micro/cli"
	"github.com/micro/go-log"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
	proto "github.com/micro/go-micro/server/debug/proto"
	"golang.org/x/net/context"
)

var (
	healthAddress = "127.0.0.1:8080"
	serverAddress string
	serverName    string
)

func main() {
	app := cmd.App()
	app.Flags = append(app.Flags, cli.StringFlag{
		Name:        "health_address",
		EnvVar:      "MICRO_HEALTH_ADDRESS",
		Usage:       "Address for the health checker. 127.0.0.1:8080",
		Value:       "127.0.0.1:8080",
		Destination: &healthAddress,
	})

	app.Action = func(c *cli.Context) {
		serverName = c.String("server_name")
		serverAddress = c.String("server_address")

		if addr := c.String("health_address"); len(addr) > 0 {
			healthAddress = addr
		}

		if len(healthAddress) == 0 {
			log.Fatal("health address not set")
		}
		if len(serverName) == 0 {
			log.Fatal("server name not set")
		}
		if len(serverAddress) == 0 {
			log.Fatal("server address not set")
		}
	}

	cmd.Init()

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

	if err := http.ListenAndServe(healthAddress, nil); err != nil {
		log.Fatal(err)
	}
}
