---
title: Go API
keywords: go-api
tags: [go-micro]
sidebar: home_sidebar
permalink: /go-api
summary: a package for defining api routes and handlers
---

# Overview

Go API is a pluggable API framework driven by service discovery to help build powerful public API gateways.

The Go API library provides api gateway routing capabilities. A microservice architecture decouples application logic into 
separate service. An api gateway provides a single entry point to consolidate these services into a unified api. The 
Go API uses routes defined in service discovery metadata to generate routing rules and serve http requests.

<img src="https://micro.mu/docs/images/go-api.png?v=1" alt="Go API" />

Go API is the basis for the [micro api](https://micro.mu/docs/api.html).

## Handlers

Handlers are http handlers used for handling requests. It uses the `http.Handler` pattern for convenience.

- [`api`](#api-handler) - Handles any HTTP request. Gives full control over the http request/response via RPC.
- [`broker`](#broker-handler) - A http handler which implements the go-micro broker interface
- [`cloudevents`](#cloudevents-handler) -  Handles CloudEvents and publishes to a message bus.
- [`event`](#event-handler) -  Handles any HTTP request and publishes to a message bus.
- [`http`](#http-handler) - Handles any HTTP request and forwards as a reverse proxy.
- [`registry`](#registry-handler) - A http handler which implements the go-micro registry interface
- [`rpc`](#rpc-handler) - Handles json and protobuf POST requests. Forwards as RPC.
- [`web`](#web-handler) - The HTTP handler with web socket support included.

## API Handler

The API handler is the default handler. It serves any HTTP requests and forwards on as an RPC request with a specific format.

- Content-Type: Any
- Body: Any
- Forward Format: [api.Request](https://github.com/micro/go-micro/blob/master/api/proto/api.proto#L11)/[api.Response](https://github.com/micro/go-micro/blob/master/api/proto/api.proto#L21)
- Path: `/[service]/[method]`
- Resolver: Path is used to resolve service and method

## Broker Handler

The broker handler is a http handler which serves the go-micro broker interface

- Content-Type: Any
- Body: Any
- Forward Format: HTTP
- Path: `/`
- Resolver: Topic is specified as a query param

Post the request and it will be published

## CloudEvents Handler

The CloudEvents handler serves HTTP and forwards the request as a CloudEvents message over a message bus using the go-micro/client.Publish method.

- Content-Type: Any
- Body: Any
- Forward Format: Request is formatted as [CloudEvents](https://github.com/cloudevents/spec) message
- Path: `/[topic]`
- Resolver: Path is used to resolve topic

## Event Handler

The event handler serves HTTP and forwards the request as a message over a message bus using the go-micro/client.Publish method.

- Content-Type: Any
- Body: Any
- Forward Format: Request is formatted as [go-api/proto.Event](https://github.com/micro/go-api/blob/master/proto/api.proto#L28L39) 
- Path: `/[topic]/[event]`
- Resolver: Path is used to resolve topic and event name

## HTTP Handler

The http handler is a http reverse proxy with built in service discovery.

- Content-Type: Any
- Body: Any
- Forward Format: HTTP Reverse proxy
- Path: `/[service]`
- Resolver: Path is used to resolve service name

## Registry Handler

The registry handler is a http handler which serves the go-micro registry interface

- Content-Type: Any
- Body: JSON
- Forward Format: HTTP
- Path: `/`
- Resolver: GET, POST, DELETE used to get service, register or deregister

## RPC Handler

The RPC handler serves json or protobuf HTTP POST requests and forwards as an RPC request.

- Content-Type: `application/json` or `application/protobuf`
- Body: JSON or Protobuf
- Forward Format: json-rpc or proto-rpc based on content
- Path: `/[service]/[method]`
- Resolver: Path is used to resolve service and method

## Web Handler

The web handler is a http reverse proxy with built in service discovery and web socket support.

- Content-Type: Any
- Body: Any
- Forward Format: HTTP Reverse proxy including web sockets
- Path: `/[service]`
- Resolver: Path is used to resolve service name


