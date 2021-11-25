package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/micro/micro/v3/client/cli/namespace"
	"github.com/micro/micro/v3/client/cli/util"
	proto "github.com/micro/micro/v3/proto/debug"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/context/metadata"
	"github.com/micro/micro/v3/service/registry"
	cbytes "github.com/micro/micro/v3/util/codec/bytes"
	"github.com/serenize/snaker"
	"github.com/urfave/cli/v2"
)

func quit(c *cli.Context, args []string) ([]byte, error) {
	os.Exit(0)
	return nil, nil
}

func help(c *cli.Context, args []string) ([]byte, error) {
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)

	fmt.Fprintln(os.Stdout, "Commands:")

	var keys []string
	for k := range commands {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		cmd := commands[k]
		fmt.Fprintln(w, "\t", cmd.name, "\t\t", cmd.usage)
	}

	w.Flush()
	return nil, nil
}

// QueryStats returns stats of specified service(s)
func QueryStats(c *cli.Context, args []string) ([]byte, error) {
	if c.String("all") == "builtin" {

		sl, err := ListServices(c, args)
		if err != nil {
			return nil, err
		}

		servList := strings.Split(string(sl), "\n")
		var builtinList []string

		for _, s := range servList {

			if util.IsBuiltInService(s) {
				builtinList = append(builtinList, s)
			}
		}

		c.Set("all", "")

		if len(builtinList) == 0 {
			return nil, errors.New("no builtin service(s) found")
		}

		return QueryStats(c, builtinList)
	}

	if c.String("all") == "custom" {
		sl, err := ListServices(c, args)
		if err != nil {
			return nil, err
		}

		servList := strings.Split(string(sl), "\n")
		var customList []string

		for _, s := range servList {

			// temporary excluding server
			if s == "server" {
				continue
			}

			if !util.IsBuiltInService(s) {
				customList = append(customList, s)
			}
		}

		c.Set("all", "")

		if len(customList) == 0 {
			return nil, errors.New("no custom service(s) found")
		}

		return QueryStats(c, customList)
	}

	if len(args) == 0 {
		return nil, cli.ShowSubcommandHelp(c)
	}

	env, err := util.GetEnv(c)
	if err != nil {
		return nil, err
	}
	ns, err := namespace.Get(env.Name)
	if err != nil {
		return nil, err
	}

	titlesList := []string{"NODE", "ADDRESS:PORT", "STARTED", "UPTIME", "MEMORY", "THREADS", "GC"}
	titles := strings.Join(titlesList, "\t")

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 1, 4, ' ', tabwriter.TabIndent)

	for _, a := range args {

		service, err := registry.DefaultRegistry.GetService(a, registry.GetDomain(ns))
		if err != nil {
			return nil, err
		}
		if len(service) == 0 {
			return nil, errors.New("Service not found")
		}

		req := client.NewRequest(service[0].Name, "Debug.Stats", &proto.StatsRequest{})

		fmt.Fprintln(w, "SERVICE\t"+service[0].Name+"\n")

		for _, serv := range service {

			fmt.Fprintln(w, "VERSION\t"+serv.Version+"\n")
			fmt.Fprintln(w, titles)

			// query health for every node
			for _, node := range serv.Nodes {
				address := node.Address
				rsp := &proto.StatsResponse{}

				var err error

				// call using client
				err = client.DefaultClient.Call(context.Background(), req, rsp, client.WithAddress(address))

				var started, uptime, memory, gc string
				if err == nil {
					started = time.Unix(int64(rsp.Started), 0).Format("Jan 2 15:04:05")
					uptime = fmt.Sprintf("%v", time.Duration(rsp.Uptime)*time.Second)
					memory = fmt.Sprintf("%.2fmb", float64(rsp.Memory)/(1024.0*1024.0))
					gc = fmt.Sprintf("%v", time.Duration(rsp.Gc))
				}

				line := fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%d\t%s\n",
					node.Id, node.Address, started, uptime, memory, rsp.Threads, gc)

				fmt.Fprintln(w, line)
			}
		}
	}

	w.Flush()

	return buf.Bytes(), nil
}

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
		return nil, cli.ShowSubcommandHelp(c)
	}

	env, err := util.GetEnv(c)
	if err != nil {
		return nil, err
	}
	ns, err := namespace.Get(env.Name)
	if err != nil {
		return nil, err
	}

	var output []string
	var srv []*registry.Service

	srv, err = registry.DefaultRegistry.GetService(args[0], registry.GetDomain(ns))
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

func ListServices(c *cli.Context, args []string) ([]byte, error) {
	var rsp []*registry.Service
	var err error

	env, err := util.GetEnv(c)
	if err != nil {
		return nil, err
	}
	ns, err := namespace.Get(env.Name)
	if err != nil {
		return nil, err
	}

	rsp, err = registry.DefaultRegistry.ListServices(registry.ListDomain(ns))
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
		return cli.ShowSubcommandHelp(c)
	}
	defer func() {
		time.Sleep(time.Millisecond * 100)
	}()
	topic := args[0]
	message := args[1]

	ct := func(o *client.MessageOptions) {
		o.ContentType = "application/json"
	}

	d := json.NewDecoder(strings.NewReader(message))
	d.UseNumber()

	var msg map[string]interface{}
	if err := d.Decode(&msg); err != nil {
		return cli.Exit(fmt.Sprintf("Error creating request %s", err), 1)
	}

	ctx := callContext(c)
	m := client.DefaultClient.NewMessage(topic, msg, ct)
	return client.Publish(ctx, m)
}

func CallService(c *cli.Context, args []string) ([]byte, error) {
	if len(args) < 2 {
		return nil, cli.ShowSubcommandHelp(c)
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
		return nil, cli.Exit(fmt.Sprintf("Error creating request %s", err), 1)
	}

	ctx := callContext(c)

	creq := client.DefaultClient.NewRequest(service, endpoint, request, client.WithContentType("application/json"))

	opts := []client.CallOption{client.WithAuthToken()}
	if timeout := c.String("request_timeout"); timeout != "" {
		duration, err := time.ParseDuration(timeout)
		if err != nil {
			return nil, cli.Exit("Invalid format for request_timeout duration. Try 500ms or 5s", 2)
		}
		opts = append(opts, client.WithRequestTimeout(duration))
	}

	if addr := c.String("address"); len(addr) > 0 {
		opts = append(opts, client.WithAddress(addr))
	}

	var err error
	if output := c.String("output"); output == "raw" {
		rsp := cbytes.Frame{}
		err = client.DefaultClient.Call(ctx, creq, &rsp, opts...)
		// set the raw output
		response = rsp.Data
	} else {
		var rsp json.RawMessage
		err = client.DefaultClient.Call(ctx, creq, &rsp, opts...)
		// set the response
		if err == nil {
			var out bytes.Buffer
			defer out.Reset()
			if err := json.Indent(&out, rsp, "", "\t"); err != nil {
				return nil, cli.Exit("Error while trying to format the response", 3)
			}
			response = out.Bytes()
		}
	}

	if err != nil {
		return nil, err
	}

	return response, nil
}

func QueryHealth(c *cli.Context, args []string) ([]byte, error) {
	if len(args) == 0 {
		return nil, errors.New("require service name")
	}

	env, err := util.GetEnv(c)
	if err != nil {
		return nil, err
	}
	ns, err := namespace.Get(env.Name)
	if err != nil {
		return nil, err
	}

	req := client.NewRequest(args[0], "Debug.Health", &proto.HealthRequest{})

	// if the address is specified then we just call it
	if addr := c.String("address"); len(addr) > 0 {
		rsp := &proto.HealthResponse{}
		err := client.DefaultClient.Call(
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
	service, err := registry.DefaultRegistry.GetService(args[0], registry.GetDomain(ns))
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
			err = client.DefaultClient.Call(
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

func getEnv(c *cli.Context, args []string) ([]byte, error) {
	env, err := util.GetEnv(c)
	if err != nil {
		return nil, err
	}
	return []byte(env.Name), nil
}

func setEnv(c *cli.Context, args []string) ([]byte, error) {
	if len(args) == 0 {
		return nil, cli.ShowSubcommandHelp(c)
	}
	return nil, util.SetEnv(args[0])
}

func listEnvs(c *cli.Context, args []string) ([]byte, error) {
	envs, err := util.GetEnvs()
	if err != nil {
		return nil, err
	}
	sort.Slice(envs, func(i, j int) bool { return envs[i].Name < envs[j].Name })
	current, err := util.GetEnv(c)
	if err != nil {
		return nil, err
	}

	byt := bytes.NewBuffer([]byte{})

	w := tabwriter.NewWriter(byt, 0, 0, 1, ' ', 0)
	for i, env := range envs {
		if i > 0 {
			fmt.Fprintf(w, "\n")
		}
		prefix := " "
		if env.Name == current.Name {
			prefix = "*"
		}
		if env.ProxyAddress == "" {
			env.ProxyAddress = "none"
		}
		fmt.Fprintf(w, "%v %v \t %v \t\t %v", prefix, env.Name, env.ProxyAddress, env.Description)
	}
	w.Flush()
	return byt.Bytes(), nil
}

func addEnv(c *cli.Context, args []string) ([]byte, error) {
	if len(args) == 0 {
		return nil, cli.ShowSubcommandHelp(c)
	}
	if len(args) == 1 {
		args = append(args, "") // default to no proxy address
	}

	return nil, util.AddEnv(util.Env{
		Name:         args[0],
		ProxyAddress: args[1],
	})
}

func delEnv(c *cli.Context, args []string) ([]byte, error) {
	if len(args) == 0 {
		return nil, cli.ShowSubcommandHelp(c)
	}
	return nil, util.DelEnv(c, args[0])
}

// TODO: stream via HTTP
func streamService(c *cli.Context, args []string) ([]byte, error) {
	if len(args) < 2 {
		return nil, cli.ShowSubcommandHelp(c)
	}
	service := args[0]
	endpoint := args[1]
	var request map[string]interface{}

	// ignore error
	json.Unmarshal([]byte(strings.Join(args[2:], " ")), &request)

	ctx := callContext(c)
	opts := []client.CallOption{client.WithAuthToken()}

	req := client.DefaultClient.NewRequest(service, endpoint, request, client.WithContentType("application/json"))
	stream, err := client.DefaultClient.Stream(ctx, req, opts...)
	if err != nil {
		if cerr := util.CliError(err); cerr.ExitCode() != 128 {
			return nil, cerr
		}
		return nil, fmt.Errorf("error calling %s.%s: %v", service, endpoint, err)
	}

	if err := stream.Send(request); err != nil {
		if cerr := util.CliError(err); cerr.ExitCode() != 128 {
			return nil, cerr
		}
		return nil, fmt.Errorf("error sending to %s.%s: %v", service, endpoint, err)
	}

	output := c.String("output")

	for {
		if output == "raw" {
			rsp := cbytes.Frame{}
			if err := stream.Recv(&rsp); err != nil && err.Error() == "EOF" {
				return nil, nil
			} else if err != nil {
				if cerr := util.CliError(err); cerr.ExitCode() != 128 {
					return nil, cerr
				}
				return nil, fmt.Errorf("error receiving from %s.%s: %v", service, endpoint, err)
			}
			fmt.Print(string(rsp.Data))
		} else {
			var response map[string]interface{}
			if err := stream.Recv(&response); err != nil && err.Error() == "EOF" {
				return nil, nil
			} else if err != nil {
				if cerr := util.CliError(err); cerr.ExitCode() != 128 {
					return nil, cerr
				}
				return nil, fmt.Errorf("error receiving from %s.%s: %v", service, endpoint, err)
			}
			b, _ := json.MarshalIndent(response, "", "\t")
			fmt.Print(string(b))
		}
	}
}

func publish(c *cli.Context, args []string) ([]byte, error) {
	if err := Publish(c, args); err != nil {
		return nil, err
	}
	return []byte(`ok`), nil
}

func queryHealth(c *cli.Context, args []string) ([]byte, error) {
	return QueryHealth(c, args)
}

func queryStats(c *cli.Context, args []string) ([]byte, error) {
	return QueryStats(c, args)
}
