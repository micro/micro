package cli

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/server/proto/health"

	"golang.org/x/net/context"
)

func formatEndpoint(v *registry.Value, r int) string {
	if len(v.Values) == 0 {
		fparts := []string{"\n", "{", "\n", "\t", "%s %s", "\n", "}"}
		if r > 0 {
			fparts = []string{"\n", "\t", "%s %s"}
			for i := 0; i < r; i++ {
				fparts[1] += "\t"
			}
		}
		return fmt.Sprintf(strings.Join(fparts, ""), strings.ToLower(v.Name), v.Type)
	}

	fparts := []string{"\n", "{", "\n", "\t", "%s %s", " {", "\n", "", "\n", "\t", "}", "\n", "}"}
	i := 7

	if r > 0 {
		fparts = []string{"\n", "\t", "%s %s", " {", "\n", "\t", "\n", "\t", "}"}
		i = 5
	}

	var app string
	for j := 0; j < r; j++ {
		if r > 0 {
			fparts[1] += "\t"
			fparts[7] += "\t"
		}
		app += "\t"
	}
	app += "\t%s"

	vals := []interface{}{strings.ToLower(v.Name), v.Type}

	for _, val := range v.Values {
		fparts[i] += app
		vals = append(vals, formatEndpoint(val, r+1))
	}

	return fmt.Sprintf(strings.Join(fparts, ""), vals...)
}

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
						ss := sortedServices{rsp}
						sort.Sort(ss)
						for _, service := range ss.services {
							fmt.Println(service.Name)
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
						if len(service) == 0 {
							fmt.Println("Service not found")
							return
						}

						fmt.Printf("service  %s\n", service[0].Name)
						for _, serv := range service {
							fmt.Println("\nversion ", serv.Version)
							fmt.Println("\nId\tAddress\tPort\tMetadata")
							for _, node := range serv.Nodes {
								var meta []string
								for k, v := range node.Metadata {
									meta = append(meta, k+"="+v)
								}
								fmt.Printf("%s\t%s\t%d\t%s\n", node.Id, node.Address, node.Port, strings.Join(meta, ","))
							}
						}

						for _, e := range service[0].Endpoints {
							var request, response string
							var meta []string
							for k, v := range e.Metadata {
								meta = append(meta, k+"="+v)
							}
							if e.Request != nil && len(e.Request.Values) > 0 {
								request = formatEndpoint(e.Request.Values[0], 0)
							} else {
								request = " {}"
							}
							if e.Response != nil && len(e.Response.Values) > 0 {
								response = formatEndpoint(e.Response.Values[0], 0)
							} else {
								response = " {}"
							}
							fmt.Printf("\nEndpoint: %s\nMetadata: %s\n", e.Name, strings.Join(meta, ","))
							fmt.Printf("Request:%s\nResponse:%s\n", request, response)
						}
					},
				},
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
				err := client.Call(context.Background(), req, &response)
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
					fmt.Println(err.Error())
					return
				}
				if service == nil || len(service) == 0 {
					fmt.Println("Service not found")
					return
				}
				req := client.NewRequest(service[0].Name, "Debug.Health", &health.Request{})
				fmt.Printf("service  %s\n\n", service[0].Name)
				for _, serv := range service {
					fmt.Println("\nversion ", serv.Version)
					fmt.Println("\nnode\t\taddress:port\t\tstatus")
					for _, node := range serv.Nodes {
						address := node.Address
						if node.Port > 0 {
							address = fmt.Sprintf("%s:%d", address, node.Port)
						}
						rsp := &health.Response{}
						err := client.CallRemote(context.Background(), address, req, rsp)
						var status string
						if err != nil {
							status = err.Error()
						} else {
							status = rsp.Status
						}
						fmt.Printf("%s\t\t%s:%d\t\t%s\n", node.Id, node.Address, node.Port, status)
					}
				}
			},
		},
	}

	return append(commands, registryCommands()...)
}
