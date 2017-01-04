# Micro API

The **micro api** is a lightweight proxy for [micro](https://github.com/micro/micro) based microservices. It conforms to the [API Gateway](http://microservices.io/patterns/apigateway.html) pattern and can be used in conjuction with [go-micro](https://github.com/micro/go-micro) based apps or any future language implementation of the [micro](https://github.com/micro/micro) toolkit.

<p align="center">
  <img src="https://raw.githubusercontent.com/micro/micro/master/doc/images/api.png" />
</p>


## Handlers

The API handles requests in three ways.

1. API Handler: /[service]/[method]
	- Request/Response: api.Request/api.Response
	- The path is used to resolve service and method.
	- Requests are handled via API services which take the request api.Request and response api.Response types. 
	- Definitions for the Request/Response can be found at [micro/api/proto](https://github.com/micro/micro/tree/master/api/proto)
	- The content type of the request/response body can be anything.
	- The default handler
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


Alternatively use the /rpc send requests directly to backend services using JSON 
- Expects params: `service`, `method`, `request`, optionally accepts `address` to target a specific host
- ```curl -d 'service=go.micro.srv.greeter' -d 'method=Say.Hello' -d 'request={"name": "Bob"}' http://localhost:8080/rpc```

## Getting started

### Install

```shell
go get github.com/micro/micro
```

### Run

```shell
micro api
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

## API

### Run Services

Start the backend service go.micro.srv.greeter

```shell
go run examples/greeter/server/main.go 
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
[github.com/micro/micro/api/proto](https://github.com/micro/micro/tree/master/api/proto).

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

A working example can be found here [Greeter Service](https://github.com/micro/micro/tree/master/examples/greeter)

## Using REST

You can serve a RESTful API by using the API as a proxy and implementing RESTful paths with libraries such as [go-restful](https://github.com/emicklei/go-restful). 
An example of a REST API service can be found at [greeter/api/rest](https://github.com/micro/micro/tree/master/examples/greeter/api/rest).

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

<img src="https://raw.githubusercontent.com/micro/micro/master/doc/images/stats.png">



