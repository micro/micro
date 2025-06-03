# Micro [![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![Discord](https://img.shields.io/badge/discord-chat-800080?style=flat-square)](https://discord.gg/UmFkPbu32m)

A Go microservices toolkit

## Overview

Micro is an ecosystem for Go microservices development. It provides the tools required for building services in the cloud. 
The core of Micro is the [Go Micro](https://github.com/micro/go-micro) framework, which developers import and use in their code to 
write services. Surrounding this we introduce a number of tools to make it easy to serve and consume services. 

## Install the CLI

Install `micro` via `go get`

```
go get github.com/micro/micro/v5@latest
```

Or via install script

```
wget -q  https://raw.githubusercontent.com/micro/micro/master/scripts/install.sh -O - | /bin/bash
```

For releases see the [latest](https://github.com/micro/micro/releases/latest) tag

## Create a service

Create your service (all setup is now automatic!):

```
micro new helloworld
```

This will:
- Create a new service in the `helloworld` directory
- Automatically run `go mod tidy` and `make proto` for you
- Show the updated project tree including generated files
- Warn you if `protoc` is not installed, with install instructions

If you need OpenAPI support, install `protoc-gen-openapi` separately:
```
go install github.com/google/gnostic/plugins/protoc-gen-openapi@latest
```

## Run the service

Run the service

```
micro run .
```

List services to see it's running and registered itself

```
micro services
```

## Describe the service

Describe the service to see available endpoints

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
                        "name": "name",
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
                        "name": "msg",
                        "type": "string",
                        "values": null
                    }
                ]
            },
            "metadata": {},
            "name": "Helloworld.Call"
        },
        {
            "request": {
                "name": "Context",
                "type": "Context",
                "values": null
            },
            "response": {
                "name": "Stream",
                "type": "Stream",
                "values": null
            },
            "metadata": {
                "stream": "true"
            },
            "name": "Helloworld.Stream"
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
            "id": "helloworld-31e55be7-ac83-4810-89c8-a6192fb3ae83",
            "address": "127.0.0.1:39963"
        }
    ]
}
```

## Call the service

Call the service using the CLI

```
micro call helloworld Helloworld.Call '{"name": "Asim"}'
```

## Create a client

Create a client to call the service

```go
package main

import (
        "context"
        "fmt"

        "go-micro.dev/v5"
)

type Request struct {
        Name string
}

type Response struct {
        Message string
}

func main() {
        client := micro.New("helloworld").Client()

        req := client.NewRequest("helloworld", "Helloworld.Call", &Request{Name: "John"})

        var rsp Response

        err := client.Call(context.TODO(), req, &rsp)
        if err != nil {
                fmt.Println(err)
                return
        }

        fmt.Println(rsp.Message)
}
```

## API Endpoints 

The API provides a fixed HTTP entrypoint for calling services.

```
curl http://localhost:8080/helloworld/Helloworld/Call -d '{"name": "John"}'
```
See /api for more details

## Web Dashboard 

Access services via the web using micro web which generates dynamic form fills

Go to [localhost:8080](http://localhost:8080)
