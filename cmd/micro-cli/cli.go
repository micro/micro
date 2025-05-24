package cli

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/micro/micro/v5/cmd/micro-cli/new"
	"github.com/micro/micro/v5/util"
	"github.com/urfave/cli/v2"
	"go-micro.dev/v5/client"
	"go-micro.dev/v5/cmd"
	"go-micro.dev/v5/codec/bytes"
	"go-micro.dev/v5/registry"
	"go-micro.dev/v5/broker"
	"go-micro.dev/v5/config"
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
	{
		Name:  "broker",
		Usage: "Broker admin commands (publish, subscribe)",
		Action: func(ctx *cli.Context, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("Usage: broker [publish|subscribe] ...")
			}
			switch args[0] {
			case "publish":
				if len(args) < 3 {
					return fmt.Errorf("Usage: broker publish [topic] [message]")
				}
				topic, msg := args[1], args[2]
				if err := broker.DefaultBroker.Publish(topic, &broker.Message{Body: []byte(msg)}); err != nil {
					return err
				}
				fmt.Println("Publish OK")
				return nil
			case "subscribe":
				if len(args) < 2 {
					return fmt.Errorf("Usage: broker subscribe [topic]")
				}
				topic := args[1]
				ch := make(chan *broker.Message, 1)
				_, err := broker.DefaultBroker.Subscribe(topic, func(p broker.Event) error {
					ch <- p.Message()
					return nil
				})
				if err != nil {
					return err
				}
				msg := <-ch
				fmt.Printf("Received: %s\n", string(msg.Body))
				return nil
			default:
				return fmt.Errorf("Unknown broker command: %s", args[0])
			}
		},
	},
	{
		Name:  "config",
		Usage: "Config admin commands (get, set, delete, list)",
		Action: func(ctx *cli.Context, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("Usage: config [get|set|delete|list] ...")
			}
			switch args[0] {
			case "get":
				if len(args) < 2 {
					return fmt.Errorf("Usage: config get [key]")
				}
				val, err := config.DefaultConfig.Get(args[1])
				if err != nil {
					return err
				}
				fmt.Println(val.String())
				return nil
			case "set":
				if len(args) < 3 {
					return fmt.Errorf("Usage: config set [key] [value]")
				}
				if err := config.DefaultConfig.Set(args[1], args[2]); err != nil {
					return err
				}
				fmt.Println("Set OK")
				return nil
			case "delete":
				if len(args) < 2 {
					return fmt.Errorf("Usage: config delete [key]")
				}
				if err := config.DefaultConfig.Delete(args[1]); err != nil {
					return err
				}
				fmt.Println("Delete OK")
				return nil
			case "list":
				vals, err := config.DefaultConfig.List()
				if err != nil {
					return err
				}
				b, _ := json.MarshalIndent(vals, "", "    ")
				fmt.Println(string(b))
				return nil
			default:
				return fmt.Errorf("Unknown config command: %s", args[0])
			}
		},
	},
	{
		Name:  "registry",
		Usage: "Registry admin commands (list, get, register, deregister)",
		Action: func(ctx *cli.Context, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("Usage: registry [list|get|register|deregister] ...")
			}
			switch args[0] {
			case "list":
				services, err := registry.DefaultRegistry.ListServices()
				if err != nil {
					return err
				}
				b, _ := json.MarshalIndent(services, "", "    ")
				fmt.Println(string(b))
				return nil
			case "get":
				if len(args) < 2 {
					return fmt.Errorf("Usage: registry get [name]")
				}
				svc, err := registry.DefaultRegistry.GetService(args[1])
				if err != nil {
					return err
				}
				b, _ := json.MarshalIndent(svc, "", "    ")
				fmt.Println(string(b))
				return nil
			case "register":
				if len(args) < 4 {
					return fmt.Errorf("Usage: registry register [name] [node_id] [address] [version]")
				}
				name, nodeID, address := args[1], args[2], args[3]
				version := ""
				if len(args) > 4 {
					version = args[4]
				}
				svc := &registry.Service{
					Name:    name,
					Version: version,
					Nodes: []*registry.Node{{
						Id:      nodeID,
						Address: address,
					}},
				}
				if err := registry.DefaultRegistry.Register(svc); err != nil {
					return err
				}
				fmt.Println("Register OK")
				return nil
			case "deregister":
				if len(args) < 3 {
					return fmt.Errorf("Usage: registry deregister [name] [node_id]")
				}
				name, nodeID := args[1], args[2]
				svc := &registry.Service{
					Name: name,
					Nodes: []*registry.Node{{Id: nodeID}},
				}
				if err := registry.DefaultRegistry.Deregister(svc); err != nil {
					return err
				}
				fmt.Println("Deregister OK")
				return nil
			default:
				return fmt.Errorf("Unknown registry command: %s", args[0])
			}
		},
	},
}

func Run(c *cli.Context) error {
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
			continue
		}

		if srv, err := util.LookupService(command); err != nil {
			fmt.Println(util.CliError(err))
		} else if srv != nil && util.ShouldRenderHelp(args) {
			fmt.Println(cli.Exit(util.FormatServiceUsage(srv, c), 0))
		} else if srv != nil {
			err := util.CallService(srv, args)
			if err != nil {
				fmt.Println(util.CliError(err))
			}
		}
	}
}

func init() {
	cmd.Register(&cli.Command{
		Name:   "cli",
		Usage:  "Launch the interactive CLI",
		Action: Run,
	})
	cmd.Register(&cli.Command{
		Name:        "new",
		Usage:       "Create a new service",
		Description: `'micro new' generates a new service skeleton. Example: 'micro new helloworld && cd helloworld'`,
		Action:      new.Run,
	})
}
