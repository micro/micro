# Runtime

The `micro runtime` is a service runtime used to manage the lifecycle of services regardless of the underlying platform.

## Overview

The micro runtime is a service lifecycle manager which acts as an abstraction over existing process and container managers. 
It's goal is to provide a programmable layer through which services can be deployed and managed without needing to 
care about the underlying deployment mechanism.

## Design

The runtime consists of two parts:

- [go-micro/runtime](https://github.com/micro/go-micro/tree/master/runtime) - a language agnostic runtime interface
- [micro/runtime](https://github.com/micro/micro/tree/master/runtime) - an rpc service which overlays the interface

Using this and the command line we have the ability to run services literally anywhere with the same local developer experience.

The micro runtime stores services to run using the Store interface. This is pluggable so we can choose how to back this storage 
whether locally or globally. It then reconciles what's running using the go-micro/runtime interface, running anything that 
doesn't exist and stopping anything that shouldn't exist. That's it!

The runtime additionally includes the ability to inject env vars and define default git sources or docker images.

## CLI

The cli experience consists of

```
# start the service foobar
micro run foobar

# stop the service
micro kill foobar

# get service status
micro get status foobar

# Add version
micro run foobar latest
```

## Config

The runtime can be configured via env vars and flags

- MICRO_RUNTIME - sets the underlying runtime; local, kubernetes, service
- MICRO_RUNTIME_PROFILE - preset env vars injected as configuration; local, kubernetes, platform
- MICRO_RUNTIME_SOURCE - set a specific source to use whether its a git repo or docker image

## Cells

In the future we want to enable the ability to run multi-lang. This will likely involve defining something 
similar to build packs which we'll call `cells`. Cells are basically an encapsulation mechanism 
that include dependencies depending on a service language and then pull that locally or pack it into 
an image. 

- Go Micro - we can pull go from source or via the micro/go-micro docker image
- Go - we'll use a Go image and again can run from source
- Node.js - prebuild a micro/node image
- ...

## Platforms

The runtime NEEDS to operate across any environment and the experience should be consistent. We should be 
able to run it as zero dep by default and issue a `micro run helloworld` command without needing to change 
this for any environment. The micro runtime service itself takes care of the details. 

When Source is specified we know to pull helloworld from a specific source. Whether its github or dockerhub.

Platforms we currently support

- Local
- Kubernetes
- WASM (TODO)
