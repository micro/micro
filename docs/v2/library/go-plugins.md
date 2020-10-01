---
title: Plugins
keywords: go-micro, framework, plugins
tags: [go-micro, framework, plugins]
sidebar: home_sidebar
permalink: /go-plugins
summary: Go Micro is a pluggable framework
---

# Overview

Go Micro is built on Go interfaces. Because of this the implementation of these interfaces are pluggable.

By default go-micro only provides a few implementation of each interface at the core but it's completely pluggable. 
There's already dozens of plugins which are available at [github.com/micro/go-plugins](https://github.com/micro/go-plugins). 
Contributions are welcome! Plugins ensure that your Go Micro services live on long after technology evolves.

### Add plugins

If you want to integrate plugins simply link them in a separate file and rebuild

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

Build your application by including the plugins file:

```shell
# assuming files main.go and plugins.go are in the top level dir
 
# For local use
go build -o service *.go
```

Flag usage of plugins:

```shell
service --registry=etcdv3 --transport=nats --broker=kafka
```

Or what's preferred is using environment variables for 12-factor apps

```
MICRO_REGISTRY=consul \
MICRO_TRANSPORT=rabbitmq \
MICRO_BROKER=kafka \
service
```

### Plugin Option

Alternatively you can set the plugin as an option to a service directly in code

```go
package main

import (
        "github.com/micro/go-micro/v2" 
        // consul registry
        "github.com/micro/go-plugins/registry/consul"
        // rabbitmq transport
        "github.com/micro/go-plugins/transport/rabbitmq"
        // kafka broker
        "github.com/micro/go-plugins/broker/kafka"
)

func main() {
	registry := consul.NewRegistry()
	broker := kafka.NewBroker()
	transport := rabbitmq.NewTransport()

        service := micro.NewService(
                micro.Name("greeter"),
                micro.Registry(registry),
                micro.Broker(broker),
                micro.Transport(transport),
        )

	service.Init()
	service.Run()
}
```

### Write Plugins

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

