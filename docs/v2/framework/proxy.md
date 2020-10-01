---
title: Service Proxy
keywords: proxy
tags: [proxy]
sidebar: home_sidebar
permalink: /proxy
summary: The micro proxy is a service-to-service proxy.
---

A service proxy is a server which acts as an intermediary for requests from one service to another.

<img src="images/proxy.svg" />

## Overview

The micro proxy provides a proxy implementation of the go-micro framework. This consolidates go-micro features into a single location which allows 
offloading service discovery, load balancing, fault tolerance, plugins, wrappers, etc to the proxy itself. Rather than updating every Go Micro 
app for infrastructure level concerns, its easier to put them in the proxy. It also allows any language to be integrated with a thin client 
rather than having to implement all the features.

## Run Proxy

Start the proxy

```shell
micro proxy
```

The server address is dynamic but can be configured as follows.

```
MICRO_SERVER_ADDRESS=localhost:9090 micro proxy
```

## Proxy Services

Now the proxy is running you can quite simply proxy requests through it.

Start your go micro app like so

```
MICRO_PROXY=go.micro.proxy go run main.go
```

Your service will lookup the proxy in discovery then use it to route any requests. If multiple proxies exist it will balance 
it's requests across them. It will also cache the proxy addresses locally.


If you would rather send requests through a single proxy specify it's address like so.

```
MICRO_PROXY=localhost:9090 go run main.go
```

Ensure the proxy is running on the address specified.

```
MICRO_SERVER_ADDRESS=localhost:9090 micro proxy
```

## Single Endpoint

Use the proxy as a front proxy for a single endpoint

```
MICRO_SERVER_NAME=helloworld \
MICRO_PROXY_ENDPOINT=localhost:10001 \
micro proxy
```

All requests to helloworld will be sent to the backend at localhost:10001

