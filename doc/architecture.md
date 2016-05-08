# Architecture

This section should explain more about how micro is constructed, how the various libraries/repos relate to each other 
and how they should be used.

## Micro

Micro is an ecosystem of tools and libraries to simplify developing microservices. By using go-micro you can essentially 
create microservices that are discoverable, can query and be queried via RPC and publish/subscribe events. The wider 
toolset provides everything necessary to run a platform.

![Micro](https://github.com/micro/micro/blob/master/doc/micro.png)
-

### CLI

The Micro CLI is a command line version of go-micro which provides a way of observing and interacting with a running 
environment.

### API

The API acts as a gateway or proxy to enable a single entry point for accessing micro services. It should be run on the edge of your infrastructure. It converts HTTP requests to RPC and forwards to the appropriate service.

<p align="center">
  <img src="https://github.com/micro/micro/blob/master/api/api.png" />
</p>

### Web UI

The UI is a web version of go-micro allowing visual interaction into an environmet. In the future it will be a way of aggregating micro web services also. It includes a way to proxy to web apps. /[name] will route to a 
a service in the registry. The Web UI adds Prefix of "go.micro.web." (which can be configured) to the name, looks 
it up in the registry and then reverse proxies to it.

<p align="center">
  <img src="https://github.com/micro/micro/blob/master/web/web.png" />
</p>

### Sidecar

The sidecar is a HTTP Interface version of go-micro. It's a way of integrating non-Go appications into a micro 
environment. 

<p align="center">
  <img src="https://github.com/micro/micro/blob/master/car/sidecar.png" />
</p>

### Bot

Bot A Hubot style bot that sits inside your microservices platform and can be interacted with via Slack, HipChat, XMPP, etc. It provides the features of the CLI via messaging. Additional commands can be added to automate common ops tasks.

<p align="center">
  <img src="https://github.com/micro/micro/blob/master/bot/bot.png" />
</p>

## Overview

Below is an overview of how services within a micro environment interact. Each time a service needs to make a 
request, it will lookup the service name within the registry then directly send a request to an instance of 
the service.

![Overview1](https://github.com/micro/micro/blob/master/doc/overview1.png)
-

![Overview2](https://github.com/micro/micro/blob/master/doc/overview2.png)
-

![Overview3](https://github.com/micro/micro/blob/master/doc/overview3.png)

## Go Micro

<p align="center">
  <img src="https://github.com/micro/go-micro/blob/master/go-micro.png" />
</p>

### Registry

The registry provides a pluggable service discovery library to find running services. Current implementations 
are consul, etcd, memory and kubernetes. The interface is easily implemented if your preferences differ.

### Selector

The selector provides a load balancing mechanism via selection. When a client makes a request to a service it 
will first query the registry for the service. This usually returns a list of running nodes representing 
the service. A selector will select one of these nodes to be used for querying. Multiple calls to the selector 
will allow balancing algorithms to be utilised. Current methods are round robin, random hashed and blacklist. 

### Broker

The broker is pluggable interface for pub/sub. Microservices are an event driven architecture where a publishing 
and subscribing to events should be a first class citizen. Current implementations include nats, rabbitmq and http 
(for development).

### Transport

Transport is a pluggable interface over point to point transfer of messages. Current implementations are http, 
rabbitmq and nats. By providing this abstraction, transports can be swapped out seamlessly.

### Client

The client provides a way to make RPC queries. It combines the registry, selector, broker and transport. It also 
provides retries, timeouts, use of context, etc.

### Server

The server is an interface to build a running microservice. It provides a way of serving RPC requests. 

## Go Platform

This is a microservice platform library built to be used with micro/go-micro. 
It provides all of the fundamental tooling required to run and manage 
a microservice environment. It's pluggable just like go-micro. It will be 
especially vital for anything above 20+ services and preferred by 
organisations. Developers looking to write standalone services should 
continue to use go-micro. 

### How does it work?

The go-platform is a client side interface for the fundamentals of a microservice platform. Each package connects to 
a service which handles that feature. Everything is an interface and pluggable which means you can choose how to 
architect your platform. Micro however provides a "platform" implementation backed by it's own services by default.

Each package can be used independently or integrated using go-micro client and handler wrappers.

### Auth 

Auth addresses authentication and authorization of services and users. The default implementation is Oauth2 with an additional policy 
engine coming soon. This is the best way to authenticate users and service to service calls using a centralised 
authority. Security is a first class citizen in a microservice platform.

### Config 

Config implements an interface for dynamic configuration. The config can be hierarchically loaded and merged from 
multiple sources e.g file, url, config service. It can and should also be namespaced so that environment specific 
config is loaded when running in dev, staging or production. The config interface is useful for business level 
configuration required by your services. It can be reloaded without needing to restart a service.

### DB (experimental) 

The DB interface is an experiment CRUD interface to simplify database access and management. The amount of CRUD boilerplate 
written and rewritten in a microservice world is immense. By offloading this to a backend service and using RPC, we 
eliminate much of that and speed up development. The platform implementation includes pluggable backends such as mysql, 
cassandra, elasticsearch and utilises the registry to lookup which nodes databases are assigned to. 

This is purely experimental at this point based on some ideas from how Google, Facebook and Twitte do database management 
internally.
 
### Discovery 

Discovery provides a high level service discovery interface on top of the go-micro registry. It utilises the watcher to 
locally cache service records and also heartbeats to a discovery service. It's akin to the Netflix Eureka 2.0 
architecture whereby we split the read and write layers of discovery into separate services.

### Event

The event package provides a way to send platform events and essentially create an event stream and record of all that's 
happening in your microservice environment. On the backend an event service aggregates the records and allows you to 
subscribe to a specific set of events. An event driven architecture is a powerful concept in a microservice environment 
and must be addressed adequately. At scale it's essential for correlating events within a distributed system e.g 
provisioning of new services, change of dynamic config, logouts for customers, tracking notifications, alerts.
 
### KV 

KV represents a simple distributed key-value interface. It's useful for sharing small fast access bits of data amonst 
instances of a service. We provide three implementations currently. Memcached, redis and a consistently hashed in distributed 
in memory system.

### Log 

Log provides a structured logging interface which allows log messages to be tagged with key-value pairs. 
The default output plugin is file which allows many centralised logging systems to be used such as the ELK stack. 

### Monitor 

The monitor provides a way to publish Status, Stats and Healtchecks to a monitoring service. Healthchecks are user defined 
checks that may be critical to a service e.g can access database, can sync from s3, etc. Monitoring in a distributed 
system is fundamentally different from the classic LAMP stack. In the old ways pings and tcp checks were regarded as enough, 
in a distributed system we require much more fine grained metrics and a monitoring service which can make sense of what 
failure means in this world.

### Metrics 

Metrics is an interface for instrumentation. We regard metrics as a superior form of observability in a distributed system over 
logging. Instrumentation is a great way to graph historic and realtime data which can be correlated and immediately 
understood. The metrics interface provides a way to create counters, gauges and histograms. We currently implement the statsd 
interface and offload to telegraf which provides an augmented statsd interface with labels.

### Router

The router builds on the registry and selector to provide rate limiting, circuit breaking and global service load balancing. 
It implements the selector interface. Stats are recorded for every request and periodically published. A centralised routing 
service aggregates these metrics from all services in the environment and makes decisions about how to route requests. 
The routing service is not a proxy. Proxies are a weak form of load balancing, we prefer smart clients which retrieve 
a list of nodes from the router and make direct connections, this means if the routing service dies or misbehaves, clients 
can continue to make request independently.
 
### Sync 

Sync is an interface for distributed synchronisation. This provides an easy way to do leadership election and locking to 
serialise access to a resource. We expect there to be multiple copies of a service running to provide fault tolerance and 
scalability but it makes it much harder to deal with transactions or serialising access. The sync package provides a 
way to regain some of these semantics.
 
### Trace 

Trace is a client side interface for distributed tracing e.g dapper, zipkin, appdash. In a microservice world, a single 
request may fan out to 20-30 services. Failure may be non deterministic and difficult to track. Distributed tracing is a 
way of tracking the lifetime of a request. The interface utilises client and server wrappers to simplify using tracing.

## Go Plugins

Go plugins is a place to share implementations of the go-micro and go-platform interfaces. 

## Code Generation

Micro provides experimental code generation to reduce the amount of boiler plate code needed to be written. An example 
of how to use code generation can be found [here](https://github.com/micro/go-micro/tree/master/examples/client/codegen).
