---
title: PubSub
keywords: go-micro, framework, pubsub
tags: [go-micro, framework, pubsub]
sidebar: home_sidebar
permalink: /go-pubsub
summary: Go Micro builds in PubSub for event driven microservices
---

# Overview

Microservices is an event driven architecture patterna and so Go Micro builds in the concept of asynchronous messaging 
using a message broker interface. It seamlessly operates on protobuf types for you. Automatically encoding and decoding 
messages as they are sent and received from the broker.

By default go-micro includes a point-to-point http broker but this can be swapped out via [go-plugins](https://github.com/micro/go-plugins).

### Publish Message

Create a new publisher with a `topic` name and service client

```go
p := micro.NewEvent("events", service.Client())
```

Publish a proto message

```go
p.Publish(context.TODO(), &proto.Event{Name: "event"})
```

### Subscribe

Create a message handler. It's signature should be `func(context.Context, v interface{}) error`.

```go
func ProcessEvent(ctx context.Context, event *proto.Event) error {
	fmt.Printf("Got event %+v\n", event)
	return nil
}
```

Register the message handler with a `topic`

```go
micro.RegisterSubscriber("events", ProcessEvent)
```

See [examples/pubsub](https://github.com/micro/examples/tree/master/pubsub) for a complete example.

