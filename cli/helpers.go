package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
	clic "github.com/micro/micro/internal/command/cli"
)

type helper func(*cli.Context, []string)

func printer(h helper) func(*cli.Context) {
	return func(c *cli.Context) {
		h(c, c.Args())
		fmt.Printf("\n")
	}
}

func listServices(c *cli.Context, args []string) {
	rsp, err := clic.ListServices(c)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Print(string(rsp))
}

func registerService(c *cli.Context, args []string) {
	rsp, err := clic.RegisterService(c, args)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Print(string(rsp))
}

func deregisterService(c *cli.Context, args []string) {
	rsp, err := clic.DeregisterService(c, args)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Print(string(rsp))
}

func getService(c *cli.Context, args []string) {
	rsp, err := clic.GetService(c, args)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Print(string(rsp))
}

func callService(c *cli.Context, args []string) {
	rsp, err := clic.CallService(c, args)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Print(string(rsp))
}

// TODO: stream via HTTP
func streamService(c *cli.Context, args []string) {
	if len(args) < 2 {
		fmt.Print("require service and method")
		return
	}
	service := args[0]
	method := args[1]
	var request map[string]interface{}
	json.Unmarshal([]byte(strings.Join(args[2:], " ")), &request)
	req := (*cmd.DefaultOptions().Client).NewRequest(service, method, request, client.WithContentType("application/json"))
	stream, err := (*cmd.DefaultOptions().Client).Stream(context.Background(), req)
	if err != nil {
		fmt.Printf("error calling %s.%s: %v", service, method, err)
		return
	}

	if err := stream.Send(request); err != nil {
		fmt.Printf("error sending to %s.%s: %v", service, method, err)
		return
	}

	for {
		var response map[string]interface{}
		if err := stream.Recv(&response); err != nil {
			fmt.Printf("error receiving from %s.%s: %v", service, method, err)
			return
		}

		b, _ := json.MarshalIndent(response, "", "\t")
		fmt.Print(string(b))

		// artificial delay
		time.Sleep(time.Millisecond * 10)
	}
}

func publish(c *cli.Context, args []string) {
	err := clic.Publish(c, args)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Print("ok")
}
func queryHealth(c *cli.Context, args []string) {
	rsp, err := clic.QueryHealth(c, args)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Print(string(rsp))
}

func queryStats(c *cli.Context, args []string) {
	rsp, err := clic.QueryStats(c, args)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Print(string(rsp))
}
