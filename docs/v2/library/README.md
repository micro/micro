---
title: Framework
keywords: go-micro, framework
tags: [go-micro, framework]
sidebar: home_sidebar
permalink: /framework
summary: Go Micro is a framework for microservices development
---

# Overview

Go Micro provides the core requirements for distributed systems development including RPC and Event driven communication. 
The micro philosophy is sane defaults with a pluggable architecture. We provide defaults to get you started quickly but everything can be easily swapped out.

## Features

Go Micro abstracts away the details of distributed systems. Here are the main features.

- **Service Discovery** - Automatic service registration and name resolution. Service discovery is at the core of micro service
development. When service A needs to speak to service B it needs the location of that service. The default discovery mechanism is
multicast DNS (mdns), a zeroconf system.

- **Load Balancing** - Client side load balancing built on service discovery. Once we have the addresses of any number of instances
of a service we now need a way to decide which node to route to. We use random hashed load balancing to provide even distribution
across the services and retry a different node if there's a problem.

- **Message Encoding** - Dynamic message encoding based on content-type. The client and server will use codecs along with content-type
to seamlessly encode and decode Go types for you. Any variety of messages could be encoded and sent from different clients. The client
and server handle this by default. This includes protobuf and json by default.

- **Request/Response** - RPC based request/response with support for bidirectional streaming. We provide an abstraction for synchronous
communication. A request made to a service will be automatically resolved, load balanced, dialled and streamed. The default
transport is [gRPC](https://grpc.io/).

- **Async Messaging** - PubSub is built in as a first class citizen for asynchronous communication and event driven architectures.
Event notifications are a core pattern in micro service development. The default messaging system is an embedded [NATS](https://nats.io/)
server.

- **Pluggable Interfaces** - Go Micro makes use of Go interfaces for each distributed system abstraction. Because of this these interfaces
are pluggable and allows Go Micro to be runtime agnostic. You can plugin any underlying technology. Find plugins in
[github.com/micro/go-plugins](https://github.com/micro/go-plugins).

## Getting started

- [Dependencies](#dependencies)
- [Installation](#installation)
- [Writing a Service](#writing-a-service)

## Dependencies

Go Micro makes use of protobuf by default. This is so we can code generate boilerplate code and also provide 
an efficient wire format for transferring data back and forth between services.

We also require some form of service discovery to resolve service names to their addresses along with 
metadata and endpoint information. See below for more info.

### Protobuf

You'll need to install protobuf to code generate API interfaces:

- [protoc-gen-micro](https://github.com/micro/micro/tree/master/cmd/protoc-gen-micro)

### Discovery

Service discovery is used to resolve service names to addresses. By default we provide a zeroconf discovery system 
using multicast DNS. This comes built-in on most operating systems. If you need something more resilient and multi-host then use etcd.

### Etcd

[Etcd](https://github.com/etcd-io/etcd) can be used as an alternative service discovery system.

- Download and run [etcd](https://github.com/etcd-io/etcd)
- Pass `--registry=etcd` to any command or the enviroment variable `MICRO_REGISTRY=etcd`

```
MICRO_REGISTRY=etcd go run main.go
```

Discovery is pluggable. Find plugins for consul, kubernetes, zookeeper and more in the [micro/go-plugins](https://github.com/micro/go-plugins) repo.

## Installation

Go Micro is a framework for Go based development. You can easily get this with the go toolchain.


Import go-micro in your service

```
import "github.com/micro/go-micro/v2"
```

We provide release tags and would recommend to stick with the latest stable releases. Making use of go modules will enable this.

```
# enable go modules
export GO111MODULE=on
# initialise go modules in your app
go mod init
# now go get
go get ./...
```

## Writing a service

Go checkout the [hello world](helloworld.html) example to get started

