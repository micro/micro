# FAQ

## What is Micro?

Micro is a set of tools and libraries to help simplify microservice development and management. It currently consists of **3** components:

- **micro** - The overarching toolkit containing a CLI, Web UI, API and Sidecar (http interface to the core).
- **go-micro** - A pluggable Go library which provides the fundamentals for writing a microservice; service discovery, client/server communication, pub/sub, etc.
- **go-platform** - A feature rich higher level pluggable Go library that sits on top of go-micro to provide a wider range of requirements for a microservice environment; tracing, monitoring, metrics, authentication, key-value, routing, etc.

There are also other libraries also like [go-plugins](https://github.com/micro/go-plugins) for implementations of each package in go-micro or go-platform and [protobuf](https://github.com/micro/protobuf), a fork of golang/protobuf, which provides experimental code generation for go-micro applications.

## How do I use Micro?

You can start by writing a microservice with [**go-micro**](https://github.com/micro/go-micro) or playing with the example [**greeter**](https://github.com/micro/micro/tree/master/examples/greeter) app. The greeter also demonstrates how to integrate non-Go applications. Micro uses proto-rpc and json-rpc by default, libraries are available for both protocols in most languages.

You can find a guide to getting started writing apps [**here**](https://github.com/micro/micro/blob/master/doc/getting-started.md) and a shorter version in the go-micro [readme](https://github.com/micro/go-micro).

Once you have an app running you can use the [**CLI**](https://github.com/micro/micro/tree/master/cli) to query it and also the [**Web UI**](https://github.com/micro/micro/tree/master/web).

There's also docker images on [Docker Hub](https://hub.docker.com/r/microhq/).

## Can I use something besides Consul?

Yes! The registry for service discovery is completely pluggable as is every other package. Consul was used as the default due to its features and simplicity.

As an example. If you would like to use etcd, import the plugin and set the command line flags on your binary.

```go
import (
	_ "github.com/micro/go-plugins/registry/etcd"
)
```

```shell
my_service --registry=etcd --registry_address=127.0.0.1:2379
```


