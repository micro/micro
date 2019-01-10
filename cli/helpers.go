package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
	clic "github.com/micro/micro/internal/command/cli"
)

type exec func(*cli.Context, []string) ([]byte, error)

func printer(e exec) func(*cli.Context) {
	return func(c *cli.Context) {
		rsp, err := e(c, c.Args())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", string(rsp))
	}
}

func listServices(c *cli.Context, args []string) ([]byte, error) {
	return clic.ListServices(c)
}

func registerService(c *cli.Context, args []string) ([]byte, error) {
	return clic.RegisterService(c, args)
}

func deregisterService(c *cli.Context, args []string) ([]byte, error) {
	return clic.DeregisterService(c, args)
}

func getService(c *cli.Context, args []string) ([]byte, error) {
	return clic.GetService(c, args)
}

func callService(c *cli.Context, args []string) ([]byte, error) {
	return clic.CallService(c, args)
}

// TODO: stream via HTTP
func streamService(c *cli.Context, args []string) ([]byte, error) {
	if len(args) < 2 {
		return nil, errors.New("require service and endpoint")
	}
	service := args[0]
	endpoint := args[1]
	var request map[string]interface{}
	err := json.Unmarshal([]byte(strings.Join(args[2:], " ")), &request)
	if err != nil {
		return nil, err
	}
	req := (*cmd.DefaultOptions().Client).NewRequest(service, endpoint, request, client.WithContentType("application/json"))
	stream, err := (*cmd.DefaultOptions().Client).Stream(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("error calling %s.%s: %v", service, endpoint, err)
	}

	if err := stream.Send(request); err != nil {
		return nil, fmt.Errorf("error sending to %s.%s: %v", service, endpoint, err)
	}

	for {
		var response map[string]interface{}
		if err := stream.Recv(&response); err != nil {
			return nil, fmt.Errorf("error receiving from %s.%s: %v", service, endpoint, err)
		}

		b, _ := json.MarshalIndent(response, "", "\t")
		fmt.Print(string(b))

		// artificial delay
		time.Sleep(time.Millisecond * 10)
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
	return clic.QueryStats(c, args)
}
