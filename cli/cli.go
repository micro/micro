package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/myodc/go-micro/client"
	"github.com/myodc/go-micro/proto/health"
	"github.com/myodc/go-micro/registry"
	"github.com/myodc/go-micro/store"
)

func registryCommands() []cli.Command {
	return []cli.Command{
		{
			Name:  "list",
			Usage: "List items in registry",
			Subcommands: []cli.Command{
				{
					Name:  "services",
					Usage: "List services in registry",
					Action: func(c *cli.Context) {
						rsp, err := registry.ListServices()
						if err != nil {
							fmt.Println(err.Error())
							return
						}
						for _, service := range rsp {
							fmt.Println(service.Name())
						}
					},
				},
			},
		},
		{
			Name:  "get",
			Usage: "Get item from registry",
			Subcommands: []cli.Command{
				{
					Name:  "service",
					Usage: "Get service from registry",
					Action: func(c *cli.Context) {
						if !c.Args().Present() {
							fmt.Println("Service required")
							return
						}
						service, err := registry.GetService(c.Args().First())
						if err != nil {
							fmt.Println(err.Error())
							return
						}
						fmt.Printf("%s\n\n", service.Name())
						for _, node := range service.Nodes() {
							fmt.Printf("%s\t%s\t%d\n", node.Id(), node.Address(), node.Port())
						}
					},
				},
			},
		},
	}
}

func storeCommands() []cli.Command {
	return []cli.Command{
		{
			Name:  "get",
			Usage: "Get item from store",
			Action: func(c *cli.Context) {
				if !c.Args().Present() {
					fmt.Println("Key required")
					return
				}
				item, err := store.Get(c.Args().First())
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				fmt.Println(string(item.Value()))
			},
		},
		{
			Name:  "del",
			Usage: "Delete item from store",
			Action: func(c *cli.Context) {
				if !c.Args().Present() {
					fmt.Println("Key required")
					return
				}
				if err := store.Del(c.Args().First()); err != nil {
					fmt.Println(err.Error())
					return
				}
			},
		},
	}
}

func Commands() []cli.Command {
	commands := []cli.Command{
		{
			Name:        "registry",
			Usage:       "Query registry",
			Subcommands: registryCommands(),
		},
		{
			Name:        "store",
			Usage:       "Query store",
			Subcommands: storeCommands(),
		},
		{
			Name:  "query",
			Usage: "Query a service method using rpc",
			Action: func(c *cli.Context) {
				if len(c.Args()) < 2 {
					fmt.Println("require service and method")
					return
				}
				service := c.Args()[0]
				method := c.Args()[1]
				var request map[string]interface{}
				var response map[string]interface{}
				json.Unmarshal([]byte(strings.Join(c.Args()[2:], " ")), &request)
				req := client.NewJsonRequest(service, method, request)
				err := client.Call(req, &response)
				if err != nil {
					fmt.Printf("error calling %s.%s: %v\n", service, method, err)
					return
				}
				b, _ := json.MarshalIndent(response, "", "\t")
				fmt.Println(string(b))
			},
		},
		{
			Name:  "health",
			Usage: "Query the health of a service",
			Action: func(c *cli.Context) {
				if !c.Args().Present() {
					fmt.Println("require service name")
					return
				}
				service, err := registry.GetService(c.Args().First())
				if err != nil {
					fmt.Printf("error querying registry: %v", err)
					return
				}
				fmt.Println("node\t\taddress:port\t\tstatus")
				req := client.NewRequest(service.Name(), "Debug.Health", &health.Request{})
				for _, node := range service.Nodes() {
					address := node.Address()
					if node.Port() > 0 {
						address = fmt.Sprintf("%s:%d", address, node.Port())
					}
					rsp := &health.Response{}
					err := client.CallRemote(address, "", req, rsp)
					var status string
					if err != nil {
						status = err.Error()
					} else {
						status = rsp.Status
					}
					fmt.Printf("%s\t\t%s:%d\t\t%s\n", node.Id(), node.Address(), node.Port(), status)
				}
			},
		},
	}

	return append(commands, registryCommands()...)
}
