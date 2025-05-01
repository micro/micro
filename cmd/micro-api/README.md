# Micro API

An API server for Go Micro

## Overview

The micro API provides a fixed HTTP entrypoint for Go Micro services. Where services may run on random ports or 
different transports, the micro API runs a http server that can serve standard http/json requests and route via 
RPC to the correct service.

## Path Based Resolution

The micro API uses path based resolution to determine request to make

Given a service named `helloworld` with endpoint `Say.Hello` the following HTTP path will work

```
/helloworld/Say/Hello
```

## Usage

Install the api

```
go get github.com/micro/micro/cmd/micro-api@latest
```

Run the api

```
micro-api
```

Call a service

```
curl http://localhost:8080/helloworld/Say/Hello -d '{"name": "John"}'
```
