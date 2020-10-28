---
title: Writing a Service
keywords: development
tags: [development]
sidebar: home_sidebar
permalink: /writing-a-go-service
summary: 
---

This is a more detailed guide to writing a service and the internals related to it. The top level [Service](https://pkg.go.dev/github.com/micro/go-micro/v2#Service) 
interface is the main component for building a service. It wraps all the underlying packages of Go Micro into a single convenient interface.

```go
type Service interface {
    Init(...Option)
    Options() Options
    Client() client.Client
    Server() server.Server
    Run() error
    String() string
}
```

### 1. Initialisation

A service is created like so using `micro.NewService`.

```go
import "github.com/micro/go-micro/v2"

service := micro.NewService() 
```

Options can be passed in during creation.

```go
service := micro.NewService(
        micro.Name("greeter"),
        micro.Version("latest"),
)
```

All the available options can be found [here](https://pkg.go.dev/github.com/micro/go-micro/v2#Option).

Go Micro also provides a way to set command line flags using `micro.Flags`.

```go
import (
        "github.com/micro/cli"
        "github.com/micro/go-micro/v2"
)

service := micro.NewService(
        micro.Flags(
                cli.StringFlag{
                        Name:  "environment",
                        Usage: "The environment",
                },
        )
)
```

To parse flags use `service.Init`. Additionally access flags use the `micro.Action` option.

```go
service.Init(
        micro.Action(func(c *cli.Context) {
                env := c.StringFlag("environment")
                if len(env) > 0 {
                        fmt.Println("Environment set to", env)
                }
        }),
)
```

Go Micro provides predefined flags which are set and parsed if `service.Init` is called. See all the flags 
[here](https://pkg.go.dev/github.com/micro/go-micro/v2/config/cmd#pkg-variables).

### 2. Defining the API

We use protobuf files to define the service API interface. This is a very convenient way to strictly define the API and 
provide concrete types for both the server and a client.

Here's an example definition.

greeter.proto

```proto
syntax = "proto3";

service Greeter {
	rpc Hello(Request) returns (Response) {}
}

message Request {
	string name = 1;
}

message Response {
	string greeting = 2;
}
```

Here we're defining a service handler called Greeter with the method Hello which takes the parameter Request type and returns Response.

### 3. Generate the API interface

You'll need the following to generate protobuf code

- [protoc](https://github.com/google/protobuf)
- [protoc-gen-go](https://github.com/golang/protobuf)
- [protoc-gen-micro](https://github.com/micro/micro/tree/master/cmd/protoc-gen-micro)

We use protoc, protoc-gen-go and protoc-gen-micro to generate the concrete go implementation for this definition.

```shell
go get github.com/golang/protobuf/{proto,protoc-gen-go}
```

```shell
go get github.com/micro/micro/v2/cmd/protoc-gen-micro@master
```

```
protoc --proto_path=$GOPATH/src:. --micro_out=. --go_out=. greeter.proto
```

The types generated can now be imported and used within a **handler** for a server or the client when making a request.

Here's part of the generated code.

```go
type Request struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

type Response struct {
	Greeting string `protobuf:"bytes,2,opt,name=greeting" json:"greeting,omitempty"`
}

// Client API for Greeter service

type GreeterClient interface {
	Hello(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
}

type greeterClient struct {
	c           client.Client
	serviceName string
}

func NewGreeterClient(serviceName string, c client.Client) GreeterClient {
	if c == nil {
		c = client.NewClient()
	}
	if len(serviceName) == 0 {
		serviceName = "greeter"
	}
	return &greeterClient{
		c:           c,
		serviceName: serviceName,
	}
}

func (c *greeterClient) Hello(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.serviceName, "Greeter.Hello", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Greeter service

type GreeterHandler interface {
	Hello(context.Context, *Request, *Response) error
}

func RegisterGreeterHandler(s server.Server, hdlr GreeterHandler) {
	s.Handle(s.NewHandler(&Greeter{hdlr}))
}
```

### 4. Implement the handler

The server requires **handlers** to be registered to serve requests. A handler is an public type with public methods 
which conform to the signature `func(ctx context.Context, req interface{}, rsp interface{}) error`.

As you can see above, a handler signature for the Greeter interface looks like so.

```go
type GreeterHandler interface {
        Hello(context.Context, *Request, *Response) error
}
```

Here's an implementation of the Greeter handler.

```go
import proto "github.com/micro/examples/service/proto"

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	rsp.Greeting = "Hello " + req.Name
	return nil
}
```


The handler is registered with your service much like a http.Handler.

```
service := micro.NewService(
	micro.Name("greeter"),
)

pb.RegisterGreeterHandler(service.Server(), new(Greeter))
```

You can also create a bidirectional streaming handler but we'll leave that for another day.

### 5. Running the service

The service can be run by calling `server.Run`. This causes the service to bind to the address in the config 
(which defaults to the first RFC1918 interface found and a random port) and listen for requests.

This will additionally Register the service with the registry on start and Deregister when issued a kill signal.

```go
if err := service.Run(); err != nil {
	log.Fatal(err)
}
```

### 6. The complete service
<br>
greeter.go

```go
package main

import (
        "log"

        "github.com/micro/go-micro/v2"
        pb "github.com/micro/examples/service/proto"

        "golang.org/x/net/context"
)

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
        rsp.Greeting = "Hello " + req.Name
        return nil
}

func main() {
        service := micro.NewService(
                micro.Name("greeter"),
        )

        service.Init()

        pb.RegisterGreeterHandler(service.Server(), new(Greeter))

        if err := service.Run(); err != nil {
                log.Fatal(err)
        }
}
```

Note. The service discovery mechanism will need to be running so the service can register to be discovered by clients and 
other services. A quick getting started for that is [here](https://github.com/micro/go-micro#getting-started).

## Writing a Client

The [client](https://pkg.go.dev/github.com/micro/go-micro/v2/client) package is used to query services. When you create a 
Service, a Client is included which matches the initialised packages used by the server.

Querying the above service is as simple as the following.

```go
// create the greeter client using the service name and client
greeter := pb.NewGreeterService("greeter", service.Client())

// request the Hello method on the Greeter handler
rsp, err := greeter.Hello(context.TODO(), &pb.Request{
	Name: "John",
})
if err != nil {
	fmt.Println(err)
	return
}

fmt.Println(rsp.Greeting)
```

`pb.NewGreeterClient` takes the service name and the client used for making requests.

The full example is can be found at [go-micro/examples/service](https://github.com/micro/examples/tree/master/service).

