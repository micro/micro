# Getting Started

This is a guide to getting started with go-micro.

## Writing a service

The [server](https://godoc.org/github.com/micro/go-micro/server) package is the main component used to build a 
server. A default server is initialised for convenience. 

### Initialisation
The server can be initialised before usage. All the available init options can be found [here](https://godoc.org/github.com/micro/go-micro/server#Option).

```go
import "github.com/micro/go-micro/server"

server.Init(
	server.Name("go.micro.srv.greeter"),
	server.Version("1.0.0"),
)
```

### Defining the service and request/response

By default go-micro uses protobuf to define the service and request/response types. This is a very convenient way to strictly define the API. 

Here's an example definition:

hello.proto
```
syntax = "proto3";

package go.micro.srv.greeter;

service Say {
	rpc Hello(Request) returns (Response) {}
}

message Request {
	string name = 1;
}

message Response {
	string msg = 1;
}
```

As you can see we're defining a service handler called Say with the method Hello which takes the parameter Request type and returns Response.

We use protoc and protoc-gen-go to generate the concrete go implementation for this definition.

```shell
protoc --go_out=. hello.proto
```

Go-micro has experimental code generation support which provides client stub methods to reduce boiler plate code. This can be used in the following way. It uses a fork of [github.com/golang/protobuf](https://github.com/golang/protobuf).

```shell
go get github.com/micro/protobuf
protoc --go_out=plugins=micro:. hello.proto
```

The types generated can now be imported and used within a **handler** for a server or the client when making a request.

Here's part of the generated code.

```go
type Request struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

type Response struct {
	Msg string `protobuf:"bytes,1,opt,name=msg" json:"msg,omitempty"`
}

// Client API for Say service

type SayClient interface {
	Hello(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
}

type sayClient struct {
	c           client.Client
	serviceName string
}

func NewSayClient(serviceName string, c client.Client) SayClient {
	if c == nil {
		c = client.NewClient()
	}
	if len(serviceName) == 0 {
		serviceName = "go.micro.srv.greeter"
	}
	return &sayClient{
		c:           c,
		serviceName: serviceName,
	}
}

func (c *sayClient) Hello(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.serviceName, "Say.Hello", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Say service

type SayHandler interface {
	Hello(context.Context, *Request, *Response) error
}

func RegisterSayHandler(s server.Server, hdlr SayHandler) {
	s.Handle(s.NewHandler(hdlr))
}
```

### Handlers
The server requires **handlers** to be registered to serve requests. A handler is an public object with public methods which conform to the signature `func(ctx context.Context, req interface{}, rsp interface{}) error`.

A **streaming** handler maintains a connection with the client and can stream back multiple responses. It has the signature `func(ctx context.Context, req interface{}, rsp func(interface{}) error) error`.

Example handler:

```go
import (
	hello "github.com/micro/micro/examples/greeter/server/proto/hello"
)

type Say struct{}

func (s *Say) Hello(ctx context.Context, req *hello.Request, rsp *hello.Response) error {
	rsp.Msg = "Hello " + req.Name
	return nil
}
```

Example streaming handler:
```go
import (
	hello "github.com/micro/micro/examples/greeter/server/proto/hello"
)

type Say struct{}

func (s *Say) StreamHello(ctx context.Context, req *hello.Request, rspFn func(*hello.Response) error) error {
	i := 0
	for {
		i++
		if err := rspFn(&rsp{Msg: "Hello " + req.Name + " response: " + i}); err != nil {
			break
		}
	}
	return nil
}
```

Registration of the handler
```
func main() {
	server.Handle(
		server.NewHandler(
			new(Say), // Create new instance of Say struct
		),
	)
}
```

### Running the server

The server can be started by calling `server.Start()`. This causes the service to bind to the address in the config (which defaults to a random all interfaces and a random port) and listen for requests.

Alternatively `server.Run()` can be called which also registers the service with the **registry** providing service name resolution and discovery.

Starting the server
```go
if err := server.Run(); err != nil {
	log.Fatal(err)
}
```

### Complete Example

```go
package main

import (
	"log"

	"github.com/micro/go-micro/server"
	hello "github.com/micro/micro/examples/greeter/server/proto/hello"
	"golang.org/x/net/context"
)

type Say struct{}

func (s *Say) Hello(ctx context.Context, req *hello.Request, rsp *hello.Response) error {
	log.Print("Received Say.Hello request")
	rsp.Msg = server.Config().Id() + ": Hello " + req.Name
	return nil
}

func main() {
	// Initialise Server
	server.Init(
		server.Name("go.micro.srv.greeter"),
	)

	// Register Handlers
	server.Handle(
		server.NewHandler(
			new(Say),
		),
	)

	// Run server
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
```

Note: The registry/discovery system is required when using server.Run() since it's used for service name resolution. Read more about that in the go-micro [README](https://github.com/micro/go-micro).

## Writing a Client

The [client](https://godoc.org/github.com/micro/go-micro/client) package is used to query services. As with the server, a default client is initialised for convenience.

Querying the above service is as simple as this.

```go
package main

import (
	"fmt"

	"github.com/micro/go-micro/client"
	hello "github.com/micro/micro/examples/greeter/server/proto/hello"

	"golang.org/x/net/context"
)

func main() {
	req := client.NewRequest("go.micro.srv.greeter", "Say.Hello", &hello.Request{
		Name: "John",
	})

	rsp := &hello.Response{}

	// Call service
	if err := client.Call(ctx, req, rsp); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(rsp.Msg)
}
```

We can also use the generated client stub methods.

```go
package main

import (
	"fmt"

	"github.com/micro/go-micro/client"
	hello "github.com/micro/micro/examples/greeter/server/proto/hello"

	"golang.org/x/net/context"
)

func main() {
	// use the generated client stub
	cl := hello.NewSayClient("go.micro.srv.greeter", client.DefaultClient)

	rsp, err := cl.Hello(ctx, &hello.Request{
		Name: "John",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(rsp.Msg)
}
```

Note: again here the client is internally using the registry for service name resolution. You can alternatively use `client.CallRemote` or `client.StreamRemote` to directly call a specific address:port. 
