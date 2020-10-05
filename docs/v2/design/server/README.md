# Server

The micro server is a standalone server that encapsulates all the features of micro in a single process.

## Overview

We've moved from Go framework to an RPC runtime but the complexity of our system makes the orchestration 
of individual services difficult. In a cloud based environment this is ok but locally or in smaller 
systems its an overhead we don't want to incur.

The micro server is a lightweight standalone system for microservices.

## Previous

We've previously made the entire stack available through a one command using `micro`. This boots the 
entire micro stack of services as processes. Our aim is still to have this as a point of coordination 
but individual processes do not serve well on machines with limited resources. 

## Solution

Our solution would be started using `micro server` and clustered by specifying a list of other 
nodes to connect to with `--nodes`. 

The server would be accessible on :8080 for http and :8081 for grpc. All calls would be routed 
to internal systems but also may act as a pass through for other services.

We'd expect services to conform to some sort of namespacing such as go.micro.service.foo where 
we'd enable calling /service/foo and combine that with the namespace. This may be an optional 
thing. What this ultimately enables is:

- /api/ routes to go.micro.api.x
- /web/ routes to go.micro.web.x
- /service/ routes to go.micro.service.x

## Deployments

Currently there is some overloading of concepts with `micro run` and the runtime service.
The syntax of the `micro run` command is `micro run [service] [version] [--source github.com/micro/services]`.

We propose to change this to `micro run [source]` for the following reasons:
- The name of the service is in the service source code.
- The version can be extracted from git.

The only thing we require is the source, which can be either:

- Local path to source
- Git url
- Relative path to the `github.com/micro/services` repo (ie. `helloworld/web`)

The runtime service can handle the extraction of the service name and version.

