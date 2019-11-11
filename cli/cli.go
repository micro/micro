// Package cli is a command line interface
package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/micro/cli"
)

var (
	prompt = "micro> "

	commands = map[string]*command{
		"quit":       &command{"quit", "Exit the CLI", quit},
		"exit":       &command{"exit", "Exit the CLI", quit},
		"call":       &command{"call", "Call a service", callService},
		"list":       &command{"list", "List services, peers or routes", list},
		"get":        &command{"get", "Get service info", getService},
		"services":   &command{"services", "List services in the network", netServices},
		"stream":     &command{"stream", "Stream a call to a service", streamService},
		"publish":    &command{"publish", "Publish a message to a topic", publish},
		"health":     &command{"health", "Get service health", queryHealth},
		"stats":      &command{"stats", "Get service stats", queryStats},
		"register":   &command{"register", "Register a service", registerService},
		"deregister": &command{"deregister", "Deregister a service", deregisterService},
	}
)

type command struct {
	name  string
	usage string
	exec  exec
}

func runc(c *cli.Context) {
	commands["help"] = &command{"help", "CLI usage", help}
	alias := map[string]string{
		"?":  "help",
		"ls": "list",
	}

	r, err := readline.New(prompt)
	if err != nil {
		fmt.Fprint(os.Stdout, err)
		os.Exit(1)
	}
	defer r.Close()

	for {
		args, err := r.Readline()
		if err != nil {
			fmt.Fprint(os.Stdout, err)
			return
		}

		args = strings.TrimSpace(args)

		// skip no args
		if len(args) == 0 {
			continue
		}

		parts := strings.Split(args, " ")
		if len(parts) == 0 {
			continue
		}

		name := parts[0]

		// get alias
		if n, ok := alias[name]; ok {
			name = n
		}

		if cmd, ok := commands[name]; ok {
			rsp, err := cmd.exec(c, parts[1:])
			if err != nil {
				println(err.Error())
				continue
			}
			println(string(rsp))
		} else {
			println("unknown command")
		}
	}
}

func HealthCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "check",
			Usage:  "Query the health of a service",
			Action: printer(queryHealth),
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "address",
					Usage:  "Set the address of the service instance to call",
					EnvVar: "MICRO_ADDRESS",
				},
			},
		},
	}
}

func NetworkCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "connect",
			Usage:  "connect to the network. specify nodes e.g connect ip:port",
			Action: printer(networkConnect),
		},
		{
			Name:   "connections",
			Usage:  "List the immediate connections to the network",
			Action: printer(networkConnections),
		},
		{
			Name:   "graph",
			Usage:  "Get the network graph",
			Action: printer(networkGraph),
		},
		{
			Name:   "nodes",
			Usage:  "List nodes in the network",
			Action: printer(netNodes),
		},
		{
			Name:   "routes",
			Usage:  "List network routes",
			Action: printer(netRoutes),
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "service",
					Usage: "Filter by service",
				},
				cli.StringFlag{
					Name:  "address",
					Usage: "Filter by address",
				},
				cli.StringFlag{
					Name:  "gateway",
					Usage: "Filter by gateway",
				},
				cli.StringFlag{
					Name:  "router",
					Usage: "Filter by router",
				},
				cli.StringFlag{
					Name:  "network",
					Usage: "Filter by network",
				},
			},
		},
		{
			Name:   "services",
			Usage:  "List network services",
			Action: printer(netServices),
		},
	}
}

func NetworkDNSCommands() []cli.Command {
	return []cli.Command{
		{
			Name:  "advertise",
			Usage: "Advertise a new node to the network",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "address",
					Usage:  "Address to register for the specified domain",
					EnvVar: "MICRO_NETWORK_DNS_ADVERTISE_ADDRESS",
				},
				cli.StringFlag{
					Name:   "domain",
					Usage:  "Domain name to register",
					EnvVar: "MICRO_NETWORK_DNS_ADVERTISE_DOMAIN",
					Value:  "network.micro.mu",
				},
				cli.StringFlag{
					Name:   "token",
					Usage:  "Bearer token for the go.micro.network.dns service",
					EnvVar: "MICRO_NETWORK_DNS_ADVERTISE_TOKEN",
				},
			},
			Action: printer(netDNSAdvertise),
		},
		{
			Name:  "remove",
			Usage: "Remove a node's record'",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "address",
					Usage:  "Address to register for the specified domain",
					EnvVar: "MICRO_NETWORK_DNS_REMOVE_ADDRESS",
				},
				cli.StringFlag{
					Name:   "domain",
					Usage:  "Domain name to remove",
					EnvVar: "MICRO_NETWORK_DNS_REMOVE_DOMAIN",
					Value:  "network.micro.mu",
				},
				cli.StringFlag{
					Name:   "token",
					Usage:  "Bearer token for the go.micro.network.dns service",
					EnvVar: "MICRO_NETWORK_DNS_REMOVE_TOKEN",
				},
			},
			Action: printer(netDNSRemove),
		},
		{
			Name:  "resolve",
			Usage: "Remove a record'",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "domain",
					Usage:  "Domain name to resolve",
					EnvVar: "MICRO_NETWORK_DNS_RESOLVE_DOMAIN",
					Value:  "network.micro.mu",
				},
				cli.StringFlag{
					Name:   "type",
					Usage:  "Domain name type to resolve",
					EnvVar: "MICRO_NETWORK_DNS_RESOLVE_TYPE",
					Value:  "A",
				},
				cli.StringFlag{
					Name:   "token",
					Usage:  "Bearer token for the go.micro.network.dns service",
					EnvVar: "MICRO_NETWORK_DNS_RESOLVE_TOKEN",
				},
			},
			Action: printer(netDNSResolve),
		},
	}
}

func RegistryCommands() []cli.Command {
	return []cli.Command{
		{
			Name:  "list",
			Usage: "List items in registry or network",
			Subcommands: []cli.Command{
				{
					Name:   "nodes",
					Usage:  "List nodes in the network",
					Action: printer(netNodes),
				},
				{
					Name:   "routes",
					Usage:  "List network routes",
					Action: printer(netRoutes),
				},
				{
					Name:   "services",
					Usage:  "List services in registry",
					Action: printer(listServices),
				},
			},
		},
		{
			Name:  "register",
			Usage: "Register an item in the registry",
			Subcommands: []cli.Command{
				{
					Name:   "service",
					Usage:  "Register a service with JSON definition",
					Action: printer(registerService),
				},
			},
		},
		{
			Name:  "deregister",
			Usage: "Deregister an item in the registry",
			Subcommands: []cli.Command{
				{
					Name:   "service",
					Usage:  "Deregister a service with JSON definition",
					Action: printer(deregisterService),
				},
			},
		},
		{
			Name:  "get",
			Usage: "Get item from registry",
			Subcommands: []cli.Command{
				{
					Name:   "service",
					Usage:  "Get service from registry",
					Action: printer(getService),
				},
			},
		},
	}
}

func Commands() []cli.Command {
	commands := []cli.Command{
		{
			Name:   "cli",
			Usage:  "Run the interactive CLI",
			Action: runc,
		},
		{
			Name:   "call",
			Usage:  "Call a service e.g micro call greeter Say.Hello '{\"name\": \"John\"}",
			Action: printer(callService),
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "address",
					Usage:  "Set the address of the service instance to call",
					EnvVar: "MICRO_ADDRESS",
				},
				cli.StringFlag{
					Name:   "output, o",
					Usage:  "Set the output format; json (default), raw",
					EnvVar: "MICRO_OUTPUT",
				},
				cli.StringSliceFlag{
					Name:   "metadata",
					Usage:  "A list of key-value pairs to be forwarded as metadata",
					EnvVar: "MICRO_METADATA",
				},
			},
		},
		{
			Name:   "services",
			Usage:  "List the services in the network",
			Action: printer(netServices),
		},
		{
			Name:   "stream",
			Usage:  "Create a service stream",
			Action: printer(streamService),
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "output, o",
					Usage:  "Set the output format; json (default), raw",
					EnvVar: "MICRO_OUTPUT",
				},
				cli.StringSliceFlag{
					Name:   "metadata",
					Usage:  "A list of key-value pairs to be forwarded as metadata",
					EnvVar: "MICRO_METADATA",
				},
			},
		},
		{
			Name:   "publish",
			Usage:  "Publish a message to a topic",
			Action: printer(publish),
			Flags: []cli.Flag{
				cli.StringSliceFlag{
					Name:   "metadata",
					Usage:  "A list of key-value pairs to be forwarded as metadata",
					EnvVar: "MICRO_METADATA",
				},
			},
		},
		{
			Name:   "stats",
			Usage:  "Query the stats of a service",
			Action: printer(queryStats),
		},
	}

	return append(commands, RegistryCommands()...)
}
