// Package web is the network web dashboard
package web

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"math/rand"
	"net"
	"net/http"
	"sort"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/config/cmd"
	pb "github.com/micro/go-micro/network/service/proto"
	"github.com/micro/go-micro/web"
)

func color() string {
	return fmt.Sprintf("%d, %d, %d", rand.Intn(255), rand.Intn(255), rand.Intn(255))
}

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

func toGraph(peer *pb.Peer, peers map[string][]string) map[string][]string {
	if peer == nil || peer.Node == nil {
		return peers
	}
	if peers == nil {
		peers = make(map[string][]string)
	}

	// get the current list
	p := peers[peer.Node.Address]

	// first append the peer
	for _, pr := range peer.Peers {
		p = append(p, pr.Node.Address)
		// save the peer list
		peers[peer.Node.Address] = p
		// now walk the peer graph
		peers = toGraph(pr, peers)
	}

	return peers
}

func Run(ctx *cli.Context) {
	c := *cmd.DefaultOptions().Client
	client := pb.NewNetworkService("go.micro.network", c)

	opts := []web.Option{
		web.Name("go.micro.web.network"),
	}

	address := ctx.String("server_address")
	if len(address) > 0 {
		opts = append(opts, web.Address(address))
	}

	// create the web service
	service := web.NewService(opts...)

	// template
	t := template.Must(template.New("layout").Parse(templateFile))
	tg := template.New("layout")
	tg = tg.Funcs(template.FuncMap{"Join": strings.Join})
	tg = tg.Funcs(template.FuncMap{"Color": color})
	tg = template.Must(tg.Parse(graphTemplate))

	// return some data
	service.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// lookup the network
		ips, _ := net.LookupHost("network.micro.mu")
		coreMap := make(map[string]bool)
		for _, ip := range ips {
			coreMap[ip+":8085"] = true
		}

		var graph *pb.Peer
		// get the network graph
		rsp, err := client.Graph(context.Background(), &pb.GraphRequest{})
		if err != nil || rsp.Root == nil {
			return
		}

		// set the root
		graph = rsp.Root

		var coreout []string
		var output []string
		for id, address := range toMap(graph, nil) {
			if _, ok := coreMap[address]; ok {
				coreout = append(coreout, id+"\t"+address)
				continue
			}
			output = append(output, id+"\t"+address)
		}

		// sort output
		sort.Strings(coreout)
		sort.Strings(output)

		// write output
		core := len(coreout)
		dev := len(output)

		b := bytes.NewBuffer(nil)

		heading := fmt.Sprintf("<p>Nodes: %d\tRoot: %s</p>", core+dev, graph.Node.Id)
		b.Write([]byte(heading))
		heading = fmt.Sprintf("<p>Core: %d\tLocale: %s</p>", core, "network.micro.mu")
		b.Write([]byte(heading))
		b.Write([]byte(strings.Join(coreout, "<br>")))
		heading = fmt.Sprintf("<p>Dev: %d\tLocale: %s<p>", dev, "global")
		b.Write([]byte(heading))
		b.Write([]byte(strings.Join(output, "<br>")))

		t.Execute(w, string(b.Bytes()))
	})

	service.HandleFunc("/graph", func(w http.ResponseWriter, r *http.Request) {
		// lookup the network
		ips, _ := net.LookupHost("network.micro.mu")
		coreMap := make(map[string]int)
		coreNodes := []string{}

		// save by index
		for i, ip := range ips {
			coreMap[ip] = i
			coreNodes = append(coreNodes, ip)
		}

		// sort the nodes (in future based on region)
		sort.Strings(coreNodes)

		var graph *pb.Peer
		// get the network graph
		rsp, err := client.Graph(context.Background(), &pb.GraphRequest{})
		if err != nil || rsp.Root == nil {
			return
		}

		// set the root
		graph = rsp.Root

		type graphT struct {
			Nodes []string
			Data  map[string][]string
		}

		graphData := new(graphT)
		// set core nodes aka labels
		graphData.Nodes = coreNodes
		graphData.Data = make(map[string][]string)

		// get the graph per node
		nodeGraph := toGraph(graph, nil)

		// range over the graph and build the data set for each
		for address, nodes := range nodeGraph {
			data := make([]string, len(coreNodes))
			address = strings.TrimSuffix(address, ":8085")

			// set all zeros
			for i := 0; i < len(coreNodes); i++ {
				data[i] = "0"
			}

			// address is one of the core nodes
			if v, ok := coreMap[address]; ok {
				// set self to high value
				data[v] = "100"
			}

			// walk all the
			for _, node := range nodes {
				node = strings.TrimSuffix(node, ":8085")

				// skip self
				if node == address {
					continue
				}
				v, ok := coreMap[node]
				if !ok {
					continue
				}
				if v >= len(data) {
					continue
				}
				// set relationship to 1
				data[v] = "100"
			}

			// save to data
			graphData.Data[address] = data
		}

		tg.Execute(w, graphData)
	})

	service.HandleFunc("/routes", func(w http.ResponseWriter, r *http.Request) {
		// get the network graph
		rsp, err := client.Routes(context.Background(), &pb.RoutesRequest{})
		if err != nil {
			return
		}

		var output []string

		wr := new(tabwriter.Writer)
		wr.Init(w, 0, 8, 2, ' ', 0)

		for _, route := range rsp.Routes {
			metric := fmt.Sprintf("%d", route.Metric)
			if route.Metric >= math.MaxInt64 || route.Metric < 0 {
				metric = "∞"
			}

			// service address gateway router network link
			val := fmt.Sprintf("<td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td>",
				route.Service, route.Address, route.Gateway,
				route.Router, route.Network, route.Link, metric)
			output = append(output, val)
		}

		// sort output
		sort.Strings(output)

		b := bytes.NewBuffer(nil)
		heading := fmt.Sprintf("<p>Routes: %d<p>", len(rsp.Routes))
		b.Write([]byte(heading))
		b.Write([]byte("<table><tr>"))
		b.Write([]byte(strings.Join(output, "</tr><tr>")))
		t.Execute(w, string(b.Bytes())+"</tr></table>")
	})

	service.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		// get the network graph
		rsp, err := client.Services(context.Background(), &pb.ServicesRequest{})
		if err != nil {
			return
		}

		var output []string

		wr := new(tabwriter.Writer)
		wr.Init(w, 0, 8, 2, ' ', 0)

		for _, service := range rsp.Services {
			output = append(output, service)
		}

		// sort output
		sort.Strings(output)

		b := bytes.NewBuffer(nil)
		heading := fmt.Sprintf("<p>Services: %d</p>", len(rsp.Services))
		b.Write([]byte(heading))
		b.Write([]byte(strings.Join(output, "<br>")))
		t.Execute(w, string(b.Bytes()))
	})

	// run the service
	service.Run()
}
