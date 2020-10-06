# API Gateway

API Gateway is useful, when you have many services and want to have single point where you authentificate, athorize users.
And pass request for specific endpoint to specific micro service via rpc call or pubsub system.

## Overview

API Gateway must provide ability not to register some handler interface to specific endpoint via some protocol like rest or grpc, graphql etc for http, unix or udp for plain sockets etc.
But as services can up and down registration must be done not once at start, but when some service want to expose it via api.

## Implemenations

Now we have api handlers:
* api
* broker
* cloudevents
* event
* file
* http
* registry
* rpc
* udp
* unix
* web

## Design

At this moment we don't have ability to register handler after api service starts.

## Proposed changes

Create proto service definition and implemets it. Needs ivestigation how the best do that.

