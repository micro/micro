package debug

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/micro/micro/v3/client/cli/namespace"
	"github.com/micro/micro/v3/client/cli/util"
	"github.com/micro/micro/v3/cmd"
	proto "github.com/micro/micro/v3/proto/debug"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/registry"
	"github.com/urfave/cli/v2"
)

func init() {
	subcommands := []*cli.Command{
		&cli.Command{
			Name:   "health",
			Usage:  `Get the service health`,
			Action: util.Print(QueryHealth),
		},
		&cli.Command{
			Name:   "stats",
			Usage:  "Query the stats of specified service(s), e.g micro stats srv1 srv2 srv3",
			Action: util.Print(queryStats),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "all",
					Usage: "to list all builtin services use --all builtin, for user's services use --all custom",
				},
			},
		},
	}

	command := &cli.Command{
		Name:        "debug",
		Usage:       "Debug a service",
		Action:      func(ctx *cli.Context) error { return nil },
		Subcommands: subcommands,
	}

	cmd.Register(command)
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

func queryHealth(c *cli.Context, args []string) ([]byte, error) {
	return QueryHealth(c, args)
}

func queryStats(c *cli.Context, args []string) ([]byte, error) {
	return QueryStats(c, args)
}
