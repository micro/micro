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

# Micro API Admin Endpoints

The Micro API provides HTTP endpoints for direct programmatic access to go-micro platform primitives and service calls. All endpoints use consistent JSON input/output and robust error handling.

## Features

- **Store API**: CRUD operations for key-value data
  - `POST /store/write` — Write a record
  - `GET /store/read?key=...&table=...` — Read a record
  - `DELETE /store/delete` — Delete a record
  - `GET /store/list?prefix=...&table=...` — List records by prefix
- **Broker API**: Publish and subscribe to topics
  - `POST /broker/publish` — Publish a message
  - `POST /broker/subscribe` — Subscribe to a topic (one message)
- **Config API**: Manage configuration values
  - `GET /config/get?key=...` — Get a config value
  - `POST /config/set` — Set a config value
  - `DELETE /config/delete` — Delete a config value
  - `GET /config/list` — List all config keys
- **Registry API**: Service discovery and management
  - `GET /registry/list` — List all services
  - `GET /registry/get?name=...` — Get a service
  - `POST /registry/register` — Register a service
  - `POST /registry/deregister` — Deregister a service
- **Service Proxy**: Call any service/endpoint via headers or path

## Usage

Start the API server:

```sh
micro api
```

### Example: Store Write

```sh
curl -X POST http://localhost:8080/store/write -d '{"key":"foo","value":"bar"}' -H 'Content-Type: application/json'
```

### Example: Broker Publish

```sh
curl -X POST http://localhost:8080/broker/publish -d '{"topic":"demo","message":"hello"}' -H 'Content-Type: application/json'
```

### Example: Config Get

```sh
curl 'http://localhost:8080/config/get?key=foo'
```

### Example: Registry List

```sh
curl 'http://localhost:8080/registry/list'
```

## Error Handling

All endpoints return JSON errors with HTTP status codes and clear messages.

## See Also
- [Web Admin UI](../micro-web/README.md)
- [CLI Admin Commands](../micro-cli/README.md)

---

For more information, see the main [Micro documentation](../../README.md).
