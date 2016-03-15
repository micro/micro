# Micro API

This is a lightweight proxy for [Micro](https://github.com/micro/micro) based microservices. It conforms to the [API Gateway](http://microservices.io/patterns/apigateway.html) pattern and can be used in conjuction with [go-micro](https://github.com/micro/go-micro) based apps or any future language implementation of the [Micro](https://github.com/micro/micro) toolchain.

## Handlers

The API handles requests in three ways.

1. /rpc
	- Sends requests directly to backend services using JSON
	- Expects params: service, method, request.
2. /[service]/[method]
	- The path is used to resolve service and method. 
	- Requests are handled via API services which take the request api.Request and response api.Response types. 
	- Definitions for the Request/Response can be found at [micro/api/proto](https://github.com/micro/micro/tree/master/api/proto)

3. /[service] reverse proxy set via `--api_handler=proxy`
	- The request will be reverse proxied to the service resolved by the first element in the path
	- This allows REST to be implemented behind the API

## Getting started

### Install the api

```bash
go get github.com/micro/micro
```

### Run the API

```bash
micro api
2016/03/15 20:53:19 Registering RPC Handler at /rpc
2016/03/15 20:53:19 Registering API Handler at /
2016/03/15 20:53:19 Listening on [::]:8080
2016/03/15 20:53:19 Listening on [::]:60971
2016/03/15 20:53:19 Broker Listening on [::]:60972
2016/03/15 20:53:19 Registering node: go.micro.api-f2ffeebf-eaef-11e5-817c-68a86d0d36b6
```

### Serve Secure TLS

The API supports serving securely with TLS certificates

```bash
micro --enable_tls --tls_cert_file=/path/to/cert --tls_key_file=/path/to/key api
```

### Set Namespace

The API defaults to serving the namespace **go.micro.api**. The combination of namespace and request path 
are used to resolve an API service and method to send the query to. 

```bash
micro --api_namespace=com.example.api
```

## Testing API

### Run the example app

Let's start the example [go-micro](https://github.com/micro/go-micro) based server.

```bash
$ go get github.com/micro/go-micro/examples/server
$ $GOPATH/bin/server 
I0525 18:17:57.574457   84421 server.go:117] Starting server go.micro.srv.example id go.micro.srv.example-fccbb6fb-0301-11e5-9f1f-68a86d0d36b6
I0525 18:17:57.574748   84421 rpc_server.go:126] Listening on [::]:62421
I0525 18:17:57.574779   84421 server.go:99] Registering node: go.micro.srv.example-fccbb6fb-0301-11e5-9f1f-68a86d0d36b6
```

### Query RPC via curl

The example server has a handler registered called Example with a method named Call. Now let's query this through the API.

```bash
$ curl -d 'service=go.micro.srv.example' \
	-d 'method=Example.Call' \
	-d 'request={"name": "Asim Aslam"}' \
	http://localhost:8080/rpc

{"msg":"go.micro.srv.example-fccbb6fb-0301-11e5-9f1f-68a86d0d36b6: Hello Asim Aslam"}
```

Alternatively let's try 'Content-Type: application/json'

```bash
$ curl -H 'Content-Type: application/json' \
	-d '{"service": "go.micro.srv.example", "method": "Example.Call", "request": {"name": "Asim Aslam"}}' \
	http://localhost:8080/rpc

{"msg":"go.micro.srv.example-fccbb6fb-0301-11e5-9f1f-68a86d0d36b6: Hello Asim Aslam"}
```

## API HTTP request translation

Micro allows you resolve HTTP Paths at the edge to individual API Services. An API service is like any other 
micro service except each method signature takes an *api.Request and *api.Response type which can be found in 
[github.com/micro/micro/api/proto](https://github.com/micro/micro/tree/master/api/proto).

The http.Request is deconstructed by the API into an api.Request and forwarded on to a backend API service. 
The api.Response is then constructed into a http.Response and returned to the client. The path of the request 
along with a namespace, is used to determine the backend service and method to call.

The default namespace for these services are **go.micro.api** but you can set your own namespace via `--api_namespace`.

Translation of URLs are as follows:

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
An example of a REST API service can be found at [greeter/api/go-restful](https://github.com/micro/micro/tree/master/examples/greeter/api/go-restful).

Starting the API with `--api_handler=proxy` will reverse proxy requests to backend services within the served API namespace (default: go.micro.api). 

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

<img src="https://github.com/micro/micro/blob/master/doc/stats.png">



