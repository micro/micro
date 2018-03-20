# micro web

The **micro web** provides a dashboard to view and query services.

## Getting Started

- [Install](#install)
- [Run](#run)
- [ACME](#acme)
- [TLS Cert](#tls-cert)
- [Screenshots](#screenshots)

### Install

```bash
go get github.com/micro/micro
```

### Run

```bash
micro web
```
Browse to localhost:8082

### ACME

Serve securely by default using ACME via letsencrypt 

```
micro --enable_acme web
```

Optionally specify a host whitelist

```
micro --enable_acme --acme_hosts=example.com,api.example.com web
```

### TLS Cert

The Web proxy supports serving securely with TLS certificates

```bash
micro --enable_tls --tls_cert_file=/path/to/cert --tls_key_file=/path/to/key web
```

## Screenshots

<img src="https://github.com/micro/docs/blob/master/images/web1.png">
-
<img src="https://github.com/micro/docs/blob/master/images/web2.png">
-
<img src="https://github.com/micro/docs/blob/master/images/web3.png">
-
<img src="https://github.com/micro/docs/blob/master/images/web4.png">

