---
title: Debugging
keywords: debug
tags: [debug]
sidebar: home_sidebar
permalink: /debugging
summary: 
---

Enable debugging of micro or go-micro very simply via the following environment variables.

## Logging

To enable debug logging

```
MICRO_LOG_LEVEL=debug
```

The log levels supported are

```
trace
debug
error
info
```

To view logs from a service

```
micro log [service]
```

## Profiling

To enable profiling via pprof

```
MICRO_DEBUG_PROFILE=http
```

This will start a http server on :6060

## Stats

To view the current runtime stats

```
micro stats [service]
```

## Health

To see if a service is running and responding to RPC queries

```
micro health [service]
```

