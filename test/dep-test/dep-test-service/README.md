# Dep Service

This is the Dep service

Generated with

```
micro new --namespace=go.micro --type=service dep-test-service
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: go.micro.service.dep
- Type: service
- Alias: dep

## Dependencies

Micro services depend on service discovery. The default is multicast DNS, a zeroconf system.

In the event you need a resilient multi-host setup we recommend etcd.

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

Run the service
```
./dep-service
```

Build a docker image
```
make docker
```

## Quirks of the tests

There is one huge gotcha that must be kept in mind.
The `t.Fatal` and other calls happen to a wrapper `t` type and not the normal `testing.T`.

This will not immediately terminate the test - it will just run through the test quickly without waiting for `try` calls. This behaviour is not intuitive in cases when the user expects the `t.Fatal` call to terminate the test and this should be fixed (perhaps just by calling return with a Fatal but that might be too easy to miss.).