# Network

The micro network service is a multi-cloud service networking solution which works across public and private environments.

## Overview

The micro network provides multi-cloud capability and builds a large scale flat network over which all services can communicate with each other. 
It makes use of our proxy, router, tunnel and network packages in go-micro to produce global routing across any environment. 

The network generates a routing table based on the local service registry and shares this amongst nodes. It builds in a router and proxy so 
any request made to any network node can be routed across the global network. It prioritises local routing first and can range up to 3 
hops in a chain if needed.

The most technical explanation currently lives at [https://micro.mu/blog/2019/12/05/building-a-microservices-network.html](https://micro.mu/blog/2019/12/05/building-a-microservices-network.html) and the wider product related focus [roadmap/network](https://github.com/m3o/development/blob/master/roadmap/network.md).

## Run Network

Start the network seed node (Runs on port :8085)

```shell
micro network
```

Start the next nodes in a different environment connecting to the first (assuming its running at 10.0.0.1:8085)

```shell
micro network --nodes=10.0.0.1:8085
```

## Network Services

You may now list the nodes, routes, services and graph

```
# list the nodes
micro network nodes

# list the routes
micro network routes

# list the services
micro network services

# print the graph
micro network graph
```

Any request now made through the network will be proxied to a service on the other side.

Set your proxy to use the network
```
MICRO_PROXY=go.micro.network go run main.go
```

Your service will direct all traffic through the network. 


## Authentication

Specify a network token to limit access to the network.

```
MICRO_NETWORK_TOKEN=foobar micro network
```

Nodes must provide a valid and matching token to join the network. The default token is "go.micro.tunnel" which allows 
any node to join and communicate between them.

