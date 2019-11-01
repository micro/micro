// Package web is the network web dashboard
package web

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/micro/cli"
	"github.com/micro/go-micro/config/cmd"
	pb "github.com/micro/go-micro/network/proto"
	"github.com/micro/go-micro/web"
)

func toMap(peer *pb.Peer, peers map[string]string) map[string]string {
	if peer == nil || peer.Node == nil {
		return peers
	}
	if peers == nil {
		peers = make(map[string]string)
	}
	peers[peer.Node.Id] = peer.Node.Address
	for _, p := range peer.Peers {
		toMap(p, peers)
	}
	return peers
}

func Run(ctx *cli.Context) {
	c := *cmd.DefaultOptions().Client
	client := pb.NewNetworkService("go.micro.network", c)

	opts := []web.Option{
		web.Name("go.micro.web.network"),
	}

	address := ctx.GlobalString("server_address")
	if len(address) > 0 {
		opts = append(opts, web.Address(address))
	}

	// create the web service
	service := web.NewService(opts...)

	// return some data
	service.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var graph *pb.Peer
		// get the network graph
		rsp, err := client.Graph(context.Background(), &pb.GraphRequest{})
		if err != nil || rsp.Root == nil {
			return
		}

		// set the root
		graph = rsp.Root

		var output []string
		for id, address := range toMap(graph, nil) {
			output = append(output, id+"\t"+address)
		}

		// sort output
		sort.Strings(output)

		// write output
		heading := fmt.Sprintf("Nodes: %d\tRoot: %s\n\n", len(output), graph.Node.Id)
		w.Write([]byte(heading))
		w.Write([]byte(strings.Join(output, "\n")))
	})

	service.HandleFunc("/routes", func(w http.ResponseWriter, r *http.Request) {
		// get the network graph
		rsp, err := client.Routes(context.Background(), &pb.RoutesRequest{})
		if err != nil {
			return
		}

		heading := fmt.Sprintf("Routes: %d\n\n", len(rsp.Routes))
		w.Write([]byte(heading))

		var output []string

		wr := new(tabwriter.Writer)
		wr.Init(w, 0, 8, 2, ' ', 0)

		for _, route := range rsp.Routes {
			// service address gateway router network link
			val := fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s",
				route.Service, route.Address, route.Gateway,
				route.Router, route.Network, route.Link)
			output = append(output, val)
		}

		// sort output
		sort.Strings(output)

		wr.Write([]byte(strings.Join(output, "\n")))
		wr.Flush()
	})

	// run the service
	service.Run()
}
