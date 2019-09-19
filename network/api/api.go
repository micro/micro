// Package api is the network api
package api

import (
	"context"
	"sync"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	goapi "github.com/micro/go-micro/api"
	pb "github.com/micro/go-micro/network/proto"
	"github.com/micro/go-micro/network/resolver"
)

type Network struct {
	client pb.NetworkService
	closed chan bool
	mtx    sync.RWMutex
	peers  map[string]string
}

func (n *Network) setCache() {
	rsp, err := n.client.ListPeers(context.TODO(), &pb.PeerRequest{
		Depth: uint32(1),
	})
	if err != nil {
		return
	}

	n.mtx.Lock()
	defer n.mtx.Unlock()

	n.peers[rsp.Peers.Node.Id] = rsp.Peers.Node.Address

	for _, peer := range rsp.Peers.Peers {
		n.peers[peer.Node.Id] = peer.Node.Address
	}
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

func (n *Network) Peers(ctx context.Context, req *map[string]interface{}, rsp *map[string]interface{}) error {
	n.mtx.RLock()
	defer n.mtx.RUnlock()

	var peers []*resolver.Record

	// make copy of peers
	for _, peer := range n.peers {
		peers = append(peers, &resolver.Record{Address: peer})
	}

	// make peer response
	peerRsp := map[string]interface{}{
		"peers": peers,
	}

	// set peer response
	*rsp = peerRsp
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
	netHandler := new(Network)
	// set the net client
	netHandler.client = netClient
	// set the handler cache
	netHandler.closed = make(chan bool)
	netHandler.peers = make(map[string]string)
	// run the cache
	go netHandler.cache()
	defer netHandler.stop()

	// create endpoint
	ep := &goapi.Endpoint{
		Name:    "Network.Peers",
		Path:    []string{"/network"},
		Method:  []string{"GET"},
		Handler: "rpc",
	}

	// register the handler
	micro.RegisterHandler(api.Server(), netHandler, goapi.WithEndpoint(ep))

	// run the api
	api.Run()
}
