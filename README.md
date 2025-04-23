# Micro

A microservices platform

## Overview

Micro is a platform for microservices development. It provides the core tools required for building services in the cloud. 
The core of Micro is the [Go Micro](https://go-micro.dev) framework, which developers import and use in their code to 
write services. Surrounding this we introduce a number of tools like a CLI and API to make it easy to serve and consume 
services. 

## Install

Micro is a single binary

```
go get github.com/micro/micro/v5@latest
```

Check the version

```
micro --version
micro version v5.0.0
```

## Usage

List your services

```
micro list services
```

Call a service

```
micro call [service] [endpoint] [request]

e.g

micro call helloworld Say.Hello '{"name": "Asim"}'
```

