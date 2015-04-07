package cli

import (
	"fmt"

	"github.com/asim/go-micro/registry"
	"github.com/asim/go-micro/store"
	"github.com/codegangsta/cli"
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
			Usage: "get item from registry",
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
	}

	return append(commands, registryCommands()...)
}
