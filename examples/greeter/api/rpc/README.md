# RPC API Example

This is an example of using an RPC based API when the api handler is set with `--handler=rpc`

## Getting Started

### Run the Micro API

```
$ micro api --handler=rpc
```

### Run the Greeter Service

```
$ go run greeter/server/main.go
```

### Run the Greeter API

```
$ go run rpc.go
Listening on [::]:64738
```

### Curl the API

Test the index
```
curl -H 'Content-Type: application/json' -d '{"name": "Asim"}' http://localhost:8080/greeter/hello
{
  "msg": "Hello Asim"
}
```
