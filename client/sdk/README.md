# Micro SDKs

Micro SDKs are gRPC generated clients

## Overview

Rather than handcrafting SDKS, we simply rely on gRPC for generated clients which can be used in multiple languages. 
The SDKs can be used against the proxy on localhost:8081 or proxy.m3o.com:443 on the platform. In the event services 
require authentication the "Authorization: Bearer" header should be set with a user or service token.
