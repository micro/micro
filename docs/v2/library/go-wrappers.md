---
title: Wrappers
keywords: go-micro, framework, wrappers
tags: [go-micro, framework, wrappers]
sidebar: home_sidebar
permalink: /go-wrappers
summary: Wrappers are middleware for Go Micro
---

# Overview

Wrappers are middleware for Go Micro. We want to create an extensible framework that includes hooks to add extra 
functionality that's not a core requirement. A lot of the time you need to execute something like auth, tracing, etc 
so this provides the ability to do that.

We use the "decorator pattern" for this.

## Usage

Here's some example usage and real code in [examples/wrapper](https://github.com/micro/examples/tree/master/wrapper).

You can find a range of wrappers in [go-plugins/wrapper](https://github.com/micro/go-plugins/tree/master/wrapper).

### Handler

Here's an example service handler wrapper which logs the incoming request

```go
// implements the server.HandlerWrapper
func logWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		fmt.Printf("[%v] server request: %s", time.Now(), req.Endpoint())
		return fn(ctx, req, rsp)
	}
}
```

It can be initialised when creating the service

```go
service := micro.NewService(
	micro.Name("greeter"),
	// wrap the handler
	micro.WrapHandler(logWrapper),
)
```

### Client

Here's an example of a client wrapper which logs requests made

```go
type logWrapper struct {
	client.Client
}

func (l *logWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	fmt.Printf("[wrapper] client request to service: %s endpoint: %s\n", req.Service(), req.Endpoint())
	return l.Client.Call(ctx, req, rsp)
}

// implements client.Wrapper as logWrapper
func logWrap(c client.Client) client.Client {
	return &logWrapper{c}
}
```

It can be initialised when creating the service

```go
service := micro.NewService(
	micro.Name("greeter"),
	// wrap the client
	micro.WrapClient(logWrap),
)
```

