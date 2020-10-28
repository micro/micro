---
author: Milos, Jake and Asim
layout:	post
title:	Building a global services network using Go, QUIC and Micro
date:	2019-12-05 09:00:00
---

Over the past 6 months we at [Micro](https://m3o.com/) have been hard at work developing a global service network to build, share and collaborate on microservices.

In this post we're going to share some of the technical details, the design decisions we made, challenges we faced and ultimately how we have succeeded in building the microservices network.

## Motivations

The power of collaborative development has largely been restricted to trusted environments within organisations. When done right, these private in-house platforms unlock incredible productivity and compounding value with every new service added.They provide an always-on runtime and known developer workflow for engineers to collaborate on and deliver new features to their customers. 

Historically, this has been quite difficult to achieve outside of organisations. When developers decide to work on new services they often have to deal with a lot of unnecessary work when it comes to making the services available to others to consume and collaborate on. Public cloud providers are too complex and the elaborate setups when hosting things yourself don’t make things easier either. At [Micro](https://m3o.com/) we felt this pain and decided to do something about it. We built a microservices network!

The micro network looks to solve these problems using a shared global network for micro services. Let’s see how we’ve made this dream a reality!

## Design

The micro network is a globally distributed network based on [go-micro](https://go-micro.dev), a Go microservices framework which enables developers to build services quickly without dealing with the complexity of distributed systems. Go Micro provides strongly opinionated interfaces that are pluggable but also come with sane defaults. This allows Go Micro services to be built once and deployed anywhere, with zero code changes.

The micro network leverages five of the core primitives: registry, transport, broker, client and server. Our default implementations can be found in each package in the [go-micro](https://github.com/micro/go-micro) framework. Community maintained plugins live in the [go-plugins](https://github.com/micro/go-plugins) repo.

The micro "network" is an overloaded term, referring both to the global network over which services discover and communicate with each other and the underpinning system consisting of peer nodes whom connect to each and establish the routes over which services communicate.

The network abstracts away the low level details of distributed system communication at large scale, across any cloud or machine, and allows anyone to build services together without thinking about where they are running. This essentially enables large scale sharing of resources and more importantly microservices.

There are four fundamental concepts that make the micro network possible. These are entirely new and built into [Go Micro](https://go-micro.dev/) as of the last 6 months:

- **Tunnel** - point to point tunnelling
- **Proxy** - transparent rpc proxying
- **Router** - route aggregation and advertising
- **Network** - multi-cloud networking built on the above three

Each of these components is just like any other [Go Micro](https://go-micro.dev/) component - pluggable, with an out of the box default implementation to get started. In our case the micro network it was important that the defaults worked at scale across the world.

Let’s dig into the details.

### Tunnel

From a high level view the micro network is an overlay network that spans the internet. All micro network nodes maintain secure tunnel connections between each other to enable the secure communication between the services running in the network. Go Micro provides a default tunnel implementation using the QUIC protocol along with custom session management.

We chose QUIC because it provides some excellent properties especially when it comes to dealing with high latency networks, an important property when dealing with [running services in large distributed networks](https://eng.uber.com/employing-quic-protocol/). QUIC runs over UDP, but by adding some connection based semantics it supports reliable packet delivery. QUIC also supports multiple streams without [head of line blocking](https://en.wikipedia.org/wiki/Head-of-line_blocking) and it’s designed to work with encryption natively. Finally, QUIC runs in userspace, not in the kernel space on conventional systems, so it can provide both a performance and extra security, too.

Micro tunnel uses [quic-go](https://github.com/lucas-clemente/quic-go) which is the most complete Go implementation of QUIC that we could find at the inception of the micro network. We are aware quic-go is a work in progress and that it can occasionally break, but we are happy to pay the early adopter cost as we believe QUIC will become the defacto standard internet communication protocol in the future, enabling large scale networks such as the micro network.

Let’s look at the Go Micro tunnel interface:

```go
// Tunnel creates a gre tunnel on top of the go-micro/transport.
// It establishes multiple streams using the Micro-Tunnel-Channel header
// and Micro-Tunnel-Session header. The tunnel id is a hash of
// the address being requested.
type Tunnel interface {
	// Address the tunnel is listening on
	Address() string
	// Connect connects the tunnel
	Connect() error
	// Close closes the tunnel
	Close() error
	// Links returns all the links the tunnel is connected to
	Links() []Link
	// Dial to a tunnel channel
	Dial(channel string, opts ...DialOption) (Session, error)
	// Accept connections on a channel
	Listen(channel string, opts ...ListenOption) (Listener, error)
}
```

It may look fairly familiar to Go developers. With Go Micro we’ve tried to maintain common interfaces in line with distributed systems development while stepping in at a lower layer to solve some of the nitty gritty details.

Most of the interface methods should hopefully be self-explanatory, but you might be wondering about channels and sessions. Channels are much like addresses, providing a way to segment different message streams over the tunnel. Listeners listen on a given channel and return a unique session when a client dials into the channel. The session is used to communicate between peers on the same tunnel channel. The Go Micro tunnel provides different communication semantics too. You can choose to use either unicast or multicast. 

<img src="https://m3o.com/docs/images/session.svg"  alt="" />

In addition tunnels enable bidirectional connections; sessions can be dialled or listened from either side. This enables the reversal of connections so anything behind a [NAT](https://en.wikipedia.org/wiki/Network_address_translation) or without a public IP can become a server.

### Router

Micro router is a critical component of the micro network. It provides the network’s routing plane. Without the router, we wouldn’t know where to send messages. It constructs a routing table based on the local service registry (a component of Go Micro). The routing table maintains the routes to the services available on the local network. With the tunnel its then also able to process messages from any other datacenter or network enabling global routing by default.

Our default routing table implementation uses a simple Go in memory map, but as with all things in Go Micro, the router and routing table are both pluggable. As we scale we’re thinking about alternative implementations and even the possibility of switching dynamically based on the size of networks.

The Go Micro router interface is as follows:

```go
// Router is an interface for a routing control plane
type Router interface {
	// The routing table
	Table() Table
	// Advertise advertises routes to the network
	Advertise() (<-chan *Advert, error)
	// Process processes incoming adverts
	Process(*Advert) error
	// Solicit advertises the whole routing table to the network
	Solicit() error
	// Lookup queries routes in the routing table
	Lookup(...QueryOption) ([]Route, error)
	// Watch returns a watcher which tracks updates to the routing table
	Watch(opts ...WatchOption) (Watcher, error)
}
```

When the router starts it automatically creates a watcher for its local registry. The micro registry emits events any time services are created, updated or deleted. The router processes these events and then applies actions to its routing table accordingly. The router itself advertises the routing table events which you can think of as a cut down version of the registry solely concerned with routing of requests where as the registry provides more feature rich information like api endpoints.

These routes are propagated as events to other routers on both the local and global network and applied by every router to their own routing table. Thus maintaining the global network routing plane.

Here’s a look at a typical route:

```go
// Route is a network route
type Route struct {
	// Service is destination service name
	Service string
	// Address is service node address
	Address string
	// Gateway is route gateway
	Gateway string
	// Network is the network name
	Network string
	// Router is router id
	Router string
	// Link is networks link
	Link string
	// Metric is the route cost
	Metric int64
}
```

What we’re primarily concerned with here is routing by service name first, finding its address if its local or a gateway if we have to go through some remote endpoint or different network. We also want to know what type of Link to use e.g whether routing through our tunnel, Cloudflare Argo tunnel or some other network implementation. And then most importantly the metric a.k.a. the cost of routing to that node. We may have many routes and we want to take routes with optimal cost to ensure lowest latency. This doesn’t always mean your request is sent to the local network though! Imagine a situation when the service running on your local network is overloaded. We will always pick the route with the lowest cost no matter where the service is running.

### Proxy

We’ve already discussed the tunnel - how messages get from point to point, and routing - detailing how to find where the services are, but then the question really is how do services actually make use of this? For this we really need a proxy.

It was important to us when building the micro network that we build something that was native to micro and capable of understanding our routing protocol. Building another VPN or IP based networking solution was not our goal. Instead we wanted to facilitate communication between services.

When a service needs to communicate with other services in the network it uses micro proxy.
 
The proxy is a native RPC proxy implementation built on the Go Micro `Client` and `Server` interfaces. It encapsulates the core means of communication for our services and provides a forwarding mechanism for requests based on service name and endpoints. Additionally it has the ability to also act as a messaging exchange for asynchronous communication since Go Micro supports both request/response and pub/sub communication. This is native to Go Micro and a powerful building block for request routing.

The interface itself is straightforward and encapsulates the complexity of proxying.

```go
// Proxy can be used as a proxy server for go-micro services
type Proxy interface {
	// ProcessMessage handles inbound messages
	ProcessMessage(context.Context, server.Message) error
	// ServeRequest handles inbound requests
	ServeRequest(context.Context, server.Request, server.Response) error
}
```

The proxy receives RPC requests and routes them to an endpoint. It asks the router for the location of the service (caching as needed) and decides based on the `Link` field in the routing table whether to send the request locally or over the tunnel across the global network. The value of the `Link` field is either `“local"` (for local services) or `“network"` if the service is accessible only via the network.

Like everything else, the proxy is something we built standalone that can work between services in one datacenter but also across many when used in conjunction with the tunnel and router.

And finally arriving at the pièce de résistance. The network interface.

### Network

Network nodes are the magic that ties all the core components together. Enabling the ability to build a truly global service network. It was really important when creating the network interface that it fit inline with our existing assumptions and understanding about Go Micro and distributed systems development. We really wanted to embrace the existing interfaces of the framework and design something with symmetry in regards to a Service. 

What we arrived at was something very similar to the [micro.Service](https://github.com/micro/go-micro/blob/master/micro.go#L16) interface itself


```go
// Network is a micro network
type Network interface {
	// Node is network node
	Node
	// Name of the network
	Name() string
	// Connect starts the resolver and tunnel server
	Connect() error
	// Close stops the tunnel and resolving
	Close() error
	// Client is micro client
	Client() client.Client
	// Server is micro server
	Server() server.Server
}

// Node is a network node
type Node interface {
	// Id is node id
	Id() string
	// Address is node bind address
	Address() string
	// Peers returns node peers
	Peers() []Node
	// Network is the network node is in
	Network() Network
}
```

As you can see, a `Network` has a Name, Client and Server, much like a `Service`, so it provides a similar method of communication. This means we can reuse a lot of the existing code base, but it also goes much further. A `Network` includes the concept of a `Node` directly in the interface, one which has peers and whom may belong to the same network or others. This means is networks are peer-to-peer while Services are largely focused on Client/Server. On a day to day basis developers stay focused on building services but these when built to communicate globally need to operate across networks made up of identical peers.

Our networks have the ability to behave as peers which route for others but also may provide some sort of service themselves. In this case it's mostly routing related information.

So how does it all work together?

Networks have a list of peer nodes to talk to. In the case of the default implementation the peer list comes from the registry with other network nodes with the same name (the name of the network itself). When a node starts it “connects" to the network by establishing its tunnel, resolving the nodes and then connecting to them. Once they’ve connected the nodes peer over two multicast sessions, one for peer announcements and the other for route advertisements. As these propagate the network begins to converge on identical routing information building a full mesh that allows for routing of services from any node to the other.

The nodes maintain keepalives, periodically advertise the full routing table and flush any events as they occur. Our core network nodes make use of multiple resolvers to find each other, including DNS and the local registry. In the case of peers that join our network, we’ve configured them to use a http resolver which gets geo-steered via Cloudflare anycast DNS and global load balanced to the most local region. From there they pull a list of nodes and connect to the ones with the lowest metric. They then repeat the same song and dance as above to continue the growth of the network and participate in service routing.

Each node maintains its own network graph based on the peer messages it receives. Peer messages contain the graph of each peer up to 3 hops which enables the ability for every node to build a local view of the network. Peers ignore anything with more than a 3 hop radius. This is to avoid potential performance problems.

We mentioned a little something about peer and route advertisements. So what message do the network nodes actually exchange? First, the network embeds the router interface through which it advertises its local routes to other network nodes. These routes are then propagated across the whole network, much like the internet. The node itself receives route advertisements from its peers and applies the advertised changes to its routing own routing table. The message types are “solicit" to ask for routes and “advert" for updates broadcast.

Network nodes send “connect" messages on start and “close" on exit. For their lifetime they are periodically broadcasting “peer" messages so that others can discover them and they all can build the network topology.

When the network is created and converges, services are then capable of sending messages across it. When a service on the network needs to communicate with some other service on the network it sends a request to the network node. The micro network node embeds micro proxy and thus has the ability to forward the request through network or locally if it deems so more fit based on the metrics it retrieves after looking up the routes in the routing table.

This as a whole forms our micro services network.

## Challenges

Building a global services is not without its challenges. We encountered many from the initial design phase right through to the present day of dealing with broken nodes, bad actors, event storms and more.

### Initial Implementation

The actual task we’d set out to accomplish was pretty monumental and we’d underestimated how much effort it would take even in an MVP phase of the first implementation.

Every time we attempted to go from design diagram to implementing code we found ourselves stuck. In theory everything made sense but no matter how many times we attempted to write code things just didn’t click.

We wrote 3-4 implementations that were essentially thrown away before figuring out the best approach was to make local networking work first and then slowly carve out specific problems to solve. So proxying, following by routing and then a network interface. Eventually when these pieces were in place we could get to multi-cloud or global networking by implementing a tunnel to handle the heavy lifting.

<center>
<img src="https://m3o.com/images/it-works.jpg" style="width: 80%; height: auto;" />
</center>

Once again, the lesson is to keep it simple, but where the task itself is complex, break it down into steps you can actually keep simple independently and then piece back together in practice.

### Multipoint Tunneling

One of the most complex pieces of code we had to write was the tunnel. It's still not where we’d like it to be but it was pretty important to us to write this from the ground up so we’d have a good understanding of how we wanted to operate globally but also have full control over the foundations of the system.

The complexity in writing network protocols really came to light in this effort, from trying to NOT reimplement tcp, crypto or anything else but also find a path to a working solution. In the end we were able to create a subset of commands which formed enough of a bidirectional and multi-session based tunnel over QUIC. We left most of the heavy lifting to QUIC but we also needed the ability to do multicast.

For us it didn’t make sense to just rely on unicast, considering the async and pubsub based semantics built into Go Micro we felt pretty adamant it needed to be part of core network routing. So with that sessions needed to be reimplemented on top of QUIC.

We’ll spare you the gory details but what’s really clear to us is that writing state machines and reliable connection protocol code is not where we want to spend the majority of our time. We have a fairly resilient working solution but our hope is to replace this with something far better in the future.

### Event Storms

When things work they work and when they break they break badly. For us everything came crashing down when we started to encounter broadcast storms caused by services being recycled. When a service is created the service registry fires a create event and when it’s shutting down it automatically deregisters from service registry which fires a delete event. Maybe you can see where this is going. As services cycled in our network they’d generate these events which leads to the routers generating new route events which are then propagated every 5 seconds to every other node on the network. 

This sounds ok if the network converges and they stop propagating events but in our case the sequence of events are observed and applied at random time intervals on every node. This in essence can lead to a broadcast storm which never stops. Understanding and resolving this is an incredibly difficult task.

For us this really led to research in BGP internet routing in which they’ve defined flap detection algorithms to handle this. We’ve read a few whitepapers to get familiar with the concepts and hacked up a simple flap detection algorithm in the router. 

At its core, the flap detection assigns a numerical cost to every route event. When a route event occurs it’s cost gets incremented. If the same event happens multiple times within a certain period of time and the accumulated cost reaches a predefined threshold the event is immediately suppressed. Suppressed events are not processed by router, but are kept in a memory cache for a redefined period of time. Meanwhile the cost of the event decays with time whilst at the same time it can keep on growing if the event keeps on flapping. If the event drops below another threshold the event is unsuppressed and can be processed by the routers. If the event remains suppressed for longer than a predefined time period it’s discarded. 

The picture below depicts how the decaying actually works.

<center>
<img src="https://m3o.com/assets/images/flap-detection.png" style="width: 80%; height: auto;" />
</center>

<small>source: http://linuxczar.net/blog/2016/01/31/flap-detection/</small>

This had a huge effect on the issues we had been experiencing in the network. The nodes were no longer hammered with crazy event storms and the network stabilised and continued to work without any interruptions. Happy days!

## Architecture

Our overall goal is to build a micro services network that manages not only communication but all aspects of running services, governance, and more. To accomplish this we started by addressing networking from the ground up for Go Micro services. Not just to communicate locally within one private network but to have the ability to do so across many networks. 

For this purpose we’ve created a global multi-cloud network that enables communication from anywhere, with anyone. This is fundamental to the architecture of the micro services network. 

Our next goal will be to tackle the runtime aspects so that we offer the ability to host services without the need to manage them. This could be imagined as the basis of a serverless microservices platform which we’re looking to launch soon.

The platform is designed to be open. Anyone should be able to run services on the platform or join the global network using their own node. What’s more, you can even pick up the open source code and build their own private networks or join theirs to our public one.

<center>
  <img src="https://github.com/micro/development/raw/f4c77580acac228c522623c217575fb266d2d4ab/images/arch.jpg" style="width: 80%; height: auto;" />
</center>
<br>

What we think is pretty cool and rather unique about the micro network is the network nodes themselves are just regular micro services like any other. Because we built everything using Go Micro they behave just like any other service. In fact what’s even more exciting is that literally *everything is a service* in the micro network. 

This holds true for all the individual components that make up the network. If you don’t want to run full network nodes, you can also run individual components of the network separately as standalone micro services such as the tunnel, router and proxy. All the components register themselves with local registry via which they can be discovered.

## Eventual success

On 29th August 2019 around 4PM we sent the first successful RPC request between our laptops across the internet using the micro network.

<center>
  <img src="https://m3o.com/assets/images/success.jpg" style="width: 80%; height: auto;" />
</center>
<br>

Since then we have squashed a lot of bugs and deployed the network nodes across the globe.
At the moment we are running the micro network in 4 cloud providers across 4 geographical regions with 3 nodes in each region.

<center>
  <img src="https://m3o.com/assets/images/radar.png" style="width: 80%; height: auto;" />
</center>

## Usage

If you're interested in testing out micro and the network just do the following.

```go
# enable go modules
export GO111MODULE=on

# download micro
go get github.com/micro/micro@master

# connect to the network
micro --peer
```

Now you're connected to the network. Start to explore what's there.

```go
# List the services in the network
micro network services

# See which nodes you're connected to
micro network connections

# List all the nodes in your network graph
micro network nodes

# See what the metrics look like to different service routes
micro network routes
```

So what does a micro network developer workflow look like? Developers write their Go code using the [Go Micro](https://github.com/micro/go-micro) framework and once they’re ready they can make their services available on the network either directly from their laptop or from anywhere the micro network node runs (more on what micro network node is later).

Here is an example of a simple service written using `go-micro`:

```go
package main

import (
	"context"
	"log"
	"time"

	hello "github.com/micro/examples/greeter/srv/proto/hello"
	"github.com/micro/go-micro"
)

type Say struct{}

func (s *Say) Hello(ctx context.Context, req *hello.Request, rsp *hello.Response) error {
	log.Print("Received Say.Hello request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

func main() {
	service := micro.NewService(
		micro.Name("helloworld"),
	)

	// optionally setup command line usage
	service.Init()

	// Register Handlers
	hello.RegisterSayHandler(service.Server(), new(Say))

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
```

Once you launch the service it automatically registers with service registry and becomes instantly accessible to everyone on the network to consume and collaborate on. All of this is completely transparent to developers. No need to deal with low level distributed systems cruft!

We’re already running a greeter service in the network so why not try giving it a call.

```
# enable proxying through the network
export MICRO_PROXY=go.micro.network

# call a service
micro call go.micro.srv.greeter Say.Hello '{"name": "John"}'
```

It works!

## Conclusion

Building distributed systems is difficult, but it turns out building the networks they communicate over is an equally, if not more difficult, problem. The classic fallacy, [the network is reliable](https://queue.acm.org/detail.cfm?id=2655736), continues to hold, as we found while building the micro network. However what’s also clear is that our world and most technology thrives through the use of networks. They underpin the very fabric of all that we’ve come to know. Our goal with the micro network is to create a new type of foundation for the open services of the future. Hopefully this post shed some light on the technical accomplishments and challenges of building such a thing.

<br />
To learn more check out the [website](https://m3o.com), follow us on [twitter](https://twitter.com/m3ocloud) or 
join the [slack](https://m3o.com/slack) community. We are hiring!

