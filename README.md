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

Create your service using [Go Micro](https://github.com/micro/go-micro)

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

Run a service

```
micro run
```

List services

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

        req := client.NewRequest("helloworld", "Say.Hello", &Request{Name: "John"})

        var rsp Response

        err := client.Call(context.TODO(), req, &rsp)
        if err != nil {
                fmt.Println(err)
                return
        }

        fmt.Println(rsp.Message)
}
```

## Make a HTTP call

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
curl http://localhost:8080/helloworld/Say/Hello -d '{"name": "John"}'
```
