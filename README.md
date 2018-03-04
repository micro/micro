# Micro [![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![GoDoc](https://godoc.org/github.com/micro/micro?status.svg)](https://godoc.org/github.com/micro/micro) [![Travis CI](https://travis-ci.org/micro/micro.svg?branch=master)](https://travis-ci.org/micro/micro) [![Go Report Card](https://goreportcard.com/badge/micro/micro)](https://goreportcard.com/report/github.com/micro/micro)


Micro is a **microservice** toolkit. Its purpose is to simplify distributed systems development.

Check out [**go-micro**](https://github.com/micro/go-micro) if you want to start writing services in Go now or [**ja-micro**](https://github.com/Sixt/ja-micro) for Java. Examples of how to use micro with other languages can be found in [examples/sidecar](https://github.com/micro/examples/tree/master/sidecar).

Learn more about Micro in the introductory blog post [https://micro.mu/blog/2016/03/20/micro.html](https://micro.mu/blog/2016/03/20/micro.html) or watch the talk from the [Golang UK Conf 2016](https://www.youtube.com/watch?v=xspaDovwk34).

Follow us on Twitter at [@MicroHQ](https://twitter.com/microhq) or join us on [Slack](http://slack.micro.mu/).

# Overview
The goal of **Micro** is to simplify distributed systems development. Micro makes writing microservices accessible to everyone, and as you scale, micro will provide the necessary tooling to manage a microservice environment.

The toolkit is composed of the following features:

- [**`api`**](https://github.com/micro/micro/tree/master/api) - An API Gateway. A single HTTP entry point. Dynamically routing HTTP requests to RPC services.

- [**`web`**](https://github.com/micro/micro/tree/master/web) - A UI and Web Gateway. Build your web apps as micro services.

- [**`cli`**](https://github.com/micro/micro/tree/master/cli) - A command line interface. Interact with your micro services. 

- [**`bot`**](https://github.com/micro/micro/tree/master/bot) - A bot for slack and hipchat. CLI equivalent via messaging.

- [**`new`**](https://github.com/micro/micro/tree/master/new) - New template generation for services.

- [**`run`**](https://github.com/micro/micro/tree/master/run) - Runtime manager. Fetch, build and run from source in one command.

- [**`sidecar`**](https://github.com/micro/micro/tree/master/car) - A go-micro proxy. All the features of go-micro over HTTP.

## Docs

For more detailed information on the architecture, installation and use of the toolkit checkout the [docs](https://micro.mu/docs).

## Getting Started

- [Writing a Service](#writing-a-service)
- [Install Micro](#install-micro)
- [Dependencies](#dependencies)
- [Example usage](#example)
- [Build with plugins](#build-with-plugins)

### Writing a service

Learn how to write and run microservices using [**go-micro**](https://github.com/micro/go-micro). 

Read the [getting started](https://micro.mu/docs/writing-a-go-service.html) guide for more details.

### Install Micro

```shell
go get -u github.com/micro/micro
```

Or via Docker

```shell
docker pull microhq/micro
```

### Dependencies

Service discovery is the only dependency of the toolkit and go-micro. We use consul as the default.

Checkout [go-plugins](https://github.com/micro/go-plugins) to swap out consul or any other plugins.

On Mac OS
```shell
brew install consul
consul agent -dev
```

For zero dependency service discovery use the built in multicast DNS plugin.

Pass `--registry=mdns` to the below commands e.g `micro --registry=mdns list services`

## Example

Let's test out the CLI

### Run a service

This is a greeter service written with go-micro. Make sure you're running service discovery.

```shell
go get github.com/micro/examples/greeter/srv && srv
```

### List services

Each service registers with discovery so we should be able to find it.

```shell
micro list services
```

Output
```
consul
go.micro.srv.greeter
```

### Get Service

Each service has a unique id, address and metadata.

```shell
micro get service go.micro.srv.greeter
```

Output
```
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

### Query service

Make an RPC query via the CLI. The query is sent in json. We support json and protobuf out of the box.

```shell
micro query go.micro.srv.greeter Say.Hello '{"name": "John"}'
```

Output
```
{
	"msg": "Hello John"
}
```

Look at the [cli doc](https://micro.mu/docs/cli.html) for more info.

Now let's test out the micro api

### Run the api

Run the greeter API. An API service logically separates frontends from backends.

```
go get github.com/micro/examples/greeter/api && api
```

### Run the micro api

The micro api is a single HTTP entry point which dynamically routes to rpc services.

```
micro api
```

### Call via API

Replicating the CLI call as a HTTP call

```
curl http://localhost:8080/greeter/say/hello?name=John
```

Output
```
{"message":"Hello John"}
```

Look at the [api doc](https://micro.mu/docs/api.html) for more info.

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

## Learn more

To learn more read the following micro content

- [Docs](https://micro.mu/docs) - documentation and guides
- [Toolkit](https://micro.mu/blog/2016/03/20/micro.html) - intro blog post about the toolkit 
- [Architecture & Design Patterns](https://micro.mu/blog/2016/04/18/micro-architecture.html) - details on micro patterns

## Community Contributions

Project		|	Description
-----		|	------
[Micro Dashboard](https://github.com/Margatroid/micro-dashboard)	|	Dashboard for microservices toolchain micro
[Ja-Micro](https://github.com/Sixt/ja-micro)	|	A micro compatible java framework for microservices

## Sponsors

Open source development of Micro is sponsored by Sixt

<a href="https://micro.mu/blog/2016/04/25/announcing-sixt-sponsorship.html"><img src="https://micro.mu/sixt_logo.png" width=150px height="auto" /></a>

