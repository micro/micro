package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	"go-micro.dev/v5/client"
	"go-micro.dev/v5/cmd"
	"go-micro.dev/v5/codec/bytes"
	"go-micro.dev/v5/errors"
	"tailscale.com/tsnet"
)

func normalize(v string) string {
	return strings.Title(v)
}

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
			endpoint = normalize(endpoint) + "." + normalize(parts[3])
		} else {
			endpoint = normalize(service) + "." + normalize(endpoint)
		}

		request, _ := io.ReadAll(r.Body)
		if len(request) == 0 {
			request = []byte(`{}`)

			if r.Method == "GET" {
				req := map[string]interface{}{}
				r.ParseForm()
				for k, v := range r.Form {
					req[k] = strings.Join(v, ",")
				}
				if len(req) > 0 {
					request, _ = json.Marshal(req)
				}
			}
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
