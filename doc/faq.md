## What is Micro?

Micro is a toolkit to help simplify microservices development. 

It currently consists of a couple components:

- **go-micro** - A pluggable Go library for writing a microservice; service discovery, client/server rpc, pub/sub, etc.
- **micro** - The overarching toolkit containing a CLI, Web UI, API and Sidecar (http interface for go-micro).

There are also other libraries also like [go-plugins](https://github.com/micro/go-plugins) for implementations of each package in go-micro or go-platform and [protobuf](https://github.com/micro/protobuf), a fork of golang/protobuf, which provides experimental code generation for go-micro

## Who's using Micro?

There's a [Users](users.md) page with a list of companies using Micro. Many more are also using it but not yet publicly listed. Feel free to add your company if you're using Micro.

## Where is the community?

There's a slack community with hundreds of members. You can invite yourself [here](http://slack.micro.mu/).

## How do I use Micro?

You can start by writing a microservice with [**go-micro**](https://github.com/micro/go-micro) or playing with the example [**greeter**](https://github.com/micro/micro/tree/master/examples/greeter) app. The greeter also demonstrates how to integrate non-Go applications. Micro uses proto-rpc and json-rpc by default, libraries are available for both protocols in most languages.

You can find a guide to getting started writing apps [**here**](getting-started.md) and a shorter version in the go-micro [readme](https://github.com/micro/go-micro).

Once you have an app running you can use the [**CLI**](https://github.com/micro/micro/tree/master/cli) to query it and also the [**Web UI**](https://github.com/micro/micro/tree/master/web).

There's also docker images on [Docker Hub](https://hub.docker.com/r/microhq/).

## Can I use something besides Consul?

Yes! The registry for service discovery is completely pluggable as is every other package. Consul was used as the default due to its features and simplicity.

### Using etcd

As an example. If you would like to use etcd, import the plugin and set the command line flags on your binary.

```go
import (
        _ "github.com/micro/go-plugins/registry/etcd"
)
```

```shell
my_service --registry=etcd --registry_address=127.0.0.1:2379
```

### Zero Dependency MDNS

Alternatively we can use multicast DNS with the built in MDNS registry for a zero dependency configuration. Just pass `--registry=mdns` to your application on startup.

## Where can I run Micro?

Micro is runtime agnostic. You can run it anywhere you like. On bare metal, on AWS, Google Cloud. On your favourite container orchestration system like Mesos or Kubernetes.

In fact there's a demo of Micro on Kubernetes. Check out the repo at [github.com/micro/kubernetes](https://github.com/micro/kubernetes) and the live demo at [web.micro.pm](http://web.micro.pm).

## What's the different between API, Web and SRV services?

<img src="arch.png" />

As part of the micro toolkit we attempt to define a set of design patterns for a scalable architecture by separating the concerns of the API, Web dashboards and backend services (SRV).

### API Services

API services are served by the micro api with the default namespace go.micro.api. The micro api conforms to the API gateway pattern. 

Learn more about it [here](https://github.com/micro/micro/tree/master/api)

### Web Services

Web services are served by the micro web with the default namespace go.micro.web. We believe in web apps as first class citizens in the microservice world therefor building web dashboards as microservices. The micro web is a reverse proxy and will forward HTTP requests to the appropriate web apps based on path to service resolution. 

Learn more about it [here](https://github.com/micro/micro/tree/master/web)

### SRV services

SRV services are basically standard RPC services, the usual kind of service you would write. We usually call them RPC or backend services as they should mainly be part of the backend architecture and never be public facing. By default we use the namespace go.micro.srv for these but you should use your domain com.example.srv. 

## How performant is it?

Performance is not a current focus of Micro. While code is written to be optimal and avoid overhead, not much time is spent on benchmarks. Comparisons to net/http or other web frameworks make no sense. Micro provides higher level requirements for microservices that include service discovery, load balancing, message encoding, etc. To compare you would need to add all these features in.

If you're still concerned with performance. The simplest way to extract the most value is simply by running with the following flags:

```
--selector=cache # enables in memory caching of discovered nodes
--client_pool_size=10 # enables the client side connection pool
```

## Does Micro support gRPC?

Yes. There are plugins for a transport, client and server in [micro/go-plugins](https://github.com/micro/go-plugins). If you want a quick start just use [micro/go-grpc](https://github.com/micro/go-grpc).

## Where Can I Learn More?

- Join the slack community - [slack.micro.mu](http://slack.micro.mu)
- Read the blog - [blog.micro.mu](https://blog.micro.mu)
- Reach out if you want to talk - [contact@micro.mu](mailto:contact@micro.mu)
