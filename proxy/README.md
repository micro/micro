# micro proxy

The **micro proxy** is a cli proxy.

The micro proxy provides a http api which serves as a proxy for the cli where an environment is not directly accessible.

## Getting Started

- [Install](#install)
- [Dependencies](#dependencies)
- [Run Proxy](#run)
- [ACME](#acme)
- [Proxy CLI](#proxy-cli)

## Usage

### Install

```shell
go get -u github.com/micro/micro
```

### Dependencies

The proxy uses go-micro which means it depends on service discovery.

Install consul

```
brew install consul
consul agent -dev
```

### Run

The micro proxy runs on port 8081 by default. 

Start the proxy

```shell
micro proxy
```

### ACME

Serve securely by default using ACME via letsencrypt 

```
MICRO_ENABLE_ACME=true micro proxy
```

Optionally specify a host whitelist

```
MICRO_ENABLE_ACME=true MICRO_ACME_HOSTS=example.com,api.example.com micro proxy 
```

## Proxy CLI

To use the proxy with the CLI specify it's address

```shell
MICRO_PROXY_ADDRESS=127.0.0.1:8081 micro list services
```

```
MICRO_PROXY_ADDRESS=127.0.0.1:8081 micro call greeter Say.Hello '{"name": "john"}'
```

