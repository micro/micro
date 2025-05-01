# Hello World

This is hello world using go-micro

## Contents

- main.go - is the main definition of the service, handler and client
- proto - contains the protobuf definition of the API

## Dependencies

- [generator](https://github.com/go-micro/generator)

## Usage

To run it

```
go run main.go
```

To rebuild the proto

```
protoc --proto_path=. --micro_out=. --go_out=. proto/greeter.proto
```
