# Micro

A Go microservices platform

## Overview

Micro is a platform for Go microservices development. It provides the tools required for building services in the cloud. 
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

Create your service using [Go Micro](https://go-micro.dev)

```go
package main

import (
        "go-micro.dev/v5"
)

type Request struct {
        Name string `json:"name"`
}

type Response struct {
        Message string `json:"message"`
}

type Say struct{}

func (h *Say) Hello(ctx context.Context, req *Request, rsp *Response) error {
        rsp.Message = "Hello " + req.Name
        return nil
}

func main() {
        // create the service
        service := micro.New("helloworld")

        // register handler
        service.Handle(new(Say))

        // run the service
        service.Run()
}
```

Run your service

```
micro run
```

List your services

```
micro services
```

Call a service

```
micro call helloworld Say.Hello '{"name": "Asim"}'
```

Describe a service

```
micro describe helloworld
```

Output

```
{
    "name": "helloworld",
    "version": "latest",
    "metadata": null,
    "endpoints": [
        {
            "request": {
                "name": "Request",
                "type": "Request",
                "values": [
                    {
                        "name": "string",
                        "type": "string",
                        "values": null
                    }
                ]
            },
            "response": {
                "name": "Response",
                "type": "Response",
                "values": [
                    {
                        "name": "string",
                        "type": "string",
                        "values": null
                    }
                ]
            },
            "metadata": {},
            "name": "Say.Hello"
        }
    ],
    "nodes": [
        {
            "metadata": {
                "broker": "http",
                "protocol": "mucp",
                "registry": "mdns",
                "server": "mucp",
                "transport": "http"
            },
            "id": "helloworld-9988def2-2ee4-45f1-9cf7-faa62535538f",
            "address": "172.17.0.1:40397"
        }
    ]
}
```

## API

Run the API

```
micro api
```

If you have [helloworld](https://github.com/micro/helloworld) running

```
curl http://localhost:8080/helloworld/Say/Hello -d '{"name": "John"}'
```

Or with headers

```
curl -H 'Micro-Service: helloworld' -H 'Micro-Endpoint: Say.Hello' http://localhost:8080/ -d '{"name": "John"}'
```

## MCP

Experimental support for [MCP](https://github.com/modelcontextprotocol)

```
micro mcp
```

Supports the tools `services`, `call`, `describe`.

## Network

To enable tailscale via the api

```
micro api --network=tailscale
```

Ensure to export your `TS_AUTHKEY`

## Plugins

Plugins can be found in [micro/plugins](https://github.com/micro/plugins) which enable various underlying interface implementations.

Note: This requires a rebuild of the binary to include those plugins. There is no clean approach to this yet.
