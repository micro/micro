package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v3/client"
	cbytes "github.com/micro/go-micro/v3/codec/bytes"
	proto "github.com/micro/go-micro/v3/debug/service/proto"
	"github.com/micro/go-micro/v3/metadata"
	"github.com/micro/go-micro/v3/registry"
	"github.com/micro/micro/v3/client/cli/namespace"
	"github.com/micro/micro/v3/client/cli/util"
	muclient "github.com/micro/micro/v3/service/client"
	muregistry "github.com/micro/micro/v3/service/registry"

	"github.com/serenize/snaker"
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

func callContext(c *cli.Context) context.Context {
	callMD := make(map[string]string)

	for _, md := range c.StringSlice("metadata") {
		parts := strings.Split(md, "=")
		if len(parts) < 2 {
			continue
		}

		key := parts[0]
		val := strings.Join(parts[1:], "=")

		// set the key/val
		callMD[key] = val
	}

	return metadata.NewContext(context.Background(), callMD)
}

func GetService(c *cli.Context, args []string) ([]byte, error) {
	if len(args) == 0 {
		return nil, errors.New("service required")
	}

	ns, err := namespace.Get(util.GetEnv(c).Name)
	if err != nil {
		return nil, err
	}

	var output []string
	var srv []*registry.Service

	reg := muregistry.DefaultRegistry

	srv, err = reg.GetService(args[0], registry.GetDomain(ns))
	if err != nil {
		return nil, err
	}
	if len(srv) == 0 {
		return nil, errors.New("Service not found")
	}

	output = append(output, "service  "+srv[0].Name)

	for _, serv := range srv {
		if len(serv.Version) > 0 {
			output = append(output, "\nversion "+serv.Version)
		}

		output = append(output, "\nID\tAddress\tMetadata")
		for _, node := range serv.Nodes {
			var meta []string
			for k, v := range node.Metadata {
				meta = append(meta, k+"="+v)
			}
			output = append(output, fmt.Sprintf("%s\t%s\t%s", node.Id, node.Address, strings.Join(meta, ",")))
		}
	}

	for _, e := range srv[0].Endpoints {
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

		output = append(output, fmt.Sprintf("\nEndpoint: %s\n", e.Name))

		// set metadata if exists
		if len(meta) > 0 {
			output = append(output, fmt.Sprintf("Metadata: %s\n", strings.Join(meta, ",")))
		}

		output = append(output, fmt.Sprintf("Request: %s\n\nResponse: %s\n", request, response))
	}

	return []byte(strings.Join(output, "\n")), nil
}

func ListServices(c *cli.Context) ([]byte, error) {
	var rsp []*registry.Service
	var err error

	ns, err := namespace.Get(util.GetEnv(c).Name)
	if err != nil {
		return nil, err
	}

	reg := muregistry.DefaultRegistry

	rsp, err = reg.ListServices(registry.ListDomain(ns))
	if err != nil {
		return nil, err
	}

	var services []string
	for _, service := range rsp {
		services = append(services, service.Name)
	}

	sort.Strings(services)

	return []byte(strings.Join(services, "\n")), nil
}

func Publish(c *cli.Context, args []string) error {
	if len(args) < 2 {
		return errors.New("require topic and message e.g micro publish event '{\"hello\": \"world\"}'")
	}
	defer func() {
		time.Sleep(time.Millisecond * 100)
	}()
	topic := args[0]
	message := args[1]

	cl := muclient.DefaultClient
	ct := func(o *client.MessageOptions) {
		o.ContentType = "application/json"
	}

	d := json.NewDecoder(strings.NewReader(message))
	d.UseNumber()

	var msg map[string]interface{}
	if err := d.Decode(&msg); err != nil {
		return err
	}

	ctx := callContext(c)
	m := cl.NewMessage(topic, msg, ct)
	return cl.Publish(ctx, m)
}

func CallService(c *cli.Context, args []string) ([]byte, error) {
	if len(args) < 2 {
		return nil, errors.New(`require service and endpoint e.g micro call greeeter Say.Hello '{"name": "john"}'`)
	}

	var req, service, endpoint string
	service = args[0]
	endpoint = args[1]

	if len(args) > 2 {
		req = strings.Join(args[2:], " ")
	}

	// empty request
	if len(req) == 0 {
		req = `{}`
	}

	var request map[string]interface{}
	var response []byte

	d := json.NewDecoder(strings.NewReader(req))
	d.UseNumber()

	if err := d.Decode(&request); err != nil {
		return nil, err
	}

	ctx := callContext(c)

	cli := muclient.DefaultClient
	creq := cli.NewRequest(service, endpoint, request, client.WithContentType("application/json"))

	var opts []client.CallOption

	if addr := c.String("address"); len(addr) > 0 {
		opts = append(opts, client.WithAddress(addr))
	}

	var err error
	if output := c.String("output"); output == "raw" {
		rsp := cbytes.Frame{}
		err = cli.Call(ctx, creq, &rsp, opts...)
		// set the raw output
		response = rsp.Data
	} else {
		var rsp json.RawMessage
		err = cli.Call(ctx, creq, &rsp, opts...)
		// set the response
		if err == nil {
			var out bytes.Buffer
			defer out.Reset()
			if err := json.Indent(&out, rsp, "", "\t"); err != nil {
				return nil, err
			}
			response = out.Bytes()
		}
	}

	if err != nil {
		return nil, fmt.Errorf("error calling %s.%s: %v", service, endpoint, err)
	}

	return response, nil
}

func QueryHealth(c *cli.Context, args []string) ([]byte, error) {
	if len(args) == 0 {
		return nil, errors.New("require service name")
	}

	ns, err := namespace.Get(util.GetEnv(c).Name)
	if err != nil {
		return nil, err
	}

	req := (muclient.DefaultClient).NewRequest(args[0], "Debug.Health", &proto.HealthRequest{})

	// if the address is specified then we just call it
	if addr := c.String("address"); len(addr) > 0 {
		rsp := &proto.HealthResponse{}
		err := (muclient.DefaultClient).Call(
			context.Background(),
			req,
			rsp,
			client.WithAddress(addr),
		)
		if err != nil {
			return nil, err
		}
		return []byte(rsp.Status), nil
	}

	// otherwise get the service and call each instance individually
	reg := muregistry.DefaultRegistry

	service, err := reg.GetService(args[0], registry.GetDomain(ns))
	if err != nil {
		return nil, err
	}

	if len(service) == 0 {
		return nil, errors.New("Service not found")
	}

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
			rsp := &proto.HealthResponse{}

			var err error

			// call using client
			err = (muclient.DefaultClient).Call(
				context.Background(),
				req,
				rsp,
				client.WithAddress(address),
			)

			var status string
			if err != nil {
				status = err.Error()
			} else {
				status = rsp.Status
			}
			output = append(output, fmt.Sprintf("%s\t\t%s\t\t%s", node.Id, node.Address, status))
		}
	}

	return []byte(strings.Join(output, "\n")), nil
}
