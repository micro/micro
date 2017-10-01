# micro web

The **micro web** provides a dashboard to view and query services as well as a reverse proxy to serve micro web applications. 
We believe in web apps as first class citizens in a microservice world.

<p align="center">
  <img src="https://github.com/micro/docs/blob/master/images/web.png" />
</p>

## API

```
- / (UI)
- /[service]
- /rpc
```

## Features

Feature	|	Description
---	|	---
UI	|	A dashboard to view and query running services
Proxy	|	A reverse proxy to micro web services (includes websocket support)

### Proxy

Micro Web has a built in HTTP reverse proxy for micro web apps. This essentially allows you 
to treat web applications as first class citizens in a microservices environment. The proxy 
will use /[service] along with the namespace (default: go.micro.web) to lookup the service 
in service discovery. It composes service name as [namespace].[name]. 

The proxy will strip /[service] from the request and forward the rest of the URL Path to the web app. It will also 
set the header "X-Micro-Web-Base-Path" to the removed path incase you need to use it for 
some reason like constructing URLs.

Example translation

Path	|	Service	|	Service Path	|	Header: X-Micro-Web-Base-Path
---	|	---	|	---	|	---
/foo	|	go.micro.web.foo	|	/	|	/foo
/foo/bar	|	go.micro.web.foo	|	/bar	|	/foo

Note: The web proxy speaks to services using HTTP. There is no ability to switch out transport.

## Getting Started

### Install
```bash
go get github.com/micro/micro
```

### Run Web UI/Proxy

```bash
micro web
```
Browse to localhost:8082

### ACME via Let's Encrypt

Serve securely by default using ACME via letsencrypt 

```
micro --enable_acme web
```

Optionally specify a host whitelist

```
micro --enable_acme --acme_hosts=example.com,api.example.com web
```

### Serve Secure TLS

The Web proxy supports serving securely with TLS certificates

```bash
micro --enable_tls --tls_cert_file=/path/to/cert --tls_key_file=/path/to/key web
```

### Set Namespace

The Web defaults to serving the namespace **go.micro.web**. The combination of namespace and request path 
are used to resolve a service to reverse proxy for.

```bash
micro web --namespace=com.example.web
```

## Stats

You can enable a stats dashboard via the `--enable_stats` flag. It will be exposed on /stats.

```shell
micro --enable_stats web
```

<img src="https://github.com/micro/docs/blob/master/images/stats.png">

## Screenshots

<img src="https://github.com/micro/docs/blob/master/images/web1.png">
-
<img src="https://github.com/micro/docs/blob/master/images/web2.png">
-
<img src="https://github.com/micro/docs/blob/master/images/web3.png">
-
<img src="https://github.com/micro/docs/blob/master/images/web4.png">

