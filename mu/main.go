package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
	"go-micro.dev/v5/client"
	"go-micro.dev/v5/cmd"
	"go-micro.dev/v5/codec/bytes"
	"go-micro.dev/v5/errors"
	"go-micro.dev/v5/registry"
	"tailscale.com/tsnet"
)

var (
	// version is set by the release action
	// this is the default for local builds
	version = "5.0.0-dev"
)

func apiHandler(c *cli.Context) error {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// assuming we're just going to parse headers
		if r.URL.Path == "/" {
			service := r.Header.Get("Micro-Service")
			endpoint := r.Header.Get("Micro-Endpoint")
			request, _ := io.ReadAll(r.Body)
			if len(request) == 0 {
				request = []byte(`{}`)
			}

			// defaulting to json
			w.Header().Set("Content-Type", "application/json")

			if len(service) == 0 || len(endpoint) == 0 {
				err := errors.New("api.error", "missing service/endpoint", 400)
				w.Header().Set("Micro-Error", err.Error())
				http.Error(w, err.Error(), 400)
				return
			}

			req := client.NewRequest(service, endpoint, &bytes.Frame{Data: request})
			var rsp bytes.Frame
			err := client.Call(r.Context(), req, &rsp)
			if err != nil {
				gerr := errors.New("api.error", err.Error(), 500)
				w.Header().Set("Micro-Error", gerr.Error())
				http.Error(w, gerr.Error(), 500)
			}

			// write the response
			w.Write(rsp.Data)
			return
		}

		parts := strings.Split(r.URL.Path, "/")

		if len(parts) < 3 {
			return
		}

		service := parts[1]
		endpoint := parts[2]

		if len(parts) == 4 {
			endpoint = endpoint + "." + parts[3]
		} else {
			endpoint = service + "." + endpoint
		}

		request, _ := io.ReadAll(r.Body)
		if len(request) == 0 {
			request = []byte(`{}`)
		}

		req := client.NewRequest(service, endpoint, &bytes.Frame{Data: request})
		var rsp bytes.Frame
		err := client.Call(r.Context(), req, &rsp)
		if err != nil {
			gerr := errors.New("api.error", err.Error(), 500)
			w.Header().Set("Micro-Error", gerr.Error())
			http.Error(w, gerr.Error(), 500)
		}

		// write the response
		w.Write(rsp.Data)

	})

	var network string
	var key string

	if c.IsSet("network") {
		network = c.Value("network").(string)
	}

	if network == "tailscale" {
		// check for TS_AUTHKEY
		key = os.Getenv("TS_AUTHKEY")
		if len(key) == 0 {
			return fmt.Errorf("missing TS_AUTHKEY")
		}

		srv := new(tsnet.Server)
		srv.AuthKey = key
		srv.Hostname = "micro"

		ln, err := srv.Listen("tcp", ":8080")
		if err != nil {
			return err
		}

		return http.Serve(ln, h)
	}

	return http.ListenAndServe(":8080", h)
}

func genProtoHandler(c *cli.Context) error {
	cmd := exec.Command("find", ".", "-name", "*.proto", "-exec", "protoc", "--proto_path=.", "--micro_out=.", "--go_out=.", `{}`, `;`)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runHandler(c *cli.Context) error {
	dir := c.Args().Get(0)
	if len(dir) == 0 {
		dir = "."
	}
	cmd := exec.Command("go", "run", dir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	cmd.App().Commands = []*cli.Command{
		{
			Name:  "api",
			Usage: "Run the API server",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "network", Value: "", Usage: "Set the network e.g --network=tailscale requires TS_AUTHKEY"},
			},
			Action: apiHandler,
		},
		{
			Name:   "run",
			Usage:  "Run a service",
			Action: runHandler,
		},
		{
			Name:  "gen",
			Usage: "Generate various things",
			Subcommands: []*cli.Command{
				{
					Name:   "proto",
					Usage:  "Generate proto requires protoc and protoc-gen-micro",
					Action: genProtoHandler,
				},
			},
		},
		{
			Name:  "services",
			Usage: "List available services",
			Action: func(ctx *cli.Context) error {
				services, err := registry.ListServices()
				if err != nil {
					return err
				}
				for _, service := range services {
					fmt.Println(service.Name)
				}
				return nil
			},
		},
		{
			Name:  "call",
			Usage: "Call a service",
			Action: func(ctx *cli.Context) error {
				args := ctx.Args()

				if args.Len() < 2 {
					return fmt.Errorf("Usage: [service] [endpoint] [request]")
				}

				service := args.Get(0)
				endpoint := args.Get(1)
				request := `{}`

				// get the request if present
				if args.Len() == 3 {
					request = args.Get(2)
				}

				req := client.NewRequest(service, endpoint, &bytes.Frame{Data: []byte(request)})
				var rsp bytes.Frame
				err := client.Call(context.TODO(), req, &rsp)
				if err != nil {
					return err
				}

				fmt.Print(string(rsp.Data))
				return nil
			},
		},
		{
			Name:  "describe",
			Usage: "Describe a service",
			Action: func(ctx *cli.Context) error {
				args := ctx.Args()

				if args.Len() != 1 {
					return fmt.Errorf("Usage: [service]")
				}

				service := args.Get(0)
				services, err := registry.GetService(service)
				if err != nil {
					return err
				}
				if len(services) == 0 {
					return nil
				}
				b, _ := json.MarshalIndent(services[0], "", "    ")
				fmt.Println(string(b))
				return nil
			},
		},
	}

	cmd.Init(
		cmd.Name("mu"),
		cmd.Version(version),
	)
}
