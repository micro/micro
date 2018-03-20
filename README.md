# Micro [![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![GoDoc](https://godoc.org/github.com/micro/micro?status.svg)](https://godoc.org/github.com/micro/micro) [![Travis CI](https://travis-ci.org/micro/micro.svg?branch=master)](https://travis-ci.org/micro/micro) [![Go Report Card](https://goreportcard.com/badge/micro/micro)](https://goreportcard.com/report/github.com/micro/micro)

Micro is a toolkit for cloud-native development.

# Overview

Micro addresses the key requirements for building cloud-native systems. It takes the microservice architecture pattern and transforms it into 
a set of tools which act as the building blocks for scalable platforms. Micro hides the complexity of distributed systems and provides 
well understood concepts to developers.

Micro builds on [go-micro](https://github.com/micro/go-micro), making it a pluggable toolkit.

## Features

The toolkit is composed of the following features:

- [**`api`**](https://github.com/micro/micro/tree/master/api) - API Gateway. A single HTTP entry point. Dynamic routing using service discovery.

- [**`cli`**](https://github.com/micro/micro/tree/master/cli) - Command line interface. Describe, query and interact directly from the terminal. 

- [**`bot`**](https://github.com/micro/micro/tree/master/bot) - Slack and hipchat bot. The CLI via messaging.

- [**`new`**](https://github.com/micro/micro/tree/master/new) - New template generation for services.

- [**`web`**](https://github.com/micro/micro/tree/master/web) - Web dashboard to interact via browser.

- [**`proxy`**](https://github.com/micro/micro/tree/master/proxy) - A cli proxy for remote environments.

## Docs

For more detailed information on the architecture, installation and use of the toolkit checkout the [docs](https://micro.mu/docs).

## Getting Started

- [Writing a Service](#writing-a-service)
- [Install Micro](#install-micro)
- [Service Discovery](#service-discovery)
- [Example usage](#example)
- [Build with plugins](#build-with-plugins)

## Writing a service

See [**go-micro**](https://github.com/micro/go-micro) to start writing services.

## Install Micro

```shell
go get -u github.com/micro/micro
```

Or via Docker

```shell
docker pull microhq/micro
```

## Service Discovery

Service discovery is the only dependency of the micro toolkit. Consul is set as the default.

Internally micro uses the [go-micro](https://github.com/micro/go-micro) registry for service discovery. This allows the toolkit to leverage 
go-micro plugins. Checkout [go-plugins](https://github.com/micro/go-plugins) to swap out consul.

### Consul

On Mac OS
```shell
brew install consul
consul agent -dev
```

### mDNS

Multicast DNS is built in for zero dependency service discovery.

Pass `--registry=mdns` or set the env var `MICRO_REGISTRY=mdns` for any command

```
## Use flag
micro --registry=mdns list services

## Use env var
MICRO_REGISTRY=mdns micro list services`
```

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

### Call service

Make an RPC call via the CLI. The query is sent in json. We support json and protobuf out of the box.

```shell
micro call go.micro.srv.greeter Say.Hello '{"name": "John"}'
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

## Plugins

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

Build binary
```shell
# For local use
go build -i -o micro ./main.go ./plugins.go

# For docker image
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -i -o micro ./main.go ./plugins.go
```

Enable with flags or env vars
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


