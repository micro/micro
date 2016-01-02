# Roadmap Overview

The micro ecosystem is rapidly growing but there's still a lot to do.

## [Micro](https://github.com/micro/micro)

1. API
  * [x] Allow requests directly to RPC services via path (/greeter/say/hello => service: greeter method: Say.Hello)  
  * Allow REST requests to RPC based services
  * Make the choice a flag/env var
  * Configurable hostnames
  * Configurable namespace for services
  * Support label based routing
  * Support weighted load balancing
2. Web
  * [x] Proxy requests to "web" micro services
  * List "web" micro services on home screen
  * CLI interface in Web UI
3. Sidecar
  * Raise awareness for non Go native app usage
  * Make it work with multiple transports
4. Examples
  * [x] greeter client/server {ruby, python, go}
  * [x] go-micro/examples
  * [x] code generation example
  * [x] geo service and api
  * [x] slack bot API service
  * [x] wrappers/middleware
  * [x] pub sub

## [Go Micro](https://github.com/micro/go-micro)

* [x] Top level initialisation

1. Middleware/Wrappers
  * [x] [Server](https://github.com/micro/go-micro/blob/master/server/server_wrapper.go)
  * [x] [Client](https://github.com/micro/go-micro/blob/master/client/client_wrapper.go)
  * [x] Example implementations
    * [x] [Client](https://github.com/micro/go-micro/tree/master/examples/client/wrapper)
    * [x] [Server](https://github.com/micro/go-micro/blob/master/examples/server/main.go#L12L28)
  * [ ] Plugins e.g. trace, monitoring, logging
2. Code generation
  * [x] Experimental generator [github.com/micro/protobuf](https://github.com/micro/protobuf)
  * [x] Example usage
    * [x] [Client](https://github.com/micro/go-micro/tree/master/examples/client/codegen)
    * [x] [Server](https://github.com/micro/go-micro/tree/master/examples/server/codegen)
  * [x] Server side generator
  * [ ] Stable interface
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
5. Transport
  * [x] Cleanup send/receive semantics - is it concurrent?
6. Codec
  * [ ] Improve codec interface
7. Bidirectional streaming
  * [x] Client
  * [x] Server
  * [x] Code generation for streaming interface
  * [x] Examples

## [Go Platform](https://github.com/micro/go-platform)

Overview
  * [ ] Define the interfaces for every package
  * [ ] Provide documentation for go-platform's usage
  * [ ] Provide easy initialisation and wrapping for go-micro client/server
  * [x] Implement trace and monitoring first

1. Discovery
  * [x] In memory catching using registry watch
  * [x] Heartbeating the registry
2. Routing
  * [ ] MPLS style label based routing 
  * [ ] Circuit breakers
  * [ ] Rate limiting
  * [ ] Weighted loadbalancing
3. Key-Value
  * [x] Implement interface
  * [x] Memcache implementation
  * [x] Redis contribution
  * [x] In memory implement
4. Trace
  * [x] Implement interface
  * [x] Pub/Sub based tracing
5. Monitor
  * [x] Implement interface
  * [x] Custom healthcheck types

## [Go Plugins](https://github.com/micro/go-plugins)

1. Provide more example implementations.
2. Improve auto loading of plugins

## [micro-services.co](https://micro-services.co)

Currently invite only

1. Let more users in
2. Cleanup UI
3. Build services to share on the site {user, login, geo, etc}
