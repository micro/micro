# micro api

The **micro api** is an API gateway for microservices. Use the API gateway [pattern](http://microservices.io/patterns/apigateway.html) to provide a 
single entry point for your services. The micro api serves HTTP and dynamically routes to the appropriate backend using service discovery.

<p align="center">
  <img src="https://github.com/micro/docs/blob/master/images/api.png" />
</p>

## Overview

The micro api is a HTTP api. Requests to the API are served over HTTP and internally routed via RPC. It builds on 
[go-micro](https://github.com/micro/go-micro), leveraging it for service discovery, load balancing, encoding and RPC based communication.

Because the micro api uses go-micro internally, this also makes it pluggable. See [go-plugins](https://github.com/micro/go-plugins) for 
support for gRPC, kubernetes, etcd, nats, rabbitmq and more.

## Getting started

- [Install](#install)
- [Run](#run)
- [ACME](#acme)
- [Examples](#examples)
- [Handlers](#handlers)
- [Resolver](#resolver)

### Install

```shell
go get -u github.com/micro/micro
```

### Run

```shell
micro api
```

### ACME 

Serve securely by default using ACME via letsencrypt

```
MICRO_ENABLE_ACME=true micro api
```

Optionally specify a host whitelist

```
MICRO_ENABLE_ACME=true \
MICRO_ACME_HOSTS=example.com,api.example.com \
micro api
```

### TLS Cert

The API supports serving securely with TLS certificates

```shell
MICRO_ENABLE_TLS=true \
MICRO_TLS_CERT_FILE=/path/to/cert \
MICRO_TLS_KEY_FILE=/path/to/key \
micro api
```

### Set Namespace

The API makes use of namespaces to logically separate backend and public facing services. The namespace and http path 
are used to resolve service name/method. The default namespace is `go.micro.api`.

```shell
MICRO_NAMESPACE=com.example.api micro api
```

## Examples

Here we have an example of a 3 tier architecture

- `micro api`: (localhost:8080) - serving as the http entry point
- `api service`: (go.micro.api.greeter) - serving a public facing api
- `backend service`: (go.micro.srv.greeter) - internally scoped service

The full example is at [examples/greeter](https://github.com/micro/examples/tree/master/greeter)

### Run Example

Prereq: Ensure you are running service discovery e.g consul agent -dev

Get examples

```
git clone https://github.com/micro/examples
```

Start the service

```shell
go run examples/greeter/srv/main.go
```

Start the API

```shell
go run examples/greeter/api/api.go
```

Start the micro api

```
micro api
```

### Query

Make a HTTP call via the micro api

```shell
curl "http://localhost:8080/greeter/say/hello?name=John"
```

The HTTP path /greeter/say/hello maps to service go.micro.api.greeter method Say.Hello

Bypass the api service and call the backend directly via /rpc

```shell
curl -d 'service=go.micro.srv.greeter' \
     -d 'method=Say.Hello' \
     -d 'request={"name": "John"}' \
     http://localhost:8080/rpc
```

Make the same call entirely as JSON

```shell
curl -H 'Content-Type: application/json' \
     -d '{"service": "go.micro.srv.greeter", "method": "Say.Hello", "request": {"name": "John"}}' \
     http://localhost:8080/rpc
```

## API

The micro api provides the following HTTP api

```
- /[service]/[method]	# HTTP paths are dynamically mapped to services
- /rpc			# Explicitly call a backend service by name and method
```

See below for examples

## Handlers

Handlers are HTTP handlers which manage request routing.

The default handler uses endpoint metadata from the registry to determine service routes. If a route match is not found it will 
fallback to the API handler. You can configure routes on registration using the [go-api](https://github.com/micro/go-api).

The API has the following configurable request handlers.

- [`api`](#api-handler) - Handles any HTTP request. Gives full control over the http request/response via RPC.
- [`event`](#event-handler) -  Handles any HTTP request and publishes to a message bus.
- [`http`](#http-handler) - Handles any HTTP request and forwards as a reverse proxy.
- [`rpc`](#rpc-handler) - Handles json and protobuf POST requests. Forwards as RPC.
- [`web`](#web-handler) - The HTTP handler with web socket support included.

Optionally bypass the handlers with the [`/rpc`](#rpc-endpoint) endpoint

### API Handler

The API handler is the default handler. It serves any HTTP requests and forwards on as an RPC request with a specific format.

- Content-Type: Any
- Body: Any
- Forward Format: [api.Request](https://github.com/micro/go-api/blob/master/proto/api.proto#L11)/[api.Response](https://github.com/micro/go-api/blob/master/proto/api.proto#L21)
- Path: `/[service]/[method]`
- Resolver: Path is used to resolve service and method
- Configure: Flag `--handler=api` or env var `MICRO_API_HANDLER=api`
- The default handler when no handler is specified

### Event Handler

The event handler serves HTTP and forwards the request as a message over a message bus using the go-micro broker.

- Content-Type: Any
- Body: Any
- Forward Format: Request is formatted as [go-api/proto.Event](https://github.com/micro/go-api/blob/master/proto/api.proto#L28L39) 
- Path: `/[topic]/[event]`
- Resolver: Path is used to resolve topic and event name
- Configure: Flag `--handler=event` or env var `MICRO_API_HANDLER=event`

### HTTP Handler

The http handler is a http reserve proxy with built in service discovery.

- Content-Type: Any
- Body: Any
- Forward Format: HTTP Reverse proxy
- Path: `/[service]`
- Resolver: Path is used to resolve service name
- Configure: Flag `--handler=http` or env var `MICRO_API_HANDLER=http`
- REST can be implemented behind the API as microservices

### RPC Handler

The RPC handler serves json or protobuf HTTP POST requests and forwards as an RPC request.

- Content-Type: `application/json` or `application/protobuf`
- Body: JSON or Protobuf
- Forward Format: json-rpc or proto-rpc based on content
- Path: `/[service]/[method]`
- Resolver: Path is used to resolve service and method
- Configure: Flag `--handler=rpc` or env var `MICRO_API_HANDLER=rpc`

### Web Handler

The web handler is a http reserve proxy with built in service discovery and web socket support.

- Content-Type: Any
- Body: Any
- Forward Format: HTTP Reverse proxy including web sockets
- Path: `/[service]`
- Resolver: Path is used to resolve service name
- Configure: Flag `--handler=web` or env var `MICRO_API_HANDLER=web`

### RPC endpoint

The /rpc endpoint let's you bypass the main handler to speak to any service directly

- Request Params
  * `service` - sets the service name
  * `method` - sets the service method
  * `request` - the request body
  * `address` - optionally specify host address to target

Example call:

```
curl -d 'service=go.micro.srv.greeter' \
     -d 'method=Say.Hello' \
     -d 'request={"name": "Bob"}' \
     http://localhost:8080/rpc
```

Find working examples in [github.com/micro/examples/api](https://github.com/micro/examples/tree/master/api)


## Resolver

Micro dynamically routes to services using a namespace value and the HTTP path.

The default namespace is `go.micro.api`. Set namespace via `--namespace` or `MICRO_NAMESPACE=`.

The resolvers used are explained below.

- [rpc](#rpc-resolver)
- [http](#http-resolver)
- [event](#event-resolver)

### RPC Resolver

RPC services have a name (go.micro.api.greeter) and a method (Greeter.Hello).

URLs are resolved as follows:

Path	|	Service	|	Method
----	|	----	|	----
/foo/bar	|	go.micro.api.foo	|	Foo.Bar
/foo/bar/baz	|	go.micro.api.foo	|	Bar.Baz
/foo/bar/baz/cat	|	go.micro.api.foo.bar	|	Baz.Cat

Versioned API URLs can easily be mapped to service names:

Path	|	Service	|	Method
----	|	----	|	----
/foo/bar	|	go.micro.api.foo	|	Foo.Bar
/v1/foo/bar	|	go.micro.api.v1.foo	|	Foo.Bar
/v1/foo/bar/baz	|	go.micro.api.v1.foo	|	Bar.Baz
/v2/foo/bar	|	go.micro.api.v2.foo	|	Foo.Bar
/v2/foo/bar/baz	|	go.micro.api.v2.foo	|	Bar.Baz

### HTTP Resolver

With the http handler we only need to deal with resolving the service name. So the resolution differs slightly to the RPC resolver.

URLS are resolved as follows:

Path	|	Service	|	Service Path
---	|	---	|	---
/foo	|	go.micro.api.foo	|	/foo
/foo/bar	|	go.micro.api.foo	|	/foo/bar
/greeter	|	go.micro.api.greeter	|	/greeter
/greeter/:name	|	go.micro.api.greeter	|	/greeter/:name

### Event Resolver

The micro api can be used as a pubsub gateway. The event handler uses the event resolver to publish events to topics.

URLS are resolved to topics and events as follows

Path	|	Topic	|	Event Name
---	|	---	|	---
/	|	go.micro.api	|	event
/user	|	go.micro.api.user	|	user
/user/login	|	go.micro.api.user	|	login
/user/login/john	|	go.micro.api.user	| 	login.john
