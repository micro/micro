# Micro [![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![GoDoc](https://godoc.org/github.com/micro/micro?status.svg)](https://godoc.org/github.com/micro/micro) [![Travis CI](https://travis-ci.org/micro/micro.svg?branch=master)](https://travis-ci.org/micro/micro) [![Go Report Card](https://goreportcard.com/badge/micro/micro)](https://goreportcard.com/report/github.com/micro/micro)

Micro is a toolkit for cloud-native development. It helps you build future-proof application platforms and services.

# Overview

Micro addresses the key requirements for building cloud-native systems. It takes the microservice architecture pattern and transforms it into 
a set of tools which act as the building blocks for scalable platforms. Micro deals with the complexity of distributed systems and provides 
simple abstractions already understood by developers.

Technology is constantly evolving. The infrastructure stack is always changing. Micro is a pluggable toolkit which addresses these issues. 
Plug in any stack or underlying technology. Build future-proof systems using micro.

## Features

The toolkit is composed of the following features:

- [**`api`**](https://github.com/micro/micro/tree/master/api) - API Gateway. A single HTTP entry point. Dynamic routing using service discovery.

- [**`bot`**](https://github.com/micro/micro/tree/master/bot) - Slack and hipchat bot. The CLI via messaging.

- [**`cli`**](https://github.com/micro/micro/tree/master/cli) - Command line interface. Describe, query and interact directly from the terminal. 

- [**`new`**](https://github.com/micro/micro/tree/master/new) - Service template generation. Get started quickly.

- [**`web`**](https://github.com/micro/micro/tree/master/web) - Web dashboard to interact via browser.

## Docs

For more detailed information on the architecture, installation and use of the toolkit checkout the [docs](https://micro.mu/docs).

## Getting Started

- [Install Micro](#install-micro)
- [Dependencies](#dependencies)
- [Writing a Service](#writing-a-service)
- [Example usage](#example)
- [Plugins](#plugins)
- [Learn More](#learn-more)
- [Community Projects](#community-projects)

## Install Micro

```shell
go get -u github.com/micro/micro
```

Or via Docker

```shell
docker pull microhq/micro
```

## Dependencies

The micro toolkit has two dependencies: 

- [Service Discovery](#service-discovery) - used for name resolution
- [Protobuf](#protobuf) - used for code generation

## Service Discovery

Service discovery is used for name resolution, routing and centralising metadata.

Micro uses the [go-micro](https://github.com/micro/go-micro) registry for service discovery. Consul is the default registry. 

See [go-plugins](https://github.com/micro/go-plugins) to swap out consul.

### Consul

Install and run consul

```shell
# install
brew install consul

# run
consul agent -dev
```

### mDNS

Multicast DNS is an alternative built in registry for zero dependency service discovery.

Pass `--registry=mdns` or set the env var `MICRO_REGISTRY=mdns` for any command

```shell
# Use flag
micro --registry=mdns list services

# Use env var
MICRO_REGISTRY=mdns micro list services`
```

## Protobuf

Protobuf is used for code generation. It reduces the amount of boilerplate code needed to be written.

```
# install protobuf
brew install protobuf

# install protoc-gen-go
go get -u github/golang/protobuf/{proto,protoc-gen-go}

# install protoc-gen-micro
go get -u github.com/micro/protoc-gen-micro
```

See [protoc-gen-micro](https://github.com/micro/protoc-gen-micro) for more details.

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
go run main.go
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
consul
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
go.micro.srv.example-437d1277-303b-11e8-9be9-f40f242f6897	192.168.1.65	53545	transport=http,broker=http,server=rpc,registry=consul

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

Now let's test out the micro api

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

## Learn more

To learn more read the following micro content

- [Docs](https://micro.mu/docs) - documentation and guides
- [Toolkit](https://micro.mu/blog/2016/03/20/micro.html) - intro blog post about the toolkit 
- [Architecture & Design Patterns](https://micro.mu/blog/2016/04/18/micro-architecture.html) - details on micro patterns
- [Presentation](https://www.youtube.com/watch?v=xspaDovwk34) - Golang UK Conf 2016
- [Twitter](https://twitter.com/microhq) - follow us on Twitter
- [Slack](http://slack.micro.mu/) - join the slack community (1000+ members)

## Community Projects

Project		|	Description
-----		|	------
[Dashboard](https://github.com/Margatroid/micro-dashboard)	|	A react based micro dashboard
[Ja-micro](https://github.com/Sixt/ja-micro)	|	A micro compatible java framework

Explore other projects at [micro.mu/explore](https://micro.mu/explore/)

## Sponsors

Open source development of Micro is sponsored by Sixt

<a href="https://micro.mu/blog/2016/04/25/announcing-sixt-sponsorship.html"><img src="https://micro.mu/sixt_logo.png" width=150px height="auto" /></a>


