# GRPC Support

Many people want gRPC support by default in micro. We need to identify what this means and the best way to address it.

## Overview

GRPC is a generic RPC framework on top of http/2. It can be considered both a framework and a protocol. GRPC Go the framework 
tightly couples the transport protocol inside its framework. This doesn't fit inline with micro where we decouple the 
client/server from the transport. Finding a point of integration has been largely difficult. 

The industry trend is towards gRPC for cloud-native infrastructure and backend services. We've supported gRPC as client/server 
which ignores our transport but also separately with a transport with a fix RPC interface. We however do not interop with 
the gRPC framework and this is something people seem to want e.g including gRPC interceptors and supporting the gRPC 
server signature, protobuf code generation, etc.

The expectation is that because we use proto IDL code generation, protobuf and RPC that we should interop with gRPC but 
as we move away from proto this becomes less relevant. Our reasons for moving away are because we need something that 
more cleanly defines our interfaces including the pubsub interface and also defines services in a better way.

If this is the case our support of gRPC should be limited to a transport.

## Client/Server

Our [client](https://github.com/micro/go-plugins/tree/master/client/grpc) and [server](https://github.com/micro/go-plugins/tree/master/server/grpc) 
implementations can be found here respectively.

The client/server leverages the gRPC-go framework beneath the covers and ignores our own transport interface. It however does not provide 
the user access to gRPC framework related tooling as we still have our own interfaces. It also requires dual maintenance of our own 
rpc client/server and the gRPC implementations.

<img src="https://micro.mu/docs/images/go-grpc.svg" />

There is some merit to saying those who want to use the gRPC framework should use it directly.

## Transport

Our [transport](https://github.com/micro/go-plugins/tree/master/transport/grpc) implementation can be found here.

The transport has a fixed interface:

```
syntax = "proto3";

package go.micro.grpc.transport;

service Transport {
	rpc Stream(stream Message) returns (stream Message) {}
}

message Message {
	map<string, string> header = 1;
	bytes body = 2;
}
```

This enables us to leverage the grpc protocol and leverage its performance but does not cater to the needs of users who appear to want 
interop with other gRPC services and visibility into what their services are calling from the gRPC tooling standpoint. It is however 
for all intents and purposes, a transport.

## Goals

- Determine the best way to leverage the gRPC protocol
- Identify the best route forward for gRPC framework interop
- Find the simplest path forward for maintainence 
- Decide whether we should drop our own protocol for this
