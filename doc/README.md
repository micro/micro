The goal of **Micro** is to provide an ecosystem of tools and services for microservice development and management. At the core, the toolkit is simple and accessible enough that anyone can easily get started writing microservices. As you scale to hundreds of services, micro will provide the fundamental tools required to manage a microservice environment.

Checkout the [roadmap](https://github.com/micro/micro/wiki/Roadmap) to see where it's all going.

# Overview

Here's a further breakdown of the main toolkit.

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

## Features

Below are the libraries and services that encompass the micro ecosystem. If you want a more detailed overview of Micro, checkout the introductory blog post [https://blog.micro.mu/2016/03/20/micro.html](https://blog.micro.mu/2016/03/20/micro.html).

### Go Micro
[Go-micro](https://github.com/micro/go-micro) is a pluggable Go framework for writing microservices.

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

### Go OS

[Go-OS](https://github.com/micro/go-os) provides pluggable libraries for integrating with higher level requirements for microservices. 
It mainly integrates functionality for distributed systems.

Feature     |   Description
-------     |   ---------
[auth](https://godoc.org/github.com/micro/go-os/auth)	|   authentication and authorisation for users and services
[config](https://godoc.org/github.com/micro/go-os/config)	|   dynamic configuration which is namespaced and versioned
[db](https://godoc.org/github.com/micro/go-os/db)		| distributed database abstraction
[discovery](https://godoc.org/github.com/micro/go-os/discovery)	|   extends the go-micro registry to add heartbeating, etc
[event](https://godoc.org/github.com/micro/go-os/event)	|	event publication, subscription and aggregation 
[kv](https://godoc.org/github.com/micro/go-os/kv)		|   simply key value layered on memcached, etcd, consul 
[log](https://godoc.org/github.com/micro/go-os/log)	|	structured logging to stdout, logstash, fluentd, pubsub
[monitor](https://godoc.org/github.com/micro/go-os/monitor)	|   add custom healthchecks measured with distributed systems in mind
[metrics](https://godoc.org/github.com/micro/go-os/metrics)	|   instrumentation and collation of counters
[router](https://godoc.org/github.com/micro/go-os/router)	|	global circuit breaking, load balancing, A/B testing
[sync](https://godoc.org/github.com/micro/go-os/sync)	|	distributed locking, leadership election, etc
[trace](https://godoc.org/github.com/micro/go-os/trace)	|	distributed tracing of request/response

### Micro OS

[Micro OS](https://github.com/micro/os) is a complete runtime for managing microservices at scale. Where Micro provides the core essentials, Micro OS goes a step further and addresses every requirement for large scale distributed systems. 

Feature		|	Description
------------	|	-------------
[Auth](https://github.com/micro/auth-srv)	|	Authentication and authorization (Oauth2)
[Config](https://github.com/micro/config-srv)	|	Dynamic configuration
[DB Proxy](https://github.com/micro/db-srv)	|	RPC based database proxy
[Discovery](https://github.com/micro/discovery-srv)	|	Service discovery read layer cache
[Events](https://github.com/micro/event-srv)	|	Event aggregation
[Monitoring](https://github.com/micro/monitor-srv)	|	Monitoring for Status, Stats and Healthchecks
[Routing](https://github.com/micro/router-srv)	|	Global service load balancing
[Tracing](https://github.com/micro/trace-srv)	|	Distributed tracing

### Go Plugins

[Go Plugins](https://github.com/micro/go-plugins) provides plugins for go-micro and go-os contributed by the community. Examples could include; circuit breakers, rate limiting. Registries built on top of Kubernetes, Zookeeper, etc. Transport using HTTP2, Zeromq, etc. Broker using Kafka, AWS SQS, etc.

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
[discovery-srv](https://github.com/micro/discovery-srv)	|	A discovery in the micro OS
[geocode-srv](https://github.com/micro/geocode-srv)	|	A geocoding service using the Google Geocoding API
[hailo-srv](https://github.com/micro/hailo-srv)	|	A service for the hailo taxi service developer api
[monitor-srv](https://github.com/micro/monitor-srv)	|	A monitoring service for Micro services
[place-srv](https://github.com/micro/place-srv)	|	A microservice to store and retrieve places (includes Google Place Search API)
[slack-srv](https://github.com/micro/slack-srv)	|	The slack bot API as a go-micro RPC service
[trace-srv](https://github.com/micro/trace-srv)	|	A distributed tracing microservice in the realm of dapper, zipkin, etc
[twitter-srv](https://github.com/micro/twitter-srv)	|	A microservice for the twitter API
[user-srv](https://github.com/micro/user-srv)	|	A microservice for user management and authentication

