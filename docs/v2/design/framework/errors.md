# Errors

Micro provides its own form of structured errors to add detail to otherwise opaque strings.

## Overview

Structured errors enable us to give an error an ID, code, and detail. This is useful when 
debugging a layered stack especially related to microservices. We want to evolve these 
structured errors into something more widely applicable beyond RPC.

## Design

The current design

```go
type Error struct {
	Id                   string
	Code                 int32
	Detail               string
	Status               string
}

func (e *Error) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}
```

Find the current implementation in [github.com/micro/go-micro/v2/errors](https://github.com/micro/go-micro/blob/master/errors/errors.go)

When output the error appears in json format

```
{"id": "go.micro.client", "code": 500, "detail": "an error occurred calling the service", "status": "internal server error"}
```

These errors are passed via the `Micro-Error` header or in grpc encapsulated within their error type.
