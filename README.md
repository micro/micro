# Micro

Micro is **a microservices toolchain** consisting of a suite of libraries and tools to write and run microservices.

Checkout [**go-micro**](https://github.com/myodc/go-micro) if you want to start writing services now.

# Overview
The goal of **Micro** is to provide a toolchain for microservice development and management. At the core, micro is simple and accessible enough that anyone can easily get started writing microservices. As you scale to hundreds of services, micro will provide the fundamental tools required to manage a microservice environment.


## Features
- Discovery
- Client/Server
- Pub/Sub
- API Gateway
- CLI
- Sidecar - for non Go native apps

## Future Features
- Config
- Routing
- Monitoring
- Tracing
- Logging

## Libraries & Tools
- [go-micro](https://github.com/myodc/go-micro) - A microservices client/server library based on http/rpc protobuf
- [api](https://github.com/myodc/micro/tree/master/api) - A lightweight gateway/proxy for Micro based services
- [cli](https://github.com/myodc/micro/tree/master/cli) - A command line tool for micro
- [sidecar](https://github.com/myodc/micro/tree/master/car) - Integrate any application into the Micro ecosystem

## Example Services
- [geo-srv](https://github.com/myodc/geo-srv) - A go-micro based geolocation tracking service using hailocab/go-geoindex
- [geo-api](https://github.com/myodc/geo-api) - A HTTP API handler for geo location tracking and search
- [greeter](https://github.com/myodc/micro/tree/master/examples/greeter) - A greeter Go service

## Getting Started

### Install

```shell
$ go get github.com/myodc/micro
```

### Usage
```shell
NAME:
   micro - A microservices toolchain

USAGE:
   micro [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
   api		Run the micro API
   registry	Query registry
   query	Query a service method using rpc
   health	Query the health of a service
   list		List items in registry
   get		Get item from registry
   sidecar	Run the micro sidecar
   help, h	Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --server_address ":0"	Bind address for the server. 127.0.0.1:8080 [$MICRO_SERVER_ADDRESS]
   --broker "http"		Broker for pub/sub. http, nats, etc [$MICRO_BROKER]
   --broker_address 		Comma-separated list of broker addresses [$MICRO_BROKER_ADDRESS]
   --registry "consul"		Registry for discovery. kubernetes, consul, etc [$MICRO_REGISTRY]
   --registry_address 		Comma-separated list of registry addresses [$MICRO_REGISTRY_ADDRESS]
   --transport "http"		Transport mechanism used; http, rabbitmq, etc [$MICRO_TRANSPORT]
   --transport_address 		Comma-separated list of transport addresses [$MICRO_TRANSPORT_ADDRESS]
   --help, -h			show help
   --version, -v		print the version
```

Read more on how to use Micro [here](https://github.com/myodc/micro/tree/master/cli)

Learn how to write and run a microservice using Go-Micro [here](https://github.com/myodc/go-micro)

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
