package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/micro/micro/v3/client/cli/util"
	cliutil "github.com/micro/micro/v3/client/cli/util"
	cbytes "github.com/micro/micro/v3/internal/codec/bytes"
	clic "github.com/micro/micro/v3/internal/command"
	"github.com/micro/micro/v3/service/client"
	"github.com/urfave/cli/v2"
)

func listServices(c *cli.Context, args []string) ([]byte, error) {
	return clic.ListServices(c)
}

func callService(c *cli.Context, args []string) ([]byte, error) {
	return clic.CallService(c, args)
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
	return nil, cliutil.SetEnv(args[0])
}

func listEnvs(c *cli.Context, args []string) ([]byte, error) {
	envs, err := cliutil.GetEnvs()
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

	return nil, cliutil.AddEnv(cliutil.Env{
		Name:         args[0],
		ProxyAddress: args[1],
	})
}

func delEnv(c *cli.Context, args []string) ([]byte, error) {
	if len(args) == 0 {
		return nil, cli.ShowSubcommandHelp(c)
	}
	return nil, cliutil.DelEnv(args[0])
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

	req := client.DefaultClient.NewRequest(service, endpoint, request, client.WithContentType("application/json"))
	stream, err := client.DefaultClient.Stream(context.Background(), req)
	if err != nil {
		if cerr := cliutil.CliError(err); cerr.ExitCode() != 128 {
			return nil, cerr
		}
		return nil, fmt.Errorf("error calling %s.%s: %v", service, endpoint, err)
	}

	if err := stream.Send(request); err != nil {
		if cerr := cliutil.CliError(err); cerr.ExitCode() != 128 {
			return nil, cerr
		}
		return nil, fmt.Errorf("error sending to %s.%s: %v", service, endpoint, err)
	}

	output := c.String("output")

	for {
		if output == "raw" {
			rsp := cbytes.Frame{}
			if err := stream.Recv(&rsp); err != nil {
				if cerr := cliutil.CliError(err); cerr.ExitCode() != 128 {
					return nil, cerr
				}
				return nil, fmt.Errorf("error receiving from %s.%s: %v", service, endpoint, err)
			}
			fmt.Print(string(rsp.Data))
		} else {
			var response map[string]interface{}
			if err := stream.Recv(&response); err != nil {
				if cerr := cliutil.CliError(err); cerr.ExitCode() != 128 {
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
	if err := clic.Publish(c, args); err != nil {
		return nil, err
	}
	return []byte(`ok`), nil
}

func queryHealth(c *cli.Context, args []string) ([]byte, error) {
	return clic.QueryHealth(c, args)
}

func queryStats(c *cli.Context, args []string) ([]byte, error) {
	return QueryStats(c, args)
}
