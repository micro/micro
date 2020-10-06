---
title: Writing a Go Function
keywords: development
tags: [development]
sidebar: home_sidebar
permalink: /writing-a-go-function
summary: 
---

This is a guide to getting started with go-micro functions. Functions are one time executing Services.

If you prefer a higher level overview of the toolkit first, checkout the introductory blog post [https://micro.mu/blog/2016/03/20/micro.html](https://micro.mu/blog/2016/03/20/micro.html).

## Writing a Function

The top level [Function](https://pkg.go.dev/github.com/micro/go-micro/v2#Function) interface is the main component for the 
function programming model in go-micro. It encapsulates the Service interface while providing one time execution.

```go
// Function is a one time executing Service
type Function interface {
	// Inherits Service interface
	Service
	// Done signals to complete execution
	Done() error
	// Handle registers an RPC handler
	Handle(v interface{}) error
	// Subscribe registers a subscriber
	Subscribe(topic string, v interface{}) error
}
```

### 1. Initialisation

A function is created like so using `micro.NewFunction`.

```go
import "github.com/micro/go-micro/v2"

function := micro.NewFunction() 
```

Options can be passed in during creation.

```go
function := micro.NewFunction(
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

function := micro.NewFunction(
        micro.Flags(
                cli.StringFlag{
                        Name:  "environment",
                        Usage: "The environment",
                },
        )
)
```

To parse flags use `function.Init`. Additionally access flags use the `micro.Action` option.

```go
function.Init(
        micro.Action(func(c *cli.Context) {
                env := c.StringFlag("environment")
                if len(env) > 0 {
                        fmt.Println("Environment set to", env)
                }
        }),
)
```

Go Micro provides predefined flags which are set and parsed if `function.Init` is called. See all the flags 
[here](https://pkg.go.dev/github.com/micro/go-micro/v2/cmd#pkg-variables).

### 2. Defining the API

We use protobuf files to define the API interface. This is a very convenient way to strictly define the API and 
provide concrete types for both the server and a client.

Here's an example definition.

greeter.proto

```proto
syntax = "proto3";

service Greeter {
	rpc Hello(HelloRequest) returns (HelloResponse) {}
}

message HelloRequest {
	string name = 1;
}

message HelloResponse {
	string greeting = 2;
}
```

Here we're defining a function handler called Greeter with the method Hello which takes the parameter HelloRequest type and returns HelloResponse.

### 3. Generate the API interface

We use protoc and protoc-gen-go to generate the concrete go implementation for this definition.

Go-micro uses code generation to provide client stub methods to reduce boiler plate code much like gRPC. It's done via a protobuf plugin 
which requires a fork of [golang/protobuf](https://github.com/golang/protobuf) that can be found at 
[github.com/micro/protobuf](https://github.com/micro/protobuf).

```shell
go get github.com/micro/protobuf/{proto,protoc-gen-go}
protoc --go_out=plugins=micro:. greeter.proto
```

The types generated can now be imported and used within a **handler** for a server or the client when making a request.

Here's part of the generated code.

```go
type HelloRequest struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

type HelloResponse struct {
	Greeting string `protobuf:"bytes,2,opt,name=greeting" json:"greeting,omitempty"`
}

// Client API for Greeter service

type GreeterClient interface {
	Hello(ctx context.Context, in *HelloRequest, opts ...client.CallOption) (*HelloResponse, error)
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

func (c *greeterClient) Hello(ctx context.Context, in *HelloRequest, opts ...client.CallOption) (*HelloResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Greeter.Hello", in)
	out := new(HelloResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Greeter service

type GreeterHandler interface {
	Hello(context.Context, *HelloRequest, *HelloResponse) error
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
        Hello(context.Context, *HelloRequest, *HelloResponse) error
}
```

Here's an implementation of the Greeter handler.

```go
import proto "github.com/micro/examples/service/proto"

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *proto.HelloRequest, rsp *proto.HelloResponse) error {
	rsp.Greeting = "Hello " + req.Name
	return nil
}
```

The handler is registered much like a http.Handler.

```
function := micro.NewFunction(
	micro.Name("greeter"),
)

proto.RegisterGreeterHandler(service.Server(), new(Greeter))
```

Alternatively the Function interface provides a simpler registration pattern.

```
function := micro.NewFunction(
        micro.Name("greeter"),
)

function.Handle(new(Greeter))
```

You can also register an async subscriber using the Subscribe method.

### 5. Running the function

The function can be run by calling `function.Run`. This causes it to bind to the address in the config 
(which defaults to the first RFC1918 interface found and a random port) and listen for requests.

This will additionally Register the function with the registry on start and Deregister when issued a kill signal.

```go
if err := function.Run(); err != nil {
	log.Fatal(err)
}
```

Upon serving a request the function will exit. You can use [micro run](https://micro.mu/docs/run.html) to manage the lifecycle 
of functions. A complete example can be found at [examples/function](https://github.com/micro/examples/tree/master/function).

### 6. The complete function
<br>
greeter.go

```go
package main

import (
        "log"

        "github.com/micro/go-micro/v2"
        proto "github.com/micro/examples/function/proto"

        "golang.org/x/net/context"
)

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *proto.HelloRequest, rsp *proto.HelloResponse) error {
        rsp.Greeting = "Hello " + req.Name
        return nil
}

func main() {
        function := micro.NewFunction(
                micro.Name("greeter"),
                micro.Version("latest"),
        )

        function.Init()

	function.Handle(new(Greeter))
	
        if err := function.Run(); err != nil {
                log.Fatal(err)
        }
}
```

Note. The service discovery mechanism will need to be running so the function can register to be discovered by those wishing 
to query it. A quick getting started for that is [here](https://github.com/micro/go-micro#getting-started).

## Writing a Client

The [client](https://pkg.go.dev/github.com/micro/go-micro/v2/client) package is used to query functions and services. When you create a 
Function, a Client is included which matches the initialised packages used by the server.

Querying the above function is as simple as the following.

```go
// create the greeter client using the service name and client
greeter := proto.NewGreeterClient("greeter", function.Client())

// request the Hello method on the Greeter handler
rsp, err := greeter.Hello(context.TODO(), &proto.HelloRequest{
	Name: "John",
})
if err != nil {
	fmt.Println(err)
	return
}

fmt.Println(rsp.Greeter)
```

`proto.NewGreeterClient` takes the function name and the client used for making requests.

The full example is can be found at [go-micro/examples/function](https://github.com/micro/examples/tree/master/function).

