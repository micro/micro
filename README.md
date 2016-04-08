# Micro [![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![GoDoc](https://godoc.org/github.com/micro/micro?status.svg)](https://godoc.org/github.com/micro/micro) [![Travis CI](https://travis-ci.org/micro/micro.svg?branch=master)](https://travis-ci.org/micro/micro) [![Roadmap](https://img.shields.io/badge/roadmap-in%20progress-lightgrey.svg)](https://github.com/micro/micro/wiki/Roadmap) [![Go Report Card](https://goreportcard.com/badge/micro/micro)](https://goreportcard.com/report/github.com/micro/micro)

Micro is a microservices toolkit. It simplifies writing and running distributed applications.

Check out [**go-micro**](https://github.com/micro/go-micro) if you want to start writing services in Go now. Examples of how to write services in other languages can be found in [examples/greeter](https://github.com/micro/micro/tree/master/examples/greeter).

Learn more about Micro in the introductory blog post [https://blog.micro.mu/2016/03/20/micro.html](https://blog.micro.mu/2016/03/20/micro.html).

Follow us on Twitter at [@MicroHQ](https://twitter.com/microhq), join the [Slack](https://micro-services.slack.com) community [here](http://micro-invites.herokuapp.com/) or 
check out the [Mailing List](https://groups.google.com/forum/#!forum/micro-services).

# Overview
The goal of **Micro** is to provide a toolkit for microservice development and management. At the core, micro is simple and accessible enough that anyone can easily get started writing microservices. As you scale to hundreds of services, micro will provide the fundamental tools required to manage a microservice environment.

![Micro](https://github.com/micro/micro/blob/master/doc/micro.png)

## Getting Started

### Writing a service

Learn how to write and run a microservice using [**go-micro**](https://github.com/micro/go-micro). 
Read the [Getting Started](https://github.com/micro/micro/wiki/Getting-Started) guide.

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

### Usage

```shell
NAME:
   micro - A microservices toolkit

USAGE:
   micro [global options] command [command options] [arguments...]
   
VERSION:
   latest
   
COMMANDS:
   api		Run the micro API
   bot		Run the micro bot
   registry	Query registry
   query	Query a service method using rpc
   stream	Query a service method using streaming rpc
   health	Query the health of a service
   list		List items in registry
   register	Register an item in the registry
   deregister	Deregister an item in the registry
   get		Get item from registry
   sidecar	Run the micro sidecar
   web		Run the micro web app
   help, h	Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --server_name 								Name of the server. go.micro.srv.example [$MICRO_SERVER_NAME]
   --server_version 								Version of the server. 1.1.0 [$MICRO_SERVER_VERSION]
   --server_id 									Id of the server. Auto-generated if not specified [$MICRO_SERVER_ID]
   --server_address 								Bind address for the server. 127.0.0.1:8080 [$MICRO_SERVER_ADDRESS]
   --server_advertise 								Used instead of the server_address when registering with discovery. 127.0.0.1:8080 [$MICRO_SERVER_ADVERTISE]
   --server_metadata [--server_metadata option --server_metadata option]	A list of key-value pairs defining metadata. version=1.0.0 [$MICRO_SERVER_METADATA]
   --broker 									Broker for pub/sub. http, nats, rabbitmq [$MICRO_BROKER]
   --broker_address 								Comma-separated list of broker addresses [$MICRO_BROKER_ADDRESS]
   --registry 									Registry for discovery. memory, consul, etcd, kubernetes [$MICRO_REGISTRY]
   --registry_address 								Comma-separated list of registry addresses [$MICRO_REGISTRY_ADDRESS]
   --selector 									Selector used to pick nodes for querying. random, roundrobin, blacklist [$MICRO_SELECTOR]
   --transport 									Transport mechanism used; http, rabbitmq, nats [$MICRO_TRANSPORT]
   --transport_address 								Comma-separated list of transport addresses [$MICRO_TRANSPORT_ADDRESS]
   --enable_tls									Enable TLS [$MICRO_ENABLE_TLS]
   --tls_cert_file 								TLS Certificate file [$MICRO_TLS_CERT_File]
   --tls_key_file 								TLS Key file [$MICRO_TLS_KEY_File]
   --api_address 								Set the api address e.g 0.0.0.0:8080 [$MICRO_API_ADDRESS]
   --proxy_address 								Proxy requests via the HTTP address specified [$MICRO_PROXY_ADDRESS]
   --sidecar_address 								Set the sidecar address e.g 0.0.0.0:8081 [$MICRO_SIDECAR_ADDRESS]
   --web_address 								Set the web UI address e.g 0.0.0.0:8082 [$MICRO_WEB_ADDRESS]
   --register_ttl "0"								Register TTL in seconds [$MICRO_REGISTER_TTL]
   --register_interval "0"							Register interval in seconds [$MICRO_REGISTER_INTERVAL]
   --api_handler 								Specify the request handler to be used for mapping HTTP requests to services. e.g api, proxy [$MICRO_API_HANDLER]
   --api_namespace 								Set the namespace used by the API e.g. com.example.api [$MICRO_API_NAMESPACE]
   --web_namespace 								Set the namespace used by the Web proxy e.g. com.example.web [$MICRO_WEB_NAMESPACE]
   --api_cors 									Comma separated whitelist of allowed origins for CORS [$MICRO_API_CORS]
   --web_cors 									Comma separated whitelist of allowed origins for CORS [$MICRO_WEB_CORS]
   --sidecar_cors 								Comma separated whitelist of allowed origins for CORS [$MICRO_SIDECAR_CORS]
   --enable_stats								Enable stats [$MICRO_ENABLE_STATS]
   --help, -h									show help
```

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

## Architecture

![Overview1](https://github.com/micro/micro/blob/master/doc/overview1.png)
-

![Overview2](https://github.com/micro/micro/blob/master/doc/overview2.png)
-

![Overview3](https://github.com/micro/micro/blob/master/doc/overview3.png)
-

### Bot

ChatOps as a first class citizen in the Micro world

<img src="https://github.com/micro/micro/blob/master/bot/slack.png">

### Dashboard

<img src="https://github.com/micro/micro/blob/master/web/web1.png">

-
<img src="https://github.com/micro/micro/blob/master/web/web2.png">

-
<img src="https://github.com/micro/micro/blob/master/web/web3.png">

-
<img src="https://github.com/micro/micro/blob/master/doc/stats.png">

## About Microservices
Microservices is an architecture pattern used to decompose a single large application in to a smaller suite of services. Generally the goal is to create light weight services of 1000 lines of code or less. Each service alone provides a particular focused solution or set of solutions. These small services can be used as the foundational building blocks in the creation of a larger system.

The concept of microservices is not new, this is the reimagination of service orientied architecture but with an approach more holistically aligned with unix processes and pipes. For those of us with extensive experience in this field we're somewhat biased and feel this is an incredibly beneficial approach to system design at large and developer productivity.

Learn more about Microservices by watching Martin Fowler's presentation [here](https://www.youtube.com/watch?v=wgdBVIX9ifA) or his blog post [here](http://martinfowler.com/articles/microservices.html).

## Microservice Requirements

The foundation of a library enabling microservices is based around the following requirements:

- Server - an ability to define handlers and serve requests 
- Client - an ability to make requests to another service
- Discovery - a mechanism by which to discover other services

These 3 components form the minimum requirements for microservices development. An ecosystem of libraries and tools can be created around them to provide a feature rich system however at the foundation only these 3 things are required to write services and communicate between them.

### Server

The server is the core component which allows you to register request handlers and serve requests. Ideally it's transport agnostic so different transports such as http, rabbitmq, etc can be chosen. On start it should register itself with discovery system so other microservices know it exists and deregister when shutting down. The server should handle encoding/decoding incoming/outgoing requests, leaving the handlers to operate on the request/response types they expect.

Example interface:
```
server.New(name, options) - instantiate new server
server.Register(handler) - register a handler with the server
server.Start() - start
server.Stop() - stop
```
### Client

Where the server allows you to serve requests, the client lets you make them to other servers. The client should support request/response and pub/sub. Part of the microservices world is event driven programming, taking action based on events, which is why pub/sub is a requirement of the client. It should also make use of the discovery system so requests can be made by service name. 

Example interface:
```
client.Request(name, request) - Make a request to another server
client.Publish(topic, message) - Publish a message on a topic
client.Subscribe(topic, channel) - Subscribe to a topic
```
### Discovery

The discovery system is really vital to microservices development. Any sort of communication between servers will first require locating it and then making the request. Discovery should support registration and retrieval of servers. It should optionally support a keepalive mechanism to remove stale servers.

Example interface:
```
discovery.Register(name, hostname, ...) - Register a server
discovery.Deregister(name, hostname, ...) - Deregister a server
discovery.Get(name) - Get the details for a server
discovery.List() - List all servers
```

### Microservices? What are they even good for?

![Micro On-Demand](https://github.com/micro/micro/blob/master/doc/ondemand.png)

## Resources

- [The Micro Blog](https://blog.micro.mu)
- [A Journey into Microservices](https://sudo.hailoapp.com/services/2015/03/09/journey-into-a-microservice-world-part-1/)
- [A Journey into a Microservice World](https://speakerdeck.com/mattheath/a-journey-into-a-microservice-world) by Matt Heath (Slides)
- [Microservices](http://martinfowler.com/articles/microservices.html) by Martin Fowler
- [Microservices: Decomposing Applications for Deployability and Scalability](http://www.slideshare.net/chris.e.richardson/microservices-decomposing-applications-for-deployability-and-scalability-jax) by Chris Richardson (Slides)
- [4 reasons why microservices resonate](http://radar.oreilly.com/2015/04/4-reasons-why-microservices-resonate.html) by Neal Ford

## Roadmap

[![Roadmap](https://img.shields.io/badge/roadmap-in%20progress-lightgrey.svg)](https://github.com/micro/micro/wiki/Roadmap)

