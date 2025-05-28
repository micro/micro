package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/urfave/cli/v2"
	"go-micro.dev/v5/client"
	"go-micro.dev/v5/cmd"
	"go-micro.dev/v5/codec/bytes"
	"go-micro.dev/v5/registry"
)

func handler(ctx context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	service, ok := r.Params.Arguments["service"].(string)
	if !ok {
		return nil, errors.New("service must be a string")
	}
	endpoint, ok := r.Params.Arguments["endpoint"].(string)
	if !ok {
		return nil, errors.New("endpoint must be a string")
	}
	request, ok := r.Params.Arguments["request"].(string)
	if !ok {
		return nil, errors.New("request must be a string")
	}

	jreq := json.RawMessage(request)

	// make the request
	req := client.NewRequest(service, endpoint, &bytes.Frame{Data: jreq})

	var rsp bytes.Frame

	if err := client.Call(ctx, req, &rsp); err != nil {
		return nil, fmt.Errorf("Call error: %v", err)
	}

	return mcp.NewToolResultText(string(rsp.Data)), nil
}

func Run(c *cli.Context) error {
	useSSE := c.Bool("sse")

	// Create MCP server
	s := server.NewMCPServer(
		"micro",
		"1.0.0",
	)

	// Add tool
	tool := mcp.NewTool("call",
		mcp.WithDescription("Call a service"),
		mcp.WithString("service",
			mcp.Required(),
			mcp.Description("Name of the service e.g helloworld"),
		),
		mcp.WithString("endpoint",
			mcp.Required(),
			mcp.Description("Name of the endpoint e.g Say.Hello"),
		),
		mcp.WithString("request",
			mcp.Required(),
			mcp.Description("JSON request body"),
		),
	)

	// Add tool handler
	s.AddTool(tool, handler)

	// Add tool
	describe := mcp.NewTool("describe",
		mcp.WithDescription("Describe a service"),
		mcp.WithString("service",
			mcp.Required(),
			mcp.Description("Name of the service e.g helloworld"),
		),
	)

	// Add describe handler
	s.AddTool(describe, func(ctx context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	// Add tool
	services := mcp.NewTool("services",
		mcp.WithDescription("List services"),
	)

	// Add describe handler
	s.AddTool(services, func(ctx context.Context, r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	if useSSE {
		// Use the SSE server from the mcp-go/server package, passing the MCP server
		return server.NewSSEServer(s).ListenAndServe()
	}

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}

	return nil
}

func init() {
	cmd.Register(&cli.Command{
		Name:   "mcp",
		Usage:  "Run MCP server on stdio or SSE (with --sse)",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "sse",
				Usage: "Run as SSE server on :8080 instead of stdio",
			},
		},
		Action: Run,
	})
}
