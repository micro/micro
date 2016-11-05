# Micro Sidecar

The sidecar is a language agnostic RPC proxy which provides HTTP endpoints to integrate any application into the Micro ecosystem.

It is similar to Netflix's sidecar [Prana](https://github.com/Netflix/Prana) or Buoyant's RPC Proxy [Linkerd](https://linkerd.io).

<p align="center">
  <img src="architecture.png" />
</p>

## Features

The sidecar has all the features of [go-micro](https://github.com/micro/go-micro). Here are the most relevant.

- Service registration and discovery
- Broker PubSub via WebSockets
- Healthchecking of services
- JSON RPC via HTTP API.
- Load balancing, retries, timeouts
- Stats UI

## Getting Started

### Install

```shell
go get github.com/micro/micro
```

### Run

The micro sidecar runs on port 8081 by default.

Starting the sidecar 

```shell
micro sidecar
```

Optionally specify app server name and address if you want to auto register an app on startup.

```shell
micro sidecar --server_name=foo --server_address=127.0.0.1:9090
```

### Serve Secure TLS

The Sidecar supports serving securely with TLS certificates

```bash
micro --enable_tls --tls_cert_file=/path/to/cert --tls_key_file=/path/to/key sidecar
```

### Host Discovery

```shell
curl http://127.0.0.1:8081/registry?service=go.micro.srv.example
{
	"name":"go.micro.srv.example",
	"nodes":[{
		"id":"go.micro.srv.example-c5718d29-da2a-11e4-be11-68a86d0d36b6",
		"address":"[::]","port":60728
	}]
}
```

### Register/Deregister a service
Register
```shell
// specify ttl as a param to expire the registration
// units ns|us|ms|s|m|h
// http://127.0.0.1:8081/registry?ttl=10s

curl -H 'Content-Type: application/json' http://127.0.0.1:8081/registry -d 
{
	"Name": "foo.bar",
	"Nodes": [{
		"Port": 9091,
		"Address": "127.0.0.1",
		"Id": "foo.bar-017da09a-734f-11e5-8136-68a86d0d36b6"
	}]
}
```

Deregister
```shell
curl -X "DELETE" -H 'Content-Type: application/json' http://127.0.0.1:8081/registry -d 
{
	"Name": "foo.bar",
	"Nodes": [{
		"Port": 9091,
		"Address": "127.0.0.1",
		"Id": "foo.bar-017da09a-734f-11e5-8136-68a86d0d36b6"
	}]
}
```

### Healthchecking

Start micro sidecar with "--healthcheck_url=" to enable the healthchecker

```shell
$ micro sidecar --server_name=foo --server_address=127.0.0.1:9090 \
	--healthcheck_url=http://127.0.0.1:9090/_status/health
I0523 12:25:36.229536   85658 car.go:184] Registering foo-6ebf29c0-013e-11e5-b55f-68a86d0d36b6
I0523 12:25:36.241680   85658 car.go:188] Starting sidecar healthchecker
```

### HTTP RPC API

Query micro services via the http rpc api.

```shell
$ curl  -d 'service=go.micro.srv.example' \
	-d 'method=Example.Call' \
	-d 'request={"name": "John"}' http://127.0.0.1:8081/rpc
{"msg":"go.micro.srv.example-c5718d29-da2a-11e4-be11-68a86d0d36b6: Hello John"}
```

### PubSub via WebSockets

Connect to the micro pub/sub broker via a websocket interface

```go
c, _, _ := websocket.DefaultDialer.Dial("ws://127.0.0.1:8081/broker?topic=foo", make(http.Header))

// optionally specify "queue=[queue name]" param to distribute traffic amongst subscribers
// websocket.DefaultDialer.Dial("ws://127.0.0.1:8081/broker?topic=foo&queue=group-1", make(http.Header))

go func() {
	for {
		_, p, err := c.ReadMessage()
		if err != nil {
			return
		}
		var msg *broker.Message
		json.Unmarshal(p, &msg)
		fmt.Println(msg.Data)
	}
}()

ticker := time.NewTicker(time.Second)

for _ = range ticker.C {
	if err := c.WriteMessage(1, []byte(`hello world`)); err != nil {
		return
	}
}
```

### Proxy CLI requests

The sidecar also acts as a proxy for the CLI

```shell
$ micro --proxy_address=127.0.0.1:8081 list services
go.micro.srv.greeter
```

## Stats Dashboard

You can enable a stats dashboard via the `--enable_stats` flag. It will be exposed on /stats.

```shell
micro --enable_stats sidecar
```

<img src="https://github.com/micro/micro/blob/master/doc/stats.png">
