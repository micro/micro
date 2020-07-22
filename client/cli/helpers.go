package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/cmd"
	cbytes "github.com/micro/go-micro/v2/codec/bytes"
	cliutil "github.com/micro/micro/v2/client/cli/util"
	clic "github.com/micro/micro/v2/internal/command/cli"
)

func listServices(c *cli.Context, args []string) ([]byte, error) {
	return clic.ListServices(c)
}

func callService(c *cli.Context, args []string) ([]byte, error) {
	return clic.CallService(c, args)
}

func getEnv(c *cli.Context, args []string) ([]byte, error) {
	env := cliutil.GetEnv(c)
	return []byte(env.Name), nil
}

func setEnv(c *cli.Context, args []string) ([]byte, error) {
	cliutil.SetEnv(args[0])
	return nil, nil
}

func listEnvs(c *cli.Context, args []string) ([]byte, error) {
	envs := cliutil.GetEnvs()
	current := cliutil.GetEnv(c)

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
		fmt.Fprintf(w, "%v %v \t %v", prefix, env.Name, env.ProxyAddress)
	}
	w.Flush()
	return byt.Bytes(), nil
}

func addEnv(c *cli.Context, args []string) ([]byte, error) {
	if len(args) == 0 {
		return nil, errors.New("name required")
	}
	if len(args) == 1 {
		args = append(args, "") // default to no proxy address
	}

	cliutil.AddEnv(cliutil.Env{
		Name:         args[0],
		ProxyAddress: args[1],
	})
	return nil, nil
}

func delEnv(c *cli.Context, args []string) ([]byte, error) {
	if len(args) == 0 {
		return nil, errors.New("name required")
	}
	cliutil.DelEnv(args[0])
	return nil, nil
}

// TODO: stream via HTTP
func streamService(c *cli.Context, args []string) ([]byte, error) {
	if len(args) < 2 {
		return nil, errors.New("require service and endpoint")
	}
	service := args[0]
	endpoint := args[1]
	var request map[string]interface{}

	// ignore error
	json.Unmarshal([]byte(strings.Join(args[2:], " ")), &request)

	req := (*cmd.DefaultCmd.Options().Client).NewRequest(service, endpoint, request, client.WithContentType("application/json"))
	stream, err := (*cmd.DefaultCmd.Options().Client).Stream(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("error calling %s.%s: %v", service, endpoint, err)
	}

	if err := stream.Send(request); err != nil {
		return nil, fmt.Errorf("error sending to %s.%s: %v", service, endpoint, err)
	}

	output := c.String("output")

	for {
		if output == "raw" {
			rsp := cbytes.Frame{}
			if err := stream.Recv(&rsp); err != nil {
				return nil, fmt.Errorf("error receiving from %s.%s: %v", service, endpoint, err)
			}
			fmt.Print(string(rsp.Data))
		} else {
			var response map[string]interface{}
			if err := stream.Recv(&response); err != nil {
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
