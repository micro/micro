# Sidecar Examples

The [micro sidecar](https://github.com/micro/micro/tree/master/car) is a language agnostic proxy which provides all the features 
of [go-micro](https://github.com/micro/go-micro) as a HTTP server. To learn more about the sidecar look [here](https://github.com/micro/micro/tree/master/car).

This directory contains examples for using the sidecar with various languages.

## Usage

Run sidecar
```
micro sidecar
```

Run sidecar with http proxy handler
```
micro sidecar --handler=proxy
```

## Examples

Each language directory {python, ruby, ...} contains examples for the following:

- Registering Service
- JSON RPC Server and Client
- HTTP Server and Client
