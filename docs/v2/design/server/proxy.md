# Proxy

The micro proxy is a microservice proxy which encapsulates the features of go-micro as a standalone proxy which can be used 
to offload the distributed systems aspects of services. 

## Overview

The proxy implements the server `router.ServeRequest` method so it can be replaced in a go-micro server. It can 
then be used by setting the env var `MICRO_PROXY` to the proxy service name or `MICRO_PROXY_ADDRESS` to the proxy 
address. Any requests are then routed through the proxy rather than using the selector/registry.

The proxy makes use of the go-micro router which maintains routing information based on the registry and acts 
as an aggregate for intelligent routing decisions. The proxy can use this directly or talk to a router service.


