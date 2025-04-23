package main

import (
	"context"
	"encoding/json"
	"fmt"

	"go-micro.dev/v5/cmd"
	"go-micro.dev/v5/client"
	"go-micro.dev/v5/codec/bytes"
	"go-micro.dev/v5/registry"
	"github.com/urfave/cli/v2"
)

func main() {
	cmd.App().Commands = []*cli.Command{{
		Name: "services",
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
		Name: "call",
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
		Name: "describe",
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
	}}

	cmd.Init(
		cmd.Name("micro"),
		cmd.Version("5.0.0"),
	)
}

