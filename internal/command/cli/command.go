package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/client"
	cbytes "github.com/micro/go-micro/v2/codec/bytes"
	proto "github.com/micro/go-micro/v2/debug/service/proto"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/service"
	"github.com/micro/micro/v2/client/cli/namespace"
	"github.com/micro/micro/v2/client/cli/util"
	"github.com/micro/micro/v2/cmd"
	inclient "github.com/micro/micro/v2/internal/client"
	dns "github.com/micro/micro/v2/service/network/dns/proto/dns"

	"github.com/olekukonko/tablewriter"
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

func getPeers(v map[string]interface{}) map[string]string {
	if v == nil {
		return nil
	}

	peers := make(map[string]string)
	node := v["node"].(map[string]interface{})
	peers[node["id"].(string)] = node["address"].(string)

	// return peers if nil
	if v["peers"] == nil {
		return peers
	}

	nodes := v["peers"].([]interface{})

	for _, peer := range nodes {
		p := getPeers(peer.(map[string]interface{}))
		for id, address := range p {
			peers[id] = address
		}
	}

	return peers
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

func RegisterService(c *cli.Context, args []string) ([]byte, error) {
	if len(args) == 0 {
		return nil, errors.New("require service definition")
	}

	req := strings.Join(args, " ")

	var srv *registry.Service

	d := json.NewDecoder(strings.NewReader(req))
	d.UseNumber()

	if err := d.Decode(&srv); err != nil {
		return nil, err
	}

	reg := *cmd.DefaultOptions().Registry
	reg.Init(service.WithClient(inclient.New(c)))
	if err := reg.Register(srv); err != nil {
		return nil, err
	}

	return []byte("ok"), nil
}

func DeregisterService(c *cli.Context, args []string) ([]byte, error) {
	if len(args) == 0 {
		return nil, errors.New("require service definition")
	}

	req := strings.Join(args, " ")

	var srv *registry.Service

	d := json.NewDecoder(strings.NewReader(req))
	d.UseNumber()

	if err := d.Decode(&srv); err != nil {
		return nil, err
	}

	reg := *cmd.DefaultOptions().Registry
	reg.Init(service.WithClient(inclient.New(c)))
	if err := reg.Deregister(srv); err != nil {
		return nil, err
	}

	return []byte("ok"), nil
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

	reg := *cmd.DefaultOptions().Registry
	reg.Init(service.WithClient(inclient.New(c)))
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

func NetworkConnect(c *cli.Context, args []string) ([]byte, error) {
	if len(args) == 0 {
		return nil, nil
	}

	cli := *cmd.DefaultOptions().Client

	request := map[string]interface{}{
		"nodes": []interface{}{
			map[string]interface{}{
				"address": args[0],
			},
		},
	}

	var rsp map[string]interface{}

	req := cli.NewRequest("go.micro.network", "Network.Connect", request, client.WithContentType("application/json"))
	err := cli.Call(context.TODO(), req, &rsp)
	if err != nil {
		return nil, err
	}

	b, _ := json.MarshalIndent(rsp, "", "\t")
	return b, nil
}

func NetworkConnections(c *cli.Context) ([]byte, error) {
	cli := *cmd.DefaultOptions().Client

	request := map[string]interface{}{
		"depth": 1,
	}

	var rsp map[string]interface{}

	req := cli.NewRequest("go.micro.network", "Network.Graph", request, client.WithContentType("application/json"))
	err := cli.Call(context.TODO(), req, &rsp)
	if err != nil {
		return nil, err
	}

	if rsp["root"] == nil {
		return nil, nil
	}

	peers := rsp["root"].(map[string]interface{})["peers"]

	if peers == nil {
		return nil, nil
	}

	b := bytes.NewBuffer(nil)
	table := tablewriter.NewWriter(b)
	table.SetHeader([]string{"NODE", "ADDRESS"})

	// root node
	for _, n := range peers.([]interface{}) {
		node := n.(map[string]interface{})["node"].(map[string]interface{})
		strEntry := []string{
			fmt.Sprintf("%s", node["id"]),
			fmt.Sprintf("%s", node["address"]),
		}
		table.Append(strEntry)
	}

	// render table into b
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()

	return b.Bytes(), nil
}

func NetworkGraph(c *cli.Context) ([]byte, error) {
	cli := *cmd.DefaultOptions().Client

	var rsp map[string]interface{}

	req := cli.NewRequest("go.micro.network", "Network.Graph", map[string]interface{}{}, client.WithContentType("application/json"))
	err := cli.Call(context.TODO(), req, &rsp)
	if err != nil {
		return nil, err
	}

	b, _ := json.MarshalIndent(rsp, "", "\t")
	return b, nil
}

func NetworkNodes(c *cli.Context) ([]byte, error) {
	cli := *cmd.DefaultOptions().Client

	var rsp map[string]interface{}

	// TODO: change to list nodes
	req := cli.NewRequest("go.micro.network", "Network.Nodes", map[string]interface{}{}, client.WithContentType("application/json"))
	err := cli.Call(context.TODO(), req, &rsp)
	if err != nil {
		return nil, err
	}

	// return if nil
	if rsp["nodes"] == nil {
		return nil, nil
	}

	b := bytes.NewBuffer(nil)
	table := tablewriter.NewWriter(b)
	table.SetHeader([]string{"ID", "ADDRESS"})

	// get nodes

	if rsp["nodes"] != nil {
		// root node
		for _, n := range rsp["nodes"].([]interface{}) {
			node := n.(map[string]interface{})
			strEntry := []string{
				fmt.Sprintf("%s", node["id"]),
				fmt.Sprintf("%s", node["address"]),
			}
			table.Append(strEntry)
		}
	}

	// render table into b
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()

	return b.Bytes(), nil
}

func NetworkRoutes(c *cli.Context) ([]byte, error) {
	cli := (*cmd.DefaultOptions().Client)

	query := map[string]string{}

	for _, filter := range []string{"service", "address", "gateway", "router", "network"} {
		if v := c.String(filter); len(v) > 0 {
			query[filter] = v
		}
	}

	request := map[string]interface{}{
		"query": query,
	}

	var rsp map[string]interface{}

	req := cli.NewRequest("go.micro.network", "Network.Routes", request, client.WithContentType("application/json"))
	err := cli.Call(context.TODO(), req, &rsp)
	if err != nil {
		return nil, err
	}

	if len(rsp) == 0 {
		return []byte(``), nil
	}

	b := bytes.NewBuffer(nil)
	table := tablewriter.NewWriter(b)
	table.SetHeader([]string{"SERVICE", "ADDRESS", "GATEWAY", "ROUTER", "NETWORK", "METRIC", "LINK"})

	routes := rsp["routes"].([]interface{})

	val := func(v interface{}) string {
		if v == nil {
			return ""
		}
		return v.(string)
	}

	var sortedRoutes [][]string

	for _, r := range routes {
		route := r.(map[string]interface{})
		service := route["service"]
		address := route["address"]
		gateway := val(route["gateway"])
		router := route["router"]
		network := route["network"]
		link := route["link"]
		metric := route["metric"]

		var metInt int64
		if metric != nil {
			metInt, _ = strconv.ParseInt(route["metric"].(string), 10, 64)
		}

		// set max int64 metric to infinity
		if metInt == math.MaxInt64 {
			metric = "∞"
		} else {
			metric = fmt.Sprintf("%d", metInt)
		}

		sortedRoutes = append(sortedRoutes, []string{
			fmt.Sprintf("%s", service),
			fmt.Sprintf("%s", address),
			fmt.Sprintf("%s", gateway),
			fmt.Sprintf("%s", router),
			fmt.Sprintf("%s", network),
			fmt.Sprintf("%s", metric),
			fmt.Sprintf("%s", link),
		})
	}

	sort.Slice(sortedRoutes, func(i, j int) bool { return sortedRoutes[i][0] < sortedRoutes[j][0] })

	table.AppendBulk(sortedRoutes)
	// render table into b
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()

	return b.Bytes(), nil
}

func NetworkServices(c *cli.Context) ([]byte, error) {
	cli := (*cmd.DefaultOptions().Client)

	var rsp map[string]interface{}

	req := cli.NewRequest("go.micro.network", "Network.Services", map[string]interface{}{}, client.WithContentType("application/json"))
	err := cli.Call(context.TODO(), req, &rsp)
	if err != nil {
		return nil, err
	}

	if len(rsp) == 0 || rsp["services"] == nil {
		return []byte(``), nil
	}

	rspSrv := rsp["services"].([]interface{})

	var services []string

	for _, service := range rspSrv {
		services = append(services, service.(string))
	}

	sort.Strings(services)

	return []byte(strings.Join(services, "\n")), nil
}

func NetworkDNSAdvertise(c *cli.Context) ([]byte, error) {
	err := networkDNSHelper("Dns.Advertise", c.String("address"), c.String("domain"), c.String("token"))
	if err != nil {
		return []byte(``), err
	}
	return []byte("Registered " + c.String("domain") + ": " + c.String("address")), nil
}

func NetworkDNSRemove(c *cli.Context) ([]byte, error) {
	err := networkDNSHelper("Dns.Remove", c.String("address"), c.String("domain"), c.String("token"))
	if err != nil {
		return []byte(``), err
	}
	return []byte("Removed " + c.String("domain") + ": " + c.String("address")), nil
}

func NetworkDNSResolve(c *cli.Context) ([]byte, error) {
	request := make(map[string]interface{})
	request["name"] = c.String("domain")
	request["type"] = c.String("type")

	cli := (*cmd.DefaultOptions().Client)
	req := cli.NewRequest("go.micro.network.dns", "Dns.Resolve", request, client.WithContentType("application/json"))
	var rsp map[string][]*dns.Record
	err := cli.Call(
		metadata.NewContext(context.Background(), map[string]string{
			"Authorization": "Bearer " + c.String("token"),
		}),
		req,
		&rsp,
		client.WithRetries(3),
	)
	if err != nil {
		return []byte(``), err
	}

	rawRecords, ok := rsp["records"]
	if !ok {
		return []byte(``), errors.New("Response did not contain any records")
	}
	var resolved []string
	for _, r := range rawRecords {
		resolved = append(resolved, r.Value)
	}

	return []byte(strings.Join(resolved, "\n")), nil
}

func networkDNSHelper(action, address, domain, token string) error {
	request := map[string]interface{}{
		"records": []*dns.Record{},
	}

	if strings.Count(address, ":") > 1 {
		request["records"] = []*dns.Record{
			&dns.Record{
				Type:  "AAAA",
				Name:  domain,
				Value: address,
				Ttl:   1,
			},
		}
	} else {
		request["records"] = []*dns.Record{
			&dns.Record{
				Type:  "A",
				Name:  domain,
				Value: address,
				Ttl:   1,
			},
		}
	}

	cli := (*cmd.DefaultOptions().Client)
	req := cli.NewRequest("go.micro.network.dns", action, request, client.WithContentType("application/json"))
	var rsp map[string]interface{}
	err := cli.Call(
		metadata.NewContext(context.Background(), map[string]string{
			"Authorization": "Bearer " + token,
		}),
		req,
		&rsp,
		client.WithRetries(3),
	)
	if err != nil {
		return err
	}
	return nil
}

func ListServices(c *cli.Context) ([]byte, error) {
	var rsp []*registry.Service
	var err error

	ns, err := namespace.Get(util.GetEnv(c).Name)
	if err != nil {
		return nil, err
	}

	reg := *cmd.DefaultOptions().Registry
	reg.Init(service.WithClient(inclient.New(c)))
	rsp, err = reg.ListServices(registry.ListDomain(ns))
	if err != nil {
		return nil, err
	}

	sort.Sort(sortedServices{rsp})

	var services []string

	for _, service := range rsp {
		services = append(services, service.Name)
	}

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

	cl := *cmd.DefaultOptions().Client
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
	cli := inclient.New(c)

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

	req := (*cmd.DefaultOptions().Client).NewRequest(args[0], "Debug.Health", &proto.HealthRequest{})

	// if the address is specified then we just call it
	if addr := c.String("address"); len(addr) > 0 {
		rsp := &proto.HealthResponse{}
		err := (*cmd.DefaultOptions().Client).Call(
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
	reg := *cmd.DefaultOptions().Registry
	reg.Init(service.WithClient(inclient.New(c)))
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
			err = (*cmd.DefaultOptions().Client).Call(
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

func QueryStats(c *cli.Context, args []string) ([]byte, error) {
	if len(args) == 0 {
		return nil, errors.New("require service name")
	}

	ns, err := namespace.Get(util.GetEnv(c).Name)
	if err != nil {
		return nil, err
	}

	reg := *cmd.DefaultOptions().Registry
	reg.Init(service.WithClient(inclient.New(c)))
	service, err := reg.GetService(args[0], registry.GetDomain(ns))
	if err != nil {
		return nil, err
	}

	if len(service) == 0 {
		return nil, errors.New("Service not found")
	}

	req := (*cmd.DefaultOptions().Client).NewRequest(service[0].Name, "Debug.Stats", &proto.StatsRequest{})

	var output []string

	// print things
	output = append(output, "service  "+service[0].Name)

	for _, serv := range service {
		// print things
		output = append(output, "\nversion "+serv.Version)
		output = append(output, "\nnode\t\taddress:port\t\tstarted\tuptime\tmemory\tthreads\tgc")

		// query health for every node
		for _, node := range serv.Nodes {
			address := node.Address
			rsp := &proto.StatsResponse{}

			var err error

			// call using client
			err = (*cmd.DefaultOptions().Client).Call(
				context.Background(),
				req,
				rsp,
				client.WithAddress(address),
			)

			var started, uptime, memory, gc string
			if err == nil {
				started = time.Unix(int64(rsp.Started), 0).Format("Jan 2 15:04:05")
				uptime = fmt.Sprintf("%v", time.Duration(rsp.Uptime)*time.Second)
				memory = fmt.Sprintf("%.2fmb", float64(rsp.Memory)/(1024.0*1024.0))
				gc = fmt.Sprintf("%v", time.Duration(rsp.Gc))
			}

			line := fmt.Sprintf("%s\t\t%s\t\t%s\t%s\t%s\t%d\t%s",
				node.Id, node.Address, started, uptime, memory, rsp.Threads, gc)

			output = append(output, line)
		}
	}

	return []byte(strings.Join(output, "\n")), nil
}
