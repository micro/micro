package util

import (
	pb "github.com/micro/micro/v3/proto/network"
	rtrPb "github.com/micro/micro/v3/proto/router"
	"github.com/micro/micro/v3/service/network"
	"github.com/micro/micro/v3/service/router"
)

// PeersToProto returns node peers graph encoded into protobuf
func PeersToProto(node network.Node, depth uint) *pb.Peer {
	// network node aka root node
	pbNode := &pb.Node{
		Id:      node.Id(),
		Address: node.Address(),
		Status: &pb.Status{
			Error: &pb.Error{
				Count: uint32(node.Status().Error().Count()),
				Msg:   node.Status().Error().Msg(),
			},
		},
	}

	// set the network name if network is not nil
	if node.Network() != nil {
		pbNode.Network = node.Network().Name()
	}

	// we will build proto topology into this
	pbPeers := &pb.Peer{
		Node:  pbNode,
		Peers: make([]*pb.Peer, 0),
	}

	for _, peer := range node.Peers() {
		pbPeer := peerProtoTopology(peer, depth)
		pbPeers.Peers = append(pbPeers.Peers, pbPeer)
	}

	return pbPeers
}

func peerProtoTopology(peer network.Node, depth uint) *pb.Peer {
	node := &pb.Node{
		Id:      peer.Id(),
		Address: peer.Address(),
		Status: &pb.Status{
			Error: &pb.Error{
				Count: uint32(peer.Status().Error().Count()),
				Msg:   peer.Status().Error().Msg(),
			},
		},
	}

	// set the network name if network is not nil
	if peer.Network() != nil {
		node.Network = peer.Network().Name()
	}

	pbPeers := &pb.Peer{
		Node:  node,
		Peers: make([]*pb.Peer, 0),
	}

	// return if we reached the end of topology or depth
	if depth == 0 || len(peer.Peers()) == 0 {
		return pbPeers
	}

	// decrement the depth
	depth--

	// iterate through peers of peers aka pops
	for _, pop := range peer.Peers() {
		peer := peerProtoTopology(pop, depth)
		pbPeers.Peers = append(pbPeers.Peers, peer)
	}

	return pbPeers
}

// RouteToProto encodes route into protobuf and returns it
func RouteToProto(route router.Route) *rtrPb.Route {
	return &rtrPb.Route{
		Service: route.Service,
		Address: route.Address,
		Gateway: route.Gateway,
		Network: route.Network,
		Router:  route.Router,
		Link:    route.Link,
		Metric:  int64(route.Metric),
	}
}

// ProtoToRoute decodes protobuf route into router route and returns it
func ProtoToRoute(route *rtrPb.Route) router.Route {
	return router.Route{
		Service: route.Service,
		Address: route.Address,
		Gateway: route.Gateway,
		Network: route.Network,
		Router:  route.Router,
		Link:    route.Link,
		Metric:  route.Metric,
	}
}
