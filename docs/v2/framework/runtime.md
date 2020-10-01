---
title: Runtime
keywords: runtime
tags: [micro, runtime]
sidebar: home_sidebar
permalink: /runtime
summary: Micro is a runtime for microservices development
---

# Overview

Micro addresses the key requirements for building scalable systems. It takes the microservice architecture pattern and transforms it into 
a set of tools which act as the building blocks of a platform. Micro deals with the complexity of distributed systems and provides 
simple abstractions already understood by developers.

Technology is constantly evolving. The infrastructure stack is always changing. Micro is a pluggable toolkit which addresses these issues. 
Plug in any stack or underlying technology. Build future-proof systems using micro.

## Features

The runtime is composed of the following features:

- **API Gateway:** A single http entry point with dynamic request routing using service discovery. The API gateway allows you to build a scalable 
microservice architecture on the backend and consolidate serving a public api on the frontend. The micro api provides powerful routing 
via discovery and pluggable handlers to serve http, grpc, websockets, publish events and more.

- **Interactive CLI:** A CLI to describe, query and interact directly with your platform and services from the terminal. The CLI 
gives you all the commands you expect to understand what's happening with your micro services. It also includes an interactive mode.

- **Service Proxy:** A transparent proxy built on [Go Micro](https://github.com/micro/go-micro) and the [MUCP](https://github.com/micro/protocol) 
protocol. Offload service discovery, load balancing, message encoding, middleware, transport and broker plugins to a single a location. 
Run it standalone or alongside your service.

- **Service Templates:** Generate new service templates to get started quickly. Micro provides predefined templates for writing micro services. 
Always start in the same way, build identical services to be more productive.

- **SlackOps Bot:** A bot which runs on your platform and lets you manage your applications from Slack itself. The micro bot enables ChatOps 
and gives you the ability to do everything with your team via messaging. It also includes ability to create slack commmands as services which 
are discovered dynamically.

- **Web Dashboard:** The web dashboard allows you to explore your services, describe their endpoints, the request and response formats and even 
query them directly. The dashboard is also includes a built in CLI like experience for developers who want to drop into the terminal on the fly.

## Install Micro

```shell
go get github.com/micro/micro/v2
```

Or via Docker

```shell
docker pull micro/micro
```

Latest release binaries

```
# MacOS
curl -fsSL https://raw.githubusercontent.com/micro/micro/master/scripts/install.sh | /bin/bash

# Linux
wget -q  https://raw.githubusercontent.com/micro/micro/master/scripts/install.sh -O - | /bin/bash

# Windows
powershell -Command "iwr -useb https://raw.githubusercontent.com/micro/micro/master/scripts/install.ps1 | iex"
```

## Dependencies

The micro runtime has two dependencies: 

- [Service Discovery](#service-discovery) - used for name resolution
- [Protobuf](#protobuf) - used for code generation

## Service Discovery

Service discovery is used for name resolution, routing and centralising metadata. All services are register themselves in 
discovery so they can be found by other services.

Micro uses the [go-micro](https://github.com/micro/go-micro) registry interface for service discovery. Multicast DNS is the default 
implementation. This enables a zeroconf setup so you don't have to change anything for local development. In case you're using docker, 
kubernetes or need something more resilient we recommend etcd.

### Etcd

Etcd is an alternative highly available service discovery mechanism

```shell
# install
brew install etcd

# run
etcd
```

Pass `--registry=etcd` or set the env var `MICRO_REGISTRY=consul` for any command

```shell
# Use flag
micro --registry=etcd list services

# Use env var
MICRO_REGISTRY=etcd micro list services`
```

Specify the address of etcd where it's not on the same host.

```
MICRO_REGISTRY_ADDRESS=10.0.0.1:2379
```

See [go-plugins](https://github.com/micro/go-plugins) for more service discovery plugins.

## Protobuf

Protobuf is used for code generation. It reduces the amount of boilerplate code needed to be written.

```
# install protobuf
brew install protobuf

# install protoc-gen-go
go get github.com/golang/protobuf/{proto,protoc-gen-go}

# install protoc-gen-micro
go get github.com/micro/micro/v2/cmd/protoc-gen-micro@master
```

See [protoc-gen-micro](https://github.com/micro/micro/tree/master/cmd/protoc-gen-micro) for more details.

## Writing a service

Micro includes new template generation to speed up writing applications

For full details on writing services see [**go-micro**](https://github.com/micro/go-micro).

### Generate template

Here we'll quickly generate an example template using `micro new`

Specify a path relative to $GOPATH

``` 
micro new github.com/micro/example
```

The command will output

```
example/
	Dockerfile	# A template docker file
	README.md	# A readme with command used
	handler/	# Example rpc handler
	main.go		# The main Go program
	proto/		# Protobuf directory
	subscriber/	# Example pubsub Subscriber
```

Compile the protobuf code using `protoc`

```
protoc --proto_path=. --micro_out=. --go_out=. proto/example/example.proto
```

Now run it like any other go application

```
# Run the micro service
micro server

# cd into your service directory
cd github.com/micro/example

# run the server
micro run --server .
```

## Example

Now we have a running application using `micro new` template generation, let's test it out.

- [List services](#list-services)
- [Get service](#get-service)
- [Call service](#call-service)
- [Run API](#run-api)
- [Call API](#call-api)

### List services

Each service registers with discovery so we should be able to find it.

```shell
micro list services
```

Output
```
etcd
go.micro.srv.example
topic:topic.go.micro.srv.example
```

The example app has registered with the fully qualified domain name `go.micro.srv.example`

### Get Service

Each service registers with a unique id, address and metadata.

```shell
micro get service go.micro.srv.example
```

Output
```
service  go.micro.srv.example

version latest

ID	Address	Port	Metadata
go.micro.srv.example-437d1277-303b-11e8-9be9-f40f242f6897	192.168.1.65	53545	transport=http,broker=http,server=rpc,registry=etcd

Endpoint: Example.Call
Metadata: stream=false

Request: {
	name string
}

Response: {
	msg string
}


Endpoint: Example.PingPong
Metadata: stream=true

Request: {}

Response: {}


Endpoint: Example.Stream
Metadata: stream=true

Request: {}

Response: {}


Endpoint: Func
Metadata: subscriber=true,topic=topic.go.micro.srv.example

Request: {
	say string
}

Response: {}


Endpoint: Example.Handle
Metadata: subscriber=true,topic=topic.go.micro.srv.example

Request: {
	say string
}

Response: {}
```

### Call service

Make an RPC call via the CLI. The query is sent as json.

```shell
micro call go.micro.srv.example Example.Call '{"name": "John"}'
```

Output
```
{
	"msg": "Hello John"
}
```

Look at the [cli doc](https://micro.mu/docs/cli.html) for more info.

Now let's test call the service via HTTP.

### Run API

The micro api is a http gateway which dynamically routes to backend services

Let's run it so we can query the example service.

```
MICRO_API_HANDLER=rpc \
MICRO_API_NAMESPACE=go.micro.srv \
micro api
```

Some info:

- `MICRO_API_HANDLER` sets the http handler
- `MICRO_API_NAMESPACE` sets the service namespace

### Call API

Make POST request to the api using json
```
curl -XPOST -H 'Content-Type: application/json' -d '{"name": "John"}' http://localhost:8080/example/call
```

Output
```
{"msg":"Hello John"}
```

See the [api doc](https://micro.mu/docs/api.html) for more info.

## Plugins

Micro is built on [go-micro](https://github.com/micro/go-micro) making it a pluggable toolkit.

Go-micro provides abstractions for distributed systems infrastructure which can be swapped out.

### Pluggable Features

The micro features which are pluggable:

- broker - pubsub message broker
- registry - service discovery 
- selector - client side load balancing
- transport - request-response or bidirectional streaming
- client - the client which manages the above features
- server - the server which manages the above features

Find plugins at [go-plugins](https://github.com/micro/go-plugins)

### Using Plugins

Integrate go-micro plugins by simply linking them in a separate file

Create a plugins.go file
```go
import (
	// etcd v3 registry
	_ "github.com/micro/go-plugins/registry/etcdv3"
	// nats transport
	_ "github.com/micro/go-plugins/transport/nats"
	// kafka broker
	_ "github.com/micro/go-plugins/broker/kafka"
)
```

### Building Binary

Rebuild the micro binary using the Go toolchain

```shell
# For local use
go build -i -o micro ./main.go ./plugins.go

# For docker image
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -i -o micro ./main.go ./plugins.go
```

### Enable Plugins

Enable the plugins with command line flags or env vars

```shell
# flags
micro --registry=etcdv3 --transport=nats --broker=kafka [command]

# env vars
MICRO_REGISTRY=etcdv3 MICRO_TRANSPORT=nats MICRO_BROKER=kafka micro [command]
```

