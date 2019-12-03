// Package web provides a dashboard for debugging and introspection of go-micro services
package web

import (
	"net/http"

	"github.com/micro/cli"
	//"github.com/micro/go-micro/config/cmd"
	//pb "github.com/micro/go-micro/network/proto"
	"github.com/micro/go-micro/web"
)

// Run starts go.micro.web.debug
func Run(ctx *cli.Context) {
	//c := *cmd.DefaultOptions().Client
	//client := pb.NewNetworkService("go.micro.network", c)

	opts := []web.Option{
		web.Name("go.micro.web.debug"),
	}

	address := ctx.GlobalString("server_address")
	if len(address) > 0 {
		opts = append(opts, web.Address(address))
	}

	service := web.NewService(opts...)
	service.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hi"))
	})
	service.Run()
}
