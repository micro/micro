package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/micro/micro/v3/client/cli/namespace"
	"github.com/micro/micro/v3/client/cli/util"
	proto "github.com/micro/micro/v3/proto/debug"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/registry"
	goregistry "github.com/micro/micro/v3/service/registry"
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

func QueryStats(c *cli.Context, args []string) ([]byte, error) {
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

	service, err := registry.DefaultRegistry.GetService(args[0], goregistry.GetDomain(ns))
	if err != nil {
		return nil, err
	}

	if len(service) == 0 {
		return nil, errors.New("Service not found")
	}

	req := client.NewRequest(service[0].Name, "Debug.Stats", &proto.StatsRequest{})

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
			err = client.DefaultClient.Call(context.Background(), req, rsp, client.WithAddress(address))

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
