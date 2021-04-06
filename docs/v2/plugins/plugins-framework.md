---
title: Framework Plugins
keywords: plugins
tags: [plugins]
sidebar: home_sidebar
permalink: /plugins-framework
summary: 
---

Micro is a pluggable toolkit and framework. The internal features can be swapped out with any [go-plugins](https://github.com/micro/go-plugins).

The toolkit has a separate plugin interface. Learn more at [micro/plugin](https://github.com/micro/micro/tree/master/plugin).

Below is info on go-micro plugin usage.

## Usage

Plugins can be added to go-micro in the following ways. By doing so they'll be available to set via command line args or environment variables.

Import the plugins in a Go program then call service.Init to parse the command line and environment variables.

```go
import (
	"github.com/micro/go-micro/v2"
	_ "github.com/micro/go-plugins/broker/rabbitmq"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	_ "github.com/micro/go-plugins/transport/nats"
)

func main() {
	service := micro.NewService(
		// Set service name
		micro.Name("my.service"),
	)

	// Parse CLI flags
	service.Init()
}
```

### Flags

Specify the plugins as flags

```shell
go run service.go --broker=rabbitmq --registry=kubernetes --transport=nats
```

### Env

Use env vars to specify the plugins

```
MICRO_BROKER=rabbitmq \
MICRO_REGISTRY=kubernetes \ 
MICRO_TRANSPORT=nats \ 
go run service.go
```

### Options

Import and set as options when creating a new service

```go
import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-plugins/registry/kubernetes"
)

func main() {
	registry := kubernetes.NewRegistry() //a default to using env vars for master API

	service := micro.NewService(
		// Set service name
		micro.Name("my.service"),
		// Set service registry
		micro.Registry(registry),
	)
}
```

## Build

An anti-pattern is modifying the `main.go` file to include plugins. Best practice recommendation is to include
plugins in a separate file and rebuild with it included. This allows for automation of building plugins and
clean separation of concerns.

Create file plugins.go

```go
package main

import (
	_ "github.com/micro/go-plugins/broker/rabbitmq"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	_ "github.com/micro/go-plugins/transport/nats"
)
```

Build with plugins.go

```shell
go build -o service main.go plugins.go
```

Run with plugins

```shell
MICRO_BROKER=rabbitmq \
MICRO_REGISTRY=kubernetes \
MICRO_TRANSPORT=nats \
service
```

## Rebuild Toolkit

If you want to integrate plugins simply link them in a separate file and rebuild

Create a plugins.go file

```go
import (
        // etcd v3 registry
        _ "github.com/micro/go-plugins/registry/etcdv3"
        // nats transport
        _ "github.com/micro/go-plugins/transport/nats"
        // kafka broker
        _ "github.com/micro/go-plugins/broker/kafka"
)
```

Build binary

```shell
// For local use
go build -i -o micro ./main.go ./plugins.go

// For docker image
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -i -o micro ./main.go ./plugins.go
```

Flag usage of plugins

```shell
micro --registry=etcdv3 --transport=nats --broker=kafka
```

## Repository

The go-micro plugins can be found in [github.com/micro/go-plugins](https://github.com/micro/go-plugins).


