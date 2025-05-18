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

Create your service and follow the guide

```
micro new helloworld
```

Tidy and make protos

```
cd helloworld
go mod tidy
make proto
```

Run the service

```
micro run .
```

List services

```
micro services
```

Call a service

```
micro call helloworld Helloworld.Call '{"name": "Asim"}'
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

## Micro API

Call services via http using the [micro-api](https://github.com/micro/micro/tree/master/cmd/micro-api)

Install the API 

```
go get github.com/micro/micro/cmd/micro-api@latest
```

Run the API

```
micro-api
```

If you have the service running

```
curl http://localhost:8080/helloworld/Helloworld/Call -d '{"name": "John"}'
```

## Micro Web

Access services via the web

## Usage

Install the web app

```
go get github.com/micro/micro/cmd/micro-web@latest
```

Run the web app

```
micro-web
```

Go to [localhost:8082](http://localhost:8082)
