---
title: Errors
keywords: go-micro, framework, errors
tags: [go-micro, framework, errors]
sidebar: home_sidebar
permalink: /go-errors
summary: Error handling and errors produced by Go Micro
---

Go Micro provides abstractions and types for most things that occur in a distributed system including errors. By 
providing a core set of errors and the ability to define detailed error types we can consistently understand 
what's going on beyond the typical Go error string.

# Overview

We define the following error type:

```go
type Error struct {
    Id     string `json:"id"`
    Code   int32  `json:"code"`
    Detail string `json:"detail"`
    Status string `json:"status"`
}
```

Anywhere in the system where you're asked to return an error from a handler or receive one from a client you should assume 
its either a go-micro error or that you should produce one. By default we return 
[errors.InternalServerError](https://pkg.go.dev/github.com/micro/go-micro/v2/errors#InternalServerError) where somethin has 
gone wrong internally and [errors.Timeout](https://pkg.go.dev/github.com/micro/go-micro/v2/errors#Timeout) where a timeout occurred.

## Usage

Let's assume some error has occurred in your handler. You should then decide what kind of error to return and do the following.


Assuming some data provided was invalid

```go
return errors.BadRequest("com.example.srv.service", "invalid field")
```

In the event an internal error occurs

```go
if err != nil {
	return errors.InternalServerError("com.example.srv.service", "failed to read db: %v", err.Error())
}
```

Now lets say you receive some error from the client

```go
pbClient := pb.NewGreeterService("go.micro.srv.greeter", service.Client())
rsp, err := pb.Client(context, req)
if err != nil {
	// parse out the error
	e := errors.Parse(err.Error())

	// inspect the value
	if e.Code == 401 {
		// unauthorised...
	}
}
```

