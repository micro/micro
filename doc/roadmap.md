# Roadmap Overview

The micro ecosystem is rapidly growing but there's still a lot to do.

## [Micro](https://github.com/micro/micro)

1. [API](https://github.com/micro/micro/tree/master/api)
  * [x] Allow requests directly to RPC services via path (/greeter/say/hello => service: greeter method: Say.Hello)
  * [x] TLS Support
  * [x] Allow namespace to be set via flags
  * [x] Apache log format
  * [x] Stats page
  * [x] Allow REST requests to RPC based services
  * [x] Make the choice a flag/env var
  * [ ] Configurable hostnames
  * [x] Configurable namespace for services
  * [ ] Support label based routing
  * [ ] Support weighted load balancing
  * Google GFE like semantics
2. [Web](https://github.com/micro/micro/tree/master/web)
  * [x] Proxy requests to "web" micro services
  * [x] List "web" micro services on home screen
  * [x] TLS Support
  * [x] Web Sockets
  * [x] Allow namespace to be set via flags
  * [x] Apache log format
  * [x] Stats page
  * [x] CLI interface in Web UI
3. [Sidecar](https://github.com/micro/micro/tree/master/car)
  * [x] TLS Support
  * [x] Apache log format
  * [x] Stats page
  * [ ] Raise awareness for non Go native app usage
  * [ ] Make it work with multiple transports
4. [CLI](https://github.com/micro/micro/tree/master/cli)
  * [x] Support querying via proxying the sidecar
  * [x] Allow connecting through the API or Web where private network isn't available
    - Done via the [Sidecar](https://github.com/micro/micro/tree/master/car#proxy-cli-requests)
5. Bot
  * [x] Implement the bot
  * [x] Feature parity with CLI
  * [x] Slack input
  * [x] Hipchat input
5. Dependencies
  * [ ] Create dependency management config for services
  * [ ] Allow push/pull from micro-services.co
6. Examples
  * [x] [greeter](https://github.com/micro/micro/tree/master/examples/greeter) client/server {ruby, python, go}
  * [x] [go-micro/examples](https://github.com/micro/go-micro/tree/master/examples)
  * [x] [pub sub](https://github.com/micro/go-micro/tree/master/examples/pubsub)
  * [x] code generation example
  * [x] geo service and api
  * [x] slack bot API service
  * [x] wrappers/middleware

## [Go Micro](https://github.com/micro/go-micro)

* [x] Top level initialisation

1. Middleware/Wrappers
  * [x] [Server](https://github.com/micro/go-micro/blob/master/server/server_wrapper.go)
  * [x] [Client](https://github.com/micro/go-micro/blob/master/client/client_wrapper.go)
  * [x] Example implementations
    * [x] [Client](https://github.com/micro/go-micro/tree/master/examples/client/wrapper)
    * [x] [Server](https://github.com/micro/go-micro/blob/master/examples/server/main.go#L12L28)
  * [x] Plugins e.g. trace, monitoring, logging
2. Code generation
  * [x] Experimental generator [github.com/micro/protobuf](https://github.com/micro/protobuf)
  * [x] Example usage
    * [x] [Client](https://github.com/micro/go-micro/tree/master/examples/client/codegen)
    * [x] [Server](https://github.com/micro/go-micro/tree/master/examples/server/codegen)
  * [x] Server side generator
  * [x] Stable interface
3. Registry
  * [x] Support Service TTLs on registration so services can be automatically removed
  * [x] Healthchecking function to renew registry lease
  * [x] Service/Node filters - known as a [Selector](https://github.com/micro/go-micro/blob/master/selector)
  * [x] Fix the watch code to return a channel with updates rather than store in memory
  * [x] Add timeout option for querying
4. Broker
  * [x] Support distributed queuing
  * [x] Support acking of messages
  * [x] Support concurrency with options
  * [x] [MQTT](https://godoc.org/github.com/micro/go-micro/broker/mqtt)
5. Transport
  * [x] Cleanup send/receive semantics - is it concurrent?
6. Bidirectional streaming
  * [x] Client
  * [x] Server
  * [x] Code generation for streaming interface
  * [x] Examples
7. TLS Support
  * [x] Registry
  * [x] Broker
  * [x] Transport
8. Selector
  * [x] [Random](https://github.com/micro/go-micro/tree/master/selector/random)
  * [x] [Round Robin](https://github.com/micro/go-micro/tree/master/selector/roundrobin)
  * [x] [Blacklist](https://github.com/micro/go-micro/tree/master/selector/blacklist)
9. Select Filters
  * [x] [Label](https://godoc.org/github.com/micro/go-micro/selector#FilterLabel)
  * [x] [Version](https://godoc.org/github.com/micro/go-micro/selector#FilterVersion)
  * [x] [Endpoint](https://godoc.org/github.com/micro/go-micro/selector#FilterEndpoint)

## [Go Platform](https://github.com/micro/go-platform)

Overview
  * [x] Define the interfaces for every package
  * [x] Provide documentation for go-platform's usage
  * [x] Implement trace and monitoring first
  * [ ] Provide easy initialisation and wrapping for go-micro client/server

1. Discovery
  * [x] In memory catching using registry watch
  * [x] Heartbeating the registry
2. Routing
  * [x] label based routing 
  * [x] Weighted loadbalancing
  * [ ] Circuit breakers
  * [ ] Rate limiting
  * Google GSLB style semantics
3. Key-Value
  * [x] Implement interface
  * [x] Memcache implementation
  * [x] Redis contribution
  * [x] In memory implement
4. Trace
  * [x] Implement interface
  * [x] Pub/Sub based tracing
  * [ ] Timing endpoints for trace service
5. Monitor
  * [x] Implement interface
  * [x] Custom healthcheck types
  * [x] Add stats/status publications
  * [ ] Monitor the health of services
6. Config
  * [x] Implement interface
7. Auth
  * [x] Implement interface
8. Logging
  * [x] Implement interface
9. Event
  * [x] Implement interface (sparse)
10. Metrics
  * [x] Implement interface

## [Go Plugins](https://github.com/micro/go-plugins)

1. [x] Provide more example implementations.
2. [ ] Improve auto loading of plugins

1. Registry
  * [x] [consul](https://godoc.org/github.com/micro/go-micro/registry/consul)
  * [x] [mdns](https://godoc.org/github.com/micro/go-micro/registry/mdns)
  * [x] [etcd](https://godoc.org/github.com/micro/go-plugins/registry/etcd)
  * [x] [nats](https://godoc.org/github.com/micro/go-plugins/registry/nats)
  * [x] [eureka](https://godoc.org/github.com/micro/go-plugins/registry/eureka)
  * [x] [gossip](https://godoc.org/github.com/micro/go-plugins/registry/gossip)
2. Transport
  * [x] [nats](https://godoc.org/github.com/micro/go-plugins/transport/nats)
  * [x] [rabbitmq](https://godoc.org/github.com/micro/go-plugins/transport/rabbitmq)
3. Broker
  * [x] [nats](https://godoc.org/github.com/micro/go-plugins/broker/nats)
  * [x] [nsq](https://godoc.org/github.com/micro/go-plugins/broker/nsq)
  * [x] [rabbitmq](https://godoc.org/github.com/micro/go-plugins/broker/rabbitmq)
  * [x] [kafka](https://godoc.org/github.com/micro/go-plugins/broker/kafka)
  * [x] [googlepubsub](https://godoc.org/github.com/micro/go-plugins/broker/googlepubsub)

## [Platform](https://github.com/micro/platform)

### Dashboards

Create simple OSS dashboards for each platform service

 * [x] [Config](https://github.com/micro/config-web)
 * [x] [Discovery](https://github.com/micro/discovery-web)
 * [x] [Events](https://github.com/micro/event-web)
 * [x] [Monitoring](https://github.com/micro/monitor-web)
 * [x] [Tracing](https://github.com/micro/trace-web)
 * [x] [Routing](https://github.com/micro/router-web)
 * [ ] Logging
 * [ ] Auth
 * [ ] API
 * [ ] Metrics

### Services

Implement version 1s
* [x] [Discovery](https://github.com/micro/discovery-srv)
* [x] [Monitoring](https://github.com/micro/monitor-srv)
* [x] [Trace](https://github.com/micro/trace-srv)
* [x] [Config](https://github.com/micro/config-srv)
* [x] [Auth](https://github.com/micro/auth-srv)
* [x] [Event](https://github.com/micro/event-srv)
* [x] [Router](https://github.com/micro/router-srv)
* [ ] Metrics
* [ ] Logging

## Deployment
* [x] [Micro on Kubernetes](https://github.com/micro/kubernetes)
* [x] Micro Docker Compose
* [ ] Platform Docker Compose

## Demo
* [x] [Web](http://web.micro.pm)
* [x] [API](http://api.micro.pm)
* [x] [Sidecar](http://proxy.micro.pm) 
