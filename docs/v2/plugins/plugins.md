---
title: Plugins
keywords: plugins
tags: [plugins]
sidebar: home_sidebar
permalink: /plugins
summary: 
---

Micro provides a pluggable architecture for all tooling. This means underlying implementations can be swapped out.

Go Micro and the Micro toolkit include separate types of plugins. Navigate from the sidebar to learn more about each.

## Overview

At a high level plugins provide the opportunity to swap out underlying infrastructure and dependencies. This means 
a microservice can be written in one way and run locally with zero dependencies but then equally run as a highly 
resilient system in the cloud when using distributed systems to underpin its usage.

## Usage

By default go-micro only provides a few implementation of each interface at the core but it's completely pluggable. 
There's already dozens of plugins which are available at [github.com/micro/go-plugins](https://github.com/micro/go-plugins). 
Contributions are welcome! Plugins ensure that your Go Micro services live on long after technology evolves.

If you want to integrate plugins simply link them in a separate file and rebuild.

Create a plugins.go file and import the plugins you want:

```go
package main

import (
        // consul registry
        _ "github.com/micro/go-plugins/registry/consul"
        // rabbitmq transport
        _ "github.com/micro/go-plugins/transport/rabbitmq"
        // kafka broker
        _ "github.com/micro/go-plugins/broker/kafka"
)
```

## Write Plugins

Plugins are a concept built on Go's interface. Each package maintains a high level interface abstraction. 
Simply implement the interface and pass it in as an option to the service.

The service discovery interface is called [Registry](https://pkg.go.dev/github.com/micro/go-micro/v2/registry#Registry). 
Anything which implements this interface can be used as a registry. The same applies to the other packages.

```go
type Registry interface {
    Register(*Service, ...RegisterOption) error
    Deregister(*Service) error
    GetService(string) ([]*Service, error)
    ListServices() ([]*Service, error)
    Watch() (Watcher, error)
    String() string
}
```

Browse [go-plugins](https://github.com/micro/go-plugins) to get a better idea of implementation details.

## Examples

### Framework

- [Consul Registry](https://github.com/micro/go-plugins/tree/master/registry/consul) - Service discovery using Consul
- [K8s Registry](https://github.com/micro/go-plugins/tree/master/registry/kubernetes) - Service discovery using Kubernetes
- [Kafka Broker](https://github.com/micro/go-plugins/tree/master/broker/kafka) - Kafka message broker

### Runtime

- [Router](https://github.com/micro/go-plugins/tree/master/micro/router) - Configurable http routing and proxying
- [AWS X-Ray](https://github.com/micro/go-plugins/tree/master/micro/trace/awsxray) - Tracing integration for AWS X-Ray
- [IP Whitelist](https://github.com/micro/go-plugins/tree/master/micro/ip_whitelist) - Whitelisting IP address access

## Repository

The open source plugins can be found at [github.com/micro/go-plugins](https://github.com/micro/go-plugins).


