# Micro API

This is a lightweight proxy for [Micro](https://github.com/myodc/micro) based microservices. It conforms to the [API Gateway](http://microservices.io/patterns/apigateway.html) pattern and can be used in conjuction with [go-micro](https://github.com/myodc/go-micro) based apps or any future language implementation of the [Micro](https://github.com/myodc/micro) toolchain.

Currently a work in progress.

### Run API
```bash
$ go get github.com/myodc/micro
$ micro api
I0523 12:23:23.413940   81384 api.go:131] API Rpc handler /rpc
I0523 12:23:23.414238   81384 api.go:143] Listening on [::]:8080
I0523 12:23:23.414272   81384 server.go:113] Starting server go.micro.api id go.micro.api-1f951765-013e-11e5-9273-68a86d0d36b6
I0523 12:23:23.414355   81384 rpc_server.go:112] Listening on [::]:51938
I0523 12:23:23.414399   81384 server.go:95] Registering node: go.micro.api-1f951765-013e-11e5-9273-68a86d0d36b6
```

### Testing API

Let's start the template [go-micro](https://github.com/myodc/go-micro) based service.
```bash
$ go get github.com/myodc/go-micro/template
$ $GOPATH/bin/template 
I0523 12:21:09.506998   77096 server.go:113] Starting server go.micro.service.template id go.micro.service.template-cfc481fc-013d-11e5-bcdc-68a86d0d36b6
I0523 12:21:09.507281   77096 rpc_server.go:112] Listening on [::]:51868
I0523 12:21:09.507329   77096 server.go:95] Registering node: go.micro.service.template-cfc481fc-013d-11e5-bcdc-68a86d0d36b6
```

The template service has a handler registered called Example with a method named Call. 
Now let's query this through the API. 
```bash
$ curl -d 'service=go.micro.service.template' -d 'method=Example.Call' -d 'request={"name": "Asim Aslam"}' http://localhost:8080/rpc
{"msg":"go.micro.service.template-cfc481fc-013d-11e5-bcdc-68a86d0d36b6: Hello Asim Aslam"}
```

Alternatively let's try 'Content-Type: application/json'
```bash
$ curl -H 'Content-Type: application/json' -d '{"service": "go.micro.service.template", "method": "Example.Call", "request": {"name": "Asim Aslam"}}' http://localhost:8080/rpc
{"msg":"go.micro.service.template-cfc481fc-013d-11e5-bcdc-68a86d0d36b6: Hello Asim Aslam"}
```

### Testing using REST based API Services

Micro allows you to handle REST based paths using rpc by providing built in handling for API Services. An API service is like any other 
micro service except each method signature takes an *api.Request and *api.Response which can be found in 
[github.com/myodc/micro/api/proto](https://github.com/myodc/micro/tree/master/api/proto).

The default namespace for these services are: go.micro.api

Translation of URLs are as follows:

/foo/bar => service: go.micro.api.foo method: Foo.Bar

/foo/bar/baz => service: go.micro.api.foo method: Bar.Baz

/foo/bar/baz/cat => service: go.micro.api.foo.bar method: Baz.Cat

A working example can be found here [Greeter Service](https://github.com/myodc/micro/tree/master/examples/greeter)
