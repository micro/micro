# Micro [![GoDoc](https://godoc.org/github.com/micro/micro?status.svg)](https://godoc.org/github.com/micro/micro) [![Travis CI](https://travis-ci.org/micro/micro.svg?branch=master)](https://travis-ci.org/micro/micro) [![Roadmap](https://img.shields.io/badge/roadmap-in%20progress-lightgrey.svg)](https://github.com/micro/micro/wiki/Roadmap)

Micro is a microservices toolkit. It simplifies writing and running distributed applications.

Checkout [**go-micro**](https://github.com/micro/go-micro) if you want to start writing services now.

Examples of how to write a service in ruby or python can be found in [here](https://github.com/micro/micro/tree/master/examples/greeter)

- [Mailing List](https://groups.google.com/forum/#!forum/micro-services) 
- [Slack](https://micro-services.slack.com) : [auto-invite](http://micro-invites.herokuapp.com/)

# Overview
The goal of **Micro** is to provide a toolkit for microservice development and management. At the core, micro is simple and accessible enough that anyone can easily get started writing microservices. As you scale to hundreds of services, micro will provide the fundamental tools required to manage a microservice environment.

![Micro](https://github.com/micro/micro/blob/master/doc/micro.png)
-

## Features

Feature		|	Description
------		|	-------
[Discovery](https://godoc.org/github.com/micro/go-micro/registry) | Find running services
[Client](https://godoc.org/github.com/micro/go-micro/client) | Query services via RPC
[Server](https://godoc.org/github.com/micro/go-micro/server) | Listen and serve RPC requests
[Pub/Sub](https://godoc.org/github.com/micro/go-micro/broker) | Publish and subscribe to events
[API Gateway](https://github.com/micro/micro/tree/master/api) | Lightweight gateway/proxy. Convert http requests to rpc
[CLI](https://github.com/micro/micro/tree/master/cli) | Command line interface
[Sidecar](https://github.com/micro/micro/tree/master/car) | Integrate any application into the Micro ecosystem
[Web UI/Proxy](https://github.com/micro/micro/tree/master/web) | A visual way to view and query services

## Example Services
Project		|	Description
-----		|	------
[greeter](https://github.com/micro/micro/tree/master/examples/greeter)	|	A greeter service (includes Go, Ruby, Python examples)
[geo-srv](https://github.com/micro/geo-srv)	|	Geolocation tracking service using hailocab/go-geoindex
[geo-api](https://github.com/micro/geo-api)	|	A HTTP API handler for geo location tracking and search
[discovery-srv](https://github.com/micro/discovery-srv)	|	A discovery in the micro platform
[geocode-srv](https://github.com/micro/geocode-srv)	|	A geocoding service using the Google Geocoding API
[hailo-srv](https://github.com/micro/hailo-srv)	|	A service for the hailo taxi service developer api
[monitoring-srv](https://github.com/micro/monitoring-srv)	|	A monitoring service for Micro services
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

## Projects

### Micro
[Micro](https://github.com/micro/micro) itself is the overarching toolkit and ecosystem

### Go Micro
[Go-micro](https://github.com/micro/go-micro) is a pluggable Go framework for writing RPC based microservices. Go micro can be used standalone but fits into the bigger Micro ecosystem.

### Go Platform
[Go-platform](https://github.com/micro/go-platform) provides higher level libraries and services that can be integrated into a go-micro service. Things like tracing, monitoring, dynamic configuration, etc. Again, pluggable like go-micro.

### Go Plugins
[Go-plugins](https://github.com/micro/go-plugins) provides a place for the community to provide their implementations of the interfaces. 
By default Micro will only support 1 or 2 implementations of each interface. Registries built on 
top of kubernetes, zookeeper, etc. Transport using http2, broker using kafka, etc.

### micro-services.co
[Micro-services.co](https://micro-services.co) is a place to share **micro** services. 

### Built in Web UI

<img src="https://github.com/micro/micro/blob/master/web/web1.png">
-
<img src="https://github.com/micro/micro/blob/master/web/web2.png">
-
<img src="https://github.com/micro/micro/blob/master/web/web3.png">

## Getting Started

### Writing a service

Learn how to write and run a microservice using [**go-micro**](https://github.com/micro/go-micro)

### Install Micro

```shell
$ go get github.com/micro/micro
```

Or via Docker

```shell
$ docker pull microhq/micro
```

### Quick start

Run consul (default discovery mechanism)
```
$ go get github.com/hashicorp/consul
$ consul agent -server -bootstrap-expect 1 -data-dir /tmp/consul
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
   micro - A microservices toolchain

USAGE:
   micro [global options] command [command options] [arguments...]
   
VERSION:
   latest
   
COMMANDS:
   api		Run the micro API
   registry	Query registry
   query	Query a service method using rpc
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
   --logtostderr								log to standard error instead of files
   --alsologtostderr								log to standard error as well as files
   --log_dir 									log files will be written to this directory instead of the default temporary directory
   --stderrthreshold 								logs at or above this threshold go to stderr
   -v 										log level for V logs
   --vmodule 									comma-separated list of pattern=N settings for file-filtered logging
   --log_backtrace_at 								when logging hits line file:N, emit a stack trace
   --help, -h									show help
   --version									print the version
   
```

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

## Resources

- [A Journey into Microservices](https://sudo.hailoapp.com/services/2015/03/09/journey-into-a-microservice-world-part-1/)
- [A Journey into a Microservice World](https://speakerdeck.com/mattheath/a-journey-into-a-microservice-world) by Matt Heath (Slides)
- [Microservices](http://martinfowler.com/articles/microservices.html) by Martin Fowler
- [Microservices: Decomposing Applications for Deployability and Scalability](http://www.slideshare.net/chris.e.richardson/microservices-decomposing-applications-for-deployability-and-scalability-jax) by Chris Richardson (Slides)
- [4 reasons why microservices resonate](http://radar.oreilly.com/2015/04/4-reasons-why-microservices-resonate.html) by Neal Ford

## Roadmap

[![Roadmap](https://img.shields.io/badge/roadmap-in%20progress-lightgrey.svg)](https://github.com/micro/micro/wiki/Roadmap)
