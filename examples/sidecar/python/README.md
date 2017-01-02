# Python

- sidecar.py: methods to call sidecar
- rpc_{client,server}.py: RPC client/server
- http_{client,server}.py: HTTP client/server

## RPC Example

Run sidecar
```shell
micro sidecar
```

Run server
```shell
# serves Say.Hello
python rpc_server.py
```

Run client
```shell
# calls go.micro.srv.greeter Say.Hello
python rpc_client.py
```

## HTTP Example

Run sidecar with proxy handler
```shell
micro sidecar --handler=proxy
```

Run server
```shell
# serves /greeter
python http_server.py
```

Run client
```shell
# calls /greeter
python http_client.py
```
