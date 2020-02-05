# Micro Server

This is `microd` a single process daemon which encapsulates all the features of micro in a http and gRPC interface.

## Overview

Caching has redis, search has elasticsearch, databases have postgres yet there's something missing for microservices. 
Micro and the go-micro framework has solved a number of problems for microservices development but we must now evolve 
too. The goal of `microd` is to act as a drop in solution for microservices development which a standard http 
and gRPC api along with code generated libraries and a slim client layer on top.

## Design

Every go-micro interface will be provided as part of a http and grpc interface e.g.

```
Broker => /broker or /micro.Broker/...
Registry => /registry or /micro.Registry/...
...
```

The http port will be 8080 and gRPC port 8081.

## Clients

To start with:

- Go
- Ruby
- Typescript

Code generated libraries will live in [github.com/micro/clients](https://github.com/micro/clients)

## Release

Coming soon..

## Usage

Coming soon..
