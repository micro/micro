# Micro - a microservices ecosystem
The goal of **Micro** is to provide an ecosystem of tools and services for microservice development and management. At the core, the toolkit is simple and accessible enough that anyone can easily get started writing microservices. As you scale to hundreds of services, micro will provide the fundamental tools required to manage a microservice environment.

Checkout the [roadmap](https://github.com/micro/micro/wiki/Roadmap) to see where it's all going.

# Overview

Below are the libraries and services that encompass the micro ecosystem.

### Micro
[Micro](https://github.com/micro/micro) itself is the overarching toolkit and ecosystem

Features:
- API
- CLI
- Web UI
- Sidecar - provides all features over go-micro through a http interface

### Go Micro
[Go-micro](https://github.com/micro/go-micro) is a pluggable Go framework for writing RPC based microservices. Go micro can be used standalone but fits into the bigger Micro ecosystem.

Features:
- Service discovery
- Client/Server
- Pub/Sub
- Codecs

### Go Platform
[Go-platform](https://github.com/micro/go-platform) provides higher level libraries and services that can be integrated into a go-micro service. Things like tracing, monitoring, dynamic configuration, etc. Again, pluggable like go-micro.

Features:
- Discovery
- Tracing
- Monitoring
- Logging
- Dynamic Config
- Key Value
- Database
- Instrumentation
- Authentication

### Go Plugins
[Go-plugins](https://github.com/micro/go-plugins) provides a place for the community to provide their implementations of the interfaces. 
By default Micro will only support 1 or 2 implementations of each interface. Registries built on 
top of kubernetes, zookeeper, etc. Transport using http2, broker using kafka, etc.

### micro-services.co
[Micro-services.co](https://micro-services.co) is a place to share **micro** services. 
