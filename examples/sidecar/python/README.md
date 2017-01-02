# Python

- sidecar.py: methods to call sidecar
- rpc_{client,server}.py: RPC client/server
- http_{client,server}.py: HTTP client/server

## RPC Example

```shell
micro sidecar
```

```shell
python rpc_server.py
```

```shell
python rpc_client.py
```

## HTTP Example

```shell
micro sidecar --handler=proxy
```

```shell
python http_server.py
```

```shell
python http_client.py
```
