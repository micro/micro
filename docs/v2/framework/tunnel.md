---
title: Service Tunnel
keywords: tunnel
tags: [tunnel]
sidebar: home_sidebar
permalink: /tunnel
summary: The micro tunnel is a point to point tunnel.
---

A service tunnel is a point to point tunnel used for accessing services in remote environments.

<img src="images/tunnel.svg" />

## Overview

The micro tunnel provides a way to access services across remote environments. This is great where you want to tunnel to staging, prod 
or expose local services externally without using something like openvpn or wireguard which would expose all things in your network.

## Run Tunnel

Start the tunnel server (Runs on port :8083)

```shell
micro tunnel
```

## Tunnel Services

Now the tunnel is running you can connect to it with a local tunnel.

```
micro tunnel --server=remote.env:8083
```

Any request now made through the tunnel will be proxied to a service on the other side.

Set your proxy to use the tunnel
```
MICRO_PROXY=go.micro.tunnel go run main.go
```

Your service will direct all traffic through the tunnel. 


## Authentication

Specify a tunnel token to limit access to who can tunnel into the environment. Tokens must match between 
tunnel clients and servers otherwise the connection is rejected.

```
MICRO_TUNNEL_TOKEN=foobar go run main.go
```

By default the token "micro" is used allowing anyone to connect via the tunnel.

