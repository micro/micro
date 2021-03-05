# OpenAPI

This directory contains OpenAPI clients for the micro api.

## Overview

The micro api is an api gateway which converts http/json to rpc. This means we can serve standard http/json to the frontend 
while using gRPC and protobuf on the backend. OpenAPI (previously Swagger) is an open standard for API definitions and usage. 
We provide specs based on the micro proto definitions.

Specs
=============

Services
--------
- [x] Auth
- [x] Broker
- [x] Config
- [x] Events
- [ ] Registry (currently having recursive looping issues)
- [x] Runtime
- [x] Signup
- [x] Store
