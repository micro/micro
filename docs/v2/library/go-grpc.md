---
title: gRPC Service
keywords: grpc
tags: [go-grpc]
sidebar: home_sidebar
permalink: /go-grpc
summary: gRPC service integration with go-micro
---

# Overview

Our gRPC service makes use of [go-micro](https://github.com/micro/go-micro) plugins to create a simpler framework for gRPC development. It interoperates with 
standard gRPC services seamlessly, including the [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway). The grpc service uses 
the go-micro client and server plugins. Because gRPC is a tightly coupled protocol and framework we ignore 
the go-micro codec and transport plugins.

Internally we make use of the gRPC framework but hide the complexity.

## Examples

Find an example greeter service in [examples/grpc](https://github.com/micro/examples/tree/master/grpc)

## Getting Started

- [Install Protobuf](#install-protobuf)
- [Service Discovery](#service-discovery)
- [Writing a Service](#writing-a-service)
- [Using with Micro](#use-with-micro)
- [Using with gRPC Gateway](#use-with-grpc-gateway)


## Install Protobuf

Protobuf is required for code generation

You'll need to install:

- [protoc-gen-micro](https://github.com/micro/micro/tree/master/cmd/protoc-gen-micro)

## Writing a Service

Go-grpc service is identical to a go-micro service. Which means you can swap out `micro.NewService` for `grpc.NewService` 
with zero other code changes.

```go
package main

import (
	"context"
	"time"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/service/grpc"
	hello "github.com/micro/go-micro/v2/service/grpc/examples/greeter/server/proto/hello"
)

type Say struct{}

func (s *Say) Hello(ctx context.Context, req *hello.Request, rsp *hello.Response) error {
	rsp.Msg = "Hello " + req.Name
	return nil
}

func main() {
	service := grpc.NewService(
		micro.Name("greeter"),
	)

	service.Init()

	hello.RegisterSayHandler(service.Server(), new(Say))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
```

## Use with Micro

You may want to use the micro toolkit with grpc services. Simply flag flip or use env vars to set the 
grpc client and server like below.

### Using env vars

```
MICRO_CLIENT=grpc MICRO_SERVER=grpc micro api
```

### Using flags

```shell
micro --client=grpc --server=grpc api
```

## Use with gRPC Gateway

The micro gRPC plugins seamlessly integrates with the gRPC ecosystem. This means the grpc-gateway can be used as per usual.

Find an example greeter api at [examples/grpc/gateway](https://github.com/micro/examples/tree/master/grpc/gateway).

