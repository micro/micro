package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/client"
	mcmd "github.com/micro/go-micro/v2/cmd"
	ccli "github.com/micro/micro/v2/client/cli"
	"github.com/micro/micro/v2/cmd"
	clic "github.com/micro/micro/v2/internal/command/cli"
	"github.com/olekukonko/tablewriter"
)

func init() {
	cmd.Register(&cli.Command{
		Name: "network",
		Subcommands: []*cli.Command{
			{
				Name:   "connect",
				Usage:  "connect to the network. specify nodes e.g connect ip:port",
				Action: ccli.Print(networkConnect),
			},
			{
				Name:   "connections",
				Usage:  "List the immediate connections to the network",
				Action: ccli.Print(networkConnections),
			},
			{
				Name:   "graph",
				Usage:  "Get the network graph",
				Action: ccli.Print(networkGraph),
			},
			{
				Name:   "nodes",
				Usage:  "List nodes in the network",
				Action: ccli.Print(networkNodes),
			},
			{
				Name:   "routes",
				Usage:  "List network routes",
				Action: ccli.Print(networkRoutes),
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "service",
						Usage: "Filter by service",
					},
					&cli.StringFlag{
						Name:  "address",
						Usage: "Filter by address",
					},
					&cli.StringFlag{
						Name:  "gateway",
						Usage: "Filter by gateway",
					},
					&cli.StringFlag{
						Name:  "router",
						Usage: "Filter by router",
					},
					&cli.StringFlag{
						Name:  "network",
						Usage: "Filter by network",
					},
				},
			},
			{
				Name:   "services",
				Usage:  "Get the network services",
				Action: ccli.Print(networkServices),
			},
			// TODO: duplicates call. Move so we reuse same stuff.
			{
				Name:   "call",
				Usage:  "Call a service e.g micro call greeter Say.Hello '{\"name\": \"John\"}",
				Action: ccli.Print(netCall),
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "address",
						Usage:   "Set the address of the service instance to call",
						EnvVars: []string{"MICRO_ADDRESS"},
					},
					&cli.StringFlag{
						Name:    "output, o",
						Usage:   "Set the output format; json (default), raw",
						EnvVars: []string{"MICRO_OUTPUT"},
					},
					&cli.StringSliceFlag{
						Name:    "metadata",
						Usage:   "A list of key-value pairs to be forwarded as metadata",
						EnvVars: []string{"MICRO_METADATA"},
					},
				},
			},
		},
	})
}

func networkConnect(c *cli.Context, args []string) ([]byte, error) {
	if len(args) == 0 {
		return nil, nil
	}

	cli := *mcmd.DefaultCmd.Options().Client

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

func networkConnections(c *cli.Context, args []string) ([]byte, error) {
	cli := *mcmd.DefaultCmd.Options().Client

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

func networkGraph(c *cli.Context, args []string) ([]byte, error) {
	cli := *mcmd.DefaultCmd.Options().Client

	var rsp map[string]interface{}

	req := cli.NewRequest("go.micro.network", "Network.Graph", map[string]interface{}{}, client.WithContentType("application/json"))
	err := cli.Call(context.TODO(), req, &rsp)
	if err != nil {
		return nil, err
	}

	b, _ := json.MarshalIndent(rsp, "", "\t")
	return b, nil
}

func networkNodes(c *cli.Context, args []string) ([]byte, error) {
	cli := *mcmd.DefaultCmd.Options().Client

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

func networkRoutes(c *cli.Context, args []string) ([]byte, error) {
	cli := (*mcmd.DefaultCmd.Options().Client)

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

func networkServices(c *cli.Context, args []string) ([]byte, error) {
	cli := (*mcmd.DefaultCmd.Options().Client)

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

// netCall calls services through the network
func netCall(c *cli.Context, args []string) ([]byte, error) {
	os.Setenv("MICRO_PROXY", "go.micro.network")
	return clic.CallService(c, args)
}
