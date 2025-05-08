package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/urfave/cli/v2"
	"go-micro.dev/v5/client"
	"go-micro.dev/v5/cmd"
	"go-micro.dev/v5/codec/bytes"
	"go-micro.dev/v5/registry"
)

var (
	// version is set by the release action
	// this is the default for local builds
	version = "5.0.0-dev"
)

type Command struct {
	Name   string
	Usage  string
	Action func(*cli.Context, []string) error
}

var commands = []Command{
	{
		Name:  "services",
		Usage: "List available services",
		Action: func(ctx *cli.Context, args []string) error {
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
		Action: func(ctx *cli.Context, args []string) error {
			if len(args) < 2 {
				return fmt.Errorf("Usage: [service] [endpoint] [request]")
			}

			service := args[0]
			endpoint := args[1]
			request := `{}`

			// get the request if present
			if len(args) >= 3 {
				request = strings.Join(args[2:], " ")
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
		Action: func(ctx *cli.Context, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("Usage: [service]")
			}

			service := args[0]
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

func process(c *cli.Context) error {
	reader := bufio.NewReader(os.Stdin)

	commandMap := map[string]Command{}
	helpUsage := []string{}

	for _, c := range commands {
		commandMap[c.Name] = c
		helpUsage = append(helpUsage, fmt.Sprintf("%-20s%s", c.Name, c.Usage))
	}

	sort.Strings(helpUsage)

	for {
		fmt.Print("micro> ") // Print the prompt

		input, _ := reader.ReadString('\n') // Read input until a newline
		input = input[:len(input)-1]        // Remove the trailing newline

		args := strings.Split(input, " ")

		if len(args) == 0 {
			continue
		}

		command := args[0]

		if command == "exit" {
			fmt.Println("Exiting...")
			return nil
		}

		if v, ok := commandMap[command]; ok {
			err := v.Action(c, args[1:])
			if err != nil {
				fmt.Println(err)
			}
			continue
		}

		if command == "help" || command == "?" {
			fmt.Println("Commands:")
			fmt.Println(strings.Join(helpUsage, "\n"))
		}
	}
}

func main() {
	// set process
	app := cmd.App()
	app.Action = process
	app.Name = "micro"
	app.Version = version
	app.Run(os.Args)
}
