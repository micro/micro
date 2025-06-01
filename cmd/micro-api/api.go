package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/urfave/cli/v2"
	"go-micro.dev/v5/client"
	"go-micro.dev/v5/cmd"
	"go-micro.dev/v5/codec/bytes"
	"go-micro.dev/v5/errors"
	"go-micro.dev/v5/registry"
	"tailscale.com/tsnet"
)

func normalize(v string) string {
	return strings.Title(v)
}

func mcpServer() *server.SSEServer {
	s := server.NewMCPServer(
		"micro",
		"1.0.0",
	)

	// Add 'call' tool
	s.AddTool(mcp.NewTool("call",
		mcp.WithDescription("Call a service"),
		mcp.WithString("service", mcp.Required(), mcp.Description("Name of the service e.g helloworld")),
		mcp.WithString("endpoint", mcp.Required(), mcp.Description("Name of the endpoint e.g Say.Hello")),
		mcp.WithString("request", mcp.Required(), mcp.Description("JSON request body")),
	), func(ctx context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		service, ok := r.Params.Arguments["service"].(string)
		if !ok {
			return nil, fmt.Errorf("service must be a string")
		}
		endpoint, ok := r.Params.Arguments["endpoint"].(string)
		if !ok {
			return nil, fmt.Errorf("endpoint must be a string")
		}
		request, ok := r.Params.Arguments["request"].(string)
		if !ok {
			return nil, fmt.Errorf("request must be a string")
		}
		jreq := json.RawMessage(request)
		req := client.NewRequest(service, endpoint, &bytes.Frame{Data: jreq})
		var rsp bytes.Frame
		if err := client.Call(ctx, req, &rsp); err != nil {
			return nil, fmt.Errorf("Call error: %v", err)
		}
		return mcp.NewToolResultText(string(rsp.Data)), nil
	})

	// Add 'describe' tool
	s.AddTool(mcp.NewTool("describe",
		mcp.WithDescription("Describe a service"),
		mcp.WithString("service", mcp.Required(), mcp.Description("Name of the service e.g helloworld")),
	), func(ctx context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		service, ok := r.Params.Arguments["service"].(string)
		if !ok {
			return nil, fmt.Errorf("service must be a string")
		}
		services, err := registry.GetService(service)
		if err != nil {
			return nil, err
		}
		if len(services) == 0 {
			return nil, fmt.Errorf("service not found")
		}
		b, _ := json.MarshalIndent(services[0], "", "    ")
		return mcp.NewToolResultText(string(b)), nil
	})

	// Add 'services' tool
	s.AddTool(mcp.NewTool("services", mcp.WithDescription("List services")), func(ctx context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		services, err := registry.ListServices()
		if err != nil {
			return nil, err
		}
		var list []string
		for _, service := range services {
			list = append(list, service.Name)
		}
		b, _ := json.MarshalIndent(list, "", "    ")
		return mcp.NewToolResultText(string(b)), nil
	})

	return server.NewSSEServer(s)
}

func init() {
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

		if r.URL.Path == "/mcp" {
			mcpServer := mcpServer()
			mcpServer.ServeHTTP(w, r)
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

	cmd.Register(&cli.Command{
		Name:  "api",
		Usage: "Run the micro api on port :8080",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "network", Value: "", Usage: "Set the network e.g --network=tailscale requires TS_AUTHKEY"},
		},
		Action: func(c *cli.Context) error {
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
		},
	})
}
