---
title: CLI
keywords: cli
tags: [cli]
sidebar: home_sidebar
summary: 
---

# micro cli

The **micro cli** is a command line interface for the micro toolkit [micro](https://github.com/micro/micro). 

## Getting Started

- [Install](#install)
- [Interactive Mode](#interactive-mode)
- [List Services](#list-services)
- [Get Service](#get-service)
- [Call Service](#call-service)
- [Service Health](#service-health)
- [Proxy Remote Environment](#proxy-remote-env)

## Install

```shell
go get github.com/micro/micro/v2
```

## Interactive Mode

To use the cli as an interactive prompt

```
micro cli
```

Remove `micro` from the below commands when in interactive mode

## Example Usage


### List Services

```shell
micro list services
```

### Get Service

```shell
micro get service go.micro.srv.example
```

Output

```
go.micro.srv.example

go.micro.srv.example-fccbb6fb-0301-11e5-9f1f-68a86d0d36b6	[::]	62421
```

### Call Service

```shell
micro call go.micro.srv.example Example.Call '{"name": "John"}'
```

Output
```
{
	"msg": "go.micro.srv.example-fccbb6fb-0301-11e5-9f1f-68a86d0d36b6: Hello John"
}
```

### Service Health

```shell
micro health go.micro.srv.example
```

Output

```
node		address:port		status
go.micro.srv.example-fccbb6fb-0301-11e5-9f1f-68a86d0d36b6		[::]:62421		ok
```

### Register/Deregister

```shell
micro register service '{"name": "foo", "version": "bar", "nodes": [{"id": "foo-1", "address": "127.0.0.1", "port": 8080}]}'
```

```shell
micro deregister service '{"name": "foo", "version": "bar", "nodes": [{"id": "foo-1", "address": "127.0.0.1", "port": 8080}]}'
```

## Proxy Remote Env

Proxy remote environments using the `micro proxy`

When developing against remote environments you may not have direct access to service discovery 
which makes it difficult to use the CLI. The `micro proxy` provides a http proxy for such scenarios.

Run the proxy in your remote environment

```
micro proxy
```

Set the env var `MICRO_PROXY_ADDRESS` so the cli knows to use the proxy

```shell
MICRO_PROXY_ADDRESS=staging.micro.mu:8081 micro list services
```

## Usage

```shell
NAME:
   micro - A cloud-native toolkit

USAGE:
   micro [global options] command [command options] [arguments...]
   
VERSION:
   0.8.0
   
COMMANDS:
    api		Run the micro API
    bot		Run the micro bot
    registry	Query registry
    call	Call a service or function
    query	Deprecated: Use call instead
    stream	Create a service or function stream
    health	Query the health of a service
    stats	Query the stats of a service
    list	List items in registry
    register	Register an item in the registry
    deregister	Deregister an item in the registry
    get		Get item from registry
    proxy	Run the micro proxy
    new		Create a new micro service by specifying a directory path relative to your $GOPATH
    web		Run the micro web app
```

