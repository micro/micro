package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"go-micro.dev/v5/client"
	"go-micro.dev/v5/codec/bytes"
	"go-micro.dev/v5/errors"
	"go-micro.dev/v5/registry"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"context"
)

var mcpSrv = mcpServer()

func APIHandler() http.Handler {
   h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	  if r.URL.Path == "/mcp" {
		  mcpSrv.ServeHTTP(w, r)
		  return
	  }

		if strings.HasPrefix(r.URL.Path, "/api/") {
			// Remove /api prefix
			path := strings.TrimPrefix(r.URL.Path, "/api")
			if path == "" || path == "/" {
				service := r.Header.Get("Micro-Service")
				endpoint := r.Header.Get("Micro-Endpoint")
				request, _ := io.ReadAll(r.Body)
				if len(request) == 0 {
					request = []byte(`{}`)
				}
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
				w.Write(rsp.Data)
				return
			}

			parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
			if len(parts) < 2 {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			service := parts[0]
			endpoint := parts[1]
			if len(parts) == 3 {
				endpoint = normalize(endpoint) + "." + normalize(parts[2])
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
			w.Write(rsp.Data)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	})
	return h
}

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
