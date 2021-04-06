package cli

import (
	"bytes"
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
	clic "github.com/micro/micro/v3/internal/command"
	proto "github.com/micro/micro/v3/proto/debug"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/registry"
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

		sl, err := clic.ListServices(c)
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
		sl, err := clic.ListServices(c)
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
