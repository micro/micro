// Package api is the network api
package api

import (
	"context"
	"errors"
	"net"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	goapi "github.com/micro/go-micro/api"
	"github.com/micro/go-micro/network"
	pb "github.com/micro/go-micro/network/proto"
	"github.com/micro/go-micro/network/resolver"
	"github.com/micro/go-micro/util/log"
)

var (
	privateBlocks []*net.IPNet
)

func init() {
	for _, b := range []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "100.64.0.0/10", "fd00::/8"} {
		if _, block, err := net.ParseCIDR(b); err == nil {
			privateBlocks = append(privateBlocks, block)
		}
	}
}

func isPrivateIP(ip net.IP) bool {
	for _, priv := range privateBlocks {
		if priv.Contains(ip) {
			return true
		}
	}
	return false
}

type node struct {
	id      string
	address string
	network string
	metric  int64
}

type Network struct {
	client pb.NetworkService
	closed chan bool

	mtx sync.RWMutex
	// indexed by address
	nodes map[string]*node
}

func (n *Network) getIP(addr string) (string, error) {
	if strings.HasPrefix(addr, "[::]") {
		return "", errors.New("ip is loopback")
	}

	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return "", err
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return "", errors.New("ip is blank")
	}

	if isPrivateIP(ip) {
		return "", errors.New("private ip")
	}

	return addr, nil
}

func (n *Network) setCache() {
	rsp, err := n.client.Graph(context.TODO(), &pb.GraphRequest{
		Depth: uint32(1),
	})
	if err != nil {
		log.Debugf("Failed to get nodes: %v\n", err)
		return
	}

	routesRsp, err := n.client.Routes(context.TODO(), &pb.RoutesRequest{})
	if err != nil {
		log.Debugf("Failed to get routes: %v\n", err)
		return
	}

	// create a map of the peers
	peers := make(map[string]string)

	setPeers := func(peer *pb.Peer) {
		if peer == nil || peer.Node == nil {
			return
		}
		ip, err := n.getIP(peer.Node.Address)
		if err == nil {
			peers[ip] = peer.Node.Id
		} else {
			log.Debugf("Error getting peer IP: %v %+v\n", err, peer.Node)
		}

		for _, p := range peer.Peers {
			ip, err := n.getIP(p.Node.Address)
			if err != nil {
				log.Debugf("Error getting peer IP: %v %+v\n", err, p.Node)
				continue
			}
			peers[ip] = p.Node.Id
		}

	}

	// set node 0
	setPeers(rsp.Root)

	// set node nodes depth 1
	for _, peer := range rsp.Root.Peers {
		setPeers(peer)
	}

	// don't proceed without peers
	if len(peers) == 0 {
		log.Debugf("Peer node list is 0... not saving")
		return
	}

	// build a route graph
	nodes := make(map[string]*node)

	for _, route := range routesRsp.Routes {
		// skip routes without a gateway
		if len(route.Gateway) == 0 {
			continue
		}

		// skip routes without a public ip
		ip, err := n.getIP(route.Gateway)
		if err != nil {
			continue
		}

		// find in the peer graph
		id, ok := peers[ip]
		if !ok {
			continue
		}

		// only get peers where the router id matches
		if route.Router != id {
			continue
		}

		// skip already saved routes
		if _, ok := nodes[ip+route.Network]; ok {
			continue
		}

		// add as gateway + network
		nodes[ip+route.Network] = &node{
			id:      id,
			address: ip,
			network: route.Network,
			metric:  route.Metric,
		}
	}

	// set the nodes
	n.mtx.Lock()
	n.nodes = nodes
	n.mtx.Unlock()

	log.Debugf("Set nodes: %+v\n", nodes)
}

func (n *Network) cache() {
	t := time.NewTicker(time.Minute)
	defer t.Stop()

	// set the cache
	n.setCache()

	for {
		select {
		case <-t.C:
			n.setCache()
		case <-n.closed:
			return
		}
	}
}

func (n *Network) stop() {
	select {
	case <-n.closed:
		return
	default:
		close(n.closed)
	}
}

// TODO: get remote IP and compare to peer list to order by nearest nodes
func (n *Network) Nodes(ctx context.Context, req *map[string]interface{}, rsp *map[string]interface{}) error {
	n.mtx.RLock()
	defer n.mtx.RUnlock()

	network := network.DefaultName

	// check if the option is passed in
	r := *req
	if name, ok := r["name"].(string); ok && len(name) > 0 {
		network = name
	}

	var nodes []*resolver.Record

	// make copy of nodes
	for _, node := range n.nodes {
		// only accept of this network name if exists
		if node.network != network {
			continue
		}
		// append to nodes
		nodes = append(nodes, &resolver.Record{
			Address:  node.address,
			Priority: node.metric,
		})
	}

	// sort by lowest priority
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].Priority < nodes[j].Priority })

	// make peer response
	nodeRsp := map[string]interface{}{
		"nodes": nodes,
	}

	// set peer response
	*rsp = nodeRsp
	return nil
}

func Run(ctx *cli.Context) {
	// create the api service
	api := micro.NewService(
		micro.Name("go.micro.api.network"),
	)

	// create the network client
	netClient := pb.NewNetworkService("go.micro.network", api.Client())

	// create new api network handler
	netHandler := &Network{
		client: netClient,
		closed: make(chan bool),
		nodes:  make(map[string]*node),
	}

	// run the cache
	go netHandler.cache()
	defer netHandler.stop()

	// create endpoint
	ep := &goapi.Endpoint{
		Name:    "Network.Nodes",
		Path:    []string{"^/network/?$"},
		Method:  []string{"GET"},
		Handler: "rpc",
	}

	// register the handler
	micro.RegisterHandler(api.Server(), netHandler, goapi.WithEndpoint(ep))

	// run the api
	api.Run()
}
