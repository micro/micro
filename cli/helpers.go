package cli

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/micro/internal/command"

	"golang.org/x/net/context"
)

func listServices(c *cli.Context) {
	rsp, err := command.ListServices(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(rsp))
}

func registerService(c *cli.Context) {
	rsp, err := command.RegisterService(c, c.Args())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(rsp))
}

func deregisterService(c *cli.Context) {
	rsp, err := command.DeregisterService(c, c.Args())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(rsp))
}

func getService(c *cli.Context) {
	rsp, err := command.GetService(c, c.Args())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(rsp))
}

func queryService(c *cli.Context) {
	rsp, err := command.QueryService(c, c.Args())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(rsp))
}

// TODO: stream via HTTP
func streamService(c *cli.Context) {
	if len(c.Args()) < 2 {
		fmt.Println("require service and method")
		return
	}
	service := c.Args()[0]
	method := c.Args()[1]
	var request map[string]interface{}
	json.Unmarshal([]byte(strings.Join(c.Args()[2:], " ")), &request)
	req := (*cmd.DefaultOptions().Client).NewJsonRequest(service, method, request)
	stream, err := (*cmd.DefaultOptions().Client).Stream(context.Background(), req)
	if err != nil {
		fmt.Printf("error calling %s.%s: %v\n", service, method, err)
		return
	}

	if err := stream.Send(request); err != nil {
		fmt.Printf("error sending to %s.%s: %v\n", service, method, err)
		return
	}

	for {
		var response map[string]interface{}
		if err := stream.Recv(&response); err != nil {
			fmt.Printf("error receiving from %s.%s: %v\n", service, method, err)
			return
		}

		b, _ := json.MarshalIndent(response, "", "\t")
		fmt.Println(string(b))

		// artificial delay
		time.Sleep(time.Millisecond * 10)
	}
}

func queryHealth(c *cli.Context) {
	rsp, err := command.QueryHealth(c, c.Args())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(rsp))
}

func queryStats(c *cli.Context) {
	rsp, err := command.QueryStats(c, c.Args())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(rsp))
}
