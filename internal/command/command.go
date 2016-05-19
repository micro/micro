package command

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/micro/cli"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/registry"

	proto "github.com/micro/go-micro/server/debug/proto"

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

func del(url string, b []byte, v interface{}) error {
	if !strings.HasPrefix(url, "http") && !strings.HasPrefix(url, "https") {
		url = "http://" + url
	}

	buf := bytes.NewBuffer(b)
	defer buf.Reset()

	req, err := http.NewRequest("DELETE", url, buf)
	if err != nil {
		return err
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	if v == nil {
		return nil
	}

	d := json.NewDecoder(rsp.Body)
	d.UseNumber()
	return d.Decode(v)
}

func get(url string, v interface{}) error {
	if !strings.HasPrefix(url, "http") && !strings.HasPrefix(url, "https") {
		url = "http://" + url
	}

	rsp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	d := json.NewDecoder(rsp.Body)
	d.UseNumber()
	return d.Decode(v)
}

func post(url string, b []byte, v interface{}) error {
	if !strings.HasPrefix(url, "http") && !strings.HasPrefix(url, "https") {
		url = "http://" + url
	}

	buf := bytes.NewBuffer(b)
	defer buf.Reset()

	rsp, err := http.Post(url, "application/json", buf)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	if v == nil {
		return nil
	}

	d := json.NewDecoder(rsp.Body)
	d.UseNumber()
	return d.Decode(v)
}

func RegisterService(c *cli.Context, args []string) ([]byte, error) {
	if len(args) == 0 {
		return nil, errors.New("require service definition")
	}

	req := strings.Join(args, " ")

	if p := c.GlobalString("proxy_address"); len(p) > 0 {
		if err := post(p+"/registry", []byte(req), nil); err != nil {
			return nil, err
		}
		return []byte("ok"), nil
	}

	var service *registry.Service

	d := json.NewDecoder(strings.NewReader(req))
	d.UseNumber()

	if err := d.Decode(&service); err != nil {
		return nil, err
	}

	if err := (*cmd.DefaultOptions().Registry).Register(service); err != nil {
		return nil, err
	}

	return []byte("ok"), nil
}

func DeregisterService(c *cli.Context, args []string) ([]byte, error) {
	if len(args) == 0 {
		return nil, errors.New("require service definition")
	}

	req := strings.Join(args, " ")

	if p := c.GlobalString("proxy_address"); len(p) > 0 {
		if err := del(p+"/registry", []byte(req), nil); err != nil {
			return nil, err
		}
		return []byte("ok"), nil
	}

	var service *registry.Service

	d := json.NewDecoder(strings.NewReader(req))
	d.UseNumber()

	if err := d.Decode(&service); err != nil {
		return nil, err
	}

	if err := (*cmd.DefaultOptions().Registry).Deregister(service); err != nil {
		return nil, err
	}

	return []byte("ok"), nil
}

func GetService(c *cli.Context, args []string) ([]byte, error) {
	if len(args) == 0 {
		return nil, errors.New("service requested")
	}

	var output []string
	var service []*registry.Service
	var err error

	if p := c.GlobalString("proxy_address"); len(p) > 0 {
		if err := get(p+"/registry?service="+args[0], &service); err != nil {
			return nil, err
		}
	} else {
		service, err = (*cmd.DefaultOptions().Registry).GetService(args[0])
	}

	if err != nil {
		return nil, err
	}

	if len(service) == 0 {
		return nil, errors.New("Service not found")
	}

	output = append(output, "service  "+service[0].Name)

	for _, serv := range service {
		output = append(output, "\nversion "+serv.Version)
		output = append(output, "\nId\tAddress\tPort\tMetadata")
		for _, node := range serv.Nodes {
			var meta []string
			for k, v := range node.Metadata {
				meta = append(meta, k+"="+v)
			}
			output = append(output, fmt.Sprintf("%s\t%s\t%d\t%s", node.Id, node.Address, node.Port, strings.Join(meta, ",")))
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

		output = append(output, fmt.Sprintf("\nEndpoint: %s\nMetadata: %s\n", e.Name, strings.Join(meta, ",")))
		output = append(output, fmt.Sprintf("Request: %s\n\nResponse: %s\n", request, response))
	}

	return []byte(strings.Join(output, "\n")), nil
}

func ListServices(c *cli.Context) ([]byte, error) {
	var rsp []*registry.Service
	var err error

	if p := c.GlobalString("proxy_address"); len(p) > 0 {
		if err := get(p+"/registry", &rsp); err != nil {
			return nil, err
		}
	} else {
		rsp, err = (*cmd.DefaultOptions().Registry).ListServices()
		if err != nil {
			return nil, err
		}
	}

	sort.Sort(sortedServices{rsp})

	var services []string

	for _, service := range rsp {
		services = append(services, service.Name)
	}

	return []byte(strings.Join(services, "\n")), nil
}

func QueryService(c *cli.Context, args []string) ([]byte, error) {
	if len(args) < 2 {
		return nil, errors.New("require service and method")
	}

	var req, service, method string
	service = args[0]
	method = args[1]

	if len(args) > 2 {
		req = strings.Join(args[2:], " ")
	}

	// empty request
	if len(req) == 0 {
		req = `{}`
	}

	var request map[string]interface{}
	var response json.RawMessage

	if p := c.GlobalString("proxy_address"); len(p) > 0 {
		request = map[string]interface{}{
			"service": service,
			"method":  method,
			"request": req,
		}

		b, err := json.Marshal(request)
		if err != nil {
			return nil, err
		}

		if err := post(p+"/rpc", b, &response); err != nil {
			return nil, err
		}

	} else {
		d := json.NewDecoder(strings.NewReader(req))
		d.UseNumber()

		if err := d.Decode(&request); err != nil {
			return nil, err
		}

		creq := (*cmd.DefaultOptions().Client).NewJsonRequest(service, method, request)
		err := (*cmd.DefaultOptions().Client).Call(context.Background(), creq, &response)
		if err != nil {
			return nil, fmt.Errorf("error calling %s.%s: %v\n", service, method, err)
		}
	}

	var out bytes.Buffer
	defer out.Reset()
	if err := json.Indent(&out, response, "", "\t"); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func QueryHealth(c *cli.Context, args []string) ([]byte, error) {
	if len(args) == 0 {
		return nil, errors.New("require service name")
	}

	service, err := (*cmd.DefaultOptions().Registry).GetService(args[0])
	if err != nil {
		return nil, err
	}

	if service == nil || len(service) == 0 {
		return nil, errors.New("Service not found")
	}

	req := (*cmd.DefaultOptions().Client).NewRequest(service[0].Name, "Debug.Health", &proto.HealthRequest{})

	var output []string

	// print things
	output = append(output, "service  "+service[0].Name)

	for _, serv := range service {
		// print things
		output = append(output, "\nversion "+serv.Version)
		output = append(output, "\nnode\t\taddress:port\t\tstatus")

		// query health for every node
		for _, node := range serv.Nodes {
			address := node.Address
			if node.Port > 0 {
				address = fmt.Sprintf("%s:%d", address, node.Port)
			}
			rsp := &proto.HealthResponse{}

			var err error

			if p := c.GlobalString("proxy_address"); len(p) > 0 {
				// call using proxy
				request := map[string]interface{}{
					"service": service[0].Name,
					"method":  "Debug.Health",
					"address": address,
				}

				b, err := json.Marshal(request)
				if err != nil {
					return nil, err
				}

				if err := post(p+"/rpc", b, &rsp); err != nil {
					return nil, err
				}
			} else {
				// call using client
				err = (*cmd.DefaultOptions().Client).CallRemote(context.Background(), address, req, rsp)
			}

			var status string
			if err != nil {
				status = err.Error()
			} else {
				status = rsp.Status
			}
			output = append(output, fmt.Sprintf("%s\t\t%s:%d\t\t%s", node.Id, node.Address, node.Port, status))
		}
	}

	return []byte(strings.Join(output, "\n")), nil
}
