package cli

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/micro/cli"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/server/proto/health"
	"github.com/serenize/snaker"

	"golang.org/x/net/context"
)

func formatEndpoint(v *registry.Value, r int) string {
	// default format is tabbed plus the value plus new line
	fparts := []string{"", "%s %s", "\n"}
	for i := 0; i < r+1; i++ {
		fparts[0] += "\t"
	}
	// its just a primitive of sorts so return
	if len(v.Values) == 0 {
		return fmt.Sprintf(strings.Join(fparts, ""), snaker.CamelToSnake(v.Name), v.Type)
	}

	// this thing has more things, it's complex
	fparts[1] += " {"

	vals := []interface{}{snaker.CamelToSnake(v.Name), v.Type}

	for _, val := range v.Values {
		fparts = append(fparts, "%s")
		vals = append(vals, formatEndpoint(val, r+1))
	}

	// at the end
	l := len(fparts) - 1
	for i := 0; i < r+1; i++ {
		fparts[l] += "\t"
	}
	fparts = append(fparts, "}\n")

	return fmt.Sprintf(strings.Join(fparts, ""), vals...)
}

func listServices(c *cli.Context) {
	rsp, err := (*cmd.DefaultOptions().Registry).ListServices()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	ss := sortedServices{rsp}
	sort.Sort(ss)
	for _, service := range ss.services {
		fmt.Println(service.Name)
	}
}

func registerService(c *cli.Context) {
	if len(c.Args()) != 1 {
		fmt.Println("require service definition")
		return
	}
	var service *registry.Service
	if err := json.Unmarshal([]byte(c.Args().First()), &service); err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := (*cmd.DefaultOptions().Registry).Register(service); err != nil {
		fmt.Println(err.Error())
		return
	}
}

func deregisterService(c *cli.Context) {
	if len(c.Args()) != 1 {
		fmt.Println("require service definition")
		return
	}
	var service *registry.Service
	if err := json.Unmarshal([]byte(c.Args().First()), &service); err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := (*cmd.DefaultOptions().Registry).Deregister(service); err != nil {
		fmt.Println(err.Error())
		return
	}
}

func getService(c *cli.Context) {
	if !c.Args().Present() {
		fmt.Println("Service required")
		return
	}
	service, err := (*cmd.DefaultOptions().Registry).GetService(c.Args().First())
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
			request = "{\n"
			for _, v := range e.Request.Values {
				request += formatEndpoint(v, 0)
			}
			request += "}"
		} else {
			request = "{}"
		}
		if e.Response != nil && len(e.Response.Values) > 0 {
			response = "{\n"
			for _, v := range e.Response.Values {
				response += formatEndpoint(v, 0)
			}
			response += "}"
		} else {
			response = "{}"
		}
		fmt.Printf("\nEndpoint: %s\nMetadata: %s\n", e.Name, strings.Join(meta, ","))
		fmt.Printf("Request: %s\n\nResponse: %s\n", request, response)
	}
}

func queryService(c *cli.Context) {
	if len(c.Args()) < 2 {
		fmt.Println("require service and method")
		return
	}
	service := c.Args()[0]
	method := c.Args()[1]
	var request map[string]interface{}
	var response map[string]interface{}
	json.Unmarshal([]byte(strings.Join(c.Args()[2:], " ")), &request)
	req := (*cmd.DefaultOptions().Client).NewJsonRequest(service, method, request)
	err := (*cmd.DefaultOptions().Client).Call(context.Background(), req, &response)
	if err != nil {
		fmt.Printf("error calling %s.%s: %v\n", service, method, err)
		return
	}
	b, _ := json.MarshalIndent(response, "", "\t")
	fmt.Println(string(b))
}

func queryHealth(c *cli.Context) {
	if !c.Args().Present() {
		fmt.Println("require service name")
		return
	}
	service, err := (*cmd.DefaultOptions().Registry).GetService(c.Args().First())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if service == nil || len(service) == 0 {
		fmt.Println("Service not found")
		return
	}
	req := (*cmd.DefaultOptions().Client).NewRequest(service[0].Name, "Debug.Health", &health.Request{})
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
			err := (*cmd.DefaultOptions().Client).CallRemote(context.Background(), address, req, rsp)
			var status string
			if err != nil {
				status = err.Error()
			} else {
				status = rsp.Status
			}
			fmt.Printf("%s\t\t%s:%d\t\t%s\n", node.Id, node.Address, node.Port, status)
		}
	}
}
