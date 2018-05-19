package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
	clic "github.com/micro/micro/internal/command/cli"
)

func die(err error) {
	if err == nil {
		return
	}

	fmt.Print(err)
	os.Exit(1)
}

type helper func(*cli.Context, []string)

func printer(h helper) func(*cli.Context) {
	return func(c *cli.Context) {
		h(c, c.Args())
		fmt.Printf("\n")
	}
}

func listServices(c *cli.Context, args []string) {
	rsp, err := clic.ListServices(c)
	die(err)
	fmt.Print(string(rsp))
}

func registerService(c *cli.Context, args []string) {
	rsp, err := clic.RegisterService(c, args)
	die(err)
	fmt.Print(string(rsp))
}

func deregisterService(c *cli.Context, args []string) {
	rsp, err := clic.DeregisterService(c, args)
	die(err)
	fmt.Print(string(rsp))
}

func getService(c *cli.Context, args []string) {
	rsp, err := clic.GetService(c, args)
	die(err)
	fmt.Print(string(rsp))
}

func callService(c *cli.Context, args []string) {
	rsp, err := clic.CallService(c, args)
	die(err)
	fmt.Print(string(rsp))
}

// TODO: stream via HTTP
func streamService(c *cli.Context, args []string) {
	if len(args) < 2 {
		fmt.Print("require service and method")
		os.Exit(1)
	}
	service := args[0]
	method := args[1]
	var request map[string]interface{}
	json.Unmarshal([]byte(strings.Join(args[2:], " ")), &request)
	req := (*cmd.DefaultOptions().Client).NewRequest(service, method, request, client.WithContentType("application/json"))
	stream, err := (*cmd.DefaultOptions().Client).Stream(context.Background(), req)
	if err != nil {
		fmt.Printf("error calling %s.%s: %v", service, method, err)
		os.Exit(1)
	}

	if err := stream.Send(request); err != nil {
		fmt.Printf("error sending to %s.%s: %v", service, method, err)
		os.Exit(1)
	}

	for {
		var response map[string]interface{}
		if err := stream.Recv(&response); err != nil {
			fmt.Printf("error receiving from %s.%s: %v", service, method, err)
			os.Exit(1)
		}

		b, _ := json.MarshalIndent(response, "", "\t")
		fmt.Print(string(b))

		// artificial delay
		time.Sleep(time.Millisecond * 10)
	}
}

func publish(c *cli.Context, args []string) {
	err := clic.Publish(c, args)
	die(err)
	fmt.Print("ok")
}
func queryHealth(c *cli.Context, args []string) {
	rsp, err := clic.QueryHealth(c, args)
	die(err)
	fmt.Print(string(rsp))
}

func queryStats(c *cli.Context, args []string) {
	rsp, err := clic.QueryStats(c, args)
	die(err)
	fmt.Print(string(rsp))
}
