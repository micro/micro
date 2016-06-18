# Micro [![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![GoDoc](https://godoc.org/github.com/micro/micro?status.svg)](https://godoc.org/github.com/micro/micro) [![Travis CI](https://travis-ci.org/micro/micro.svg?branch=master)](https://travis-ci.org/micro/micro) [![Roadmap](https://img.shields.io/badge/roadmap-in%20progress-lightgrey.svg)](https://github.com/micro/micro/wiki/Roadmap) [![Go Report Card](https://goreportcard.com/badge/micro/micro)](https://goreportcard.com/report/github.com/micro/micro)


Micro is a microservices toolkit. It simplifies writing and running distributed applications.

Check out [**go-micro**](https://github.com/micro/go-micro) if you want to start writing services in Go now. Examples of how to write services in other languages can be found in [examples/greeter](https://github.com/micro/micro/tree/master/examples/greeter).

Learn more about Micro in the introductory blog post [https://blog.micro.mu/2016/03/20/micro.html](https://blog.micro.mu/2016/03/20/micro.html).

Follow us on Twitter at [@MicroHQ](https://twitter.com/microhq), join the [Slack](https://micro-services.slack.com) community [here](http://micro-invites.herokuapp.com/) or 
check out the [Mailing List](https://groups.google.com/forum/#!forum/micro-services).

# Overview
The goal of **Micro** is to provide a toolkit for microservice development and management. At the core, micro is simple and accessible enough that anyone can easily get started writing microservices. As you scale to hundreds of services, micro will provide the fundamental tools required to manage a microservice environment.

**Go Micro** - A pluggable RPC framework for writing microservices in Go. It provides libraries for 
service discovery, client side load balancing, encoding, synchronous and asynchronous communication.

**API** - An API Gateway that serves HTTP and routes requests to appropriate micro services. 
It acts as a single entry point and can either be used as a reverse proxy or translate HTTP requests to RPC.

**Web** - A web dashboard and reverse proxy for micro web applications. We believe that 
web apps should be built as microservices and therefore treated as a first class citizen in a microservice world. It behaves much the like the API 
reverse proxy but also includes support for web sockets.

**Sidecar** - The Sidecar provides all the features of go-micro as a HTTP service. While we love Go and 
believe it's a great language to build microservices, you may also want to use other languages, so the Sidecar provides a way to integrate 
your other apps into the Micro world.

**CLI** - A straight forward command line interface to interact with your micro services. 
It also allows you to leverage the Sidecar as a proxy where you may not want to directly connect to the service registry.

**Bot** A Hubot style bot that sits inside your microservices platform and can be interacted with via Slack, HipChat, XMPP, etc. 
It provides the features of the CLI via messaging. Additional commands can be added to automate common ops tasks.

## Getting Started

### Writing a service

Learn how to write and run a microservice using [**go-micro**](https://github.com/micro/go-micro). 

Read the [Getting Started](https://github.com/micro/micro/wiki/Getting-Started) guide or the blog post on 
[Writing microservices with Go-Micro](https://blog.micro.mu/2016/03/28/go-micro.html).

### Install Micro

```shell
$ go get github.com/micro/micro
```

Or via Docker

```shell
$ docker pull microhq/micro
```

### Quick start

We need service discovery, so let's spin up Consul (default discovery mechanism; checkout [go-plugins](https://github.com/micro/go-plugins) to switch it out).

```
$ go get github.com/hashicorp/consul
$ consul agent -dev -advertise=127.0.0.1
```

Alternatively you can use multicast DNS with the built in MDNS registry. Just pass `--registry=mdns` to the below commands e.g. `server --registry=mdns` or `micro --registry=mdns list services`.

Run the greeter example app
```shell
$ go get github.com/micro/micro/examples/greeter/server
$ server
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
go.micro.srv.greeter

Id	Address	Port	Metadata
go.micro.srv.greeter-154a6487-7d7e-11e5-882a-34363b77bace	[::]	57067	

Endpoint: Say.Hello
Metadata: stream=false
Request:
{
	name string
}
Response:
{
	msg string
}

Endpoint: Debug.Health
Metadata: stream=false
Request: {}
Response:
{
	status string
}
```

Query service
```shell
$ micro query go.micro.srv.greeter Say.Hello '{"name": "John"}'
{
	"msg": "go.micro.srv.greeter-154a6487-7d7e-11e5-882a-34363b77bace: Hello John"
}

```

Read more on how to use the Micro [CLI](https://github.com/micro/micro/tree/master/cli)

## The Ecosystem

The overarching project [github.com/micro](https://github.com/micro) is a microservice ecosystem which consists of a number of tools and libraries. Each of which can either be used totally independently, plugged into your architecture or combined as a whole to provide a completely distributed systems platform.

It currently consists of the following.

### Go Micro

[Go-micro](https://github.com/micro/go-micro) is a pluggable Go client framework for writing microservices. Go-micro can be used standalone and should be the starting point for writing applications.

Feature		|	Description
-------		|	-----------
[Registry](https://godoc.org/github.com/micro/go-micro/registry)	|	Service discovery
[Client](https://godoc.org/github.com/micro/go-micro/client)	|	RPC Client
[Codec](https://godoc.org/github.com/micro/go-micro/codec)	|	Request/Response Encoding
[Selector](https://godoc.org/github.com/micro/go-micro/selector)	|	Load balancing 
[Server](https://godoc.org/github.com/micro/go-micro/server)	|	RPC Server
[Broker](https://godoc.org/github.com/micro/go-micro/broker)	|	Asynchronous Messaging
[Transport](https://godoc.org/github.com/micro/go-micro/transport)	|	Synchronous Messaging

### Micro

[Micro](https://github.com/micro/micro) provides entry points into a running system with an API Gateway, Web UI, HTTP Sidecar and CLI. Micro can be used to manage the public facing aspect of your services and will normally run at the edge of your infrastructure.

Feature		|	Description
------		|	-------
[API Gateway](https://github.com/micro/micro/tree/master/api) | Lightweight gateway/proxy. Convert http requests to rpc
[CLI](https://github.com/micro/micro/tree/master/cli) | Command line interface
[Sidecar](https://github.com/micro/micro/tree/master/car) | HTTP proxy for non Go-micro apps
[Web UI/Proxy](https://github.com/micro/micro/tree/master/web) | A visual way to view and query services
[Bot](https://github.com/micro/micro/tree/master/bot) | A bot that sits inside Slack, HipChat, etc for ChatOps

### Go Platform

[Go-platform](https://github.com/micro/go-platform) provides pluggable libraries for integrating with higher level requirements for microservices. 
It mainly integrates functionality for distributed systems.

Feature     |   Description
-------     |   ---------
[auth](https://godoc.org/github.com/micro/go-platform/auth)	|   authentication and authorisation for users and services
[config](https://godoc.org/github.com/micro/go-platform/config)	|   dynamic configuration which is namespaced and versioned
[db](https://godoc.org/github.com/micro/go-platform/db)		| distributed database abstraction
[discovery](https://godoc.org/github.com/micro/go-platform/discovery)	|   extends the go-micro registry to add heartbeating, etc
[event](https://godoc.org/github.com/micro/go-platform/event)	|	platform event publication, subscription and aggregation 
[kv](https://godoc.org/github.com/micro/go-platform/kv)		|   simply key value layered on memcached, etcd, consul 
[log](https://godoc.org/github.com/micro/go-platform/log)	|	structured logging to stdout, logstash, fluentd, pubsub
[monitor](https://godoc.org/github.com/micro/go-platform/monitor)	|   add custom healthchecks measured with distributed systems in mind
[metrics](https://godoc.org/github.com/micro/go-platform/metrics)	|   instrumentation and collation of counters
[router](https://godoc.org/github.com/micro/go-platform/router)	|	global circuit breaking, load balancing, A/B testing
[sync](https://godoc.org/github.com/micro/go-platform/sync)	|	distributed locking, leadership election, etc
[trace](https://godoc.org/github.com/micro/go-platform/trace)	|	distributed tracing of request/response

### Platform

[Platform](https://github.com/micro/platform) is a complete runtime for managing microservices at scale. Where Micro provides the core essentials, the platform goes a step further and addresses every requirement for large scale distributed system deployments. 

Feature		|	Description
------------	|	-------------
[Auth](https://github.com/micro/auth-srv)	|	Authentication and authorization (Oauth2)
[Config](https://github.com/micro/config-srv)	|	Dynamic configuration
[DB Proxy](https://github.com/micro/db-srv)	|	RPC based database proxy
[Discovery](https://github.com/micro/discovery-srv)	|	Service discovery read layer cache
[Events](https://github.com/micro/event-srv)	|	Platform event aggregation
[Monitoring](https://github.com/micro/monitor-srv)	|	Monitoring for Status, Stats and Healthchecks
[Routing](https://github.com/micro/router-srv)	|	Global service load balancing
[Tracing](https://github.com/micro/trace-srv)	|	Distributed tracing

### Go Plugins

[Go Plugins](https://github.com/micro/go-plugins) provides plugins for go-micro and go-platform contributed by the community. Examples could include; circuit breakers, rate limiting. Registries built on top of Kubernetes, Zookeeper, etc. Transport using HTTP2, Zeromq, etc. Broker using Kafka, AWS SQS, etc.

Example plugins

Plugin	|	Description
-----	|	------
[NATS](https://godoc.org/github.com/micro/go-plugins/transport/nats)	|	Synchronous transport with the NATS message bus
[Etcd](https://godoc.org/github.com/micro/go-plugins/registry/etcd)	|	Service discovery using etcd
[BSON-RPC](https://godoc.org/github.com/micro/go-plugins/codec/bsonrpc)	|	Request/Response encoding using bson-rpc

## Example Services
Project		|	Description
-----		|	------
[greeter](https://github.com/micro/micro/tree/master/examples/greeter)	|	A greeter service (includes Go, Ruby, Python examples)
[geo-srv](https://github.com/micro/geo-srv)	|	Geolocation tracking service using hailocab/go-geoindex
[geo-api](https://github.com/micro/geo-api)	|	A HTTP API handler for geo location tracking and search
[discovery-srv](https://github.com/micro/discovery-srv)	|	A discovery in the micro platform
[geocode-srv](https://github.com/micro/geocode-srv)	|	A geocoding service using the Google Geocoding API
[hailo-srv](https://github.com/micro/hailo-srv)	|	A service for the hailo taxi service developer api
[monitor-srv](https://github.com/micro/monitor-srv)	|	A monitoring service for Micro services
[place-srv](https://github.com/micro/place-srv)	|	A microservice to store and retrieve places (includes Google Place Search API)
[slack-srv](https://github.com/micro/slack-srv)	|	The slack bot API as a go-micro RPC service
[trace-srv](https://github.com/micro/trace-srv)	|	A distributed tracing microservice in the realm of dapper, zipkin, etc
[twitter-srv](https://github.com/micro/twitter-srv)	|	A microservice for the twitter API
[user-srv](https://github.com/micro/user-srv)	|	A microservice for user management and authentication

## Community Contributions

Project		|	Description
-----		|	------
[Micro Dashboard](https://github.com/Margatroid/micro-dashboard)	|	Dashboard for microservices toolchain micro

## Sponsors

<a href="https://www.sixt.com"><img src="https://micro.mu/sixt_logo.png" width=150px height="auto" /></a>

## Roadmap

[![Roadmap](https://img.shields.io/badge/roadmap-in%20progress-lightgrey.svg)](https://github.com/micro/micro/wiki/Roadmap)

## License

Apache 2.0

## Contributing

1. [Join](http://slack.micro.mu/) the Slack to discuss
2. Look at existing coding style
3. Submit PR
4. ?
5. Profit

We're looking for implementations equivalent to [go-micro](https://github.com/micro/go-micro) in different languages. 
Come join the Slack to discuss.
