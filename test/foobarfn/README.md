# Foobarfn Function

This is the Foobarfn function

Generated with

```
micro new --namespace=go.micro --type=function foobarfn
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: go.micro.function.foobarfn
- Type: function
- Alias: foobarfn

## Dependencies

Micro functions depend on service discovery. The default is etcd.

```
# install etcd
brew install etcd

# run etcd
etcd
```

## Usage

A Makefile is included for convenience

Build the binary

```
make build
```

Run the function once
```
./foobarfn-function
```

Build a docker image
```
make docker
```