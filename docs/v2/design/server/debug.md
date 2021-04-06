# Debugging

Debugging is one of the core requirements for any programming tools and absolutely critical for microservices.

## Overview

Micro provides a built in debugging experience through stats, logs and tracing. Debugging is a separate concern 
from monitoring. Debugging is a form of observability that's deeply integrated into the go-micro framework. 
The idea here is to mimic Go's runtime tooling or thereabouts e.g runtime stats, stdout and stderr logs, debug 
stack traces. Monitoring is built ON debugging.

Core concerns:

- Stats - runtime stats including cpu, mem, go routines, request rate, error rate
- Logs - the output of stdout and stderr
- Trace - The instrumented stack for request, response and messages

## Architecture

We include a go-micro/debug package which adds the above concerns as a core tenant of the framework and platform. 
Debug is then embedded as a handler into every service with endpoints `Debug.{Log, Stats, Trace}` and may even 
be extended to Status, Info, etc. Our goal will be to mimic programming tools when it comes to Debug and not 
overly extend of complicate. Simplicity is key.

Every service maintains an memory buffers for stats, logs and tracing. This is so any service can be directly 
queried to retrieve that information. This is a zero dep experience which works locally and in production.

This data is then scraped by the go.micro.debug service which is part of our micro runtime. This provides an 
aggregated view of all services.

## Centralisation

Because debug concerns are interfaces like any other, they will be implemented and backed by centralised systems 
such as netdata for stats, kubernetes for logs and jaeger for tracing. While we can write our own storage 
and implementations for these, the ability to offload to systems in the platform is useful.
