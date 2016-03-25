package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/registry"
	"github.com/micro/micro/internal/command"

	"golang.org/x/net/context"
)

func post(url string, b []byte, v interface{}) error {
	if !strings.HasPrefix(url, "http") && !strings.HasPrefix(url, "https") {
		url = "http://" + url
	}

	buf := bytes.NewBuffer(b)
	defer buf.Reset()

	rsp, err := http.Post(url, "application/json", buf)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	bu, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	if v == nil {
		return nil
	}

	return json.Unmarshal(bu, v)
}

func del(url string, b []byte, v interface{}) error {
	if !strings.HasPrefix(url, "http") && !strings.HasPrefix(url, "https") {
		url = "http://" + url
	}

	buf := bytes.NewBuffer(b)
	defer buf.Reset()

	req, err := http.NewRequest("DELETE", url, buf)
	if err != nil {
		return err
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	bu, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	if v == nil {
		return nil
	}

	return json.Unmarshal(bu, v)
}

func listServices(c *cli.Context) {
	rsp, err := command.ListServices(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(rsp))
}

func registerService(c *cli.Context) {
	if len(c.Args()) != 1 {
		fmt.Println("require service definition")
		return
	}

	if p := c.GlobalString("proxy_address"); len(p) > 0 {
		if err := post(p+"/registry", []byte(c.Args().First()), nil); err != nil {
			fmt.Println(err.Error())
		}
		return
	}

	var service *registry.Service

	if err := json.Unmarshal([]byte(c.Args().First()), &service); err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := (*cmd.DefaultOptions().Registry).Register(service); err != nil {
		fmt.Println(err.Error())
	}
}

func deregisterService(c *cli.Context) {
	if len(c.Args()) != 1 {
		fmt.Println("require service definition")
		return
	}

	if p := c.GlobalString("proxy_address"); len(p) > 0 {
		if err := del(p+"/registry", []byte(c.Args().First()), nil); err != nil {
			fmt.Println(err.Error())
		}
		return
	}

	var service *registry.Service
	if err := json.Unmarshal([]byte(c.Args().First()), &service); err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := (*cmd.DefaultOptions().Registry).Deregister(service); err != nil {
		fmt.Println(err.Error())
		return
	}
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
