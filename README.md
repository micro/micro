# Micro [![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![GoDoc](https://godoc.org/github.com/micro/micro?status.svg)](https://godoc.org/github.com/micro/micro) [![Travis CI](https://travis-ci.org/micro/micro.svg?branch=master)](https://travis-ci.org/micro/micro) [![Roadmap](https://img.shields.io/badge/roadmap-in%20progress-lightgrey.svg)](https://github.com/micro/docs/blob/master/roadmap.md) [![Go Report Card](https://goreportcard.com/badge/micro/micro)](https://goreportcard.com/report/github.com/micro/micro)


Micro is a **microservice** toolkit. Its purpose is to simplify distributed systems development.

Check out [**go-micro**](https://github.com/micro/go-micro) if you want to start writing services in Go now. Examples of how to use micro with other languages can be found in [examples/sidecar](https://github.com/micro/examples/tree/master/sidecar).

Learn more about Micro in the introductory blog post [https://blog.micro.mu/2016/03/20/micro.html](https://blog.micro.mu/2016/03/20/micro.html) or watch the talk from the [Golang UK Conf 2016](https://www.youtube.com/watch?v=xspaDovwk34).

Follow us on Twitter at [@MicroHQ](https://twitter.com/microhq), join the [Slack](https://micro-services.slack.com) community [here](http://slack.micro.mu/) or 
check out the [Mailing List](https://groups.google.com/forum/#!forum/microhq).

# Overview
The goal of **Micro** is to simplify distributed systems development. At the core micro should be accessible enough to anyone who wants to get started writing microservices. As you scale to hundreds of services, micro will provide the necessary tooling to manage a microservice environment.

The toolkit is composed of the following components:

- **Go Micro** - A pluggable RPC framework for writing microservices in Go. It provides libraries for 
service discovery, client side load balancing, encoding, synchronous and asynchronous communication.

- **API** - An API Gateway that serves HTTP and routes requests to appropriate micro services. 
It acts as a single entry point and can either be used as a reverse proxy or translate HTTP requests to RPC.

- **Web** - A web dashboard and reverse proxy for micro web applications. We believe that 
web apps should be built as microservices and therefore treated as a first class citizen in a microservice world. It behaves much the like the API 
reverse proxy but also includes support for web sockets.

- **Sidecar** - A language agnostic RPC proxy which provides all the features of go-micro as a HTTP endpoints. While we love Go and 
believe it's a great language to build microservices, you may also want to use other languages, so the Sidecar provides a way to integrate 
your other apps into the Micro world. It can be used as a standalone way to build fault tolerant microservices.

- **CLI** - A straight forward command line interface to interact with your micro services. 
It also allows you to leverage the Sidecar as a proxy where you may not want to directly connect to the service registry.

- **Bot** - A Hubot style bot that sits inside your microservices platform and can be interacted with via Slack, HipChat, XMPP, etc. 
It provides the features of the CLI via messaging. Additional commands can be added to automate common ops tasks.

Read the [docs](https://github.com/micro/docs) for more details.

## Getting Started

### Writing a service

Learn how to write and run a microservice using [**go-micro**](https://github.com/micro/go-micro). 

Read the [Getting Started](https://github.com/micro/docs/blob/master/getting-started.md) guide or the blog post on 
[Writing microservices with Go-Micro](https://blog.micro.mu/2016/03/28/go-micro.html).

### Install Micro

```shell
go get github.com/micro/micro
```

Or via Docker

```shell
docker pull microhq/micro
```

### Quick start

We need service discovery, so let's spin up Consul (the default); checkout [go-plugins](https://github.com/micro/go-plugins) to swap it out.

Install consul
```
brew install consul
```

Run consul
```
consul agent -dev -advertise=127.0.0.1
```

Or we can use multicast DNS for zero dependency discovery

Pass `--registry=mdns` to the below commands e.g `micro --registry=mdns list services`

Run the greeter service

```shell
go get github.com/micro/examples/greeter/srv && srv
```

List services
```shell
$ micro list services
consul
go.micro.srv.greeter
```

Get Service
```shell
$ micro get service go.micro.srv.greeter
service  go.micro.srv.greeter

version 1.0.0

Id	Address	Port	Metadata
go.micro.srv.greeter-34c55534-368b-11e6-b732-68a86d0d36b6	192.168.1.66	62525	server=rpc,registry=consul,transport=http,broker=http

Endpoint: Say.Hello
Metadata: stream=false

Request: {
	name string
}

Response: {
	msg string
}
```

Query service
```shell
$ micro query go.micro.srv.greeter Say.Hello '{"name": "John"}'
{
	"msg": "Hello John"
}
```

Read the [docs](https://github.com/micro/docs) to learn more.

## Build with plugins

If you want to integrate plugins simply link them in a separate file and rebuild

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

Build binary
```shell
// For local use
go build -i -o micro ./main.go ./plugins.go

// For docker image
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -i -o micro ./main.go ./plugins.go
```

Flag usage of plugins
```shell
micro --registry=etcdv3 --transport=nats --broker=kafka
```

## Community Contributions

Project		|	Description
-----		|	------
[Micro Dashboard](https://github.com/Margatroid/micro-dashboard)	|	Dashboard for microservices toolchain micro

## Sponsors

Open source development of Micro is sponsored by Sixt

<a href="https://blog.micro.mu/2016/04/25/announcing-sixt-sponsorship.html"><img src="https://micro.mu/sixt_logo.png" width=150px height="auto" /></a>

