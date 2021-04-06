---
title: V2 to V3 Upgrade Guide
keywords: micro
tags: [micro]
sidebar: home_sidebar
permalink: /v2-to-v3-upgrade-guide
summary: Upgrading your project from Go-Micro v2 to Micro v3
---

## Upgrading your project from Go-Micro v2 to Micro v3

### Introduction

As part of v3 we've simplified the development model and consolidated the micro and go-micro libraries. This means that you, as the developer, only need to worry about leveraging the libraries offered in micro. Users of the now deprecated [go-micro](https://github.com/asim/go-micro) framework can follow this guide to upgrade their projects to use Micro v3.

### Service

In go-micro, you would create a service by importing `github.com/micro/go-micro` and using the following syntax:

```go
  srv := micro.NewService(
    micro.Name("go.micro.service.foo")
  )
```

In Micro V3, we've removed namespaces from service names, so the above service would now be named simply **foo**. We've also created a package specifically for services: `github.com/micro/micro/v3/service`. Dotted service names can still be used, however the full service name will be required when calling the API, such as: "http://localhost:8080/go.micro.service.foo/Bar", unless the API is configured with the namespace flag.

The above service can now be defined as:
```go
  srv := service.New(
    service.Name("foo")
  )
```

### Packages

In go-micro, modules such as config could only be accessed via the service object, for example `srv.Options().Config.Read()`. In Micro this has been made easier by providing public functions in each package which can be accessed by importing any package, nested within "github.com/micro/micro/v3/service". Here is an example of how you can load config in Micro:

```go
package main

import (
  "github.com/micro/micro/v3/service"
  "github.com/micro/micro/v3/service/config"
  "github.com/micro/micro/v3/service/logger"
)

func main() {
  // Create the service
  srv := service.New(
    service.Name("config-example"),
  )

  // Load the config
  val, err := config.Get("mykey")
  if err != nil {
    logger.Fatalf("Could not load mykey: %v", err)
  } else {
    logger.Infof("mykey = %v", val)
  }

  // Run the service
  if err := srv.Run(); err != nil {
    logger.Fatal(err)
  }
}

```

Other examples of packages which can be imported in this manner are:
* [Auth](https://github.com/micro/micro/tree/master/service/auth)
* [Broker](https://github.com/micro/micro/tree/master/service/broker)  
* [Client](https://github.com/micro/micro/tree/master/service/client)
* [Config](https://github.com/micro/micro/tree/master/service/config)
* [Context](https://github.com/micro/micro/tree/master/service/context)
* [Errors](https://github.com/micro/micro/tree/master/service/errors)
* [Logger](https://github.com/micro/micro/tree/master/service/logger)
* [Store](https://github.com/micro/micro/tree/master/service/store)


### Running services

Go Micro services interfaced directly with the underlying infrastructure, meaning they could be run using `go run`. In Micro V3 this has been abstracted using the `micro server`. The Micro server provides a level of abstraction to the underlying infrastructure, keeping your services unbound from their dependancies. 

In practice, this means you need to be connected to a micro environment when running your services. You can spin up this environment locally using the `micro server` command, or connect to the free platform provided by M3O using `micro env set platform`. You can read more about running a service [in our getting started guide](/getting-started#running-a-service).

*Note: don't forget to run your micro v3 services using `micro run .` instead of `go run .`*
