---
title: Web Dashboard
keywords: web, dashboard
tags: [web, dashboard]
sidebar: home_sidebar
permalink: /web
summary: Micro Web provides a dashboard to visualise and explore services
---

The web dashboard provides a visual tool for explorings services and a built-in web proxy for 
web based micro services.

## Usage

```bash
micro web
```
Browse to localhost:8082

## Use ACME

The micro web dashboard supports ACME via Let's Encrypt. It automatically gets a TLS certificate for your domain.

```
micro --enable_acme web
```

Optionally specify a host whitelist

```
micro --enable_acme --acme_hosts=example.com,api.example.com web
```

## Set TLS Certificate

The dashboard supports serving securely with TLS certificates

```bash
micro --enable_tls --tls_cert_file=/path/to/cert --tls_key_file=/path/to/key web
```

## Web Services

The web dashboard has a built-in proxy for web services. This is the idea of building web applications 
as micro services which you can do via the [go-micro/web](https://pkg.go.dev/github.com/micro/go-micro/v2/web) package.

### Routing

Web services are much like API services in the sense that they are namespaced. The default namespace is "go.micro.web".

When a request such as `/foo` hits the web proxy, it will route to the service `go.micro.web.foo`. This is what 
your service should be called; `namespace + path`.

## Screenshots

<img src="images/web1.png">

<img src="images/web2.png">

<img src="images/web3.png">

<img src="images/web4.png">


