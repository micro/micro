---
title: Run Service
keywords: runtime
tags: [micro, runtime]
sidebar: home_sidebar
permalink: /run-service
summary: Run a service with the micro runtime
---

The micro runtime provides a way to manage the lifecycle of services without the complexity of orchestration systems. 
By default it provides a simple local process manager which starts your service, watches for changes and rebuilds 
as required.

You can start the service with one command

```
micro run service
```

## Usage

To run and manage your service locally do the following.

```
# cd to your service directory e.g examples/greeter/srv
cd examples/greeter/srv

# run the service
micro run service

# edit a file
sed -i '1 i\// Package main' main.go
```

Watch as your service is started, sees the file change, rebuilds and starts again.

## TODO

We'll be adding the ability to send the command to run a service to the `micro runtime` service as well as 
querying status, killing the service, etc. 

