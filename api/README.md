# micro api

The **micro api** is an API gateway for microservices. Use the API gateway [pattern](http://microservices.io/patterns/apigateway.html) to provide a 
single entry point for your services. The micro api serves HTTP and dynamically routes to the appropriate backend service.

<p align="center">
  <img src="https://github.com/micro/docs/blob/master/images/api.png" />
</p>

The micro api builds on [go-micro](https://github.com/micro/go-micro), leveraging it for service discovery, load balancing, encoding and 
RPC based communication. Requests to the API are served over HTTP and internally routed via RPC. Because the micro api uses go-micro internally, 
this also makes it pluggable, so feel free to switch out consul service discovery for the kubernetes api or RPC for gRPC.

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
	- Set via `--handler=proxy`.

Alternatively use the /rpc endpoint to speak to any service directly
- Expects params: `service`, `method`, `request`, optionally accepts `address` to target a specific host

```
curl -d 'service=go.micro.srv.greeter' \
	-d 'method=Say.Hello' \
	-d 'request={"name": "Bob"}' \
	http://localhost:8080/rpc
```

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

Below is an example of querying a service through the API

```
git clone https://github.com/micro/examples
```

### Run Example

Start the backend service go.micro.srv.greeter

```shell
go run examples/greeter/srv/main.go
```

Start the API service go.micro.api.greeter

```shell
go run examples/greeter/api/api.go
```

### Query

Make a HTTP call

```shell
curl "http://localhost:8080/greeter/say/hello?name=Asim+Aslam"
```

Make an RPC call via the /rpc

```shell
curl -d 'service=go.micro.srv.greeter' \
	-d 'method=Say.Hello' \
	-d 'request={"name": "Asim Aslam"}' \
	http://localhost:8080/rpc
```

Make an RPC call via /rpc with content-type set to json

```shell
$ curl -H 'Content-Type: application/json' \
	-d '{"service": "go.micro.srv.greeter", "method": "Say.Hello", "request": {"name": "Asim Aslam"}}' \
	http://localhost:8080/rpc
```

## API Request Mapping

Micro allows you resolve HTTP URL Paths at the edge to individual API Services. An API service is like any other 
micro service except each method signature takes an *api.Request and *api.Response type which can be found in 
[github.com/micro/go-api/proto](https://github.com/micro/go-api/blob/master/proto/api.proto).

The http.Request is deconstructed by the API into an api.Request and forwarded on to a backend API service. 
The api.Response is then constructed into a http.Response and returned to the client. The path of the request 
along with a namespace, is used to determine the backend service and method to call.

The default namespace for these services are **go.micro.api** but you can set your own namespace via `--namespace`.

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

A working example can be found here [Greeter Service](https://github.com/micro/examples/tree/master/greeter)

## Using REST

You can serve a RESTful API by using the API as a proxy and implementing RESTful paths with libraries such as [go-restful](https://github.com/emicklei/go-restful). 
An example of a REST API service can be found at [greeter/api/rest](https://github.com/micro/examples/tree/master/greeter/api/rest).

Starting the API with `--handler=proxy` will reverse proxy requests to backend services within the served API namespace (default: go.micro.api). 

Example

Path	|	Service	|	Service Path
---	|	---	|	---
/greeter	|	go.micro.api.greeter	|	/greeter
/greeter/:name	|	go.micro.api.greeter	|	/greeter/:name


Note: Using this method means directly speaking HTTP with the backend service. This eliminates the ability to switch transports.

## Stats Dashboard

You can enable a stats dashboard via the `--enable_stats` flag. It will be exposed on /stats.

```shell
micro --enable_stats api
```

<img src="https://github.com/micro/docs/blob/master/images/stats.png">



