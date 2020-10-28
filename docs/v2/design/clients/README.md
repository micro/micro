# Clients

Micro has mainly focused on the Go programming language via the [Go Micro](https://github.com/micro/go-micro) framework. We're 
now going to move into multi-language via clients built against the Micro itself. This will involve code generating 
via a gRPC api and then building the clients on top. 

Languages we want to support from day 1 include {go, py, js, rb, java}

## API

The API for each service will live in [github.com/micro/clients](https://github.com/micro/clients). We'll have a proto dir 
which provides a proto per service that micro provides.

```
clients/
	proto/
		registry/registry.proto
		broker/broker.proto
		store/store.proto
		client/client.proto
		...
		
```

Alternatively can potentially introduce a consolidated API for ease of use like so:

```
syntax = "proto3";

package server;

service Client {
	rpc Call(Request) returns (Response) {};
	rpc Stream(stream Request) returns (stream Response) {};
	rpc Publish(Message) returns (Empty) {};
	rpc Subscribe(Topic) returns (stream Message) {};
	rpc Register(Service) returns (Empty) {};
	rpc Deregister(Service) returns (Empty) {};
}

```

We'll then code generate via this api and have gRPC clients that can be used in any language. Although the goal then 
is to level up and build idiomatic clients around this in every language to provide a truly beautiful developer experience. 
It's clear that gRPC has its benefits but its clients are not great beyond Go. From a microservices perspective 
enabling that via higher level clients makes the most sense.

## Serving

Looking at the developer experience for adopting an api, cache, search, database or anything else it's clear the experience 
needs to be a drop-in server and then providing client libraries in any language. This is sort of a frictionless thing 
which augments the app development experience without having to totally buy into a framework.

Micro can now be booted using a single command which provides a vastly superior developer experience.

```
micro server
```

This launches all the services which exist as interfaces in go-micro but will also have equivalent protos to call the services. 
All requests can be routed through the fixed entry point of `localhost:8081` which is the micro proxy. This will forward 
the request to the appropriate place.

Services that are not micro native can use a single command to be registered in discovery

```
# micro service --name [service name] --endpoint [address of service] [command to exec]
micro service --name helloworld --endpoint localhost:9090 go run main.go
```

Endpoint can include the protocol of the service e.g `http://localhost:9090` or `grpc://localhost:9090` otherwise we default to mucp/http

## Routing

The protos will define a package name which is prefixed to the service method in grpc 

```
package go.micro.registry;

service Registry {
	rpc GetService(...) returns (...) {};
}
```

The above would translate to `/go.micro.registry.Registry/GetService`. The proxy will use `go.micro.registry` as the service it calls.
