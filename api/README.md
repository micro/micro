# micro api

The **micro api** is an API gateway for microservices. Use the API gateway [pattern](http://microservices.io/patterns/apigateway.html) to provide a 
single entry point for your services. The micro api serves HTTP and dynamically routes to the appropriate backend service.

<p align="center">
  <img src="https://github.com/micro/docs/blob/master/images/api.png" />
</p>

## How it works

The micro api builds on [go-micro](https://github.com/micro/go-micro), leveraging it for service discovery, load balancing, encoding and 
RPC based communication. Requests to the API are served over HTTP and internally routed via RPC. 

Because the micro api uses go-micro internally, this also makes it pluggable, so feel free to switch out consul service discovery for the 
kubernetes api or RPC for gRPC.

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

The API has three types of configurable request handlers.

1. API Handler: /[service]/[method]
	- Request/Response: api.Request/api.Response
	- The path is used to resolve service and method.
	- Requests are handled via API services which take the request api.Request and response api.Response types. 
	- Definitions for the Request/Response can be found at [go-api/proto](https://github.com/micro/go-api/blob/master/proto/api.proto)
	- The content type of the request/response body can be anything.
	- The default fallback handler where routes are not available.
	- Set via `--handler=api`
2. RPC Handler: /[service]/[method]
	- Request/Response: json/protobuf
	- An alternative to the default handler which uses the go-micro client to forward the request body as an RPC request.
	- Allows API handlers to be defined with concrete Go types.
	- Useful where you do not need full control of headers or request/response.
	- Can be used to run a single layer of backend services rather than additional API services.
	- Supported content-type `application/json` and `application/protobuf`.
	- Set via `--handler=rpc`
3. Reverse Proxy: /[service]
	- Request/Response: http
	- The request will be reverse proxied to the service resolved by the first element in the path
	- This allows REST to be implemented behind the API
	- Set via `--handler=proxy`
4. Event Handler: /[topic]/[event]
	- Async handler publishes request to message broker as an event
	- Request is formatted as [go-api/proto.Event](https://github.com/micro/go-api/blob/master/proto/api.proto#L28L39)
	- Set via `--handler=event`

Alternatively use the /rpc endpoint to speak to any service directly
- Expects params: `service`, `method`, `request`, optionally accepts `address` to target a specific host

```
curl -d 'service=go.micro.srv.greeter' \
	-d 'method=Say.Hello' \
	-d 'request={"name": "Bob"}' \
	http://localhost:8080/rpc
```

Find working examples in [github.com/micro/examples/api](https://github.com/micro/examples/tree/master/api)

### API Handler Request/Response Proto

The API Handler which is also the default handler expects services to use specific request and response protos, available 
at [go-api/proto](https://github.com/micro/go-api/blob/master/proto/api.proto). This allows the micro api to deconstruct 
a HTTP request into RPC and back to HTTP.

## Getting started

### Install

```shell
go get github.com/micro/micro
```

### Run

```shell
micro api
```

### ACME via Let's Encrypt

Serve securely by default using ACME via letsencrypt

```
micro --enable_acme api
```

Optionally specify a host whitelist

```
micro --enable_acme --acme_hosts=example.com,api.example.com api
```

### Serve Secure TLS

The API supports serving securely with TLS certificates

```shell
micro --enable_tls --tls_cert_file=/path/to/cert --tls_key_file=/path/to/key api
```

### Set Namespace

The API defaults to serving the namespace **go.micro.api**. The combination of namespace and request path 
are used to resolve an API service and method to send the query to.

```shell
micro api --namespace=com.example.api
```

## Examples

Here we have an example of a 3 tier architecture

- micro api (localhost:8080) - serving as the http entry point
- api service (go.micro.api.greeter) - serving a public facing api
- backend service (go.micro.srv.greeter) - internally scoped service

The full working example is [here](https://github.com/micro/examples/tree/master/greeter)

### Run Example

Prereq: Ensure you are running service discovery e.g consul agent -dev

Get examples

```
git clone https://github.com/micro/examples
```

Start the service go.micro.srv.greeter

```shell
go run examples/greeter/srv/main.go
```

Start the API service go.micro.api.greeter

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
curl "http://localhost:8080/greeter/say/hello?name=Asim+Aslam"
```

The HTTP path /greeter/say/hello maps to service go.micro.api.greeter method Say.Hello

Bypass the api service and call the backend directly via /rpc

```shell
curl -d 'service=go.micro.srv.greeter' \
	-d 'method=Say.Hello' \
	-d 'request={"name": "Asim Aslam"}' \
	http://localhost:8080/rpc
```

Make the same call entirely as JSON

```shell
$ curl -H 'Content-Type: application/json' \
	-d '{"service": "go.micro.srv.greeter", "method": "Say.Hello", "request": {"name": "Asim Aslam"}}' \
	http://localhost:8080/rpc
```

## Request Mapping

Micro dynamically routes to services using a fixed namespace and the HTTP path.

The default namespace for these services is **go.micro.api** but the namespace can be set via the `--namespace` flag. 

### API per Service

We promote a pattern of creating an API service per backend service for public facing traffic. This logically separates the concerns 
of serving an API frontend and backend services. 

### RPC Mapping

URLs are mapped as follows:

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

### REST Mapping

You can serve a RESTful API by using the API as a reverse proxy and implementing RESTful paths with libraries such as [go-restful](https://github.com/emicklei/go-restful). 
An example of a REST API service can be found at [greeter/api/rest](https://github.com/micro/examples/tree/master/greeter/api/rest).

Running the micro api with `--handler=proxy` will reverse proxy requests to services  within the API namespace.

Path	|	Service	|	Service Path
---	|	---	|	---
/foo	|	go.micro.api.foo	|	/foo
/foo/bar	|	go.micro.api.foo	|	/foo/bar
/greeter	|	go.micro.api.greeter	|	/greeter
/greeter/:name	|	go.micro.api.greeter	|	/greeter/:name

Using this handler means speaking HTTP directly with the backend service, ignoring any go-micro transport plugins.

## Stats Dashboard

Enable a stats dashboard via the `--enable_stats` flag. It will be exposed on /stats.

```shell
micro --enable_stats api
```

<img src="https://github.com/micro/docs/blob/master/images/stats.png">

