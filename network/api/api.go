// Package api is the network api
package api

import (
	"context"
	"errors"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	goapi "github.com/micro/go-micro/api"
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

type Network struct {
	client pb.NetworkService
	closed chan bool

	mtx   sync.RWMutex
	peers map[string]string
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
	rsp, err := n.client.ListPeers(context.TODO(), &pb.PeerRequest{
		Depth: uint32(1),
	})
	if err != nil {
		log.Debugf("Failed to get peers: %v\n", err)
		return
	}

	n.mtx.Lock()
	defer n.mtx.Unlock()

	setPeers := func(peer *pb.Peer) {
		if peer == nil || peer.Node == nil {
			return
		}
		ip, err := n.getIP(peer.Node.Address)
		if err == nil {
			n.peers[ip] = peer.Node.Id
		} else {
			log.Debugf("Error getting peer IP: %v %+v\n", err, peer.Node)
		}

		for _, p := range peer.Peers {
			ip, err := n.getIP(p.Node.Address)
			if err != nil {
				log.Debugf("Error getting peer IP: %v %+v\n", err, p.Node)
				continue
			}
			n.peers[ip] = p.Node.Id
		}

	}

	// set node 0
	setPeers(rsp.Peers)

	// set node peers depth 1
	for _, peer := range rsp.Peers.Peers {
		setPeers(peer)
	}

	log.Debugf("Set peers: %+v\n", n.peers)
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
func (n *Network) Peers(ctx context.Context, req *map[string]interface{}, rsp *map[string]interface{}) error {
	n.mtx.RLock()
	defer n.mtx.RUnlock()

	var peers []*resolver.Record

	// make copy of peers
	for peer, _ := range n.peers {
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
	netHandler := &Network{
		client: netClient,
		closed: make(chan bool),
		peers:  make(map[string]string),
	}

	// run the cache
	go netHandler.cache()
	defer netHandler.stop()

	// create endpoint
	ep := &goapi.Endpoint{
		Name:    "Network.Peers",
		Path:    []string{"^/network/?$"},
		Method:  []string{"GET"},
		Handler: "rpc",
	}

	// register the handler
	micro.RegisterHandler(api.Server(), netHandler, goapi.WithEndpoint(ep))

	// run the api
	api.Run()
}
