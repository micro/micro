# Ruby

- sidecar.rb: methods to call sidecar
- rpc_{client,server}.rb: RPC client/server
- http_{client,server}.rb: HTTP client/server

## RPC Example

Run sidecar
```shell
micro sidecar
```

Run server
```shell
# serves Say.Hello
ruby rpc_server.rb
```

Run client
```shell
# calls go.micro.srv.greeter Say.Hello
ruby rpc_client.rb
```

## HTTP Example

Run sidecar with proxy handler
```shell
micro sidecar --handler=proxy
```

Run server
```shell
# serves /greeter
ruby http_server.rb
```

Run client
```shell
# calls /greeter
ruby http_client.rb
```
