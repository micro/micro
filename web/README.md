# Micro Web

Micro web provides a visual point of entry for the micro environment and should replicate 
the features of the CLI.

It also includes a reverse proxy to route requests to micro web 
apps. /[name] will proxy to the service [namespace].[name]. The default namespace is 
go.micro.web.

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

### Serve Secure TLS

The Web proxy supports serving securely with TLS certificates

```bash
micro --enable_tls --tls_cert_file=/path/to/cert --tls_key_file=/path/to/key web
```

### Set Namespace

The Web defaults to serving the namespace **go.micro.web**. The combination of namespace and request path 
are used to resolve a service to reverse proxy for.

```bash
micro --web_namespace=com.example.web
```

## Screenshots

<img src="https://github.com/micro/micro/blob/master/web/web1.png">
-
<img src="https://github.com/micro/micro/blob/master/web/web2.png">
-
<img src="https://github.com/micro/micro/blob/master/web/web3.png">

