package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"os"

	"tailscale.com/tsnet"
	"github.com/urfave/cli/v2"
	"go-micro.dev/v5/cmd"
	"go-micro.dev/v5/client"
	"go-micro.dev/v5/errors"
	"go-micro.dev/v5/codec/bytes"
)

func main() {
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


	app := cmd.App()

	app.Flags = []cli.Flag{
		&cli.StringFlag{Name: "network", Value: "", Usage: "Set the network e.g --network=tailscale requires TS_AUTHKEY"},
	}
	app.Action = func(c *cli.Context) error {
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

	cmd.Init(
		cmd.Name("micro-api"),
	)
}
