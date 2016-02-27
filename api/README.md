# Micro API

This is a lightweight proxy for [Micro](https://github.com/micro/micro) based microservices. It conforms to the [API Gateway](http://microservices.io/patterns/apigateway.html) pattern and can be used in conjuction with [go-micro](https://github.com/micro/go-micro) based apps or any future language implementation of the [Micro](https://github.com/micro/micro) toolchain.

Currently a work in progress.

## Getting started

### Install the api

```bash
go get github.com/micro/micro
```

### Run the API

```bash
micro --logtostderr api
I0523 12:23:23.413940   81384 api.go:131] API Rpc handler /rpc
I0523 12:23:23.414238   81384 api.go:143] Listening on [::]:8080
I0523 12:23:23.414272   81384 server.go:113] Starting server go.micro.api id go.micro.api-1f951765-013e-11e5-9273-68a86d0d36b6
I0523 12:23:23.414355   81384 rpc_server.go:112] Listening on [::]:51938
I0523 12:23:23.414399   81384 server.go:95] Registering node: go.micro.api-1f951765-013e-11e5-9273-68a86d0d36b6
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

Let's start the example [go-micro](https://github.com/micro/go-micro) based server.
```bash
$ go get github.com/micro/go-micro/examples/server
$ $GOPATH/bin/server 
I0525 18:17:57.574457   84421 server.go:117] Starting server go.micro.srv.example id go.micro.srv.example-fccbb6fb-0301-11e5-9f1f-68a86d0d36b6
I0525 18:17:57.574748   84421 rpc_server.go:126] Listening on [::]:62421
I0525 18:17:57.574779   84421 server.go:99] Registering node: go.micro.srv.example-fccbb6fb-0301-11e5-9f1f-68a86d0d36b6
```

The example server has a handler registered called Example with a method named Call. 
Now let's query this through the API. 
```bash
$ curl -d 'service=go.micro.srv.example' -d 'method=Example.Call' -d 'request={"name": "Asim Aslam"}' http://localhost:8080/rpc
{"msg":"go.micro.srv.example-fccbb6fb-0301-11e5-9f1f-68a86d0d36b6: Hello Asim Aslam"}
```

Alternatively let's try 'Content-Type: application/json'
```bash
$ curl -H 'Content-Type: application/json' -d '{"service": "go.micro.srv.example", "method": "Example.Call", "request": {"name": "Asim Aslam"}}' http://localhost:8080/rpc
{"msg":"go.micro.srv.example-fccbb6fb-0301-11e5-9f1f-68a86d0d36b6: Hello Asim Aslam"}
```

## Testing using REST based API Services

Micro allows you to handle REST based paths using rpc by providing built in handling for API Services. An API service is like any other 
micro service except each method signature takes an *api.Request and *api.Response which can be found in 
[github.com/micro/micro/api/proto](https://github.com/micro/micro/tree/master/api/proto).

The default namespace for these services are: go.micro.api

Translation of URLs are as follows:

/foo/bar => service: go.micro.api.foo method: Foo.Bar

/foo/bar/baz => service: go.micro.api.foo method: Bar.Baz

/foo/bar/baz/cat => service: go.micro.api.foo.bar method: Baz.Cat

A working example can be found here [Greeter Service](https://github.com/micro/micro/tree/master/examples/greeter)
