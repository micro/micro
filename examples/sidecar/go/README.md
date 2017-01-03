# Go

This Go example uses vanilla net/http and the sidecar

- sidecar.go: methods to call sidecar
- rpc_{client,server}.go: RPC client/server
- http_{client,server}.go: HTTP client/server

## RPC Example

Run sidecar
```shell
micro sidecar
```

Run server
```shell
# serves Say.Hello
go run rpc_server.go sidecar.go
```

Run client
```shell
# calls go.micro.srv.greeter Say.Hello
go run rpc_client.go sidecar.go
```

## HTTP Example

Run sidecar with proxy handler
```shell
micro sidecar --handler=proxy
```

Run server
```shell
# serves /greeter
go run http_server.go sidecar.go
```

Run client
```shell
# calls /greeter
go run http_client.go sidecar.go
```
