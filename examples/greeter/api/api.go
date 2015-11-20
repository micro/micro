package main

import (
	"code.google.com/p/go.net/context"
	"encoding/json"
	"strings"

	log "github.com/golang/glog"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/server"
	api "github.com/micro/micro/api/proto"
	hello "github.com/micro/micro/examples/greeter/server/proto/hello"
)

type Say struct{}

func (s *Say) Hello(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Info("Received Say.Hello API request")

	name, ok := req.Get["name"]
	if !ok || len(name.Values) == 0 {
		return errors.BadRequest("go.micro.api.greeter", "Name cannot be blank")
	}

	request := client.NewRequest("go.micro.srv.greeter", "Say.Hello", &hello.Request{
		Name: strings.Join(name.Values, " "),
	})

	response := &hello.Response{}

	if err := client.Call(ctx, request, response); err != nil {
		return err
	}

	rsp.StatusCode = 200
	b, _ := json.Marshal(map[string]string{
		"message": response.Msg,
	})
	rsp.Body = string(b)

	return nil
}

func main() {
	// optionally setup command line usage
	cmd.Init()

	// Initialise Server
	server.Init(
		server.Name("go.micro.api.greeter"),
	)

	// Register Handlers
	server.Handle(
		server.NewHandler(
			new(Say),
		),
	)

	// Run server
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
